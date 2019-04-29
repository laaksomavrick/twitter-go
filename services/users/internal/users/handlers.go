package users

import (
	"encoding/json"
	"fmt"
	"twitter-go/services/users/internal/core"
)

func CreateHandler(u *core.Users) func([]byte) interface{} {
	return func(msg []byte) interface{} {
		var dto CreateUserDto

		if err := json.Unmarshal(msg, &dto); err != nil {
			//TODO-1: err handling?
			panic(err)
		}

		var foo string

		_ = u.Cassandra.Session.Query("SELECT now() FROM system.local;").Scan(&foo)

		fmt.Println(foo)

		return dto
	}
}
