package main

import (
	"log"
	"os"
	"time"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/metrics"
)

func main() {
	// Set up the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create a new client without a token (register a new agent)
	options := client.DefaultClientOptions()
	options.Logger = logger
	options.Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiTVlBR0VOVDQiLCJ2ZXJzaW9uIjoidjIuMi4wIiwicmVzZXRfZGF0ZSI6IjIwMjQtMDMtMjQiLCJpYXQiOjE3MTIxNzQxMjUsInN1YiI6ImFnZW50LXRva2VuIn0.b4ICe3ILHtC8P2g4ycOtTbsXO1x8obBAFX6f7TNtltoJ7cImvUYHgEIC8yc8REbXtiu3JTgcbYAuqgP6qoEOzKJlIn_vHuhOZ3AX6OjOQdi6hnlqYq0kF0Vn36CV15Pp8ulgObKGx9zB1SzedLUV5ud77bGZNUNQv5MW8VKGBpqwN3Kv_Eh9dzdyIuKXvD6hMTgC7FlAbVUJE7itThOAnHvX7BBzJY6aiRGdSCbuh07YDRyQ-_28JB4cFC1byXLJ50ZC-3Oh7zdnQFEYXPX3Akv9ntRSQGugXdCiHmDyNdoB_29I6-fYwF1kGw3tXNUvwf6QnndTme5zTREQ0XYZPg"

	metricsReporter := metrics.NewMetricsClient("http://192.168.1.33:8086", "238nUuJVX9CzDqdsU7wvINk8ByIG-3MykZqwUSTEwBIeLgBKNTbwV8x_lik_4t1oXSTfj8OqRPwzmPS8y3tsdg==", "spacetraders", "spacetraders")

	client, cerr := client.NewClient(options, metricsReporter)
	if cerr != nil {
		logger.Fatalf("Failed to create client and register agent: %v", cerr)
	}

	agent, err := client.GetAgent()

	if err != nil {
		logger.Fatalf("Failed to get agent: %v", err)
	}

	metric, merr := metrics.NewMetricBuilder().
		Namespace("spacetraders").
		Tag("agent", agent.Symbol).
		Field("credits", agent.Credits).
		Timestamp(time.Now()).
		Build()

	if merr != nil {
		logger.Fatalf("Failed to build metric: %v", err)
	}

	client.WriteMetric(metric)

	paginator, err := client.ListSystems()
	if err != nil {
		logger.Fatalf("Failed to initiate listing systems: %v", err)
	}

	allPages, err := paginator.FetchAllPages()
	if err != nil {
		logger.Fatalf("Failed to fetch all pages: %v", err)
	}

	for _, contract := range allPages {
		logger.Printf("Contract: %+v", contract.Symbol)
	}

	logger.Printf("Total number of items: %d", len(allPages))

}
