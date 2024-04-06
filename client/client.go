package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jjkirkpatrick/spacetraders-client/api"
	"github.com/jjkirkpatrick/spacetraders-client/metrics"
	"github.com/jjkirkpatrick/spacetraders-client/models"
	"golang.org/x/time/rate"
)

// Client represents the SpaceTraders API client
type Client struct {
	baseURL         string
	token           string
	httpClient      *resty.Client
	context         context.Context
	retryCount      int
	retryDelay      time.Duration
	metricsReporter metrics.MetricsReporter // Use the interface
	logger          *log.Logger
}

// ClientOptions represents the configuration options for the SpaceTraders API client
type ClientOptions struct {
	BaseURL           string
	Token             string
	RequestsPerSecond float32
	RetryCount        int
	RetryDelay        time.Duration
	Logger            *log.Logger
}

// DefaultClientOptions returns the default configuration options for the SpaceTraders API client
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		BaseURL:           "https://api.spacetraders.io/v2",
		RetryCount:        3,
		Logger:            log.New(os.Stdout, "", log.LstdFlags),
		RequestsPerSecond: 2,
		RetryDelay:        1 * time.Second,
	}
}

// NewClient creates a new instance of the SpaceTraders API client
func NewClient(options ClientOptions, metricsReporter metrics.MetricsReporter) (*Client, error) {
	if options.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	if metricsReporter == nil {
		metricsReporter = &metrics.NoOpMetricsReporter{}
	}

	client := &Client{
		baseURL:         options.BaseURL,
		token:           options.Token,
		httpClient:      resty.New(),
		context:         context.Background(),
		retryCount:      options.RetryCount,
		retryDelay:      options.RetryDelay,
		metricsReporter: metricsReporter,
		logger:          options.Logger,
	}

	client.httpClient.SetRateLimiter(rate.NewLimiter(rate.Limit(options.RequestsPerSecond), 10))

	return client, nil
}

