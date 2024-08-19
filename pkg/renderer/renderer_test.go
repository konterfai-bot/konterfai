package renderer_test

import (
	"context"
	"testing"

	"codeberg.org/konterfai/konterfai/pkg/renderer"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRenderer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Functions Suite")
}

var _ = Describe("Renderer", func() {
	var ctx context.Context
	var r *renderer.Renderer
	var rd renderer.RenderData
	BeforeEach(func() {
		ctx = context.Background()
		r = renderer.NewRenderer(ctx, []string{})
		rd = renderer.RenderData{
			NewsAnchor:    "newsAnchor",
			Headline:      "headline",
			Content:       "content",
			FollowUpLink:  "followUpLink",
			HeadlineLinks: []string{"headLineLink0", "headLineLink1", "headLineLink2", "headLineLink3", "headLineLink4", "headLineLink5", "headLineLink6", "headLineLink7", "headLineLink8", "headLineLink9"},
			RandomTopics:  []renderer.RandomTopic{{Topic: "topic", Link: "link"}},
			Year:          "year",
			CurrentYear:   "currentYear",
			MetaData:      renderer.MetaData{Description: "description", Keywords: "keywords", Charset: "charset"},
			LanguageCode:  "languageCode",
		}
	})

	Context("NewRenderer", func() {
		It("should return a new renderer", func() {
			Expect(renderer.NewRenderer(ctx, []string{})).NotTo(BeNil())
		})
	})

	Context("RenderInRandomTemplate", func() {
		It("should render a random template", func() {
			renderedTemplate, err := r.RenderInRandomTemplate(ctx, rd)
			Expect(err).NotTo(HaveOccurred())
			Expect(renderedTemplate).NotTo(BeEmpty())
			Expect(renderedTemplate).NotTo(BeNil())
		})
	})
})
