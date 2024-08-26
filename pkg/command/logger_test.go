package command_test

import (
	"codeberg.org/konterfai/konterfai/pkg/command"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger", func() {
	When("Set logger is run", func() {
		It("Sets the logger when both arguments are empty", func() {
			logger, err := command.SetLogger("", "")
			Expect(err).To(BeNil())
			Expect(logger).NotTo(BeNil())
		})

		It("Sets the logger when format is text and level is info", func() {
			logger, err := command.SetLogger("text", "info")
			Expect(err).To(BeNil())
			Expect(logger).NotTo(BeNil())
		})

		It("Sets the logger when format is json and level is debug", func() {
			logger, err := command.SetLogger("json", "debug")
			Expect(err).To(BeNil())
			Expect(logger).NotTo(BeNil())
		})

		It("Sets the logger when format is off and level is warn", func() {
			logger, err := command.SetLogger("off", "warn")
			Expect(err).To(BeNil())
			Expect(logger).NotTo(BeNil())
		})

		It("Sets the logger when format is text and level is error", func() {
			logger, err := command.SetLogger("text", "error")
			Expect(err).To(BeNil())
			Expect(logger).NotTo(BeNil())
		})
	})
})
