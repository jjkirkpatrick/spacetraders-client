package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := newMetricsApp(ctx)
	if err != nil {
		slog.Error("Failed to create app", "error", err)
		os.Exit(1)
	}
	defer app.client.Close(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	systems, err := entities.ListSystems(app.client)
	if err != nil {
		slog.Error("Failed to list systems", "error", err)
		os.Exit(1)
	}

	// Create a worker pool with 10 workers
	numWorkers := 10
	systemsChan := make(chan *entities.System, len(systems))
	errChan := make(chan error, len(systems))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case system, ok := <-systemsChan:
					if !ok {
						return
					}
					waypoints, _, err := system.ListWaypoints("", "")
					if err != nil {
						errChan <- fmt.Errorf("failed to list waypoints for system %s: %w", system.Symbol, err)
						continue
					}

					for _, waypoint := range waypoints {
						slog.Info("Waypoint", "system", system.Symbol, "waypoint.name", waypoint.Symbol)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	// Send systems to workers in a separate goroutine
	go func() {
		for _, system := range systems {
			select {
			case systemsChan <- system:
			case <-ctx.Done():
				return
			}
		}
		close(systemsChan)
	}()

	// Handle shutdown
	select {
	case err := <-errChan:
		slog.Error("Application error", "error", err)
		cancel()
	case <-sigChan:
		slog.Info("Shutting down gracefully...")
		cancel()
	}

	// Wait for all workers to complete
	wg.Wait()
	close(errChan)

	// Check for any remaining errors
	for err := range errChan {
		slog.Error("Error processing waypoints", "error", err)
	}
}
