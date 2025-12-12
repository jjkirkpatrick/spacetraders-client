package main

import (
	"context"
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
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	ctx := context.Background()

	// Create client with telemetry
	options := client.DefaultClientOptions()
	options.Symbol = "RESET-DEMO"
	options.Faction = "COSMIC"
	options.RequestQueueSize = 100
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-reset-handler"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"

	c, err := client.NewClient(options)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}
	defer c.Close(ctx)

	// Set up combined logging
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-reset-handler", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	// Get tracer
	tracer = otel.GetTracerProvider().Tracer("spacetraders-reset-handler")

	// Create root span
	ctx, rootSpan := tracer.Start(ctx, "game_session")
	defer rootSpan.End()

	// Create cancellable context
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		slog.InfoContext(ctx, "Shutdown signal received", "signal", sig.String())
		cancel()
	}()

	// Start monitoring goroutines
	var wg sync.WaitGroup

	wg.Add(1)
	go monitorGameReset(ctx, c, &wg)

	wg.Add(1)
	go makeRequests(ctx, c, &wg)

	slog.InfoContext(ctx, "Game reset handler running, press Ctrl+C to exit")
	wg.Wait()

	slog.InfoContext(ctx, "Application shutdown complete")
}

func monitorGameReset(ctx context.Context, c *client.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, span := tracer.Start(ctx, "game_reset_monitor")
	defer span.End()

	slog.InfoContext(ctx, "Game reset monitor started")

	for {
		resetDetected := c.WaitForGameReset(ctx)

		if !resetDetected {
			slog.InfoContext(ctx, "Game reset monitor shutting down")
			return
		}

		slog.ErrorContext(ctx, "Game reset detected",
			"action", "re-registration required",
			"status", "exiting",
		)

		os.Exit(1)
	}
}

func makeRequests(ctx context.Context, c *client.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, span := tracer.Start(ctx, "request_loop")
	defer span.End()

	slog.InfoContext(ctx, "Request loop started")

	requestCount := 0
	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "Request loop shutting down", "total_requests", requestCount)
			return
		default:
			if c.IsGameReset() {
				slog.WarnContext(ctx, "Game reset detected in request loop, stopping")
				return
			}

			requestCount++

			ctx, reqSpan := tracer.Start(ctx, "get_agent_request")
			agent, err := entities.GetAgent(c)
			if err != nil {
				reqSpan.RecordError(err)
				slog.ErrorContext(ctx, "Failed to get agent",
					"request_num", requestCount,
					"error", err,
				)
				reqSpan.End()
			} else {
				slog.InfoContext(ctx, "Agent status",
					"request_num", requestCount,
					"symbol", agent.Symbol,
					"credits", agent.Credits,
				)
				reqSpan.End()
			}

			time.Sleep(2 * time.Second)
		}
	}
}
