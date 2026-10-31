package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// maxBodyBytes caps request bodies so a client cannot exhaust memory.
const maxBodyBytes = 1 << 20

// WriteJSON writes a JSON response with the provided HTTP status code.
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode response", "err", err)
	}
}

// WriteErrorJSON writes an error payload in a consistent JSON shape.
func WriteErrorJSON(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, map[string]string{"error": message})
}

// serverError logs the real cause and returns a generic message, so internals never leak to clients.
func serverError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	slog.Error(msg, "err", err, "path", r.URL.Path)
	WriteErrorJSON(w, http.StatusInternalServerError, msg)
}

// decodeJSON reads a size-limited JSON body into dst.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if dec.More() {
		return errors.New("unexpected data after JSON body")
	}
	return nil
}
