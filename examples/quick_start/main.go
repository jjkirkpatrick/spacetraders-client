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
	"go.opentelemetry.io/otel/trace"
)

type GameState struct {
	Agent      *entities.Agent
	HomeSystem string
	Contracts  []*entities.Contract `json:"contracts"`
	Ships      []*entities.Ship     `json:"ships"`
}

var tracer trace.Tracer

func main() {
	ctx := context.Background()

	// Create a new client with a token
	options := client.DefaultClientOptions()
	options.Symbol = "BlueCa-99"
	options.Faction = "COSMIC"
	options.LogLevel = slog.LevelInfo

	// Initialize telemetry with the new public options
	options.TelemetryOptions = client.DefaultTelemetryOptions()
	options.TelemetryOptions.ServiceName = "spacetraders-quickstart"
	options.TelemetryOptions.ServiceVersion = "1.0.0"
	options.TelemetryOptions.OTLPEndpoint = "localhost:4317"

	gameState := &GameState{}

	c, cerr := client.NewClient(options)
	if cerr != nil {
		slog.Error("Failed to create client", "error", cerr)
		os.Exit(1)
	}
	defer c.Close(ctx)

	// Initialize slog with combined handler (console + OTLP/Loki)
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	combinedHandler := telemetry.NewCombinedSlogHandler("spacetraders-quickstart", slog.LevelInfo, consoleHandler)
	slog.SetDefault(slog.New(combinedHandler))

	// Get a tracer for creating spans
	tracer = otel.GetTracerProvider().Tracer("spacetraders-quickstart")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		slog.Info("Received shutdown signal, cleaning up...")
		c.Close(ctx)
		os.Exit(0)
	}()

	// Phase 1: Initialize game state (discrete trace)
	agent, contracts, currentSystem := initializeGameState(ctx, c, gameState)

	// Phase 2: Setup mining (discrete trace)
	ship, asteroid := setupMining(ctx, c, gameState, currentSystem)

	// Phase 3: Mining loop - each iteration is its own trace
	activeContracts := gameState.getActiveContracts()
	runMiningLoop(ctx, gameState, ship, asteroid, activeContracts)

	slog.Info("Game session completed successfully",
		"agent", agent.Symbol,
		"contracts_completed", len(contracts),
	)
}

// initializeGameState loads agent, contracts and home system - single discrete trace
func initializeGameState(ctx context.Context, c *client.Client, gameState *GameState) (*entities.Agent, []*entities.Contract, *entities.System) {
	ctx, span := tracer.Start(ctx, "initialize_game_state")
	defer span.End()

	// Fetch agent
	slog.InfoContext(ctx, "Fetching agent information")
	agent, err := entities.GetAgent(c)
	if err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Failed to fetch agent", "error", err)
		os.Exit(1)
	}
	span.SetAttributes(
		attribute.String("agent.symbol", agent.Symbol),
		attribute.Int64("agent.credits", agent.Credits),
	)
	slog.InfoContext(ctx, "Agent loaded",
		"symbol", agent.Symbol,
		"credits", agent.Credits,
		"headquarters", agent.Headquarters,
	)

	gameState.Agent = agent
	gameState.HomeSystem = getSystemNameFromHomeSystem(agent)

	// Fetch contracts
	slog.InfoContext(ctx, "Fetching contracts")
	contracts, err := entities.ListContracts(c)
	if err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Failed to fetch contracts", "error", err)
		os.Exit(1)
	}
	slog.InfoContext(ctx, "Contracts loaded", "count", len(contracts))
	gameState.Contracts = contracts

	// Accept contracts
	for _, contract := range contracts {
		updatedAgent, _, err := contract.Accept()
		if err != nil && strings.Contains(err.Error(), "has already been accepted") {
			slog.InfoContext(ctx, "Contract already accepted", "contract_id", contract.ID)
			continue
		}
		if err != nil {
			span.RecordError(err)
			slog.ErrorContext(ctx, "Failed to accept contract", "contract_id", contract.ID, "error", err)
			os.Exit(1)
		}
		gameState.Agent = updatedAgent
		slog.InfoContext(ctx, "Contract accepted", "contract_id", contract.ID)
	}

	// Fetch home system
	slog.InfoContext(ctx, "Fetching home system", "system", gameState.HomeSystem)
	currentSystem, err := entities.GetSystem(c, gameState.HomeSystem)
	if err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Failed to fetch home system", "error", err)
		os.Exit(1)
	}
	slog.InfoContext(ctx, "Home system loaded", "symbol", currentSystem.Symbol, "type", currentSystem.Type)

	span.SetAttributes(
		attribute.Int("contracts.count", len(contracts)),
		attribute.String("home_system", currentSystem.Symbol),
	)

	return agent, contracts, currentSystem
}

