package hallucinator

import (
	"context"
	"fmt"
	"html/template"
	"math/rand"

	"codeberg.org/konterfai/konterfai/pkg/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"
	"codeberg.org/konterfai/konterfai/pkg/renderer/v0"
)

// GetHallucinationCount returns the current hallucination count.
func (h *Hallucinator) GetHallucinationCount(ctx context.Context) int {
	_, span := tracer.Start(ctx, "Hallucinator.GetHallucinationCount")
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
	h.CleanHallucinations(ctx)
	if h.GetHallucinationCount(ctx) < 1 {
		hallucination, err := h.renderer.RenderInRandomTemplate(ctx,
			renderer.RenderData{
				NewsAnchor:   textblocks.RandomNewsPaperName(ctx),
				Headline:     Dream404String,
				Content:      DreamString,
				FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, BackToStartString)), //nolint: gosec
				RandomTopics: h.generateRandomTopicLinks(ctx),
				Year:         functions.PickRandomYear(ctx),
				MetaData: renderer.MetaData{
					Description: DreamString,
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
			Content:      template.HTML(h.clutterTextWithRandomHref(ctx, h.hallucinations[0].Text)), //nolint: gosec
			FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, ContinueString)),                //nolint: gosec
			RandomTopics: h.generateRandomTopicLinks(ctx),
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
	if h.GetHallucinationCount(ctx) < 1 {
		hallucination, err := h.renderer.RenderInRandomTemplate(ctx,
			renderer.RenderData{
				NewsAnchor:   textblocks.RandomNewsPaperName(ctx),
				Headline:     Dream404String,
				Content:      DreamString,
				FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, BackToStartString)), //nolint: gosec
				RandomTopics: h.generateRandomTopicLinks(ctx),
				Year:         functions.PickRandomYear(ctx),
				MetaData: renderer.MetaData{
					Description: DreamString,
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
	randomIndex := rand.Intn(h.GetHallucinationCount(ctx)) //nolint: gosec
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
			Content:      template.HTML(h.clutterTextWithRandomHref(ctx, currentHallucination)), //nolint: gosec
			FollowUpLink: template.HTML(h.generateFollowUpLink(ctx, ContinueString)),            //nolint: gosec
			RandomTopics: h.generateRandomTopicLinks(ctx),
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
	h.CleanHallucinations(ctx)

	return hallucination
}

// AppendHallucination appends a hallucination to the list of hallucinations.
func (h *Hallucinator) AppendHallucination(ctx context.Context, hallucination Hallucination) {
	ctx, span := tracer.Start(ctx, "Hallucinator.AppendHallucination")
	defer span.End()

	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.hallucinations = append(h.hallucinations, hallucination)
	h.setHallucinationCount(ctx)
}

// CleanHallucinations cleans the list of hallucinations and removes hallucinations with requestCount 0.
func (h *Hallucinator) CleanHallucinations(ctx context.Context) {
	// This function does not have a lock on the hallucinations list. It is expected that the caller has locked the list.
	// This happens in PopHallucination and PopRandomHallucination.
	ctx, span := tracer.Start(ctx, "Hallucinator.CleanHallucinations")
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
	_, span := tracer.Start(ctx, "Hallucinator.setHallucinationCount")
	defer span.End()

	h.hallucinationCountLock.Lock()
	defer h.hallucinationCountLock.Unlock()
	h.hallucinationCount = len(h.hallucinations)
}
