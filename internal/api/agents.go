package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/internal/models"
)

// GetFunc is a function type that sends a GET request to the specified endpoint
type GetFunc func(endpoint string, queryParams map[string]string, result interface{}) *models.APIError
type PostFunc func(endpoint string, payload interface{}, queryParams map[string]string, result interface{}) *models.APIError
type PutFunc func(endpoint string, payload interface{}, queryParams map[string]string, result interface{}) *models.APIError
type DeleteFunc func(endpoint string) *models.APIError
type PatchFunc func(endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError

// GetAgent retrieves the agent's details
func GetAgent(get GetFunc) (*models.Agent, *models.APIError) {
	endpoint := "/my/agent"

	var response struct {
		Data models.Agent `json:"data"`
	}

	err := get(endpoint, nil, &response)

	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// ListAgentsResponse represents the response from the ListAgents endpoint
type listAgentsResponse struct {
	Data []*models.Agent `json:"data"`
	Meta models.Meta     `json:"meta"`
}

// ListAgents retrieves a list of agents with pagination
func ListAgents(get GetFunc, meta *models.Meta) ([]*models.Agent, *models.Meta, *models.APIError) {
	endpoint := "/agents"

	var response listAgentsResponse

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

// GetPublicAgentResponse represents the response from the GetPublicAgent endpoint
type GetPublicAgentResponse struct {
	Data *models.Agent `json:"data"`
}

// GetPublicAgent retrieves the details of a public agent
func GetPublicAgent(get GetFunc, agentSymbol string) (*models.Agent, *models.APIError) {
	endpoint := fmt.Sprintf("/agents/%s", agentSymbol)

	var response GetPublicAgentResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
