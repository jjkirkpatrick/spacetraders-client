package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jjkirkpatrick/spacetraders-client/internal/cache"
	"github.com/jjkirkpatrick/spacetraders-client/internal/metrics"
	"github.com/jjkirkpatrick/spacetraders-client/models"
	"github.com/phuslu/log"
	"golang.org/x/time/rate"
)

// Client represents the SpaceTraders API client
type Client struct {
	context         context.Context
	baseURL         string
	token           string
	httpClient      *resty.Client
	retryDelay      time.Duration
	AgentSymbol     string
	MetricsReporter metrics.MetricsReporter
	CacheClient     *cache.Cache
	Logger          *log.Logger
}

// ClientOptions represents the configuration options for the SpaceTraders API client
type ClientOptions struct {
	BaseURL           string
	Symbol            string
	Faction           string
	Email             string
	RequestsPerSecond float32
	LogLevel          log.Level
	RetryDelay        time.Duration
}

// DefaultClientOptions returns the default configuration options for the SpaceTraders API client
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		BaseURL:           "https://api.spacetraders.io/v2",
		RequestsPerSecond: 2,
		RetryDelay:        1 * time.Second,
	}
}

type Glog struct {
	Logger log.Logger
}

var glog = &Glog{log.Logger{
	Level:      log.InfoLevel,
	Caller:     1,
	TimeFormat: "15:04:05.999999",
	Writer:     &log.ConsoleWriter{ColorOutput: true, Formatter: Logformat},
}}

// NewClient creates a new instance of the SpaceTraders API client
func NewClient(options ClientOptions) (*Client, error) {
	if options.Symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	client := &Client{
		baseURL:         options.BaseURL,
		httpClient:      resty.New(),
		context:         context.Background(),
		retryDelay:      options.RetryDelay,
		AgentSymbol:     options.Symbol,
		MetricsReporter: &metrics.NoOpMetricsReporter{},
		CacheClient:     cache.NewCache(),
		Logger:          &glog.Logger,
	}

	client.Logger.SetLevel(options.LogLevel)

	err := client.getOrRegisterToken(options.Faction, options.Symbol, options.Email)

	if err != nil {
		return nil, err
	}

	client.httpClient.SetRateLimiter(rate.NewLimiter(rate.Limit(options.RequestsPerSecond), 10))

	client.Logger.Debug().Msgf("New SpaceTraders client initialized with baseURL: %s, rateLimit: %f requests/second", client.baseURL, options.RequestsPerSecond)
	return client, nil
}

func (c *Client) MetricBuilder() *metrics.MetricBuilder {
	return metrics.NewMetricBuilder()
}

func (c *Client) ConfigureMetricsClient(url, token, org, bucket string) {
	c.MetricsReporter = metrics.NewMetricsClient(url, token, org, bucket)
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
	for {
		c.Logger.Trace().Msgf("Sending request: %s %s", method, c.baseURL+endpoint)
		resp, err = request.Execute(method, c.baseURL+endpoint)
		metric, _ := metrics.NewMetricBuilder().
			Namespace("api_request").
			Tag("method", method).
			Field("count", 1).
			Timestamp(time.Now()).
			Build()
		c.MetricsReporter.WritePoint(metric)
		if err != nil {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
			continue
		}

		if resp.IsError() {
			if isRateLimitError(resp.StatusCode()) {
				handleRateLimit(resp, c.Logger)
				metric, _ := metrics.NewMetricBuilder().
					Namespace("api_request_error").
					Tag("method", method).
					Tag("error_type", "rate_limit").
					Tag("status_code", fmt.Sprintf("%d", resp.StatusCode())).
					Field("count", 1).
					Timestamp(time.Now()).
					Build()
				c.MetricsReporter.WritePoint(metric)

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
					c.MetricsReporter.WritePoint(metric)

				}
				if apiError != nil {
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
		c.MetricsReporter.WritePoint(metric)
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
			logger.Debug().Msgf("Rate limit exceeded. Waiting until reset: %v", waitDuration)
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
	c.MetricsReporter.WritePoint(metric)
}

func (c *Client) GetToken() string {
	return c.token
}
