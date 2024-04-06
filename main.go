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

	client, cerr := client.NewClient(options)
	if cerr != nil {
		logger.Fatalf("Failed to create client and register agent: %v", cerr)
	}

	_, _, err := client.AcceptContract("cluk89lvh36rbs60c4i01lvhe")

	if err != nil {
		logger.Fatalf("Failed to get contract details: %d %s %s", err.Code, err.Message, err.Data)
	}

}
