package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/jjkirkpatrick/spacetraders-client/internal/telemetry"
	"github.com/jjkirkpatrick/spacetraders-client/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type GameState struct {
	Agent      *entities.Agent
	HomeSystem string
	Contracts  []*entities.Contract `json:"contracts"`
	Ships      []*entities.Ship     `json:"ships"`
}

func main() {
	ctx := context.Background()

	// Initialize slog with pretty printing
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Create a new client with a token
	options := client.DefaultClientOptions()
	options.Symbol = "BlueCa-99"
	options.Faction = "COSMIC"
	options.LogLevel = slog.LevelInfo
	options.TelemetryConfig = &telemetry.Config{
		ServiceName:    "spacetraders-metrics",
		ServiceVersion: "1.0.0",
		OTLPEndpoint:   "localhost:4317",
	}

	gameState := &GameState{}

	client, cerr := client.NewClient(options)
	if cerr != nil {
		slog.Error("Failed to create client", "error", cerr)
		os.Exit(1)
	}
	defer client.Close(ctx)

	// Get a tracer
	tracer := otel.GetTracerProvider().Tracer("quickstart-example")

	// Create a root span for the game session
	ctx, rootSpan := tracer.Start(ctx, "game_session")
	defer rootSpan.End()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Info("Shutting down gracefully...")
		rootSpan.End()
		client.Close(ctx)
		os.Exit(0)
	}()

	// Get agent information with tracing
	ctx, agentSpan := tracer.Start(ctx, "get_agent")
	agent, err := entities.GetAgent(client)
	if err != nil {
		agentSpan.RecordError(err)
		agentSpan.SetAttributes(attribute.String("error", err.Error()))
		slog.Error("Failed to get agent", "error", err)
		os.Exit(1)
	}
	agentSpan.SetAttributes(
		attribute.String("agent.symbol", agent.Symbol),
		attribute.Int64("agent.credits", agent.Credits),
	)
	agentSpan.End()

	gameState.Agent = agent
	gameState.HomeSystem = getSystemNameFromHomeSystem(gameState.Agent)

	// Get all contracts with tracing
	ctx, contractsSpan := tracer.Start(ctx, "list_contracts")
	contracts, err := entities.ListContracts(client)
	if err != nil {
		contractsSpan.RecordError(err)
		slog.Error("Failed to list contracts", "error", err)
		os.Exit(1)
	}
	contractsSpan.SetAttributes(attribute.Int("contracts.count", len(contracts)))
	contractsSpan.End()

	gameState.Contracts = contracts

	// Accept contracts with tracing
	ctx, acceptSpan := tracer.Start(ctx, "accept_contracts")
	for _, contract := range gameState.Contracts {
		agent, _, err := contract.Accept()
		if err != nil && strings.Contains(err.Error(), "has already been accepted") {
			continue
		}
		gameState.Agent = agent
		if err != nil {
			acceptSpan.RecordError(err)
			slog.Error("Failed to accept contract", "error", err)
			os.Exit(1)
		}
	}
	acceptSpan.End()

	// Get active contract terms
	activeContractTrms := gameState.getActiveContracts()

	// Print active contract details
	gameState.printActiveContractDetails(activeContractTrms)

	currentSystem, err := entities.GetSystem(client, gameState.HomeSystem)
	if err != nil {
		slog.Error("Failed to get system", "error", err)
		os.Exit(1)
	}

	// Find the engineered asteroid in the system
	astroid, err := currentSystem.GetWaypointsWithTrait("", "ENGINEERED_ASTEROID")
	if err != nil {
		slog.Error("Failed to find engineered asteroid", "error", err)
		os.Exit(1)
	}
	// Get mining ship
	miningShipSymbol, err := gameState.getMiningShip(client)
	if err != nil {
		slog.Info("No mining ship found, attempting to find a shipyard", "error", err)
		shipyards, err := currentSystem.GetWaypointsWithTrait("SHIPYARD", "")
		if err != nil || len(shipyards) == 0 {
			slog.Error("Failed to find shipyard", "error", err)
			os.Exit(1)
		}
		var shipyard *models.Shipyard
		for _, shipyardInfo := range shipyards {
			tempShipyard, err := currentSystem.GetShipyard(shipyardInfo.Symbol)
			if err != nil {
				slog.Info("Failed to get shipyard details",
					"shipyard", shipyardInfo.Symbol,
					"error", err)
				continue
			}
			for _, shipType := range tempShipyard.ShipTypes {
				if shipType.Type == models.ShipMiningDrone {
					shipyard = tempShipyard
					break
				}
			}
			if shipyard != nil {
				break
			}
		}
		if shipyard == nil {
			slog.Error("Failed to find a shipyard selling SHIP_MINING_DRONE")
			os.Exit(1)
		}
		_, ship, _, err := entities.PurchaseShip(client, "SHIP_MINING_DRONE", shipyard.Symbol)
		if err != nil {
			slog.Error("Failed to purchase mining ship", "error", err)
			os.Exit(1)
		}
		miningShipSymbol = ship.Symbol
		slog.Info("Purchased mining ship", "symbol", miningShipSymbol)
	}

	ship, err := entities.GetShip(client, miningShipSymbol)
	if err != nil {
		slog.Error("Failed to get ship", "error", err)
		os.Exit(1)
	}

	// Mine resources and deliver contracts
	err = gameState.mineAndDeliverContracts(ship, astroid[0], activeContractTrms)
	if err != nil {
		slog.Error("Failed to mine and deliver contracts", "error", err)
		os.Exit(1)
	}

	slog.Info("Program completed successfully")
}

