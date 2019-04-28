package core

import (
	"os"
)

// Config defines the shape of the configuration values used across the api
type Config struct {
	Env        string
	Port       string
	AmqpURL    string
	AmqpPort   string
	LogLevel   string
	HmacSecret []byte
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
		os.Setenv("AMQP_URL", "amqp://rabbitmq:rabbitmq@localhost")
	}

	if os.Getenv("AMQP_PORT") == "" {
		os.Setenv("AMQP_PORT", "5672")
	}

	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LOG_LEVEL", "debug")
	}

	if os.Getenv("HMAC_SECRET") == "" {
		// TODO: real hmac secret; read from file
		// TODO: should only need to read this once on app start, not every req
		// -> put hmacSecret in config
		// if keyData, e := ioutil.ReadFile("test/hmacTestKey"); e == nil {
		// 	hmacSampleSecret = keyData
		// } else {
		// 	panic(e)
		// }
		os.Setenv("HMAC_SECRET", "hmacsecret")
	}

	return &Config{
		Env:        os.Getenv("GO_ENV"),
		Port:       os.Getenv("PORT"),
		AmqpURL:    os.Getenv("AMQP_URL"),
		AmqpPort:   os.Getenv("AMQP_PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
		HmacSecret: []byte(os.Getenv("HMAC_SECRET")),
	}
}
