package statisticsserver_test

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/command"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"codeberg.org/konterfai/konterfai/pkg/statisticsserver"
	"github.com/oklog/run"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Statisticsserver", func() {
	var (
		ctx    context.Context
		Host   string
		Port   int
		st     *statistics.Statistics
		logger *slog.Logger
	)

	BeforeEach(func() {
		ctx = context.Background()
		logger, _ = command.SetLogger("off", "")
	})

	Context("NewStatisticsServer", func() {
		It("should return a new statistics server", func() {
			ss := statisticsserver.NewStatisticsServer(ctx, logger, Host, Port, st)
			Expect(ss).NotTo(BeNil())
		})
	})

	Context("Start", func() {
		var (
			ss  *statisticsserver.StatisticsServer
			err error
		)
		BeforeEach(func() {
			Host = "localhost"
			Port = 8081
			logger, _ = command.SetLogger("off", "")
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
			ss = statisticsserver.NewStatisticsServer(ctx, logger, Host, Port, st)
			syncer := make(chan error)
			gr := run.Group{}
			gr.Add(func() error {
				select {
				case <-ctx.Done():
					return nil
				case syncer <- ss.Serve(ctx):
					return <-syncer
				}
			}, func(error) {
				Expect(err).NotTo(HaveOccurred())
			})

			go func() {
				err = gr.Run()
				Expect(err).NotTo(HaveOccurred())
			}()
		})

		It("should start the server", func() {
			time.Sleep(1 * time.Second)
			ctx.Done()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reply with a 200 status code and body content", func() {
			httpClient := http.Client{
				Timeout: 5 * time.Second,
			}
			resp, err := httpClient.Get("http://localhost:8081")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			bodyData, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(bodyData)).To(ContainSubstring("this is just a dummy string"))
			ctx.Done()
		})

		It("should reply with a 200 status code and body content on the metrics endpoint", func() {
			httpClient := http.Client{
				Timeout: 5 * time.Second,
			}
			resp, err := httpClient.Get("http://localhost:8081/metrics")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			bodyData, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(bodyData)).To(ContainSubstring("konterfai_requests_total"))
			ctx.Done()
		})
	})
})