func waitForCooldown(ship *entities.Ship) {
	_, cerr := ship.FetchCooldown()
	if cerr != nil {
		slog.Error("Failed to get ship cooldown", "error", cerr)
		os.Exit(1)
	}

	if ship.Cooldown.RemainingSeconds > 0 {
		slog.Info("Waiting for cooldown to finish",
			"remaining_seconds", ship.Cooldown.RemainingSeconds)
		time.Sleep(time.Duration(ship.Cooldown.RemainingSeconds) * time.Second)
	}
}

func getSystemNameFromHomeSystem(agent *entities.Agent) string {
	parts := strings.Split(agent.Headquarters, "-")
	if len(parts) >= 2 {
		return parts[0] + "-" + parts[1]
	}
	return ""
}

func (gs *GameState) getActiveContracts() []entities.Contract {
	var activeContractTrms []entities.Contract
	for _, contract := range gs.Contracts {
		if contract.Accepted && !contract.Fulfilled {
			activeContractTrms = append(activeContractTrms, *contract)
		}
	}
	return activeContractTrms
}

func (gs *GameState) printActiveContractDetails(activeContractTrms []entities.Contract) {
	for _, contract := range activeContractTrms {
		slog.Info("Contract", "id", contract.ID)
		for _, deliver := range contract.Terms.Deliver {
			slog.Info("Deliver", "trade_symbol", deliver.TradeSymbol, "destination_symbol", deliver.DestinationSymbol, "units_required", deliver.UnitsRequired, "units_fulfilled", deliver.UnitsFulfilled)
		}
	}
}

func (gs *GameState) getMiningShip(client *client.Client) (string, error) {
	slog.Info("Getting mining ship")
	allShips, err := entities.ListShips(client)
	if err != nil {
		return "", fmt.Errorf("failed to list ships: %v", err)
	}
	gs.Ships = allShips

	for _, ship := range gs.Ships {
		if ship.Registration.Role == models.Excavator {
			return ship.Symbol, nil
		}
	}
	return "", fmt.Errorf("no mining ship found")
}

func (gs *GameState) mineAndDeliverContracts(miningShip *entities.Ship, astroid *models.Waypoint, activeContractTrms []entities.Contract) error {
	slog.Info("Mining and delivering contracts")
	for {
		// Check if all contracts are fulfilled
		allContractsFulfilled := true
		for _, contract := range activeContractTrms {
			if !contract.Fulfilled {
				allContractsFulfilled = false
				break
			}
		}
		if allContractsFulfilled {
			slog.Info("All contracts fulfilled.")
			break
		}

		//Navigate to the asteroid
		slog.Info("Navigating to asteroid")
		err := gs.navigateToWaypoint(miningShip, *astroid)
		if err != nil {
			return fmt.Errorf("failed to navigate to asteroid: %v", err)
		}

		slog.Info("Mining resources")

		_, err = miningShip.Orbit()
		if err != nil {
			return fmt.Errorf("failed to orbit ship: %v", err)
		}

		// Mine resources
		err = gs.mineResources(miningShip, activeContractTrms)
		if err != nil {
			return fmt.Errorf("failed to mine resources: %v", err)
		}

		slog.Info("Selling trade goods")
		// Sell trade goods not required for contracts
		err = gs.jettisonUnwantedCargo(miningShip, activeContractTrms)
		if err != nil {
			return fmt.Errorf("failed to sell trade goods: %v", err)
		}

		slog.Info("Delivering contract goods")

		// Deliver goods to contract destinations
		err = gs.deliverContractGoods(miningShip, activeContractTrms)
		if err != nil {
			return fmt.Errorf("failed to deliver contract goods: %v", err)
		}
	}
	return nil
}

