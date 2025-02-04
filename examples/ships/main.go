package main

import (
	"log/slog"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
)

func main() {
	// Set up the logger
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Create a new client with a token
	options := client.DefaultClientOptions()
	options.Symbol = "ships-example"
	options.Faction = "COSMIC"

	client, cerr := client.NewClient(options)
	if cerr != nil {
		slog.Error("Failed to create client", "error", cerr)
		os.Exit(1)
	}

	// Fetch systems from the API
	ships, err := entities.ListShips(client)
	if err != nil {
		slog.Error("Failed to list ships", "error", err)
		os.Exit(1)
	}

	for _, ship := range ships {
		slog.Info("Ship", "symbol", ship.Symbol)
	}

	ship := ships[0]

	// All receiver functions of the Ship type will update the returned data automatically
	// Calling ship.Dock() will update the ship's Nav to "DOCKED" assuming a successful API call
	_, err = ship.Dock()
	if err != nil {
		slog.Error("Failed to dock", "error", err)
		os.Exit(1)
	}

	// All receiver functions of the Ship type will also explicitly return the data that the corresponding API call returns
	nav, err := ship.Dock()
	if err != nil {
		slog.Error("Failed to get nav", "error", err)
		os.Exit(1)
	}

	// the value of Nav will be the same as the value of ship.Nav
	slog.Info("Navigation status", "nav", nav, "ship_nav_status", ship.Nav)
}
