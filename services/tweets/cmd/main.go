package main

import (
	"twitter-go/services/common/amqp"
	"twitter-go/services/common/cassandra"
	"twitter-go/services/common/logger"
	"twitter-go/services/tweets/internal/core"
	"twitter-go/services/tweets/internal/tweets"
)

func main() {

	logger.Init()

	config := core.NewConfig()

	amqp, err := amqp.NewClient(config.AmqpURL, config.AmqpPort)

	if err != nil {
		panic(err)
	}

	cassandra, err := cassandra.NewClient(config.CassandraURL, config.CassandraKeyspace)

	if err != nil {
		panic(err)
	}

	service := core.NewTweetsService(amqp, cassandra, config)

	service.Init(tweets.Repliers)
}
