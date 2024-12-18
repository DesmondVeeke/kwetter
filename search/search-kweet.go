package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Hits struct {
		Hits []struct {
			Source Kweet `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Kweet struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

func SearchKweetsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	// Build Elasticsearch query
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"content": query,
			},
		},
	}
	queryBody, _ := json.Marshal(esQuery)

	// Send request to Elasticsearch
	esURL := os.Getenv("ELASTICSEARCH_URL") + "/kweets/_search"
	resp, err := http.Post(esURL, "application/json", bytes.NewBuffer(queryBody))
	if err != nil {
		http.Error(w, "Failed to search kweets", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse response
	body, _ := ioutil.ReadAll(resp.Body)
	var searchResponse SearchResponse
	json.Unmarshal(body, &searchResponse)

	// Return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(searchResponse.Hits.Hits)
}
