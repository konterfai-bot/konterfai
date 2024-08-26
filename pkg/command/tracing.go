package command

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SetTraceProvider returns a new trace provider with the given endpoint and service name.
func SetTraceProvider(ctx context.Context, logger *slog.Logger, endpoint, serviceName string) error {
	if endpoint == "" {
		logger.InfoContext(ctx, "tracing is disabled")
		otel.SetTracerProvider(trace.NewTracerProvider(
			trace.WithSampler(trace.NeverSample()),
		))

		return nil
	}

	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("failed to create grpc connection (%v)", err))

		return err
	}

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("failed to create trace exporter (%v)", err))

		return err
	}

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		logger.ErrorContext(ctx, fmt.Sprintf("failed to create resource (%v)", err))

		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)

	return nil
}
