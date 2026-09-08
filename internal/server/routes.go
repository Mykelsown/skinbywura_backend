// Package server contains the route registration and HTTP handlers for the API.
package server

import (
	"encoding/json"
	"net/http"
)

// stats is the JSON payload returned by the health check endpoint.
type stats struct {
	Status string `json:"status"`
}

// loadRoute registers all API endpoints and returns the mux that serves them.
func loadRoute() http.Handler {
	mux := http.NewServeMux()
	// Health checks let deployments verify that the HTTP process is responsive.
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		// Return a JSON response so clients can inspect the service state.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		state := stats{Status: "ok"}

		// Encode the response directly to the HTTP writer.
		err := json.NewEncoder(w).Encode(state)
		if err != nil {
			http.Error(w, "wasn't able to write", http.StatusBadRequest)
			return
		}
	})
	return mux
}
