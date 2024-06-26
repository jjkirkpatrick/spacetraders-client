package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

type listSystemsResponse struct {
	Data []*models.System `json:"data"`
	Meta models.Meta      `json:"meta"`
}

// ListSystems retrieves a list of systems
func ListSystems(get GetFunc, meta *models.Meta) ([]*models.System, *models.Meta, *models.APIError) {
	endpoint := "/systems"

	var response listSystemsResponse

	queryParams := map[string]string{
		"page":  fmt.Sprintf("%d", meta.Page),
		"limit": fmt.Sprintf("%d", meta.Limit),
	}

	err := get(endpoint, queryParams, &response)
	if err != nil {
		return nil, nil, err
	}

	return response.Data, &response.Meta, nil
}

type getSystemResponse struct {
	Data models.System `json:"data"`
	Meta models.Meta   `json:"meta"`
}

// GetSystem retrieves the details of a specific system
func GetSystem(get GetFunc, systemSymbol string) (*models.System, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s", systemSymbol)

	var response getSystemResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// ListWaypointsInSystem retrieves a list of waypoints in a specific system
func ListWaypointsInSystem(get GetFunc, meta *models.Meta, systemSymbol string, trait models.WaypointTrait, waypointType models.WaypointType) ([]*models.Waypoint, *models.Meta, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints", systemSymbol)

	var response struct {
		Data []*models.Waypoint `json:"data"`
		Meta models.Meta        `json:"meta"`
	}

	queryParams := map[string]string{
		"page":  fmt.Sprintf("%d", meta.Page),
		"limit": fmt.Sprintf("%d", meta.Limit),
	}

	if trait != "" {
		queryParams["traits"] = string(trait)
	}

	if waypointType != "" {
		queryParams["type"] = string(waypointType)
	}

	err := get(endpoint, queryParams, &response)
	if err != nil {
		return nil, nil, err
	}

	return response.Data, &response.Meta, nil
}

// GetWaypoint retrieves the details of a specific waypoint
func GetWaypoint(get GetFunc, systemSymbol, waypointSymbol string) (*models.Waypoint, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s", systemSymbol, waypointSymbol)

	var response struct {
		Data models.Waypoint `json:"data"`
		Meta models.Meta     `json:"meta"`
	}
	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// GetMarket retrieves the market details of a specific waypoint
func GetMarket(get GetFunc, systemSymbol, waypointSymbol string) (*models.Market, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/market", systemSymbol, waypointSymbol)

	var response struct {
		Data models.Market `json:"data"`
		Meta models.Meta   `json:"meta"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// GetShipyard retrieves the shipyard details of a specific waypoint
func GetShipyard(get GetFunc, systemSymbol, waypointSymbol string) (*models.Shipyard, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/shipyard", systemSymbol, waypointSymbol)

	var response struct {
		Data models.Shipyard `json:"data"`
		Meta models.Meta     `json:"meta"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// GetJumpGate retrieves the jump gate details of a specific waypoint
func GetJumpGate(get GetFunc, systemSymbol, waypointSymbol string) (*models.JumpGate, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/jump-gate", systemSymbol, waypointSymbol)

	var response struct {
		Data models.JumpGate `json:"data"`
		Meta models.Meta     `json:"meta"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		apiErr := err
		return nil, apiErr
	}

	return &response.Data, nil
}

// GetConstructionSite retrieves the construction site details of a specific waypoint
func GetConstructionSite(get GetFunc, systemSymbol, waypointSymbol string) (*models.ConstructionSite, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/construction", systemSymbol, waypointSymbol)

	var response struct {
		Data models.ConstructionSite `json:"data"`
		Meta models.Meta             `json:"meta"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// SupplyConstructionSite supplies a construction site with the required materials
func SupplyConstructionSite(post PostFunc, systemSymbol, waypointSymbol string, request models.SupplyConstructionSiteRequest) (*models.SupplyConstructionSiteResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/construction/supply", systemSymbol, waypointSymbol)

	var response models.SupplyConstructionSiteResponse

	err := post(endpoint, request, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func FindMarketsForGood(get GetFunc, systemSymbol string, goodSymbol string) ([]*models.Market, *models.APIError) {
	var allWaypoints []*models.Waypoint
	meta := &models.Meta{Page: 1, Limit: 20}
	for {
		waypoints, metaPtr, err := ListWaypointsInSystem(get, meta, systemSymbol, models.TraitMarketplace, "")
		if err != nil {
			return nil, err
		}
		allWaypoints = append(allWaypoints, waypoints...)
		if metaPtr.Page*metaPtr.Limit >= metaPtr.Total {
			break
		}
		meta.Page++
	}

	var marketsBuyingGood []*models.Market

	for _, waypoint := range allWaypoints {
		market, err := GetMarket(get, systemSymbol, waypoint.Symbol)
		if err != nil {
			continue // Skip waypoints where we can't get market data
		}

		for _, good := range market.Imports {
			if good.Symbol == models.GoodSymbol(goodSymbol) {
				marketsBuyingGood = append(marketsBuyingGood, &models.Market{
					Symbol:   waypoint.Symbol,
					Exports:  market.Exports,
					Imports:  market.Imports,
					Exchange: market.Exchange,
				})
				break // Found the good, no need to check other imports
			}
		}
	}

	return marketsBuyingGood, nil
}
