package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

type listContractResponse struct {
	Data []*models.Contract `json:"data"`
	Meta models.Meta        `json:"meta"`
}

func ListContracts(get GetFunc, limit, page int) ([]*models.Contract, error) {
	endpoint := fmt.Sprintf("/my/contracts?limit=%d&page=%d", limit, page)

	var response listContractResponse

	err := get(endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %v", err)
	}

	return response.Data, nil
}

func GetContract(get GetFunc, contractId string) (*models.Contract, error) {
	endpoint := fmt.Sprintf("/my/contracts/%s", contractId)

	var response struct {
		Data models.Contract `json:"data"`
	}

	err := get(endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent details: %v", err)
	}

	return &response.Data, nil
}

func AcceptContract(post PostFunc, contractId string) (*models.Agent, *models.Contract, error) {
	endpoint := fmt.Sprintf("/my/contracts/%s/accept", contractId)

	var response struct {
		Data struct {
			Agent    *models.Agent    `json:"agent"`
			Contract *models.Contract `json:"contract"`
		}
	}

	err := post(endpoint, nil, &response)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent details: %v", err)
	}

	return response.Data.Agent, response.Data.Contract, nil

}

func DeliverContractCargo(post PostFunc, contractId string, body models.DeliverContractCargoRequest) (*models.Contract, *models.Cargo, error) {
	endpoint := fmt.Sprintf("/my/contracts/%s/deliver", contractId)

	var response struct {
		Data struct {
			Contract *models.Contract `json:"contract"`
			Cargo    *models.Cargo    `json:"cargo"`
		}
	}

	err := post(endpoint, body, &response)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent details: %v", err)
	}

	return response.Data.Contract, response.Data.Cargo, nil
}

func FulfillContract(post PostFunc, contractId string) (*models.Agent, *models.Contract, error) {
	endpoint := fmt.Sprintf("/my/contracts/%s/fulfill", contractId)

	var response struct {
		Data struct {
			Agent    *models.Agent    `json:"agent"`
			Contract *models.Contract `json:"contract"`
		}
	}

	err := post(endpoint, nil, &response)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get agent details: %v", err)
	}

	return response.Data.Agent, response.Data.Contract, nil
}
