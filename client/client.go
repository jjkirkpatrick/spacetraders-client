package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
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
	rateLimitGauge     metric.Float64ObservableGauge
	remainingRequests  metric.Int64ObservableGauge
	resetTimeGauge     metric.Float64ObservableGauge
}

type RateLimiter struct {
	staticLimiter *rate.Limiter
	burstLimiter  *rate.Limiter
	mu            sync.Mutex
	// Track API-provided limits
	limitPerSecond float64
	limitBurst     int
	// Track remaining requests
	remaining int64
	resetTime time.Time
	// Add a channel to coordinate waiting for reset
	resetChan chan struct{}
}

func NewRateLimiter(staticRate, burstRate float64) *RateLimiter {
	return &RateLimiter{
		staticLimiter:  rate.NewLimiter(rate.Limit(staticRate), 2), // Allow bursting up to 2 tokens
		burstLimiter:   rate.NewLimiter(rate.Limit(burstRate), 30), // Fallback burst limiter for spikes
		limitPerSecond: staticRate,
		limitBurst:     30,
		remaining:      30,
		resetTime:      time.Now().Add(time.Second),
		resetChan:      make(chan struct{}),
	}
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
	rl.mu.Lock()
	// If we have no remaining requests, we need to wait for reset
	if rl.remaining <= 0 {
		resetDuration := time.Until(rl.resetTime)
		if resetDuration > 0 {
			rl.mu.Unlock()
			// Add a smaller buffer to ensure we're past the reset
			waitDuration := resetDuration + 15*time.Millisecond
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitDuration):
				// After waiting, reacquire lock and reset remaining
				rl.mu.Lock()
				rl.remaining = int64(rl.limitBurst)
			}
		}
	}

	// Use the static limiter to maintain steady rate
	err := rl.staticLimiter.Wait(ctx)
	if err != nil {
		rl.mu.Unlock()
		return err
	}

	// Only decrement remaining if we're getting close to the limit
	// This allows full utilization of the 2/s rate while still protecting against bursts
	if rl.remaining < int64(rl.limitBurst/2) {
		rl.remaining--
	}

	rl.mu.Unlock()
	return nil
}

