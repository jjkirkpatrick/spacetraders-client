# System Operations Guide

This guide covers system and waypoint operations using the `entities` package.

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

// Get a system
system, err := entities.GetSystem(c, "X1-ABC")
```

## System Functions

### ListSystems

Fetches all systems in the game (paginated automatically).

```go
func ListSystems(c *client.Client) ([]*System, error)
```

**Example:**
```go
systems, err := entities.ListSystems(c)
if err != nil {
    log.Fatalf("Failed to list systems: %v", err)
}
fmt.Printf("Found %d systems\n", len(systems))
```

### GetSystem

Retrieves detailed information about a specific system.

```go
func GetSystem(c *client.Client, symbol string) (*System, error)
```

**Example:**
```go
system, err := entities.GetSystem(c, "X1-ABC")
if err != nil {
    log.Fatalf("Failed to get system: %v", err)
}
fmt.Printf("System: %s, Type: %s\n", system.Symbol, system.Type)
```

## Waypoint Functions

### ListWaypoints

Lists all waypoints in a system, optionally filtered by trait and type.

```go
func (s *System) ListWaypoints(trait models.WaypointTrait, waypointType models.WaypointType) ([]*models.Waypoint, *models.Meta, error)
```

**Example:**
```go
// List all waypoints
waypoints, meta, err := system.ListWaypoints("", "")

// Filter by type
waypoints, meta, err := system.ListWaypoints("", "ASTEROID")

// Filter by trait
waypoints, meta, err := system.ListWaypoints("MARKETPLACE", "")
```

### GetWaypointsWithTrait

Helper function to get waypoints with specific trait and/or type.

```go
func (s *System) GetWaypointsWithTrait(trait string, waypointType string) ([]*models.Waypoint, error)
```

**Example:**
```go
// Find all waypoints with shipyards
shipyards, err := system.GetWaypointsWithTrait("SHIPYARD", "")

// Find engineered asteroids
asteroids, err := system.GetWaypointsWithTrait("", "ENGINEERED_ASTEROID")

// Find marketplaces on orbital stations
markets, err := system.GetWaypointsWithTrait("MARKETPLACE", "ORBITAL_STATION")
```

### FetchWaypoint

Fetches detailed information about a specific waypoint.

```go
func (s *System) FetchWaypoint(symbol string) (*models.Waypoint, error)
```

**Example:**
```go
waypoint, err := system.FetchWaypoint("X1-ABC-ASTEROID1")
if err != nil {
    log.Fatalf("Failed to fetch waypoint: %v", err)
}
fmt.Printf("Waypoint: %s, Type: %s\n", waypoint.Symbol, waypoint.Type)
```

## Market Functions

### GetMarket

Retrieves market information for a waypoint with a marketplace trait.

```go
func (s *System) GetMarket(waypointSymbol string) (*models.Market, error)
```

**Example:**
```go
market, err := system.GetMarket("X1-ABC-STATION")
if err != nil {
    log.Fatalf("Failed to get market: %v", err)
}

fmt.Println("Trade Goods:")
for _, good := range market.TradeGoods {
    fmt.Printf("  %s: Buy %d, Sell %d\n", good.Symbol, good.PurchasePrice, good.SellPrice)
}
```

## Shipyard Functions

### GetShipyard

Retrieves shipyard information for a waypoint with a shipyard trait.

```go
func (s *System) GetShipyard(waypointSymbol string) (*models.Shipyard, error)
```

**Example:**
```go
shipyard, err := system.GetShipyard("X1-ABC-SHIPYARD")
if err != nil {
    log.Fatalf("Failed to get shipyard: %v", err)
}

fmt.Println("Available Ships:")
for _, ship := range shipyard.Ships {
    fmt.Printf("  %s: %d credits\n", ship.Type, ship.PurchasePrice)
}
```

## Jump Gate Functions

### GetJumpGate

Retrieves jump gate information including connected systems.

```go
func (s *System) GetJumpGate(waypointSymbol string) (*models.JumpGate, error)
```

**Example:**
```go
jumpGate, err := system.GetJumpGate("X1-ABC-JUMPGATE")
if err != nil {
    log.Fatalf("Failed to get jump gate: %v", err)
}

