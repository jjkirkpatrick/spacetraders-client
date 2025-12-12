package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/jjkirkpatrick/spacetraders-client/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	ctx := context.Background()

	// Configure client with telemetry
	options := client.DefaultClientOptions()
	options.Symbol = "OTEL-DEMO"          // Your agent symbol
	options.Faction = "COSMIC"            // Starting faction
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-otel-demo"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317" // OTEL collector endpoint

	// Create the client (initializes telemetry)
	c, err := client.NewClient(options)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}
	defer c.Close(ctx)

	// Set up combined logging (console + OTLP/Loki)
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-otel-demo", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	// Get a tracer for creating spans
	tracer := otel.GetTracerProvider().Tracer("spacetraders-otel-demo")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Info("Shutting down...")
		c.Close(ctx)
		os.Exit(0)
	}()

	// Create a root span for the operation
	ctx, span := tracer.Start(ctx, "get_agent_info")
	defer span.End()

	// Get agent data - this will generate metrics and be part of the trace
	slog.InfoContext(ctx, "Fetching agent information")
	agent, apiErr := entities.GetAgent(c)
	if apiErr != nil {
		span.RecordError(apiErr)
		slog.ErrorContext(ctx, "Failed to get agent", "error", apiErr)
		os.Exit(1)
	}

	// Add agent info to the span
	span.SetAttributes(
		attribute.String("agent.symbol", agent.Symbol),
		attribute.Int64("agent.credits", agent.Credits),
		attribute.String("agent.headquarters", agent.Headquarters),
	)

	// Log with trace context (trace_id and span_id auto-injected)
	slog.InfoContext(ctx, "Agent retrieved successfully",
		"symbol", agent.Symbol,
		"credits", agent.Credits,
		"headquarters", agent.Headquarters,
		"ship_count", agent.ShipCount,
	)

	slog.InfoContext(ctx, "Demo complete - check Grafana for traces (Tempo), logs (Loki), and metrics")
}
