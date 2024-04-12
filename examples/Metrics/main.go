package main

import (
	"time"

	"github.com/phuslu/log"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/entities"
	"github.com/jjkirkpatrick/spacetraders-client/internal/metrics"
)

func main() {
	// Set up the logger to output to standard output with standard flags
	log.DefaultLogger = log.Logger{
		Level:      log.InfoLevel,
		Caller:     1,
		TimeFormat: "15:04:05",
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			EndWithMessage: true,
			Formatter:      client.Logformat,
		},
	}

	// Initialize default client options and assign a logger
	options := client.DefaultClientOptions()

	// Set the symbol and faction for the client
	options.Symbol = "metrics-example"
	options.Faction = "COSMIC"

	// Create a new client instance with the specified options
	client, cerr := client.NewClient(options)
	if cerr != nil {
		// If client creation fails, log the error and exit
		log.Fatal().Msgf("Failed to create client: %v", cerr)
	}

	// Configure the metrics client with InfluxDB details
	client.ConfigureMetricsClient(
		"http://influxdb:8086", // InfluxDB URL
		"238n==",               // InfluxDB Token
		"spacetraders",         // InfluxDB Organization
		"spacetraders",         // InfluxDB Bucket
	)

	// Retrieve the agent details
	agent, err := entities.GetAgent(client)

	if err != nil {
		// If retrieving the agent fails, log the error and exit
		log.Fatal().Msgf("Failed to get agent: %v", err)
	}

	// Build a metric for the agent's credits
	metric, _ := metrics.NewMetricBuilder().
		Namespace("agent_credits").       // Metric namespace
		Tag("agentSymbol", agent.Symbol). // Tag with the agent's symbol
		Field("total", agent.Credits).    // Field for the total credits
		Timestamp(time.Now()).            // Current timestamp
		Build()

	// Write the metric to the metrics reporter
	client.MetricsReporter.WritePoint(metric)

}
