package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

func GetFaction(get GetFunc, factionSymbol string) (*models.Faction, error) {
	endpoint := fmt.Sprintf("/factions/%s", factionSymbol)

	var response struct {
		Data models.Faction `json:"data"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %v", err)
	}

	return &response.Data, nil
}

type listFactionsResponse struct {
	Data []*models.Faction `json:"data"`
	Meta models.Meta       `json:"meta"`
}

func ListFactions(get GetFunc, limit, page int) ([]*models.Faction, error) {
	endpoint := fmt.Sprintf("/factions/f?limit=%d&page=%d", limit, page)

	var response listFactionsResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %v", err)
	}

	return response.Data, nil
}
