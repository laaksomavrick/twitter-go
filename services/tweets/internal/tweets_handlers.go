package internal

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

// CreateHandler handles creating a tweet record
func CreateHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {

		var tweet Tweet

		if err := json.Unmarshal(msg, &tweet); err != nil {
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		tweet.prepareForInsert()

		repo := NewRepository(s.Cassandra)
		if err := repo.Insert(tweet); err != nil {
			return nil, err
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
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		repo := NewRepository(s.Cassandra)
		tweets, err := repo.GetAll(req.Username)
		if err != nil {
			return nil, err
		}

		body, _ := json.Marshal(tweets)

		return &amqp.OkResponse{Body: body}, nil

	}
}
