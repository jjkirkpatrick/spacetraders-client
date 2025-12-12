package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Config holds the configuration for OpenTelemetry setup
type Config struct {
	// ServiceName identifies your service in telemetry data (required)
	ServiceName string

	// ServiceVersion is the version of your service
	ServiceVersion string

	// Environment describes the deployment environment (e.g., "development", "production")
	Environment string

	// OTLPEndpoint is the gRPC endpoint for your OpenTelemetry collector (required)
	// Example: "localhost:4317" or "otel-collector.example.com:4317"
	OTLPEndpoint string

	// MetricInterval controls how frequently metrics are exported
	MetricInterval time.Duration

	// TraceSampleRate controls the fraction of traces to sample (0.0 to 1.0)
	// 1.0 means sample all traces, 0.1 means sample 10% of traces
	TraceSampleRate float64

	// EnableMetrics enables metric collection (default: true)
	EnableMetrics bool

	// EnableTracing enables distributed tracing (default: true)
	EnableTracing bool

	// EnableLogging enables log export to OTLP (default: true)
	EnableLogging bool

	// AdditionalAttrs are custom resource attributes added to all telemetry
	AdditionalAttrs []attribute.KeyValue

	// GRPCDialOptions allows customization of the gRPC connection
	GRPCDialOptions []grpc.DialOption
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() Config {
	return Config{
		Environment:     "development",
		MetricInterval:  5 * time.Second,
		TraceSampleRate: 1.0, // Sample all traces by default
		EnableMetrics:   true,
		EnableTracing:   true,
		EnableLogging:   true,
		GRPCDialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	}
}

// Providers holds the initialized OpenTelemetry providers
type Providers struct {
	MeterProvider  *sdkmetric.MeterProvider
	TracerProvider *sdktrace.TracerProvider
	LoggerProvider *sdklog.LoggerProvider
	Resource       *resource.Resource

	// Internal: gRPC connection for cleanup
	conn *grpc.ClientConn
}

// InitTelemetry initializes OpenTelemetry with the provided configuration.
// It sets up metrics, tracing, and logging exporters based on the config.
func InitTelemetry(ctx context.Context, cfg Config) (*Providers, error) {
	if cfg.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}
	if cfg.OTLPEndpoint == "" {
		return nil, fmt.Errorf("OTLP endpoint is required")
	}

	// Create resource with service information
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

	// Initialize gRPC connection with timeout
	dialOpts := cfg.GRPCDialOptions
	if len(dialOpts) == 0 {
		dialOpts = DefaultConfig().GRPCDialOptions
	}

	// Add a timeout context for dialing
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, cfg.OTLPEndpoint, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to %s: %w", cfg.OTLPEndpoint, err)
	}

	providers := &Providers{
		Resource: res,
		conn:     conn,
	}

	// Initialize metrics
	if cfg.EnableMetrics {
		metricExp, err := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithGRPCConn(conn),
		)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to create metric exporter: %w", err)
		}

		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(
				sdkmetric.NewPeriodicReader(
					metricExp,
					sdkmetric.WithInterval(cfg.MetricInterval),
				),
			),
		)
		providers.MeterProvider = mp
		otel.SetMeterProvider(mp)
	}

	// Initialize tracing
	if cfg.EnableTracing {
		traceExp, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithGRPCConn(conn),
		)
		if err != nil {
			providers.shutdownPartial(ctx)
			return nil, fmt.Errorf("failed to create trace exporter: %w", err)
		}

		// Configure sampler based on sample rate
		var sampler sdktrace.Sampler
		if cfg.TraceSampleRate >= 1.0 {
			sampler = sdktrace.AlwaysSample()
		} else if cfg.TraceSampleRate <= 0.0 {
			sampler = sdktrace.NeverSample()
		} else {
			sampler = sdktrace.TraceIDRatioBased(cfg.TraceSampleRate)
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
			sdktrace.WithBatcher(traceExp),
			sdktrace.WithSampler(sampler),
		)
		providers.TracerProvider = tp
		otel.SetTracerProvider(tp)

		// Set up trace context propagation
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	}

	// Initialize logging
	if cfg.EnableLogging {
		logExp, err := otlploggrpc.New(ctx,
			otlploggrpc.WithGRPCConn(conn),
		)
		if err != nil {
			providers.shutdownPartial(ctx)
			return nil, fmt.Errorf("failed to create log exporter: %w", err)
		}

		lp := sdklog.NewLoggerProvider(
			sdklog.WithResource(res),
			sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
		)
		providers.LoggerProvider = lp
		global.SetLoggerProvider(lp)
	}

	return providers, nil
}

// shutdownPartial shuts down any initialized providers (used during init errors)
func (p *Providers) shutdownPartial(ctx context.Context) {
	if p.MeterProvider != nil {
		p.MeterProvider.Shutdown(ctx)
	}
	if p.TracerProvider != nil {
		p.TracerProvider.Shutdown(ctx)
	}
	if p.LoggerProvider != nil {
		p.LoggerProvider.Shutdown(ctx)
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

// Shutdown gracefully shuts down all OpenTelemetry providers.
// Call this when your application terminates to ensure all telemetry is flushed.
func (p *Providers) Shutdown(ctx context.Context) error {
	var errs []error

	if p.MeterProvider != nil {
		if err := p.MeterProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
		}
	}

	if p.TracerProvider != nil {
		if err := p.TracerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("tracer provider shutdown: %w", err))
		}
	}

	if p.LoggerProvider != nil {
		if err := p.LoggerProvider.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("logger provider shutdown: %w", err))
		}
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("gRPC connection close: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}

// GetTracer returns a tracer from the TracerProvider.
// Returns nil if tracing is not enabled.
func (p *Providers) GetTracer(name string) interface{} {
	if p.TracerProvider == nil {
		return nil
	}
	return p.TracerProvider.Tracer(name)
}

// GetMeter returns a meter from the MeterProvider.
// Returns nil if metrics are not enabled.
func (p *Providers) GetMeter(name string) interface{} {
	if p.MeterProvider == nil {
		return nil
	}
	return p.MeterProvider.Meter(name)
}
