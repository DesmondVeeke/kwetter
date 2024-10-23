package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	// Handle requests for /api/users/
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		// Forward request to the user service
		proxyRequest(w, r, "http://localhost:8081")
	})

	// Handle requests for /api/products/
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		// Forward request to the product service
		proxyRequest(w, r, "http://localhost:8082")
	})

	// Handle requests for health check
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Optional, default is 200 OK
		w.Write([]byte("Health check ok!"))
	})

	// Start the API Gateway server
	log.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyRequest(w http.ResponseWriter, r *http.Request, serviceURL string) {
	// Create a new request to the destination service
	req, err := http.NewRequest(r.Method, serviceURL+r.URL.Path[len("/api/"):], r.Body) // Adjust the URL path
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Forward the headers from the incoming request
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Send the request to the service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy the service response back to the client
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
