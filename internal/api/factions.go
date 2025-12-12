package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

func GetFaction(get GetFunc, factionSymbol string) (*models.Faction, *models.APIError) {
	endpoint := fmt.Sprintf("/factions/%s", factionSymbol)

	var response struct {
		Data models.Faction `json:"data"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

type ListFactionsResponse struct {
	Data []*models.Faction `json:"data"`
	Meta models.Meta       `json:"meta"`
}

// ListAgents retrieves a list of agents with pagination
func ListFactions(get GetFunc, meta *models.Meta) ([]*models.Faction, *models.Meta, *models.APIError) {
	endpoint := "/factions"

	var response models.ListFactionsResponse

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
