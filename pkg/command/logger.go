package command

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("codeberg.org/konterfai/konterfai/pkg/command")

// SetLogger sets the logger for the command package.
func SetLogger(format, level string) (*slog.Logger, error) {
	// NOTE: When runnning in CI always set the logger to off, otherwise it will blow up woodpecker!
	//       If a test is failing, you have to debug it locally.
	_, span := tracer.Start(context.Background(), "SetLogger")
	defer span.End()

	var opts *slog.HandlerOptions
	switch strings.ToLower(level) {
	case "debug":
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	case "info":
		opts = &slog.HandlerOptions{Level: slog.LevelInfo}
	case "warn":
		opts = &slog.HandlerOptions{Level: slog.LevelWarn}
	case "error":
		fallthrough
	default:
		opts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}

	switch strings.ToLower(format) {
	case "off":
		return slog.New(slog.NewTextHandler(io.Discard, opts)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, opts)), nil
	case "text":
		fallthrough
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, opts)), nil
	}
}
