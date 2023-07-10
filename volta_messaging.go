package volta

import (
	"context"
	"encoding/xml"
	"github.com/rabbitmq/amqp091-go"
	"math/rand"
	"time"
)

func (m *App) AddConsumer(routingKey string, handlers ...Handler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.handlers == nil {
		m.handlers = make(map[string][]Handler)
	}

	m.handlers[routingKey] = handlers
}

// Consume consumes messages from the queue with the given routing key.
// Handlers are the functions that will be executed when a message is received.
// Handlers are executed in the order they are passed.
func (m *App) consume(routingKey string, handlers ...Handler) error {
	connection, err := amqp091.Dial(m.config.RabbitMQ)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	messages, err := channel.Consume(routingKey, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	handlersWithMiddlewares := make([]Handler, 0)
	handlersWithMiddlewares = append(handlersWithMiddlewares, m.middlewares...)
	handlersWithMiddlewares = append(handlersWithMiddlewares, handlers...)

	go func() {
		for message := range messages {
			go func(msg amqp091.Delivery, h []Handler, channel *amqp091.Channel) {
				h[0](&Ctx{App: m, Delivery: msg, handlers: h, Channel: channel})
			}(message, handlersWithMiddlewares, channel)
		}
	}()

	return nil
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
func (m *App) Publish(name, exchange string, body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	connection, err := amqp091.Dial(m.config.RabbitMQ)
	if err != nil {
		return err
	}

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

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
func (m *App) Request(name string, body []byte) (data []byte, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	connection, err := amqp091.Dial(m.config.RabbitMQ)
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

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
		"",
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
func (m *App) PublishJSON(name, exchange string, body interface{}) error {
	data, err := m.config.Marshal(body)
	if err != nil {
		return err
	}

	return m.Publish(name, exchange, data)
}

// RequestJSON publishes a message to the exchange with the given name, body will be marshaled to JSON
// and exchange is the exchange name.
// Waits for response.
// Unmarshals response to response interface.
func (m *App) RequestJSON(name string, body interface{}, response interface{}) error {
	data, err := m.config.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := m.Request(name, data)
	if err != nil {
		return err
	}

	err = m.config.Unmarshal(resp, &response)
	if err != nil {
		return err
	}

	return nil
}

func (m *App) PublishXML(name, exchange string, body interface{}) error {
	data, err := xml.Marshal(body)
	if err != nil {
		return err
	}

	return m.Publish(name, exchange, data)
}

func (m *App) RequestXML(name string, body interface{}, response interface{}) error {
	data, err := xml.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := m.Request(name, data)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(resp, &response)
	if err != nil {
		return err
	}

	return nil
}
