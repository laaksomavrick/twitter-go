package amqp

import (
	"encoding/json"
	"errors"
	"twitter-go/services/common/logger"

	"github.com/streadway/amqp"
)

// Client wraps common amqp operations
type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewClient constructs a new instance of a client
func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, errors.New("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, errors.New("Failed to open a channel in RabbitMQ")
	}

	if err != nil {
		return nil, errors.New("Failed to initialize the rpc_queue in RabbitMQ")
	}

	return &Client{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// SendRPC send a direct reply message to the given routingKey,
// receiving and returning a response
func (client *Client) SendRPC(routingKey string, payload map[string]interface{}) (res []byte, err error) {
	msgs, err := client.Channel.Consume(
		"amq.rabbitmq.reply-to", // queue
		"",                      // consumer
		true,                    // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "SendRPC",
			Message: "An error occurred creating a consumer",
			Data:    map[string]interface{}{},
		})
		return res, err
	}

	bytes, err := json.Marshal(payload)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "SendRPC",
			Message: "An error occurred parsing the payload to a byte array",
			Data: map[string]interface{}{
				"payload": payload,
			},
		})
		return res, err
	}

	err = client.Channel.Publish(
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			ReplyTo:     "amq.rabbitmq.reply-to",
			Body:        bytes,
		})

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "SendRPC",
			Message: "Failed to publish a message",
			Data: map[string]interface{}{
				"routingKey": routingKey,
				"body":       payload,
			},
		})
		return res, err
	}

	for d := range msgs {
		res = d.Body
		break
	}

	return res, nil
}

// ConsumeRPC applies a given function as a callback on a given routingKey for processing,
// directly replying with the result
func (client *Client) ConsumeRPC(routingKey string, callback func([]byte) interface{}) {
	q, err := client.Channel.QueueDeclare(
		routingKey, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "ConsumeRPC",
			Message: "An error occurred initializing a queue",
			Data: map[string]interface{}{
				"routingKey": routingKey,
			},
		})
		return
	}

	err = client.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Caller:  "ConsumeRPC",
			Message: "An error occurred setting QoS",
			Data: map[string]interface{}{
				"routingKey": routingKey,
			},
		})
		return
	}

	msgs, err := client.Channel.Consume(
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
			Caller:  "ConsumeRPC",
			Message: "An error occurred registering a consumer",
			Data: map[string]interface{}{
				"routingKey": routingKey,
				"queueName":  q.Name,
			},
		})
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			payload := callback(d.Body)

			bytes, err := json.Marshal(payload)

			if err != nil {
				logger.Error(logger.Loggable{
					Caller:  "ConsumeRPC",
					Message: "An error occurred parsing the payload to a byte array",
					Data: map[string]interface{}{
						"payload": payload,
					},
				})
				return
			}

			err = client.Channel.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        bytes,
				})

			if err != nil {
				logger.Error(logger.Loggable{
					Caller:  "ConsumeRPC",
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

	<-forever
}
