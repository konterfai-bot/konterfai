package command

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.26.0"
)

// GetTraceProvider returns a new trace provider with the given endpoint and service name.
func GetTraceProvider(ctx context.Context, endpoint, serviceName string) *trace.TracerProvider {
	if endpoint == "" {
		fmt.Println("tracing is disabled")
		return trace.NewTracerProvider(
			trace.WithSampler(trace.NeverSample()),
		)
	}
	exporter, err := otlptrace.New(
		ctx,
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithHeaders(map[string]string{
				"Content-Type": "application/json",
			}),
			otlptracehttp.WithInsecure(),
		),
	)

	if err != nil {
		fmt.Printf("failed to create trace exporter: %v\n", err)
		return nil
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
	)
	otel.SetTracerProvider(tp)
	return tp
}
