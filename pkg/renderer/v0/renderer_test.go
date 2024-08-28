package renderer_test

import (
	"context"
	"log/slog"
	"testing"

	"codeberg.org/konterfai/konterfai/pkg/command"
	"codeberg.org/konterfai/konterfai/pkg/renderer/v0"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRenderer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Renderer Suite")
}

var _ = Describe("Renderer", func() {
	var (
		ctx    context.Context
		r      *renderer.Renderer
		rd     renderer.RenderData
		logger *slog.Logger
	)
	BeforeEach(func() {
		ctx = context.Background()
		logger, _ = command.SetLogger("off", "")
		r = renderer.NewRenderer(ctx, logger, []string{})
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
			Expect(renderer.NewRenderer(ctx, logger, []string{})).NotTo(BeNil())
		})
	})

	Context("RenderInRandomTemplate", func() {
		It("should render a random template", func() {
			renderedTemplate, err := r.RenderInRandomTemplate(ctx, rd)
			Expect(err).NotTo(HaveOccurred())
			Expect(renderedTemplate).NotTo(BeEmpty())
			Expect(renderedTemplate).NotTo(BeNil())
		})

		It("should render a random template if headlineLinks are provided from the renderer", func() {
			rd.HeadlineLinks = []string{}
			r = renderer.NewRenderer(ctx, logger, []string{"headLineLink0", "headLineLink1", "headLineLink2", "headLineLink3", "headLineLink4", "headLineLink5", "headLineLink6", "headLineLink7", "headLineLink8", "headLineLink9"})
			renderedTemplate, err := r.RenderInRandomTemplate(ctx, rd)
			Expect(err).NotTo(HaveOccurred())
			Expect(renderedTemplate).NotTo(BeEmpty())
			Expect(renderedTemplate).NotTo(BeNil())
		})

		It("should not render if headlineLinks is empty", func() {
			rd.HeadlineLinks = []string{}
			renderedTemplate, err := r.RenderInRandomTemplate(ctx, rd)
			Expect(err).To(HaveOccurred())
			Expect(renderedTemplate).To(Equal(""))
		})

		It("should not render if headlineLinks is nil", func() {
			rd.HeadlineLinks = nil
			renderedTemplate, err := r.RenderInRandomTemplate(ctx, rd)
			Expect(err).To(HaveOccurred())
			Expect(renderedTemplate).To(Equal(""))
		})

		It("should throw an error if the template is not found", func() {
			r.SetTemplates([]string{})
			renderedTemplate, err := r.RenderInRandomTemplate(ctx, rd)
			Expect(err).To(HaveOccurred())
			Expect(renderedTemplate).To(Equal(""))
		})

		It("should throw an error if the template contains invalid fields", func() {
			r.SetTemplates([]string{"{{.InvalidField}}"})
			renderedTemplate, err := r.RenderInRandomTemplate(ctx, rd)
			Expect(err).To(HaveOccurred())
			Expect(renderedTemplate).To(Equal(""))
		})
	})
})