// setupMining finds asteroid and ensures we have a mining ship - single discrete trace
func setupMining(ctx context.Context, c *client.Client, gameState *GameState, currentSystem *entities.System) (*entities.Ship, *models.Waypoint) {
	ctx, span := tracer.Start(ctx, "setup_mining")
	defer span.End()

	// Find engineered asteroid
	slog.InfoContext(ctx, "Searching for engineered asteroid")
	asteroids, err := currentSystem.GetWaypointsWithTrait("", "ENGINEERED_ASTEROID")
	if err != nil || len(asteroids) == 0 {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Failed to find engineered asteroid", "error", err)
		os.Exit(1)
	}
	asteroid := asteroids[0]
	slog.InfoContext(ctx, "Engineered asteroid found", "waypoint", asteroid.Symbol)

	// Get or purchase mining ship
	miningShipSymbol, err := gameState.getMiningShip(ctx, c)
	if err != nil {
		slog.InfoContext(ctx, "No mining ship found, searching for shipyard", "reason", err.Error())

		shipyards, err := currentSystem.GetWaypointsWithTrait("SHIPYARD", "")
		if err != nil || len(shipyards) == 0 {
			span.RecordError(err)
			slog.ErrorContext(ctx, "Failed to find shipyard", "error", err)
			os.Exit(1)
		}

		var shipyard *models.Shipyard
		for _, shipyardInfo := range shipyards {
			tempShipyard, err := currentSystem.GetShipyard(shipyardInfo.Symbol)
			if err != nil {
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
			slog.ErrorContext(ctx, "No shipyard found selling SHIP_MINING_DRONE")
			os.Exit(1)
		}

		_, ship, _, err := entities.PurchaseShip(c, "SHIP_MINING_DRONE", shipyard.Symbol)
		if err != nil {
			span.RecordError(err)
			slog.ErrorContext(ctx, "Failed to purchase mining ship", "error", err)
			os.Exit(1)
		}
		miningShipSymbol = ship.Symbol
		slog.InfoContext(ctx, "Purchased mining ship", "symbol", miningShipSymbol)
	}

	ship, err := entities.GetShip(c, miningShipSymbol)
	if err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Failed to fetch ship details", "error", err)
		os.Exit(1)
	}

	span.SetAttributes(
		attribute.String("ship.symbol", ship.Symbol),
		attribute.String("asteroid.symbol", asteroid.Symbol),
	)
	slog.InfoContext(ctx, "Mining setup complete",
		"ship", ship.Symbol,
		"asteroid", asteroid.Symbol,
	)

	return ship, asteroid
}

// runMiningLoop executes mining iterations - each iteration is its own trace
func runMiningLoop(ctx context.Context, gameState *GameState, ship *entities.Ship, asteroid *models.Waypoint, activeContracts []entities.Contract) {
	iteration := 0
	for {
		// Check if all contracts fulfilled
		allFulfilled := true
		for _, contract := range activeContracts {
			if !contract.Fulfilled {
				allFulfilled = false
				break
			}
		}
		if allFulfilled {
			slog.Info("All contracts fulfilled")
			break
		}

		iteration++
		// Each mining iteration gets its own trace
		executeMiningIteration(ctx, gameState, ship, asteroid, activeContracts, iteration)
	}
}