// NewClientWithAgentRegistration creates a new instance of the SpaceTraders API client and registers a new agent
func NewClientWithAgentRegistration(options ClientOptions, faction, symbol, email string, metricsReporter metrics.MetricsReporter) (*Client, error) {

	if metricsReporter == nil {
		metricsReporter = &metrics.NoOpMetricsReporter{}
	}

	client := &Client{
		baseURL:         options.BaseURL,
		httpClient:      resty.New(),
		context:         context.Background(),
		retryCount:      options.RetryCount,
		retryDelay:      options.RetryDelay,
		metricsReporter: metricsReporter,
		logger:          options.Logger,
	}

	client.httpClient.SetRateLimiter(rate.NewLimiter(rate.Limit(options.RequestsPerSecond), 10))

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
func (c *Client) Get(endpoint string, queryParams map[string]string, result interface{}) *models.APIError {
	return c.sendRequest("GET", endpoint, nil, queryParams, result)
}

// Post sends a POST request to the specified endpoint with optional query parameters
func (c *Client) Post(endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
	return c.sendRequest("POST", endpoint, body, queryParams, result)
}

// Put sends a PUT request to the specified endpoint with optional query parameters
func (c *Client) Put(endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
	return c.sendRequest("PUT", endpoint, body, queryParams, result)
}

// Delete sends a DELETE request to the specified endpoint with optional query parameters
func (c *Client) Delete(endpoint string, queryParams map[string]string, result interface{}) *models.APIError {
	return c.sendRequest("DELETE", endpoint, nil, queryParams, result)
}

// Patch sends a PATCH request to the specified endpoint with optional query parameters
func (c *Client) Patch(endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
	return c.sendRequest("PATCH", endpoint, body, queryParams, result)
}

func (c *Client) sendRequest(method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {

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
	var apiError *models.APIError
	var err error

	backoff := c.retryDelay
	for i := 0; i <= c.retryCount; i++ {
		resp, err = request.Execute(method, c.baseURL+endpoint)
		metric, _ := metrics.NewMetricBuilder().
			Namespace("api_request").
			Tag("method", method).
			Field("count", 1).
			Timestamp(time.Now()).
			Build()
		c.metricsReporter.WritePoint(metric)
		if err != nil {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
			continue
		}

		if resp.IsError() {
			if isRateLimitError(resp.StatusCode()) {
				handleRateLimit(resp, c.logger)
				metric, _ := metrics.NewMetricBuilder().
					Namespace("api_request_error").
					Tag("method", method).
					Tag("error_type", "rate_limit").
					Tag("status_code", fmt.Sprintf("%d", resp.StatusCode())).
					Field("count", 1).
					Timestamp(time.Now()).
					Build()
				c.metricsReporter.WritePoint(metric)

				continue
			} else {
				apiError = parseAPIError(resp)
				if apiError != nil {
					metric, _ := metrics.NewMetricBuilder().
						Namespace("api_request_error").
						Tag("method", method).
						Tag("error_type", "api_error").
						Tag("status_code", fmt.Sprintf("%d", resp.StatusCode())).
						Field("count", 1).
						Timestamp(time.Now()).
						Build()
					c.metricsReporter.WritePoint(metric)

				}
				if apiError != nil || i == c.retryCount {
					break // Break if we have a parsed error or are on the last retry
				}
			}
		} else {
			return nil // Success
		}

		// Apply random jitter to backoff
		backoff = applyJitter(backoff)
		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff
	}

	if apiError == nil && err != nil {
		apiError = &models.APIError{Message: err.Error(), Code: 500}
		metric, _ := metrics.NewMetricBuilder().
			Namespace("api_request_error").
			Tag("method", method).
			Tag("error_type", "unknown_error").
			Tag("status_code", "500").
			Field("count", 1).
			Timestamp(time.Now()).
			Build()
		c.metricsReporter.WritePoint(metric)
	}
	return apiError
}

// Helper function to apply jitter to backoff
func applyJitter(backoff time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return backoff + jitter/2
}

// Helper function to handle rate limit errors
func handleRateLimit(resp *resty.Response, logger *log.Logger) {
	resetTime := resp.Header().Get("x-ratelimit-reset")
	if resetTime != "" {
		if resetTimestamp, parseErr := strconv.ParseInt(resetTime, 10, 64); parseErr == nil {
			waitDuration := time.Until(time.Unix(resetTimestamp, 0))
			logger.Printf("Rate limit exceeded. Waiting until reset: %v", waitDuration)
			time.Sleep(waitDuration)
		}
	}
}

// Helper function to parse API error from response
func parseAPIError(resp *resty.Response) *models.APIError {
	var errorWrapper struct {
		Error models.APIError `json:"error"`
	}
	err := json.Unmarshal(resp.Body(), &errorWrapper)
	if err != nil {
		return &models.APIError{Message: "failed to parse API error response", Code: resp.StatusCode()}
	}
	return &errorWrapper.Error
}

// Adjust isRateLimitError to include more transient errors if needed
func isRateLimitError(statusCode int) bool {
	return statusCode == 429 || statusCode == 502
}

func (c *Client) WriteMetric(metric metrics.Metric) {
	c.metricsReporter.WritePoint(metric)
}

func (c *Client) GetToken() string {
	return c.token
}

// GetAgent retrieves the agent's details
func (c *Client) GetAgent() (*models.Agent, *models.APIError) {
	return api.GetAgent(c.Get)
}

func (c *Client) GetPublicAgent(agentSymbol string) (*models.Agent, *models.APIError) {
	return api.GetPublicAgent(c.Get, agentSymbol)
}

func (c *Client) ListAgents() (*Paginator[*models.Agent], *models.APIError) {
	fetchFunc := func(meta models.Meta) ([]*models.Agent, models.Meta, *models.APIError) {
		metaPtr := &meta
		agents, metaPtr, err := api.ListAgents(c.Get, metaPtr)
		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return agents, *metaPtr, err
		}
		if metaPtr != nil {
			return agents, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return agents, defaultMeta, nil
		}
	}
	return NewPaginator[*models.Agent](fetchFunc), nil
}

func (c *Client) ListContracts() (*Paginator[*models.Contract], *models.APIError) {
	fetchFunc := func(meta models.Meta) ([]*models.Contract, models.Meta, *models.APIError) {
		metaPtr := &meta
		contracts, metaPtr, err := api.ListContracts(c.Get, metaPtr)
		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return contracts, *metaPtr, err
		}
		if metaPtr != nil {
			return contracts, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return contracts, defaultMeta, nil
		}
	}
	return NewPaginator[*models.Contract](fetchFunc), nil
}

