package examples

import (
	"fmt"
	"log"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/client"
)

func getContractAndAccept() {
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

	// Retrieve agent details
	agent, err := client.GetAgent()
	if err != nil {
		logger.Fatalf("Failed to retrieve agent details: %v", err)
	}
	fmt.Printf("Agent Symbol: %s\n", agent.Symbol)
	fmt.Printf("Agent Headquarters: %s\n", agent.Headquarters)

	// List available contracts
	const pageSize = 5
	const pageNumber = 1

	contracts, err := client.ListContracts(pageSize, pageNumber)
	if err != nil {
		logger.Fatalf("Failed to list contracts: %v", err)
	}

	fmt.Printf("Available Contracts (Page %d):\n", pageNumber)
	for _, contract := range contracts {
		fmt.Printf("- Contract ID: %s\n", contract.ID)
	}

	// Get details of a specific contract
	const contractID = "cluk89lvh36rbs60c4i01lvhe"

	contract, err := client.GetContract(contractID)
	if err != nil {
		logger.Fatalf("Failed to get contract details: %v", err)
	}
	fmt.Printf("\nContract Details (ID: %s):\n", contract.ID)
	fmt.Printf("- Type: %s\n", contract.Type)
	fmt.Printf("- Faction: %s\n", contract.FactionSymbol)
	fmt.Printf("- Accepted: %t\n", contract.Accepted)

	// Accept the contract
	agent, contract, err = client.AcceptContract(contractID)
	if err != nil {
		logger.Fatalf("Failed to accept contract: %v", err)
	}
	fmt.Println("\nContract Accepted!")
	fmt.Printf("Agent Symbol: %s\n", agent.Symbol)
	fmt.Printf("Contract ID: %s\n", contract.ID)
	fmt.Printf("Contract Accepted: %t\n", contract.Accepted)
}
