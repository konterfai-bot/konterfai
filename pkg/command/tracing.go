package command

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// SetTraceProvider returns a new trace provider with the given endpoint and service name.
func SetTraceProvider(ctx context.Context, endpoint, serviceName string) {
	if endpoint == "" {
		fmt.Println("tracing is disabled")
		otel.SetTracerProvider(trace.NewTracerProvider(
			trace.WithSampler(trace.NeverSample()),
		))
		return
	}

	conn, err := grpc.NewClient(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)

	if err != nil {
		fmt.Printf("failed to create grpc connection: %v\n", err)
		return
	}

	exporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithGRPCConn(conn),
	)

	if err != nil {
		fmt.Printf("failed to create trace exporter: %v\n", err)
		return
	}

	resources, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)

	if err != nil {
		fmt.Printf("failed to create resource: %v\n", err)
		return
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)
}
