package client

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// RequestExecutor is an interface for executing API requests
// This allows us to mock the client in tests
type RequestExecutor interface {
	executeRequest(method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError
}

// apiRequest represents a request to be processed by the queue
type apiRequest struct {
	method      string
	endpoint    string
	body        interface{}
	queryParams map[string]string
	result      interface{}
	responseCh  chan apiResponse
	// Timestamps for metrics
	enqueuedAt time.Time
	startedAt  time.Time
	finishedAt time.Time
}

// apiResponse represents the response from a processed request
type apiResponse struct {
	err         *models.APIError
	queueTime   time.Duration // Time spent in queue
	processTime time.Duration // Time spent processing
}

// RequestQueue manages a queue of API requests to be processed at a controlled rate
type RequestQueue struct {
	requests     chan apiRequest
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	executor     RequestExecutor
	processingCh chan struct{} // Channel to control processing rate

	// Metrics tracking
	mu                sync.RWMutex
	totalQueueTime    time.Duration
	totalProcessTime  time.Duration
	requestsProcessed int64
}

// Maximum number of retries for rate-limited requests
const maxRetries = 3

// NewRequestQueue creates a new request queue with the specified buffer size
func NewRequestQueue(ctx context.Context, executor RequestExecutor, bufferSize int) *RequestQueue {
	queueCtx, cancel := context.WithCancel(ctx)

	queue := &RequestQueue{
		requests:     make(chan apiRequest, bufferSize),
		ctx:          queueCtx,
		cancel:       cancel,
		executor:     executor,
		processingCh: make(chan struct{}, 1), // Buffer of 1 to allow non-blocking sends
	}

	// Start the worker goroutine
	queue.startWorker()

	return queue
}

// startWorker starts the worker goroutine that processes requests
func (q *RequestQueue) startWorker() {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		q.processRequests()
	}()
}

// processRequests processes requests from the queue at a controlled rate
func (q *RequestQueue) processRequests() {
	// Create a ticker to control the rate of processing
	// Default to 450ms per request (slightly faster than 2 requests per second)
	baseTickerInterval := 450 * time.Millisecond
	currentTickerInterval := baseTickerInterval
	ticker := time.NewTicker(currentTickerInterval)
	defer ticker.Stop()

	// Track consecutive rate limit errors to adjust processing rate
	consecutiveRateLimitErrors := 0
	consecutiveSuccesses := 0

	for {
		select {
		case <-q.ctx.Done():
			// Context cancelled, stop processing
			return
		case <-ticker.C:
			// Time to process a request if available
			select {
			case req := <-q.requests:
				// Record when processing started
				req.startedAt = time.Now()
				queueTime := req.startedAt.Sub(req.enqueuedAt)

				// Process the request with retries for rate limit errors
				var err *models.APIError
				var processTime time.Duration

				// Try the request with retries
				for retryCount := 0; retryCount <= maxRetries; retryCount++ {
					// Execute the request
					err = q.executor.executeRequest(req.method, req.endpoint, req.body, req.queryParams, req.result)

					// If successful or not a rate limit error, break out of retry loop
					if err == nil || err.Code != 429 {
						break
					}

					// This is a rate limit error
					consecutiveRateLimitErrors++
					consecutiveSuccesses = 0

					// Record retry metric if client has telemetry enabled
					if client, ok := q.executor.(*Client); ok && client.meter != nil {
						client.retryCounter.Add(client.context, 1, metric.WithAttributes(
							attribute.String("agent", client.AgentSymbol),
							attribute.String("endpoint", req.endpoint),
							attribute.String("method", req.method),
							attribute.Int("retry_count", retryCount),
						))
					}

					// Adjust ticker interval if we're getting too many rate limit errors
					if consecutiveRateLimitErrors >= 2 {
						// Increase the interval by 20% each time, up to 2x the base interval
						newInterval := currentTickerInterval * 6 / 5
						if newInterval > baseTickerInterval*2 {
							newInterval = baseTickerInterval * 2
						}

						if newInterval != currentTickerInterval {
							currentTickerInterval = newInterval
							ticker.Reset(currentTickerInterval)

							// Log the adjustment
							if client, ok := q.executor.(*Client); ok {
								client.Logger.Info("Adjusting request processing rate due to rate limits",
									"new_interval", currentTickerInterval.String(),
									"consecutive_errors", consecutiveRateLimitErrors)
							}
						}
					}

					// This is a rate limit error, prepare to retry
					if retryCount < maxRetries {
						// Calculate backoff time - start with 500ms and increase exponentially
						// Also use the retryAfter value from the API if available
						backoff := time.Duration(500*time.Millisecond) * time.Duration(1<<retryCount) // 500ms, 1s, 2s, etc.

						// Check if the API provided a retryAfter value
						if err.Data != nil {
							if retryAfter, ok := err.Data["retryAfter"].(float64); ok && retryAfter > 0 {
								// Convert to duration (API returns milliseconds)
								apiBackoff := time.Duration(retryAfter * float64(time.Millisecond))
								// Use the API's suggestion if it's reasonable
								if apiBackoff < 5*time.Second {
									backoff = apiBackoff
								}
							}

							// If we have reset information, use that for a more accurate backoff
							if resetStr, ok := err.Data["reset"].(string); ok {
								if resetTime, parseErr := time.Parse(time.RFC3339, resetStr); parseErr == nil {
									resetBackoff := time.Until(resetTime) + 50*time.Millisecond
									if resetBackoff > 0 && resetBackoff < backoff {
										backoff = resetBackoff
									}
								}
							}
						}

						// Add a small jitter to prevent thundering herd
						jitter := time.Duration(rand.Int63n(int64(50 * time.Millisecond)))
						backoff += jitter

						// Log the retry
						if client, ok := q.executor.(*Client); ok {
							client.Logger.Info("Rate limit exceeded, retrying request",
								"endpoint", req.endpoint,
								"method", req.method,
								"retry", retryCount+1,
								"backoff", backoff.String())
						}

						// Wait before retrying
						select {
						case <-q.ctx.Done():
							// Context cancelled during backoff, stop processing
							err = &models.APIError{
								Code:    499, // Client closed request
								Message: "request cancelled during retry backoff: client is shutting down",
							}
							break
						case <-time.After(backoff):
							// Continue to retry
						}
					}
				}

				// If we didn't get a rate limit error this time, track consecutive successes
				if err == nil || err.Code != 429 {
					consecutiveRateLimitErrors = 0
					consecutiveSuccesses++

					// If we've had multiple successful requests, gradually decrease the ticker interval
					if consecutiveSuccesses >= 5 && currentTickerInterval > baseTickerInterval {
						// Decrease by 5% each time
						newInterval := currentTickerInterval * 95 / 100
						if newInterval < baseTickerInterval {
							newInterval = baseTickerInterval
						}

						if newInterval != currentTickerInterval {
							currentTickerInterval = newInterval
							ticker.Reset(currentTickerInterval)

							// Log the adjustment
							if client, ok := q.executor.(*Client); ok {
								client.Logger.Info("Adjusting request processing rate after successful requests",
									"new_interval", currentTickerInterval.String(),
									"consecutive_successes", consecutiveSuccesses)
							}
						}
					}
				}

				// Record when processing finished
				req.finishedAt = time.Now()
				processTime = req.finishedAt.Sub(req.startedAt)

				// Update metrics
				q.mu.Lock()
				q.totalQueueTime += queueTime
				q.totalProcessTime += processTime
				q.requestsProcessed++
				q.mu.Unlock()

				// Send the response back to the caller
				req.responseCh <- apiResponse{
					err:         err,
					queueTime:   queueTime,
					processTime: processTime,
				}

				// Signal that processing is complete
				select {
				case q.processingCh <- struct{}{}:
				default:
					// Non-blocking send
				}
			default:
				// No request available, continue
			}
		}
	}
}

