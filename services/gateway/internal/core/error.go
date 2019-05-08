package core

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	UnprocessableEntity = "Unprocessable request sent."
	BadRequest          = "Bad request sent."
	Forbidden           = "Forbidden."
)

// ErrorResponse defines the shape of the default error response served by the application
type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// ErrorsResponse defines the shape of the default error response served by the application
type ErrorsResponse struct {
	Status int                    `json:"status"`
	Errors map[string]interface{} `json:"errors"`
}

// TODO-8: unify err and errs response
// Error -> Message ?

// EncodeJSONError issues an ErrorResponse payload to the client
func EncodeJSONError(w http.ResponseWriter, err error, status int) {
	errResp := &ErrorResponse{
		Status: status,
		Error:  err.Error(), //TODO-10: error should be derivative from status, not user provided ***
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}

// EncodeJSONErrors issues an ErrorsResponse payload to the client
func EncodeJSONErrors(w http.ResponseWriter, errs url.Values, status int) {
	errsMap := map[string]interface{}{
		"errors": errs,
	}
	errResp := &ErrorsResponse{
		Status: status,
		Errors: errsMap,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}
