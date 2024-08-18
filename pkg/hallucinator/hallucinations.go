package hallucinator

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"

	"codeberg.org/konterfai/konterfai/pkg/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"
	"codeberg.org/konterfai/konterfai/pkg/renderer"
)

// GetHallucinationCount returns the current hallucination count.
func (h *Hallucinator) GetHallucinationCount(ctx context.Context) int {
	ctx, span := tracer.Start(ctx, "Hallucinator.GetHallucinationCount")
	defer span.End()

	h.hallucinationCountLock.Lock()
	defer h.hallucinationCountLock.Unlock()
	return h.hallucinationCount
}

// DecreaseHallucinationRequestCount decreases the request count of a hallucination by 1.
func (h *Hallucinator) DecreaseHallucinationRequestCount(ctx context.Context, id int) {
	ctx, span := tracer.Start(ctx, "Hallucinator.DecreaseHallucinationRequestCount")
	defer span.End()

	go func() {
		h.hallucinationLock.Lock()
		defer h.hallucinationLock.Unlock()
		if id < 0 || h.GetHallucinationCount(ctx) <= id {
			return
		}
		h.hallucinations[id].RequestCount--
	}()
}

// PopHallucination withdraws the first hallucination from the list of hallucinations.
func (h *Hallucinator) PopHallucination(ctx context.Context) string {
	ctx, span := tracer.Start(ctx, "Hallucinator.PopHallucination")
	defer span.End()

	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.cleanHallucinations(ctx)
	if h.GetHallucinationCount(ctx) < 1 {
		hallucination, err := h.renderer.RenderInRandomTemplate(ctx,
			renderer.RenderData{
				NewsAnchor:   textblocks.RandomNewsPaperName(ctx),
				Headline:     dream404String,
				Content:      dreamString,
				FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, backToStartString)),
				RandomTopics: h.generateRandomTopicLinks(ctx, 10),
				Year:         functions.PickRandomYear(ctx),
				MetaData: renderer.MetaData{
					Description: dreamString,
					Keywords:    textblocks.RandomKeywords(ctx, 10),
					Charset:     functions.PickRandomStringFromSlice(ctx, &dictionaries.Charsets),
				},
				LanguageCode: functions.PickRandomStringFromSlice(ctx, &dictionaries.LanguageCodes),
			})
		if err != nil {
			return fmt.Sprintf("Could not render template, error: %v", err)
		}
		return hallucination
	}
	currentHallucination := h.hallucinations[0].Text
	var metaDescription string
	if len(currentHallucination) < 255 {
		metaDescription = currentHallucination
	} else {
		metaDescription = currentHallucination[:255]
	}
	h.DecreaseHallucinationRequestCount(ctx, 0)
	hallucination, err := h.renderer.RenderInRandomTemplate(ctx,
		renderer.RenderData{
			NewsAnchor:   textblocks.RandomNewsPaperName(ctx),
			Headline:     textblocks.RandomHeadline(ctx),
			Content:      template.HTML(h.clutterTextWithRandomHref(ctx, h.hallucinations[0].Text)),
			FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, continueString)),
			RandomTopics: h.generateRandomTopicLinks(ctx, 10),
			Year:         functions.PickRandomYear(ctx),
			MetaData: renderer.MetaData{
				Description: metaDescription,
				Keywords:    textblocks.RandomKeywords(ctx, 10),
				Charset:     functions.PickRandomStringFromSlice(ctx, &dictionaries.Charsets),
			},
			LanguageCode: functions.PickRandomStringFromSlice(ctx, &dictionaries.LanguageCodes),
		})
	if err != nil {
		return fmt.Sprintf("Could not render template, error: %v", err)
	}
	return hallucination
}

