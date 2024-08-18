package command

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"
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

			&cli.IntFlag{
				Name:        "statistics-port",
				Usage:       "The port to listen on for statistics.",
				Value:       8081,
				DefaultText: "8081",
			},
			&cli.StringFlag{
				Name:        "hallucinator-url",
				Usage:       "The FQDN of the hallicinator, e.g. http://localhost:8080",
				Value:       "http://localhost:8080",
				DefaultText: "http://localhost:8080",
			},
			&cli.DurationFlag{
				Name:        "generate-interval",
				Usage:       "The interval in seconds to generate a new hallucination",
				Value:       5 * time.Second,
				DefaultText: "5",
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
