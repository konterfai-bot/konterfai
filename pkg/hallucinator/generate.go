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
	"runtime"
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
			h.hallucinatorURL,
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
	requestURL, err := url.JoinPath(h.ollamaAddress, "/api/chat")
	if err != nil {
		fmt.Printf("could not join url path (%v)", err)
		defer os.Exit(1)
		runtime.Goexit()
	}
	prompt := h.generatePrompt(ctx)
	fmt.Printf("generating hallucination with prompt: \"%s\"\n", prompt)
	requestBody := ollamaJSONRequest{
		Model: h.ollamaModel, Messages: []ollamaMessage{{Role: "user", Content: prompt}},
		Options: ollamaOptions{Temperature: h.aiTemperature, Seed: h.aiSeed},
	}
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("could not marshal request body")
		defer os.Exit(1)
		runtime.Goexit()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(requestBodyJSON))
	if err != nil {
		fmt.Printf("could not create request (%v)\n", err)

		return Hallucination{}, err
	}
	res, err := h.httpClient.Do(req)
	if err != nil {
		fmt.Printf("could not get hallucination from ollama (%v)\n", err)

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
	pl := concatOllamaMessages(resBody)
	if err := res.Body.Close(); err != nil {
		fmt.Printf("could not close response body (%v)\n", err)
	}

	return Hallucination{Text: pl, Prompt: prompt, RequestCount: h.hallucinationRequestCount}, nil
}

// concatOllamaMessages concatenates Ollama messages.
func concatOllamaMessages(responseBody []byte) string {
	responses := strings.Split(string(responseBody), "\n")
	var payload []string
	for _, message := range responses {
		m := ollamaResponse{}
		if err := json.Unmarshal([]byte(message), &m); err != nil {
			continue
		}
		if msg := strings.Trim(m.Message.Content, " "); msg != "" && msg != "\n" {
			payload = append(payload, msg)
		}
		if m.Done {
			break
		}
	}

	return strings.Join(payload, " ")
}

// generatePrompt generates a prompt for the Hallucinator.
func (h *Hallucinator) generatePrompt(ctx context.Context) string {
	ctx, span := tracer.Start(ctx, "Hallucinator.generatePrompt")
	defer span.End()
	words := ""
	for range h.promptWordCount {
		rnd := rand.Intn(100) % 3 //nolint: gosec
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
func (h *Hallucinator) generateRandomTopicLinks(ctx context.Context) []renderer.RandomTopic {
	ctx, span := tracer.Start(ctx, "Hallucinator.generateRandomTopicLinks")
	defer span.End()
	topics := make([]renderer.RandomTopic, 0, 10)
	for range 10 {
		topics = append(topics, renderer.RandomTopic{
			Topic: textblocks.RandomTopic(ctx),
			Link: links.RandomLink(ctx,
				h.hallucinatorURL,
				h.hallucinatorLinkMaxSubdirectories,
				h.hallucinatorLinkMaxVariables,
				h.hallucinatorLinkHasVariablesProbability,
			),
		})
	}

	return topics
}
