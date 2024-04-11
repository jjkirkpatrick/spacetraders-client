# Ship Operations Guide

This guide provides an overview of how to interact with ship-related functionalities using the provided methods in `ships.go`. Each function is designed to perform specific operations related to ships in the SpaceTraders API.

## Getting Started

Before you can use any of the ship functions, ensure you have a client instance created and authenticated. Refer to the client setup guide for instructions on initializing and authenticating your client.

## Function Descriptions

### ListShips

- **Purpose**: Fetches a list of all ships the agent owns.
- **Usage**: `ListShips(client *client.Client) ([]*Ship, error)`

### GetShip

- **Purpose**: Retrieves detailed information about a specific ship by its symbol.
- **Usage**: `GetShip(client *client.Client, symbol string) (*Ship, error)`

### PurchaseShip

- **Purpose**: Allows the purchase of a new ship.
- **Usage**: `PurchaseShip(client *client.Client, shipType string, waypoint string) (*models.Agent, *Ship, *models.Transaction, error)`

### Orbit

- **Purpose**: Commands a ship to enter orbit around the current waypoint.
- **Usage**: `func (s *Ship) Orbit() (*models.ShipNav, error)`

### Dock

- **Purpose**: Commands a ship to dock at the current waypoint.
- **Usage**: `func (s *Ship) Dock() (*models.ShipNav, error)`

### FetchCargo

- **Purpose**: Retrieves the current cargo load of a ship.
- **Usage**: `func (s *Ship) FetchCargo() (*models.Cargo, error)`

### Refine

- **Purpose**: Processes raw materials into a refined state.
- **Usage**: `func (s *Ship) Refine(produce string) (*models.Produced, *models.Consumed, error)`

### Chart

- **Purpose**: Generates a navigation chart for a ship.
- **Usage**: `func (s *Ship) Chart() (*models.Chart, *models.Waypoint, error)`

### FetchCooldown

- **Purpose**: Checks the cooldown status of a ship's operations.
- **Usage**: `func (s *Ship) FetchCooldown() error`

### Survey

- **Purpose**: Conducts a survey of the surrounding area for resources.
- **Usage**: `func (s *Ship) Survey() ([]models.Survey, error)`

### Extract

- **Purpose**: Extracts resources from the current location.
- **Usage**: `func (s *Ship) Extract() (*models.Extraction, error)`

### Siphon

- **Purpose**: Siphons resources from another entity.
- **Usage**: `func (s *Ship) Siphon() (*models.Extraction, error)`

### ExtractWithSurvey

- **Purpose**: Extracts resources based on a previously conducted survey.
- **Usage**: `func (s *Ship) ExtractWithSurvey(survey models.Survey) (*models.Extraction, error)`

### Jettison

- **Purpose**: Jettisons cargo into space.
- **Usage**: `func (s *Ship) Jettison(goodSymbol models.GoodSymbol, units int) (*models.Cargo, error)`

### Jump

- **Purpose**: Commands a ship to jump to another system.
- **Usage**: `func (s *Ship) Jump(systemSymbol string) (*models.ShipNav, *models.ShipCooldown, *models.Transaction, *models.Agent, error)`

### Navigate

- **Purpose**: Navigates a ship to a specified waypoint.
- **Usage**: `func (s *Ship) Navigate(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, []models.Event, error)`

### SetFlightMode

- **Purpose**: Sets the flight mode of a ship.
- **Usage**: `func (s *Ship) SetFlightMode(flightmode models.FlightMode) error`

### FetchNavigationStatus

- **Purpose**: Retrieves the current navigation status of a ship.
- **Usage**: `func (s *Ship) FetchNavigationStatus() (*models.ShipNav, error)`

### Warp

- **Purpose**: Commands a ship to warp to a new location.
- **Usage**: `func (s *Ship) Warp(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, error)`

### SellCargo

- **Purpose**: Sells cargo from a ship's inventory.
- **Usage**: `func (s *Ship) SellCargo(goodSymbol models.GoodSymbol, units int) (*models.Agent, *models.Cargo, *models.Transaction, error)`

### NegotiateContract

- **Purpose**: Negotiates a new contract for the ship.
- **Usage**: `func (s *Ship) NegotiateContract() (*models.Contract, error)`

### GetMounts

- **Purpose**: Retrieves available mounts for the ship.
- **Usage**: `func (s *Ship) GetMounts() (*models.MountSymbol, string, string, int, []string, models.ShipRequirements, error)`

### InstallMount

- **Purpose**: Installs a mount on the ship.
- **Usage**: `func (s *Ship) InstallMount(mountSymbol models.MountSymbol) (*models.Agent, []models.Mount, *models.Cargo, *models.Transaction, error)`

### RemoveMount

- **Purpose**: Removes a mount from the ship.
- **Usage**: `func (s *Ship) RemoveMount(mountSymbol models.MountSymbol) (*models.Agent, []models.Mount, *models.Cargo, *models.Transaction, error)`

### GetScrapPrice

- **Purpose**: Retrieves the current scrap price of the ship.
- **Usage**: `func (s *Ship) GetScrapPrice() (*models.Transaction, error)`

### ScrapShip

- **Purpose**: Scraps the ship for resources.
- **Usage**: `func (s *Ship) ScrapShip() (*models.Transaction, error)`

### GetRepairPrice

- **Purpose**: Retrieves the current repair price of the ship.
- **Usage**: `func (s *Ship) GetRepairPrice() (*models.Transaction, error)`

### RepairShip

- **Purpose**: Repairs the ship.
- **Usage**: `func (s *Ship) RepairShip() (*models.Ship, *models.Transaction, error)`



