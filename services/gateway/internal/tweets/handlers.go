package tweets

import (
	"encoding/json"
	"errors"
	"net/http"
	"twitter-go/services/common/amqp"
	"twitter-go/services/gateway/internal/core"
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

		if errs := createTweetDto.Validate(); len(errs) > 0 {
			core.EncodeJSONErrors(w, errs, http.StatusBadRequest)
			return
		}

		if jwtUsername != createTweetDto.Username {
			core.EncodeJSONError(w, errors.New(core.Forbidden), http.StatusForbidden)
			return
		}

		res, rpcErr := s.Amqp.RPCRequest(amqp.CreateTweetKey, createTweetDto)

		if rpcErr != nil {
			core.EncodeJSONError(w, errors.New(rpcErr.Message), rpcErr.Status)
			return
		}

		tweet := make(map[string]interface{})
		if err := json.Unmarshal(res, &tweet); err != nil {
			core.EncodeJSONError(w, errors.New(core.InternalServerError), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(tweet)
	}
}
