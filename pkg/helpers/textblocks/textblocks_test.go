package textblocks_test

import (
	"context"
	"testing"

	"codeberg.org/konterfai/konterfai/pkg/helpers/textblocks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTextblocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Functions Suite")
}

var _ = Describe("Textblocks", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("RandomHeadline", func() {
		It("should return a random headline", func() {
			Expect(textblocks.RandomHeadline(ctx)).NotTo(BeEmpty())
		})

		It("should match the expected format", func() {
			Expect(textblocks.RandomHeadline(ctx)).To(MatchRegexp(`.*:.* .* .*`))
		})
	})

	Context("RandomKeywords", func() {
		It("should return n random keywords", func() {
			Expect(textblocks.RandomKeywords(ctx, 3)).NotTo(BeEmpty())
		})

		It("should match the expected format", func() {
			Expect(textblocks.RandomKeywords(ctx, 3)).To(MatchRegexp(`.*,.*,.*`))
		})
	})

	Context("RandomNewsPaperName", func() {
		It("should return a random newspaper name", func() {
			Expect(textblocks.RandomNewsPaperName(ctx)).NotTo(BeEmpty())
		})

		It("should match the expected format", func() {
			Expect(textblocks.RandomNewsPaperName(ctx)).To(MatchRegexp(`.* .*`))
		})
	})
})
