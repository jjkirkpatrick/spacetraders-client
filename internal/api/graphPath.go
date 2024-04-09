package api

import (
	"container/heap"
	"math"

	"github.com/jjkirkpatrick/spacetraders-client/internal/models"
)

func GetRouteToDestination(get GetFunc, graph models.Graph, destination string, shipSymbol string) (*models.Route, error) {
	// Get ship details
	ship, err := GetShip(get, shipSymbol)
	if err != nil {
		return nil, err
	}

	// Get current ship location
	startLocation := ship.Nav.WaypointSymbol

	allWaypoints, aerr := getAllWaypoints(get, ship.Nav.SystemSymbol)
	if aerr != nil {
		return nil, aerr
	}

	// Find the optimal route using Dijkstra's algorithm
	steps, totalTime := findOptimalRoute(graph, allWaypoints, startLocation, destination, ship.Fuel.Current, ship.Fuel.Capacity)

	return &models.Route{StartLocation: startLocation, EndLocation: destination, Steps: steps, TotalTime: totalTime}, nil
}

func getAllWaypoints(get GetFunc, systemSymbol string) ([]*models.Waypoint, error) {
	var allWaypoints []*models.Waypoint

	meta := &models.Meta{Page: 1, Limit: 20}
	for {
		waypoints, metaPtr, err := ListWaypointsInSystem(get, meta, systemSymbol, "", "")
		if err != nil {
			return nil, err
		}
		for _, waypoint := range waypoints {
			allWaypoints = append(allWaypoints, &models.Waypoint{
				Symbol: waypoint.Symbol,
				X:      waypoint.X,
				Y:      waypoint.Y,
				Traits: waypoint.Traits,
			})
		}
		if metaPtr.Page*metaPtr.Limit >= metaPtr.Total {
			break
		}
		meta.Page++
	}

	return allWaypoints, nil
}

func BuildGraph(get GetFunc, systemSymbol string, engineSpeed int) (models.Graph, error) {
	graph := make(models.Graph)

	allWaypoints, err := getAllWaypoints(get, systemSymbol)
	if err != nil {
		return nil, err
	}

	for _, startWaypoint := range allWaypoints {
		for _, endWaypoint := range allWaypoints {
			if startWaypoint.Symbol == endWaypoint.Symbol {
				continue
			}

			distance := calculateDistance(startWaypoint.X, startWaypoint.Y, endWaypoint.X, endWaypoint.Y)

			for _, flightMode := range []models.FlightMode{models.FlightModeDrift, models.FlightModeCruise, models.FlightModeBurn} {
				if isWithinRange(distance, flightMode, 100) {
					fuelRequired := calculateFuelRequired(distance, flightMode)
					travelTime := calculateTravelTime(distance, flightMode, engineSpeed)

					if _, ok := graph[startWaypoint.Symbol]; !ok {
						graph[startWaypoint.Symbol] = make(map[string]map[models.FlightMode]*models.Edge)
					}
					if _, ok := graph[startWaypoint.Symbol][endWaypoint.Symbol]; !ok {
						graph[startWaypoint.Symbol][endWaypoint.Symbol] = make(map[models.FlightMode]*models.Edge)
					}

					graph[startWaypoint.Symbol][endWaypoint.Symbol][flightMode] = &models.Edge{
						Distance:     distance,
						FuelRequired: fuelRequired,
						TravelTime:   travelTime,
					}
				}
			}
		}
	}

	return graph, nil
}

func hasMarketplaceTrait(waypoint *models.Waypoint) bool {
	for _, trait := range waypoint.Traits {
		if trait.Symbol == models.TraitMarketplace {
			return true
		}
	}
	return false
}

func isWithinRange(distance float64, flightMode models.FlightMode, fuelCapacity int) bool {
	var fuelCost float64
	switch flightMode {
	case models.FlightModeCruise:
		fuelCost = math.Round(distance)
	case models.FlightModeDrift:
		fuelCost = 1
	case models.FlightModeBurn:
		fuelCost = math.Max(2, 2*math.Round(distance))
	default:
		fuelCost = math.Round(distance) // Default to CRUISE mode if flight mode is unknown
	}
	return int(fuelCost) <= fuelCapacity
}

func calculateDistance(x1, y1, x2, y2 int) float64 {
	// Calculate Euclidean distance and round the result before returning
	return math.Round(math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2)))
}

func calculateFuelRequired(distance float64, flightMode models.FlightMode) int {
	var fuel float64
	switch flightMode {
	case models.FlightModeDrift:
		fuel = 1 // Drift mode always incurs a cost of 1
	case models.FlightModeCruise:
		fuel = math.Round(distance) // Cruise mode rounds the distance
	case models.FlightModeBurn:
		fuel = math.Max(2, 2*math.Round(distance)) // Burn mode doubles the rounded distance, minimum cost of 2
	default:
		fuel = math.Round(distance) // Default to rounding the distance if flight mode is unknown
	}

	return int(fuel)
}

