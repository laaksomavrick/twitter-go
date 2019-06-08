package config

import (
	"os"
)

// ServiceConfig defines the shape of the configuration values used in backend services
type ServiceConfig struct {
	*Config
	CassandraURL      string
	CassandraKeyspace string
	HMACSecret        string
}

// NewServiceConfig initializes the default service configuration values
func NewServiceConfig() *ServiceConfig {
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

	if os.Getenv("CASSANDRA_URL") == "" {
		os.Setenv("CASSANDRA_URL", "127.0.0.1")
	}

	if os.Getenv("CASSANDRA_KEYSPACE") == "" {
		os.Setenv("CASSANDRA_KEYSPACE", "twtr")
	}

	if os.Getenv("HMAC_SECRET") == "" {
		os.Setenv("HMAC_SECRET", "hmacsecret")
	}

	return &ServiceConfig{
		Config: &Config{
			Env:      os.Getenv("GO_ENV"),
			Port:     os.Getenv("PORT"),
			AmqpURL:  os.Getenv("AMQP_URL"),
			AmqpPort: os.Getenv("AMQP_PORT"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
		CassandraURL:      os.Getenv("CASSANDRA_URL"),
		CassandraKeyspace: os.Getenv("CASSANDRA_KEYSPACE"),
		HMACSecret:        os.Getenv("HMAC_SECRET"),
	}
}
