package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		proxyRequest(w, r, "http://localhost:8081")
	})

	http.HandleFunc("/api/write-kweet", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received..")
		proxyRequest(w, r, "http://write-kweet:8082/write-kweet")
	})

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health check ok!"))
	})

	log.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyRequest(w http.ResponseWriter, r *http.Request, serviceURL string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		log.Println("Error reading request body:", err)
		return
	}
	defer r.Body.Close()

	// Forward the request without appending the path again
	forwardURL := serviceURL
	log.Printf("Forwarding request to: %s", forwardURL)

	req, err := http.NewRequest(r.Method, forwardURL, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Println("Error creating request:", err)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusBadGateway)
		log.Println("Error forwarding request:", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response: %d", resp.StatusCode)

	// Forward response headers and body
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println("Error copying response body:", err)
	}
}
