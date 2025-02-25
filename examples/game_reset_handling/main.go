package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
)

func main() {
	// Create a client with default options
	options := client.DefaultClientOptions()
	options.Symbol = "BLUE1"       // Replace with your agent symbol
	options.Faction = "COSMIC"     // Replace with your faction
	options.RequestQueueSize = 100 // Set request queue size

	c, err := client.NewClient(options)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close(context.Background())

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Start a goroutine to monitor for game resets
	var wg sync.WaitGroup
	wg.Add(1)
	go monitorGameReset(ctx, c, &wg)

	// Start making API requests in a loop
	wg.Add(1)
	go makeRequests(ctx, c, &wg)

	fmt.Println("Application running. Press Ctrl+C to exit.")
	wg.Wait()
	fmt.Println("Application shutdown complete.")
}

// monitorGameReset monitors for game resets and handles them
func monitorGameReset(ctx context.Context, c *client.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Wait for either a game reset or context cancellation
		resetDetected := c.WaitForGameReset(ctx)

		if !resetDetected {
			// Context was cancelled, exit the loop
			fmt.Println("Game reset monitor shutting down...")
			return
		}

		// Game reset was detected
		fmt.Println("\n!!! GAME RESET DETECTED !!!")
		fmt.Println("The SpaceTraders game has been reset.")
		fmt.Println("You need to re-register your agent.")
		fmt.Println("Exiting application...")

		// Cancel the context to signal all goroutines to shut down
		// This is a more severe action that will terminate the application
		os.Exit(1)
	}
}

// makeRequests makes API requests in a loop
func makeRequests(ctx context.Context, c *client.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	requestCount := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Request loop shutting down...")
			return
		default:
			// Check if a game reset has been detected (non-blocking)
			if c.IsGameReset() {
				fmt.Println("Request loop detected game reset, stopping...")
				return
			}

			// Make an API request
			requestCount++
			fmt.Printf("Making request #%d...\n", requestCount)

			agent, err := entities.GetAgent(c)
			if err != nil {
				fmt.Printf("Error getting agent: %v\n", err)
				// Continue making requests even if there's an error
				// The game reset detection will handle token version mismatch errors
			} else {
				fmt.Printf("Agent %s has %d credits\n", agent.Symbol, agent.Credits)
			}

			// Wait a bit before making the next request
			time.Sleep(2 * time.Second)
		}
	}
}
