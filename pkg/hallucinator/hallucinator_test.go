package hallucinator_test

import (
	"context"
	"net/url"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			10,
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
})
