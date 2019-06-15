package followers

import (
	"encoding/json"
	"errors"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/gateway/internal/core"
	"twitter-go/services/gateway/internal/users"
)

func FollowUserHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtUsername := core.GetUsernameFromRequest(r)
		followUserDto := &FollowUserDto{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(followUserDto); err != nil {
			core.EncodeJSONError(w, errors.New(core.UnprocessableEntity), http.StatusBadRequest)
			return
		}

		followUserDto.Username = jwtUsername

		if errs := followUserDto.Validate(); len(errs) > 0 {
			core.EncodeJSONErrors(w, errs, http.StatusUnprocessableEntity)
			return
		}

		existsUserDto := &users.ExistsUserDto{
			Username: followUserDto.FollowingUsername,
		}

		// Check that the user exists
		_, errorResponse := s.Amqp.DirectRequest(amqp.ExistsUserKey, []string{followUserDto.Username}, existsUserDto)

		if errorResponse != nil {
			core.EncodeJSONError(w, errors.New(errorResponse.Message), errorResponse.Status)
			return
		}

		// Follow the user
		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.FollowUserKey, []string{jwtUsername}, followUserDto)

		if errorResponse != nil {
			core.EncodeJSONError(w, errors.New(errorResponse.Message), errorResponse.Status)
			return
		}

		relationship := make(map[string]interface{})

		if err := json.Unmarshal(okResponse.Body, &relationship); err != nil {
			core.EncodeJSONError(w, errors.New(core.UnprocessableEntity), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(relationship)
	}
}
