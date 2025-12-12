package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/jjkirkpatrick/spacetraders-client/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type MetricsApp struct {
	client      *client.Client
	meter       metric.Meter
	tracer      trace.Tracer
	creditGauge metric.Float64ObservableGauge
}

func newMetricsApp(ctx context.Context) (*MetricsApp, error) {
	options := client.DefaultClientOptions()
	options.Symbol = "METRICS-DEMO"
	options.Faction = "COSMIC"

	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-metrics"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"
	options.TelemetryOptions.MetricInterval = 1 * time.Second

	spaceClient, err := client.NewClient(options)
	if err != nil {
		return nil, err
	}

	// Set up combined logging
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-metrics", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	app := &MetricsApp{
		client: spaceClient,
		meter:  otel.GetMeterProvider().Meter("spacetraders-metrics"),
		tracer: otel.GetTracerProvider().Tracer("spacetraders-metrics"),
	}

	app.creditGauge, err = app.meter.Float64ObservableGauge("agent_credits",
		metric.WithDescription("Current credit balance of the agent"),
		metric.WithUnit("credits"),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *MetricsApp) setupCreditGaugeCallback(agent *entities.Agent) error {
	_, err := app.meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
		o.ObserveFloat64(app.creditGauge, float64(agent.Credits),
			metric.WithAttributes(
				attribute.String("agent", agent.Symbol),
			),
		)
		return nil
	}, app.creditGauge)
	return err
}

func (app *MetricsApp) run(ctx context.Context, iterations int) error {
	ctx, span := app.tracer.Start(ctx, "metrics_collection_loop")
	defer span.End()

	slog.InfoContext(ctx, "Starting metrics collection loop", "iterations", iterations)

	for i := 0; i < iterations; i++ {
		ctx, iterSpan := app.tracer.Start(ctx, "iteration")
		iterSpan.SetAttributes(attribute.Int("iteration", i+1))

		slog.InfoContext(ctx, "Starting iteration", "iteration", i+1, "total", iterations)

		agent, err := entities.GetAgent(app.client)
		if err != nil {
			iterSpan.RecordError(err)
			slog.ErrorContext(ctx, "Failed to get agent",
				"iteration", i+1,
				"error", err,
			)
			iterSpan.End()
			return err
		}

		slog.InfoContext(ctx, "Agent data retrieved",
			"iteration", i+1,
			"symbol", agent.Symbol,
			"credits", agent.Credits,
		)

		if err := app.setupCreditGaugeCallback(agent); err != nil {
			iterSpan.RecordError(err)
			slog.ErrorContext(ctx, "Failed to setup credit gauge callback",
				"iteration", i+1,
				"error", err,
			)
			iterSpan.End()
			return err
		}

		slog.InfoContext(ctx, "Credit gauge updated",
			"iteration", i+1,
			"agent", agent.Symbol,
			"credits", agent.Credits,
		)

		iterSpan.End()
	}

	slog.InfoContext(ctx, "Metrics collection loop completed", "iterations", iterations)
	return nil
}

func main() {
	ctx := context.Background()

	slog.Info("Initializing metrics application")

	app, err := newMetricsApp(ctx)
	if err != nil {
		slog.Error("Failed to create metrics application", "error", err)
		os.Exit(1)
	}
	defer app.client.Close(ctx)

	// Create root span
	ctx, rootSpan := app.tracer.Start(ctx, "metrics_demo")
	defer rootSpan.End()

	slog.InfoContext(ctx, "Metrics application initialized")

	// Wait for collector connection
	slog.InfoContext(ctx, "Waiting for OTLP collector connection", "delay_seconds", 2)
	time.Sleep(2 * time.Second)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	slog.InfoContext(ctx, "Starting metrics collection")

	errChan := make(chan error, 1)
	go func() {
		errChan <- app.run(ctx, 20) // Run 20 iterations
	}()

	select {
	case err := <-errChan:
		if err != nil {
			slog.ErrorContext(ctx, "Application error", "error", err)
			os.Exit(1)
		}
		slog.InfoContext(ctx, "Application completed successfully")
	case sig := <-sigChan:
		slog.InfoContext(ctx, "Shutdown signal received", "signal", sig.String())
	}

	slog.InfoContext(ctx, "Metrics demo complete")
}
