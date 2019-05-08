package tweets

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

		fmt.Println(r.Context().Value("username"))

		json.NewEncoder(w).Encode(map[string]string{
			"hello": "1",
		})
	}
}
