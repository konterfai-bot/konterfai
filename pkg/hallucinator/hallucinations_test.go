package hallucinator_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"codeberg.org/konterfai/konterfai/pkg/command"
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
			10,
			10,
			10,
			st,
		)
	})

	Context("GetHallucinationCount", func() {
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
			hals []hallucinator.Hallucination
		)

		It("should return the current hallucination count when there is a fresh instance", func() {
			Expect(h.GetHallucinationCount(ctx)).To(Equal(0))
		})

		It("should return a default text when no hallucinations are available and using PopHallucination", func() {
			Expect(h.GetHallucinationCount(ctx)).To(Equal(0))
			hal := h.PopHallucination(ctx)
			Expect(hal).To(ContainSubstring(hallucinator.Dream404String))
			Expect(hal).To(ContainSubstring(hallucinator.DreamString))
			Expect(hal).To(ContainSubstring(hallucinator.BackToStartString))
		})

		It("should return a default text when no hallucinations are available and using PopRandomHallucination", func() {
			Expect(h.GetHallucinationCount(ctx)).To(Equal(0))
			hal := h.PopRandomHallucination(ctx)
			Expect(hal).To(ContainSubstring(hallucinator.Dream404String))
			Expect(hal).To(ContainSubstring(hallucinator.DreamString))
			Expect(hal).To(ContainSubstring(hallucinator.BackToStartString))
		})

		It("should return 10 hallucinations when the hallucinations when using PopHallucination", func() {
			for i := range 9 {
				hal := hallucinator.Hallucination{
					RequestCount: 1,
					Prompt:       fmt.Sprintf("dummy hallucination prompt %0.2d", i),
					Text:         fmt.Sprintf("dummy hallucination text %0.2d", i),
				}
				hals = append(hals, hal)
				h.AppendHallucination(ctx, hal)
			}
			h.AppendHallucination(ctx, hallucinator.Hallucination{
				RequestCount: 1,
				Prompt:       "dummy hallucination prompt 10 with length > 255",
				Text:         longHallucinationText,
			})

			Expect(h.GetHallucinationCount(ctx)).To(Equal(10))

			for range 10 {
				hal := h.PopHallucination(ctx)
				Expect(hal).To(ContainSubstring("dummy hallucination text"))
				Expect(hal).To(ContainSubstring(hallucinator.ContinueString))
			}
			h.CleanHallucinations(ctx)
			c := h.GetHallucinationCount(ctx)
			Expect(c).To(BeNumerically("<", 10))
		})

		It("should return 10 hallucinations when the hallucinations when using PopRandomHallucination", func() {
			for i := range 9 {
				hal := hallucinator.Hallucination{
					RequestCount: 1,
					Prompt:       fmt.Sprintf("dummy hallucination prompt %0.2d", i),
					Text:         fmt.Sprintf("dummy hallucination text %0.2d", i),
				}
				hals = append(hals, hal)
				h.AppendHallucination(ctx, hal)
			}
			h.AppendHallucination(ctx, hallucinator.Hallucination{
				RequestCount: 1,
				Prompt:       "dummy hallucination prompt 10 with length > 255",
				Text:         longHallucinationText,
			})

			d := h.GetHallucinationCount(ctx)
			Expect(d).To(Equal(10))

			for range 10 {
				hal := h.PopRandomHallucination(ctx)
				Expect(hal).To(ContainSubstring("dummy hallucination text"))
				Expect(hal).To(ContainSubstring(hallucinator.ContinueString))
			}
			h.CleanHallucinations(ctx)
			c := h.GetHallucinationCount(ctx)
			Expect(c).To(BeNumerically("<", 10))
		})

		It("does not fail when decreasing the hallucination count and the id is < 0", func() {
			h.DecreaseHallucinationRequestCount(ctx, -1)
		})

		It("sets the hallucinationMinimalLength to math.MaxInt when the value is < 1", func() {
			h := hallucinator.NewHallucinator(
				ctx,
				logger,
				5,
				10,
				10,
				10,
				0,
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
			Expect(h.GetHallucinationCount(ctx)).To(Equal(0))
			for i := range 9 {
				hal := hallucinator.Hallucination{
					RequestCount: 1,
					Prompt:       fmt.Sprintf("dummy hallucination prompt %0.2d", i),
					Text:         fmt.Sprintf("dummy hallucination text %0.2d", i),
				}
				hals = append(hals, hal)
				h.AppendHallucination(ctx, hal)
			}
			h.AppendHallucination(ctx, hallucinator.Hallucination{
				RequestCount: 1,
				Prompt:       "dummy hallucination prompt 10 with length > 255",
				Text:         longHallucinationText,
			})

			Expect(h.GetHallucinationCount(ctx)).To(Equal(10))
		})
	})
})
