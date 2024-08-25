package hallucinator_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MockHttpClient is a mock implementation of the http.Client interface.
type MockHttpClient struct {
	mock.Mock
}

// Do is a mock implementation of the http.Client Do method.
func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

var _ = Describe("Hallucinator", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		h      *hallucinator.Hallucinator
		st     *statistics.Statistics
	)

	BeforeEach(func() {
		ctx, cancel = func() (context.Context, context.CancelFunc) {
			return context.WithCancel(context.Background())
		}()
		defer cancel()
		st = statistics.NewStatistics(ctx, "this is just a dummy string")
		h = hallucinator.NewHallucinator(
			ctx,
			5,
			1,
			10,
			10,
			10,
			10,
			10,
			10,
			10,
			url.URL{
				Scheme: "http",
				Host:   "localhost:8080",
			},
			"dummy",
			"dummy",
			10,
			10,
			10,
			st,
		)
	})

	It("Start the hallucinator", func() {
		go func() {
			err := h.Start(ctx)
			Expect(err).To(BeNil())
		}()
		time.Sleep(1 * time.Second)
		ctx.Done()
	})

	It("Starts the hallucinator with mocks", func() {
		mockHttpClient := new(MockHttpClient)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"message": {"content": "This is a valid response"}}`)),
		}, nil)
		h.HTTPClient = mockHttpClient
		go func() {
			err := h.Start(ctx)
			Expect(err).To(BeNil())
		}()
		time.Sleep(2 * time.Second)
	})
})
