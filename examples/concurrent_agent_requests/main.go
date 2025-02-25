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
	options.Symbol = "BLUE1"        // Replace with your agent symbol
	options.Faction = "COSMIC"      // Replace with your faction
	options.RequestQueueSize = 1000 // Set a larger request queue size to handle all requests

	// Initialize telemetry with the new public options
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-concurrent-agent-requests"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"
	options.TelemetryOptions.MetricInterval = 1 * time.Second

	c, err := client.NewClient(options)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer c.Close(context.Background())

	// Configuration
	totalRequests := 20
	numGoroutines := 1
	requestsPerGoroutine := totalRequests / numGoroutines

	fmt.Println("This example demonstrates making concurrent API requests using the request queue.")
	fmt.Println("With our optimized rate limiting, we expect to achieve close to 2 requests per second")
	fmt.Println("while still avoiding rate limit errors.")
	fmt.Println()

	// Channels for collecting results and errors
	results := make(chan *entities.Agent, totalRequests)
	errors := make(chan error, totalRequests)

	// WaitGroup to wait for all goroutines to complete
	var wg sync.WaitGroup

	fmt.Printf("Starting %d GetAgent requests across %d goroutines...\n", totalRequests, numGoroutines)
	startTime := time.Now()

	// Launch goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				requestID := routineID*requestsPerGoroutine + j

				// Get agent information
				agent, err := entities.GetAgent(c)
				if err != nil {
					fmt.Printf("Error in goroutine %d, request %d: %v\n", routineID, requestID, err)
					errors <- fmt.Errorf("request %d failed: %w", requestID, err)
					continue
				}

				// Send successful result to channel
				results <- agent

				if j%5 == 0 {
					fmt.Printf("Goroutine %d completed %d/%d requests\n", routineID, j+1, requestsPerGoroutine)
				}
			}

			fmt.Printf("Goroutine %d completed all %d requests\n", routineID, requestsPerGoroutine)
		}(i)
	}

	// Start a goroutine to close channels when all work is done
	go func() {
		wg.Wait()
		close(results)
		close(errors)
		fmt.Printf("\nAll goroutines completed\n")
	}()

	// Collect and count results
	successCount := 0
	errorCount := 0

	// Process results as they come in
	for agent := range results {
		successCount++
		if successCount <= 5 || successCount > totalRequests-5 {
			fmt.Printf("Agent %s has %d credits\n", agent.Symbol, agent.Credits)
		} else if successCount == 6 {
			fmt.Println("... (showing only first and last 5 results) ...")
		}
	}

	// Process errors
	for err := range errors {
		errorCount++
		if errorCount <= 5 {
			fmt.Printf("Error: %v\n", err)
		} else if errorCount == 6 {
			fmt.Println("... (more errors omitted) ...")
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nSummary:\n")
	fmt.Printf("- Total requests: %d\n", totalRequests)
	fmt.Printf("- Successful requests: %d\n", successCount)
	fmt.Printf("- Failed requests: %d\n", errorCount)
	fmt.Printf("- Time taken: %v\n", duration)
	fmt.Printf("- Average time per request: %v\n", duration/time.Duration(totalRequests))

	if successCount == totalRequests {
		fmt.Println("\nSUCCESS: All requests completed without errors!")
	} else {
		fmt.Printf("\nWARNING: %d requests failed\n", errorCount)
	}

	fmt.Println("\nNote: Even though requests were made concurrently across multiple goroutines,")
	fmt.Println("the request queue ensured they were processed at a controlled rate to respect API rate limits.")
}
