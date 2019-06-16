package internal

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/service"
)

// Repliers maps routing keys to handlers
var Repliers = service.Repliers{
	service.Replier{
		RoutingKey: amqp.CreateUserKey,
		Handler:    CreateHandler,
	},
	service.Replier{
		RoutingKey: amqp.AuthorizeUserKey,
		Handler:    AuthorizeHandler,
	},
	service.Replier{
		RoutingKey: amqp.ExistsUserKey,
		Handler:    ExistsHandler,
	},
}
