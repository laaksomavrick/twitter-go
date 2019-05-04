package users

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/auth"
	"twitter-go/services/users/internal/core"
)

// CreateHandler handles creating a user record
func CreateHandler(u *core.Users) func([]byte) interface{} {
	return func(msg []byte) interface{} {
		var user User

		if err := json.Unmarshal(msg, &user); err != nil {
			return amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		if err := user.prepareForInsert(); err != nil {
			return amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		repo := NewUsersRepository(u.Cassandra)
		if err := repo.Insert(user); err != nil {
			return err
		}

		accessToken, err := auth.GenerateToken(user.Username, u.Config.HmacSecret)
		if err != nil {
			return amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		user.AccessToken = accessToken

		user.sanitize()

		return user
	}

}

// AuthorizeHandler handles authorizing a user given their username and password
func AuthorizeHandler(u *core.Users) func([]byte) interface{} {
	return func(msg []byte) interface{} {

		var authorizeDto AuthorizeDto

		if err := json.Unmarshal(msg, &authorizeDto); err != nil {
			return amqp.RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
		}

		// find user from given username
		repo := NewUsersRepository(u.Cassandra)
		userRecord, amqpErr := repo.FindByUsername(authorizeDto.Username)
		if amqpErr != nil {
			return amqpErr
		}

		// compare password against hash
		if err := userRecord.compareHashAndPassword(authorizeDto.Password); err != nil {
			return amqp.RPCError{Message: "Invalid password provided", Status: http.StatusBadRequest}
		}

		// return new accessToken and refreshToken from record
		accessToken, err := auth.GenerateToken(authorizeDto.Username, u.Config.HmacSecret)
		if err != nil {
			return amqp.RPCError{Message: "Something went wrong", Status: http.StatusInternalServerError}
		}

		authorized := AuthorizeResponse{
			RefreshToken: userRecord.RefreshToken,
			AccessToken:  accessToken,
		}

		return authorized
	}
}
