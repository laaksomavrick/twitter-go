package users

import "encoding/json"

func CreateUser(msg []byte) interface{} {

	var dto CreateUserDto

	if err := json.Unmarshal(msg, &dto); err != nil {
		//TODO-1: err handling?
		panic(err)
	}

	return dto

	// return map[string]interface{}{
	// 	"hello": "from rpc :)",
	// }
}
