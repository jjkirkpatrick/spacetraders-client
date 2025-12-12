# Ship Operations Guide

This guide covers ship-related operations using the `entities` package.

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

// Get your ships
ships, err := entities.ListShips(c)
ship := ships[0]  // Work with first ship
```

## Fleet Management

### ListShips

Fetches all ships owned by the agent.

```go
func ListShips(c *client.Client) ([]*Ship, error)
```

### GetShip

Retrieves a specific ship by symbol.

```go
func GetShip(c *client.Client, symbol string) (*Ship, error)
```

### PurchaseShip

Purchases a new ship from a shipyard.

```go
func PurchaseShip(c *client.Client, shipType string, waypoint string) (*models.Agent, *Ship, *models.Transaction, error)
```

**Example:**
```go
agent, ship, transaction, err := entities.PurchaseShip(c, "SHIP_MINING_DRONE", "X1-ABC-SHIPYARD")
if err != nil {
    log.Fatalf("Failed to purchase ship: %v", err)
}
fmt.Printf("Purchased %s for %d credits\n", ship.Symbol, transaction.TotalPrice)
```

## Navigation

### Orbit

Commands a ship to enter orbit around the current waypoint.

```go
func (s *Ship) Orbit() (*models.ShipNav, error)
```

### Dock

Commands a ship to dock at the current waypoint.

```go
func (s *Ship) Dock() (*models.ShipNav, error)
```

### Navigate

Navigates a ship to a specified waypoint within the same system.

```go
func (s *Ship) Navigate(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, []models.Event, error)
```

**Example:**
```go
fuel, nav, events, err := ship.Navigate("X1-ABC-ASTEROID")
if err != nil {
    log.Fatalf("Navigation failed: %v", err)
}
fmt.Printf("Arriving at %s, ETA: %s\n", nav.WaypointSymbol, nav.Route.Arrival)
```

### Warp

Commands a ship to warp to a waypoint in another system.

```go
func (s *Ship) Warp(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, error)
```

### Jump

Commands a ship to jump to another system using a jump gate.

```go
func (s *Ship) Jump(systemSymbol string) (*models.ShipNav, *models.ShipCooldown, *models.Transaction, *models.Agent, error)
```

### SetFlightMode

Sets the flight mode of a ship (CRUISE, BURN, DRIFT, STEALTH).

```go
func (s *Ship) SetFlightMode(flightmode models.FlightMode) error
```

**Example:**
```go
err := ship.SetFlightMode(models.FlightModeBurn)
```

### FetchNavigationStatus

Retrieves the current navigation status of a ship.

```go
func (s *Ship) FetchNavigationStatus() (*models.ShipNav, error)
```

### GetRouteToDestination

Calculates an optimal route to a destination waypoint using pathfinding.

```go
func (s *Ship) GetRouteToDestination(destination string) (*models.PathfindingRoute, error)
```

**Example:**
```go
route, err := ship.GetRouteToDestination("X1-ABC-DESTINATION")
if err != nil {
    log.Fatalf("Pathfinding failed: %v", err)
}
for _, step := range route.Steps {
    fmt.Printf("Step: %s via %s\n", step.Waypoint, step.FlightMode)
}
```

### CalculateFuelRequired

Calculates fuel required for a given distance and flight mode.

```go
func (s *Ship) CalculateFuelRequired(distance float64, flightMode models.FlightMode) int
```

### CalculateTravelTime

Calculates travel time in seconds for a given distance and flight mode.

```go
func (s *Ship) CalculateTravelTime(distance float64, flightMode models.FlightMode) int
```

## Mining & Extraction

### Extract

Extracts resources from the current location (asteroid field).

```go
func (s *Ship) Extract() (*models.Extraction, error)
```

**Example:**
```go
extraction, err := ship.Extract()
if err != nil {
    log.Fatalf("Extraction failed: %v", err)
}
fmt.Printf("Extracted %d units of %s\n", extraction.Yield.Units, extraction.Yield.Symbol)
```

### ExtractWithSurvey

Extracts resources using a previously conducted survey for better yields.

```go
func (s *Ship) ExtractWithSurvey(survey models.Survey) (*models.Extraction, error)
```

### Survey

Conducts a survey of the surrounding area to find resource deposits.

```go
func (s *Ship) Survey() ([]models.Survey, error)
```

### Siphon

Siphons resources from gas giants.

```go
func (s *Ship) Siphon() (*models.Extraction, error)
```

### Refine

Processes raw materials into refined goods.

```go
func (s *Ship) Refine(produce string) (*models.Produced, *models.Consumed, error)
```

## Cargo Management

### FetchCargo

Retrieves the current cargo contents of a ship.

```go
func (s *Ship) FetchCargo() (*models.Cargo, error)
```

**Example:**
```go
cargo, err := ship.FetchCargo()
if err != nil {
    log.Fatalf("Failed to fetch cargo: %v", err)
}
fmt.Printf("Cargo: %d/%d units\n", cargo.Units, cargo.Capacity)
for _, item := range cargo.Inventory {
    fmt.Printf("  %s: %d units\n", item.Symbol, item.Units)
}
```

### Jettison

Jettisons cargo into space.

```go
func (s *Ship) Jettison(goodSymbol models.GoodSymbol, units int) (*models.Cargo, error)
```

### TransferCargo

Transfers cargo to another ship.

```go
func (s *Ship) TransferCargo(goodSymbol models.GoodSymbol, units int, shipSymbol string) (*models.Cargo, error)
```

## Trading

### SellCargo

Sells cargo at a marketplace.

```go
func (s *Ship) SellCargo(goodSymbol models.GoodSymbol, units int) (*models.Agent, *models.Cargo, *models.Transaction, error)
```

**Example:**
```go
agent, cargo, transaction, err := ship.SellCargo(models.GoodSymbol("IRON_ORE"), 10)
if err != nil {
    log.Fatalf("Sale failed: %v", err)
}
fmt.Printf("Sold for %d credits\n", transaction.TotalPrice)
```

### PurchaseCargo

Purchases cargo from a marketplace.

```go
func (s *Ship) PurchaseCargo(goodSymbol models.GoodSymbol, units int) (*models.Agent, *models.Cargo, *models.Transaction, error)
```

## Fuel

### Refuel

Refuels the ship at a marketplace or from cargo.

```go
func (s *Ship) Refuel(amount int, fromCargo bool) (*models.Agent, *models.FuelDetails, *models.Transaction, error)
```

**Example:**
```go
// Refuel to full from marketplace
agent, fuel, transaction, err := ship.Refuel(0, false)

