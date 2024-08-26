package command_test

import (
	"codeberg.org/konterfai/konterfai/pkg/command"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Root", func() {
	Describe("Initialize", func() {
		It("should return no error", func() {
			err := command.Initialize()
			// Usually, you would use Expect(err).To(BeNil()) here.
			// However go will pass a -test.* flag here which is parsed globally and will end up in urfave/cli.
			// According to https://github.com/golang/go/issues/46869 this is expected behavior.
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("flag provided but not defined: -test."))
		})
	})
})
