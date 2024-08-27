package webserver_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/command"
	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"codeberg.org/konterfai/konterfai/pkg/webserver"
	"github.com/oklog/run"
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
		hal                            *hallucinator.Hallucinator
		st                             *statistics.Statistics
		baseUrl                        url.URL
		HttpOkProbability, Uncertainty float64
		errorCacheSize                 int
		logger                         *slog.Logger
	)
	BeforeEach(func() {
		ctx = context.Background()
		host = "localhost"
		port = 8080
		hal = hallucinator.NewHallucinator(
			ctx,
			logger,
			5,
			1,
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
			10,
			10,
			10,
			st,
		)
		st = nil
		baseUrl = url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
		}
		HttpOkProbability = 0.5
		Uncertainty = 0.5
		errorCacheSize = 10
		logger, _ = command.SetLogger("off", "")
	})

	Context("NewWebserver", func() {
		It("should return a new webserver", func() {
			ws := webserver.NewWebServer(ctx, logger, host, port, hal, st, baseUrl, HttpOkProbability, Uncertainty, errorCacheSize)
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
			host = "localhost"
			port = 8080
			st = statistics.NewStatistics(ctx, logger, "this is just a dummy string")
			isRobotsTxt := false
			for range 10 {
				st.AppendRequest(ctx, statistics.Request{
					UserAgent:   "test",
					IPAddress:   "127.0.0.1",
					Timestamp:   time.Now(),
					IsRobotsTxt: isRobotsTxt,
					Size:        0,
				})
				isRobotsTxt = !isRobotsTxt
			}
			st.AppendRequest(ctx, statistics.Request{
				UserAgent:   "test",
				IPAddress:   "127.0.0.2",
				Timestamp:   time.Now(),
				IsRobotsTxt: false,
				Size:        0,
			})
			logger, _ = command.SetLogger("off", "")
			ws = webserver.NewWebServer(ctx, logger, host, port, hal, st, baseUrl, HttpOkProbability, Uncertainty, errorCacheSize)
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

		It("should reply with body content", func() {
			// status code is not deterministic, it could be anything
			httpClient := http.Client{
				Timeout: 5 * time.Second,
			}
			resp, err := httpClient.Get("http://localhost:8080")
			Expect(err).NotTo(HaveOccurred())
			bodyData, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bodyData)).To(BeNumerically(">", 0))
			ctx.Done()
		})

		It("should reply with body content on a subdir", func() {
			// status code is not deterministic, it could be anything
			httpClient := http.Client{
				Timeout: 5 * time.Second,
			}
			for range 10 {
				resp, err := httpClient.Get("http://localhost:8080/foobar")
				Expect(err).NotTo(HaveOccurred())
				bodyData, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(bodyData)).To(BeNumerically(">", 0))
			}
			ctx.Done()
		})

		It("should reply with body content on the robots.txt endpoint", func() {
			// status code is not deterministic, it could be anything
			httpClient := http.Client{
				Timeout: 5 * time.Second,
			}
			resp, err := httpClient.Get("http://localhost:8080/robots.txt")
			Expect(err).NotTo(HaveOccurred())
			bodyData, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bodyData)).To(BeNumerically(">", 0))
			ctx.Done()
		})
	})
})
