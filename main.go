package main

import (
	"log"
	"os"

	"github.com/jjkirkpatrick/spacetraders-client/client"
)

func main() {
	// Set up the logger
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Create a new client without a token (register a new agent)
	options := client.DefaultClientOptions()
	options.Logger = logger
	options.Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGlmaWVyIjoiTVlBR0VOVDQiLCJ2ZXJzaW9uIjoidjIuMi4wIiwicmVzZXRfZGF0ZSI6IjIwMjQtMDMtMjQiLCJpYXQiOjE3MTIxNzQxMjUsInN1YiI6ImFnZW50LXRva2VuIn0.b4ICe3ILHtC8P2g4ycOtTbsXO1x8obBAFX6f7TNtltoJ7cImvUYHgEIC8yc8REbXtiu3JTgcbYAuqgP6qoEOzKJlIn_vHuhOZ3AX6OjOQdi6hnlqYq0kF0Vn36CV15Pp8ulgObKGx9zB1SzedLUV5ud77bGZNUNQv5MW8VKGBpqwN3Kv_Eh9dzdyIuKXvD6hMTgC7FlAbVUJE7itThOAnHvX7BBzJY6aiRGdSCbuh07YDRyQ-_28JB4cFC1byXLJ50ZC-3Oh7zdnQFEYXPX3Akv9ntRSQGugXdCiHmDyNdoB_29I6-fYwF1kGw3tXNUvwf6QnndTme5zTREQ0XYZPg"

	client, err := client.NewClient(options)
	if err != nil {
		logger.Fatalf("Failed to create client and register agent: %v", err)
	}

	contracts, err := client.ListContracts()

	if err != nil {
		logger.Fatalf("Failed to list contracts: %v", err)
	}

	for _, contract := range contracts.Data {
		logger.Printf("Contract ID: %s\n", contract.ID)
	}

}
