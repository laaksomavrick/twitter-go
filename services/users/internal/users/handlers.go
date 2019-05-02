package users

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/auth"
	"twitter-go/services/users/internal/core"
)

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
