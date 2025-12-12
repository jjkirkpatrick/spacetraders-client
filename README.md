# SpaceTraders API Client

A Go client library for the [SpaceTraders API](https://api.spacetraders.io/v2). Provides rate-limited request handling, automatic pagination, OpenTelemetry observability (metrics, traces, logs), and game reset detection.

> **Warning**: This is a work in progress. Not all endpoints have been tested. There will be bugs and missing features, and the API may change over time.

## Installation

```bash
go get github.com/jjkirkpatrick/spacetraders-client
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"

    "github.com/jjkirkpatrick/spacetraders-client/client"
    "github.com/jjkirkpatrick/spacetraders-client/entities"
)

func main() {
    ctx := context.Background()

    // Create client with default options
    options := client.DefaultClientOptions()
    options.Symbol = "YOUR-AGENT-SYMBOL"
    options.Faction = "COSMIC"

    c, err := client.NewClient(options)
    if err != nil {
        slog.Error("Failed to create client", "error", err)
        os.Exit(1)
    }
    defer c.Close(ctx)

    // Get agent information
    agent, err := entities.GetAgent(c)
    if err != nil {
        slog.Error("Failed to get agent", "error", err)
        os.Exit(1)
    }

    fmt.Printf("Agent: %s, Credits: %d\n", agent.Symbol, agent.Credits)
}
```

## Client Configuration

```go
options := client.DefaultClientOptions()
options.Symbol = "YOUR-AGENT-SYMBOL"    // Required: Your agent symbol
options.Faction = "COSMIC"               // Required for new agents
options.LogLevel = slog.LevelInfo        // Optional: Log level (default: Info)
options.RequestQueueSize = 200           // Optional: Queue size (default: 100)

c, err := client.NewClient(options)
```

## OpenTelemetry Integration

The client supports full OpenTelemetry observability including metrics, traces, and logs. This allows you to monitor your application using Grafana, Prometheus, Jaeger, Loki, or any OTLP-compatible backend.

### Enabling Telemetry

```go
options := client.DefaultClientOptions()
options.Symbol = "YOUR-AGENT-SYMBOL"
options.Faction = "COSMIC"

// Enable OpenTelemetry
options.TelemetryOptions = client.DefaultTelemetryOptions()
options.TelemetryOptions.ServiceName = "my-spacetraders-app"
options.TelemetryOptions.ServiceVersion = "1.0.0"
options.TelemetryOptions.OTLPEndpoint = "localhost:4317"  // Your OTLP collector
options.TelemetryOptions.Environment = "development"

// Optional: Add custom attributes to all telemetry
options.TelemetryOptions.AdditionalAttributes = map[string]string{
    "deployment": "us-west",
    "team": "platform",
}

c, err := client.NewClient(options)
```

### Setting Up Logging with Loki

The client provides a public `telemetry` package with slog handlers that send logs to both console and OTLP (for Loki):

```go
import (
    "log/slog"
    "os"

    "github.com/jjkirkpatrick/spacetraders-client/client"
    "github.com/jjkirkpatrick/spacetraders-client/telemetry"
)

func main() {
    ctx := context.Background()

    // Create client with telemetry enabled
    options := client.DefaultClientOptions()
    options.Symbol = "YOUR-AGENT-SYMBOL"
    options.Faction = "COSMIC"
    options.TelemetryOptions = client.DefaultTelemetryOptions()
    options.TelemetryOptions.ServiceName = "my-app"
    options.TelemetryOptions.OTLPEndpoint = "localhost:4317"

    c, err := client.NewClient(options)
    if err != nil {
        slog.Error("Failed to create client", "error", err)
        os.Exit(1)
    }
    defer c.Close(ctx)

    // Set up combined logging (console + OTLP/Loki)
    consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
    combinedHandler := telemetry.NewCombinedSlogHandler("my-app", slog.LevelInfo, consoleHandler)
    slog.SetDefault(slog.New(combinedHandler))

    // Logs now go to both console AND Loki via OTLP
    slog.Info("Application started", "agent", options.Symbol)
}
```

### Creating Traces

```go
import "go.opentelemetry.io/otel"

// Get a tracer
tracer := otel.GetTracerProvider().Tracer("my-app")

// Create a span
ctx, span := tracer.Start(ctx, "my_operation")
defer span.End()

// Add attributes to the span
span.SetAttributes(
    attribute.String("agent", agent.Symbol),
    attribute.Int64("credits", agent.Credits),
)

// Log with trace context (trace_id and span_id auto-injected)
slog.InfoContext(ctx, "Operation completed", "result", "success")
```

### Available Metrics

The client automatically exports the following metrics:

| Metric | Type | Description |
|--------|------|-------------|
| `api_requests_total` | Counter | Total API requests made |
| `api_request_duration_seconds` | Histogram | Request duration |
| `api_errors_total` | Counter | Total API errors |
| `api_retries_total` | Counter | Total request retries |
| `api_rate_limit` | Gauge | Current rate limit |
| `api_remaining_requests` | Gauge | Requests remaining before rate limit |
| `api_queue_length` | Gauge | Requests waiting in queue |
| `api_queue_wait_time_seconds` | Histogram | Time spent waiting in queue |
| `api_queue_process_time_seconds` | Histogram | Time to process requests |

## Rate Limiting and Request Queue

The SpaceTraders API enforces rate limits (2 requests per second with burst capability). The client automatically handles this through a centralized request queue.

### How It Works

1. **Automatic Rate Limiting**: All requests are queued and processed at a controlled rate
2. **Concurrent-Safe**: Multiple goroutines can safely make API calls
3. **Automatic Retries**: Rate limit errors (429) are automatically retried with exponential backoff
4. **Adaptive Rate**: The queue adjusts its processing rate based on API responses

### Example: Concurrent Requests

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()

        // Requests are automatically queued and rate-limited
        agent, err := entities.GetAgent(c)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }

        fmt.Printf("Agent %s has %d credits\n", agent.Symbol, agent.Credits)
    }()
}

