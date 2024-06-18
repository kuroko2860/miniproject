package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"os"

// 	"sync"
// 	"time"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// 	pb "kuroko.com/goserver/object"
// )

// var (
// 	conn             *grpc.ClientConn
// 	wg               sync.WaitGroup
// 	client           pb.ObjectClient
// 	predefinedPoints []Object
// )

// func testHello() {
// 	// gRPC connection
// 	conn, err := grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	client := pb.NewObjectClient(conn)
// 	var wg sync.WaitGroup
// 	numRequests := 100000
// 	concurrency := 8000 // Adjust this based on your needs

// 	startTime := time.Now()

// 	for i := 0; i < numRequests; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			client.Hello(context.Background(), &pb.HelloRequest{})
// 		}()
// 		if i%concurrency == 0 {
// 			wg.Wait() // Wait for a batch of goroutines to finish
// 		}
// 	}
// 	wg.Wait() // Wait for all remaining goroutines

// 	endTime := time.Now()
// 	totalTime := endTime.Sub(startTime)

// 	rps := float64(numRequests) / totalTime.Seconds()
// 	avgLatency := totalTime.Seconds() / float64(numRequests)

// 	fmt.Printf("\nTotal requests: %d\n", numRequests)
// 	fmt.Printf("Total time: %s\n", totalTime)
// 	fmt.Printf("Requests per second (RPS): %.2f\n", rps)
// 	fmt.Printf("Average latency: %.2f ms\n", avgLatency*1000)
// }

// func testCreateObject() {
// 	types := []string{"car", "bike", "truck", "bus"}
// 	l_types := len(types)
// 	colors := []string{"red", "green", "blue", "yellow"}
// 	l_colors := len(colors)
// 	statuses := []string{"moving", "static"}
// 	l_statuses := len(statuses)

// 	var numRequests int32 = 100000
// 	var concurrency int32 = 8000 // Adjust this based on your needs

// 	startTime := time.Now()
// 	type_idx := 0
// 	color_idx := 0
// 	status_idx := 0
// 	var i int32
// 	for i = 1; i <= numRequests; i++ {
// 		wg.Add(1)
// 		go func(j int32) {
// 			defer wg.Done()
// 			obj := predefinedPoints[j%8000]
// 			objectReq := &pb.ObjectRequest{
// 				Id:        obj.Id,
// 				Type:      types[type_idx],
// 				Color:     colors[color_idx],
// 				Lat:       obj.Lat,
// 				Lng:       obj.Lng,
// 				Status:    statuses[status_idx],
// 				Timestamp: time.Now().Unix(),
// 			}
// 			// log.Println("id", objectReq.Id, " time:", objectReq.Timestamp)

// 			// Increment lat/lng slightly for each object
// 			obj.Lng = obj.Lng + obj.DirX
// 			obj.Lat = obj.Lat + obj.DirY
// 			_, err := client.CreateObject(context.Background(), objectReq)
// 			if err != nil {
// 				log.Printf("Could not create object: %v", err)
// 			}
// 		}(i)
// 		if i%concurrency == 0 {
// 			type_idx = (type_idx + 1) % l_types
// 			color_idx = (color_idx + 1) % l_colors
// 			status_idx = (status_idx + 1) % l_statuses
// 			wg.Wait() // Wait for a batch of goroutines to finish
// 		}
// 	}
// 	wg.Wait() // Wait for all remaining goroutines

// 	endTime := time.Now()
// 	totalTime := endTime.Sub(startTime)

// 	rps := float64(numRequests) / totalTime.Seconds()
// 	avgLatency := totalTime.Seconds() / float64(numRequests)

// 	fmt.Printf("Total requests: %d\n", numRequests)
// 	fmt.Printf("Total time: %s\n", totalTime)
// 	fmt.Printf("Requests per second (RPS): %.2f\n", rps)
// 	fmt.Printf("Average latency: %.2f ms\n", avgLatency*1000)
// }

// func openLogFile(path string) (*os.File, error) {
// 	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return logFile, nil
// }

// type Object struct {
// 	Id   int32
// 	Lng  float32 // kinh do
// 	Lat  float32 // vi do
// 	DirX float32
// 	DirY float32
// }

// func main() {
// 	// f, err := os.Create("cpu_profile.prof")

// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// defer f.Close()

// 	// err = pprof.StartCPUProfile(f)

// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// defer pprof.StopCPUProfile()
// 	file, err := openLogFile("./mylog.log")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.SetOutput(file)
// 	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

// 	// gRPC connection
// 	conn, err = grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	client = pb.NewObjectClient(conn)

// 	// Create mock data
// 	predefinedPoints = make([]Object, 8000)
// 	for i := 1; i <= 8000; i++ {
// 		obj := Object{
// 			Id:   int32(i),
// 			Lng:  rand.Float32()*(109.5-102.1) + 102.1,
// 			Lat:  rand.Float32()*(15.3-8.5) + 8.5,
// 			DirX: (rand.Float32()*(0.9-0.1) + 0.1) * 0.01,
// 			DirY: (rand.Float32()*(0.9-0.1) + 0.1) * 0.01,
// 		}
// 		predefinedPoints[i-1] = obj
// 		// log.Println(obj)
// 	}

// 	// Running test
// 	durationInMinutes := 10
// 	timeout := time.After(time.Duration(durationInMinutes) * time.Minute) // Create a channel that will receive a value after 10 minutes
// 	ticker := time.Tick(3 * time.Second)                                  // Create a ticker that ticks every 3 second (adjust as needed)

//		for {
//			select {
//			case <-timeout:
//				fmt.Printf("%d minutes have passed. Exiting loop.", durationInMinutes)
//				return
//			case <-ticker:
//				testCreateObject()
//				// os.Exit(0)
//				// testHello()
//			}
//		}
//	}
