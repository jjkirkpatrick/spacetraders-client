package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jjkirkpatrick/spacetraders-client/api"
	"github.com/jjkirkpatrick/spacetraders-client/models"
	"golang.org/x/time/rate"
)

// ClientOptions represents the configuration options for the SpaceTraders API client
type ClientOptions struct {
	BaseURL           string
	Token             string
	RequestsPerSecond float64
	RetryCount        int
	RetryDelay        time.Duration
	Logger            *log.Logger
}

// DefaultClientOptions returns the default configuration options for the SpaceTraders API client
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		BaseURL:           "https://api.spacetraders.io/v2",
		RequestsPerSecond: 2,
		RetryCount:        3,
		RetryDelay:        time.Second,
		Logger:            log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Client represents the SpaceTraders API client
type Client struct {
	baseURL    string
	token      string
	httpClient *resty.Client
	limiter    *rate.Limiter
	context    context.Context
	retryCount int
	retryDelay time.Duration
	logger     *log.Logger
}

// NewClient creates a new instance of the SpaceTraders API client
func NewClient(options ClientOptions) (*Client, error) {
	if options.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	return &Client{
		baseURL:    options.BaseURL,
		token:      options.Token,
		httpClient: resty.New(),
		limiter:    rate.NewLimiter(rate.Limit(options.RequestsPerSecond), 10),
		context:    context.Background(),
		retryCount: options.RetryCount,
		retryDelay: options.RetryDelay,
		logger:     options.Logger,
	}, nil
}

// NewClientWithAgentRegistration creates a new instance of the SpaceTraders API client and registers a new agent
func NewClientWithAgentRegistration(options ClientOptions, faction, symbol, email string) (*Client, error) {

	client := &Client{
		baseURL:    options.BaseURL,
		httpClient: resty.New(),
		limiter:    rate.NewLimiter(rate.Limit(options.RequestsPerSecond), 10),
		context:    context.Background(),
		retryCount: options.RetryCount,
		retryDelay: options.RetryDelay,
		logger:     options.Logger,
	}

	err := client.RegisterNewAgent(faction, symbol, email)
	if err != nil {
		return nil, fmt.Errorf("failed to register new agent: %v", err)
	}

	return client, nil
}

// SetBaseURL sets the base URL for the API client
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// SetToken sets the authentication token for the API client
func (c *Client) SetToken(token string) {
	c.token = token
}

// Get sends a GET request to the specified endpoint with optional query parameters
func (c *Client) Get(endpoint string, queryParams map[string]string, result interface{}) error {
	return c.sendRequest("GET", endpoint, nil, queryParams, result)
}

// Post sends a POST request to the specified endpoint with optional query parameters
func (c *Client) Post(endpoint string, body interface{}, queryParams map[string]string, result interface{}) error {
	return c.sendRequest("POST", endpoint, body, queryParams, result)
}

// Put sends a PUT request to the specified endpoint with optional query parameters
func (c *Client) Put(endpoint string, body interface{}, queryParams map[string]string, result interface{}) error {
	return c.sendRequest("PUT", endpoint, body, queryParams, result)
}

// Delete sends a DELETE request to the specified endpoint with optional query parameters
func (c *Client) Delete(endpoint string, queryParams map[string]string, result interface{}) error {
	return c.sendRequest("DELETE", endpoint, nil, queryParams, result)
}

// Patch sends a PATCH request to the specified endpoint with optional query parameters
func (c *Client) Patch(endpoint string, body interface{}, queryParams map[string]string, result interface{}) error {
	return c.sendRequest("PATCH", endpoint, body, queryParams, result)
}

