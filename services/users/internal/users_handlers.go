package internal

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/auth"
	"twitter-go/services/common/service"
)

// CreateHandler handles creating a user record
func CreateHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var user User

		if err := json.Unmarshal(msg, &user); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"user": user})
		}

		repo := NewRepository(s.Cassandra)

		exists, err := repo.Exists(user.Username)

		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"user": user})
		}

		if exists == true {
			return nil, &amqp.ErrorResponse{Message: "User already exists", Status: http.StatusConflict}
		}

		if err := user.prepareForInsert(); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"user": user})
		}

		if err := repo.Insert(user); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"user": user})
		}

		accessToken, err := auth.GenerateToken(user.Username, s.Config.HMACSecret)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"user": user})
		}

		user.AccessToken = accessToken
		user.sanitize()

		body, _ := json.Marshal(user)

		return &amqp.OkResponse{Body: body}, nil
	}

}

// AuthorizeHandler handles authorizing a user given their username and password
func AuthorizeHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var authorizeDto AuthorizeDto

		if err := json.Unmarshal(msg, &authorizeDto); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"authorizeDto": authorizeDto})
		}

		// find user from given username
		repo := NewRepository(s.Cassandra)
		userRecord, err := repo.FindByUsername(authorizeDto.Username)
		if err != nil {
			return nil, &amqp.ErrorResponse{Message: "Not found", Status: http.StatusNotFound}
		}

		// compare password against hash
		if err := userRecord.compareHashAndPassword(authorizeDto.Password); err != nil {
			return nil, &amqp.ErrorResponse{Message: "Invalid password provided", Status: http.StatusUnprocessableEntity}
		}

		// return new accessToken and refreshToken from record
		accessToken, err := auth.GenerateToken(authorizeDto.Username, s.Config.HMACSecret)
		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"authorizeDto": authorizeDto})
		}

		authorized := AuthorizeResponse{
			RefreshToken: userRecord.RefreshToken,
			AccessToken:  accessToken,
		}

		body, _ := json.Marshal(authorized)

		return &amqp.OkResponse{Body: body}, nil
	}
}

// ExistsHandler verifies a user exists in the system
func ExistsHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		var existsDto ExistsDto

		if err := json.Unmarshal(msg, &existsDto); err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"existsDto": existsDto})
		}

		repo := NewRepository(s.Cassandra)

		exists, err := repo.Exists(existsDto.Username)

		if err != nil {
			return amqp.HandleInternalServiceError(err, map[string]interface{}{"username": existsDto.Username})
		}

		if exists == false {
			return nil, &amqp.ErrorResponse{Message: "User not found", Status: http.StatusNotFound}
		}

		existsResponse := ExistsResponse{
			Exists: true,
		}

		body, _ := json.Marshal(existsResponse)

		return &amqp.OkResponse{Body: body}, nil
	}
}
