package internal

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

// Repliers maps routing keys to handlers
var Repliers = service.Repliers{
	service.Replier{
		RoutingKey: amqp.CreateTweetKey,
		Handler:    CreateHandler,
	},
	service.Replier{
		RoutingKey: amqp.GetAllUserTweetsKey,
		Handler:    GetAllHandler,
	},
}
