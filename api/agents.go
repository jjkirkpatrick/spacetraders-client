package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

// GetFunc is a function type that sends a GET request to the specified endpoint
type GetFunc func(endpoint string, result interface{}) error
type PostFunc func(endpoint string, payload interface{}, result interface{}) error
type PutFunc func(endpoint string, payload interface{}, result interface{}) error
type DeleteFunc func(endpoint string) error
type PatchFunc func(endpoint string, payload interface{}, result interface{}) error

// GetAgent retrieves the agent's details
func GetAgent(get GetFunc) (*models.Agent, error) {
	endpoint := "/my/agent"

	var response struct {
		Data models.Agent `json:"data"`
	}

	err := get(endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent details: %v", err)
	}

	return &response.Data, nil
}

// ListAgentsResponse represents the response from the ListAgents endpoint
type listAgentsResponse struct {
	Data []*models.Agent `json:"data"`
	Meta models.Meta     `json:"meta"`
}

// ListAgents retrieves a list of agents with pagination
func ListAgents(get GetFunc, limit, page int) ([]*models.Agent, error) {
	endpoint := fmt.Sprintf("/agents?limit=%d&page=%d", limit, page)

	var response listAgentsResponse

	err := get(endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list agents: %v", err)
	}

	return response.Data, nil
}

// GetPublicAgentResponse represents the response from the GetPublicAgent endpoint
type GetPublicAgentResponse struct {
	Data *models.Agent `json:"data"`
}

// GetPublicAgent retrieves the details of a public agent
func GetPublicAgent(get GetFunc, agentSymbol string) (*models.Agent, error) {
	endpoint := fmt.Sprintf("/agents/%s", agentSymbol)

	var response GetPublicAgentResponse

	err := get(endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get public agent: %v", err)
	}

	return response.Data, nil
}
