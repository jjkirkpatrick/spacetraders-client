# Concurrent Agent Requests Example

This example demonstrates how to make multiple concurrent API requests using the SpaceTraders client's request queue system. It shows how to:

1. Make 100 `GetAgent` calls across 20 goroutines
2. Collect and process all responses
3. Handle any errors that might occur
4. Verify that all requests complete successfully

## Key Features Demonstrated

- **Concurrent API Requests**: Shows how to make API calls from multiple goroutines simultaneously
- **Request Queue**: Demonstrates how the client's request queue manages API rate limits automatically
- **Error Handling**: Properly captures and reports any errors that occur during API calls
- **Result Collection**: Uses channels to collect and process results from concurrent operations

## How to Run

1. Set your SpaceTraders API token as an environment variable:
   ```bash
   # On Windows PowerShell
   $env:SPACETRADERS_TOKEN="your_token_here"
   
   # On Linux/macOS
   export SPACETRADERS_TOKEN="your_token_here"
   ```

2. Edit the `main.go` file to set your agent symbol and faction:
   ```go
   options.Symbol = "YOUR_AGENT_SYMBOL" // Replace with your agent symbol
   options.Faction = "COSMIC"           // Replace with your faction
   ```

3. Run the example:
   ```bash
   go run main.go
   ```

## Expected Output

The example will:
1. Start 20 goroutines, each making 5 `GetAgent` requests
2. Show progress updates as requests complete
3. Display the first and last 5 successful responses
4. Provide a summary of all requests, including success/failure counts and timing information

## Understanding the Results

Even though the requests are initiated concurrently from 20 different goroutines, the client's request queue ensures they are processed at a controlled rate to respect the SpaceTraders API rate limits (typically 2 requests per second).

This means the total execution time will be at least:
- 100 requests ÷ 2 requests per second = 50 seconds

The actual time will be slightly longer due to processing overhead and the time it takes to start all goroutines.

This example demonstrates how the request queue effectively manages concurrent requests while ensuring API rate limits are respected, preventing rate limit errors even under high concurrency. 