// PopRandomHallucination withdraws a random hallucination from the list of hallucinations.
func (h *Hallucinator) PopRandomHallucination(ctx context.Context) string {
	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.cleanHallucinations(ctx)
	if h.GetHallucinationCount(ctx) < 1 {
		hallucination, err := h.renderer.RenderInRandomTemplate(ctx,
			renderer.RenderData{
				NewsAnchor:   textblocks.RandomNewsPaperName(ctx),
				Headline:     dream404String,
				Content:      dreamString,
				FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, backToStartString)),
				RandomTopics: h.generateRandomTopicLinks(ctx, 10),
				Year:         functions.PickRandomYear(ctx),
				MetaData: renderer.MetaData{
					Description: dreamString,
					Keywords:    textblocks.RandomKeywords(ctx, 10),
					Charset:     functions.PickRandomStringFromSlice(ctx, &dictionaries.Charsets),
				},
				LanguageCode: functions.PickRandomStringFromSlice(ctx, &dictionaries.LanguageCodes),
			})
		if err != nil {
			return fmt.Sprintf("Could not render template, error: %v", err)
		}
		return hallucination
	}
	randomIndex := rand.Intn(h.GetHallucinationCount(ctx))
	h.DecreaseHallucinationRequestCount(ctx, randomIndex)
	currentHallucination := h.hallucinations[randomIndex].Text
	var metaDescription string
	if len(currentHallucination) < 255 {
		metaDescription = currentHallucination
	} else {
		metaDescription = currentHallucination[:255]
	}
	hallucination, err := h.renderer.RenderInRandomTemplate(ctx,
		renderer.RenderData{
			NewsAnchor:   textblocks.RandomNewsPaperName(ctx),
			Headline:     textblocks.RandomHeadline(ctx),
			Content:      template.HTML(h.clutterTextWithRandomHref(ctx, currentHallucination)),
			FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, continueString)),
			RandomTopics: h.generateRandomTopicLinks(ctx, 10),
			Year:         functions.PickRandomYear(ctx),
			MetaData: renderer.MetaData{
				Description: metaDescription,
				Keywords:    textblocks.RandomKeywords(ctx, 10),
				Charset:     functions.PickRandomStringFromSlice(ctx, &dictionaries.Charsets),
			},
			LanguageCode: functions.PickRandomStringFromSlice(ctx, &dictionaries.LanguageCodes),
		})
	if err != nil {
		return fmt.Sprintf("Could not render template, error: %v", err)
	}
	return hallucination
}

// appendHallucination appends a hallucination to the list of hallucinations.
func (h *Hallucinator) appendHallucination(ctx context.Context, hallucination Hallucination) {
	ctx, span := tracer.Start(ctx, "Hallucinator.appendHallucination")
	defer span.End()

	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.hallucinations = append(h.hallucinations, hallucination)
	h.setHallucinationCount(ctx)
}

// cleanHallucinations cleans the list of hallucinations and removes hallucinations with requestCount 0.
func (h *Hallucinator) cleanHallucinations(ctx context.Context) {
	// This function does not have a lock on the hallucinations list. It is expected that the caller has locked the list.
	// This happens in PopHallucination and PopRandomHallucination.
	ctx, span := tracer.Start(ctx, "Hallucinator.cleanHallucinations")
	defer span.End()

	if h.GetHallucinationCount(ctx) < 1 {
		return
	}
	newHallucinations := []Hallucination{}
	for _, hallucination := range h.hallucinations {
		if hallucination.RequestCount > 0 {
			newHallucinations = append(newHallucinations, hallucination)
		}
	}
	h.hallucinations = newHallucinations
	h.setHallucinationCount(ctx)
}

// setHallucinationCount sets the hallucination count from the length of the hallucination slice.
func (h *Hallucinator) setHallucinationCount(ctx context.Context) {
	// This function does not have a lock on the hallucinations list. It is expected that the caller has locked the list.
	// This happens in AppendHallucination and CleanHallucinations.
	ctx, span := tracer.Start(ctx, "Hallucinator.setHallucinationCount")
	defer span.End()

	h.hallucinationCountLock.Lock()
	defer h.hallucinationCountLock.Unlock()
	h.hallucinationCount = len(h.hallucinations)
}
