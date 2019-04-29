package core

import (
	"os"
	"twitter-go/services/common/config"
)

// UsersConfig defines the shape of the configuration values used in the users service
type UsersConfig struct {
	*config.Config
}

// NewConfig returns the default configuration values used across the api
func NewConfig() *UsersConfig {
	// set defaults - these can be overwritten via command line

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

	return &UsersConfig{
		Config: &config.Config{
			Env:      os.Getenv("GO_ENV"),
			Port:     os.Getenv("PORT"),
			AmqpURL:  os.Getenv("AMQP_URL"),
			AmqpPort: os.Getenv("AMQP_PORT"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
	}
}
