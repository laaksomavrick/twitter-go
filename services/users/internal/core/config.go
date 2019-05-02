package core

import (
	"os"
	"twitter-go/services/common/config"
)

// UsersConfig defines the shape of the configuration values used in the users service
type UsersConfig struct {
	*config.Config
	CassandraURL      string
	CassandraKeyspace string
	HmacSecret        string
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

	// TODO-4: move cassandra env to another common config

	if os.Getenv("CASSANDRA_URL") == "" {
		os.Setenv("CASSANDRA_URL", "127.0.0.1")
	}

	if os.Getenv("CASSANDRA_KEYSPACE") == "" {
		os.Setenv("CASSANDRA_KEYSPACE", "twtr")
	}

	if os.Getenv("HMAC_SECRET") == "" {
		// TODO-5: real hmac secret; read from file
		// if keyData, e := ioutil.ReadFile("test/hmacTestKey"); e == nil {
		// 	hmacSampleSecret = keyData
		// } else {
		// 	panic(e)
		// }
		os.Setenv("HMAC_SECRET", "hmacsecret")
	}

	return &UsersConfig{
		Config: &config.Config{
			Env:      os.Getenv("GO_ENV"),
			Port:     os.Getenv("PORT"),
			AmqpURL:  os.Getenv("AMQP_URL"),
			AmqpPort: os.Getenv("AMQP_PORT"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
		CassandraURL:      os.Getenv("CASSANDRA_URL"),
		CassandraKeyspace: os.Getenv("CASSANDRA_KEYSPACE"),
		HmacSecret:        os.Getenv("HMAC_SECRET"),
	}
}
