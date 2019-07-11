package internal

import (
	"encoding/json"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/service"
	"twitter-go/services/common/types"
)

func FollowUserHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var followUser types.FollowUser

		if err := json.Unmarshal(msg, &followUser); err != nil {
			return amqp.HandleInternalServiceError(err, nil)
		}

		logger.Info(logger.Loggable{Message: "Following user", Data: followUser})

		repo := NewRepository(s.Cassandra)

		err := repo.FollowUser(followUser.Username, followUser.FollowingUsername)
		if err != nil {
			return amqp.HandleInternalServiceError(err, followUser)
		}

		body, err := json.Marshal(followUser)
		if err != nil {
			return amqp.HandleInternalServiceError(err, followUser)
		}

		logger.Info(logger.Loggable{Message: "Following user ok", Data: nil})

		return &amqp.OkResponse{Body: body}, nil
	}
}

func GetUserFollowersHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var getUserFollowers types.GetUserFollowers

		if err := json.Unmarshal(msg, &getUserFollowers); err != nil {
			return amqp.HandleInternalServiceError(err, getUserFollowers)
		}

		logger.Info(logger.Loggable{Message: "Getting user followers", Data: getUserFollowers})

		repo := NewRepository(s.Cassandra)

		followers, err := repo.GetUserFollowers(getUserFollowers.Username)
		if err != nil {
			return amqp.HandleInternalServiceError(err, getUserFollowers)
		}

		body, err := json.Marshal(followers)
		if err != nil {
			return amqp.HandleInternalServiceError(err, getUserFollowers)
		}

		logger.Info(logger.Loggable{Message: "Get user followers ok", Data: followers})

		return &amqp.OkResponse{Body: body}, nil
	}
}