wg.Wait()
```

## Game Reset Handling

The SpaceTraders game undergoes periodic resets which invalidate existing tokens. The client automatically detects these resets.

### Detection Methods

```go
// Non-blocking check
if c.IsGameReset() {
    // Handle game reset
}

// Blocking wait with context
if c.WaitForGameReset(ctx) {
    // Game reset detected
}
```

### Example: Monitoring for Resets

```go
// Start a goroutine to monitor for game resets
go func() {
    if c.WaitForGameReset(ctx) {
        slog.Error("Game reset detected! Re-registration required.")
        os.Exit(1)
    }
}()

// In your main loop
for {
    if c.IsGameReset() {
        break
    }

    // Make API requests...
}
```

## API Documentation

For detailed information on specific operations, see the guides in the `Docs/` folder:

- [Agent Operations](Docs/Agent.md) - Agent information and public agent queries
- [Ship Operations](Docs/Ships.md) - Ship management, navigation, mining, trading
- [System Operations](Docs/Systems.md) - Systems, waypoints, markets, shipyards
- [Contract Operations](Docs/Contracts.md) - Contract management and fulfillment
- [Faction Operations](Docs/Factions.md) - Faction information and listings

## Examples

The `examples/` directory contains working examples:

- `quick_start/` - Complete mining bot with contracts
- `otel_minimal/` - Minimal OpenTelemetry setup
- `concurrent/` - Concurrent request handling
- `ships/` - Ship operations
- `caching/` - System caching
- `game_reset_handling/` - Game reset detection

## Docker Setup for Observability

To run a local observability stack (Grafana, Tempo, Loki, Prometheus), use this docker-compose:

```yaml
version: '3.8'

services:
  # OpenTelemetry Collector
  otel-collector:
    image: otel/opentelemetry-collector-contrib:latest
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4317:4317"   # OTLP gRPC
      - "4318:4318"   # OTLP HTTP

  # Grafana
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana

  # Tempo (Traces)
  tempo:
    image: grafana/tempo:latest
    command: ["-config.file=/etc/tempo.yaml"]
    volumes:
      - ./tempo.yaml:/etc/tempo.yaml
    ports:
      - "3200:3200"

  # Loki (Logs)
  loki:
    image: grafana/loki:latest
    ports:
      - "3100:3100"

  # Prometheus (Metrics)
  prometheus:
    image: prom/prometheus:latest
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"

volumes:
  grafana-data:
```

See the `examples/` directory for complete configuration files.

## Graceful Shutdown

Always close the client to ensure telemetry data is flushed:

```go
defer c.Close(context.Background())
```
