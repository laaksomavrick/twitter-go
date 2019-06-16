package internal

import (
	"encoding/json"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

func FollowUserHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var followUser FollowUser

		if err := json.Unmarshal(msg, &followUser); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"followUser": followUser})
		}

		repo := NewRepository(s.Cassandra)

		err := repo.FollowUser(followUser.Username, followUser.FollowingUsername)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"followUser": followUser})
		}

		body, err := json.Marshal(followUser)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"followUser": followUser})
		}

		return &amqp.OkResponse{Body: body}, nil
	}
}

func GetUserFollowersHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var getUserFollowers GetUserFollowers

		if err := json.Unmarshal(msg, &getUserFollowers); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"getUserFollowers": getUserFollowers})
		}

		repo := NewRepository(s.Cassandra)

		followers, err := repo.GetUserFollowers(getUserFollowers.Username)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"getUserFollowers": getUserFollowers})
		}

		body, err := json.Marshal(followers)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"getUserFollowers": getUserFollowers})
		}

		return &amqp.OkResponse{Body: body}, nil
	}
}
