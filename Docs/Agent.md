# Agent Operations Guide

This guide covers agent-related operations using the `entities` package.

## Getting Started

```go
import (
    "github.com/jjkirkpatrick/spacetraders-client/client"
    "github.com/jjkirkpatrick/spacetraders-client/entities"
)

// Create and configure client
options := client.DefaultClientOptions()
options.Symbol = "YOUR-AGENT"
options.Faction = "COSMIC"

c, err := client.NewClient(options)
if err != nil {
    // handle error
}
defer c.Close(ctx)
```

## Functions

### GetAgent

Retrieves detailed information about the authenticated agent.

```go
func GetAgent(c *client.Client) (*Agent, error)
```

**Example:**
```go
agent, err := entities.GetAgent(c)
if err != nil {
    log.Fatalf("Failed to get agent: %v", err)
}

fmt.Printf("Symbol: %s\n", agent.Symbol)
fmt.Printf("Credits: %d\n", agent.Credits)
fmt.Printf("Headquarters: %s\n", agent.Headquarters)
fmt.Printf("Ship Count: %d\n", agent.ShipCount)
```

### ListPublicAgents

Fetches a paginated list of all public agents in the game.

```go
func ListPublicAgents(c *client.Client) ([]*Agent, error)
```

**Example:**
```go
agents, err := entities.ListPublicAgents(c)
if err != nil {
    log.Fatalf("Failed to list agents: %v", err)
}

for _, agent := range agents {
    fmt.Printf("Agent: %s, Credits: %d\n", agent.Symbol, agent.Credits)
}
```

### GetPublicAgent

Retrieves detailed information about a specific public agent by symbol.

```go
func GetPublicAgent(c *client.Client, symbol string) (*Agent, error)
```

**Example:**
```go
agent, err := entities.GetPublicAgent(c, "SOME-AGENT")
if err != nil {
    log.Fatalf("Failed to get public agent: %v", err)
}

fmt.Printf("Agent %s has %d credits\n", agent.Symbol, agent.Credits)
```

## Agent Structure

The `Agent` entity contains the following fields from the API:

| Field | Type | Description |
|-------|------|-------------|
| `Symbol` | string | Unique agent identifier |
| `Headquarters` | string | Agent's headquarters waypoint |
| `Credits` | int64 | Current credit balance |
| `StartingFaction` | string | Faction the agent started with |
| `ShipCount` | int | Number of ships owned |
