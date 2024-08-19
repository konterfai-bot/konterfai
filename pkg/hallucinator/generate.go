package hallucinator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/links"
	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"
	"codeberg.org/konterfai/konterfai/pkg/renderer"
)

// generateFollowUpLink returns a follow-up link for the Hallucinator.
func (h *Hallucinator) generateFollowUpLink(ctx context.Context, continueText string) string {
	ctx, span := tracer.Start(ctx, "Hallucinator.generateFollowUpLink")
	defer span.End()

	return fmt.Sprintf("<br/><br/><a href=\"%s\">%s</a>",
		links.RandomLink(
			ctx,
			h.hallucinatorUrl,
			h.hallucinatorLinkMaxSubdirectories,
			h.hallucinatorLinkMaxVariables,
			h.hallucinatorLinkHasVariablesProbability,
		),
		continueText,
	)
}

// generateHallucination generates a hallucination from the Ollama API.
func (h *Hallucinator) generateHallucination(ctx context.Context) (Hallucination, error) {
	ctx, span := tracer.Start(ctx, "Hallucinator.generateHallucination")
	defer span.End()
	requestUrl, err := url.JoinPath(h.ollamaAddress, "/api/chat")
	if err != nil {
		fmt.Println("could not join url path")
		os.Exit(1)
	}

	prompt := h.generatePrompt(ctx)
	fmt.Printf("generating hallucination with prompt: \"%s\"\n", prompt)

	requestBody := ollamaJsonRequest{
		Model: h.ollamaModel,
		Messages: []ollamaMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Options: ollamaOptions{
			Temperature: h.aiTemperature,
			Seed:        h.aiSeed,
		},
	}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("could not marshal request body")
		os.Exit(1)
	}

	res, err := h.httpClient.Post(requestUrl, "application/json", bytes.NewReader(requestBodyJson))
	if err != nil {
		fmt.Println(fmt.Errorf("could not get hallucination from ollama (%v)", err))
		return Hallucination{}, err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("ollama did not return 200 OK")
		return Hallucination{}, errors.New("ollama did not return 200 OK")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("could not read response body")
		return Hallucination{}, err
	}
	responses := strings.Split(string(resBody), "\n")
	payload := []string{}
	for _, message := range responses {
		m := ollamaResponse{}
		err := json.Unmarshal([]byte(message), &m)
		if err != nil {
			continue
		}
		msg := strings.Trim(m.Message.Content, " ")
		if msg != "" && msg != "\n" {
			payload = append(payload, msg)
		}
		if m.Done {
			break
		}
	}
	pl := strings.Join(payload, " ")
	return Hallucination{
		Text:         pl,
		Prompt:       prompt,
		RequestCount: h.hallucinationRequestCount,
	}, nil
}

// generatePrompt generates a prompt for the Hallucinator.
func (h *Hallucinator) generatePrompt(ctx context.Context) string {
	ctx, span := tracer.Start(ctx, "Hallucinator.generatePrompt")
	defer span.End()
	words := ""
	for i := 0; i < h.promptWordCount; i++ {
		rnd := rand.Intn(100) % 3
		switch rnd {
		case 0:
			words += functions.PickRandomStringFromSlice(ctx, &dictionaries.Verbs) + " "
		case 1:
			words += functions.PickRandomStringFromSlice(ctx, &dictionaries.Cities) + " "
		case 2:
			fallthrough
		default:
			words += functions.PickRandomStringFromSlice(ctx, &dictionaries.Nouns) + " "
		}
	}
	return fmt.Sprintf(functions.PickRandomStringFromSlice(ctx, &dictionaries.Prompts),
		functions.PickRandomStringFromSlice(ctx, &dictionaries.ArticleTypes),
		words,
		h.hallucinationWordCount,
		functions.PickRandomStringFromSlice(ctx, &dictionaries.Languages),
	)
}

// generateRandomTopicLinks generates random topic links.
func (h *Hallucinator) generateRandomTopicLinks(ctx context.Context, count int) []renderer.RandomTopic {
	ctx, span := tracer.Start(ctx, "Hallucinator.generateRandomTopicLinks")
	defer span.End()
	var topics []renderer.RandomTopic
	for i := 0; i < count; i++ {
		topics = append(topics, renderer.RandomTopic{
			Topic: textblocks.RandomTopic(ctx),
			Link: links.RandomLink(ctx,
				h.hallucinatorUrl,
				h.hallucinatorLinkMaxSubdirectories,
				h.hallucinatorLinkMaxVariables,
				h.hallucinatorLinkHasVariablesProbability,
			),
		})
	}
	return topics
}
