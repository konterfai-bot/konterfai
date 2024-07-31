package main

import (
	"fmt"
	"os"

	"codeberg.org/konterfai/konterfai/pkg/command"
)

func main() {
	if err := command.Initialize(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
