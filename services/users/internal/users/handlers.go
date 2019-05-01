package users

import (
	"encoding/json"
	"fmt"
	"log"
	"twitter-go/services/users/internal/core"
)

func CreateHandler(u *core.Users) func([]byte) interface{} {
	return func(msg []byte) interface{} {
		var user User

		if err := json.Unmarshal(msg, &user); err != nil {
			//TODO-1: err handling?
			log.Fatal(err)
		}

		// bcrypt password

		// set refresh token

		// insert
		repo := NewUsersRepository(u.Cassandra)
		err := repo.Insert(user)
		if err != nil {
			//TODO-1: err handling?
			log.Fatal(err)
		}

		// set access token (expiry, etc)

		var foo string

		_ = u.Cassandra.Session.Query("SELECT now() FROM system.local;").Scan(&foo)

		fmt.Println(foo)

		return user
	}
}
