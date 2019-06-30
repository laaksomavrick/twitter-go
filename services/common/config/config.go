package config

import (
	"os"
)

// Config defines the shape of the configuration values used across the api
type Config struct {
	Env      string
	Port     string
	AmqpURL  string
	AmqpPort string
	LogLevel string
}

// NewConfig returns the default configuration values used across the api
func NewConfig() *Config {
	// set defaults - these can be overwritten via command line

	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "development")
	}

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "3000")
	}

	if os.Getenv("AMQP_URL") == "" {
		os.Setenv("AMQP_URL", "amqp://rabbitmq:rabbitmq@localhost:5672")
	}

	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "debug")
	}

	return &Config{
		Env:      os.Getenv("GO_ENV"),
		Port:     os.Getenv("PORT"),
		AmqpURL:  os.Getenv("AMQP_URL"),
		AmqpPort: os.Getenv("AMQP_PORT"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
