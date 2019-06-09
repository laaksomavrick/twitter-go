package internal

import (
	"encoding/json"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

// CreateHandler handles creating a tweet record
func CreateHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var tweet Tweet

		if err := json.Unmarshal(msg, &tweet); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"tweet": tweet})
		}

		tweet.prepareForInsert()

		repo := NewRepository(s.Cassandra)
		if err := repo.Insert(tweet); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"tweet": tweet})
		}

		s.Amqp.PublishToTopic(amqp.CreatedTweetKey, []string{tweet.Username}, tweet)

		body, _ := json.Marshal(tweet)

		return &amqp.OkResponse{Body: body}, nil
	}
}

// GetAllHandler handles returning all tweets for a given username
func GetAllHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var req GetAllUserTweets

		if err := json.Unmarshal(msg, &req); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"getAllUserTweets": req})
		}

		repo := NewRepository(s.Cassandra)

		tweets, err := repo.GetAll(req.Username)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"getAllUserTweets": req})
		}

		body, _ := json.Marshal(tweets)

		return &amqp.OkResponse{Body: body}, nil
	}
}
