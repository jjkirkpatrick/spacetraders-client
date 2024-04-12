package main

import (
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

	options.Symbol = "ships-example"
	options.Faction = "COSMIC"

	client, cerr := client.NewClient(options)
	if cerr != nil {
		logger.Fatalf("Failed to create client: %v", cerr)
	}

	// Fetch systems from the API
	ships, err := entities.ListShips(client)

	if err != nil {
		logger.Fatalf("Failed to list ships: %v", err)
	}

	for _, ship := range ships {
		logger.Printf("Ship: %v\n", ship.Symbol)

	}

	ship := ships[0]

	// All reciever functions of the Ship type will update the returned data automatically
	// Calling ship.Dock() will update the ship's Nav to "DOCKED" assuming a successful API call
	_, err = ship.Dock()
	if err != nil {
		logger.Fatalf("Failed to dock: %v", err)
	}

	// All reciever functions of the Ship type will also explicitly return the data that the corresponding API call returns
	nav, err := ship.Dock()
	if err != nil {
		logger.Fatalf("Failed to get nav: %v", err)
	}

	// the value of Nav will be the same as the value of ship.Nav
	logger.Printf("Nav: %v\n", nav)
	logger.Printf("Ship NavStatus: %v\n", ship.Nav)

}
