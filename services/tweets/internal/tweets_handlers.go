package internal

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

// CreateHandler handles creating a tweet record
func CreateHandler(s *service.Service) func([]byte) interface{} {
	return func(msg []byte) interface{} {

		var tweet Tweet

		if err := json.Unmarshal(msg, &tweet); err != nil {
			return amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		tweet.prepareForInsert()

		repo := NewRepository(s.Cassandra)
		if err := repo.Insert(tweet); err != nil {
			return err
		}

		s.Amqp.PublishToTopic(amqp.CreatedTweetKey, []string{tweet.Username}, tweet)

		return tweet
	}
}

// GetAllHandler handles returning all tweets for a given username
func GetAllHandler(s *service.Service) func([]byte) interface{} {
	return func(msg []byte) interface{} {
		var req GetAllUserTweets

		if err := json.Unmarshal(msg, &req); err != nil {
			return amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		repo := NewRepository(s.Cassandra)
		tweets, err := repo.GetAll(req.Username)
		if err != nil {
			return err
		}

		return tweets
	}
}
