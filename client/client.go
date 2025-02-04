package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jjkirkpatrick/spacetraders-client/internal/cache"
	"github.com/jjkirkpatrick/spacetraders-client/internal/telemetry"
	"github.com/jjkirkpatrick/spacetraders-client/models"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TelemetryOptions represents the configuration options for OpenTelemetry
type TelemetryOptions struct {
	// ServiceName is the name of your service (required if telemetry is enabled)
	ServiceName string
	// ServiceVersion is the version of your service (optional)
	ServiceVersion string
	// Environment is the deployment environment (e.g., "development", "production")
	Environment string
	// OTLPEndpoint is the endpoint for your OpenTelemetry collector (required if telemetry is enabled)
	OTLPEndpoint string
	// MetricInterval is how often metrics are exported (defaults to 15s)
	MetricInterval time.Duration
	// AdditionalAttributes are any extra attributes to add to all telemetry
	AdditionalAttributes map[string]string
	// GRPCDialOptions are additional options for the gRPC connection to the collector
	GRPCDialOptions []grpc.DialOption
}

// ClientOptions represents the configuration options for the SpaceTraders API client
type ClientOptions struct {
	BaseURL           string
	Symbol            string
	Faction           string
	Email             string
	RequestsPerSecond float32
	LogLevel          slog.Level
	RetryDelay        time.Duration
	// Telemetry configuration (optional)
	TelemetryOptions *TelemetryOptions
}

// Client represents the SpaceTraders API client
type Client struct {
	context     context.Context
	baseURL     string
	token       string
	httpClient  *resty.Client
	retryDelay  time.Duration
	AgentSymbol string
	CacheClient *cache.Cache
	Logger      *slog.Logger
	RateLimiter *RateLimiter

	// Telemetry (metrics only)
	telemetryProviders *telemetry.Providers
	meter              metric.Meter
	requestCounter     metric.Int64Counter
	requestDuration    metric.Float64Histogram
	errorCounter       metric.Int64Counter
}

type RateLimiter struct {
	staticLimiter *rate.Limiter
	burstLimiter  *rate.Limiter
	staticReset   time.Duration
	burstReset    time.Duration
}

func NewRateLimiter(staticRate, burstRate float64) *RateLimiter {
	// Adjusting limiters to handle both static and burst rates with updated capacities and durations
	return &RateLimiter{
		staticLimiter: rate.NewLimiter(rate.Limit(staticRate), 2), // 2 requests per second
		burstLimiter:  rate.NewLimiter(rate.Limit(burstRate), 30), // 30 requests per 60 seconds
		staticReset:   time.Second,
		burstReset:    60 * time.Second,
	}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	// Prioritize the static pool over the burst pool
	if rl.staticLimiter.Allow() {
		return nil
	}
	if rl.burstLimiter.Allow() {
		return nil
	}
	return rl.staticLimiter.Wait(ctx)
}

// DefaultClientOptions returns the default configuration options for the SpaceTraders API client
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		BaseURL:           "https://api.spacetraders.io/v2",
		RequestsPerSecond: 2,
		RetryDelay:        1 * time.Second,
		LogLevel:          slog.LevelInfo,
		// Telemetry is disabled by default
		TelemetryOptions: nil,
	}
}

// DefaultTelemetryOptions returns the default configuration options for OpenTelemetry
func DefaultTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Environment:    "development",
		MetricInterval: 15 * time.Second,
		GRPCDialOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		},
	}
}

