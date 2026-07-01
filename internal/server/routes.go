package server

import (
	"encoding/json"
	"net/http"
)

type stats struct {
	Status string `json:"status"`
}

func loadRoute() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		
		state := stats{"ok"}

		err := json.NewEncoder(w).Encode(state)
		if err != nil {
			http.Error(w, "wasn't able to write", 400)
			return 
		}
	})
	return mux
}