package tweets

import (
	"encoding/json"
	"errors"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/gateway/internal/core"
	"twitter-go/services/gateway/internal/users"

	"github.com/gorilla/mux"
)

// CreateHandler handles creating a new tweet
func CreateHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createTweetDto := &CreateTweetDto{}
		jwtUsername := core.GetUsernameFromRequest(r)

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(createTweetDto); err != nil {
			core.EncodeJSONError(w, errors.New(core.BadRequest), http.StatusBadRequest)
			return
		}

		if jwtUsername == "" {
			core.EncodeJSONError(w, errors.New(core.Forbidden), http.StatusForbidden)
			return
		}

		createTweetDto.Username = jwtUsername

		if errs := createTweetDto.Validate(); len(errs) > 0 {
			core.EncodeJSONErrors(w, errs, http.StatusBadRequest)
			return
		}

		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.CreateTweetKey, []string{createTweetDto.Username}, createTweetDto)

		if errorResponse != nil {
			core.EncodeJSONError(w, errors.New(errorResponse.Message), errorResponse.Status)
			return
		}

		tweet := make(map[string]interface{})
		if err := json.Unmarshal(okResponse.Body, &tweet); err != nil {
			core.EncodeJSONError(w, errors.New(core.InternalServerError), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(tweet)
	}
}

func GetAllUserTweets(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		getAllUserTweetsDto := &GetAllUserTweetsDto{
			Username: username,
		}

		existsUserDto := &users.ExistsUserDto{
			Username: username,
		}

		if errs := getAllUserTweetsDto.Validate(); len(errs) > 0 {
			core.EncodeJSONErrors(w, errs, http.StatusBadRequest)
			return
		}

		// Check that the user exists
		_, errorResponse := s.Amqp.DirectRequest(amqp.ExistsUserKey, []string{getAllUserTweetsDto.Username}, existsUserDto)

		if errorResponse != nil {
			core.EncodeJSONError(w, errors.New(errorResponse.Message), errorResponse.Status)
			return
		}

		// Get that user's tweets
		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.GetAllUserTweetsKey, []string{getAllUserTweetsDto.Username}, getAllUserTweetsDto)

		if errorResponse != nil {
			core.EncodeJSONError(w, errors.New(errorResponse.Message), errorResponse.Status)
			return
		}

		tweets := make([]map[string]interface{}, 0)

		if err := json.Unmarshal(okResponse.Body, &tweets); err != nil {
			core.EncodeJSONError(w, errors.New(core.InternalServerError), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(tweets)
	}
}
