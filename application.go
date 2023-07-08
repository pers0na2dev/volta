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
	baseConnection *amqp091.Connection

	// Global Middlewares
	middlewareMutex sync.RWMutex
	middlewares     []Handler

	// Exchanges
	exchangeMutex sync.RWMutex
	exchanges     map[string]Exchange

	// Queues
	queueMutex sync.RWMutex
	queues     map[string]Queue

	// Handlers
	handlerMutex sync.RWMutex
	handlers     map[string][]Handler
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

	return app
}

// Listen starts the application, registers the error handler and connects to RabbitMQ
func (m *App) Listen() {
	// Connect to RabbitMQ
	color.Cyan("Connecting to RabbitMQ...\n")
	err := error(nil)
	m.baseConnection, err = amqp091.Dial(m.config.RabbitMQ)
	if err != nil {
		panic(fmt.Sprintf("volta: Problem with connecting to RabbitMQ: %s", err.Error()))
	}
	color.HiWhite("RabbitMQ: %s\n", m.config.RabbitMQ)

	// Register exchanges
	color.Cyan("\nRegistering exchanges...\n")
	for _, exchange := range m.exchanges {
		err := m.declareExchange(exchange)
		if err != nil {
			panic(fmt.Sprintf("volta: Problem with declaring exchange %s: %s", exchange.Name, err.Error()))
		}

		color.HiWhite("Exchange \"%s\" registered", exchange.Name)
	}

	// Register queues
	color.Cyan("\nRegistering queues...\n")
	for _, queue := range m.queues {
		err := m.declareQueue(queue)
		if err != nil {
			panic(fmt.Sprintf("volta: Problem with declaring queue %s: %s", queue.Name, err.Error()))
		}

		color.HiWhite("Queue \"%s\" registered", queue.Name)
	}

	// Register consumers
	color.Cyan("\nRegistering consumers...\n")
	for rk, handlers := range m.handlers {
		m.consume(rk, handlers...)

		color.HiWhite("Consumer \"%s\" registered", rk)
	}

	// Check for connection active
	go func() {
		color.HiWhite("\nConnection watcher registered")
		for {
			if m.baseConnection.IsClosed() {
				panic("volta: Connection to RabbitMQ is closed")
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
	m.middlewareMutex.Lock()
	defer m.middlewareMutex.Unlock()

	m.middlewares = append(m.middlewares, middlewares...)
}
