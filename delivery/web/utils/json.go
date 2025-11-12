package utils

import (
	"encoding/json"
	"net/http"
)

// writeJSON is a helper to write JSON responses.
// type any = interface{} any is an alias for interface{} and is equivalent to interface{} in all ways.
func WriteJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json") // Set content type header
	w.WriteHeader(status)

	// Set the status code
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func JsonError(err error) string {
	type response struct {
		Error string `json:"error"`
	}

	rsp := response{Error: err.Error()}
	jsonData, _ := json.Marshal(rsp)

	return string(jsonData)
}
