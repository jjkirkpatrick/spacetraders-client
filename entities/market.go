package entities

import (
	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/internal/api"
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

// GetSupplyChain retrieves the supply chain information showing which exports map to which imports
func GetSupplyChain(c *client.Client) (*models.SupplyChainResponse, error) {
	response, err := api.GetSupplyChain(c.Get)
	if err != nil {
		return nil, err
	}

	return response, nil
}
