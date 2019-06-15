package internal

import (
	"encoding/json"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

func GetMyFeedHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var getMyFeed GetMyFeed

		if err := json.Unmarshal(msg, &getMyFeed); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"getMyFeed": getMyFeed})
		}

		repo := NewRepository(s.Cassandra)
		feed, err := repo.GetFeed(getMyFeed.Username)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"getMyFeed": getMyFeed})
		}

		body, err := json.Marshal(feed)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"feed": feed})
		}

		return &amqp.OkResponse{Body: body}, nil
	}
}
