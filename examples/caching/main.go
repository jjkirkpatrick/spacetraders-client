package main

import (
	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/phuslu/log"
)

func main() {

	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeFormat: "15:04:05",
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
			Formatter:      client.Logformat,
		},
	}

	// Create a new client with a token
	options := client.DefaultClientOptions()

	options.Symbol = "caching-example"
	options.Faction = "COSMIC"
	options.LogLevel = log.InfoLevel

	client, cerr := client.NewClient(options)
	if cerr != nil {
		log.Fatal().Msgf("Failed to create client: %v", cerr)
	}

	// Fetch systems from the API
	systems, err := entities.ListSystems(client)
	if err != nil {
		log.Fatal().Msgf("Failed to list systems: %v", err)
	}

	// Store the fetched systems in the cache
	client.CacheClient.Set("systems", systems, 0)

	// Print the number of systems fetched
	log.Info().Msgf("Number of systems fetched: %v", len(systems))

	// Print the symbol of each system
	for _, system := range systems {
		log.Info().Msgf("System symbol: %v", system.Symbol)
	}

	// Retrieve the cached systems
	convertedSystems := []*entities.System{}
	if cachedSystems, found := client.CacheClient.Get("systems"); found {
		convertedSystems = cachedSystems.([]*entities.System)
	}

	// Print the number of converted systems retrieved from the cache
	log.Info().Msgf("Number of converted systems retrieved from cache: %v", len(convertedSystems))

}
