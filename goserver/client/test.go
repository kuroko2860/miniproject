package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// "runtime/pprof"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "kuroko.com/goserver/object"
)

var (
	conn   *grpc.ClientConn
	err    error
	wg     sync.WaitGroup
	client pb.ObjectClient
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
	concurrency := 8000 // Adjust this based on your needs

	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.Hello(context.Background(), &pb.HelloRequest{})
		}()
		if i%concurrency == 0 {
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
	types := []string{"car", "bike", "truck", "bus"}
	l_types := len(types)
	colors := []string{"red", "green", "blue", "yellow"}
	l_colors := len(colors)
	statuses := []string{"moving", "static"}
	l_statuses := len(statuses)

	var numRequests int32 = 100000
	var concurrency int32 = 50000 // Adjust this based on your needs
	// Vietnam's approximate bounding box (latitude and longitude)
	minLat := float32(8.5)
	// maxLat := float32(15.3)
	minLng := float32(102.1)
	// maxLng := float32(109.5)

	// Starting coordinates
	lat := minLat
	lng := minLng

	startTime := time.Now()
	type_idx := 0
	color_idx := 0
	status_idx := 0
	var i int32
	for i = 1; i <= numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			objectReq := &pb.ObjectRequest{
				Id:        i%concurrency + 1,
				Type:      types[type_idx],
				Color:     colors[color_idx],
				Lat:       lat,
				Lng:       lng,
				Status:    statuses[status_idx],
				Timestamp: time.Now().Unix(),
			}
			// log.Println("id", objectReq.Id, " time:", objectReq.Timestamp)

			// Increment lat/lng slightly for each object
			lat = lat + 0.01*float32(i%10)
			lng = lng + 0.01*float32(i%10)
			_, err := client.CreateObject(context.Background(), objectReq)
			if err != nil {
				log.Printf("Could not create object: %v", err)
			}
		}()
		if i%concurrency == 0 {
			type_idx = (type_idx + 1) % l_types
			color_idx = (color_idx + 1) % l_colors
			status_idx = (status_idx + 1) % l_statuses
			wg.Wait() // Wait for a batch of goroutines to finish
		}
	}
	wg.Wait() // Wait for all remaining goroutines

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	rps := float64(numRequests) / totalTime.Seconds()
	avgLatency := totalTime.Seconds() / float64(numRequests)

	fmt.Printf("Total requests: %d\n", numRequests)
	fmt.Printf("Total time: %s\n", totalTime)
	fmt.Printf("Requests per second (RPS): %.2f\n", rps)
	fmt.Printf("Average latency: %.2f ms\n", avgLatency*1000)
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}

type Object struct {
	Lat  float32
	Lng  float32
	DirX float32
	DirY float32
}

func main() {
	// f, err := os.Create("cpu_profile.prof")

	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()

	// err = pprof.StartCPUProfile(f)

	// if err != nil {
	// 	panic(err)
	// }

	// defer pprof.StopCPUProfile()
	// file, err := openLogFile("./mylog.log")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(file)
	// log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// gRPC connection
	conn, err = grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewObjectClient(conn)

	// Create mock data
	var predefinedPoints = map[int32]Object{}
	for i := 1; i <= 100000; i++ {
		predefinedPoints[int32(i)] = Object{
			Lat:  10.0,
			Lng:  10.0,
			DirX: 0.01,
			DirY: 0.01,
		}
	}

	// Running test
	durationInMinutes := 10
	timeout := time.After(time.Duration(durationInMinutes) * time.Minute) // Create a channel that will receive a value after 10 minutes
	ticker := time.Tick(3 * time.Second)                                  // Create a ticker that ticks every 3 second (adjust as needed)

	for {
		select {
		case <-timeout:
			fmt.Printf("%d minutes have passed. Exiting loop.", durationInMinutes)
			return
		case <-ticker:
			testCreateObject()
			// testHello()
		}
	}
}
