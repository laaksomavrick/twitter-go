package users

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/users/internal/core"
)

// Repliers maps routing keys to handlers
var Repliers = core.Repliers{
	core.Replier{
		RoutingKey: amqp.CreateUserKey,
		Handler:    CreateHandler,
	},
	core.Replier{
		RoutingKey: amqp.AuthorizeUserKey,
		Handler:    AuthorizeHandler,
	},
}
