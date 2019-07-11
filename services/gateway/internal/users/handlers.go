package users

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/types"
	"twitter-go/services/gateway/internal/core"
)

// CreateUserHandler provides a HandlerFunc for creating a new user.
func CreateUserHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createUserDto := types.CreateUser{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&createUserDto); err != nil {
			core.Error(w, http.StatusBadRequest, core.BadRequest)
			return
		}

		if err := createUserDto.Validate(); err != nil {
			core.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		logger.Info(logger.Loggable{Message: "Create user request", Data: nil})

		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.CreateUserKey, []string{}, createUserDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		user := make(map[string]interface{})

		if err := json.Unmarshal(okResponse.Body, &user); err != nil {
			core.Error(w, http.StatusUnprocessableEntity, core.UnprocessableEntity)
			return
		}

		core.Ok(w, user)
	}
}

// AuthorizeHandler provides a HandlerFunc for the app authorization flow
func AuthorizeHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticateUserDto := types.AuthenticateUser{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&authenticateUserDto); err != nil {
			core.Error(w, http.StatusBadRequest, core.BadRequest)
			return
		}

		if err := authenticateUserDto.Validate(); err != nil {
			core.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		logger.Info(logger.Loggable{Message: "Authorize user request", Data: nil})

		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.AuthorizeUserKey, []string{authenticateUserDto.Username}, authenticateUserDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		authorization := make(map[string]interface{})
		if err := json.Unmarshal(okResponse.Body, &authorization); err != nil {
			core.Error(w, http.StatusUnprocessableEntity, core.UnprocessableEntity)
			return
		}

		core.Ok(w, authorization)
	}
}
