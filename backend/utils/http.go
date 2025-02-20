package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data    any    `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SendError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{
		Data:    nil,
		Success: false,
		Message: msg,
	}
	json.NewEncoder(w).Encode(response)
}

func SendMessage(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Response{
		Data:    data,
		Success: true,
		Message: "Success",
	}
	json.NewEncoder(w).Encode(response)
}
