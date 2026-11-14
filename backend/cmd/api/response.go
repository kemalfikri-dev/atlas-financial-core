package main

import (
	"encoding/json"
	"net/http"
)

type ErrorDetail struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func errorWriter(w http.ResponseWriter, status int, code string, message string) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Status:  status,
			Code:    code,
			Message: message,
		},
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)

}
