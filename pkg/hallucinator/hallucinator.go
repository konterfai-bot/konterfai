package hallucinator

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/links"
	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"
	"codeberg.org/konterfai/konterfai/pkg/renderer"
)

// Hallucinator is the structure for the Hallucinator.
type Hallucinator struct {
	Interval                                time.Duration
	hallucinations                          []Hallucination
	hallucinationCount                      int
	hallucinationLock                       sync.Mutex
	hallucinationCountLock                  sync.Mutex
	hallucinationCacheSize                  int
	hallucinationRequestCount               int
	hallucinationWordCount                  int
	hallucinationLinkPercentage             int
	hallucinatorLinkMaxSubdirectories       int
	hallucinatorLinkHasVariablesProbability float64
	hallucinatorLinkMaxVariables            int
	hallucinatorUrl                         url.URL
	ollamaAddress                           string
	ollamaModel                             string
	aiTemperature                           float64
	aiSeed                                  int
	promptWordCount                         int

	httpClient *http.Client
	renderer   *renderer.Renderer
}

// NewHallucinator creates a new Hallucinator instance.
func NewHallucinator(interval time.Duration,
	hallucinationCacheSize int,
	hallucinatorPromptWordCount int,
	hallucinationRequestCount int,
	hallucinationWordCount int,
	hallucinationLinkPercentage int,
	hallucinatorLinkMaxSubdirectories int,
	hallucinatorLinkHasVariablesProbability float64,
	hallucinatorLinkMaxVariables int,
	hallucinatorUrl url.URL,
	ollamaAddress string,
	ollamaModel string,
	ollamaRequestTimeOut time.Duration,
	aiTemperature float64,
	aiSeed int,
) *Hallucinator {
	headLineLinks := [10]string{}
	for i := 0; i < len(headLineLinks); i++ {
		headLineLinks[i] = links.RandomLink(
			hallucinatorUrl,
			hallucinatorLinkMaxSubdirectories,
			hallucinatorLinkMaxVariables,
			hallucinatorLinkHasVariablesProbability,
		)
	}
	return &Hallucinator{
		Interval:                                interval,
		hallucinations:                          []Hallucination{},
		hallucinationCount:                      0,
		hallucinationCacheSize:                  hallucinationCacheSize,
		hallucinationRequestCount:               hallucinationRequestCount,
		hallucinationWordCount:                  hallucinationWordCount,
		hallucinationLinkPercentage:             hallucinationLinkPercentage,
		hallucinatorLinkMaxSubdirectories:       hallucinatorLinkMaxSubdirectories,
		hallucinatorLinkHasVariablesProbability: hallucinatorLinkHasVariablesProbability,
		hallucinatorLinkMaxVariables:            hallucinatorLinkMaxVariables,
		hallucinatorUrl:                         hallucinatorUrl,
		ollamaAddress:                           ollamaAddress,
		ollamaModel:                             ollamaModel,
		aiTemperature:                           aiTemperature,
		aiSeed:                                  aiSeed,
		promptWordCount:                         hallucinatorPromptWordCount,

		httpClient: &http.Client{
			Timeout: ollamaRequestTimeOut,
		},
		renderer: renderer.NewRenderer(headLineLinks[:]),
	}
}

