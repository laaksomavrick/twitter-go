package hello

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/gateway/internal/core"
)

// Index returns the status of all the services for the api
func Index(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&map[string]string{
			"hello": "world",
		})
	}
}
