package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/client"
)

func example_create_client_without_token() {
	// Set up the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create a new client without a token (register a new agent)
	options := client.DefaultClientOptions()
	options.Logger = logger

	faction := "COSMIC"
	symbol := "MYAGENT"
	email := "example@example.com"

	client, err := client.NewClientWithAgentRegistration(options, faction, symbol, email)
	if err != nil {
		logger.Fatalf("Failed to create client and register agent: %v", err)
	}

	// Use the client to make API requests
	var result map[string]interface{}
	err = client.Get("/my/agent", &result)
	if err != nil {
		logger.Fatalf("Failed to retrieve agent details: %v", err)
	}

	fmt.Println("Agent details:", result)
}
