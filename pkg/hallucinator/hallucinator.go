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
	"go.opentelemetry.io/otel"
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
	hallucinatorURL                         url.URL
	ollamaAddress                           string
	ollamaModel                             string
	aiTemperature                           float64
	aiSeed                                  int
	promptWordCount                         int

	httpClient *http.Client
	renderer   *renderer.Renderer
	statistics *statistics.Statistics
}

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/hallucinator")

// NewHallucinator creates a new Hallucinator instance.
func NewHallucinator(ctx context.Context, interval time.Duration,
	hallucinationCacheSize int,
	hallucinatorPromptWordCount int,
	hallucinationRequestCount int,
	hallucinationWordCount int,
	hallucinationLinkPercentage int,
	hallucinatorLinkMaxSubdirectories int,
	hallucinatorLinkHasVariablesProbability float64,
	hallucinatorLinkMaxVariables int,
	hallucinatorURL url.URL,
	ollamaAddress string,
	ollamaModel string,
	ollamaRequestTimeOut time.Duration,
	aiTemperature float64,
	aiSeed int,
	statistics *statistics.Statistics,
) *Hallucinator {
	ctx, span := tracer.Start(ctx, "Hallucinator.NewHallucinator")
	defer span.End()

	headLineLinks := [10]string{}
	for i := 0; i < len(headLineLinks); i++ {
		headLineLinks[i] = links.RandomLink(ctx,
			hallucinatorURL,
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
		hallucinatorURL:                         hallucinatorURL,
		ollamaAddress:                           ollamaAddress,
		ollamaModel:                             ollamaModel,
		aiTemperature:                           aiTemperature,
		aiSeed:                                  aiSeed,
		promptWordCount:                         hallucinatorPromptWordCount,

		httpClient: &http.Client{
			Timeout: ollamaRequestTimeOut,
		},
		renderer:   renderer.NewRenderer(ctx, headLineLinks[:]),
		statistics: statistics,
	}
}

// Start starts the Hallucinator.
func (h *Hallucinator) Start(ctx context.Context) error {
	// No need to trace this function as it is the entry point and an endless loop.
	promptNeedsUpdate := false
	for {
		if h.GetHallucinationCount(ctx) < h.hallucinationCacheSize {
			promptNeedsUpdate = true
			fmt.Printf("hallucinations cache has empty slots, generating more... [%d/%d]\n",
				len(h.hallucinations)+1, h.hallucinationCacheSize)
			hal, err := h.generateHallucination(ctx)
			if err != nil {
				functions.SleepWithContext(ctx, h.Interval)
				fmt.Printf("could not generate hallucination (%v)\n", err)

				continue
			}
			if h.isValidResult(ctx, hal.Text) {
				h.AppendHallucination(ctx, hal)

				// Update Prometheus metrics
				statistics.PromptsGeneratedTotal.Inc()
			} else {
				fmt.Println("invalid hallucination, skipping...")
			}

			continue
		}
		fmt.Println("hallucinations cache is full, waiting for next interval...")
		if promptNeedsUpdate {
			go func() {
				prompts := map[string]int{}
				h.hallucinationLock.Lock()
				for idx := range h.GetHallucinationCount(ctx) {
					h.hallucinationCountLock.Lock()
					prompts[h.hallucinations[idx].Prompt] = h.hallucinations[idx].RequestCount
					h.hallucinationCountLock.Unlock()
				}
				h.hallucinationLock.Unlock()
				h.statistics.UpdatePrompts(ctx, prompts)
			}()
			promptNeedsUpdate = false
		}
		functions.SleepWithContext(ctx, h.Interval)
	}
}

// clutterTextWithRandomHref clutters the given text with random hrefs.
func (h *Hallucinator) clutterTextWithRandomHref(ctx context.Context, text string) string {
	ctx, span := tracer.Start(ctx, "Hallucinator.clutterTextWithRandomHref")
	defer span.End()

	textSlice := strings.Split(text, " ")
	percentile := len(textSlice) * h.hallucinationLinkPercentage / 100
	generated := map[int]bool{}
	for {
		if len(generated) >= percentile {
			break
		}
		i := rand.Intn(len(textSlice)) //nolint: gosec
		if !generated[i] {
			textSlice[i] = fmt.Sprintf("<a href=\"%s\">%s</a>",
				links.RandomLink(ctx,
					h.hallucinatorURL,
					h.hallucinatorLinkMaxSubdirectories,
					h.hallucinatorLinkMaxVariables,
					h.hallucinatorLinkHasVariablesProbability,
				), textSlice[i])
			generated[i] = true
		}
	}

	return strings.Join(textSlice, " ")
}