func waitForCooldown(ctx context.Context, ship *entities.Ship) {
	_, cerr := ship.FetchCooldown()
	if cerr != nil {
		slog.ErrorContext(ctx, "Failed to fetch ship cooldown", "ship", ship.Symbol, "error", cerr)
		os.Exit(1)
	}

	if ship.Cooldown.RemainingSeconds > 0 {
		slog.InfoContext(ctx, "Waiting for cooldown",
			"ship", ship.Symbol,
			"remaining_seconds", ship.Cooldown.RemainingSeconds,
		)
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

func (gs *GameState) printActiveContractDetails(ctx context.Context, activeContractTrms []entities.Contract) {
	for _, contract := range activeContractTrms {
		slog.InfoContext(ctx, "Active contract", "contract_id", contract.ID, "type", contract.Type)
		for _, deliver := range contract.Terms.Deliver {
			slog.InfoContext(ctx, "Contract delivery requirement",
				"contract_id", contract.ID,
				"trade_symbol", deliver.TradeSymbol,
				"destination", deliver.DestinationSymbol,
				"units_required", deliver.UnitsRequired,
				"units_fulfilled", deliver.UnitsFulfilled,
			)
		}
	}
}

func (gs *GameState) getMiningShip(ctx context.Context, c *client.Client) (string, error) {
	slog.InfoContext(ctx, "Searching for mining ship in fleet")
	allShips, err := entities.ListShips(c)
	if err != nil {
		return "", fmt.Errorf("failed to list ships: %v", err)
	}
	gs.Ships = allShips
	slog.InfoContext(ctx, "Fleet loaded", "ship_count", len(allShips))

	for _, ship := range gs.Ships {
		if ship.Registration.Role == models.Excavator {
			slog.InfoContext(ctx, "Mining ship found", "symbol", ship.Symbol)
			return ship.Symbol, nil
		}
	}
	return "", fmt.Errorf("no mining ship found in fleet")
}

// executeMiningIteration performs one complete mining cycle - its own discrete trace
func executeMiningIteration(ctx context.Context, gs *GameState, ship *entities.Ship, asteroid *models.Waypoint, activeContracts []entities.Contract, iteration int) {
	// Fresh context for this trace (not nested under parent)
	ctx, span := tracer.Start(context.Background(), "mining_iteration")
	defer span.End()

	span.SetAttributes(
		attribute.Int("iteration", iteration),
		attribute.String("ship", ship.Symbol),
		attribute.String("asteroid", asteroid.Symbol),
	)

	slog.InfoContext(ctx, "Starting mining iteration",
		"iteration", iteration,
		"ship", ship.Symbol,
	)

	// Navigate to asteroid
	slog.InfoContext(ctx, "Navigating to asteroid", "target", asteroid.Symbol)
	if err := gs.navigateToWaypoint(ctx, ship, *asteroid); err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Navigation failed", "error", err)
		return
	}

	// Enter orbit
	if _, err := ship.Orbit(); err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Failed to orbit", "error", err)
		return
	}
	slog.InfoContext(ctx, "Ship in orbit", "ship", ship.Symbol)

	// Mine resources
	if err := gs.mineResources(ctx, ship, activeContracts); err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Mining failed", "error", err)
		return
	}

	// Jettison unwanted cargo
	if err := gs.jettisonUnwantedCargo(ctx, ship, activeContracts); err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Jettison failed", "error", err)
		return
	}

	// Deliver goods
	if err := gs.deliverContractGoods(ctx, ship, activeContracts); err != nil {
		span.RecordError(err)
		slog.ErrorContext(ctx, "Delivery failed", "error", err)
		return
	}

	slog.InfoContext(ctx, "Mining iteration complete", "iteration", iteration)
}

