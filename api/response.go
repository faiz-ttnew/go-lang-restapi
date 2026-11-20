package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func writeResponse(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(w, status, APIResponse{
		Message: message,
		Data:    data,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("Could not write JSON response: %v", err)
	}
}
