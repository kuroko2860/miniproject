package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "kuroko.com/goserver/object"
)

var (
	conn             *grpc.ClientConn
	wg               sync.WaitGroup
	client           pb.ObjectClient
	predefinedPoints []Object
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

	var numRequests int32 = 200000
	var concurrency int32 = 1000 // Adjust this based on your needs

	startTime := time.Now()

	var i int32
	wg.Add(int(concurrency))
	for i = 1; i <= concurrency; i++ {
		go func(j int32) {
			defer wg.Done()
			obj_id := rand.Int31n(100000) + 1
			obj := predefinedPoints[obj_id-1]
			for t := 0; t < int(numRequests/concurrency); t++ {
				objectReq := &pb.ObjectRequest{
					Id:        obj.Id,
					Type:      obj.Type,
					Color:     obj.Color,
					Lat:       obj.Lat,
					Lng:       obj.Lng,
					Status:    obj.Status,
					Timestamp: time.Now().Unix(),
				}
				// log.Println("id", objectReq.Id, " time:", objectReq.Timestamp)

				// Increment lat/lng slightly for each object
				obj.Lng = obj.Lng + obj.DirX
				obj.Lat = obj.Lat + obj.DirY
				_, err := client.CreateObject(context.Background(), objectReq)
				if err != nil {
					log.Printf("Could not create object: %v\n", err)
				}
			}
		}(i)
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
	Id     int32
	Lng    float32 // kinh do
	Lat    float32 // vi do
	Type   string
	Color  string
	Status string
	DirX   float32
	DirY   float32
}

func main() {
	types := []string{"car", "bike", "truck", "bus"}
	l_types := len(types)
	colors := []string{"red", "green", "blue", "yellow"}
	l_colors := len(colors)
	// statuses := []string{"moving", "static"}
	// l_statuses := len(statuses)

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
	file, err := openLogFile("./mylog.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// gRPC connection
	conn, err = grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client = pb.NewObjectClient(conn)
	fmt.Println("Creating mock data...")
	// Create mock data
	predefinedPoints = make([]Object, 100000)
	for i := 1; i <= 100000; i++ {
		obj := Object{
			Id:     int32(i),
			Type:   types[rand.Intn(l_types)],
			Color:  colors[rand.Intn(l_colors)],
			Status: "moving",
			Lng:    rand.Float32()*(109.5-102.1) + 102.1,
			Lat:    rand.Float32()*(23.3-8.5) + 8.5,
			DirX:   (rand.Float32()*(0.9-0.1) + 0.1) * 0.001,
			DirY:   (rand.Float32()*(0.9-0.1) + 0.1) * 0.001,
		}
		predefinedPoints[i-1] = obj
		// log.Println(obj)
	}

	fmt.Println("Created 100.000 mock object...")
	fmt.Println("Running test...")
	// Running test
	durationInMinutes := 10
	timeout := time.After(time.Duration(durationInMinutes) * time.Minute) // Create a channel that will receive a value after 10 minutes
	ticker := time.Tick(3 * time.Second)                                  // Create a ticker that ticks every 3 second (adjust as needed)

	for {
		select {
		case <-timeout:
			fmt.Printf("%d minutes have passed. Done testing.", durationInMinutes)
			return
		case <-ticker:
			testCreateObject()
			// os.Exit(0)
			// testHello()
		}
	}
}
