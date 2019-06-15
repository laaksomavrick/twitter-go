package feeds

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/gateway/internal/core"
)

func GetMyFeed(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foo := map[string]string{"hello": "world"}
		json.NewEncoder(w).Encode(foo)
	}
}
