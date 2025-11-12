package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// writeJSON is a helper to write JSON responses.
// type any = interface{} any is an alias for interface{} and is equivalent to interface{} in all ways.
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	fmt.Printf("WriteJSON called: status=%d, data=%+v\n", status, data)

	w.Header().Set("Content-Type", "application/json") // Set content type header
	w.WriteHeader(status)

	// Set the status code
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Printf("WriteJSON ERROR: Failed to encode JSON: %v\n", err)
		return err
	}

	fmt.Printf("WriteJSON success: wrote status %d\n", status)
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