func (rl *RateLimiter) updateLimits(limitPerSecond float64, limitBurst int, remaining int64, resetTime time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limitPerSecond > 0 && limitPerSecond != rl.limitPerSecond {
		rl.limitPerSecond = limitPerSecond
		// Allow bursting up to 2 requests, which matches the API's per-second rate
		rl.staticLimiter = rate.NewLimiter(rate.Limit(limitPerSecond), 2)
	}

	if limitBurst > 0 {
		rl.limitBurst = limitBurst
	}

	rl.remaining = remaining
	rl.resetTime = resetTime

	// If we're at 0 remaining, start a timer to reset
	if rl.remaining <= 0 && !resetTime.IsZero() {
		go func() {
			waitDuration := time.Until(resetTime) + time.Millisecond
			time.Sleep(waitDuration)
			rl.mu.Lock()
			rl.remaining = int64(rl.limitBurst)
			rl.mu.Unlock()
			// Notify any waiters
			select {
			case rl.resetChan <- struct{}{}:
			default:
			}
		}()
	}
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
		MetricInterval: 1 * time.Second,
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
		RateLimiter: NewRateLimiter(2, 30),
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

		client.requestDuration, merr = client.meter.Float64Histogram("api_request_duration_seconds",
			metric.WithDescription("Duration of API requests in seconds"),
			metric.WithUnit("s"),
			metric.WithExplicitBucketBoundaries(0.01, 0.05, 0.1, 0.5, 1, 2.5, 5, 10),
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

		client.rateLimitGauge, merr = client.meter.Float64ObservableGauge("api_rate_limit",
			metric.WithDescription("Current API rate limit settings"),
			metric.WithUnit("{requests_per_second}"),
		)
		if merr != nil {
			return nil, fmt.Errorf("failed to create rate limit gauge: %w", merr)
		}

		client.remainingRequests, merr = client.meter.Int64ObservableGauge("api_remaining_requests",
			metric.WithDescription("Number of API requests remaining before rate limit"),
			metric.WithUnit("{requests}"),
		)
		if merr != nil {
			return nil, fmt.Errorf("failed to create remaining requests gauge: %w", merr)
		}

		client.resetTimeGauge, merr = client.meter.Float64ObservableGauge("api_rate_limit_reset",
			metric.WithDescription("Time until rate limit reset in seconds"),
			metric.WithUnit("s"),
		)
		if merr != nil {
			return nil, fmt.Errorf("failed to create reset time gauge: %w", merr)
		}

		_, err := client.meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
			o.ObserveFloat64(client.rateLimitGauge, client.RateLimiter.limitPerSecond,
				metric.WithAttributes(
					attribute.String("type", "static"),
					attribute.String("agent", client.AgentSymbol),
				))
			o.ObserveInt64(client.remainingRequests, client.RateLimiter.remaining,
				metric.WithAttributes(
					attribute.String("type", "static"),
					attribute.String("agent", client.AgentSymbol),
				))
			resetTime := client.RateLimiter.resetTime
			if !resetTime.IsZero() {
				o.ObserveFloat64(client.resetTimeGauge, time.Until(resetTime).Seconds(),
					metric.WithAttributes(
						attribute.String("agent", client.AgentSymbol),
					))
			}
			return nil
		}, client.rateLimitGauge, client.remainingRequests, client.resetTimeGauge)
		if err != nil {
			return nil, fmt.Errorf("failed to register metric callbacks: %w", err)
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
	var rateLimit *RateLimitResponse

	// Wait for rate limit token - this will block until we can make the request
	if err := c.RateLimiter.Wait(c.context); err != nil {
		c.Logger.Error("Client Log: Rate limiter error", "error", err)
		return &models.APIError{Message: err.Error(), Code: 429}
	}

	// Make the request
	resp, err = request.Execute(method, c.baseURL+endpoint)
	duration := time.Since(startTime)

	statusCode := 500
	if resp != nil {
		statusCode = resp.StatusCode()

		// Extract rate limit information from headers if available
		if remaining := resp.Header().Get("x-ratelimit-remaining"); remaining != "" {
			if rem, parseErr := strconv.ParseInt(remaining, 10, 64); parseErr == nil {
				rateLimit = &RateLimitResponse{
					Remaining: rem,
				}
			}
		}
		if reset := resp.Header().Get("x-ratelimit-reset"); reset != "" {
			if resetTime, parseErr := time.Parse(time.RFC3339, reset); parseErr == nil {
				if rateLimit == nil {
					rateLimit = &RateLimitResponse{}
				}
				rateLimit.Reset = resetTime
			}
		}
	}

	// Record metrics with rate limit information
	c.recordMetrics(method, endpoint, duration, statusCode, err, rateLimit)

	// If successful, return immediately
	if err == nil && !resp.IsError() {
		return nil
	}

	// Handle rate limit response
	if resp != nil && resp.StatusCode() == 429 {
		// Parse the rate limit error details
		apiError = parseAPIError(resp)
		if apiError != nil && apiError.Data != nil {
			// Update rate limit information from error response
			if limitPerSecond, ok := apiError.Data["limitPerSecond"].(float64); ok {
				if limitBurst, ok := apiError.Data["limitBurst"].(float64); ok {
					if rateLimit == nil {
						rateLimit = &RateLimitResponse{}
					}
					rateLimit.LimitPerSecond = limitPerSecond
					rateLimit.LimitBurst = int(limitBurst)

					if remaining, ok := apiError.Data["remaining"].(float64); ok {
						rateLimit.Remaining = int64(remaining)
					}

					if resetTimeStr, ok := apiError.Data["reset"].(string); ok {
						if resetTime, parseErr := time.Parse(time.RFC3339, resetTimeStr); parseErr == nil {
							rateLimit.Reset = resetTime
						}
					}

					c.Logger.Debug("Updating rate limits from API response",
						"limitPerSecond", limitPerSecond,
						"limitBurst", int(limitBurst),
						"remaining", rateLimit.Remaining,
						"reset", rateLimit.Reset)

					// Update our rate limiter with the new information
					c.RateLimiter.updateLimits(
						rateLimit.LimitPerSecond,
						rateLimit.LimitBurst,
						rateLimit.Remaining,
						rateLimit.Reset,
					)

					// Wait for the specified retry time
					if retryAfter, ok := apiError.Data["retryAfter"].(float64); ok {
						waitDuration := time.Duration(retryAfter * float64(time.Millisecond))
						c.Logger.Debug("Rate limit exceeded, waiting as specified by API",
							"wait_duration", waitDuration)
						time.Sleep(waitDuration)

						// Retry the request once after waiting
						resp, err = request.Execute(method, c.baseURL+endpoint)
						if err == nil && !resp.IsError() {
							return nil
						}
					}
				}
			}
		}
	}

	// If we still have an error, return it
	if resp.IsError() {
		apiError = parseAPIError(resp)
		c.Logger.Error("Client Log: API Request resulted in error",
			"error", apiError.Error(),
			"data", apiError.Data)
		return apiError
	}

	return nil
}

func (c *Client) recordMetrics(method, endpoint string, duration time.Duration, statusCode int, err error, rateLimit *RateLimitResponse) {
	if c.meter == nil {
		return // Telemetry is disabled
	}

	attrs := []attribute.KeyValue{
		attribute.String("agent", c.AgentSymbol),
		attribute.String("endpoint", endpoint),
		attribute.String("method", method),
		attribute.Int("status_code", statusCode),
	}

	c.requestCounter.Add(c.context, 1, metric.WithAttributes(attrs...))
	c.requestDuration.Record(c.context, duration.Seconds(), metric.WithAttributes(attrs...))

	// Record rate limit metrics if available
	if rateLimit != nil {
		c.RateLimiter.updateLimits(
			rateLimit.LimitPerSecond,
			rateLimit.LimitBurst,
			rateLimit.Remaining,
			rateLimit.Reset,
		)
	}

	// Record errors with enhanced context
	if err != nil || statusCode >= 400 {
		errorAttrs := append(attrs,
			attribute.Int("error_code", statusCode),
		)
		if err != nil {
			errorAttrs = append(errorAttrs,
				attribute.String("error_type", "client"),
				attribute.String("error_message", err.Error()),
			)
		} else {
			errorAttrs = append(errorAttrs,
				attribute.String("error_type", "server"),
			)
		}
		c.errorCounter.Add(c.context, 1, metric.WithAttributes(errorAttrs...))
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

// RateLimitResponse represents the rate limit information from the API
type RateLimitResponse struct {
	LimitPerSecond float64
	LimitBurst     int
	Remaining      int64
	Reset          time.Time
}
