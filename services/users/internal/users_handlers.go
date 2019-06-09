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
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		repo := NewRepository(s.Cassandra)

		exists, err := repo.Exists(user.Username)

		if err != nil {
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		if exists == true {
			return nil, &amqp.ErrorResponse{Message: "User already exists", Status: http.StatusConflict}
		}

		if err := user.prepareForInsert(); err != nil {
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		if err := repo.Insert(user); err != nil {
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		accessToken, err := auth.GenerateToken(user.Username, s.Config.HMACSecret)
		if err != nil {
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
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
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
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
			return nil, &amqp.ErrorResponse{Message: "Something went wrong", Status: http.StatusInternalServerError}
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
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		repo := NewRepository(s.Cassandra)

		exists, err := repo.Exists(existsDto.Username)

		if err != nil {
			return nil, &amqp.ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
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
