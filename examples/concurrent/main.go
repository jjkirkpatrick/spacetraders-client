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
	options.Symbol = "CONCURRENT-DEMO"
	options.Faction = "COSMIC"
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-concurrent"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"
	options.TelemetryOptions.MetricInterval = 1 * time.Second

	spaceClient, err := client.NewClient(options)
	if err != nil {
		return nil, err
	}

	// Set up combined logging
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-concurrent", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	app := &MetricsApp{
		client: spaceClient,
		meter:  otel.GetMeterProvider().Meter("spacetraders-concurrent"),
		tracer: otel.GetTracerProvider().Tracer("spacetraders-concurrent"),
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
		slog.Error("Failed to create application", "error", err)
		os.Exit(1)
	}
	defer app.client.Close(ctx)

	// Create root span
	ctx, rootSpan := app.tracer.Start(ctx, "concurrent_processing")
	defer rootSpan.End()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Fetch systems
	ctx, fetchSpan := app.tracer.Start(ctx, "fetch_systems")
	slog.InfoContext(ctx, "Fetching systems for concurrent processing")
	systems, err := entities.ListSystems(app.client)
	if err != nil {
		fetchSpan.RecordError(err)
		slog.ErrorContext(ctx, "Failed to list systems", "error", err)
		os.Exit(1)
	}
	fetchSpan.SetAttributes(attribute.Int("systems.count", len(systems)))
	slog.InfoContext(ctx, "Systems loaded", "count", len(systems))
	fetchSpan.End()

	// Create worker pool
	numWorkers := 10
	systemsChan := make(chan *entities.System, len(systems))
	errChan := make(chan error, len(systems))
	var wg sync.WaitGroup

	ctx, workerSpan := app.tracer.Start(ctx, "process_waypoints")
	workerSpan.SetAttributes(attribute.Int("workers", numWorkers))
	slog.InfoContext(ctx, "Starting worker pool", "workers", numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
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
						slog.ErrorContext(ctx, "Worker failed to fetch waypoints",
							"worker_id", workerID,
							"system", system.Symbol,
							"error", err,
						)
						continue
					}

					slog.InfoContext(ctx, "System waypoints fetched",
						"worker_id", workerID,
						"system", system.Symbol,
						"waypoint_count", len(waypoints),
					)

					for _, waypoint := range waypoints {
						slog.DebugContext(ctx, "Waypoint discovered",
							"system", system.Symbol,
							"waypoint", waypoint.Symbol,
							"type", waypoint.Type,
						)
					}
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	// Send systems to workers
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
		slog.ErrorContext(ctx, "Worker error occurred", "error", err)
		cancel()
	case <-sigChan:
		slog.InfoContext(ctx, "Shutdown signal received, stopping workers...")
		cancel()
	}

	// Wait for workers
	wg.Wait()
	workerSpan.End()
	close(errChan)

	// Log remaining errors
	errorCount := 0
	for err := range errChan {
		errorCount++
		if errorCount <= 5 {
			slog.ErrorContext(ctx, "Processing error", "error", err)
		}
	}
	if errorCount > 5 {
		slog.WarnContext(ctx, "Additional errors occurred", "count", errorCount-5)
	}

	slog.InfoContext(ctx, "Concurrent processing complete")
}
