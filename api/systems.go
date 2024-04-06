package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

// ListSystems retrieves a list of systems
func ListSystems(get GetFunc, meta *models.Meta) ([]*models.ListSystemsResponse, *models.Meta, *models.APIError) {
	endpoint := "/systems"

	var response struct {
		Data []*models.ListSystemsResponse `json:"data"`
		Meta models.Meta                   `json:"meta"`
	}

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

// GetSystem retrieves the details of a specific system
func GetSystem(get GetFunc, systemSymbol string) (*models.GetSystemResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s", systemSymbol)

	var response models.GetSystemResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ListWaypointsInSystem retrieves a list of waypoints in a specific system
func ListWaypointsInSystem(get GetFunc, meta *models.Meta, systemSymbol string, trait models.WaypointTrait, waypointType models.WaypointType) ([]*models.ListWaypointsResponse, *models.Meta, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints", systemSymbol)

	var response struct {
		Data []*models.ListWaypointsResponse `json:"data"`
		Meta models.Meta                     `json:"meta"`
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
func GetWaypoint(get GetFunc, systemSymbol, waypointSymbol string) (*models.GetWaypointResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s", systemSymbol, waypointSymbol)

	var response models.GetWaypointResponse
	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMarket retrieves the market details of a specific waypoint
func GetMarket(get GetFunc, systemSymbol, waypointSymbol string) (*models.GetMarketResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/market", systemSymbol, waypointSymbol)

	var response models.GetMarketResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetShipyard retrieves the shipyard details of a specific waypoint
func GetShipyard(get GetFunc, systemSymbol, waypointSymbol string) (*models.GetShipyardResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/shipyard", systemSymbol, waypointSymbol)

	var response models.GetShipyardResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetJumpGate retrieves the jump gate details of a specific waypoint
func GetJumpGate(get GetFunc, systemSymbol, waypointSymbol string) (*models.GetJumpGatesResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/jump-gate", systemSymbol, waypointSymbol)

	var response models.GetJumpGatesResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		apiErr := err
		return nil, apiErr
	}

	return &response, nil
}

// GetConstructionSite retrieves the construction site details of a specific waypoint
func GetConstructionSite(get GetFunc, systemSymbol, waypointSymbol string) (*models.GetConstructionSitesResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/construction", systemSymbol, waypointSymbol)

	var response models.GetConstructionSitesResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SupplyConstructionSite supplies a construction site with the required materials
func SupplyConstructionSite(post PostFunc, systemSymbol, waypointSymbol string, payload interface{}) (*models.SupplyConstructionSiteResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/systems/%s/waypoints/%s/construction/supply", systemSymbol, waypointSymbol)

	var response models.SupplyConstructionSiteResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
