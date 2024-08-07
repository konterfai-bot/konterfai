package hallucinator

import (
	"codeberg.org/konterfai/konterfai/pkg/helpers/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"
	"codeberg.org/konterfai/konterfai/pkg/renderer"
	"fmt"
	"html/template"
	"math/rand"
)

// GetHallucinationCount returns the current hallucination count.
func (h *Hallucinator) GetHallucinationCount() int {
	h.hallucinationCountLock.Lock()
	defer h.hallucinationCountLock.Unlock()
	return h.hallucinationCount
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
			MetaData: renderer.MetaData{
				Description: dreamString,
				Keywords:    textblocks.RandomKeywords(10),
				Charset:     functions.PickRandomStringFromSlice(&dictionaries.Charsets),
			},
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
	h.hallucinations[0].RequestCount--
	hallucination, err := h.renderer.RenderInRandomTemplate(renderer.RenderData{
		NewsAnchor:   textblocks.RandomNewsPaperName(),
		Headline:     textblocks.RandomHeadline(),
		Content:      template.HTML(h.clutterTextWithRandomHref(h.hallucinations[0].Text)),
		FollowUpLink: template.HTML(h.generateFollowUpLink(continueString)),
		RandomTopics: h.generateRandomTopicLinks(10),
		Year:         functions.PickRandomYear(),
		MetaData: renderer.MetaData{
			Description: metaDescription,
			Keywords:    textblocks.RandomKeywords(10),
			Charset:     functions.PickRandomStringFromSlice(&dictionaries.Charsets),
		},
	})
	if err != nil {
		return fmt.Sprintf("Could not render template, error: %v", err)
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
			MetaData: renderer.MetaData{
				Description: dreamString,
				Keywords:    textblocks.RandomKeywords(10),
				Charset:     functions.PickRandomStringFromSlice(&dictionaries.Charsets),
			},
		})
		if err != nil {
			return fmt.Sprintf("Could not render template, error: %v", err)
		}
		return hallucination
	}
	randomIndex := rand.Intn(len(h.hallucinations))
	h.hallucinations[randomIndex].RequestCount--
	currentHallucination := h.hallucinations[randomIndex].Text
	var metaDescription string
	if len(currentHallucination) < 255 {
		metaDescription = currentHallucination
	} else {
		metaDescription = currentHallucination[:255]
	}
	hallucination, err := h.renderer.RenderInRandomTemplate(renderer.RenderData{
		NewsAnchor:   textblocks.RandomNewsPaperName(),
		Headline:     textblocks.RandomHeadline(),
		Content:      template.HTML(h.clutterTextWithRandomHref(currentHallucination)),
		FollowUpLink: template.HTML(h.generateFollowUpLink(continueString)),
		RandomTopics: h.generateRandomTopicLinks(10),
		Year:         functions.PickRandomYear(),
		MetaData: renderer.MetaData{
			Description: metaDescription,
			Keywords:    textblocks.RandomKeywords(10),
			Charset:     functions.PickRandomStringFromSlice(&dictionaries.Charsets),
		},
	})
	if err != nil {
		return fmt.Sprintf("Could not render template, error: %v", err)
	}
	return hallucination
}

// appendHallucination appends a hallucination to the list of hallucinations.
func (h *Hallucinator) appendHallucination(hallucination Hallucination) {
	h.hallucinationLock.Lock()
	defer h.hallucinationLock.Unlock()
	h.hallucinations = append(h.hallucinations, hallucination)
	h.setHallucinationCount()
}

// cleanHallucinations cleans the list of hallucinations and removes hallucinations with requestCount 0.
func (h *Hallucinator) cleanHallucinations() {
	// This function does not have a lock on the hallucinations list. It is expected that the caller has locked the list.
	// This happens in PopHallucination and PopRandomHallucination.
	if h.GetHallucinationCount() < 1 {
		return
	}
	newHallucinations := []Hallucination{}
	for _, hallucination := range h.hallucinations {
		if hallucination.RequestCount > 0 {
			newHallucinations = append(newHallucinations, hallucination)
		}
	}
	h.hallucinations = newHallucinations
	h.setHallucinationCount()
}

// setHallucinationCount sets the hallucination count from the length of the hallucination slice.
func (h *Hallucinator) setHallucinationCount() {
	// This function does not have a lock on the hallucinations list. It is expected that the caller has locked the list.
	// This happens in AppendHallucination and CleanHallucinations.
	h.hallucinationCountLock.Lock()
	defer h.hallucinationCountLock.Unlock()
	h.hallucinationCount = len(h.hallucinations)
}