fmt.Println("Connected Systems:")
for _, connection := range jumpGate.Connections {
    fmt.Printf("  %s (distance: %d)\n", connection.Symbol, connection.Distance)
}
```

## Construction Functions

### GetConstructionSite

Retrieves construction site information for a waypoint under construction.

```go
func (s *System) GetConstructionSite(waypointSymbol string) (*models.ConstructionSite, error)
```

### SupplyConstructionSite

Supplies materials to a construction site from a ship.

```go
func (s *System) SupplyConstructionSite(shipSymbol string, waypointSymbol string, good models.GoodSymbol, quantity int) error
```

**Example:**
```go
err := system.SupplyConstructionSite("AGENT-1", "X1-ABC-CONSTRUCTION", models.GoodSymbol("IRON"), 100)
if err != nil {
    log.Fatalf("Failed to supply construction site: %v", err)
}
```

## Utility Functions

### CalculateDistanceBetweenWaypoints

Calculates the Euclidean distance between two waypoints.

```go
func CalculateDistanceBetweenWaypoints(x1, y1, x2, y2 int) float64
```

**Example:**
```go
distance := entities.CalculateDistanceBetweenWaypoints(
    waypoint1.X, waypoint1.Y,
    waypoint2.X, waypoint2.Y,
)
fmt.Printf("Distance: %.2f units\n", distance)
```

## Common Waypoint Types

| Type | Description |
|------|-------------|
| `PLANET` | A planet |
| `GAS_GIANT` | A gas giant (for siphoning) |
| `MOON` | A moon |
| `ORBITAL_STATION` | An orbital station |
| `JUMP_GATE` | A jump gate for inter-system travel |
| `ASTEROID_FIELD` | An asteroid field for mining |
| `ASTEROID` | An individual asteroid |
| `ENGINEERED_ASTEROID` | An engineered asteroid with guaranteed resources |
| `ASTEROID_BASE` | A base on an asteroid |
| `NEBULA` | A nebula |
| `DEBRIS_FIELD` | A debris field |
| `GRAVITY_WELL` | A gravity well |
| `ARTIFICIAL_GRAVITY_WELL` | An artificial gravity well |
| `FUEL_STATION` | A fuel station |

## Common Waypoint Traits

| Trait | Description |
|-------|-------------|
| `MARKETPLACE` | Has a marketplace for trading |
| `SHIPYARD` | Has a shipyard for purchasing ships |
| `UNCHARTED` | Not yet charted |
| `UNDER_CONSTRUCTION` | Currently under construction |
| `OUTPOST` | An outpost |
| `STRIPPED` | Resources have been stripped |
| `OVERCROWDED` | Overcrowded location |
| `HIGH_TECH` | High-tech goods available |
| `CORRUPT` | Corrupt officials |
| `BUREAUCRATIC` | Bureaucratic processes |
| `TRADING_HUB` | Major trading hub |
| `INDUSTRIAL` | Industrial production |
| `BLACK_MARKET` | Black market available |
| `RESEARCH_FACILITY` | Research facility |
| `MILITARY_BASE` | Military base |
| `SURVEILLANCE_OUTPOST` | Surveillance outpost |
| `EXPLORATION_OUTPOST` | Exploration outpost |
| `MINERAL_DEPOSITS` | Contains mineral deposits |
| `COMMON_METAL_DEPOSITS` | Common metal deposits |
| `PRECIOUS_METAL_DEPOSITS` | Precious metal deposits |
| `RARE_METAL_DEPOSITS` | Rare metal deposits |
| `METHANE_POOLS` | Methane pools |
| `ICE_CRYSTALS` | Ice crystal deposits |
| `EXPLOSIVE_GASES` | Explosive gas deposits |
