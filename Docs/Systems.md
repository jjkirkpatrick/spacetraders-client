# System Operations Guide

This guide provides an overview of how to interact with system-related functionalities using the provided methods in `systems.go`. Each function is designed to perform specific operations related to systems in the SpaceTraders API.

## Getting Started

Before you can use any of the system functions, ensure you have a client instance created and authenticated. Refer to the client setup guide for instructions on initializing and authenticating your client.

## Function Descriptions

### ListSystems

- **Purpose**: Fetches a list of all systems.
- **Usage**: `ListSystems(client *client.Client) ([]*System, error)`

### GetSystem

- **Purpose**: Retrieves detailed information about a specific system by its symbol.
- **Usage**: `GetSystem(client *client.Client, symbol string) (*System, error)`

### ListWaypoints

- **Purpose**: Lists all waypoints within a system, optionally filtered by trait and type.
- **Usage**: `(s *System) ListWaypoints(trait models.WaypointTrait, waypointType models.WaypointType) ([]*models.Waypoint, *models.Meta, error)`

### FetchWaypoint

- **Purpose**: Fetches detailed information about a specific waypoint within a system.
- **Usage**: `(s *System) FetchWaypoint(symbol string) (*models.Waypoint, error)`

### GetMarket

- **Purpose**: Retrieves market information for a specific waypoint within a system.
- **Usage**: `(s *System) GetMarket(waypointSymbol string) (*models.Market, error)`

### GetShipyard

- **Purpose**: Retrieves shipyard information for a specific waypoint within a system.
- **Usage**: `(s *System) GetShipyard(waypointSymbol string) (*models.Shipyard, error)`

### GetJumpGate

- **Purpose**: Retrieves jump gate information for a specific waypoint within a system.
- **Usage**: `(s *System) GetJumpGate(waypointSymbol string) (*models.JumpGate, error)`

### GetConstructionSite

- **Purpose**: Retrieves construction site information for a specific waypoint within a system.
- **Usage**: `(s *System) GetConstructionSite(waypointSymbol string) (*models.ConstructionSite, error)`

### SupplyConstructionSite

- **Purpose**: Supplies a construction site with required goods from a ship.
- **Usage**: `(s *System) SupplyConstructionSite(shipSymbol string, waypointSymbol string, good models.GoodSymbol, quantity int) error`

Each system operation allows you to interact with the SpaceTraders universe's systems, providing functionalities such as listing all available systems, retrieving detailed information about a specific system, and managing waypoints within a system. These operations are crucial for navigating and understanding the vast universe of SpaceTraders.
