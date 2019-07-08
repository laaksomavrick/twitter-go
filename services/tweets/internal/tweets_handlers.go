package internal

import (
	"encoding/json"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/service"
)

// CreateHandler handles creating a tweet record
func CreateHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var tweet Tweet

		if err := json.Unmarshal(msg, &tweet); err != nil {
			return amqp.HandleInternalServiceError(err, nil)
		}

		logger.Info(logger.Loggable{Message: "Creating tweet", Data: tweet})

		tweet.prepareForInsert()

		repo := NewRepository(s.Cassandra)
		if err := repo.Insert(tweet); err != nil {
			return amqp.HandleInternalServiceError(err, tweet)
		}

		logger.Info(logger.Loggable{Message: "Publishing tweet created event", Data: tweet})

		s.Amqp.PublishToTopic(amqp.CreatedTweetKey, []string{tweet.Username}, tweet)

		body, _ := json.Marshal(tweet)

		logger.Info(logger.Loggable{Message: "Create tweet ok", Data: nil})

		return &amqp.OkResponse{Body: body}, nil
	}
}

// GetAllHandler handles returning all tweets for a given username
func GetAllHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var req GetAllUserTweets

		if err := json.Unmarshal(msg, &req); err != nil {
			return amqp.HandleInternalServiceError(err, req)
		}

		logger.Info(logger.Loggable{Message: "Getting all tweets", Data: req})

		repo := NewRepository(s.Cassandra)

		tweets, err := repo.GetAll(req.Username)
		if err != nil {
			return amqp.HandleInternalServiceError(err, req)
		}

		body, _ := json.Marshal(tweets)

		logger.Info(logger.Loggable{Message: "Get all tweets ok", Data: tweets})

		return &amqp.OkResponse{Body: body}, nil
	}
}
