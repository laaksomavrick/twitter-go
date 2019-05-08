package tweets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"twitter-go/services/gateway/internal/core"
)

// CreateHandler handles creating a new tweet
func CreateHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.Context().Value("username"))

		json.NewEncoder(w).Encode(map[string]string{
			"hello": "1",
		})
	}
}
