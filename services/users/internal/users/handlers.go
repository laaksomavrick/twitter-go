package users

import (
	"encoding/json"
	"twitter-go/services/common/auth"
	"twitter-go/services/users/internal/core"
)

func CreateHandler(u *core.Users) func([]byte) interface{} {
	return func(msg []byte) interface{} {
		var user User

		if err := json.Unmarshal(msg, &user); err != nil {
			//TODO-1: err handling?
			// log.Fatal(err)
			return nil
		}

		if err := user.prepareForInsert(); err != nil {
			//TODO-1: err handling?
			// log.Fatal(err)
			return nil
		}

		repo := NewUsersRepository(u.Cassandra)
		if err := repo.Insert(user); err != nil {
			//TODO-1: err handling?
			// log.Fatal(err)
			return nil
		}

		accessToken, err := auth.GenerateToken(user.Username, u.Config.HmacSecret)
		if err != nil {
			//TODO-1: err handling?
			// log.Fatal(err)
			return nil
		}
		user.AccessToken = accessToken

		user.sanitize()

		return user
	}

}
