package feeds

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/gateway/internal/core"
)

func GetMyFeed(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtUsername := core.GetUsernameFromRequest(r)

		getFeedDto := &GetFeedDto{Username: jwtUsername}

		json.NewEncoder(w).Encode(getFeedDto)
	}
}
