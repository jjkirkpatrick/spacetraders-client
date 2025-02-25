# Game Reset Handling Example

This example demonstrates how to handle game resets in the SpaceTraders API. The SpaceTraders game undergoes periodic resets (typically weekly or bi-weekly during alpha), which invalidate existing tokens. When this happens, you need to re-register your agent.

## Key Features Demonstrated

- **Game Reset Detection**: Shows how to detect when the SpaceTraders game has been reset
- **Non-blocking Monitoring**: Uses a dedicated goroutine to monitor for game resets without blocking the main application
- **Graceful Shutdown**: Properly handles application shutdown when a game reset is detected
- **Multiple Detection Methods**: Demonstrates both blocking and non-blocking ways to check for game resets

## How It Works

The SpaceTraders client now includes a `GameResetCh` channel that receives a notification when a token version mismatch error is detected. This error indicates that the game has been reset and the token is no longer valid.

The example demonstrates two ways to check for game resets:

1. **Blocking Method** (`WaitForGameReset`): Waits until either a game reset is detected or the context is cancelled
2. **Non-blocking Method** (`IsGameReset`): Checks if a game reset has been detected without blocking

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
1. Start making API requests in a loop
2. Monitor for game resets in a separate goroutine
3. If a game reset is detected, it will print a message and exit the application

## Integrating Game Reset Handling in Your Application

To handle game resets in your own application:

1. Check for game resets using either the blocking or non-blocking method:
   ```go
   // Non-blocking check
   if client.IsGameReset() {
       // Handle game reset (e.g., re-register agent, exit application)
   }
   
   // Or, blocking wait with a timeout
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   if client.WaitForGameReset(ctx) {
       // Handle game reset
   }
   ```

2. When a game reset is detected, you should:
   - Stop making API requests
   - Notify the user that the game has been reset
   - Either exit the application or re-register the agent with a new token

This approach ensures your application can gracefully handle game resets without continuously encountering errors. 