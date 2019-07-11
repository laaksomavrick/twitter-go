package internal

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/auth"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/service"
	"twitter-go/services/common/types"
)

// CreateHandler handles creating a user record
func CreateHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var user types.User

		if err := json.Unmarshal(msg, &user); err != nil {
			return amqp.HandleInternalServiceError(err, nil)
		}

		repo := NewRepository(s.Cassandra)

		exists, err := repo.Exists(user.Username)

		if err != nil {
			return amqp.HandleInternalServiceError(err, user)
		}

		if exists == true {
			return nil, &amqp.ErrorResponse{Message: "User already exists", Status: http.StatusConflict}
		}

		if err := user.PrepareForInsert(); err != nil {
			return amqp.HandleInternalServiceError(err, user)
		}

		if err := repo.Insert(user); err != nil {
			return amqp.HandleInternalServiceError(err, user)
		}

		accessToken, err := auth.GenerateToken(user.Username, s.Config.HMACSecret)
		if err != nil {
			return amqp.HandleInternalServiceError(err, user)
		}

		user.AccessToken = accessToken
		user.Sanitize()

		// Don't want to log user password ;)
		logger.Info(logger.Loggable{Message: "Create user ok", Data: user})

		body, _ := json.Marshal(user)

		return &amqp.OkResponse{Body: body}, nil
	}

}

// AuthorizeHandler handles authorizing a user given their username and password
func AuthorizeHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var authorizeDto types.Authorize

		if err := json.Unmarshal(msg, &authorizeDto); err != nil {
			return amqp.HandleInternalServiceError(err, nil)
		}

		logger.Info(
			logger.Loggable{
				Message: "Authorizing user",
				Data:    map[string]interface{}{"username": authorizeDto.Username},
			},
		)

		// find user from given username
		repo := NewRepository(s.Cassandra)
		userRecord, err := repo.FindByUsername(authorizeDto.Username)
		if err != nil {
			return nil, &amqp.ErrorResponse{Message: "Not found", Status: http.StatusNotFound}
		}

		// compare password against hash
		if err := userRecord.CompareHashAndPassword(authorizeDto.Password); err != nil {
			return nil, &amqp.ErrorResponse{Message: "Invalid password provided", Status: http.StatusUnprocessableEntity}
		}

		// return new accessToken and refreshToken from record
		accessToken, err := auth.GenerateToken(authorizeDto.Username, s.Config.HMACSecret)
		if err != nil {
			return amqp.HandleInternalServiceError(err, authorizeDto)
		}

		authorized := types.Authorized{
			RefreshToken: userRecord.RefreshToken,
			AccessToken:  accessToken,
		}

		body, _ := json.Marshal(authorized)

		logger.Info(
			logger.Loggable{
				Message: "Authorize user ok",
				Data:    map[string]interface{}{"username": authorizeDto.Username},
			},
		)

		return &amqp.OkResponse{Body: body}, nil
	}
}

// ExistsHandler verifies a user exists in the system
func ExistsHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var existsDto types.DoesExist

		if err := json.Unmarshal(msg, &existsDto); err != nil {
			return amqp.HandleInternalServiceError(err, nil)
		}

		logger.Info(logger.Loggable{Message: "Checking that user exists", Data: existsDto})

		repo := NewRepository(s.Cassandra)

		exists, err := repo.Exists(existsDto.Username)

		if err != nil {
			return amqp.HandleInternalServiceError(err, existsDto)
		}

		if exists == false {
			return nil, &amqp.ErrorResponse{Message: "User not found", Status: http.StatusNotFound}
		}

		existsResponse := types.Exists{
			Exists: true,
		}

		body, _ := json.Marshal(existsResponse)

		logger.Info(logger.Loggable{Message: "User exists ok", Data: existsResponse})

		return &amqp.OkResponse{Body: body}, nil
	}
}
