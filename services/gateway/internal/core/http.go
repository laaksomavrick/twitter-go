package core

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
}

// Ok issues a 200 response in a uniform format across the api
func Ok(w http.ResponseWriter, data interface{}) {
	status := 200
	response := apiResponse{
		Status: status,
		Data:   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Error issues a 4xx or 5xx response in a uniform format across the api
func Error(w http.ResponseWriter, status int, err string) {
	response := apiResponse{
		Status: status,
		Error:  err,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