func calculateTravelTime(distance float64, flightMode models.FlightMode, engineSpeed int) int {
	var multiplier float64
	switch flightMode {
	case models.FlightModeCruise:
		multiplier = 25
	case models.FlightModeDrift:
		multiplier = 250
	case models.FlightModeBurn:
		multiplier = 12.5
	default:
		multiplier = 25 // Default to Cruise mode if flight mode is unknown
	}

	travelTime := math.Round(math.Round(math.Max(1, distance))*(multiplier/float64(engineSpeed)) + 15)
	return int(travelTime)
}

func findOptimalRoute(graph models.Graph, allWaypoints []*models.Waypoint, start, end string, currentFuel, fuelCapacity int) ([]models.RouteStep, int) {
	// Create a map to store the shortest distance to each waypoint
	shortestDistances := make(map[string]int)
	for waypoint := range graph {
		shortestDistances[waypoint] = math.MaxInt32
	}
	shortestDistances[start] = 0

	// Create a map to store the previous waypoint in the shortest path
	previous := make(map[string]string)

	// Create a priority queue to store waypoints to visit
	pq := &PriorityQueue{}
	heap.Push(pq, &Item{
		value:    start,
		priority: 0,
	})

	flightModes := make(map[string]models.FlightMode)
	fuelLevels := make(map[string]int)
	fuelLevels[start] = currentFuel

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		current := item.value

		// If we have reached the end waypoint, we can stop searching
		if current == end {
			break
		}

		// Explore neighboring waypoints
		for neighbor, edges := range graph[current] {
			bestFlightMode := models.FlightModeDrift
			bestTravelTime := math.MaxInt32

			for flightMode, edge := range edges {
				// Calculate the fuel required to reach the neighbor using the current flight mode
				fuelToNeighbor := edge.FuelRequired

				// Check if there is enough fuel to reach the neighbor using the current flight mode
				if fuelLevels[current] >= fuelToNeighbor {
					// Calculate the tentative distance to the neighbor through the current waypoint and flight mode
					tentativeDistance := shortestDistances[current] + edge.TravelTime

					// If the tentative distance is shorter than the current shortest distance to the neighbor,
					// update the shortest distance, the previous waypoint, and the best flight mode
					if tentativeDistance < shortestDistances[neighbor] {
						shortestDistances[neighbor] = tentativeDistance
						previous[neighbor] = current
						bestFlightMode = flightMode
						bestTravelTime = tentativeDistance
						fuelLevels[neighbor] = fuelLevels[current] - fuelToNeighbor
					} else if tentativeDistance == shortestDistances[neighbor] {
						// If the tentative distance is the same as the current shortest distance,
						// prioritize paths through waypoints with a market
						if hasMarketplace(allWaypoints, neighbor) && !hasMarketplace(allWaypoints, previous[neighbor]) {
							previous[neighbor] = current
							bestFlightMode = flightMode
							bestTravelTime = tentativeDistance
							fuelLevels[neighbor] = fuelLevels[current] - fuelToNeighbor
						} else if flightMode == models.FlightModeCruise {
							// If both waypoints have a market or neither have a market,
							// prioritize CRUISE flight mode over DRIFT
							bestFlightMode = flightMode
							bestTravelTime = tentativeDistance
							fuelLevels[neighbor] = fuelLevels[current] - fuelToNeighbor
						}
					}
				}
			}

			if bestTravelTime != math.MaxInt32 {
				heap.Push(pq, &Item{
					value:    neighbor,
					priority: bestTravelTime,
				})
				flightModes[neighbor] = bestFlightMode
			}
		}

		// Refuel at the current waypoint if it has a market
		if hasMarketplace(allWaypoints, current) {
			fuelLevels[current] = fuelCapacity
		}
	}

	// Reconstruct the shortest path from start to end
	path := []models.RouteStep{}
	current := end
	totalTime := shortestDistances[end]
	for current != start {
		path = append([]models.RouteStep{{Waypoint: current, FlightMode: flightModes[current]}}, path...)
		current = previous[current]
	}

	return path, totalTime
}

func hasMarketplace(allWaypoints []*models.Waypoint, waypointSymbol string) bool {
	for _, waypoint := range allWaypoints {
		if waypoint.Symbol == waypointSymbol {
			for _, trait := range waypoint.Traits {
				if trait.Symbol == models.TraitMarketplace {
					return true
				}
			}
			break
		}
	}
	return false
}

// Item represents an item in the priority queue
type Item struct {
	value    string
	priority int
}

// PriorityQueue represents a priority queue of items
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
