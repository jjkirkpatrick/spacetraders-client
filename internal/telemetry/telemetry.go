package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds the configuration for OpenTelemetry setup
type Config struct {
	ServiceName     string
	ServiceVersion  string
	Environment     string
	OTLPEndpoint    string
	MetricInterval  time.Duration
	AdditionalAttrs []attribute.KeyValue
	GRPCDialOptions []grpc.DialOption
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		Environment:    "development",
		MetricInterval: 1 * time.Second,
		GRPCDialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		},
	}
}

// Providers holds the initialized OpenTelemetry providers
type Providers struct {
	MeterProvider *metric.MeterProvider
	Resource      *resource.Resource
}

// InitTelemetry initializes OpenTelemetry with the provided configuration
func InitTelemetry(ctx context.Context, cfg Config) (*Providers, error) {
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// Create resource with service information and additional attributes
	attrs := append([]attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
		semconv.ServiceVersion(cfg.ServiceVersion),
		semconv.DeploymentEnvironment(cfg.Environment),
	}, cfg.AdditionalAttrs...)

	res, err := resource.New(ctx,
		resource.WithAttributes(attrs...),
		resource.WithContainer(),
		resource.WithHost(),
		resource.WithOS(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize gRPC connection
	dialOpts := cfg.GRPCDialOptions
	if len(dialOpts) == 0 {
		dialOpts = DefaultConfig().GRPCDialOptions
	}
	conn, err := grpc.DialContext(ctx, cfg.OTLPEndpoint, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	// Initialize metric exporter
	metricExp, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// Create MeterProvider
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(
			metric.NewPeriodicReader(
				metricExp,
				metric.WithInterval(cfg.MetricInterval),
			),
		),
	)

	// Set global provider
	otel.SetMeterProvider(mp)

	return &Providers{
		MeterProvider: mp,
		Resource:      res,
	}, nil
}

// Shutdown gracefully shuts down the OpenTelemetry providers
func (p *Providers) Shutdown(ctx context.Context) error {
	if err := p.MeterProvider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown meter provider: %w", err)
	}
	return nil
}

// Example usage:
/*
func main() {
    ctx := context.Background()

    cfg := telemetry.DefaultConfig()
    cfg.ServiceName = "my-service"
    cfg.ServiceVersion = "1.0.0"
    cfg.OTLPEndpoint = "localhost:4317"

    providers, err := telemetry.InitTelemetry(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer providers.Shutdown(ctx)

    // Get a logger
    logger := providers.LoggerProvider.Logger("my-service")

    // Create a span
    tracer := otel.GetTracerProvider().Tracer("my-service")
    ctx, span := tracer.Start(ctx, "my-operation")
    defer span.End()

    // Log with trace context
    logger.Info(ctx, "Operation in progress",
        sdklogs.String("component", "main"),
        sdklogs.Int("attempt", 1),
    )
}
*/
