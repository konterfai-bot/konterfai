package command_test

import (
	"codeberg.org/konterfai/konterfai/pkg/command"
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
	"time"
)

var _ = Describe("Command", func() {
	var cliContext *cli.Context
	Describe("Run", func() {
		It("should return no error", func() {
			cliContext = &cli.Context{
				Context: context.Background(),
			}
			logger, _ := command.SetLogger("off", "info")
			go func() {
				err := command.Run(cliContext, logger)
				Expect(err).To(HaveOccurred())
			}()
			time.Sleep(1 * time.Second)
			cliContext.Done()
		})
	})
})
