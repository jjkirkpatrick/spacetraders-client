package client

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/internal/models"
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

// TokenFile represents the structure of the token file
type TokenFile struct {
	Tokens map[string]string `json:"tokens"`
}

// GetOrRegisterToken retrieves the token for the given symbol from the token file or registers a new agent if the token doesn't exist
func (c *Client) getOrRegisterToken(faction, symbol, email string) error {
	// Check if a token exists for the given symbol
	token, err := c.getTokenFromFile(symbol)
	if err != nil {
		fmt.Println("Error getting token from file:", err)
		return err
	}

	if token != "" {
		// Token found, set it in the client
		c.token = token
		return nil
	}

	// Token not found, register a new agent
	registerReq := RegisterRequest{
		Faction: faction,
		Symbol:  symbol,
		Email:   email,
	}

	var registerResp RegisterResponse
	apiErr := c.Post("/register", registerReq, nil, &registerResp)
	if apiErr != nil {
		return apiErr
	}

	fmt.Println("Agent registered successfully:", registerResp.Data)

	// Update the token file with the new token
	err = c.updateTokenFile(symbol, registerResp.Data.Token)
	if err != nil {
		return err
	}

	c.token = registerResp.Data.Token
	return nil
}

// getTokenFromFile retrieves the token for the given symbol from the token file
func (c *Client) getTokenFromFile(symbol string) (string, error) {
	file, err := os.Open("tokens.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Token file doesn't exist, create an empty one
			err = c.createEmptyTokenFile()
			if err != nil {
				return "", err
			}
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	var tokenFile TokenFile
	err = json.NewDecoder(file).Decode(&tokenFile)
	if err != nil {
		return "", err
	}

	token, exists := tokenFile.Tokens[symbol]
	if !exists {
		return "", nil // Token does not exist for the given symbol
	}

	return token, nil
}

// updateTokenFile updates the token file with the new token for the given symbol
func (c *Client) updateTokenFile(symbol, token string) error {
	fmt.Println("Updating token file with new token for symbol", symbol)
	file, err := os.OpenFile("tokens.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var tokenFile TokenFile
	err = json.NewDecoder(file).Decode(&tokenFile)
	if err != nil && err != io.EOF {
		return err
	}
	if tokenFile.Tokens == nil {
		tokenFile.Tokens = make(map[string]string)
	}

	tokenFile.Tokens[symbol] = token

	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Truncate the file to clear its contents
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	// Write the updated token file
	err = json.NewEncoder(file).Encode(tokenFile)
	if err != nil {
		return err
	}

	return nil
}

// createEmptyTokenFile creates an empty token file
func (c *Client) createEmptyTokenFile() error {
	file, err := os.Create("tokens.json")
	if err != nil {
		return err
	}
	defer file.Close()

	tokenFile := TokenFile{
		Tokens: make(map[string]string),
	}

	err = json.NewEncoder(file).Encode(tokenFile)
	if err != nil {
		return err
	}

	return nil
}
