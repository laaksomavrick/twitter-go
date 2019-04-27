package users

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/gateway/internal/core"
)

// CreateHandler handles creating a new user.
func CreateHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createUserDto := &CreateUserDto{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(createUserDto); err != nil {
			panic(err)
		}

		if errs := createUserDto.Validate(); len(errs) > 0 {
			core.EncodeJSONErrors(w, errs, http.StatusBadRequest)
			return
		}

		// TODO: amqp to the presently non-existent user service :)

		json.NewEncoder(w).Encode(map[string]string{
			"hello": "world",
		})
	}
}
