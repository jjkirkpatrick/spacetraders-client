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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type MetricsApp struct {
	client      *client.Client
	meter       metric.Meter
	creditGauge metric.Float64ObservableGauge
}

func newMetricsApp(ctx context.Context) (*MetricsApp, error) {
	options := client.DefaultClientOptions()
	options.Symbol = "METRICS_TEST1"
	options.Faction = "COSMIC"

	// Initialize telemetry with the new public options
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-metrics"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"
	options.TelemetryOptions.MetricInterval = 1 * time.Second

	spaceClient, err := client.NewClient(options)
	if err != nil {
		return nil, err
	}

	app := &MetricsApp{
		client: spaceClient,
		meter:  otel.GetMeterProvider().Meter("spacetraders-metrics"),
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

func (app *MetricsApp) run(ctx context.Context) error {
	app.client.Logger.Info("Starting metrics collection loop")

	// Get agent details and collect metrics
	for i := 0; i < 2000; i++ {
		app.client.Logger.Info("Starting iteration", "iteration", i)

		// No need for artificial delays - the rate limiter in the client will handle pacing
		agent, err := entities.GetAgent(app.client)
		if err != nil {
			app.client.Logger.Error("Failed to get agent", "iteration", i, "error", err)
			return err
		}
		app.client.Logger.Info("Retrieved agent", "iteration", i, "agent.symbol", agent.Symbol, "agent.credits", agent.Credits)

		// Setup credit gauge
		if err := app.setupCreditGaugeCallback(agent); err != nil {
			app.client.Logger.Error("Failed to set up credit gauge callback", "iteration", i, "error", err)
			return err
		}
		app.client.Logger.Info("Gauge callback set for agent", "iteration", i, "agent.symbol", agent.Symbol)

		app.client.Logger.Info("Iteration completed", "iteration", i)
	}
	app.client.Logger.Info("Finished metrics collection loop")

	return nil
}

func main() {
	ctx := context.Background()

	slog.Info("Initializing metrics application")
	app, err := newMetricsApp(ctx)
	if err != nil {
		slog.Error("Failed to create metrics app", "error", err)
		os.Exit(1)
	}
	app.client.Logger.Info("Metrics application initialized successfully")
	defer app.client.Close(ctx)

	app.client.Logger.Info("Waiting 2 seconds for collector connection")
	time.Sleep(2 * time.Second)

	app.client.Logger.Info("Setting up graceful shutdown handler")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	app.client.Logger.Info("Starting metrics collector goroutine")
	errChan := make(chan error, 1)
	go func() {
		app.client.Logger.Info("Metrics collector goroutine started")
		errChan <- app.run(ctx)
	}()

	app.client.Logger.Info("Entering main event loop, awaiting errors or shutdown signals")
	select {
	case err := <-errChan:
		if err != nil {
			app.client.Logger.Error("Application error", "error", err)
			os.Exit(1)
		}
		app.client.Logger.Info("Application completed successfully")
	case s := <-sigChan:
		app.client.Logger.Info("Received shutdown signal", "signal", s)
		app.client.Logger.Info("Shutting down gracefully...")
	}
}
