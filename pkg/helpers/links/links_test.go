package links_test

import (
	"codeberg.org/konterfai/konterfai/pkg/helpers/links"
	"context"
	"net/url"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLinks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Links Suite")
}

var _ = Describe("Links", func() {
	var ctx context.Context
	totalTests := 1000
	url := url.URL{Scheme: "https", Host: "example.com"}
	BeforeEach(func() {
		ctx = context.Background()
	})

	Context("RandomLink", func() {
		It("should return the original input when subdierectories and variablescount are 0", func() {
			for i := 0; i < totalTests; i++ {
				link := links.RandomLink(ctx, url, 0, 0, 0)
				Expect(link).NotTo(BeEmpty())
				Expect(link).To(HavePrefix("https://example.com/"))
			}
		})

		It("should return a random link with a subdirectory and without variables", func() {
			for i := 0; i < totalTests; i++ {
				link := links.RandomLink(ctx, url, 1, 0, 1)
				Expect(link).NotTo(BeEmpty())
				Expect(link).To(HavePrefix("https://example.com/"))
				Expect(link).To(MatchRegexp(`https://example.com/.+`))
			}
		})

		It("should return a random link without subdirectories and with a variable", func() {
			for i := 0; i < totalTests; i++ {
				link := links.RandomLink(ctx, url, 0, 1, 1)
				Expect(link).NotTo(BeEmpty())
				Expect(link).To(HavePrefix("https://example.com/"))
				Expect(link).To(MatchRegexp(`https://example.com/\?.+=.+`))
			}
		})

		It("should return a random link with variables", func() {
			for i := 0; i < totalTests; i++ {
				link := links.RandomLink(ctx, url, 1, 1, 1)
				Expect(link).NotTo(BeEmpty())
				Expect(link).To(HavePrefix("https://example.com/"))
				Expect(link).To(MatchRegexp(`https://example.com/.+/?.+=.+`))
			}
		})
	})

	Context("RandomSimpleLink", func() {
		It("should return a random simple link", func() {
			for i := 0; i < totalTests; i++ {
				link := links.RandomSimpleLink(ctx, url)
				Expect(link).NotTo(BeEmpty())
				Expect(link).To(HavePrefix("https://example.com/"))
				Expect(link).To(MatchRegexp(`https://example.com/.+-.+`))
			}
		})
	})
})
