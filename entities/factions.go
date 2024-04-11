package entities

import (
	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/internal/api"
	"github.com/jjkirkpatrick/spacetraders-client/internal/models"
)

type Faction struct {
	models.Faction
	client *client.Client
}

func ListFactions(c *client.Client) ([]*Faction, error) {
	fetchFunc := func(meta models.Meta) ([]*Faction, models.Meta, error) {
		metaPtr := &meta
		factions, metaPtr, err := api.ListFactions(c.Get, metaPtr)

		var convertedFactions []*Faction
		for _, modelFaction := range factions {
			convertedFaction := &Faction{
				Faction: *modelFaction,
				client:  c,
			}
			convertedFactions = append(convertedFactions, convertedFaction)
		}

		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return convertedFactions, *metaPtr, err.AsError()
		}
		if metaPtr != nil {
			return convertedFactions, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return convertedFactions, defaultMeta, nil
		}
	}
	return client.NewPaginator[*Faction](fetchFunc).FetchAllPages()
}

func GetFaction(c *client.Client, symbol string) (*Faction, error) {
	faction, err := api.GetFaction(c.Get, symbol)
	if err != nil {
		return nil, err
	}

	agentEntity := &Faction{
		Faction: *faction,
		client:  c,
	}

	return agentEntity, nil
}
