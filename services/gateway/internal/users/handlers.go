package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/gateway/internal/core"
)

// CreateHandler handles creating a new user.
func CreateHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createUserDto := &CreateUserDto{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(createUserDto); err != nil {
			core.EncodeJSONError(w, errors.New(core.UnprocessableEntity), http.StatusBadRequest)
			return
		}

		if errs := createUserDto.Validate(); len(errs) > 0 {
			core.EncodeJSONErrors(w, errs, http.StatusBadRequest)
			return
		}

		res, err := s.Amqp.RPCRequest(amqp.CreateUserKey, createUserDto)

		if err != nil {
			core.EncodeJSONError(w, errors.New(err.Message), err.Status)
			return
		}

		user := make(map[string]interface{})
		if err := json.Unmarshal(res, &user); err != nil {
			handleError(
				w,
				err,
				"Users.CreateHandler",
				"An error occurred processing an rpc response",
				http.StatusInternalServerError,
			)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

// AuthorizeHandler handles the username and password authorization flow
func AuthorizeHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authenticateUserDto := &AuthenticateUserDto{}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(authenticateUserDto); err != nil {
			core.EncodeJSONError(w, errors.New(core.UnprocessableEntity), http.StatusBadRequest)
			return
		}

		if errs := authenticateUserDto.Validate(); len(errs) > 0 {
			core.EncodeJSONErrors(w, errs, http.StatusBadRequest)
			return
		}

		res, err := s.Amqp.RPCRequest(amqp.AuthorizeUserKey, authenticateUserDto)

		if err != nil {
			core.EncodeJSONError(w, errors.New(err.Message), err.Status)
			return
		}

		authorization := make(map[string]interface{})
		if err := json.Unmarshal(res, &authorization); err != nil {
			handleError(
				w,
				err,
				"Users.AuthorizeHandler",
				"An error occurred processing an rpc response",
				http.StatusInternalServerError,
			)
			return
		}

		json.NewEncoder(w).Encode(authorization)

	}
}

func handleError(w http.ResponseWriter, err error, caller string, msg string, status int) {
	logger.Error(logger.Loggable{
		Caller:  caller,
		Message: msg,
		Data: map[string]interface{}{
			"error": err.Error(),
		},
	})
	core.EncodeJSONError(w, err, status)
}
