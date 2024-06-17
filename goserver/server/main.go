package main

import (
	"context"
	"database/sql"

	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof"

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
			func() { // Process in a goroutine to avoid blocking the server
				err := sendToPostGIS(objMap.(*sync.Map)) // Pass the whole list
				if err != nil {
					log.Fatal("Error storing objects in PostGIS:", err)
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
		log.Fatal(err)
	}
	defer tx.Rollback() // Ensure rollback on error
	stmt, err := tx.Prepare("INSERT INTO objects (object_id, type, color, location, status, created_at) VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($4, $5),4326), $6, $7)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	objMap.Range(func(_, value any) bool {
		obj := value.(*object.ObjectRequest)

		if _, err = stmt.Exec(obj.Id, obj.Type, obj.Color, obj.Lng, obj.Lat, obj.Status, time.Unix(obj.Timestamp, 0)); err != nil {
			log.Fatal(err)
		}

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
	var err error
	go func() {
		http.ListenAndServe(":10001", nil)
	}()
	// db connect
	dbConnString := "postgres://postgres:2862003(())aa@localhost:5432/miniproject?sslmode=disable"
	db, err = sql.Open("postgres", dbConnString)
	db.SetMaxOpenConns(50)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("The database is connected")
	defer db.Close()

	// create server
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	serverRegistrar := grpc.NewServer()
	service := &myObjectServer{}
	object.RegisterObjectServer(serverRegistrar, service)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatal("impossible to serve:", err)
	}

}
