package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
)

func main() {
	// Set up the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create a new client with a token
	options := client.DefaultClientOptions()
	options.Logger = logger

	options.Symbol = "caching-test"
	options.Faction = "COSMIC"

	client, cerr := client.NewClient(options)
	if cerr != nil {
		logger.Fatalf("Failed to create client: %v", cerr)
	}

	// Fetch systems from the API
	systems, err := entities.ListSystems(client)
	if err != nil {
		logger.Fatalf("Failed to list systems: %v", err)
	}

	// Store the fetched systems in the cache
	client.CacheClient.Set("systems", systems, 0)

	// Print the number of systems fetched
	fmt.Printf("Number of systems fetched: %v\n", len(systems))

	// Print the symbol of each system
	for _, system := range systems {
		fmt.Printf("System symbol: %v\n", system.Symbol)
	}

	// Retrieve the cached systems
	convertedSystems := []*entities.System{}
	if cachedSystems, found := client.CacheClient.Get("systems"); found {
		convertedSystems = cachedSystems.([]*entities.System)
	}

	// Print the number of converted systems retrieved from the cache
	fmt.Printf("Number of converted systems retrieved from cache: %v\n", len(convertedSystems))

}
