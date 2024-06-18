package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/rs/cors"
)

const (
	port         = ":9090"
	dbConnString = "postgres://postgres:2862003(())aa@localhost:5432/miniproject?sslmode=disable" // Replace with your credentials
)

var db *sql.DB

// ... (Your ObjectData struct) ...

func countObjectsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	polygonWKT := r.URL.Query().Get("polygon")

	// fmt.Println(start, end, polygonWKT)

	// Get optional filter parameters (type, color, status)
	objectType := r.URL.Query().Get("type")
	objectColor := r.URL.Query().Get("color")
	objectStatus := r.URL.Query().Get("status")

	// Call the PostGIS function with optional parameters
	var count int64
	query := `
        SELECT count_objects_within_polygon_and_time_range(
            $1, $2, $3, $4, $5, $6
        )
    `
	err := db.QueryRow(query, start, end, polygonWKT, objectType, objectColor, objectStatus).Scan(&count)
	// ... (error handling) ...
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	// Construct the JSON response
	response := struct {
		ObjectCount int64 `json:"object_count"`
	}{
		ObjectCount: count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	var err error
	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Check database connection
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	router := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := router.Handler(http.HandlerFunc(countObjectsHandler))
	http.Handle("/api/objects/count", handler)

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