// sendRequest is a helper method to send requests with rate limiting, automatic retry on rate limit errors, pagination support, and query parameters
func (c *Client) sendRequest(method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) error {
	err := c.limiter.Wait(c.context)
	if err != nil {
		return fmt.Errorf("rate limit exceeded: %v", err)
	}

	request := c.httpClient.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(c.token).
		SetResult(result)

	if body != nil {
		request.SetBody(body)
	}

	if queryParams != nil {
		request.SetQueryParams(queryParams)
	}

	var resp *resty.Response
	for i := 0; i <= c.retryCount; i++ {
		switch method {
		case "GET":
			resp, err = request.Get(c.baseURL + endpoint)
		case "POST":
			resp, err = request.Post(c.baseURL + endpoint)
		case "PUT":
			resp, err = request.Put(c.baseURL + endpoint)
		case "DELETE":
			resp, err = request.Delete(c.baseURL + endpoint)
		case "PATCH":
			resp, err = request.Patch(c.baseURL + endpoint)
		default:
			return fmt.Errorf("unsupported HTTP method: %s", method)
		}

		if err != nil {
			if resp != nil && isRateLimitError(resp.StatusCode()) {
				if i < c.retryCount {
					c.logger.Printf("Rate limit exceeded. Retrying in %v...", c.retryDelay)
					time.Sleep(c.retryDelay)
					continue
				}
			}
			return err
		}

		if resp.IsError() {
			apiError := fmt.Errorf("API error: %s", resp.Status()+" "+string(resp.Body()))
			return apiError
		}

		break
	}

	return nil
}

// isRateLimitError checks if the given status code represents a rate limit error
func isRateLimitError(statusCode int) bool {
	return statusCode == 429 || (statusCode >= 502 && statusCode <= 599)
}

func (c *Client) GetToken() string {
	return c.token
}

// GetAgent retrieves the agent's details
func (c *Client) GetAgent() (*models.Agent, error) {
	return api.GetAgent(c.Get)
}

func (c *Client) GetPublicAgent(agentSymbol string) (*models.Agent, error) {
	return api.GetPublicAgent(c.Get, agentSymbol)
}

func (c *Client) ListAgents() (*Paginator[*models.Agent], error) {
	fetchFunc := func(meta models.Meta) ([]*models.Agent, models.Meta, error) {
		// Since api.ListAgents expects a pointer to models.Meta, create a pointer from the value.
		metaPtr := &meta
		// Call api.ListAgents with a pointer to meta.
		agents, metaPtr, err := api.ListAgents(c.Get, metaPtr)
		// Dereference metaPtr when returning to match the expected return types.
		return agents, *metaPtr, err
	}
	// Initialize the paginator with the fetch function.
	return NewPaginator[*models.Agent](fetchFunc).FetchFirstPage()
}

func (c *Client) ListContracts() (*Paginator[*models.Contract], error) {
	fetchFunc := func(meta models.Meta) ([]*models.Contract, models.Meta, error) {
		// Since api.ListAgents expects a pointer to models.Meta, create a pointer from the value.
		metaPtr := &meta
		// Call api.ListAgents with a pointer to meta.
		agents, metaPtr, err := api.ListContracts(c.Get, metaPtr)
		// Dereference metaPtr when returning to match the expected return types.
		return agents, *metaPtr, err
	}
	// Initialize the paginator with the fetch function.
	return NewPaginator[*models.Contract](fetchFunc).FetchFirstPage()
}

func (c *Client) GetContract(contractId string) (*models.Contract, error) {
	return api.GetContract(c.Get, contractId)
}

func (c *Client) AcceptContract(contractId string) (*models.Agent, *models.Contract, error) {
	agent, contract, err := api.AcceptContract(c.Post, contractId)
	return agent, contract, err
}

func (c *Client) DeliverContractCargo(contractId string, body models.DeliverContractCargoRequest) (*models.Contract, *models.Cargo, error) {
	contract, cargo, err := api.DeliverContractCargo(c.Post, contractId, body)
	return contract, cargo, err
}

func (c *Client) FulfilContract(contractId string) (*models.Agent, *models.Contract, error) {
	agent, contract, err := api.FulfillContract(c.Post, contractId)
	return agent, contract, err
}
