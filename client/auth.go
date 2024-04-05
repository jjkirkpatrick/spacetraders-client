package client

import (
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

// RegisterRequest represents the request payload for registering a new agent
type RegisterRequest struct {
	Faction string `json:"faction"`
	Symbol  string `json:"symbol"`
	Email   string `json:"email,omitempty"`
}

// RegisterResponse represents the response payload for registering a new agent
type RegisterResponse struct {
	Data struct {
		Agent    models.Agent    `json:"agent"`
		Contract models.Contract `json:"contract"`
		Faction  models.Faction  `json:"faction"`
		Ship     models.Ship     `json:"ship"`
		Token    string          `json:"token"`
	} `json:"data"`
}

// RegisterNewAgent registers a new agent with the specified faction and symbol
func (c *Client) RegisterNewAgent(faction, symbol, email string) error {
	registerReq := RegisterRequest{
		Faction: faction,
		Symbol:  symbol,
		Email:   email,
	}

	var registerResp RegisterResponse
	err := c.Post("/register", registerReq, nil, &registerResp)
	if err != nil {
		return err
	}

	c.token = registerResp.Data.Token
	return nil
}
