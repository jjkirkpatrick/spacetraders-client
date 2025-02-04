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
	options.Symbol = "caching-example"
	options.Faction = "COSMIC"
	options.LogLevel = slog.LevelInfo

	client, cerr := client.NewClient(options)
	if cerr != nil {
		slog.Error("Failed to create client", "error", cerr)
		os.Exit(1)
	}

	// Fetch systems from the API
	systems, err := entities.ListSystems(client)
	if err != nil {
		slog.Error("Failed to list systems", "error", err)
		os.Exit(1)
	}

	// Store the fetched systems in the cache
	client.CacheClient.Set("systems", systems, 0)

	// Print the number of systems fetched
	slog.Info("Systems fetched", "count", len(systems))

	// Print the symbol of each system
	for _, system := range systems {
		slog.Info("System", "symbol", system.Symbol)
	}

	// Retrieve the cached systems
	convertedSystems := []*entities.System{}
	if cachedSystems, found := client.CacheClient.Get("systems"); found {
		convertedSystems = cachedSystems.([]*entities.System)
	}

	// Print the number of converted systems retrieved from the cache
	slog.Info("Systems retrieved from cache", "count", len(convertedSystems))
}