// Enqueue adds a request to the queue and returns a channel for the response
func (q *RequestQueue) Enqueue(method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
	// Create a response channel
	responseCh := make(chan apiResponse, 1)

	// Create the request with current timestamp
	req := apiRequest{
		method:      method,
		endpoint:    endpoint,
		body:        body,
		queryParams: queryParams,
		result:      result,
		responseCh:  responseCh,
		enqueuedAt:  time.Now(),
	}

	// Add the request to the queue
	select {
	case q.requests <- req:
		// Request added to queue
	case <-q.ctx.Done():
		// Context cancelled, return error
		return &models.APIError{
			Code:    499, // Client closed request
			Message: "request cancelled: client is shutting down",
		}
	}

	// Wait for the response
	select {
	case resp := <-responseCh:
		// Record queue metrics in the client if needed
		if client, ok := q.executor.(*Client); ok && client.meter != nil {
			attrs := []attribute.KeyValue{
				attribute.String("agent", client.AgentSymbol),
				attribute.String("endpoint", endpoint),
				attribute.String("method", method),
			}
			client.queueWaitTime.Record(client.context, resp.queueTime.Seconds(), metric.WithAttributes(attrs...))
			client.queueProcessTime.Record(client.context, resp.processTime.Seconds(), metric.WithAttributes(attrs...))
		}
		return resp.err
	case <-q.ctx.Done():
		return &models.APIError{
			Code:    499, // Client closed request
			Message: "request cancelled: client is shutting down",
		}
	}
}

// Shutdown gracefully shuts down the request queue
func (q *RequestQueue) Shutdown() {
	// Signal the worker to stop
	q.cancel()

	// Wait for the worker to finish
	q.wg.Wait()
}

// QueueLength returns the current number of requests in the queue
func (q *RequestQueue) QueueLength() int {
	return len(q.requests)
}

// GetMetrics returns the current queue metrics
func (q *RequestQueue) GetMetrics() (avgQueueTime, avgProcessTime time.Duration, requestsProcessed int64) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.requestsProcessed > 0 {
		avgQueueTime = q.totalQueueTime / time.Duration(q.requestsProcessed)
		avgProcessTime = q.totalProcessTime / time.Duration(q.requestsProcessed)
	}

	return avgQueueTime, avgProcessTime, q.requestsProcessed
}