// PopHallucination withdraws the first hallucination from the list of hallucinations.
func (h *Hallucinator) PopHallucination() string {
	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.cleanHallucinations()
	if len(h.hallucinations) == 0 {
		hallucination, err := h.renderer.RenderInRandomTemplate(renderer.RenderData{
			NewsAnchor:   textblocks.RandomNewsPaperName(),
			Headline:     textblocks.RandomHeadline(),
			Content:      dreamString,
			FollowUpLink: template.HTML(h.generateFollowUpLink(backToStartString)),
			RandomTopics: h.generateRandomTopicLinks(10),
			Year:         functions.PickRandomYear(),
		})
		if err != nil {
			return "Could not render...."
		}
		return hallucination
	}
	h.hallucinations[0].RequestCount--
	hallucination, err := h.renderer.RenderInRandomTemplate(renderer.RenderData{
		NewsAnchor:   textblocks.RandomNewsPaperName(),
		Headline:     textblocks.RandomHeadline(),
		Content:      template.HTML(h.clutterTextWithRandomHref(h.hallucinations[0].Text)),
		FollowUpLink: template.HTML(h.generateFollowUpLink(continueString)),
		RandomTopics: h.generateRandomTopicLinks(10),
		Year:         functions.PickRandomYear(),
	})
	if err != nil {
		return "Could not render...."
	}
	return hallucination
}

// PopRandomHallucination withdraws a random hallucination from the list of hallucinations.
func (h *Hallucinator) PopRandomHallucination() string {
	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.cleanHallucinations()
	if len(h.hallucinations) == 0 {
		hallucination, err := h.renderer.RenderInRandomTemplate(renderer.RenderData{
			NewsAnchor:   textblocks.RandomNewsPaperName(),
			Headline:     textblocks.RandomHeadline(),
			Content:      dreamString,
			FollowUpLink: template.HTML(h.generateFollowUpLink(backToStartString)),
			RandomTopics: h.generateRandomTopicLinks(10),
			Year:         functions.PickRandomYear(),
		})
		if err != nil {
			return "Could not render...."
		}
		return hallucination
	}
	randomIndex := rand.Intn(len(h.hallucinations))
	h.hallucinations[randomIndex].RequestCount--
	hallucination, err := h.renderer.RenderInRandomTemplate(renderer.RenderData{
		NewsAnchor:   textblocks.RandomNewsPaperName(),
		Headline:     textblocks.RandomHeadline(),
		Content:      template.HTML(h.clutterTextWithRandomHref(h.hallucinations[randomIndex].Text)),
		FollowUpLink: template.HTML(h.generateFollowUpLink(continueString)),
		RandomTopics: h.generateRandomTopicLinks(10),
		Year:         functions.PickRandomYear(),
	})
	if err != nil {
		return "Could not render...."
	}
	return hallucination
}

// Start starts the Hallucinator.
func (h *Hallucinator) Start(ctx context.Context) error {
	for {
		if h.GetHallucinationCount() < h.hallucinationCacheSize {
			fmt.Printf("hallucinations cache has empty slots, generating more... [%d/%d]\n", len(h.hallucinations)+1, h.hallucinationCacheSize)
			hal, err := h.generateHallucination()
			if err != nil {
				fmt.Println(fmt.Errorf("could not generate hallucination (%v)", err))
			} else {
				if h.isValidResult(hal) {
					h.appendHallucination(Hallucination{
						Text:         hal,
						RequestCount: h.hallucinationRequestCount,
					})
				} else {
					fmt.Println("invalid hallucination, skipping...")
				}
			}
		} else {
			fmt.Println("hallucinations cache is full, waiting for next interval...")
			functions.SleepWithContext(ctx, h.Interval)
		}
	}
}

// clutterTextWithRandomHref clutters the given text with random hrefs.
func (h *Hallucinator) clutterTextWithRandomHref(text string) string {
	textSlice := strings.Split(text, " ")
	percentile := len(textSlice) * h.hallucinationLinkPercentage / 100
	generated := map[int]bool{}
	for {
		if len(generated) >= percentile {
			break
		}
		i := rand.Intn(len(textSlice))
		if !generated[i] {
			textSlice[i] = fmt.Sprintf("<a href=\"%s\">%s</a>",
				links.RandomLink(
					h.hallucinatorUrl,
					h.hallucinatorLinkMaxSubdirectories,
					h.hallucinatorLinkMaxVariables,
					h.hallucinatorLinkHasVariablesProbability,
				), textSlice[i])
			generated[i] = true
		}
	}
	return strings.Join(textSlice, " ")
}
