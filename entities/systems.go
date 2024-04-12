package entities

import (
	"math"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/internal/api"
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

type System struct {
	models.System
	client *client.Client
}

func ListSystems(c *client.Client) ([]*System, error) {
	fetchFunc := func(meta models.Meta) ([]*System, models.Meta, error) {
		metaPtr := &meta

		systems, metaPtr, err := api.ListSystems(c.Get, metaPtr)

		var convertedSystems []*System
		for _, modelSystem := range systems {
			convertedSystem := &System{
				System: *modelSystem, // Directly embed the modelContract
				client: c,
			}
			convertedSystems = append(convertedSystems, convertedSystem)
		}

		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
				metaPtr = &defaultMeta
			}
			return convertedSystems, *metaPtr, err.AsError()
		}
		if metaPtr != nil {

			return convertedSystems, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
			return convertedSystems, defaultMeta, nil
		}
	}
	return client.NewPaginator[*System](fetchFunc).FetchAllPages()
}

func GetSystem(c *client.Client, symbol string) (*System, error) {
	system, err := api.GetSystem(c.Get, symbol)
	if err != nil {
		return nil, err
	}

	systemEntity := &System{
		System: *system,
		client: c,
	}

	return systemEntity, nil
}

func (s *System) ListWaypoints(trait models.WaypointTrait, waypointType models.WaypointType) ([]*models.Waypoint, *models.Meta, error) {
	var allWaypoints []*models.Waypoint
	meta := models.Meta{Page: 1, Limit: 20, Total: 0}

	for {
		waypoints, _, err := api.ListWaypointsInSystem(s.client.Get, &meta, s.Symbol, trait, waypointType)
		if err != nil {
			return nil, nil, err
		}
		allWaypoints = append(allWaypoints, waypoints...)
		if len(waypoints) < meta.Limit {
			break
		}
		meta.Page++
	}

	return allWaypoints, &meta, nil
}

func (s *System) FetchWaypoint(symbol string) (*models.Waypoint, error) {
	waypoint, err := api.GetWaypoint(s.client.Get, s.Symbol, symbol)
	if err != nil {
		return nil, err
	}

	return waypoint, nil
}

func (s *System) GetWaypointsWithTrait(trait string, waypointType string) ([]*models.Waypoint, error) {
	waypoints, _, err := s.ListWaypoints(models.WaypointTrait(trait), models.WaypointType(waypointType))
	if err != nil {
		return nil, err
	}

	return waypoints, nil
}

func (s *System) GetMarket(waypointSymbol string) (*models.Market, error) {
	market, err := api.GetMarket(s.client.Get, s.Symbol, waypointSymbol)
	if err != nil {
		return nil, err
	}

	return market, nil
}

func (s *System) GetShipyard(waypointSymbol string) (*models.Shipyard, error) {
	shipyard, err := api.GetShipyard(s.client.Get, s.Symbol, waypointSymbol)
	if err != nil {
		return nil, err
	}

	return shipyard, nil
}

func (s *System) GetJumpGate(waypointSymbol string) (*models.JumpGate, error) {
	jumpGate, err := api.GetJumpGate(s.client.Get, s.Symbol, waypointSymbol)
	if err != nil {
		return nil, err
	}

	return jumpGate, nil
}

func (s *System) GetConstructionSite(waypointSymbol string) (*models.ConstructionSite, error) {
	projects, err := api.GetConstructionSite(s.client.Get, s.Symbol, waypointSymbol)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *System) SupplyConstructionSite(shipSymbol string, waypointSymbol string, good models.GoodSymbol, quantity int) error {
	payload := models.SupplyConstructionSiteRequest{
		ShipSymbol:  shipSymbol,
		TradeSymbol: good,
		Units:       quantity,
	}

	_, err := api.SupplyConstructionSite(s.client.Post, s.Symbol, waypointSymbol, payload)
	if err != nil {
		return err
	}

	return nil
}

// utiltiy functions

// Calculate the distance between two waypoints
func CalculateDistanceBetweenWaypoints(x1, y1, x2, y2 int) float64 {
	// Calculate Euclidean distance and round the result before returning
	return math.Round(math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2)))
}
