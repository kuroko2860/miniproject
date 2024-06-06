package main

import (
	"context"
	"database/sql"

	"fmt"
	"log"
	"net"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	object "kuroko.com/goserver/object"
)

var (
	// ... other global variables ...
	db    *sql.DB       // Global database connection
	cache = &sync.Map{} // Map to store objects with timestamps
)

type myObjectServer struct {
	object.UnimplementedObjectServer
}

func (s myObjectServer) CreateObject(ctx context.Context, req *object.ObjectRequest) (*object.ObjectResponse, error) {
	currentTime := time.Now().Unix()
	// record++
	// fmt.Println(record)
	req.Timestamp = &currentTime

	// Get or create the object list for this timestamp in the cache
	key := currentTime
	actual, _ := cache.LoadOrStore(key, &sync.Map{})
	objMap := actual.(*sync.Map)
	// fmt.Println(req.Id)
	objMap.Store(req.Id, req)
	// Schedule cleanup after 1 second to send to the database
	time.AfterFunc(1*time.Second, func() {
		if objMap, ok := cache.LoadAndDelete(key); ok {
			go func() { // Process in a goroutine to avoid blocking the server
				err := sendToPostGIS(objMap.(*sync.Map)) // Pass the whole list
				if err != nil {
					log.Printf("Error storing objects in PostGIS: %v", err)
					// You might want to handle the error more robustly here (e.g., retry)
				}
			}()
		}
	})

	return &object.ObjectResponse{
		Ack: req.Id,
	}, nil
}

func (s myObjectServer) Hello(ctx context.Context, req *object.HelloRequest) (*object.HelloResponse, error) {
	return &object.HelloResponse{
		Message: "Hello World!",
	}, nil
}

func sendToPostGIS(objMap *sync.Map) error {
	tx, err := db.BeginTx(context.Background(), nil) // Use a transaction
	if err != nil {
		return err
	}
	defer tx.Rollback() // Ensure rollback on error

	objMap.Range(func(_, value interface{}) bool {
		obj := value.(*object.ObjectRequest)
		createStmt := `INSERT INTO objects (object_id, type, color, location, status, created_at) VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5),4326), $6, $7)`
		_, err = tx.Exec(createStmt, obj.Id, obj.Type, obj.Color, obj.Lat, obj.Lng, obj.Status, time.Unix(*obj.Timestamp, 0))

		return err == nil // Stop iteration on error else continue
	})

	return tx.Commit() // Commit the transaction if everything is successful
}

// func gen_data() {
// 	const numObjects = 10000
// 	const maxID = 100000
// 	types := []string{"car", "bike", "truck", "bus"}
// 	colors := []string{"red", "green", "blue", "yellow"}
// 	statuses := []string{"moving", "static"}

// 	// Starting coordinates
// 	lat := rand.Float32()*180 - 90
// 	lng := rand.Float32()*360 - 180

// 	// Open the file in write mode (create if not exists)
// 	file, err := os.Create("../test/object_requests.bin")
// 	if err != nil {
// 		log.Fatalf("failed to open file: %v", err)
// 	}
// 	defer file.Close()

// 	for i := 0; i < numObjects; i++ {
// 		objectReq := &object.ObjectRequest{
// 			Id:     fmt.Sprintf("%d", rand.Intn(maxID)+1),
// 			Type:   types[rand.Intn(len(types))],
// 			Color:  colors[rand.Intn(len(colors))],
// 			Lat:    lat,
// 			Lng:    lng,
// 			Status: statuses[rand.Intn(len(statuses))],
// 		}

// 		// Increment lat/lng slightly for each object
// 		lat += rand.Float32() * 0.001
// 		lng += rand.Float32() * 0.001

// 		// Serialize the object
// 		data, err := proto.Marshal(objectReq)
// 		if err != nil {
// 			log.Fatalf("marshaling error: %v", err)
// 		}
// 		buf := make([]byte, 4)
// 		binary.LittleEndian.PutUint32(buf, uint32(len(data)))

// 		if _, err := file.Write(buf); err != nil {
// 			log.Fatalf("failed to write message size: %v", err)
// 		}
// 		if _, err := file.Write(data); err != nil {
// 			log.Fatalf("failed to write message data: %v", err)
// 		}
// 	}
// 	fmt.Println("Binary data written to object_requests.bin")
// }

func main() {
	// gen_data()
	// db connect
	var err error
	dbConnString := "postgres://postgres:2862003(())aa@localhost:5432/miniproject?sslmode=disable"
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("The database is connected")
	defer db.Close()

	// create server
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &myObjectServer{}
	object.RegisterObjectServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
