package command

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"codeberg.org/konterfai/konterfai/pkg/hallucinator"
	"codeberg.org/konterfai/konterfai/pkg/statistics"
	"codeberg.org/konterfai/konterfai/pkg/statisticsserver"
	"codeberg.org/konterfai/konterfai/pkg/webserver"
	"github.com/oklog/run"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel"
)

// Initialize is the entry point for initializing the konterfAI cli command.
func Initialize() error {
	app := &cli.App{
		Name:  "konterfai",
		Usage: "Run konterfAI the anti-AI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Usage:       "The address to listen on",
				Value:       "127.0.0.1",
				DefaultText: "127.0.0.1",
			},
			&cli.IntFlag{
				Name:        "port",
				Usage:       "The port to listen on",
				Value:       8080,
				DefaultText: "8080",
			},
			&cli.StringFlag{
				Name:        "hallucinator-url",
				Usage:       "The FQDN of the hallicinator, e.g. http://localhost:8080",
				Value:       "http://localhost:8080",
				DefaultText: "http://localhost:8080",
			},
			&cli.IntFlag{
				Name:        "statistics-port",
				Usage:       "The port to listen on for statistics.",
				Value:       8081,
				DefaultText: "8081",
			},
			&cli.DurationFlag{
				Name:        "generate-interval",
				Usage:       "The interval in seconds to generate a new hallucination",
				Value:       2 * time.Second,
				DefaultText: "2",
			},
			&cli.IntFlag{
				Name: "hallucination-cache-size",
				Usage: "The number of hallucinations to cache." +
					" Use high numbers for slow CPUs/GPUs and low" +
					" numbers if you have vast amount of CPU-/GPU-time to spare.",
				Value:       10,
				DefaultText: "10",
			},
			&cli.IntFlag{
				Name:        "hallucination-prompt-word-count",
				Usage:       "The number of words to use for hallucination prompts",
				Value:       5,
				DefaultText: "5",
			},
			&cli.IntFlag{
				Name:        "hallucination-word-count",
				Usage:       "The number of words to use for hallucinations.",
				Value:       500,
				DefaultText: "500",
			},
			&cli.IntFlag{
				Name: "hallucination-request-count",
				Usage: "Counter how many times the same hallucination should be presented." +
					" Use a high number here to reduce CPU-/GPU-load.",
				Value:       5,
				DefaultText: "5",
			},
			&cli.IntFlag{
				Name:        "hallucinator-link-percentage",
				Usage:       "The percentage of links to add to the hallucination measured by total words.",
				Value:       10,
				DefaultText: "10",
			},
			&cli.IntFlag{
				Name:        "hallucinator-link-max-subdirectory-depth",
				Usage:       "The maximum number of subdirectories for a link in the hallucination.",
				Value:       5,
				DefaultText: "5",
			},
			&cli.Float64Flag{
				Name:        "hallucinator-link-has-variables-probability",
				Usage:       "The probability of a link having variables.",
				Value:       0.5,
				DefaultText: "0.5",
			},
			&cli.IntFlag{
				Name:        "hallucinator-link-max-variables",
				Usage:       "The maximum number of variables for a link in the hallucination.",
				Value:       5,
				DefaultText: "5",
			},
			&cli.StringFlag{
				Name:        "ollama-address",
				Usage:       "The address of the ollama service",
				Value:       "http://localhost:11434",
				DefaultText: "http://localhost:11434",
			},
			&cli.StringFlag{
				Name:        "ollama-model",
				Usage:       "The model to use for hallucinations",
				Value:       "qwen2:0.5b",
				DefaultText: "qwen2:0.5b",
			},
			&cli.DurationFlag{
				Name:        "ollama-request-timeout",
				Usage:       "The timeout for the ollama service",
				Value:       60 * time.Second,
				DefaultText: "60s",
			},
			&cli.Float64Flag{
				Name: "ai-temperature",
				Usage: "The temperature for the AI. Use a high number for more randomness" +
					" and a low number for more coherence.",
				Value:       30.0,
				DefaultText: "30.0",
			},
			&cli.IntFlag{
				Name:        "ai-seed",
				Usage:       "The seed for the AI",
				Value:       0,
				DefaultText: "0",
			},
			&cli.Float64Flag{
				Name:        "webserver-200-probability",
				Usage:       "The probability of returning a 200 status code for a request.",
				Value:       0.95,
				DefaultText: "0.95",
			},
			&cli.IntFlag{
				Name:        "webserver-error-cache-size",
				Usage:       "The number of error codes to cache.",
				Value:       1000,
				DefaultText: "1000",
			},
			&cli.Float64Flag{
				Name:  "random-uncertainty",
				Usage: "The uncertainty for the random generator. Use a high number for more randomness",
				Value: 0.1,
			},
			&cli.StringFlag{
				Name:  "tracing-endpoint",
				Usage: "The endpoint for the tracing server (open telemetry). If empty, tracing is disabled.",
				Value: "",
			},
		},
		Action: func(c *cli.Context) error {
			return Run(c)
		},
	}
	return app.Run(os.Args)
}

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/command")

// Run is the entry point for running konterfAI.
func Run(c *cli.Context) error {
	ctx, cancel := func() (context.Context, context.CancelFunc) {
		return context.WithCancel(c.Context)
	}()
	defer cancel()

	SetTraceProvider(ctx, c.String("tracing-endpoint"), "konterfai")

	fmt.Println(generateHeader(c, true))
	syncer := make(chan error)

	st := statistics.NewStatistics(ctx, generateHeader(c, false))

	hcUrl, err := url.Parse(c.String("hallucinator-url"))
	if err != nil {
		fmt.Println("could not parse hallucinator-url")
		return err
	}

	hal := hallucinator.NewHallucinator(
		c.Duration("generate-interval"),
		c.Int("hallucination-cache-size"),
		c.Int("hallucination-prompt-word-count"),
		c.Int("hallucination-request-count"),
		c.Int("hallucination-word-count"),
		c.Int("hallucinator-link-percentage"),
		c.Int("hallucinator-link-max-subdirectory-depth"),
		c.Float64("hallucinator-link-has-variables-probability"),
		c.Int("hallucinator-link-max-variables"),
		*hcUrl,
		c.String("ollama-address"),
		c.String("ollama-model"),
		c.Duration("ollama-request-timeout"),
		c.Float64("ai-temperature"),
		c.Int("ai-seed"),
		st,
	)

	gr := run.Group{}
	gr.Add(func() error {
		select {
		case <-ctx.Done():
			return nil
		case syncer <- hal.Start(ctx):
			return <-syncer
		}
	}, func(_ error) {
		fmt.Println("shutting down hallucinator")
		cancel()
	})

	gr.Add(func() error {
		ws := webserver.NewWebServer(c.String("address"),
			c.Int("port"),
			hal,
			st,
			*hcUrl,
			c.Float64("webserver-200-probability"),
			c.Float64("random-uncertainty"),
			c.Int("webserver-error-cache-size"),
		)
		select {
		case <-ctx.Done():
			return nil
		case syncer <- ws.Serve():
			return <-syncer
		}
	}, func(_ error) {
		fmt.Println("shutting down webserver")
		cancel()
	})

	gr.Add(func() error {
		ss := statisticsserver.NewStatisticsServer(c.String("address"),
			c.Int("statistics-port"),
			st,
		)
		select {
		case <-ctx.Done():
			return nil
		case syncer <- ss.Serve():
			return <-syncer
		}
	}, func(_ error) {
		fmt.Println("shutting down statistics server")
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