func (c *Client) GetContract(contractId string) (*models.Contract, *models.APIError) {
	return api.GetContract(c.Get, contractId)
}

func (c *Client) AcceptContract(contractId string) (*models.Agent, *models.Contract, *models.APIError) {
	agent, contract, err := api.AcceptContract(c.Post, contractId)
	return agent, contract, err
}

func (c *Client) DeliverContractCargo(contractId string, body models.DeliverContractCargoRequest) (*models.Contract, *models.Cargo, *models.APIError) {
	contract, cargo, err := api.DeliverContractCargo(c.Post, contractId, body)
	return contract, cargo, err
}

func (c *Client) FulfilContract(contractId string) (*models.Agent, *models.Contract, *models.APIError) {
	agent, contract, err := api.FulfillContract(c.Post, contractId)
	return agent, contract, err
}

func (c *Client) ListSystems() (*Paginator[*models.System], *models.APIError) {
	fetchFunc := func(meta models.Meta) ([]*models.System, models.Meta, *models.APIError) {
		metaPtr := &meta
		systems, metaPtr, err := api.ListSystems(c.Get, metaPtr)
		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return systems, *metaPtr, err
		}
		if metaPtr != nil {
			return systems, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return systems, defaultMeta, nil
		}
	}
	return NewPaginator[*models.System](fetchFunc), nil
}

func (c *Client) GetSystem(systemSymbol string) (*models.System, *models.APIError) {
	return api.GetSystem(c.Get, systemSymbol)
}
func (c *Client) ListWaypointsInSystem(systemSymbol string, trait models.WaypointTrait, waypointType models.WaypointType) (*Paginator[*models.Waypoint], *models.APIError) {
	fetchFunc := func(meta models.Meta) ([]*models.Waypoint, models.Meta, *models.APIError) {
		metaPtr := &meta
		waypoint, metaPtr, err := api.ListWaypointsInSystem(c.Get, metaPtr, systemSymbol, trait, waypointType)
		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return waypoint, *metaPtr, err
		}
		if metaPtr != nil {
			return waypoint, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return waypoint, defaultMeta, nil
		}
	}
	return NewPaginator[*models.Waypoint](fetchFunc), nil
}

func (c *Client) GetWaypoint(systemSymbol, waypointSymbol string) (*models.Waypoint, *models.APIError) {
	return api.GetWaypoint(c.Get, systemSymbol, waypointSymbol)
}

func (c *Client) GetMarket(systemSymbol, waypointSymbol string) (*models.Market, *models.APIError) {
	return api.GetMarket(c.Get, systemSymbol, waypointSymbol)
}

func (c *Client) GetShipyard(systemSymbol, waypointSymbol string) (*models.Shipyard, *models.APIError) {
	return api.GetShipyard(c.Get, systemSymbol, waypointSymbol)
}

func (c *Client) GetJumpGate(systemSymbol, waypointSymbol string) (*models.JumpGate, *models.APIError) {
	return api.GetJumpGate(c.Get, systemSymbol, waypointSymbol)
}

func (c *Client) GetConstructionSite(systemSymbol, waypointSymbol string) (*models.ConstructionSite, *models.APIError) {
	return api.GetConstructionSite(c.Get, systemSymbol, waypointSymbol)
}

func (c *Client) SupplyConstructionSite(systemSymbol, waypointSymbol string, payload models.SupplyConstructionSiteRequest) (*models.ConstructionSite, *models.APIError) {
	return api.SupplyConstructionSite(c.Post, systemSymbol, waypointSymbol, payload)
}

// Functions from fleet.go

