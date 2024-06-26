package entities

import (
	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/internal/api"
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

type Agent struct {
	models.Agent
	Client *client.Client
}

func ListPublicAgents(c *client.Client) ([]*Agent, error) {
	fetchFunc := func(meta models.Meta) ([]*Agent, models.Meta, error) {
		metaPtr := &meta
		agents, metaPtr, err := api.ListAgents(c.Get, metaPtr)

		var convertedAgents []*Agent
		for _, modelAgent := range agents {
			convertedShip := &Agent{
				Agent:  *modelAgent,
				Client: c,
			}
			convertedAgents = append(convertedAgents, convertedShip)
		}

		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
				metaPtr = &defaultMeta
			}
			return convertedAgents, *metaPtr, err.AsError()
		}
		if metaPtr != nil {
			return convertedAgents, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
			return convertedAgents, defaultMeta, nil
		}
	}
	return client.NewPaginator[*Agent](fetchFunc).FetchAllPages()
}

func GetAgent(c *client.Client) (*Agent, error) {
	agent, err := api.GetAgent(c.Get)
	if err != nil {
		return nil, err
	}

	agentEntity := &Agent{
		Agent:  *agent,
		Client: c,
	}

	return agentEntity, nil
}

func GetPublicAgent(c *client.Client, symbol string) (*Agent, error) {
	agent, err := api.GetPublicAgent(c.Get, symbol)
	if err != nil {
		return nil, err
	}

	agentEntity := &Agent{
		Agent:  *agent,
		Client: c,
	}

	return agentEntity, nil
}
