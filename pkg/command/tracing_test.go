package command_test

import (
	"context"
	"testing"

	"codeberg.org/konterfai/konterfai/pkg/command"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tracing", func() {
	var (
		ctx      context.Context
		endpoint string
		service  string
	)

	BeforeEach(func() {
		ctx = context.Background()
		endpoint = "http://localhost:4317"
		service = "konterfai"
	})

	It("should not return an error if endpoint is set properly", func() {
		err := command.SetTraceProvider(ctx, endpoint, service)
		Expect(err).To(BeNil())
	})

	It("should not return an error if endpoint is empty (tracing disabled)", func() {
		err := command.SetTraceProvider(ctx, "", service)
		Expect(err).To(BeNil())
	})
})
