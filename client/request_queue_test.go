package client

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/models"
	"github.com/stretchr/testify/assert"
)

// mockExecutor is a mock implementation of the RequestExecutor interface for testing
type mockExecutor struct {
	executeRequestFunc func(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError
}

func (m *mockExecutor) executeRequest(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
	return m.executeRequestFunc(ctx, method, endpoint, body, queryParams, result)
}

func TestRequestQueue_Enqueue(t *testing.T) {
	// Create a mock executor
	mockExec := &mockExecutor{
		executeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
			// Simulate a successful request
			return nil
		},
	}

	// Create a request queue with the mock executor
	ctx := context.Background()
	queue := NewRequestQueue(ctx, mockExec, 10)
	defer queue.Shutdown()

	// Test a single request
	var result interface{}
	err := queue.Enqueue("GET", "/test", nil, nil, &result)
	assert.Nil(t, err)
}

func TestRequestQueue_ConcurrentRequests(t *testing.T) {
	// Create a mock executor with a delay to simulate API processing
	mockExec := &mockExecutor{
		executeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
			// Simulate processing time
			time.Sleep(100 * time.Millisecond)
			return nil
		},
	}

	// Create a request queue with the mock executor
	ctx := context.Background()
	queue := NewRequestQueue(ctx, mockExec, 20)
	defer queue.Shutdown()

	// Test multiple concurrent requests
	var wg sync.WaitGroup
	requestCount := 10
	wg.Add(requestCount)

	startTime := time.Now()

	for i := 0; i < requestCount; i++ {
		go func(i int) {
			defer wg.Done()
			var result interface{}
			err := queue.Enqueue("GET", "/test", nil, nil, &result)
			assert.Nil(t, err)
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)

	// All requests should be processed at a rate of 2 per second (500ms each)
	// So 10 requests should take at least 5 seconds
	// But we'll be a bit lenient in the test
	assert.True(t, duration >= 4*time.Second, "Requests should be rate-limited")
}

func TestRequestQueue_Shutdown(t *testing.T) {
	// Create a mock executor
	mockExec := &mockExecutor{
		executeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
			// Simulate a successful request with some delay
			time.Sleep(50 * time.Millisecond)
			return nil
		},
	}

	// Create a request queue with the mock executor
	ctx := context.Background()
	queue := NewRequestQueue(ctx, mockExec, 10)

	// Enqueue a few requests
	for i := 0; i < 5; i++ {
		go func() {
			var result interface{}
			_ = queue.Enqueue("GET", "/test", nil, nil, &result)
		}()
	}

	// Give some time for requests to be enqueued
	time.Sleep(100 * time.Millisecond)

	// Shutdown the queue
	queue.Shutdown()

	// Try to enqueue after shutdown
	var result interface{}
	err := queue.Enqueue("GET", "/test", nil, nil, &result)

	// Should return an error
	assert.NotNil(t, err)
	assert.Equal(t, 499, err.Code)
}

func TestRequestQueue_RateLimitHandling(t *testing.T) {
	// Create a counter for requests
	var requestCount int
	var mu sync.Mutex

	// Create a mock executor that simulates rate limit errors
	mockExec := &mockExecutor{
		executeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
			mu.Lock()
			requestCount++
			count := requestCount
			mu.Unlock()

			// Every third request will hit a rate limit
			if count%3 == 0 {
				return &models.APIError{
					Code:    429,
					Message: "Rate limit exceeded",
					Data: map[string]interface{}{
						"limitPerSecond": 2.0,
						"limitBurst":     30.0,
						"remaining":      0.0,
						"reset":          time.Now().Add(1 * time.Second).Format(time.RFC3339),
						"retryAfter":     1000.0,
					},
				}
			}
			return nil
		},
	}

	// Create a request queue with the mock executor
	ctx := context.Background()
	queue := NewRequestQueue(ctx, mockExec, 10)
	defer queue.Shutdown()

	// Test handling of rate limit errors
	var wg sync.WaitGroup
	successCount := 0
	var successMu sync.Mutex

	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var result interface{}
			err := queue.Enqueue("GET", "/test", nil, nil, &result)
			if err == nil {
				successMu.Lock()
				successCount++
				successMu.Unlock()
			}
		}()
	}

	wg.Wait()

	// We should have some successful requests despite rate limiting
	assert.True(t, successCount > 0)
}

func TestRequestQueue_EnqueueWithContext(t *testing.T) {
	mockExec := &mockExecutor{
		executeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
			return nil
		},
	}

	ctx := context.Background()
	queue := NewRequestQueue(ctx, mockExec, 10)
	defer queue.Shutdown()

	// Create context with labels
	reqCtx := WithMetricLabels(context.Background(), map[string]string{
		"tree_name":   "mining",
		"action_name": "extract",
	})

	var result interface{}
	err := queue.EnqueueWithContext(reqCtx, "GET", "/test", nil, nil, &result)
	assert.Nil(t, err)
}

func TestRequestQueue_ContextPropagation(t *testing.T) {
	var capturedLabels map[string]string

	mockExec := &mockExecutor{
		executeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}, queryParams map[string]string, result interface{}) *models.APIError {
			capturedLabels = GetMetricLabels(ctx)
			return nil
		},
	}

	ctx := context.Background()
	queue := NewRequestQueue(ctx, mockExec, 10)
	defer queue.Shutdown()

	reqCtx := WithMetricLabels(context.Background(), map[string]string{
		"tree_name":   "trader",
		"action_name": "sell_cargo",
		"ship_role":   "hauler",
	})
	var result interface{}
	err := queue.EnqueueWithContext(reqCtx, "GET", "/test", nil, nil, &result)

	assert.Nil(t, err)
	assert.Equal(t, "trader", capturedLabels["tree_name"])
	assert.Equal(t, "sell_cargo", capturedLabels["action_name"])
	assert.Equal(t, "hauler", capturedLabels["ship_role"])
}
