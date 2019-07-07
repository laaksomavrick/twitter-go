package followers

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/gateway/internal/core"
	"twitter-go/services/gateway/internal/users"
)

// FollowUserHandler provides a HandlerFunc for following a user
func FollowUserHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtUsername := core.GetUsernameFromRequest(r)
		followUserDto := &FollowUserDto{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(followUserDto); err != nil {
			core.Error(w, http.StatusBadRequest, core.BadRequest)
			return
		}

		followUserDto.Username = jwtUsername

		if err := followUserDto.Validate(); err != nil {
			core.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		existsUserDto := &users.ExistsUserDto{
			Username: followUserDto.FollowingUsername,
		}

		// Check that the user exists
		_, errorResponse := s.Amqp.DirectRequest(amqp.ExistsUserKey, []string{followUserDto.Username}, existsUserDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		// Follow the user
		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.FollowUserKey, []string{jwtUsername}, followUserDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		relationship := make(map[string]interface{})

		if err := json.Unmarshal(okResponse.Body, &relationship); err != nil {
			core.Error(w, http.StatusUnprocessableEntity, core.UnprocessableEntity)
			return
		}

		core.Ok(w, relationship)
	}
}
