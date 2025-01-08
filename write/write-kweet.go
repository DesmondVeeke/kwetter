package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Kweet represents the structure of a kweet
type Kweet struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

var db *sql.DB

func createTableIfNotExists() {
	query := `
        CREATE TABLE IF NOT EXISTS kweets (
            id SERIAL PRIMARY KEY,
            user_name TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW()
        );
    `

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Table 'kweets' ensured to exist")
}

func initDB() {
	// Load environment variables
	connStr := fmt.Sprintf("host=postgres-kweets user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER_KWEET"),
		os.Getenv("POSTGRES_PASSWORD_KWEET"),
		os.Getenv("POSTGRES_DB_KWEET"))

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the database connection
	if err = db.Ping(); err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}

	// Create the table if it doesn't exist
	createTableIfNotExists()
	log.Println("Connected to PostgreSQL successfully")
}

// WriteKweetHandler handles incoming requests to write a kweet
func WriteKweetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the kweet from request body
	var kweet Kweet
	if err := json.NewDecoder(r.Body).Decode(&kweet); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Insert the kweet into the database
	query := "INSERT INTO kweets (user_name, content) VALUES ($1, $2)"
	_, err := db.Exec(query, kweet.User, kweet.Content)
	if err != nil {
		log.Printf("Failed to save kweet: %v", err)
		http.Error(w, "Failed to save kweet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Kweet successfully created!"))
	log.Println("Kweet created successfully")
}

func main() {
	// Initialize the database
	initDB()

	// Set up HTTP handler
	http.HandleFunc("/write-kweet", WriteKweetHandler)

	// Start the server
	port := "8082"
	log.Printf("Write-Kweet service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
