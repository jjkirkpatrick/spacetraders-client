package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/models"
	"github.com/phuslu/log"
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
	c.Logger.Debug("Attempting to get or register token", "faction", faction, "symbol", symbol, "email", email)

	if faction == "" || symbol == "" {
		return fmt.Errorf("faction and symbol must be set")
	}

	validFactions := map[string]bool{
		"COSMIC": true, "VOID": true, "GALACTIC": true, "QUANTUM": true,
		"DOMINION": true, "ASTRO": true, "CORSAIRS": true, "OBSIDIAN": true,
		"AEGIS": true, "UNITED": true, "SOLITARY": true, "COBALT": true,
		"OMEGA": true, "ECHO": true, "LORDS": true, "CULT": true,
		"ANCIENTS": true, "SHADOW": true, "ETHEREAL": true,
	}

	if _, ok := validFactions[faction]; !ok {
		return fmt.Errorf("invalid faction: %s", faction)
	}

	// Check if a token exists for the given symbol
	token, err := c.getTokenFromFile(symbol)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get token from file")
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
	log.Debug().Msgf("Updating token file with new token for symbol %s", symbol)

	// Read the current contents of the file
	fileContent, err := ioutil.ReadFile("tokens.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Token file doesn't exist, create a new one with the current token
			tokenFile := TokenFile{
				Tokens: map[string]string{
					symbol: token,
				},
			}
			return c.writeTokenFile(tokenFile)
		}
		return err
	}

	var tokenFile TokenFile
	err = json.Unmarshal(fileContent, &tokenFile)
	if err != nil {
		return err
	}

	// Update the token map with the new token
	if tokenFile.Tokens == nil {
		tokenFile.Tokens = make(map[string]string)
	}
	tokenFile.Tokens[symbol] = token

	// Write the updated token file
	return c.writeTokenFile(tokenFile)
}

func (c *Client) writeTokenFile(tokenFile TokenFile) error {
	file, err := os.OpenFile("tokens.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(tokenFile)
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
