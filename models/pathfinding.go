package models

type PathfindingRoute struct {
	StartLocation string
	EndLocation   string
	Steps         []RouteStep
	TotalTime     int
}

type RouteStep struct {
	Waypoint   string
	FlightMode FlightMode
}

type Edge struct {
	Distance       float64
	FuelRequired   int
	TravelTime     int
	HasMarketplace bool
}

type Graph map[string]map[string]map[FlightMode]*Edge
