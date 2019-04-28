package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"twitter-go/services/gateway/internal/core"

	"twitter-go/services/common/amqp"
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

		client, err := amqp.NewClient("amqp://rabbitmq:rabbitmq@localhost:5672")
		if err != nil {
			log.Fatalf("%s", err)
		}

		res, err := client.SendRPC("rpc_queue", map[string]interface{}{"number": 3})
		if err != nil {
			log.Fatalf("%s", err)
		}

		fmt.Println(res)

		fmtRes := make(map[string]interface{})
		if err := json.Unmarshal(res, &fmtRes); err != nil {
			log.Fatalf("%s", err)
		}

		json.NewEncoder(w).Encode(fmtRes)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
