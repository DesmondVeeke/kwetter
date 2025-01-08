package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	// Handler for the /api/auth/ endpoint
	http.HandleFunc("/api/auth/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Auth request received.")
		// Strip /api/auth from the path
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/auth")
		proxyRequest(w, r, "http://keycloak:8080")
	})

	// Handler for the /api/write-kweet endpoint
	http.HandleFunc("/api/write-kweet", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Write Kweet request received.")
		//Strip /api/write-kweet from the path
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/write-kweet")
		proxyRequest(w, r, "http://write-kweet:8082/write-kweet")
	})

	// Health check endpoint
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health check ok!"))
	})

	log.Println("API Gateway running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyRequest(w http.ResponseWriter, r *http.Request, serviceURL string) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		log.Println("Error reading request body:", err)
		return
	}
	defer r.Body.Close()

	// Create the forward request URL
	forwardURL := serviceURL + r.URL.Path
	log.Printf("Forwarding request to: %s", forwardURL)

	req, err := http.NewRequest(r.Method, forwardURL, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Println("Error creating request:", err)
		return
	}

	// Copy headers from the original request
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

	// Handle redirects (rewrite URLs in response headers if necessary)
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect || resp.StatusCode == http.StatusPermanentRedirect {
		redirectURL := resp.Header.Get("Location")
		if redirectURL != "" {
			parsedURL, err := url.Parse(redirectURL)
			if err == nil && strings.HasPrefix(parsedURL.Path, "/admin") {
				parsedURL.Host = "localhost:8080"
				resp.Header.Set("Location", parsedURL.String())
			}
		}
	}

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