// NewClient creates a new instance of the SpaceTraders API client
func NewClient(options ClientOptions) (*Client, error) {
	if options.Symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	// Configure basic slog handler
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: options.LogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	})

	// Create initial client with basic logging
	client := &Client{
		baseURL:     options.BaseURL,
		httpClient:  resty.New(),
		context:     context.Background(),
		retryDelay:  options.RetryDelay,
		AgentSymbol: options.Symbol,
		CacheClient: cache.NewCache(),
		Logger:      slog.New(handler),
		RateLimiter: NewRateLimiter(2, 0.5),
	}

	// Initialize telemetry if configured
	if options.TelemetryOptions != nil {
		// Convert public options to internal config
		telemetryConfig := telemetry.Config{
			ServiceName:    options.TelemetryOptions.ServiceName,
			ServiceVersion: options.TelemetryOptions.ServiceVersion,
			Environment:    options.TelemetryOptions.Environment,
			OTLPEndpoint:   options.TelemetryOptions.OTLPEndpoint,
			MetricInterval: options.TelemetryOptions.MetricInterval,
		}

		// Convert additional attributes to KeyValue pairs
		if options.TelemetryOptions.AdditionalAttributes != nil {
			attrs := make([]attribute.KeyValue, 0, len(options.TelemetryOptions.AdditionalAttributes))
			for k, v := range options.TelemetryOptions.AdditionalAttributes {
				attrs = append(attrs, attribute.String(k, v))
			}
			telemetryConfig.AdditionalAttrs = attrs
		}

		// Add gRPC dial options if provided
		if options.TelemetryOptions.GRPCDialOptions != nil {
			telemetryConfig.GRPCDialOptions = options.TelemetryOptions.GRPCDialOptions
		}

		providers, terr := telemetry.InitTelemetry(client.context, telemetryConfig)
		if terr != nil {
			return nil, fmt.Errorf("failed to initialize telemetry: %w", terr)
		}
		client.telemetryProviders = providers

		// Initialize metrics and tracer
		client.meter = otel.GetMeterProvider().Meter("spacetraders-client")

		var merr error
		client.requestCounter, merr = client.meter.Int64Counter("api_requests_total",
			metric.WithDescription("Total number of API requests made"),
			metric.WithUnit("{requests}"),
		)
		if merr != nil {
			return nil, fmt.Errorf("failed to create request counter: %w", merr)
		}

		client.requestDuration, merr = client.meter.Float64Histogram("api_request_duration",
			metric.WithDescription("Duration of API requests"),
			metric.WithUnit("ms"),
		)
		if merr != nil {
			return nil, fmt.Errorf("failed to create request duration histogram: %w", merr)
		}

		client.errorCounter, merr = client.meter.Int64Counter("api_errors_total",
			metric.WithDescription("Total number of API errors"),
			metric.WithUnit("{errors}"),
		)
		if merr != nil {
			return nil, fmt.Errorf("failed to create error counter: %w", merr)
		}
	}

	if apiError := client.getOrRegisterToken(options.Faction, options.Symbol, options.Email); apiError != nil {
		return nil, apiError
	}

	client.Logger.Info("New SpaceTraders client initialized",
		"baseURL", client.baseURL,
		"rateLimit", options.RequestsPerSecond)
	return client, nil
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
	startTime := time.Now()

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
	for {
		c.Logger.Debug("Client Log: Sending request",
			"method", method,
			"url", c.baseURL+endpoint,
			"body", request.Body,
			"params", request.QueryParam)

		err = c.RateLimiter.Wait(c.context)
		if err != nil {
			c.Logger.Error("Client Log: Rate limiter error", "error", err)
			return &models.APIError{Message: err.Error(), Code: 429}
		}

		resp, err = request.Execute(method, c.baseURL+endpoint)
		duration := time.Since(startTime)

		statusCode := 500
		if resp != nil {
			statusCode = resp.StatusCode()
		}
		c.recordMetrics(method, endpoint, duration, statusCode, err)

		if err == nil && !resp.IsError() {
			return nil // Success
		}

		if resp.StatusCode() == 429 {
			c.Logger.Warn("Client Log: Rate limit exceeded, handling rate limit")
			handleRateLimit(resp, c.Logger)
			continue
		}

		if resp.IsError() {
			apiError = parseAPIError(resp)
			c.Logger.Error("Client Log: API Request resulted in error",
				"error", apiError.Error(),
				"data", apiError.Data)
			return apiError
		}

		backoff = applyJitter(backoff)
		time.Sleep(backoff)
		backoff *= 2
	}
}

func (c *Client) recordMetrics(method, endpoint string, duration time.Duration, statusCode int, err error) {
	if c.meter == nil {
		return // Telemetry is disabled
	}

	attrs := []attribute.KeyValue{
		attribute.String("agent", c.AgentSymbol),
		attribute.String("endpoint", endpoint),
		attribute.String("method", method),
	}

	c.requestCounter.Add(c.context, 1, metric.WithAttributes(attrs...))
	c.requestDuration.Record(c.context, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))

	if err != nil || statusCode >= 400 {
		errorAttrs := append(attrs,
			attribute.Int("status_code", statusCode),
		)
		if err != nil {
			errorAttrs = append(errorAttrs, attribute.String("error", err.Error()))
		}
		c.errorCounter.Add(c.context, 1, metric.WithAttributes(errorAttrs...))
	}
}

// Helper function to apply jitter to backoff
func applyJitter(backoff time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return backoff + jitter/2
}

// Helper function to handle rate limit errors
func handleRateLimit(resp *resty.Response, logger *slog.Logger) {
	resetTime := resp.Header().Get("x-ratelimit-reset")
	if resetTime != "" {
		if resetTimestamp, parseErr := strconv.ParseInt(resetTime, 10, 64); parseErr == nil {
			waitDuration := time.Until(time.Unix(resetTimestamp, 0))
			logger.Debug("Rate limit exceeded, waiting until reset",
				"wait_duration", waitDuration)
			time.Sleep(waitDuration)
		}
	}
}

// Helper function to parse API error from response
func parseAPIError(resp *resty.Response) *models.APIError {
	var errorWrapper struct {
		Error struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		} `json:"error"`
	}

	err := json.Unmarshal(resp.Body(), &errorWrapper)
	if err != nil {
		return &models.APIError{
			Message: "failed to parse API error response",
			Code:    resp.StatusCode(),
		}
	}

	return &models.APIError{
		Code:    errorWrapper.Error.Code,
		Message: errorWrapper.Error.Message,
		Data:    errorWrapper.Error.Data,
	}
}

// Adjust isRateLimitError to include more transient errors if needed
func isRateLimitError(statusCode int) bool {
	return statusCode == 429 || statusCode == 502
}

func (c *Client) GetToken() string {
	return c.token
}

// Close gracefully shuts down the client and its telemetry providers
func (c *Client) Close(ctx context.Context) error {
	if c.telemetryProviders != nil {
		return c.telemetryProviders.Shutdown(ctx)
	}
	return nil
}