func (c *Client) ListShips(systemSymbol string) (*Paginator[*models.Ship], *models.APIError) {
	fetchFunc := func(meta models.Meta) ([]*models.Ship, models.Meta, *models.APIError) {
		metaPtr := &meta
		ships, metaPtr, err := api.ListShips(c.Get, metaPtr)
		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return ships, *metaPtr, err
		}
		if metaPtr != nil {
			return ships, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return ships, defaultMeta, nil
		}
	}
	return NewPaginator[*models.Ship](fetchFunc), nil
}

func (c *Client) PurchaseShip(payload *models.PurchaseShipRequest) (*models.PurchaseShipResponse, *models.APIError) {
	return api.PurchaseShip(c.Post, payload)
}

func (c *Client) GetShip(ShipSymbol string) (*models.Ship, *models.APIError) {
	return api.GetShip(c.Get, ShipSymbol)
}

func (c *Client) GetShipCargo(ShipSymbol string) (*models.Cargo, *models.APIError) {
	return api.GetShipCargo(c.Get, ShipSymbol)
}

func (c *Client) OrbitShip(ShipSymbol string, payload *models.OrbitRequest) (*models.ShipNav, *models.APIError) {
	return api.OrbitShip(c.Post, ShipSymbol, payload)
}

func (c *Client) ShipRefine(ShipSymbol string, payload *models.RefineRequest) (*models.ShipRefineResponse, *models.APIError) {
	return api.ShipRefine(c.Post, ShipSymbol, payload)
}

func (c *Client) CreateChart(ShipSymbol string) (*models.CreateChartResponse, *models.APIError) {
	return api.CreateChart(c.Post, ShipSymbol)
}

func (c *Client) GetShipCooldown(ShipSymbol string) (*models.ShipCooldown, *models.APIError) {
	return api.GetShipCooldown(c.Get, ShipSymbol)
}

func (c *Client) DockShip(ShipSymbol string) (*models.ShipNav, *models.APIError) {
	return api.DockShip(c.Post, ShipSymbol)
}

func (c *Client) CreateSurvey(ShipSymbol string) (*models.CreateSurveyResponse, *models.APIError) {
	return api.CreateSurvey(c.Post, ShipSymbol)
}

func (c *Client) ExtractResources(ShipSymbol string, payload *models.Survey) (*models.ExtractionResponse, *models.APIError) {
	return api.ExtractResources(c.Post, ShipSymbol, payload)
}

func (c *Client) SiphonResources(ShipSymbol string) (*models.SiphonResponse, *models.APIError) {
	return api.SiphonResources(c.Post, ShipSymbol)
}

func (c *Client) ExtractResourcesWithSurvey(ShipSymbol string, payload *models.ExtractWithSurveyRequest) (*models.ExtractionResponse, *models.APIError) {
	return api.ExtractResourcesWithSurvey(c.Post, ShipSymbol, payload)
}

func (c *Client) JettisonCargo(ShipSymbol string, payload *models.JettisonRequest) (*models.JettisonResponse, *models.APIError) {
	return api.JettisonCargo(c.Post, ShipSymbol, payload)
}

func (c *Client) JumpShip(ShipSymbol string, payload *models.JumpShipRequest) (*models.JumpShipResponse, *models.APIError) {
	return api.JumpShip(c.Post, ShipSymbol, payload)
}

func (c *Client) NavigateShip(ShipSymbol string, payload *models.NavigateRequest) (*models.NavigateResponse, *models.APIError) {
	return api.NavigateShip(c.Post, ShipSymbol, payload)
}

func (c *Client) PatchShipNav(ShipSymbol string, payload *models.NavUpdateRequest) (*models.PatchShipNacResponse, *models.APIError) {
	return api.PatchShipNav(c.Patch, ShipSymbol, payload)
}

func (c *Client) GetShipNav(ShipSymbol string) (*models.ShipNav, *models.APIError) {
	return api.GetShipNav(c.Get, ShipSymbol)
}

func (c *Client) WarpShip(ShipSymbol string, payload *models.WarpRequest) (*models.WarpResponse, *models.APIError) {
	return api.WarpShip(c.Post, ShipSymbol, payload)
}