func (gs *GameState) navigateToWaypoint(miningShip *entities.Ship, waypointSymbol models.Waypoint) error {

	// Use the path finding system to get the route to the target waypoint
	route, Rerr := miningShip.GetRouteToDestination(waypointSymbol.Symbol)
	if Rerr != nil {
		return fmt.Errorf("failed to get route to destination: %v", Rerr)
	}

	// Log each step of the route
	for _, step := range route.Steps {
		slog.Info("Navigating to waypoint", "waypoint", step.Waypoint)
		//ensure ship is in orbit
		_, err := miningShip.Orbit()

		if err != nil {
			return fmt.Errorf("failed to orbit ship: %v", err)
		}

		err = miningShip.SetFlightMode(step.FlightMode)

		if err != nil {
			return fmt.Errorf("failed to update ship navigation: %v", err)
		}

		_, _, _, err = miningShip.Navigate(step.Waypoint)
		if err != nil {
			return fmt.Errorf("failed to navigate to waypoint %s: %v", step.Waypoint, err)
		}

		arrivalTime := miningShip.Nav.Route.Arrival
		slog.Info("Navigating to", "waypoint", waypointSymbol.Symbol, "arrival_time", arrivalTime, "using_flight_mode", step.FlightMode)

		arrivalTimeParsed, stateErr := time.Parse(time.RFC3339, arrivalTime)
		if stateErr != nil {
			return fmt.Errorf("failed to parse arrival time: %v", stateErr)
		}

		time.Sleep(time.Until(arrivalTimeParsed.Add(1 * time.Second)))

		gs.dockAndRefuelShip(miningShip)

	}

	err := miningShip.SetFlightMode(models.FlightModeCruise)
	if err != nil {
		return fmt.Errorf("failed to update ship navigation: %v", err)
	}

	refuelErr := gs.dockAndRefuelShip(miningShip)
	if refuelErr != nil {
		return fmt.Errorf("failed to refuel ship: %v", refuelErr)
	}

	return nil
}

func (gs *GameState) dockAndRefuelShip(miningShip *entities.Ship) error {
	slog.Info("Initiating refueling process for ship", "symbol", miningShip.Symbol)
	_, err := miningShip.Dock()
	if err != nil {
		return fmt.Errorf("failed to dock ship: %v", err)
	}

	_, _, _, err = miningShip.Refuel(0, false)
	if err != nil {
		return fmt.Errorf("failed to refuel ship: %v", err)
	}
	slog.Info("Refueling completed successfully for ship", "symbol", miningShip.Symbol)

	return nil
}

func (gs *GameState) mineResources(miningShip *entities.Ship, activeContractTrms []entities.Contract) error {
	cargo, err := miningShip.FetchCargo()
	if err != nil {
		return fmt.Errorf("failed to get ship cargo: %v", err)
	}

	for _, contract := range activeContractTrms {
		slog.Info("Checking contract requirements for contract ID", "id", contract.ID)
		for _, deliver := range contract.Terms.Deliver {
			slog.Info("Requirement", "units_required", deliver.UnitsRequired, "trade_symbol", deliver.TradeSymbol, "destination_symbol", deliver.DestinationSymbol)
			unitsAvailable := 0
			for _, cargoItem := range cargo.Inventory {
				if cargoItem.Symbol == deliver.TradeSymbol {
					unitsAvailable += cargoItem.Units
				}
			}
			if unitsAvailable >= deliver.UnitsRequired-deliver.UnitsFulfilled {
				slog.Info("Cargo contains enough", "trade_symbol", deliver.TradeSymbol, "to fulfill the contract requirement.")
			} else {
				for unitsAvailable < deliver.UnitsRequired {
					slog.Info("Not enough", "trade_symbol", deliver.TradeSymbol, "in cargo to fulfill the contract requirement. Required:", deliver.UnitsRequired-activeContractTrms[0].Terms.Deliver[0].UnitsFulfilled, "Available:", unitsAvailable)
					if cargo.Units >= cargo.Capacity {
						slog.Info("Cargo hold is full. Unable to extract more resources.")
						break
					}
					slog.Info("Attempting to extract more resources...")
					waitForCooldown(miningShip)
					_, err := miningShip.Extract()
					if err != nil {
						return fmt.Errorf("failed to extract resources: %v", err)
					}
					// Update cargo after extraction
					cargo, err = miningShip.FetchCargo()
					if err != nil {
						return fmt.Errorf("failed to get ship cargo: %v", err)
					}

					for _, item := range cargo.Inventory {
						if item.Symbol != deliver.TradeSymbol {
							slog.Info("Jettisoning", "name", item.Name, "as it is not required for the current contract.")
							_, jettisonErr := miningShip.Jettison(models.GoodSymbol(item.Symbol), item.Units)
							if jettisonErr != nil {
								return fmt.Errorf("failed to jettison %s: %v", item.Name, jettisonErr)
							}
							slog.Info("Successfully jettisoned", "name", item.Name)
						}
					}
					// Refresh cargo after jettisoning unnecessary items
					cargo, err = miningShip.FetchCargo()
					if err != nil {
						return fmt.Errorf("failed to refresh ship cargo: %v", err)
					}

					// Update unitsAvailable after extraction
					unitsAvailable = 0
					for _, cargoItem := range cargo.Inventory {
						if cargoItem.Symbol == deliver.TradeSymbol {
							unitsAvailable += cargoItem.Units
						}
					}

					slog.Info("Cargo Summary")
					slog.Info("Capacity", "capacity", cargo.Capacity, "units", cargo.Units)
					slog.Info("Inventory")
					for _, item := range cargo.Inventory {
						slog.Info("-", "name", item.Name, "units", item.Units)
					}
				}
				if unitsAvailable >= deliver.UnitsRequired {
					slog.Info("Now have enough", "trade_symbol", deliver.TradeSymbol, "to fulfill the contract requirement. Required:", deliver.UnitsRequired-activeContractTrms[0].Terms.Deliver[0].UnitsFulfilled, "Available:", unitsAvailable)
				}
			}
		}
	}
	return nil
}

