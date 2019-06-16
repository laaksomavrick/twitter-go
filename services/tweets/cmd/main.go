package main

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/config"
	"twitter-go/services/common/logger"
	"twitter-go/services/common/service"
	"twitter-go/services/tweets/internal"
)

func main() {

	logger.Init()

	config := config.NewServiceConfig()

	amqp, err := amqp.NewClient(config.AmqpURL, config.AmqpPort)

	if err != nil {
		panic(err)
	}

	cassandra, err := cassandra.NewClient(config.CassandraURL, config.CassandraKeyspace)

	if err != nil {
		panic(err)
	}

	svc := service.NewService("Tweets", amqp, cassandra, config)

	svc.Init(internal.Repliers, service.Consumers{})
}
