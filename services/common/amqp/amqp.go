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

const exchange = "twtr"

// TODO-6: handle disconnects from rmqp

// Client wraps common amqp operations
type Client struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	replyToDelivery <-chan amqp.Delivery
}

// OkResponse represents the expected type for an amqp reply
type OkResponse struct {
	Body []byte
}

// ErrorResponse represents the shape of an amqp response error
type ErrorResponse struct {
	Message string
	Status  int
}

func HandleInternalServiceError(err error, data interface{}) (*OkResponse, *ErrorResponse) {
	logger.Error(logger.Loggable{
		Message: err.Error(),
		Data:    data,
	})
	return nil, &ErrorResponse{Message: "Internal server error", Status: http.StatusInternalServerError}
}

// NewClient constructs a new instance of a client
func NewClient(url string) (*Client, error) {
	dial := fmt.Sprintf("%s", url)
	conn, err := amqp.Dial(dial)

	if err != nil {
		logger.Error(logger.Loggable{
			Message: err.Error(),
			Data: map[string]interface{}{
				"url": url,
			},
		})
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

// DirectRequest send a direct reply message to the given routingKey,
// receiving and returning a response
func (client *Client) DirectRequest(routingKey string, routingKeyValues []string, payload interface{}) (*OkResponse, *ErrorResponse) {
	bytes, err := json.Marshal(payload)
	interpolatedRoutingKey := interpolateRoutingKey(routingKey, routingKeyValues)

	if err != nil {
		logger.Error(logger.Loggable{
			Message: "An error occurred parsing the payload to a byte array",
			Data: map[string]interface{}{
				"payload": payload,
			},
		})
		return nil, &ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	corrID := uuid.New().String()

	err = client.channel.Publish(
		exchange,               // exchange
		interpolatedRoutingKey, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       "amq.rabbitmq.reply-to",
			Body:          bytes,
		})

	if err != nil {
		logger.Error(logger.Loggable{
			Message: "Failed to publish a message",
			Data: map[string]interface{}{
				"routingKey": interpolatedRoutingKey,
				"body":       payload,
			},
		})
		return nil, &ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	for d := range client.replyToDelivery {
		if corrID == d.CorrelationId {
			var errorResponse ErrorResponse
			_ = json.Unmarshal(d.Body, &errorResponse)
			// if replyToDelivery is a valid error message
			if errorResponse.Status != 0 {
				return nil, &errorResponse
			}
			okResponse := &OkResponse{Body: d.Body}
			return okResponse, nil
		}
	}

	return nil, nil
}

// DirectReply applies a given function as a callback on a given routingKey for processing,
// directly replying with the result
func (client *Client) DirectReply(routingKey string, callback func([]byte) (*OkResponse, *ErrorResponse)) {
	client.declareExchange(routingKey)

	q, err := client.declareQueue(routingKey)
	if err != nil {
		return
	}

	err = client.bindQueue(routingKey, q)
	if err != nil {
		return
	}

	err = client.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		logger.Error(logger.Loggable{
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
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		logger.Error(logger.Loggable{
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
			okResponse, errorResponse := callback(d.Body)

			var bytes []byte

			if okResponse != nil {
				bytes = okResponse.Body
			} else if errorResponse != nil {
				bytes, err = json.Marshal(errorResponse)
			} else {
				logger.Error(logger.Loggable{
					Message: "An error occurred retrieving a response from an amqp replyFunc",
					Data: map[string]interface{}{
						"okResponse":    okResponse,
						"errorResponse": errorResponse,
					},
				})
			}

			if err != nil {
				logger.Error(logger.Loggable{
					Message: "An error occurred parsing the payload to a byte array",
					Data: map[string]interface{}{
						"okResponse":    okResponse,
						"errorResponse": errorResponse,
					},
				})
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
					Message: "An error occurred publishing a message",
					Data: map[string]interface{}{
						"routingKey":    routingKey,
						"queueName":     q.Name,
						"okResponse":    okResponse,
						"errorResponse": errorResponse,
					},
				})
				return
			}

			d.Ack(false)
		}
	}()
}

// PublishToTopic publishes a message to the given routing key on a topic exchange
func (client *Client) PublishToTopic(routingKey string, keyValues []string, payload interface{}) *ErrorResponse {
	bytes, err := json.Marshal(payload)
	interpolatedRoutingKey := interpolateRoutingKey(routingKey, keyValues)

	if err != nil {
		logger.Error(logger.Loggable{
			Message: "An error occurred parsing the payload to a byte array",
			Data: map[string]interface{}{
				"payload": payload,
			},
		})
		return &ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	client.declareExchange(routingKey)

	err = client.channel.Publish(
		exchange,               // exchange
		interpolatedRoutingKey, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})

	if err != nil {
		logger.Error(logger.Loggable{
			Message: "Failed to publish a message",
			Data: map[string]interface{}{
				"routingKey": interpolatedRoutingKey,
				"body":       payload,
			},
		})
		return &ErrorResponse{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	return nil

}

// ConsumeFromTopic calls a callback for all messages sent to a given routingKey on a topic exchange
func (client *Client) ConsumeFromTopic(routingKey string, callback func([]byte)) {
	client.declareExchange(routingKey)

	q, err := client.declareQueue(routingKey)
	if err != nil {
		return
	}

	err = client.bindQueue(routingKey, q)
	if err != nil {
		return
	}

	msgs, err := client.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			callback(d.Body)
		}
	}()
}

func (client *Client) declareExchange(routingKey string) {
	err := client.channel.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Message: "An error occurred declaring an exchange",
			Data: map[string]interface{}{
				"routingKey": routingKey,
			},
		})
	}
}

func (client *Client) declareQueue(routingKey string) (*amqp.Queue, error) {
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
			Message: "An error occurred initializing a queue",
			Data: map[string]interface{}{
				"routingKey": routingKey,
			},
		})
		return nil, err
	}

	return &q, nil
}

func (client *Client) bindQueue(routingKey string, q *amqp.Queue) error {
	err := client.channel.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)

	if err != nil {
		logger.Error(logger.Loggable{
			Message: "An error occurred binding a queue",
			Data: map[string]interface{}{
				"routingKey": routingKey,
			},
		})
		return err
	}
	return nil
}