func (gs *GameState) jettisonUnwantedCargo(miningShip *entities.Ship, activeContractTrms []entities.Contract) error {
	cargo, err := miningShip.FetchCargo()
	if err != nil {
		return fmt.Errorf("failed to get ship cargo: %v", err)
	}

	slog.Info("Starting to sell trade goods not required for the active contract.")

	for _, item := range cargo.Inventory {
		// Check if the item is not part of the active contract requirements
		if !gs.isItemRequiredForContracts(models.GoodSymbol(item.Symbol), activeContractTrms) {
			slog.Info("No markets found buying", "name", item.Name, "jettisoning cargo.")
			_, jettisonErr := miningShip.Jettison(models.GoodSymbol(item.Symbol), item.Units)
			if jettisonErr != nil {
				return fmt.Errorf("failed to jettison %s: %v", item.Name, jettisonErr)
			}
			slog.Info("Successfully jettisoned", "name", item.Name)
		}
	}

	slog.Info("Finished dumping trade goods.")
	return nil
}

func (gs *GameState) isItemRequiredForContracts(itemSymbol models.GoodSymbol, activeContractTrms []entities.Contract) bool {
	for _, contract := range activeContractTrms {
		for _, deliver := range contract.Terms.Deliver {
			if deliver.TradeSymbol == string(itemSymbol) {
				return true
			}
		}
	}
	return false
}

func (gs *GameState) deliverContractGoods(miningShip *entities.Ship, activeContractTrms []entities.Contract) error {
	for _, contract := range activeContractTrms {
		for _, deliver := range contract.Terms.Deliver {
			if deliver.UnitsFulfilled < deliver.UnitsRequired {
				slog.Info("Delivering goods for contract ID", "id", contract.ID)
				naverr := gs.navigateToWaypoint(miningShip, models.Waypoint{Symbol: deliver.DestinationSymbol})
				if naverr != nil {
					return fmt.Errorf("failed to navigate to destination: %v", naverr)
				}

				_, err := miningShip.Dock()
				if err != nil {
					return fmt.Errorf("failed to dock ship: %v", err)
				}

				cargo, err := miningShip.FetchCargo()
				if err != nil {
					return fmt.Errorf("failed to get ship cargo: %v", err)
				}

				unitsOfRequiredItem := 0
				for _, item := range cargo.Inventory {
					if item.Symbol == deliver.TradeSymbol {
						unitsOfRequiredItem = item.Units
						break
					}
				}

				_, _, err = contract.DeliverCargo(miningShip, models.GoodSymbol(deliver.TradeSymbol), unitsOfRequiredItem)
				if err != nil {
					return fmt.Errorf("failed to deliver contract: %v", err)
				}
				slog.Info("Delivered", "units", deliver.UnitsRequired-deliver.UnitsFulfilled, "trade_symbol", deliver.TradeSymbol, "for contract ID", "id", contract.ID)
			}
		}

		if !contract.Fulfilled {
			// Check if all required deliveries have been made
			allDeliveriesMade := true
			for _, deliver := range contract.Terms.Deliver {
				if deliver.UnitsFulfilled < deliver.UnitsRequired {
					allDeliveriesMade = false
					break
				}
			}
			if allDeliveriesMade {
				_, _, err := contract.Fulfill()
				if err != nil {
					return fmt.Errorf("failed to fulfill contract: %v", err)
				}
				slog.Info("Fulfilled contract ID", "id", contract.ID)
			} else {
				slog.Info("Cannot fulfill contract ID", "id", contract.ID, "as not all deliveries have been made.")
			}
		}
	}
	return nil
}
