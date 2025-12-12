package api

import (
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

// GetSupplyChain retrieves the supply chain information showing which exports map to which imports
func GetSupplyChain(get GetFunc) (*models.SupplyChainResponse, *models.APIError) {
	endpoint := "/market/supply-chain"

	var response models.SupplyChainResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
