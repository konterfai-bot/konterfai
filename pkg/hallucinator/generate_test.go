package hallucinator_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/command"
	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Generate", func() {
	const (
		longHallucinationText = "dummy hallucination text " +
			" Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod" +
			" tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et" +
			" justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum" +
			" dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod" +
			" tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam" +
			" et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum" +
			" dolor sit amet."
	)

	var (
		ctx    context.Context
		cancel context.CancelFunc
		h      *hallucinator.Hallucinator
		st     *statistics.Statistics
		logger *slog.Logger
	)

	BeforeEach(func() {
		ctx, cancel = func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		}()
		defer cancel()
		logger, _ = command.SetLogger("off", "")
		st = statistics.NewStatistics(ctx, logger, "this is just a dummy string")
		h = hallucinator.NewHallucinator(
			ctx,
			logger,
			5,
			10,
			10,
			10,
			500,
			10,
			10,
			10,
			10,
			10,
			url.URL{
				Scheme: "http",
				Host:   "localhost:8080",
			},
			"http://localhost:11434",
			"dummy",
			1,
			10,
			10,
			st,
		)
	})

	It("should return an error if backend is unreachable", func() {
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).To(HaveOccurred())
		Expect(hal).To(Equal(hallucinator.Hallucination{}))
	})

	It("should return an error if ollama does not return 200 OK", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).To(HaveOccurred())
		Expect(hal).To(Equal(hallucinator.Hallucination{}))
	})

	It("should return an error if ollama does not return 200 OK and body is nil", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
		}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).To(HaveOccurred())
		Expect(hal).To(Equal(hallucinator.Hallucination{}))
	})

	It("should return an error if ollama does not return 200 OK and body is empty", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(``)),
		}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).To(HaveOccurred())
		Expect(hal).To(Equal(hallucinator.Hallucination{}))
	})

	It("returns a hallucination if ollama returns 200 OK", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"message": {"content": "` + longHallucinationText + `"}}`)),
		}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(hal).NotTo(Equal(hallucinator.Hallucination{}))
	})

	It("should return a valid hallucination", func() {
		ollamaResponse := hallucinator.OllamaResponse{
			Message: hallucinator.OllamaMessage{
				Role:    "test",
				Content: longHallucinationText,
			},
			Done: true,
		}
		ollamaResponseJSON, err := json.Marshal(ollamaResponse)
		Expect(err).NotTo(HaveOccurred())
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(string(ollamaResponseJSON))),
		}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(hal.Text).To(Equal(ollamaResponse.Message.Content))
	})

	It("should return an error if the hallucination matches a regexp", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"message": {"content": "Sorry, but I can't assist with that,foobar"}}`)),
		}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).To(HaveOccurred())
		Expect(hal).To(Equal(hallucinator.Hallucination{}))
	})

	It("should return an error if the hallucination response is malformed", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"message": {"content": "This content is valid"`)),
		}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("unexpected end of JSON input"))
		Expect(hal).To(Equal(hallucinator.Hallucination{}))
	})

	It("should return an error if the hallucination is < 500 characters", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"message": {"content": "This is a short response"}}`)),
		}, nil)
		h.HTTPClient = mockHttpClient
		hal, err := h.GenerateHallucination(ctx)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError("ollama returned a hallucination that is too short"))
		Expect(hal).To(Equal(hallucinator.Hallucination{}))
	})
})
