package main

import (
	"context"
	"log/slog"
	"os"

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
	options.Symbol = "CACHE-DEMO"
	options.Faction = "COSMIC"
	options.LogLevel = slog.LevelInfo
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-caching"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"

	c, cerr := client.NewClient(options)
	if cerr != nil {
		slog.Error("Failed to create client", "error", cerr)
		os.Exit(1)
	}
	defer c.Close(ctx)

	// Set up combined logging (console + OTLP)
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-caching", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	// Get tracer
	tracer := otel.GetTracerProvider().Tracer("spacetraders-caching")

	// Create root span
	ctx, rootSpan := tracer.Start(ctx, "caching_demo")
	defer rootSpan.End()

	slog.InfoContext(ctx, "Starting caching demonstration")

	// Fetch systems from the API
	ctx, fetchSpan := tracer.Start(ctx, "fetch_systems")
	slog.InfoContext(ctx, "Fetching systems from API")
	systems, err := entities.ListSystems(c)
	if err != nil {
		fetchSpan.RecordError(err)
		slog.ErrorContext(ctx, "Failed to list systems", "error", err)
		os.Exit(1)
	}
	fetchSpan.SetAttributes(attribute.Int("systems.count", len(systems)))
	slog.InfoContext(ctx, "Systems fetched from API", "count", len(systems))
	fetchSpan.End()

	// Store in cache
	ctx, cacheSpan := tracer.Start(ctx, "cache_systems")
	c.CacheClient.Set("systems", systems, 0)
	slog.InfoContext(ctx, "Systems stored in cache", "count", len(systems))
	cacheSpan.End()

	// Display first few systems
	displayCount := 5
	if len(systems) < displayCount {
		displayCount = len(systems)
	}
	for i := 0; i < displayCount; i++ {
		slog.InfoContext(ctx, "System info",
			"symbol", systems[i].Symbol,
			"type", systems[i].Type,
			"x", systems[i].X,
			"y", systems[i].Y,
		)
	}
	if len(systems) > displayCount {
		slog.InfoContext(ctx, "Additional systems not displayed", "remaining", len(systems)-displayCount)
	}

	// Retrieve from cache
	ctx, retrieveSpan := tracer.Start(ctx, "retrieve_from_cache")
	retrieveSpan.SetAttributes(attribute.String("cache.key", "systems"))
	convertedSystems := []*entities.System{}
	if cachedSystems, found := c.CacheClient.Get("systems"); found {
		convertedSystems = cachedSystems.([]*entities.System)
		slog.InfoContext(ctx, "Systems retrieved from cache",
			"count", len(convertedSystems),
			"cache_hit", true,
		)
	} else {
		slog.WarnContext(ctx, "Cache miss for systems key", "cache_hit", false)
	}
	retrieveSpan.End()

	slog.InfoContext(ctx, "Caching demonstration complete",
		"fetched", len(systems),
		"cached", len(convertedSystems),
	)
}
