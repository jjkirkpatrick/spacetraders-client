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
	options.Symbol = "SHIPS-DEMO"
	options.Faction = "COSMIC"
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-ships"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"

	c, cerr := client.NewClient(options)
	if cerr != nil {
		slog.Error("Failed to create client", "error", cerr)
		os.Exit(1)
	}
	defer c.Close(ctx)

	// Set up combined logging
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-ships", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	// Get tracer
	tracer := otel.GetTracerProvider().Tracer("spacetraders-ships")

	// Create root span
	ctx, rootSpan := tracer.Start(ctx, "ships_demo")
	defer rootSpan.End()

	slog.InfoContext(ctx, "Starting ships demonstration")

	// Fetch ships from the API
	ctx, listSpan := tracer.Start(ctx, "list_ships")
	slog.InfoContext(ctx, "Fetching fleet")
	ships, err := entities.ListShips(c)
	if err != nil {
		listSpan.RecordError(err)
		slog.ErrorContext(ctx, "Failed to list ships", "error", err)
		os.Exit(1)
	}
	listSpan.SetAttributes(attribute.Int("ships.count", len(ships)))
	slog.InfoContext(ctx, "Fleet loaded", "ship_count", len(ships))
	listSpan.End()

	if len(ships) == 0 {
		slog.WarnContext(ctx, "No ships found in fleet")
		return
	}

	// Display ships
	for _, ship := range ships {
		slog.InfoContext(ctx, "Ship info",
			"symbol", ship.Symbol,
			"role", ship.Registration.Role,
			"nav_status", ship.Nav.Status,
			"fuel", ship.Fuel.Current,
			"fuel_capacity", ship.Fuel.Capacity,
		)
	}

	// Work with first ship
	ship := ships[0]
	slog.InfoContext(ctx, "Selected ship for operations", "symbol", ship.Symbol)

	// Dock the ship
	ctx, dockSpan := tracer.Start(ctx, "dock_ship")
	dockSpan.SetAttributes(attribute.String("ship.symbol", ship.Symbol))
	slog.InfoContext(ctx, "Docking ship", "symbol", ship.Symbol)

	nav, err := ship.Dock()
	if err != nil {
		dockSpan.RecordError(err)
		slog.ErrorContext(ctx, "Failed to dock ship",
			"symbol", ship.Symbol,
			"error", err,
		)
		os.Exit(1)
	}

	slog.InfoContext(ctx, "Ship docked successfully",
		"symbol", ship.Symbol,
		"status", nav.Status,
		"waypoint", nav.WaypointSymbol,
	)
	dockSpan.End()

	// Demonstrate that receiver functions update ship state
	slog.InfoContext(ctx, "Ship navigation state after dock",
		"nav_from_response", nav.Status,
		"nav_from_ship", ship.Nav.Status,
		"states_match", nav.Status == ship.Nav.Status,
	)

	slog.InfoContext(ctx, "Ships demonstration complete")
}
