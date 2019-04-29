package main

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/users/internal/core"
)

func main() {
	config := core.NewConfig()

	amqp, err := amqp.NewClient(config.AmqpURL, config.AmqpPort)
	if err != nil {
		panic(err)
	}

	users := core.NewUsers(amqp, config)

	users.Init()

	forever := make(chan bool)
	<-forever
}
