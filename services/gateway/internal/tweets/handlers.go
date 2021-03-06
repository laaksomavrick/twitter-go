package tweets

import (
	"encoding/json"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/types"
	"twitter-go/services/gateway/internal/core"

	"github.com/gorilla/mux"
)

// CreateTweetHandler provides a HandlerFunc for creating a tweet
func CreateTweetHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createTweetDto := types.CreateTweet{}
		jwtUsername := core.GetUsernameFromRequest(r)

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&createTweetDto); err != nil {
			core.Error(w, http.StatusBadRequest, core.BadRequest)
			return
		}

		if jwtUsername == "" {
			core.Error(w, http.StatusForbidden, core.Forbidden)
			return
		}

		createTweetDto.Username = jwtUsername

		if err := createTweetDto.Validate(); err != nil {
			core.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		logger.Info(logger.Loggable{Message: "Create tweet request", Data: createTweetDto})

		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.CreateTweetKey, []string{createTweetDto.Username}, createTweetDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		tweet := make(map[string]interface{})
		if err := json.Unmarshal(okResponse.Body, &tweet); err != nil {
			core.Error(w, http.StatusInternalServerError, core.InternalServerError)
			return
		}

		core.Ok(w, tweet)
	}
}

// GetAllUserTweetsHandler provides a HandlerFunc for retrieving all tweets made by a user
func GetAllUserTweetsHandler(s *core.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		getAllUserTweetsDto := types.GetAllUserTweets{
			Username: username,
		}

		existsUserDto := types.DoesExist{
			Username: username,
		}

		if err := getAllUserTweetsDto.Validate(); err != nil {
			core.Error(w, http.StatusBadRequest, core.BadRequest)
			return
		}

		// Check that the user exists
		_, errorResponse := s.Amqp.DirectRequest(amqp.ExistsUserKey, []string{getAllUserTweetsDto.Username}, existsUserDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		logger.Info(logger.Loggable{Message: "Get all user tweets request", Data: getAllUserTweetsDto})

		// Get that user's tweets
		okResponse, errorResponse := s.Amqp.DirectRequest(amqp.GetAllUserTweetsKey, []string{getAllUserTweetsDto.Username}, getAllUserTweetsDto)

		if errorResponse != nil {
			core.Error(w, errorResponse.Status, errorResponse.Message)
			return
		}

		tweets := make([]map[string]interface{}, 0)

		if err := json.Unmarshal(okResponse.Body, &tweets); err != nil {
			core.Error(w, http.StatusInternalServerError, core.InternalServerError)
			return
		}

		core.Ok(w, tweets)
	}
}
