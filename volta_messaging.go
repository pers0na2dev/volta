package volta

import (
	"context"
	"encoding/xml"
	"math/rand"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func (a *App) AddConsumer(routingKey string, handlers ...Handler) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.handlers == nil {
		a.handlers = make(map[string][]Handler)
	}

	a.handlers[routingKey] = handlers
}

// JSONConsumer is a helper function that creates a handler that will unmarshal the request body to the given type.
func JSONConsumer[Data any](callback func(ctx *Ctx, body Data) error) Handler {
	return func(ctx *Ctx) error {
		var body Data
		if err := ctx.BindJSON(&body); err != nil {
			return ctx.App.onBindError(ctx, err)
		}
		return callback(ctx, body)
	}
}

// XMLConsumer is a helper function that creates a handler that will unmarshal the request body to the given type.
func XMLConsumer[Data any](callback func(ctx *Ctx, body Data) error) Handler {
	return func(ctx *Ctx) error {
		var body Data
		if err := ctx.BindXML(&body); err != nil {
			return ctx.App.onBindError(ctx, err)
		}
		return callback(ctx, body)
	}
}

// Consume consumes messages from the queue with the given routing key.
// Handlers are the functions that will be executed when a message is received.
// Handlers are executed in the order they are passed.
func (a *App) consume(routingKey string, handlers ...Handler) error {
	connection, err := amqp091.Dial(a.config.RabbitMQ)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	messages, err := channel.Consume(routingKey, randomString(12), false, false, false, false, nil)
	if err != nil {
		return err
	}

	handlersWithMiddlewares := make([]Handler, 0)
	handlersWithMiddlewares = append(handlersWithMiddlewares, a.middlewares...)
	handlersWithMiddlewares = append(handlersWithMiddlewares, handlers...)

	go func() {
		for message := range messages {
			go func(msg amqp091.Delivery, h []Handler, channel *amqp091.Channel) {
				h[0](&Ctx{App: a, Delivery: msg, handlers: h, Channel: channel})
			}(message, handlersWithMiddlewares, channel)
		}

		defer connection.Close()
		defer channel.Close()
	}()

	return nil
}

// ConsumeNative consumes messages from the specified routing key using the AMQP 0.9.1 protocol.
// It returns a channel of message deliveries and an error if any occurred.
func (a *App) ConsumeNative(routingKey string) (<-chan amqp091.Delivery, error) {
	connection, err := amqp091.Dial(a.config.RabbitMQ)
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	messages, err := channel.Consume(routingKey, randomString(12), false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// Publish publishes a message to the exchange with the given name, body is the message body
// and exchange is the exchange name.
// No wait for response.
// Returns error if something went wrong.
func (a *App) Publish(name, exchange string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	connection, err := amqp091.Dial(a.config.RabbitMQ)
	if err != nil {
		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.PublishWithContext(
		ctx,
		exchange,
		name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

// Request publishes a message to the exchange with the given name, body is the message body
// and exchange is the exchange name.
// Waits for response.
func (a *App) Request(name, exchange string, body []byte) (data []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Timeout)*time.Second)
	defer cancel()

	connection, err := amqp091.Dial(a.config.RabbitMQ)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	corrId := randomString(32)

	err = channel.PublishWithContext(
		ctx,
		exchange,
		name,
		false,
		false,
		amqp091.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          body,
		})
	if err != nil {
		return nil, err
	}

	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for message := range messages {
		if message.CorrelationId == corrId {
			return message.Body, nil
		}
	}

	return nil, nil
}

// PublishJSON publishes a message to the exchange with the given name, body will be marshaled to JSON
// and exchange is the exchange name.
// No wait for response.
func (a *App) PublishJSON(name, exchange string, body interface{}) error {
	data, err := a.config.Marshal(body)
	if err != nil {
		return err
	}

	return a.Publish(name, exchange, data)
}

// RequestJSON publishes a message to the exchange with the given name, body will be marshaled to JSON
// and exchange is the exchange name.
// Waits for response.
// Unmarshals response to response interface.
func (a *App) RequestJSON(name, exchange string, body interface{}, response interface{}) error {
	data, err := a.config.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := a.Request(name, exchange, data)
	if err != nil {
		return err
	}

	return a.config.Unmarshal(resp, &response)
}

func (a *App) PublishXML(name, exchange string, body interface{}) error {
	data, err := xml.Marshal(body)
	if err != nil {
		return err
	}

	return a.Publish(name, exchange, data)
}

func (a *App) RequestXML(name, exchange string, body interface{}, response interface{}) error {
	data, err := xml.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := a.Request(name, exchange, data)
	if err != nil {
		return err
	}

	return xml.Unmarshal(resp, &response)
}
