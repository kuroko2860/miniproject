package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "kuroko.com/goserver/object"
)

func main() {

	const maxID = 100000
	types := []string{"car", "bike", "truck", "bus"}
	colors := []string{"red", "green", "blue", "yellow"}
	statuses := []string{"moving", "static"}

	// Starting coordinates
	lat := rand.Float32()*180 - 90
	lng := rand.Float32()*360 - 180

	// gRPC connection
	conn, err := grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewObjectClient(conn)
	var wg sync.WaitGroup
	numRequests := 10000
	concurrency := 100 // Adjust this based on your needs

	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// start := time.Now()

			objectReq := &pb.ObjectRequest{
				// ... populate object request ...
				Id:     fmt.Sprintf("%d", rand.Intn(maxID)+1),
				Type:   types[rand.Intn(len(types))],
				Color:  colors[rand.Intn(len(colors))],
				Lat:    lat,
				Lng:    lng,
				Status: statuses[rand.Intn(len(statuses))],
			}
			// Increment lat/lng slightly for each object
			lat += rand.Float32() * 0.001
			lng += rand.Float32() * 0.001
			_, err := client.CreateObject(context.Background(), objectReq)
			if err != nil {
				log.Printf("Could not create object: %v", err)
			}

			// elapsed := time.Since(start)
			// fmt.Printf("Request %d took %s\n", i+1, elapsed)
		}()

		if i > 0 && i%concurrency == 0 {
			wg.Wait() // Wait for a batch of goroutines to finish
		}
	}
	wg.Wait() // Wait for all remaining goroutines

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	rps := float64(numRequests) / totalTime.Seconds()
	avgLatency := totalTime.Seconds() / float64(numRequests)

	fmt.Printf("\nTotal requests: %d\n", numRequests)
	fmt.Printf("Total time: %s\n", totalTime)
	fmt.Printf("Requests per second (RPS): %.2f\n", rps)
	fmt.Printf("Average latency: %.2f ms\n", avgLatency*1000)

}
