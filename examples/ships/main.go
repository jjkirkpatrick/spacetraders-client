package main

import (
	"github.com/phuslu/log"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
)

func main() {
	// Set up the logger
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

	options.Symbol = "ships-example"
	options.Faction = "COSMIC"

	client, cerr := client.NewClient(options)
	if cerr != nil {
		log.Fatal().Msgf("Failed to create client: %v", cerr)
	}

	// Fetch systems from the API
	ships, err := entities.ListShips(client)

	if err != nil {
		log.Fatal().Msgf("Failed to list ships: %v", err)
	}

	for _, ship := range ships {
		log.Info().Msgf("Ship: %v", ship.Symbol)

	}

	ship := ships[0]

	// All reciever functions of the Ship type will update the returned data automatically
	// Calling ship.Dock() will update the ship's Nav to "DOCKED" assuming a successful API call
	_, err = ship.Dock()
	if err != nil {
		log.Fatal().Msgf("Failed to dock: %v", err)
	}

	// All reciever functions of the Ship type will also explicitly return the data that the corresponding API call returns
	nav, err := ship.Dock()
	if err != nil {
		log.Fatal().Msgf("Failed to get nav: %v", err)
	}

	// the value of Nav will be the same as the value of ship.Nav
	log.Info().Msgf("Nav: %v", nav)
	log.Info().Msgf("Ship NavStatus: %v", ship.Nav)

}
