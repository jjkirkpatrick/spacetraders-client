package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/jjkirkpatrick/spacetraders-client/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	ctx := context.Background()

	// Create client with telemetry
	options := client.DefaultClientOptions()
	options.Symbol = "CONCURRENT-AGENT"
	options.Faction = "COSMIC"
	options.RequestQueueSize = 1000

	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-concurrent-agent"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"
	options.TelemetryOptions.MetricInterval = 1 * time.Second

	c, err := client.NewClient(options)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}
	defer c.Close(ctx)

	// Set up combined logging
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-concurrent-agent", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	// Get tracer
	tracer := otel.GetTracerProvider().Tracer("spacetraders-concurrent-agent")

	// Create root span
	ctx, rootSpan := tracer.Start(ctx, "concurrent_agent_requests")
	defer rootSpan.End()

	// Configuration
	totalRequests := 20
	numGoroutines := 4
	requestsPerGoroutine := totalRequests / numGoroutines

	slog.InfoContext(ctx, "Starting concurrent agent request demonstration",
		"total_requests", totalRequests,
		"goroutines", numGoroutines,
		"requests_per_goroutine", requestsPerGoroutine,
	)

	// Channels for results
	results := make(chan *entities.Agent, totalRequests)
	errors := make(chan error, totalRequests)
	var wg sync.WaitGroup

	startTime := time.Now()

	// Launch goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()

			// Create span for this goroutine
			_, goroutineSpan := tracer.Start(ctx, "goroutine_requests")
			goroutineSpan.SetAttributes(attribute.Int("goroutine_id", routineID))
			defer goroutineSpan.End()

			for j := 0; j < requestsPerGoroutine; j++ {
				requestID := routineID*requestsPerGoroutine + j

				agent, err := entities.GetAgent(c)
				if err != nil {
					slog.ErrorContext(ctx, "Request failed",
						"goroutine_id", routineID,
						"request_id", requestID,
						"error", err,
					)
					errors <- err
					continue
				}

				results <- agent

				if (j+1)%5 == 0 {
					slog.InfoContext(ctx, "Goroutine progress",
						"goroutine_id", routineID,
						"completed", j+1,
						"total", requestsPerGoroutine,
					)
				}
			}

			slog.InfoContext(ctx, "Goroutine completed all requests",
				"goroutine_id", routineID,
				"requests", requestsPerGoroutine,
			)
		}(i)
	}

	// Close channels when done
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	successCount := 0
	for agent := range results {
		successCount++
		if successCount <= 3 || successCount > totalRequests-3 {
			slog.InfoContext(ctx, "Agent response received",
				"request_num", successCount,
				"symbol", agent.Symbol,
				"credits", agent.Credits,
			)
		} else if successCount == 4 {
			slog.InfoContext(ctx, "Continuing to process responses...",
				"showing", "first 3 and last 3 only",
			)
		}
	}

	// Count errors
	errorCount := 0
	for err := range errors {
		errorCount++
		if errorCount <= 3 {
			slog.ErrorContext(ctx, "Request error", "error", err)
		}
	}

	duration := time.Since(startTime)
	rootSpan.SetAttributes(
		attribute.Int("requests.total", totalRequests),
		attribute.Int("requests.success", successCount),
		attribute.Int("requests.failed", errorCount),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)

	slog.InfoContext(ctx, "Concurrent request test complete",
		"total_requests", totalRequests,
		"successful", successCount,
		"failed", errorCount,
		"duration", duration.String(),
		"avg_per_request", (duration / time.Duration(totalRequests)).String(),
	)

	if successCount == totalRequests {
		slog.InfoContext(ctx, "All requests completed successfully")
	} else {
		slog.WarnContext(ctx, "Some requests failed", "failed_count", errorCount)
	}
}
