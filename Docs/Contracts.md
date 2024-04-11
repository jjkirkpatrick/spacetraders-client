# Contract Operations Guide

This guide provides an overview of how to interact with contract-related functionalities using the provided methods in `contracts.go`. Each function is designed to perform specific operations related to contracts in the SpaceTraders API.

## Getting Started

Before you can use any of the contract functions, ensure you have a client instance created and authenticated. Refer to the client setup guide for instructions on initializing and authenticating your client.

## Function Descriptions

### ListContracts

- **Purpose**: Fetches a list of all contracts available to the player.
- **Usage**: `ListContracts(client *client.Client) ([]*Contract, error)`

### GetContract

- **Purpose**: Retrieves detailed information about a specific contract by its symbol.
- **Usage**: `GetContract(client *client.Client, symbol string) (*Contract, error)`

### Accept

- **Purpose**: Accepts a contract, allowing the player to start fulfilling its requirements.
- **Usage**: `(c *Contract) Accept() (*models.Agent, *models.Contract, error)`

### DeliverCargo

- **Purpose**: Delivers cargo to fulfill a contract's requirements.
- **Usage**: `(c *Contract) DeliverCargo(shop Ship, tradeGood models.GoodSymbol, units int) (*models.Contract, *models.Cargo, error)`

### Fulfill

- **Purpose**: Marks a contract as fulfilled, completing it and receiving the rewards.
- **Usage**: `(c *Contract) Fulfill() (*models.Agent, *models.Contract, error)`

Each contract operation allows you to interact with the SpaceTraders universe's contracts, providing functionalities such as listing all available contracts, retrieving detailed information about a specific contract, accepting contracts, delivering cargo to fulfill contract requirements, and marking contracts as fulfilled. These operations are crucial for progressing and gaining rewards within the SpaceTraders universe.
