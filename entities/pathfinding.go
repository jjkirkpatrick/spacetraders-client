package entities

import (
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

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
