package entities

import (
	"container/heap"
	"math"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/internal/api"
	"github.com/jjkirkpatrick/spacetraders-client/models"
	"github.com/phuslu/log"
)

type Ship struct {
	models.Ship
	Client *client.Client
	Graph  models.Graph
}

func ListShips(c *client.Client) ([]*Ship, error) {
	fetchFunc := func(meta models.Meta) ([]*Ship, models.Meta, error) {
		metaPtr := &meta

		// Check if ships are in cache
		ships, metaPtr, err := api.ListShips(c.Get, metaPtr)

		var convertedShips []*Ship
		for _, modelShip := range ships {
			convertedShip := &Ship{
				Ship:   *modelShip, // Directly embed the modelShip
				Client: c,
			}
			graph, err := convertedShip.buildGraph()
			if err != nil {
				return nil, models.Meta{}, err
			}
			convertedShip.Graph = *graph
			convertedShips = append(convertedShips, convertedShip)
		}

		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
				metaPtr = &defaultMeta
			}
			return convertedShips, *metaPtr, err.AsError()
		}
		if metaPtr != nil {
			// Store ships in cache
			return convertedShips, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
			return convertedShips, defaultMeta, nil
		}
	}
	return client.NewPaginator[*Ship](fetchFunc).FetchAllPages()
}

func GetShip(c *client.Client, symbol string) (*Ship, error) {
	ship, err := api.GetShip(c.Get, symbol)
	if err != nil {
		return nil, err
	}

	shipEntity := &Ship{
		Ship:   *ship,
		Client: c,
	}

	graph, graphErr := shipEntity.buildGraph()
	if graphErr != nil {
		return nil, graphErr
	}
	shipEntity.Graph = *graph

	return shipEntity, nil
}

