package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/jjkirkpatrick/spacetraders-client/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	ctx := context.Background()

	// Create client with telemetry
	options := client.DefaultClientOptions()
	options.Symbol = "QUEUE-DEMO"
	options.Faction = "COSMIC"
	options.RequestQueueSize = 100
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-queue-demo"
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
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-queue-demo", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	// Get tracer
	tracer := otel.GetTracerProvider().Tracer("spacetraders-queue-demo")

	// Create root span
	ctx, rootSpan := tracer.Start(ctx, "request_queue_demo")
	defer rootSpan.End()

	// Configuration
	requestCount := 10

	slog.InfoContext(ctx, "Starting request queue demonstration",
		"concurrent_requests", requestCount,
		"queue_size", options.RequestQueueSize,
	)

	var wg sync.WaitGroup
	startTime := time.Now()

	// Make multiple concurrent requests
	ctx, concurrentSpan := tracer.Start(ctx, "concurrent_requests")
	concurrentSpan.SetAttributes(attribute.Int("request_count", requestCount))

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func(requestID int) {
			defer wg.Done()

			agent, err := entities.GetAgent(c)
			if err != nil {
				slog.ErrorContext(ctx, "Request failed",
					"request_id", requestID,
					"error", err,
				)
				return
			}

			slog.InfoContext(ctx, "Request completed",
				"request_id", requestID,
				"agent", agent.Symbol,
				"credits", agent.Credits,
			)

			// Small delay to demonstrate concurrent starts
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	wg.Wait()
	concurrentSpan.End()

	duration := time.Since(startTime)
	rootSpan.SetAttributes(
		attribute.Int("requests.total", requestCount),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)

	slog.InfoContext(ctx, "Request queue demonstration complete",
		"requests", requestCount,
		"duration", duration.String(),
	)

	slog.InfoContext(ctx, "Note: All requests were rate-limited by the queue despite concurrent starts")
}
