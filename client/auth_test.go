package client

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jjkirkpatrick/spacetraders-client/internal/cache"
)

func TestGetOrRegisterToken(t *testing.T) {
	// Setup
	options := DefaultClientOptions()
	options.BaseURL = "https://stoplight.io/mocks/spacetraders/spacetraders/96627693"
	options.Symbol = "TestAgent"
	options.Faction = "COSMIC"

	client := &Client{
		baseURL:     options.BaseURL,
		httpClient:  resty.New(),
		context:     context.Background(),
		retryDelay:  options.RetryDelay,
		CacheClient: cache.NewCache(),
		Logger:      slog.Default(),
		RateLimiter: NewRateLimiter(2.0, 10.0),
	}

	// Test for valid faction and symbol
	err := client.getOrRegisterToken(options.Faction, options.Symbol, "email@example.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test for invalid faction
	err = client.getOrRegisterToken(options.Faction, "", "email@example.com")
	if err == nil {
		t.Errorf("Expected an error for invalid faction, got nil")
	}

	// Test for empty faction or symbol
	err = client.getOrRegisterToken("", options.Symbol, "email@example.com")
	if err == nil {
		t.Errorf("Expected an error for empty faction or symbol, got nil")
	}
}

func TestTokenFileOperations(t *testing.T) {
	// Setup
	options := DefaultClientOptions()
	options.BaseURL = "https://stoplight.io/mocks/spacetraders/spacetraders/96627693"
	options.Symbol = "TestAgent"
	options.Faction = "COSMIC"

	client := &Client{
		baseURL:     options.BaseURL,
		httpClient:  resty.New(),
		context:     context.Background(),
		retryDelay:  options.RetryDelay,
		CacheClient: cache.NewCache(),
		Logger:      slog.Default(),
		RateLimiter: NewRateLimiter(2.0, 10.0),
	}

	// Ensure token file is clean before tests
	_ = os.Remove("tokens.json")

	// Test 1: File should be created on a new agent that's not been used.
	t.Run("Create Token File on New Agent", func(t *testing.T) {
		err := client.getOrRegisterToken("COSMIC", "NewAgent", "newagent@example.com")
		if err != nil {
			t.Fatalf("Failed to register new agent: %v", err)
		}

		// Check if file exists
		if _, err := os.Stat("tokens.json"); os.IsNotExist(err) {
			t.Fatalf("Token file was not created")
		}
	})

	// Test 2: File should be updated with the agent when registered
	t.Run("Update Token File on Agent Registration", func(t *testing.T) {
		initialAgent := "NewAgent"
		newAgent := "SecondAgent"

		err := client.getOrRegisterToken("VOID", newAgent, "secondagent@example.com")
		if err != nil {
			t.Fatalf("Failed to register second agent: %v", err)
		}

		file, err := os.Open("tokens.json")
		if err != nil {
			t.Fatalf("Failed to open token file: %v", err)
		}
		defer file.Close()

		var tokenFile TokenFile
		err = json.NewDecoder(file).Decode(&tokenFile)
		if err != nil {
			t.Fatalf("Failed to decode token file: %v", err)
		}

		if _, ok := tokenFile.Tokens[initialAgent]; !ok {
			t.Errorf("Initial agent token not found in token file")
		}

		if _, ok := tokenFile.Tokens[newAgent]; !ok {
			t.Errorf("New agent token not found in token file")
		}
	})

	// Test 4: The token file should allow for multiple agents i.e., an array of agents
	t.Run("Token File Supports Multiple Agents", func(t *testing.T) {
		agents := []struct {
			faction string
			symbol  string
			email   string
		}{
			{"COSMIC", "AgentOne", "agentone@example.com"},
			{"VOID", "AgentTwo", "agenttwo@example.com"},
			{"GALACTIC", "AgentThree", "agentthree@example.com"},
		}

		for _, agent := range agents {
			err := client.getOrRegisterToken(agent.faction, agent.symbol, agent.email)
			if err != nil {
				t.Fatalf("Failed to register agent %s: %v", agent.symbol, err)
			}
		}

		file, err := os.Open("tokens.json")
		if err != nil {
			t.Fatalf("Failed to open token file: %v", err)
		}
		defer file.Close()

		var tokenFile TokenFile
		err = json.NewDecoder(file).Decode(&tokenFile)
		if err != nil {
			t.Fatalf("Failed to decode token file: %v", err)
		}

		for _, agent := range agents {
			if _, ok := tokenFile.Tokens[agent.symbol]; !ok {
				t.Errorf("Token for agent %s not found in token file", agent.symbol)
			}
		}

		if len(tokenFile.Tokens) < 3 {
			t.Errorf("Token file does not support multiple agents. Expected at least 3 tokens, got %d", len(tokenFile.Tokens))
		}
	})

	// Cleanup
	_ = os.Remove("tokens.json")
}
