package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
)

func main() {
	// Create a client with default options
	options := client.DefaultClientOptions()
	options.Symbol = "BLUE1"       // Replace with your agent symbol
	options.Faction = "COSMIC"     // Replace with your faction
	options.RequestQueueSize = 100 // Set the request queue size

	c, err := client.NewClient(options)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close(context.Background())

	// Demonstrate concurrent API calls
	var wg sync.WaitGroup
	requestCount := 10

	fmt.Println("Starting concurrent API calls...")
	startTime := time.Now()

	// Make multiple concurrent requests
	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Get agent information
			agent, err := entities.GetAgent(c)
			if err != nil {
				fmt.Printf("Error getting agent (request %d): %v\n", i, err)
				return
			}

			fmt.Printf("Request %d: Agent %s has %d credits\n", i, agent.Symbol, agent.Credits)

			// Small delay between starting goroutines to demonstrate they're concurrent
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	// Wait for all requests to complete
	wg.Wait()
	duration := time.Since(startTime)

	fmt.Printf("\nAll %d requests completed in %v\n", requestCount, duration)
	fmt.Println("Note: Requests were processed through the queue at a controlled rate")
	fmt.Println("Even though the goroutines started concurrently, the API calls were rate-limited")
}
