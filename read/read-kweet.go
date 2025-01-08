package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
)

type Kweet struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Content string `json:"content"`
}

var db *sql.DB

func init() {
	connStr := os.Getenv("POSTGRES_CONN")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func ReadKweetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract kweet ID from query params
	kweetID := r.URL.Query().Get("id")
	if kweetID == "" {
		http.Error(w, "Missing kweet ID", http.StatusBadRequest)
		return
	}

	// Query the database
	var kweet Kweet
	query := "SELECT id, user_name, content FROM kweets WHERE id = $1"
	err := db.QueryRow(query, kweetID).Scan(&kweet.ID, &kweet.User, &kweet.Content)
	if err == sql.ErrNoRows {
		http.Error(w, "Kweet not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to read kweet", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kweet)
}
