package statisticsserver_test

import (
	"context"
	"github.com/oklog/run"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"codeberg.org/konterfai/konterfai/pkg/statisticsserver"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Statisticsserver", func() {
	var (
		ctx        context.Context
		Host       string
		Port       int
		Statistics *statistics.Statistics
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("NewStatisticsServer", func() {
		It("should return a new statistics server", func() {
			ss := statisticsserver.NewStatisticsServer(ctx, Host, Port, Statistics)
			Expect(ss).NotTo(BeNil())
		})
	})

	Context("Start", func() {
		var (
			ss  *statisticsserver.StatisticsServer
			err error
		)
		BeforeEach(func() {
			ss = statisticsserver.NewStatisticsServer(ctx, Host, Port, Statistics)
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
	})
})
