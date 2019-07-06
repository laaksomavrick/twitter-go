package main

import (
	"log"
	"os"
	"time"
	"twitter-go/services/common/cassandra"

	"github.com/streadway/amqp"
)

func main() {
	var amqpConn *amqp.Connection
	var cassandraClient *cassandra.Client

	amqpURL := os.Getenv("AMQP_URL")
	cassandraURL := os.Getenv("CASSANDRA_URL")
	cassandraKeyspace := os.Getenv("CASSANDRA_KEYSPACE")

	if amqpURL == "" {
		log.Fatalf("AMQP_URL not supplied, got %s", amqpURL)
	}

	if cassandraURL == "" {
		log.Fatalf("CASSANDRA_URL not supplied, got %s", cassandraURL)
	}

	if cassandraKeyspace == "" {
		log.Fatalf("CASSANDRA_KEYSPACE not supplied, got %s", cassandraKeyspace)
	}

	log.Printf("AMQP_URL: %s", amqpURL)
	log.Printf("CASSANDRA_URL: %s", cassandraURL)

	// Give ourselves a minute (10 * 6s) - should we make this configurable?
	for i := 0; i < 10; i++ {
		if amqpConn == nil {
			log.Print("Attempting to dial rabbit")
			amqpConn, _ = amqp.Dial(amqpURL)
			if amqpConn != nil {
				log.Print("Rabbit OK")
			}
		}

		if cassandraClient == nil {
			log.Print("Attempting to connect to cassandra")
			cassandraClient, _ = cassandra.NewClient(cassandraURL, cassandraKeyspace)
			if cassandraClient != nil {
				log.Print("Cassandra OK")
			}
		}

		if amqpConn != nil && cassandraClient != nil {
			log.Print("Connected to rabbitmq and cassandra, exiting")
			os.Exit(0)
		}

		log.Printf("Sleeping...")

		time.Sleep(6 * time.Second)
	}

}
