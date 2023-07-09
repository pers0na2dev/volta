package volta

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

type App struct {
	// Configuration
	config Config

	// RabbitMQ connection
	connectRetries int
	baseConnection *amqp091.Connection
	mutex          sync.Mutex

	// Global Middlewares
	middlewares []Handler

	// Exchanges
	exchanges map[string]Exchange

	// Queues
	queues map[string]Queue

	// Handlers
	handlers map[string][]Handler
}

// New creates a new App instance
func New(config Config) *App {
	// Create a new App instance
	app := &App{config: config}

	// Set the configuration to the given one
	if config.RabbitMQ == "" {
		app.config.RabbitMQ = DefaultConfig.RabbitMQ
	}
	if config.Marshal == nil {
		app.config.Marshal = DefaultConfig.Marshal
	}
	if config.Unmarshal == nil {
		app.config.Unmarshal = DefaultConfig.Unmarshal
	}
	if config.Timeout == 0 {
		app.config.Timeout = DefaultConfig.Timeout
	}
	if config.ConnectRetries == 0 {
		app.config.ConnectRetries = DefaultConfig.ConnectRetries
	}
	if config.ConnectRetryInterval == 0 {
		app.config.ConnectRetryInterval = DefaultConfig.ConnectRetryInterval
	}

	return app
}

func (m *App) initExchanges() {
	color.Cyan("\nRegistering exchanges...\n")

	for _, exchange := range m.exchanges {
		err := m.declareExchange(exchange)
		if err != nil {
			panic(fmt.Sprintf("volta: Problem with declaring exchange %s: %s", exchange.Name, err.Error()))
		}

		color.HiWhite("Exchange \"%s\" registered", exchange.Name)
	}
}

func (m *App) initQueues() {
	color.Cyan("\nRegistering queues...\n")

	for _, queue := range m.queues {
		if queue.Exchange != "" {
			err := m.declareQueue(queue)
			if err != nil {
				panic(fmt.Sprintf("volta: Problem with declaring queue %s: %s", queue.Name, err.Error()))
			}

			color.HiWhite("Queue \"%s\" registered", queue.Name)
		} else {
			color.HiRed("Queue \"%s\" skipped (no exchange)", queue.Name)
		}
	}
}

func (m *App) initConsumers() {
	color.Cyan("\nRegistering consumers...\n")
	for rk, handlers := range m.handlers {
		m.consume(rk, handlers...)

		color.HiWhite("Consumer \"%s\" registered", rk)
	}

}

func (m *App) connect() {
	var err error
	m.baseConnection, err = amqp091.Dial(m.config.RabbitMQ)
	if err != nil {
		color.HiRed("volta: Problem with connecting to RabbitMQ: %s", err.Error())
		m.connectRetries++
		if m.connectRetries > m.config.ConnectRetries {
			panic("volta: Problem with connecting to RabbitMQ")
		}

		time.Sleep(time.Duration(m.config.ConnectRetryInterval) * time.Second)

		m.connect()
	}
}

// Listen starts the application, registers the error handler and connects to RabbitMQ
func (m *App) Listen() {
	// Connect to RabbitMQ
	m.connect()

	// Register exchanges
	m.initExchanges()

	// Register queues
	m.initQueues()

	// Register consumers
	m.initConsumers()

	// Check for connection active
	go func() {
		color.HiWhite("\nConnection watcher registered")
		for {
			if m.baseConnection.IsClosed() {
				color.HiRed("Connection to RabbitMQ lost, reconnecting...")

				m.Listen()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// Infinite loop
	forever := make(chan bool)
	<-forever
}

// Close closes the connection to RabbitMQ
func (m *App) Close() {
	err := m.baseConnection.Close()
	if err != nil {
		panic(fmt.Sprintf("volta: Problem with closing the connection to RabbitMQ: %s", err.Error()))
	}
}

func (m *App) Use(middlewares ...Handler) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.middlewares = append(m.middlewares, middlewares...)
}
