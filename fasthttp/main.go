package main

import (
	"context"
	// "encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
)

// Object represents the structure of the object datatype Object struct {
type Object struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Color    string `json:"color"`
	Location struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"location"`
	Status string `json:"status"`
	Dump   string `json:"dump"` // Dump data to get 10kb request data
}

var (
	redisClient *redis.Client
)

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Adjust the address if needed
	})

	// Set up fasthttp server
	server := &fasthttp.Server{
		Handler:      requestHandler,
		Name:         "FastHTTP Server",
		TCPKeepalive: true,
		Concurrency:  10000,
	}

	// Start the server
	log.Println("change Starting server on :8080...")
	if err := server.ListenAndServe(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		getHealthCheck(ctx)
	case "/objects":
		postObject(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func getHealthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("Hello World")
}

func postObject(ctx *fasthttp.RequestCtx) {
	// Perform the Redis operation in a separate Goroutine
	go func(data []byte) {
		if err := redisClient.RPush(context.Background(), "objects", data).Err(); err != nil {
			log.Printf("Failed to store object in Redis: %v", err)
		}
	}(ctx.PostBody())

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("OK")
}