func (c *Client) SellCargo(ShipSymbol string, payload *models.SellCargoRequest) (*models.SellCargoResponse, *models.APIError) {
	return api.SellCargo(c.Post, ShipSymbol, payload)
}

func (c *Client) ScanSystems(ShipSymbol string) (*models.ScanSystemsResponse, *models.APIError) {
	return api.ScanSystems(c.Post, ShipSymbol)
}

func (c *Client) ScanWaypoints(ShipSymbol string) (*models.ScanWaypointsResponse, *models.APIError) {
	return api.ScanWaypoints(c.Post, ShipSymbol)
}

func (c *Client) ScanShips(ShipSymbol string) (*models.ScanShipsResponse, *models.APIError) {
	return api.ScanShips(c.Post, ShipSymbol)
}

func (c *Client) RefuelShip(ShipSymbol string, payload *models.RefuelShipRequest) (*models.RefuelShipResponse, *models.APIError) {
	return api.RefuelShip(c.Post, ShipSymbol, payload)
}

func (c *Client) PurchaseCargo(ShipSymbol string, payload *models.PurchaseCargoRequest) (*models.PurchaseCargoResponse, *models.APIError) {
	return api.PurchaseCargo(c.Post, ShipSymbol, payload)
}

func (c *Client) TransferCargo(ShipSymbol string, payload *models.TransferCargoRequest) (*models.TransferCargoResponse, *models.APIError) {
	return api.TransferCargo(c.Post, ShipSymbol, payload)
}

func (c *Client) NegotiateContract(ShipSymbol string) (*models.NegotiateContractResponse, *models.APIError) {
	return api.NegotiateContract(c.Post, ShipSymbol)
}

func (c *Client) GetMounts(ShipSymbol string) (*models.GetMountsResponse, *models.APIError) {
	return api.GetMounts(c.Get, ShipSymbol)
}

func (c *Client) InstallMount(ShipSymbol string, payload *models.InstallMountRequest) (*models.InstallMountResponse, *models.APIError) {
	return api.InstallMount(c.Post, ShipSymbol, payload)
}

func (c *Client) RemoveMount(ShipSymbol string, payload *models.RemoveMountRequest) (*models.RemoveMountResponse, *models.APIError) {
	return api.RemoveMount(c.Post, ShipSymbol, payload)
}

func (c *Client) GetScrapShip(ShipSymbol string) (*models.GetScrapShipResponse, *models.APIError) {
	return api.GetScrapShip(c.Get, ShipSymbol)
}

func (c *Client) ScrapShip(ShipSymbol string) (*models.ScrapShipResponse, *models.APIError) {
	return api.ScrapShip(c.Post, ShipSymbol)
}

func (c *Client) GetRepairShip(ShipSymbol string) (*models.GetRepairShipResponse, *models.APIError) {
	return api.GetRepairShip(c.Get, ShipSymbol)
}

func (c *Client) RepairShip(ShipSymbol string) (*models.RepairShipResponse, *models.APIError) {
	return api.RepairShip(c.Post, ShipSymbol)
}

// Functions from factions.go

// GetFaction retrieves the faction's details
// API Docs: https://spacetraders.stoplight.io/docs/spacetraders/a50decd0f9483-get-faction
func (c *Client) GetFaction(factionSymbol string) (*models.Faction, *models.APIError) {
	faction, err := api.GetFaction(c.Get, factionSymbol)
	if err != nil {
		return nil, &models.APIError{Message: err.Error()}
	}
	return faction, nil
}

// ListFactions retrieves a list of factions with pagination
// API Docs: https://spacetraders.stoplight.io/docs/spacetraders/93c5d5e6ad5b0-list-factions
func (c *Client) ListFactions() (*Paginator[*models.Faction], *models.APIError) {
	fetchFunc := func(meta models.Meta) ([]*models.Faction, models.Meta, *models.APIError) {
		metaPtr := &meta
		factions, metaPtr, err := api.ListFactions(c.Get, metaPtr)
		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return factions, *metaPtr, err
		}
		if metaPtr != nil {
			return factions, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return factions, defaultMeta, nil
		}
	}
	return NewPaginator[*models.Faction](fetchFunc), nil
}
