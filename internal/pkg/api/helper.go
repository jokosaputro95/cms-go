package api

import (
	"encoding/json"
	"net/http"
)

// SendSuccess mengirimkan respons sukses dalam format JSON
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}, meta interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := NewSuccessResponse(message, data, meta)
	json.NewEncoder(w).Encode(response)
}

// SendError mengirimkan respons error dalam format JSON
func SendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := NewErrorResponse(message)
	json.NewEncoder(w).Encode(response)
}