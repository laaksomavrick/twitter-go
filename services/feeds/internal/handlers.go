package internal

import (
	"encoding/json"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/service"
	"twitter-go/services/common/types"
)

// GetMyFeedHandler handles requests to retrieve an activty feed for a particular user
func GetMyFeedHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var getMyFeed types.GetMyFeed

		if err := json.Unmarshal(msg, &getMyFeed); err != nil {
			return amqp.HandleInternalServiceError(err, nil)
		}

		logger.Info(logger.Loggable{Message: "Retrieving user feed", Data: getMyFeed})

		repo := NewRepository(s.Cassandra)
		feed, err := repo.GetFeed(getMyFeed.Username)
		if err != nil {
			return amqp.HandleInternalServiceError(err, getMyFeed)
		}

		body, err := json.Marshal(feed)
		if err != nil {
			return amqp.HandleInternalServiceError(err, feed)
		}

		logger.Info(logger.Loggable{Message: "Retrieving user feed ok", Data: feed})

		return &amqp.OkResponse{Body: body}, nil
	}
}

// AddTweetToFeedHandler handles broadcasts to add a tweet to the feed of
// all followers of a particular user
func AddTweetToFeedHandler(s *service.Service) func([]byte) {
	return func(msg []byte) {
		var addTweetToFeed types.AddTweetToFeed

		if err := json.Unmarshal(msg, &addTweetToFeed); err != nil {
			logger.Error(logger.Loggable{
				Message: err.Error(),
			})
			return
		}

		logger.Info(logger.Loggable{Message: "Adding tweet to feeds", Data: addTweetToFeed})

		// find all users subscribed to tweetUsername
		getUserFollowers := types.GetUserFollowers{Username: addTweetToFeed.TweetUsername}

		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.GetAllUserFollowers, []string{getUserFollowers.Username}, getUserFollowers)

		if errorResponse != nil {
			logger.Error(logger.Loggable{
				Message: errorResponse.Message,
				Data:    getUserFollowers,
			})
			return
		}

		followers := types.Followers{}

		if err := json.Unmarshal(okResponse.Body, &followers); err != nil {
			logger.Error(logger.Loggable{
				Message: err.Error(),
			})
			return
		}

		logger.Info(logger.Loggable{Message: "Received getAllUserFollowers response", Data: followers})

		// for each user, upsert the tweet to their feed
		repo := NewRepository(s.Cassandra)

		for _, follower := range followers {
			feedItem := types.FeedItem{
				Username:  addTweetToFeed.TweetUsername,
				ID:        addTweetToFeed.TweetID,
				Content:   addTweetToFeed.TweetContent,
				CreatedAt: addTweetToFeed.TweetCreatedAt,
			}
			err := repo.WriteToFeed(follower.Username, feedItem)
			if err != nil {
				logger.Error(logger.Loggable{
					Message: err.Error(),
					Data:    feedItem,
				})
			}
		}

		logger.Info(logger.Loggable{Message: "Adding tweet to feeds ok", Data: nil})
	}
}
