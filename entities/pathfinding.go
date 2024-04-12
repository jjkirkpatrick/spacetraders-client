package entities

import (
	"container/heap"
	"math"

	"github.com/jjkirkpatrick/spacetraders-client/models"
	"github.com/phuslu/log"
)

func findOptimalRoute(ship *Ship, allWaypoints []*models.Waypoint, destination string) ([]models.RouteStep, int) {
	log.Debug().Msgf("Finding optimal route from %s to %s with %d fuel and %d fuel capacity", ship.Nav.WaypointSymbol, destination, ship.Fuel.Current, ship.Fuel.Capacity)

	// Create a map to store the shortest distance to each waypoint
	shortestDistances := make(map[string]int)
	for waypoint := range ship.Graph {
		shortestDistances[waypoint] = math.MaxInt32
	}
	shortestDistances[ship.Nav.WaypointSymbol] = 0

	// Create a map to store the previous waypoint in the shortest path
	previous := make(map[string]string)

	// Create a priority queue to store waypoints to visit
	pq := &PriorityQueue{}
	heap.Push(pq, &Item{
		value:    ship.Nav.WaypointSymbol,
		priority: 0,
	})

	flightModes := make(map[string]models.FlightMode)
	fuelLevels := make(map[string]int)
	fuelLevels[ship.Nav.WaypointSymbol] = ship.Fuel.Current

	visited := make(map[string]bool)

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		current := item.value
		log.Debug().Msgf("Current waypoint: %s", current)

		if visited[current] {
			log.Trace().Msgf("Waypoint %s already visited, skipping", current)
			continue
		}
		visited[current] = true

		// If we have reached the end waypoint, we can stop searching
		if current == destination {
			log.Trace().Msgf("Reached end waypoint %s, stopping search", destination)
			break
		}

		log.Trace().Msgf("Exploring neighbors of waypoint %s", current)
		// Explore neighboring waypoints
		for neighbor, edges := range ship.Graph[current] {
			log.Trace().Msgf("Checking neighbor waypoint %s", neighbor)
			bestFlightMode := models.FlightModeDrift
			bestTravelTime := math.MaxInt32

			for flightMode, edge := range edges {
				log.Trace().Msgf("Checking flight mode %s to neighbor %s", flightMode, neighbor)
				// Calculate the fuel required to reach the neighbor using the current flight mode
				fuelToNeighbor := edge.FuelRequired

				// Check if there is enough fuel to reach the neighbor using the current flight mode
				if fuelLevels[current] >= fuelToNeighbor {
					log.Trace().Msgf("Enough fuel (%d) to reach neighbor %s using flight mode %s (requires %d fuel)", fuelLevels[current], neighbor, flightMode, fuelToNeighbor)
					// Calculate the tentative distance to the neighbor through the current waypoint and flight mode
					tentativeDistance := shortestDistances[current] + edge.TravelTime

					// If the tentative distance is shorter than the current shortest distance to the neighbor,
					// update the shortest distance, the previous waypoint, and the best flight mode
					if tentativeDistance < shortestDistances[neighbor] {
						log.Trace().Msgf("Found shorter path to neighbor %s through waypoint %s using flight mode %s (tentative distance: %d, current shortest: %d)", neighbor, current, flightMode, tentativeDistance, shortestDistances[neighbor])
						shortestDistances[neighbor] = tentativeDistance
						previous[neighbor] = current
						bestFlightMode = flightMode
						bestTravelTime = tentativeDistance
						fuelLevels[neighbor] = fuelLevels[current] - fuelToNeighbor
					} else if tentativeDistance == shortestDistances[neighbor] {
						log.Trace().Msgf("Found path to neighbor %s through waypoint %s using flight mode %s with same distance as current shortest (%d)", neighbor, current, flightMode, tentativeDistance)
						// If the tentative distance is the same as the current shortest distance,
						// prioritize paths through waypoints with a market
						if hasMarketplace(allWaypoints, neighbor) && !hasMarketplace(allWaypoints, previous[neighbor]) {
							log.Trace().Msgf("Prioritizing path to neighbor %s through waypoint %s because it has a marketplace and previous waypoint %s does not", neighbor, current, previous[neighbor])
							previous[neighbor] = current
							bestFlightMode = flightMode
							bestTravelTime = tentativeDistance
							fuelLevels[neighbor] = fuelLevels[current] - fuelToNeighbor
						} else if flightMode == models.FlightModeCruise {
							log.Trace().Msgf("Prioritizing CRUISE flight mode over DRIFT for path to neighbor %s through waypoint %s", neighbor, current)
							// If both waypoints have a market or neither have a market,
							// prioritize CRUISE flight mode over DRIFT
							bestFlightMode = flightMode
							bestTravelTime = tentativeDistance
							fuelLevels[neighbor] = fuelLevels[current] - fuelToNeighbor
						}
					}
				} else {
					log.Trace().Msgf("Not enough fuel (%d) to reach neighbor %s using flight mode %s (requires %d fuel)", fuelLevels[current], neighbor, flightMode, fuelToNeighbor)
				}
			}

			if bestTravelTime != math.MaxInt32 && !visited[neighbor] {
				log.Trace().Msgf("Adding neighbor %s to priority queue with best travel time %d using flight mode %s", neighbor, bestTravelTime, bestFlightMode)
				heap.Push(pq, &Item{
					value:    neighbor,
					priority: bestTravelTime,
				})
				flightModes[neighbor] = bestFlightMode
			}
		}

		// Refuel at the current waypoint if it has a market
		if hasMarketplace(allWaypoints, current) {
			log.Trace().Msgf("Refueling at waypoint %s with marketplace, setting fuel to max capacity %d", current, ship.Fuel.Capacity)
			fuelLevels[current] = ship.Fuel.Capacity
		} else {
			log.Trace().Msgf("Waypoint %s does not have a marketplace, skipping refuel", current)
		}
	}
	// Reconstruct the shortest path from start to end
	path := []models.RouteStep{}
	current := destination
	totalTime := shortestDistances[destination]
	for current != ship.Nav.WaypointSymbol {
		path = append([]models.RouteStep{{Waypoint: current, FlightMode: flightModes[current]}}, path...)
		current = previous[current]
	}
	log.Debug().Msgf("Optimal route found: %v", path)
	return path, totalTime
}

func hasMarketplace(allWaypoints []*models.Waypoint, waypointSymbol string) bool {
	log.Debug().Msgf("Checking for marketplace at waypoint %s", waypointSymbol)
	for _, waypoint := range allWaypoints {
		if waypoint.Symbol == waypointSymbol {
			log.Debug().Msgf("Found waypoint %s, checking for marketplace trait", waypointSymbol)
			for _, trait := range waypoint.Traits {
				if trait.Symbol == models.TraitMarketplace {
					log.Debug().Msgf("Marketplace found at waypoint %s", waypointSymbol)
					return true
				}
			}
			log.Debug().Msgf("No marketplace found at waypoint %s", waypointSymbol)
			break
		}
	}
	log.Debug().Msgf("Waypoint %s not found in allWaypoints", waypointSymbol)
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
