package main

import (
	"context"
	"encoding/json"
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
		Handler:     requestHandler,
		Concurrency: 3000,
		Name:        "FastHTTP Server",
	}

	// Start the server
	log.Println("Starting server on :8080...")
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
	var obj Object

	// Parse JSON body
	if err := json.Unmarshal(ctx.PostBody(), &obj); err != nil {
		ctx.Error("Invalid request body", fasthttp.StatusBadRequest)
		return
	}

	// Convert object to JSON
	objJSON, err := json.Marshal(obj)
	if err != nil {
		ctx.Error("Failed to marshal object", fasthttp.StatusInternalServerError)
		return
	}

	// Push object JSON to Redis list
	if err := redisClient.RPush(context.Background(), "objects", objJSON).Err(); err != nil {
		ctx.Error("Failed to store object in Redis", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("OK")
}
