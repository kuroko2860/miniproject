package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

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

var rdb *redis.Client
var jobQueue chan Object
var numWorkers = 10000

func worker(jobs <-chan Object) {
	for data := range jobs {
		// Store data in Redis
		err := storeObjectInRedis(data)
		if err != nil {
			fmt.Println("Error storing in Redis:", err)
		}
	}
}
func main() {
	// Initialize Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Address of the Redis container within Docker network
		Password: "",           // No password set
		DB:       0,            // Use default DB
	})

	// Create worker pool and job queue
	jobQueue = make(chan Object, 10000) // Buffered channel to prevent blocking
	for i := 0; i < numWorkers; i++ {
		go worker(jobQueue)
	}

	// Create a new Gin router
	router := gin.Default()

	// Define route handler for POST requests to "/object"
	router.POST("/objects", addObject)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("Go Server is running on :%s", port)
	log.Fatal(server.ListenAndServe())
}

// addObject handles POST requests to add an object
func addObject(c *gin.Context) {
	// Parse JSON request body into Object struct
	var obj Object
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(c.Request.Body)
		return
	}

	// Store object data in Redis asynchronously
	jobQueue <- obj

	// Respond immediately
	c.JSON(http.StatusOK, gin.H{"message": "Object data is being processed"})
}

// storeObjectInRedis pushes object data to a Redis list
func storeObjectInRedis(obj Object) error {
	// Marshal object data to JSON
	objJSON, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Push object data to Redis list
	_, err = rdb.RPush(context.Background(), "objects", objJSON).Result()
	if err != nil {
		return err
	}

	return nil
}
