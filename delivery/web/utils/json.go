package utils

import (
	"encoding/json"
	"net/http"
)

// writeJSON is a helper to write JSON responses.
// type any = interface{} any is an alias for interface{} and is equivalent to interface{} in all ways.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json") // Set content type header

	w.WriteHeader(status)           // Set the status code
	json.NewEncoder(w).Encode(data) // Encode and write the JSON response
}
