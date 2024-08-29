package renderer_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"codeberg.org/konterfai/konterfai/pkg/command"
	"codeberg.org/konterfai/konterfai/pkg/renderer/v1"
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
		logger *slog.Logger
		rd     *renderer.Renderer
	)

	BeforeEach(func() {
		ctx = context.Background()
		logger, _ = command.SetLogger("off", "")
		rd = renderer.NewRenderer(ctx, logger, []string{})
	})

	Context("NewRenderer", func() {
		It("should return a new Renderer", func() {
			Expect(rd).NotTo(BeNil())
			Expect(rd).To(BeAssignableToTypeOf(&renderer.Renderer{}))
		})
	})

	Context("RenderRandomSite", func() {
		It("should render a random site", func() {
			data, err := rd.RenderRandomSite(ctx)
			fmt.Println(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(data).NotTo(BeEmpty())
		})
	})
})
