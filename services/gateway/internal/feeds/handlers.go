package feeds

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/types"
	"twitter-go/services/gateway/internal/core"
)

// GetMyFeedHandler provides a HandlerFunc for retrieving a user's feed
func GetMyFeedHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtUsername := core.GetUsernameFromRequest(r)

		getFeedDto := types.GetMyFeed{Username: jwtUsername}

		logger.Info(logger.Loggable{Message: "Getting user feed", Data: getFeedDto})

		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.GetMyFeedKey, []string{getFeedDto.Username}, getFeedDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		feed := make([]map[string]interface{}, 0)

		if err := json.Unmarshal(okResponse.Body, &feed); err != nil {
			core.Error(w, http.StatusInternalServerError, core.InternalServerError)
			return
		}

		core.Ok(w, feed)
	}
}
