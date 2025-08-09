package api

import (
	"encoding/json"
	"net/http"
)

// Response adalah struct standar untuk format respons API
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorDetail menyimpan informasi spesifik tentang error.
// Details menggunakan map untuk fleksibilitas.
type ErrorDetail struct {
	Type    string                 `json:"type"`
	Details map[string]interface{} `json:"details,omitempty"`
}


// ErrorResponse digunakan untuk format respons error yang lebih detail.
type ErrorResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Error   ErrorDetail `json:"error"`
}

// --- Fungsi Helper untuk Mengirimkan Response ---
// SendJSON mengirimkan respons API dalam format JSON.
func SendJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// SendSuccess membuat dan mengirimkan respons sukses.
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}, meta interface{}) {
	response := &Response{
		Status:  true,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
	SendJSON(w, statusCode, response)
}

// SendError membuat dan mengirimkan respons error umum.
func SendError(w http.ResponseWriter, statusCode int, message string) {
	response := &Response{
		Status:  false,
		Message: message,
	}
	SendJSON(w, statusCode, response)
}

// SendDetailedError membuat dan mengirimkan respons error dengan detail.
func SendDetailedError(w http.ResponseWriter, statusCode int, message, errorType string, details map[string]interface{}) {
	response := &ErrorResponse{
		Status:  false,
		Message: message,
		Error: ErrorDetail{
			Type:    errorType,
			Details: details,
		},
	}
	SendJSON(w, statusCode, response)
}