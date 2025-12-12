# Contract Operations Guide

This guide covers contract-related operations using the `entities` package.

## Getting Started

```go
import (
    "github.com/jjkirkpatrick/spacetraders-client/client"
    "github.com/jjkirkpatrick/spacetraders-client/entities"
    "github.com/jjkirkpatrick/spacetraders-client/models"
)

// Create client
c, err := client.NewClient(options)
defer c.Close(ctx)

// Get contracts
contracts, err := entities.ListContracts(c)
```

## Functions

### ListContracts

Fetches all contracts available to the agent.

```go
func ListContracts(c *client.Client) ([]*Contract, error)
```

**Example:**
```go
contracts, err := entities.ListContracts(c)
if err != nil {
    log.Fatalf("Failed to list contracts: %v", err)
}

for _, contract := range contracts {
    fmt.Printf("Contract: %s, Type: %s, Accepted: %v\n",
        contract.ID, contract.Type, contract.Accepted)
}
```

### GetContract

Retrieves a specific contract by ID.

```go
func GetContract(c *client.Client, contractID string) (*Contract, error)
```

**Example:**
```go
contract, err := entities.GetContract(c, "contract-id-here")
if err != nil {
    log.Fatalf("Failed to get contract: %v", err)
}
fmt.Printf("Contract: %s, Fulfilled: %v\n", contract.ID, contract.Fulfilled)
```

### Accept

Accepts a contract, allowing the agent to start fulfilling its requirements.

```go
func (c *Contract) Accept() (*Agent, *Contract, error)
```

**Example:**
```go
agent, updatedContract, err := contract.Accept()
if err != nil {
    log.Fatalf("Failed to accept contract: %v", err)
}
fmt.Printf("Contract accepted! New credit balance: %d\n", agent.Credits)
```

### DeliverCargo

Delivers cargo from a ship to fulfill contract requirements.

```go
func (c *Contract) DeliverCargo(ship *Ship, tradeGood models.GoodSymbol, units int) (*Contract, *models.Cargo, error)
```

**Example:**
```go
// Get the ship
ship, err := entities.GetShip(c, "AGENT-1")
if err != nil {
    log.Fatalf("Failed to get ship: %v", err)
}

// Deliver cargo to the contract
updatedContract, cargo, err := contract.DeliverCargo(ship, models.GoodSymbol("IRON_ORE"), 100)
if err != nil {
    log.Fatalf("Failed to deliver cargo: %v", err)
}

fmt.Printf("Delivered! Units fulfilled: %d\n",
    updatedContract.Terms.Deliver[0].UnitsFulfilled)
```

### Fulfill

Marks a contract as fulfilled once all delivery requirements are met.

```go
func (c *Contract) Fulfill() (*models.Agent, *models.Contract, error)
```

**Example:**
```go
agent, fulfilledContract, err := contract.Fulfill()
if err != nil {
    log.Fatalf("Failed to fulfill contract: %v", err)
}
fmt.Printf("Contract fulfilled! Reward: %d credits\n",
    fulfilledContract.Terms.Payment.OnFulfilled)
```

## Contract Workflow

Here's a typical workflow for completing a contract:

```go
// 1. Get available contracts
contracts, err := entities.ListContracts(c)
if err != nil {
    log.Fatalf("Failed to list contracts: %v", err)
}

// 2. Find and accept a contract
for _, contract := range contracts {
    if !contract.Accepted {
        agent, contract, err := contract.Accept()
        if err != nil {
            log.Printf("Failed to accept contract %s: %v", contract.ID, err)
            continue
        }
        fmt.Printf("Accepted contract: %s\n", contract.ID)
        break
    }
}

// 3. Check delivery requirements
for _, delivery := range contract.Terms.Deliver {
    fmt.Printf("Need to deliver %d units of %s to %s\n",
        delivery.UnitsRequired,
        delivery.TradeSymbol,
        delivery.DestinationSymbol)
}

// 4. Mine/buy the required goods and deliver them
// ... (mining/trading code)

// 5. Navigate to delivery destination and deliver
ship, _ := entities.GetShip(c, "AGENT-1")
ship.Navigate(contract.Terms.Deliver[0].DestinationSymbol)
// Wait for arrival...
ship.Dock()

// Deliver the cargo
contract.DeliverCargo(ship,
    models.GoodSymbol(contract.Terms.Deliver[0].TradeSymbol),
    contract.Terms.Deliver[0].UnitsRequired)

// 6. Fulfill the contract once all deliveries are complete
agent, contract, err := contract.Fulfill()
if err != nil {
    log.Fatalf("Failed to fulfill: %v", err)
}
fmt.Printf("Contract complete! Earned %d credits\n",
    contract.Terms.Payment.OnFulfilled)
```

## Contract Structure

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Unique contract identifier |
| `FactionSymbol` | string | Faction offering the contract |
| `Type` | string | Contract type (PROCUREMENT, TRANSPORT, SHUTTLE) |
| `Terms` | Terms | Contract terms including deliveries and payment |
| `Accepted` | bool | Whether the contract has been accepted |
| `Fulfilled` | bool | Whether the contract has been fulfilled |
| `Expiration` | string | When the contract expires (RFC3339) |
| `DeadlineToAccept` | string | Deadline to accept the contract |

## Contract Terms

The `Terms` field contains:

| Field | Type | Description |
|-------|------|-------------|
| `Deadline` | string | Deadline to fulfill the contract |
| `Payment.OnAccepted` | int | Credits received on acceptance |
| `Payment.OnFulfilled` | int | Credits received on fulfillment |
| `Deliver` | []Deliver | List of delivery requirements |

## Delivery Requirements

Each delivery requirement contains:

| Field | Type | Description |
|-------|------|-------------|
| `TradeSymbol` | string | The good to deliver |
| `DestinationSymbol` | string | Waypoint to deliver to |
| `UnitsRequired` | int | Total units needed |
| `UnitsFulfilled` | int | Units already delivered |