func PurchaseShip(c *client.Client, shipType string, waypoint string) (*models.Agent, *Ship, *models.Transaction, error) {
	purchaseShipRequest := &models.PurchaseShipRequest{
		ShipType:       models.ShipType(shipType),
		WaypointSymbol: waypoint,
	}

	response, err := api.PurchaseShip(c.Post, purchaseShipRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	shipEntity := &Ship{
		Ship:   response.Data.Ship,
		Client: c,
	}

	graph, graphErr := shipEntity.buildGraph()
	if graphErr != nil {
		return nil, nil, nil, graphErr
	}
	shipEntity.Graph = *graph

	c.CacheClient.Delete("all_ships")

	return &response.Data.Agent, shipEntity, &response.Data.Transaction, nil
}

func (s *Ship) Orbit() (*models.ShipNav, error) {
	//check if ship is already orbiting to avoid unnecessary API calls
	if s.Nav.Status == models.NavStatusInOrbit {
		return &s.Nav, nil
	}

	nav, err := api.OrbitShip(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Nav = *nav

	return nav, nil
}

func (s *Ship) Dock() (*models.ShipNav, error) {
	//check if ship is already docked to avoid unnecessary API calls
	if s.Nav.Status == models.NavStatusDocked {
		return &s.Nav, nil
	}

	nav, err := api.DockShip(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Nav = *nav

	return nav, nil
}

func (s *Ship) FetchCargo() (*models.Cargo, error) {
	cargo, err := api.GetShipCargo(s.Client.Get, s.Symbol)
	if err != nil {
		return nil, err
	}

	s.Cargo = *cargo

	return cargo, nil
}

func (s *Ship) Refine(produce string) (*models.Produced, *models.Consumed, error) {
	refineRequest := &models.RefineRequest{
		Produce: produce,
	}

	response, err := api.ShipRefine(s.Client.Post, s.Symbol, refineRequest)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Produced, &response.Data.Consumed, nil
}

func (s *Ship) Chart() (*models.Chart, *models.Waypoint, error) {
	nav, err := api.CreateChart(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	return &nav.Data.Chart, &nav.Data.Waypoint, nil
}

func (s *Ship) FetchCooldown() (*models.ShipCooldown, error) {
	cooldown, err := api.GetShipCooldown(s.Client.Get, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cooldown = *cooldown

	return cooldown, nil
}

func (s *Ship) Survey() ([]models.Survey, error) {
	response, err := api.CreateSurvey(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cooldown = response.Data.Cooldown

	return response.Data.Surveys, nil
}

func (s *Ship) Extract() (*models.Extraction, error) {
	response, err := api.ExtractResources(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Extraction, nil
}

func (s *Ship) Siphon() (*models.Extraction, error) {
	response, err := api.SiphonResources(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Extraction, nil
}

func (s *Ship) ExtractWithSurvey(survey models.Survey) (*models.Extraction, error) {
	extractWithSurveyRequest := &models.ExtractWithSurveyRequest{
		Signature:  survey.Signature,
		Symbol:     survey.Symbol,
		Deposits:   survey.Deposits,
		Expiration: survey.Expiration,
		Size:       survey.Size,
	}

	response, err := api.ExtractResourcesWithSurvey(s.Client.Post, s.Symbol, extractWithSurveyRequest)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Extraction, nil
}

func (s *Ship) Jettison(goodSymbol models.GoodSymbol, units int) (*models.Cargo, error) {
	jettisonRequest := &models.JettisonRequest{
		Symbol: goodSymbol,
		Units:  units,
	}

	response, err := api.JettisonCargo(s.Client.Post, s.Symbol, jettisonRequest)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Cargo, nil
}

func (s *Ship) Jump(systemSymbol string) (*models.ShipNav, *models.ShipCooldown, *models.Transaction, *models.Agent, error) {
	jumpRequest := &models.JumpShipRequest{
		WaypointSymbol: systemSymbol,
	}

	response, err := api.JumpShip(s.Client.Post, s.Symbol, jumpRequest)
	if err != nil {
		return nil, nil, nil, nil, err.AsError()
	}

	s.Nav = response.Data.Nav
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Nav, &response.Data.Cooldown, &response.Data.Transaction, &response.Data.Agent, nil
}

func (s *Ship) Navigate(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, []models.Event, error) {
	navigateRequest := &models.NavigateRequest{
		WaypointSymbol: waypointSymbol,
	}

	response, err := api.NavigateShip(s.Client.Post, s.Symbol, navigateRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	s.Fuel = response.Data.Fuel
	s.Nav = response.Data.Nav

	return &response.Data.Fuel, &response.Data.Nav, response.Data.Events, nil
}

func (s *Ship) SetFlightMode(flightmode models.FlightMode) error {
	flightModeRequest := &models.NavUpdateRequest{
		FlightMode: flightmode,
	}

	response, err := api.PatchShipNav(s.Client.Patch, s.Symbol, flightModeRequest)
	if err != nil {
		return err.AsError()
	}

	s.Nav.FlightMode = response.Data.FlightMode
	s.Nav.Status = response.Data.Status
	s.Nav.Route = response.Data.Route
	s.Nav.SystemSymbol = response.Data.SystemSymbol
	s.Nav.WaypointSymbol = response.Data.WaypointSymbol

	return nil
}

func (s *Ship) FetchNavigationStatus() (*models.ShipNav, error) {
	response, err := api.GetShipNav(s.Client.Get, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return response, nil
}

func (s *Ship) Warp(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, error) {
	warpRequest := &models.WarpRequest{
		WaypointSymbol: waypointSymbol,
	}

	response, err := api.WarpShip(s.Client.Post, s.Symbol, warpRequest)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Fuel = response.Data.Fuel
	s.Nav = response.Data.Nav

	return &response.Data.Fuel, &response.Data.Nav, nil
}

func (s *Ship) SellCargo(goodSymbol models.GoodSymbol, units int) (*models.Agent, *models.Cargo, *models.Transaction, error) {
	sellRequest := &models.SellCargoRequest{
		Symbol: goodSymbol,
		Units:  units,
	}

	response, err := api.SellCargo(s.Client.Post, s.Symbol, sellRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) ScanSystems() (*models.ShipCooldown, []models.System, error) {
	response, err := api.ScanSystems(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Cooldown = response.Data.Cooldown

	return &response.Data.Cooldown, response.Data.Systems, nil
}

func (s *Ship) ScanWaypoints() (*models.ShipCooldown, []models.Waypoint, error) {
	response, err := api.ScanWaypoints(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Cooldown = response.Data.Cooldown

	return &response.Data.Cooldown, response.Data.Waypoints, nil
}

func (s *Ship) Refuel(amount int, fromCargo bool) (*models.Agent, *models.FuelDetails, *models.Transaction, error) {
	refuelRequest := &models.RefuelShipRequest{
		FromCargo: fromCargo,
	}
	if amount == 0 {
		refuelRequest.Units = s.Fuel.Capacity
	} else {
		refuelRequest.Units = amount
	}

	response, err := api.RefuelShip(s.Client.Post, s.Symbol, refuelRequest)
	if err != nil {
		log.Error().Msgf("Error refueling ship %s: %v", s.Symbol, err.Data)
		return nil, nil, nil, err.AsError()
	}

	s.Fuel = response.Data.Fuel

	return &response.Data.Agent, &response.Data.Fuel, &response.Data.Transaction, nil
}

func (s *Ship) PurchaseCargo(goodSymbol models.GoodSymbol, units int) (*models.Agent, *models.Cargo, *models.Transaction, error) {
	purchaseRequest := &models.PurchaseCargoRequest{
		Symbol: goodSymbol,
		Units:  units,
	}

	response, err := api.PurchaseCargo(s.Client.Post, s.Symbol, purchaseRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) TransferCargo(goodSymbol models.GoodSymbol, units int, shipSymbol string) (*models.Cargo, error) {
	transferRequest := &models.TransferCargoRequest{
		TradeSymbol: goodSymbol,
		Units:       units,
		ShipSymbol:  shipSymbol,
	}

	response, err := api.TransferCargo(s.Client.Post, s.Symbol, transferRequest)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Cargo, nil
}

func (s *Ship) NegotiateContract() (*models.Contract, error) {

	response, err := api.NegotiateContract(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return &response.Data.Contract, nil
}

func (s *Ship) GetMounts() (*models.MountSymbol, string, string, int, []string, models.ShipRequirements, error) {
	response, err := api.GetMounts(s.Client.Get, s.Symbol)
	if err != nil {
		return nil, "", "", 0, nil, models.ShipRequirements{}, err.AsError()
	}

	return &response.Data.Symbol, response.Data.Name, response.Data.Description, response.Data.Strength, response.Data.Depsits, response.Data.Requirements, nil
}

func (s *Ship) InstallMount(mountSymbol models.MountSymbol) (*models.Agent, []models.ShipMount, *models.Cargo, *models.Transaction, error) {
	installRequest := &models.InstallMountRequest{
		Symbol: mountSymbol,
	}

	response, err := api.InstallMount(s.Client.Post, s.Symbol, installRequest)
	if err != nil {
		return nil, nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, response.Data.Mounts, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) RemoveMount(mountSymbol models.MountSymbol) (*models.Agent, []models.ShipMount, *models.Cargo, *models.Transaction, error) {
	removeRequest := &models.RemoveMountRequest{
		Symbol: mountSymbol,
	}

	response, err := api.RemoveMount(s.Client.Post, s.Symbol, removeRequest)
	if err != nil {
		return nil, nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, response.Data.Mounts, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) GetScrapPrice() (*models.Transaction, error) {
	response, err := api.GetScrapShip(s.Client.Get, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return &response.Data.Transaction, nil
}

func (s *Ship) ScrapShip() (*models.Transaction, error) {
	response, err := api.ScrapShip(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Client.CacheClient.Delete("all_ships")

	return &response.Data.Transaction, nil
}

func (s *Ship) GetRepairPrice() (*models.Transaction, error) {
	response, err := api.GetRepairShip(s.Client.Get, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return &response.Data.Transaction, nil
}

func (s *Ship) RepairShip() (*models.Ship, *models.Transaction, error) {
	response, err := api.RepairShip(s.Client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Ship = response.Data.Ship

	return &response.Data.Ship, &response.Data.Transaction, nil
}

func (s *Ship) GetRouteToDestination(destination string) (*models.PathfindingRoute, error) {
	log.Debug().Msgf("Getting route for ship %s", s.Symbol)

	// Find the optimal route using Dijkstra's algorithm
	steps, totalTime := s.findOptimalRoute(destination)
	return &models.PathfindingRoute{StartLocation: s.Nav.WaypointSymbol, EndLocation: destination, Steps: steps, TotalTime: totalTime}, nil
}

func (s *Ship) buildGraph() (*models.Graph, error) {
	log.Trace().Msgf("Building graph for ship %s", s.Symbol)

	// Attempt to retrieve the graph from cache first
	cachedGraph, found := s.Client.CacheClient.Get(s.Nav.SystemSymbol)
	if found {
		graph, ok := cachedGraph.(models.Graph)
		if ok {
			s.Graph = graph
			return &graph, nil
		}
	}

	// Retrieve the system and waypoints from cache or API
	system, err := s.getSystemFromCache()
	if err != nil {
		return nil, err
	}

	allWaypoints, err := s.getWaypointsFromCache(system)
	if err != nil {
		return nil, err
	}

	graph := make(models.Graph)

	// Build the graph
	for _, startWaypoint := range allWaypoints {
		// Create self-edge for each waypoint
		if _, ok := graph[startWaypoint.Symbol]; !ok {
			graph[startWaypoint.Symbol] = make(map[string]map[models.FlightMode]*models.Edge)
		}
		if _, ok := graph[startWaypoint.Symbol][startWaypoint.Symbol]; !ok {
			graph[startWaypoint.Symbol][startWaypoint.Symbol] = make(map[models.FlightMode]*models.Edge)
		}
		graph[startWaypoint.Symbol][startWaypoint.Symbol][models.FlightModeCruise] = &models.Edge{
			Distance:       0,
			FuelRequired:   0,
			TravelTime:     0,
			HasMarketplace: hasMarketplace(allWaypoints, startWaypoint.Symbol),
		}

		for _, endWaypoint := range allWaypoints {
			if startWaypoint.Symbol == endWaypoint.Symbol {
				continue
			}

			distance := CalculateDistanceBetweenWaypoints(startWaypoint.X, startWaypoint.Y, endWaypoint.X, endWaypoint.Y)

			for _, flightMode := range []models.FlightMode{models.FlightModeDrift, models.FlightModeCruise, models.FlightModeBurn} {
				hasMarketPlace := hasMarketplace(allWaypoints, endWaypoint.Symbol)

				fuelRequired := s.CalculateFuelRequired(distance, flightMode)
				if !hasMarketplace(allWaypoints, endWaypoint.Symbol) {
					fuelRequired *= 2
				}
				travelTime := s.CalculateTravelTime(distance, flightMode)

				if _, ok := graph[startWaypoint.Symbol]; !ok {
					graph[startWaypoint.Symbol] = make(map[string]map[models.FlightMode]*models.Edge)
				}
				if _, ok := graph[startWaypoint.Symbol][endWaypoint.Symbol]; !ok {
					graph[startWaypoint.Symbol][endWaypoint.Symbol] = make(map[models.FlightMode]*models.Edge)
				}

				graph[startWaypoint.Symbol][endWaypoint.Symbol][flightMode] = &models.Edge{
					Distance:       distance,
					FuelRequired:   fuelRequired,
					TravelTime:     travelTime,
					HasMarketplace: hasMarketPlace,
				}
			}
		}
	}

	s.Graph = graph
	// Cache the newly built graph
	s.Client.CacheClient.Set(s.Nav.SystemSymbol, graph, 0)

	return &graph, nil
}

// Helper functions

func (s *Ship) getSystemFromCache() (*System, error) {
	cachedSystem, found := s.Client.CacheClient.Get("system_" + s.Nav.SystemSymbol)
	if found {
		system, _ := cachedSystem.(*System)
		return system, nil
	}

	system, err := GetSystem(s.Client, s.Nav.SystemSymbol)
	if err != nil {
		return nil, err
	}
	s.Client.CacheClient.Set("system_"+s.Nav.SystemSymbol, system, 0)

	return system, nil
}

func (s *Ship) getWaypointsFromCache(system *System) ([]*models.Waypoint, error) {
	cachedWaypoints, found := s.Client.CacheClient.Get("waypoints_" + s.Nav.SystemSymbol)
	if found {
		allWaypoints, _ := cachedWaypoints.([]*models.Waypoint)
		return allWaypoints, nil
	}
	allWaypoints, _, err := system.ListWaypoints("", "")
	if err != nil {
		return nil, err
	}
	s.Client.CacheClient.Set("waypoints_"+s.Nav.SystemSymbol, allWaypoints, 0)

	return allWaypoints, nil
}

func (s *Ship) CalculateFuelRequired(distance float64, flightMode models.FlightMode) int {
	var fuel float64
	switch flightMode {
	case models.FlightModeDrift:
		fuel = 1
	case models.FlightModeCruise:
		fuel = math.Round(distance)
	case models.FlightModeBurn:
		fuel = math.Max(2, 2*math.Round(distance))
	default:
		fuel = math.Round(distance)
	}
	return int(fuel)
}

func (s *Ship) CalculateTravelTime(distance float64, flightMode models.FlightMode) int {
	var multiplier float64
	switch flightMode {
	case models.FlightModeCruise:
		multiplier = 25
	case models.FlightModeDrift:
		multiplier = 250
	case models.FlightModeBurn:
		multiplier = 12.5
	default:
		multiplier = 25
	}
	travelTime := math.Round(math.Round(math.Max(1, distance))*(multiplier/float64(s.Engine.Speed)) + 15)
	return int(travelTime)
}

func (s *Ship) findOptimalRoute(destination string) ([]models.RouteStep, int) {

	//check if the ship has a 0 fuel capacity if so return a path to drift to the destination
	if s.Fuel.Capacity == 0 {
		return []models.RouteStep{{
			Waypoint:     destination,
			FlightMode:   models.FlightModeDrift,
			ShouldRefuel: false,
		}}, 0
	}

	// Create a map to store the shortest distance to each waypoint
	shortestDistances := make(map[string]int)
	for waypoint := range s.Graph {
		shortestDistances[waypoint] = math.MaxInt32
	}
	shortestDistances[s.Nav.WaypointSymbol] = 0

	// Create a map to store the previous waypoint in the shortest path
	previous := make(map[string]string)

	// Create a map to store the flight mode used to reach each waypoint
	flightModes := make(map[string]models.FlightMode)

	// Create a priority queue to store waypoints to visit
	pq := make(PriorityQueue, 0)
	pq = append(pq, &Item{
		value:    s.Nav.WaypointSymbol,
		priority: 0,
	})

	for len(pq) > 0 {
		current := heap.Pop(&pq).(*Item).value

		// If we have reached the destination waypoint, we can stop searching
		if current == destination {
			break
		}

		// Explore neighboring waypoints
		for neighbor, edges := range s.Graph[current] {
			// Skip waypoints without a marketplace unless it's the destination
			if neighbor != destination {
				if neighborEdges, ok := s.Graph[neighbor][neighbor]; ok {
					if edge, ok := neighborEdges[models.FlightModeCruise]; ok && edge != nil {
						if !edge.HasMarketplace {
							continue
						}
					}
				}
			}

			for flightMode, edge := range edges {
				fuelRequired := edge.FuelRequired
				travelTime := edge.TravelTime

				// Check if the ship has enough fuel to reach the neighbor waypoint
				if s.Fuel.Current >= fuelRequired {
					tentativeDistance := shortestDistances[current] + travelTime

					if tentativeDistance < shortestDistances[neighbor] {
						shortestDistances[neighbor] = tentativeDistance
						previous[neighbor] = current
						flightModes[neighbor] = flightMode

						heap.Push(&pq, &Item{
							value:    neighbor,
							priority: tentativeDistance,
						})
					}
				}
			}
		}
	}

	path := []models.RouteStep{}
	current := destination
	totalTime := 0

	for current != s.Nav.WaypointSymbol {
		previousWaypoint := previous[current]
		shouldRefuel := false

		if edges, ok := s.Graph[current][current]; ok {
			if edge, ok := edges[models.FlightModeCruise]; ok && edge != nil {
				shouldRefuel = edge.HasMarketplace
			}
		}

		// Check if the flight mode is set and has a valid edge
		if flightMode, ok := flightModes[current]; ok {
			if edge, ok := s.Graph[previousWaypoint][current][flightMode]; ok && edge != nil {
				path = append([]models.RouteStep{{
					Waypoint:     current,
					FlightMode:   flightMode,
					ShouldRefuel: shouldRefuel,
				}}, path...)

				totalTime += edge.TravelTime
			}
		}

		current = previousWaypoint
	}

	return path, totalTime
}
