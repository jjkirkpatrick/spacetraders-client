package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/jjkirkpatrick/spacetraders-client/models"
	"github.com/phuslu/log"
)

type GameState struct {
	Agent      *entities.Agent
	HomeSystem string
	Contracts  []*entities.Contract `json:"contracts"`
	Ships      []*entities.Ship     `json:"ships"`
}

func main() {

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
	options.Symbol = "BlueCaloria-1"
	options.Faction = "COSMIC"
	options.LogLevel = log.InfoLevel

	gameState := &GameState{}

	client, cerr := client.NewClient(options)
	if cerr != nil {
		log.Fatal().Msgf("Failed to create client: %v", cerr)
	}

	client.ConfigureMetricsClient(
		"http://192.168.1.33:8086",
		"238nUuJVX9CzDqdsU7wvINk8ByIG-3MykZqwUSTEwBIeLgBKNTbwV8x_lik_4t1oXSTfj8OqRPwzmPS8y3tsdg==",
		"spacetraders",
		"spacetraders",
	)

	agent, err := entities.GetAgent(client)
	if err != nil {
		log.Fatal().Msgf("Failed to get agent: %v", err)
	}
	gameState.Agent = agent
	gameState.HomeSystem = getSystemNameFromHomeSystem(gameState.Agent)

	// Get all contracts
	contracts, err := entities.ListContracts(client)
	if err != nil {
		log.Fatal().Msgf("Failed to list contracts: %v", err)
	}

	gameState.Contracts = contracts

	for _, contract := range gameState.Contracts {
		log.Info().Msgf("Contract: %s", contract.ID)
		agent, _, err := contract.Accept()
		if err != nil && strings.Contains(err.Error(), "has already been accepted") {
			continue
		}
		gameState.Agent = agent
		if err != nil {
			log.Fatal().Msgf("Failed to accept contract: %v", err)
		}
	}

	// Get active contract terms
	activeContractTrms := gameState.getActiveContracts()

	// Print active contract details
	gameState.printActiveContractDetails(activeContractTrms)

	currentSystem, err := entities.GetSystem(client, gameState.HomeSystem)
	if err != nil {
		log.Fatal().Msgf("Failed to get system: %v", err)
	}

	//shipyards, err := currentSystem.GetWaypointsWithTrait("SHIPYARD", "")
	// Find the engineered asteroid in the system
	astroid, err := currentSystem.GetWaypointsWithTrait("", "ENGINEERED_ASTEROID")
	if err != nil {
		log.Fatal().Msgf("Failed to find engineered asteroid: %v", err)
	}
	// Get mining ship
	miningShipSymbol, err := gameState.getMiningShip(client)
	if err != nil {
		log.Info().Msgf("No mining ship found, attempting to find a shipyard: %v", err)
		shipyards, err := currentSystem.GetWaypointsWithTrait("SHIPYARD", "")
		if err != nil || len(shipyards) == 0 {
			log.Fatal().Msgf("Failed to find shipyard: %v", err)
		}
		var shipyard *models.Shipyard
		for _, shipyardInfo := range shipyards {
			tempShipyard, err := currentSystem.GetShipyard(shipyardInfo.Symbol)
			if err != nil {
				log.Info().Msgf("Failed to get shipyard details for %s: %v", shipyardInfo.Symbol, err)
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
			log.Fatal().Msg("Failed to find a shipyard selling SHIP_MINING_DRONE")
		}
		_, ship, _, err := entities.PurchaseShip(client, "SHIP_MINING_DRONE", shipyard.Symbol)
		if err != nil {
			log.Fatal().Msgf("Failed to purchase mining ship: %v", err)
		}
		miningShipSymbol = ship.Symbol
		log.Info().Msgf("Purchased mining ship: %s", miningShipSymbol)
	}

	ship, err := entities.GetShip(client, miningShipSymbol)
	if err != nil {
		log.Fatal().Msgf("Failed to get ship: %v", err)
	}

	// Mine resources and deliver contracts
	err = gameState.mineAndDeliverContracts(ship, astroid[0], activeContractTrms)
	if err != nil {
		log.Fatal().Msgf("Failed to mine and deliver contracts: %v", err)
	}

	log.Info().Msg("Program completed successfully.")
}

func waitForCooldown(ship *entities.Ship) {
	_, cerr := ship.FetchCooldown()
	if cerr != nil {
		log.Fatal().Msgf("Failed to get ship cooldown: %v", cerr)
	}

	if ship.Cooldown.RemainingSeconds > 0 {
		log.Info().Msgf("Waiting for cooldown to finish Remaining Seconds: %d", ship.Cooldown.RemainingSeconds)
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
		log.Info().Msgf("Contract: %s", contract.ID)
		for _, deliver := range contract.Terms.Deliver {
			log.Info().Msgf("Deliver: %s %s %d %d", deliver.TradeSymbol, deliver.DestinationSymbol, deliver.UnitsRequired, deliver.UnitsFulfilled)
		}
	}
}

func (gs *GameState) getMiningShip(client *client.Client) (string, error) {
	log.Info().Msg("Getting mining ship")
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
	log.Info().Msg("Mining and delivering contracts")
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
			log.Info().Msg("All contracts fulfilled.")
			break
		}

		//Navigate to the asteroid
		log.Info().Msg("Navigating to asteroid")
		err := gs.navigateToWaypoint(miningShip, *astroid)
		if err != nil {
			return fmt.Errorf("failed to navigate to asteroid: %v", err)
		}

		log.Info().Msg("Mining resources")

		_, err = miningShip.Orbit()
		if err != nil {
			return fmt.Errorf("failed to orbit ship: %v", err)
		}

		// Mine resources
		err = gs.mineResources(miningShip, activeContractTrms)
		if err != nil {
			return fmt.Errorf("failed to mine resources: %v", err)
		}

		log.Info().Msg("Selling trade goods")
		// Sell trade goods not required for contracts
		err = gs.jettisonUnwantedCargo(miningShip, activeContractTrms)
		if err != nil {
			return fmt.Errorf("failed to sell trade goods: %v", err)
		}

		log.Info().Msg("Delivering contract goods")

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
		log.Info().Msgf("Navigating to waypoint %s", step.Waypoint)
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
		log.Info().Msgf("Navigating to %s. Arrival Time: %s, using Flight mode:  %s\n", waypointSymbol.Symbol, arrivalTime, step.FlightMode)

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
	log.Info().Msgf("Initiating refueling process for ship: %s", miningShip.Symbol)
	_, err := miningShip.Dock()
	if err != nil {
		return fmt.Errorf("failed to dock ship: %v", err)
	}

	_, _, _, err = miningShip.Refuel(0, false)
	if err != nil {
		return fmt.Errorf("failed to refuel ship: %v", err)
	}
	log.Info().Msgf("Refueling completed successfully for ship: %s", miningShip.Symbol)

	return nil
}

