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

func testHello() {
	// gRPC connection
	conn, err := grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewObjectClient(conn)
	var wg sync.WaitGroup
	numRequests := 100000
	concurrency := 1000 // Adjust this based on your needs

	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Hello(context.Background(), &pb.HelloRequest{})
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

func testCreateObject() {

	const maxID = 100000
	types := []string{"car", "bike", "truck", "bus"}
	l_types := len(types)
	colors := []string{"red", "green", "blue", "yellow"}
	l_colors := len(colors)
	statuses := []string{"moving", "static"}
	l_statuses := len(statuses)

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
	numRequests := 100000
	// concurrency := 100000 // Adjust this based on your needs
	loop := 1
	requestPerLoop := numRequests / loop

	startTime := time.Now()
	for t := 0; t < loop; t++ {
		for i := 0; i < requestPerLoop; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				objectReq := &pb.ObjectRequest{
					Id:        fmt.Sprintf("%d", rand.Intn(maxID)+1),
					Type:      types[rand.Intn(l_types)],
					Color:     colors[rand.Intn(l_colors)],
					Lat:       lat,
					Lng:       lng,
					Status:    statuses[rand.Intn(l_statuses)],
					Timestamp: time.Now().Unix(),
				}
				// Increment lat/lng slightly for each object
				lat += rand.Float32() * 0.01
				lng += rand.Float32() * 0.01
				_, err := client.CreateObject(context.Background(), objectReq)
				if err != nil {
					log.Printf("Could not create object: %v", err)
				}
			}()

			// if i > 0 && i%concurrency == 0 {
			// 	wg.Wait() // Wait for a batch of goroutines to finish
			// }
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

func main() {
	// testHello()
	testCreateObject()
}
