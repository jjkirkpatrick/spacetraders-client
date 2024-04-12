package api

import (
	"container/heap"
	"math"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

func FindOptimalRoute(graph models.Graph, allWaypoints []*models.Waypoint, start, end string, currentFuel, fuelCapacity int) ([]models.RouteStep, int) {
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
