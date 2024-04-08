package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

type listContractResponse struct {
	Data []*models.Contract `json:"data"`
	Meta models.Meta        `json:"meta"`
}

func ListContracts(get GetFunc, meta *models.Meta) ([]*models.Contract, *models.Meta, *models.APIError) {
	endpoint := "/my/contracts"

	var response listContractResponse

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

func GetContract(get GetFunc, contractId string) (*models.Contract, *models.APIError) {
	endpoint := fmt.Sprintf("/my/contracts/%s", contractId)

	var response struct {
		Data models.Contract `json:"data"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func AcceptContract(post PostFunc, contractId string) (*models.Agent, *models.Contract, *models.APIError) {
	endpoint := fmt.Sprintf("/my/contracts/%s/accept", contractId)

	var response struct {
		Data struct {
			Agent    *models.Agent    `json:"agent"`
			Contract *models.Contract `json:"contract"`
		}
	}

	err := post(endpoint, nil, nil, &response)

	if err != nil {
		return nil, nil, err
	}

	return response.Data.Agent, response.Data.Contract, nil

}

func DeliverContractCargo(post PostFunc, contractId string, body models.DeliverContractCargoRequest) (*models.Contract, *models.Cargo, *models.APIError) {
	endpoint := fmt.Sprintf("/my/contracts/%s/deliver", contractId)

	fmt.Println("Delivering Contract Cargo with Request Body:", body)

	var response struct {
		Data struct {
			Contract *models.Contract `json:"contract"`
			Cargo    *models.Cargo    `json:"cargo"`
		}
	}

	err := post(endpoint, body, nil, &response)

	if err != nil {
		return nil, nil, err
	}

	return response.Data.Contract, response.Data.Cargo, nil
}

func FulfillContract(post PostFunc, contractId string) (*models.Agent, *models.Contract, *models.APIError) {
	endpoint := fmt.Sprintf("/my/contracts/%s/fulfill", contractId)

	var response struct {
		Data struct {
			Agent    *models.Agent    `json:"agent"`
			Contract *models.Contract `json:"contract"`
		}
	}

	err := post(endpoint, nil, nil, &response)

	if err != nil {
		return nil, nil, err
	}

	return response.Data.Agent, response.Data.Contract, nil
}
