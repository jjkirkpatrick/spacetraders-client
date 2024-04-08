package main

import (
	"log"
	"os"
	"strings"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/metrics"
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

type GameState struct {
	Agent      models.Agent
	HomeSystem string
	Contracts  []models.Contract              `json:"contracts"`
	Waypoints  []models.Waypoint              `json:"waypoints"`
	ShipYards  []models.ListWaypointsResponse `json:"shipyards"`
	Ships      []*models.Ship                 `json:"ships"`
}

func pathfindingExample_Test() {
	// Set up the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create a new client with a token
	options := client.DefaultClientOptions()
	options.Logger = logger
	options.Token = "Space Traders Token"

	metricsReporter := metrics.NewMetricsClient(
		"http://192.168.1.33:8086",
		"Influx token",
		"spacetraders",
		"spacetraders",
	)

	gameState := &GameState{}

	client, cerr := client.NewClient(options, metricsReporter)
	if cerr != nil {
		logger.Fatalf("Failed to create client: %v", cerr)
	}

	agent, err := client.GetAgent()
	if err != nil {
		logger.Fatalf("Failed to get agent: %v", err)
	}
	gameState.Agent = *agent
	gameState.HomeSystem = getSystemNameFromHomeSystem(gameState.Agent)

	ships, err := client.GetShip("MYAGENT4-3")

	minerGraph, gerr := client.BuildGraph("X1-KS68", ships.Engine.Speed)
	if gerr != nil {
		logger.Fatalf("Failed to build graph: %v", gerr)
	}

	// Get path bettween X1-KS68-F41 and X1-KS68-H44
	route, pathErr := client.GetRouteToDestination(minerGraph, "X1-KS68-H44", "MYAGENT4-3")

	if pathErr != nil {
		logger.Fatalf("Failed to get route: %v", pathErr)
	}

	for _, step := range route.Steps {
		logger.Printf("Waypoint: %s, FlightMode: %s", step.Waypoint, step.FlightMode)
	}

}

func getSystemNameFromHomeSystem(agent models.Agent) string {
	parts := strings.Split(agent.Headquarters, "-")
	if len(parts) >= 2 {
		return parts[0] + "-" + parts[1]
	}
	return ""
}