// Refuel 100 units from cargo
agent, fuel, transaction, err := ship.Refuel(100, true)
```

## Scanning

### ScanSystems

Scans for nearby systems.

```go
func (s *Ship) ScanSystems() (*models.ShipCooldown, []models.System, error)
```

### ScanWaypoints

Scans for waypoints in the current system.

```go
func (s *Ship) ScanWaypoints() (*models.ShipCooldown, []models.Waypoint, error)
```

## Cooldowns

### FetchCooldown

Checks the cooldown status of a ship's operations.

```go
func (s *Ship) FetchCooldown() (*models.ShipCooldown, error)
```

**Example:**
```go
cooldown, err := ship.FetchCooldown()
if err != nil {
    log.Fatalf("Failed to fetch cooldown: %v", err)
}
if cooldown.RemainingSeconds > 0 {
    fmt.Printf("Cooldown: %d seconds remaining\n", cooldown.RemainingSeconds)
    time.Sleep(time.Duration(cooldown.RemainingSeconds) * time.Second)
}
```

## Charting

### Chart

Charts the current waypoint, making it visible to other agents.

```go
func (s *Ship) Chart() (*models.Chart, *models.Waypoint, error)
```

## Contracts

### NegotiateContract

Negotiates a new contract at the current location.

```go
func (s *Ship) NegotiateContract() (*models.Contract, error)
```

## Ship Modifications

### GetMounts

Retrieves the ship's current mounts.

```go
func (s *Ship) GetMounts() ([]models.ShipMount, error)
```

### InstallMount

Installs a mount on the ship.

```go
func (s *Ship) InstallMount(mountSymbol models.MountSymbol) (*models.Agent, []models.ShipMount, *models.Cargo, *models.Transaction, error)
```

### RemoveMount

Removes a mount from the ship.

```go
func (s *Ship) RemoveMount(mountSymbol models.MountSymbol) (*models.Agent, []models.ShipMount, *models.Cargo, *models.Transaction, error)
```

## Ship Maintenance

### GetRepairPrice

Gets the repair cost for the ship.

```go
func (s *Ship) GetRepairPrice() (*models.Transaction, error)
```

### RepairShip

Repairs the ship.

```go
func (s *Ship) RepairShip() (*models.Ship, *models.Transaction, error)
```

### GetScrapPrice

Gets the scrap value of the ship.

```go
func (s *Ship) GetScrapPrice() (*models.Transaction, error)
```

### ScrapShip

Scraps the ship for credits.

```go
func (s *Ship) ScrapShip() (*models.Transaction, error)
```
