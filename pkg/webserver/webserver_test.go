package webserver_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/oklog/run"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"codeberg.org/konterfai/konterfai/pkg/webserver"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWebserver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webserver Suite")
}

var _ = Describe("Webserver", func() {
	var (
		ctx                            context.Context
		host                           string
		port                           int
		hallucinator                   *hallucinator.Hallucinator
		statistics                     *statistics.Statistics
		baseUrl                        url.URL
		HttpOkProbability, Uncertainty float64
		errorCacheSize                 int
	)
	BeforeEach(func() {
		ctx = context.Background()
		host = "localhost"
		port = 8080
		hallucinator = nil
		statistics = nil
		baseUrl = url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
		}
		HttpOkProbability = 0.5
		Uncertainty = 0.5
		errorCacheSize = 10
	})

	Context("NewWebserver", func() {
		It("should return a new webserver", func() {
			ws := webserver.NewWebServer(ctx, host, port, hallucinator, statistics, baseUrl, HttpOkProbability, Uncertainty, errorCacheSize)
			Expect(ws).NotTo(BeNil())
			Expect(ws.Host).To(Equal(host))
			Expect(ws.Port).To(Equal(port))
			Expect(ws.HTTPBaseURL).To(Equal(baseUrl))
			Expect(ws.HTTPOkProbability).To(Equal(HttpOkProbability))
			Expect(ws.Uncertainty).To(Equal(Uncertainty))
		})
	})

	Context("Serve", func() {
		var (
			ws  *webserver.WebServer
			err error
		)
		BeforeEach(func() {
			ws = webserver.NewWebServer(ctx, host, port, hallucinator, statistics, baseUrl, HttpOkProbability, Uncertainty, errorCacheSize)
			syncer := make(chan error)
			gr := run.Group{}
			gr.Add(func() error {
				select {
				case <-ctx.Done():
					return nil
				case syncer <- ws.Serve(ctx):
					return <-syncer
				}
			}, func(err error) {
				Expect(err).ToNot(HaveOccurred())
			})

			go func() {
				err = gr.Run()
				Expect(err).ToNot(HaveOccurred())
			}()
		})

		It("should start the web server", func() {
			time.Sleep(1 * time.Second)
			ctx.Done()
			Expect(err).ToNot(HaveOccurred())
		})

		// TODO: In the future here could be more integration tests for actually checking the responses
	})
})
