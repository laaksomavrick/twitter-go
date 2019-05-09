package core

import (
	"os"
	"twitter-go/services/common/config"
)

// TweetsServiceConfig defines the shape of the configuration values used in the tweets service
type TweetsServiceConfig struct {
	*config.Config
	CassandraURL      string
	CassandraKeyspace string
}

// NewConfig returns the default configuration values used in the service
func NewConfig() *TweetsServiceConfig {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "development")
	}

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3000")
	}

	if os.Getenv("AMQP_URL") == "" {
		os.Setenv("AMQP_URL", "amqp://rabbitmq:rabbitmq@localhost")
	}

	if os.Getenv("AMQP_PORT") == "" {
		os.Setenv("AMQP_PORT", "5672")
	}

	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "debug")
	}

	// TODO-4: move cassandra env to another common config

	if os.Getenv("CASSANDRA_URL") == "" {
		os.Setenv("CASSANDRA_URL", "127.0.0.1")
	}

	if os.Getenv("CASSANDRA_KEYSPACE") == "" {
		os.Setenv("CASSANDRA_KEYSPACE", "twtr")
	}

	return &TweetsServiceConfig{
		Config: &config.Config{
			Env:      os.Getenv("GO_ENV"),
			Port:     os.Getenv("PORT"),
			AmqpURL:  os.Getenv("AMQP_URL"),
			AmqpPort: os.Getenv("AMQP_PORT"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
		CassandraURL:      os.Getenv("CASSANDRA_URL"),
		CassandraKeyspace: os.Getenv("CASSANDRA_KEYSPACE"),
	}
}