func (gs *GameState) navigateToWaypoint(ctx context.Context, miningShip *entities.Ship, waypointSymbol models.Waypoint) error {
	route, Rerr := miningShip.GetRouteToDestination(waypointSymbol.Symbol)
	if Rerr != nil {
		return fmt.Errorf("failed to get route to destination: %v", Rerr)
	}

	slog.InfoContext(ctx, "Route calculated",
		"destination", waypointSymbol.Symbol,
		"steps", len(route.Steps),
	)

	for i, step := range route.Steps {
		slog.InfoContext(ctx, "Navigating to waypoint",
			"step", i+1,
			"waypoint", step.Waypoint,
			"flight_mode", step.FlightMode,
		)

		if _, err := miningShip.Orbit(); err != nil {
			return fmt.Errorf("failed to orbit ship: %v", err)
		}

		if err := miningShip.SetFlightMode(step.FlightMode); err != nil {
			return fmt.Errorf("failed to set flight mode: %v", err)
		}

		if _, _, _, err := miningShip.Navigate(step.Waypoint); err != nil {
			return fmt.Errorf("failed to navigate to waypoint %s: %v", step.Waypoint, err)
		}

		arrivalTime := miningShip.Nav.Route.Arrival
		arrivalTimeParsed, stateErr := time.Parse(time.RFC3339, arrivalTime)
		if stateErr != nil {
			return fmt.Errorf("failed to parse arrival time: %v", stateErr)
		}

		waitDuration := time.Until(arrivalTimeParsed.Add(1 * time.Second))
		slog.InfoContext(ctx, "In transit",
			"destination", step.Waypoint,
			"wait_seconds", int(waitDuration.Seconds()),
		)
		time.Sleep(waitDuration)

		gs.dockAndRefuelShip(ctx, miningShip)
	}

	if err := miningShip.SetFlightMode(models.FlightModeCruise); err != nil {
		return fmt.Errorf("failed to reset flight mode: %v", err)
	}

	return gs.dockAndRefuelShip(ctx, miningShip)
}

func (gs *GameState) dockAndRefuelShip(ctx context.Context, miningShip *entities.Ship) error {
	if _, err := miningShip.Dock(); err != nil {
		return fmt.Errorf("failed to dock ship: %v", err)
	}

	fuelBefore := miningShip.Fuel.Current
	if _, _, _, err := miningShip.Refuel(0, false); err != nil {
		return fmt.Errorf("failed to refuel ship: %v", err)
	}

	slog.InfoContext(ctx, "Ship refueled",
		"ship", miningShip.Symbol,
		"fuel_before", fuelBefore,
		"fuel_after", miningShip.Fuel.Current,
	)

	return nil
}

func (gs *GameState) mineResources(ctx context.Context, miningShip *entities.Ship, activeContractTrms []entities.Contract) error {
	cargo, err := miningShip.FetchCargo()
	if err != nil {
		return fmt.Errorf("failed to get ship cargo: %v", err)
	}

	for _, contract := range activeContractTrms {
		for _, deliver := range contract.Terms.Deliver {
			slog.InfoContext(ctx, "Mining for contract requirement",
				"contract_id", contract.ID,
				"trade_symbol", deliver.TradeSymbol,
				"units_required", deliver.UnitsRequired,
				"units_fulfilled", deliver.UnitsFulfilled,
			)

			unitsAvailable := 0
			for _, cargoItem := range cargo.Inventory {
				if cargoItem.Symbol == deliver.TradeSymbol {
					unitsAvailable += cargoItem.Units
				}
			}

			unitsNeeded := deliver.UnitsRequired - deliver.UnitsFulfilled
			if unitsAvailable >= unitsNeeded {
				slog.InfoContext(ctx, "Cargo has enough for contract",
					"trade_symbol", deliver.TradeSymbol,
					"available", unitsAvailable,
					"needed", unitsNeeded,
				)
				continue
			}

			for unitsAvailable < deliver.UnitsRequired {
				slog.InfoContext(ctx, "Mining additional resources",
					"trade_symbol", deliver.TradeSymbol,
					"available", unitsAvailable,
					"needed", unitsNeeded,
				)

				if cargo.Units >= cargo.Capacity {
					slog.WarnContext(ctx, "Cargo hold full",
						"used", cargo.Units,
						"capacity", cargo.Capacity,
					)
					break
				}

				waitForCooldown(ctx, miningShip)

				extraction, err := miningShip.Extract()
				if err != nil {
					return fmt.Errorf("failed to extract resources: %v", err)
				}

				slog.InfoContext(ctx, "Extracted resources",
					"symbol", extraction.Yield.Symbol,
					"units", extraction.Yield.Units,
				)

				cargo, err = miningShip.FetchCargo()
				if err != nil {
					return fmt.Errorf("failed to get ship cargo: %v", err)
				}

				// Jettison unwanted items immediately
				for _, item := range cargo.Inventory {
					if item.Symbol != deliver.TradeSymbol {
						slog.InfoContext(ctx, "Jettisoning unwanted cargo",
							"item", item.Name,
							"units", item.Units,
						)
						_, jettisonErr := miningShip.Jettison(models.GoodSymbol(item.Symbol), item.Units)
						if jettisonErr != nil {
							return fmt.Errorf("failed to jettison %s: %v", item.Name, jettisonErr)
						}
					}
				}

				cargo, err = miningShip.FetchCargo()
				if err != nil {
					return fmt.Errorf("failed to refresh cargo: %v", err)
				}

				unitsAvailable = 0
				for _, cargoItem := range cargo.Inventory {
					if cargoItem.Symbol == deliver.TradeSymbol {
						unitsAvailable += cargoItem.Units
					}
				}

				slog.InfoContext(ctx, "Cargo status",
					"used", cargo.Units,
					"capacity", cargo.Capacity,
					"target_units", unitsAvailable,
				)
			}

			if unitsAvailable >= unitsNeeded {
				slog.InfoContext(ctx, "Target resources collected",
					"trade_symbol", deliver.TradeSymbol,
					"collected", unitsAvailable,
				)
			}
		}
	}
	return nil
}

