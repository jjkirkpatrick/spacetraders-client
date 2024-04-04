package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/client"
)

func mexample_create_client() {
	// Set up the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create a new client with a token
	options := client.DefaultClientOptions()
	options.Token = "your-token-here"
	options.Logger = logger

	client, err := client.NewClient(options)
	if err != nil {
		logger.Fatalf("Failed to create client: %v", err)
	}

	// Use the client to make API requests
	var result map[string]interface{}
	err = client.Get("/my/agent", &result)
	if err != nil {
		logger.Fatalf("Failed to retrieve agent details: %v", err)
	}

	fmt.Println("Agent details:", result)
}
