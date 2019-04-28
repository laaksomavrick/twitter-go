package main

import (
	"fmt"
	"log"
	"twitter-go/services/common/amqp"
)

func main() {
	log.Print("Users started")
	client, err := amqp.NewClient("amqp://rabbitmq:rabbitmq@localhost", "5672")
	if err != nil {
		log.Fatalf("%s", err)
	}
	client.RPCReply("rpc_queue", func(msg []byte) interface{} {
		fmt.Print("replying")
		return map[string]interface{}{
			"hello": "from rpc :)",
		}
	})
}
