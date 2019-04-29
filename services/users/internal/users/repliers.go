package users

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/users/internal/core"
)

var Repliers = core.Repliers{
	core.Replier{
		RoutingKey: amqp.CreateUserKey,
		Handler:    CreateHandler,
	},
}
