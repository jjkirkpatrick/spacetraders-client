# Agent Operations Guide

This guide provides an overview of how to interact with agent-related functionalities using the provided methods in `agents.go`. Each function is designed to perform specific operations related to agents in the SpaceTraders API.

## Getting Started

Before you can use any of the agent functions, ensure you have a client instance created and authenticated. Refer to the client setup guide for instructions on initializing and authenticating your client.

## Function Descriptions

### ListPublicAgents

- **Purpose**: Fetches a list of all public agents.
- **Usage**: `ListPublicAgents(client *client.Client) ([]*Agent, error)`

### GetAgent

- **Purpose**: Retrieves detailed information about the authenticated agent.
- **Usage**: `GetAgent(client *client.Client) (*Agent, error)`

### GetPublicAgent

- **Purpose**: Retrieves detailed information about a specific public agent by its symbol.
- **Usage**: `GetPublicAgent(client *client.Client, symbol string) (*Agent, error)`

Each agent operation allows you to interact with the SpaceTraders universe's agents, providing functionalities such as listing all available public agents, retrieving detailed information about the authenticated agent, and retrieving detailed information about a specific public agent. These operations are crucial for understanding the dynamics and affiliations within the SpaceTraders universe.
