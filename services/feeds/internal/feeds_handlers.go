package internal

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

func GetMyFeedHandler(s *service.Service) func([]byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
	return func(msg []byte) (*amqp.OkResponse, *amqp.ErrorResponse) {
		return &amqp.OkResponse{Body: make([]byte, 0)}, nil
	}
}
