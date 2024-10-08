package command

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"codeberg.org/konterfai/konterfai/pkg/statisticsserver"
	"codeberg.org/konterfai/konterfai/pkg/webserver"
	"github.com/oklog/run"
	"github.com/urfave/cli/v2"
)

// Run is the entry point for running konterfAI.
func Run(c *cli.Context, logger *slog.Logger) error { //nolint: funlen
	ctx, cancel := func() (context.Context, context.CancelFunc) {
		return context.WithCancel(c.Context)
	}()
	defer cancel()

	err := SetTraceProvider(ctx, logger, c.String("tracing-endpoint"), "konterfai")
	if err != nil {
		return err
	}

	fmt.Println(generateHeader(c, true)) //nolint:forbidigo
	syncer := make(chan error)
	st := statistics.NewStatistics(ctx, logger, generateHeader(c, false))
	hcURL, err := url.Parse(c.String("hallucinator-url"))
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("could not parse hallucinator-url (%v)", err))

		return err
	}
	hal := hallucinator.NewHallucinator(ctx, logger, c.Duration("generate-interval"),
		c.Int("hallucination-cache-size"), c.Int("hallucination-prompt-word-count"),
		c.Int("hallucination-request-count"), c.Int("hallucination-minimal-length"),
		c.Int("hallucination-word-count"), c.Int("hallucinator-link-percentage"),
		c.Int("hallucinator-link-max-subdirectory-depth"),
		c.Float64("hallucinator-link-has-variables-probability"), c.Int("hallucinator-link-max-variables"),
		*hcURL, c.String("ollama-address"), c.String("ollama-model"),
		c.Duration("ollama-request-timeout"), c.Float64("ai-temperature"), c.Int("ai-seed"), st)
	gr := run.Group{}
	gr.Add(func() error {
		select {
		case <-ctx.Done():
			return nil
		case syncer <- hal.Start(ctx):
			return <-syncer
		}
	}, func(_ error) {
		logger.InfoContext(ctx, "shutting down hallucinator")
		cancel()
	})
	gr.Add(func() error {
		ws := webserver.NewWebServer(ctx, logger, c.String("address"), c.Int("port"), hal, st, *hcURL,
			c.Float64("webserver-200-probability"), c.Float64("random-uncertainty"),
			c.Int("webserver-error-cache-size"))
		select {
		case <-ctx.Done():
			return nil
		case syncer <- ws.Serve(ctx):
			return <-syncer
		}
	}, func(_ error) {
		logger.InfoContext(ctx, "shutting down webserver")
		cancel()
	})
	gr.Add(func() error {
		ss := statisticsserver.NewStatisticsServer(ctx, logger, c.String("address"), c.Int("statistics-port"), st)
		select {
		case <-ctx.Done():
			return nil
		case syncer <- ss.Serve(ctx):
			return <-syncer
		}
	}, func(_ error) {
		logger.InfoContext(ctx, "shutting down statistics server")
		cancel()
	})

	return gr.Run()
}

// gernerateHeader prints the header of the konterfAI cli command.
func generateHeader(c *cli.Context, withHeadline bool) string {
	var header string
	if withHeadline {
		header += strings.Join([]string{
			fmt.Sprintln("#############################"),
			fmt.Sprintln("# konterfAI - the anti-AI-AI #"),
			fmt.Sprintln("#############################"),
			fmt.Sprintln(),
			fmt.Sprintln("Configuration:"),
		}, "")
	}
	header += strings.Join([]string{
		fmt.Sprintln("\t- Address: \t\t\t\t", c.String("address")),
		fmt.Sprintln("\t- Port: \t\t\t\t", c.Int("port")),
		fmt.Sprintln("\t- Statistics Port: \t\t\t", c.Int("statistics-port")),
		fmt.Sprintln("\t- Generate Interval: \t\t\t", c.Duration("generate-interval")),
		fmt.Sprintln("\t- Hallucination Cache Size: \t\t", c.Int("hallucination-cache-size")),
		fmt.Sprintln("\t- Hallucination Prompt Word Count: \t", c.Int("hallucination-prompt-word-count")),
		fmt.Sprintln("\t- Hallucination Word Count: \t\t", c.Int("hallucination-word-count")),
		fmt.Sprintln("\t- Hallucination Request Count:  \t", c.Int("hallucination-request-count")),
		fmt.Sprintln("\t- Ollama Address: \t\t\t", c.String("ollama-address")),
		fmt.Sprintln("\t- Ollama Model: \t\t\t", c.String("ollama-model")),
		fmt.Sprintln("\t- AI Temperature: \t\t\t", c.Float64("ai-temperature")),
		fmt.Sprintln("\t- AI Seed: \t\t\t\t", c.Int("ai-seed")),
		fmt.Sprintln("\t- Hallucinator URL: \t\t\t", c.String("hallucinator-url")),
		fmt.Sprintln("\t- Log Level: \t\t\t\t", c.String("log-level")),
		fmt.Sprintln("\t- Log Format: \t\t\t\t", c.String("log-format")),
	}, "")

	if c.String("tracing-endpoint") != "" {
		header += strings.Join([]string{
			fmt.Sprintln("\t- Tracing Endpoint: \t\t\t", c.String("tracing-endpoint")),
		}, "")
	}

	header += strings.Join([]string{
		fmt.Sprintln(),
	},
		"",
	)

	return header
}