func (gs *GameState) mineResources(miningShip *entities.Ship, activeContractTrms []entities.Contract) error {
	cargo, err := miningShip.FetchCargo()
	if err != nil {
		return fmt.Errorf("failed to get ship cargo: %v", err)
	}

	for _, contract := range activeContractTrms {
		log.Info().Msgf("Checking contract requirements for contract ID: %s", contract.ID)
		for _, deliver := range contract.Terms.Deliver {
			log.Info().Msgf("Requirement: %d units of %s to be delivered to %s", deliver.UnitsRequired, deliver.TradeSymbol, deliver.DestinationSymbol)
			unitsAvailable := 0
			for _, cargoItem := range cargo.Inventory {
				if cargoItem.Symbol == deliver.TradeSymbol {
					unitsAvailable += cargoItem.Units
				}
			}
			if unitsAvailable >= deliver.UnitsRequired-deliver.UnitsFulfilled {
				log.Info().Msgf("Cargo contains enough %s to fulfill the contract requirement.", deliver.TradeSymbol)
			} else {
				for unitsAvailable < deliver.UnitsRequired {
					log.Info().Msgf("Not enough %s in cargo to fulfill the contract requirement. Required: %d, Available: %d", deliver.TradeSymbol, deliver.UnitsRequired-activeContractTrms[0].Terms.Deliver[0].UnitsFulfilled, unitsAvailable)
					if cargo.Units >= cargo.Capacity {
						log.Info().Msg("Cargo hold is full. Unable to extract more resources.")
						break
					}
					log.Info().Msg("Attempting to extract more resources...")
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
							log.Info().Msgf("Jettisoning %s as it is not required for the current contract.", item.Name)
							_, jettisonErr := miningShip.Jettison(models.GoodSymbol(item.Symbol), item.Units)
							if jettisonErr != nil {
								return fmt.Errorf("failed to jettison %s: %v", item.Name, jettisonErr)
							}
							log.Info().Msgf("Successfully jettisoned %s.", item.Name)
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

					log.Info().Msgf("Cargo Summary:")
					log.Info().Msgf("Capacity: %d, Units: %d", cargo.Capacity, cargo.Units)
					log.Info().Msg("Inventory:")
					for _, item := range cargo.Inventory {
						log.Info().Msgf("- %s: %d units", item.Name, item.Units)
					}
				}
				if unitsAvailable >= deliver.UnitsRequired {
					log.Info().Msgf("Now have enough %s to fulfill the contract requirement. Required: %d, Available: %d", deliver.TradeSymbol, deliver.UnitsRequired-activeContractTrms[0].Terms.Deliver[0].UnitsFulfilled, unitsAvailable)
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

	log.Info().Msg("Starting to sell trade goods not required for the active contract.")

	for _, item := range cargo.Inventory {
		// Check if the item is not part of the active contract requirements
		if !gs.isItemRequiredForContracts(models.GoodSymbol(item.Symbol), activeContractTrms) {
			log.Info().Msgf("No markets found buying %s, jettisoning cargo.", item.Name)
			_, jettisonErr := miningShip.Jettison(models.GoodSymbol(item.Symbol), item.Units)
			if jettisonErr != nil {
				return fmt.Errorf("failed to jettison %s: %v", item.Name, jettisonErr)
			}
			log.Info().Msgf("Successfully jettisoned %s.", item.Name)
		}
	}

	log.Info().Msg("Finished dumping trade goods.")
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
				log.Info().Msgf("Delivering goods for contract ID: %s", contract.ID)
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
				log.Info().Msgf("Delivered %d units of %s for contract ID: %s", deliver.UnitsRequired-deliver.UnitsFulfilled, deliver.TradeSymbol, contract.ID)
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
				log.Info().Msgf("Fulfilled contract ID: %s", contract.ID)
			} else {
				log.Info().Msgf("Cannot fulfill contract ID: %s as not all deliveries have been made.", contract.ID)
			}
		}
	}
	return nil
}
