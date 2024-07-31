package hallucinator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/helpers/dictionaries"
	"codeberg.org/konterfai/konterfai/pkg/helpers/functions"
	"codeberg.org/konterfai/konterfai/pkg/helpers/links"
	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"
	"codeberg.org/konterfai/konterfai/pkg/renderer"
)

// generateFollowUpLink returns a follow-up link for the Hallucinator.
func (h *Hallucinator) generateFollowUpLink(continueText string) string {
	return fmt.Sprintf("<br/><br/><a href=\"%s\">%s</a>",
		links.RandomLink(
			h.hallucinatorUrl,
			h.hallucinatorLinkMaxSubdirectories,
			h.hallucinatorLinkMaxVariables,
			h.hallucinatorLinkHasVariablesProbability,
		),
		continueText,
	)
}

// generateHallucination generates a hallucination from the Ollama API.
func (h *Hallucinator) generateHallucination() (string, error) {
	requestUrl, err := url.JoinPath(h.ollamaAddress, "/api/chat")
	if err != nil {
		fmt.Println("could not join url path")
		os.Exit(1)
	}

	prompt := h.generatePrompt()
	fmt.Printf("generating hallucination with prompt: \"%s\"\n", prompt)

	requestBody := OllamaJsonRequest{
		Model: h.ollamaModel,
		Messages: []OllamaMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Options: OllamaOptions{
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
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println("ollama did not return 200 OK")
		return "", errors.New("ollama did not return 200 OK")
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("could not read response body")
		return "", err
	}
	responses := strings.Split(string(resBody), "\n")
	payload := []string{}
	for _, message := range responses {
		m := OllamaResponse{}
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
	return pl, nil
}

// generatePrompt generates a prompt for the Hallucinator.
func (h *Hallucinator) generatePrompt() string {
	prompt := "write me a story about "
	for i := 0; i < h.promptWordCount; i++ {
		rnd := rand.Intn(100) % 3
		switch rnd {
		case 0:
			prompt += functions.PickRandomFromSlice(&dictionaries.Verbs) + " "
		case 1:
			prompt += functions.PickRandomFromSlice(&dictionaries.Cities) + " "
		case 2:
			fallthrough
		default:
			prompt += functions.PickRandomFromSlice(&dictionaries.Nouns) + " "
		}
	}
	return strings.Trim(prompt, " ") + fmt.Sprintf(", write at least %d words about this. Do not add any comments.", h.hallucinationWordCount)
}

// generateRandomTopicLinks generates random topic links.
func (h *Hallucinator) generateRandomTopicLinks(count int) []renderer.RandomTopic {
	var topics []renderer.RandomTopic
	for i := 0; i < count; i++ {
		topics = append(topics, renderer.RandomTopic{
			Topic: textblocks.RandomTopic(),
			Link:  links.RandomLink(h.hallucinatorUrl, h.hallucinatorLinkMaxSubdirectories, h.hallucinatorLinkMaxVariables, h.hallucinatorLinkHasVariablesProbability),
		})
	}
	return topics
}
