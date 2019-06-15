package feeds

import (
	"encoding/json"
	"errors"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/gateway/internal/core"
)

func GetMyFeed(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtUsername := core.GetUsernameFromRequest(r)

		getFeedDto := &GetFeedDto{Username: jwtUsername}

		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.GetMyFeedKey, []string{getFeedDto.Username}, getFeedDto)

		if errorResponse != nil {
			core.EncodeJSONError(w, errors.New(errorResponse.Message), errorResponse.Status)
			return
		}

		feed := make([]map[string]interface{}, 0)

		if err := json.Unmarshal(okResponse.Body, &feed); err != nil {
			core.EncodeJSONError(w, errors.New(core.InternalServerError), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(feed)
	}
}