func (gs *GameState) jettisonUnwantedCargo(ctx context.Context, miningShip *entities.Ship, activeContractTrms []entities.Contract) error {
	cargo, err := miningShip.FetchCargo()
	if err != nil {
		return fmt.Errorf("failed to get ship cargo: %v", err)
	}

	for _, item := range cargo.Inventory {
		if !gs.isItemRequiredForContracts(models.GoodSymbol(item.Symbol), activeContractTrms) {
			slog.InfoContext(ctx, "Jettisoning cargo not needed for contracts",
				"item", item.Name,
				"symbol", item.Symbol,
				"units", item.Units,
			)
			_, jettisonErr := miningShip.Jettison(models.GoodSymbol(item.Symbol), item.Units)
			if jettisonErr != nil {
				return fmt.Errorf("failed to jettison %s: %v", item.Name, jettisonErr)
			}
		}
	}

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

func (gs *GameState) deliverContractGoods(ctx context.Context, miningShip *entities.Ship, activeContractTrms []entities.Contract) error {
	for _, contract := range activeContractTrms {
		for _, deliver := range contract.Terms.Deliver {
			if deliver.UnitsFulfilled < deliver.UnitsRequired {
				slog.InfoContext(ctx, "Delivering goods for contract",
					"contract_id", contract.ID,
					"destination", deliver.DestinationSymbol,
					"trade_symbol", deliver.TradeSymbol,
				)

				if err := gs.navigateToWaypoint(ctx, miningShip, models.Waypoint{Symbol: deliver.DestinationSymbol}); err != nil {
					return fmt.Errorf("failed to navigate to destination: %v", err)
				}

				if _, err := miningShip.Dock(); err != nil {
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

				if _, _, err := contract.DeliverCargo(miningShip, models.GoodSymbol(deliver.TradeSymbol), unitsOfRequiredItem); err != nil {
					return fmt.Errorf("failed to deliver contract: %v", err)
				}

				slog.InfoContext(ctx, "Cargo delivered to contract",
					"contract_id", contract.ID,
					"trade_symbol", deliver.TradeSymbol,
					"units_delivered", unitsOfRequiredItem,
				)
			}
		}

		if !contract.Fulfilled {
			allDeliveriesMade := true
			for _, deliver := range contract.Terms.Deliver {
				if deliver.UnitsFulfilled < deliver.UnitsRequired {
					allDeliveriesMade = false
					break
				}
			}
			if allDeliveriesMade {
				if _, _, err := contract.Fulfill(); err != nil {
					return fmt.Errorf("failed to fulfill contract: %v", err)
				}
				slog.InfoContext(ctx, "Contract fulfilled", "contract_id", contract.ID)
			} else {
				slog.InfoContext(ctx, "Contract not yet complete, continuing mining", "contract_id", contract.ID)
			}
		}
	}
	return nil
}
