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
	db    *sql.DB       // Global database connection
	cache = &sync.Map{} // Map to store objects with timestamps
)

type myObjectServer struct {
	object.UnimplementedObjectServer
}

func (s myObjectServer) CreateObject(ctx context.Context, req *object.ObjectRequest) (*object.ObjectResponse, error) {
	// Get or create the object list for this timestamp in the cache
	key := req.Timestamp
	actual, _ := cache.LoadOrStore(key, &sync.Map{})
	objMap := actual.(*sync.Map)
	objMap.Store(req.Id, req)
	// Schedule cleanup after 1 second to send to the database, not blocking the go routine
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

func sendToPostGIS(objMap *sync.Map) error {
	tx, err := db.BeginTx(context.Background(), nil) // Use a transaction
	if err != nil {
		return err
	}
	defer tx.Rollback() // Ensure rollback on error

	objMap.Range(func(_, value interface{}) bool {
		obj := value.(*object.ObjectRequest)
		createStmt := `INSERT INTO objects (object_id, type, color, location, status, created_at) VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5),4326), $6, $7)`
		_, err = tx.Exec(createStmt, obj.Id, obj.Type, obj.Color, obj.Lat, obj.Lng, obj.Status, time.Unix(obj.Timestamp, 0))

		return err == nil // Stop iteration on error else continue
	})

	return tx.Commit() // Commit the transaction if everything is successful
}

func (s myObjectServer) Hello(ctx context.Context, req *object.HelloRequest) (*object.HelloResponse, error) {
	return &object.HelloResponse{
		Message: "Hello World!",
	}, nil
}

func main() {
	// db connect
	var err error
	dbConnString := "postgres://postgres:2862003(())aa@localhost:5432/miniproject?sslmode=disable"
	db, err = sql.Open("postgres", dbConnString)
	db.SetMaxOpenConns(50)
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
