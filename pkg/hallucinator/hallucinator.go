package hallucinator

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/links"
	"codeberg.org/konterfai/konterfai/pkg/renderer"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
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
	statistics *statistics.Statistics
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
	statistics *statistics.Statistics,
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
		renderer:   renderer.NewRenderer(headLineLinks[:]),
		statistics: statistics,
	}
}

// Start starts the Hallucinator.
func (h *Hallucinator) Start(ctx context.Context) error {
	promptNeedsUpdate := false
	for {
		if h.GetHallucinationCount() < h.hallucinationCacheSize {
			promptNeedsUpdate = true
			fmt.Printf("hallucinations cache has empty slots, generating more... [%d/%d]\n", len(h.hallucinations)+1, h.hallucinationCacheSize)
			hal, err := h.generateHallucination()
			if err != nil {
				fmt.Println(fmt.Errorf("could not generate hallucination (%v)", err))
			} else {
				if h.isValidResult(hal.Text) {
					h.appendHallucination(hal)
				} else {
					fmt.Println("invalid hallucination, skipping...")
				}
			}
		} else {
			fmt.Println("hallucinations cache is full, waiting for next interval...")
			if promptNeedsUpdate {
				go func() {
					prompts := map[string]int{}
					h.hallucinationLock.Lock()
					for idx := range h.GetHallucinationCount() {
						h.hallucinationCountLock.Lock()
						prompts[h.hallucinations[idx].Prompt] = h.hallucinations[idx].RequestCount
						h.hallucinationCountLock.Unlock()
					}
					h.hallucinationLock.Unlock()
					h.statistics.UpdatePrompts(prompts)
					promptNeedsUpdate = false
				}()
			}
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
