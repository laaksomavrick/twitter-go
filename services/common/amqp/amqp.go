package amqp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"twitter-go/services/common/logger"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// TODO-6: handle disconnects from rmqp

// Client wraps common amqp operations
type Client struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	replyToDelivery <-chan amqp.Delivery
}

// RPCError represents the shape of an error from an rpc request
type RPCError struct {
	Message string
	Status  int
}

// NewClient constructs a new instance of a client
func NewClient(url string, port string) (*Client, error) {
	dial := fmt.Sprintf("%s:%s", url, port)
	conn, err := amqp.Dial(dial)

	if err != nil {
		return nil, errors.New("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, errors.New("Failed to open a channel in RabbitMQ")
	}

	delivery, err := ch.Consume(
		"amq.rabbitmq.reply-to", // queue
		"",                      // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)

	if err != nil {
		return nil, errors.New("Failed to create reply-to consumer")
	}

	return &Client{
		conn:            conn,
		channel:         ch,
		replyToDelivery: delivery,
	}, nil
}

// RPCRequest send a direct reply message to the given routingKey,
// receiving and returning a response
func (client *Client) RPCRequest(routingKey string, payload interface{}) (res []byte, rpcError *RPCError) {
	bytes, err := json.Marshal(payload)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "RPCRequest",
			Message: "An error occurred parsing the payload to a byte array",
			Data: map[string]interface{}{
				"payload": payload,
			},
		})
		return res, &RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	corrID := uuid.New().String()

	err = client.channel.Publish(
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       "amq.rabbitmq.reply-to",
			Body:          bytes,
		})

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "RPCRequest",
			Message: "Failed to publish a message",
			Data: map[string]interface{}{
				"routingKey": routingKey,
				"body":       payload,
			},
		})
		return res, &RPCError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	for d := range client.replyToDelivery {
		if corrID == d.CorrelationId {
			var replyError RPCError
			_ = json.Unmarshal(d.Body, &replyError)
			// if replyToDelivery is a valid error message
			if replyError.Status != 0 {
				return res, &replyError
			}
			return d.Body, nil
		}
	}

	return res, nil
}

// RPCReply applies a given function as a callback on a given routingKey for processing,
// directly replying with the result
func (client *Client) RPCReply(routingKey string, callback func([]byte) interface{}) {
	q, err := client.channel.QueueDeclare(
		routingKey, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "RPCReply",
			Message: "An error occurred initializing a queue",
			Data: map[string]interface{}{
				"routingKey": routingKey,
			},
		})
		return
	}

	err = client.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "RPCReply",
			Message: "An error occurred setting QoS",
			Data: map[string]interface{}{
				"routingKey": routingKey,
			},
		})
		return
	}

	msgs, err := client.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "RPCReply",
			Message: "An error occurred registering a consumer",
			Data: map[string]interface{}{
				"routingKey": routingKey,
				"queueName":  q.Name,
			},
		})
		return
	}

	go func() {
		for d := range msgs {
			payload := callback(d.Body)

			bytes, err := json.Marshal(payload)

			if err != nil {
				logger.Error(logger.Loggable{
					Caller:  "RPCReply",
					Message: "An error occurred parsing the payload to a byte array",
					Data: map[string]interface{}{
						"payload": payload,
					},
				})
				return
			}

			err = client.channel.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          bytes,
				})

			if err != nil {
				logger.Error(logger.Loggable{
					Caller:  "RPCReply",
					Message: "An error occurred publishing a message",
					Data: map[string]interface{}{
						"routingKey": routingKey,
						"queueName":  q.Name,
						"payload":    payload,
					},
				})
				return
			}

			d.Ack(false)
		}
	}()
}
