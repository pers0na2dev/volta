package volta

import (
	"errors"
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

	// Error handlers
	onBindError OnBindError
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

func (a *App) initExchanges() error {
	if !a.config.DisableLogging {
		color.Cyan("\nRegistering exchanges...\n")
	}

	for _, exchange := range a.exchanges {
		err := a.declareExchange(exchange)
		if err != nil {
			return errors.New(fmt.Sprintf("volta: Problem with declaring exchange %s: %s", exchange.Name, err.Error()))
		}

		if !a.config.DisableLogging {
			color.HiWhite("Exchange \"%s\" registered", exchange.Name)
		}
	}

	return nil
}

func (a *App) initQueues() error {
	if !a.config.DisableLogging {
		color.Cyan("\nRegistering queues...\n")
	}
	for _, queue := range a.queues {
		if queue.Exchange != "" {
			err := a.declareQueue(queue)
			if err != nil {
				return errors.New(fmt.Sprintf("volta: Problem with declaring queue %s: %s", queue.Name, err.Error()))
			}

			if !a.config.DisableLogging {
				color.HiWhite("Queue \"%s\" registered", queue.Name)
			}
		} else {
			if !a.config.DisableLogging {
				color.HiRed("Queue \"%s\" skipped (no exchange)", queue.Name)
			}
		}
	}

	return nil
}

func (a *App) initConsumers() error {
	if !a.config.DisableLogging {
		color.Cyan("\nRegistering consumers...\n")
	}
	for rk, handlers := range a.handlers {
		if err := a.consume(rk, handlers...); err != nil {
			return errors.New(fmt.Sprintf("volta: Problem with consuming queue %s: %s", rk, err.Error()))
		} else {
			if !a.config.DisableLogging {
				color.HiWhite("Consumer \"%s\" registered", rk)
			}
		}
	}

	return nil
}

func (a *App) connect() (err error) {
	if !a.config.DisableLogging {
		color.Cyan("Connecting to RabbitMQ...\n")
	}
	a.baseConnection, err = amqp091.Dial(a.config.RabbitMQ)
	if err != nil {
		if !a.config.DisableLogging {
			color.HiRed("volta: Problem with connecting to RabbitMQ: %s", err.Error())
		}
		a.connectRetries++
		if a.connectRetries > a.config.ConnectRetries {
			return errors.New("volta: Problem with connecting to RabbitMQ")
		}

		time.Sleep(time.Duration(a.config.ConnectRetryInterval) * time.Second)

		connError := a.connect()
		if connError != nil {
			return connError
		}
	}

	return nil
}

// Listen starts the application, registers the error handler and connects to RabbitMQ
func (a *App) Listen() error {
	// Connect to RabbitMQ
	if err := a.connect(); err != nil {
		return err
	}

	// Register exchanges
	if err := a.initExchanges(); err != nil {
		return err
	}

	// Register queues
	if err := a.initQueues(); err != nil {
		return err
	}

	// Register consumers
	if err := a.initConsumers(); err != nil {
		return err
	}

	// Check for connection active
	go func() {
		if !a.config.DisableLogging {
			color.HiWhite("\nConnection watcher registered")
		}
		for {
			if a.baseConnection.IsClosed() {
				if !a.config.DisableLogging {
					color.HiRed("Connection to RabbitMQ lost, reconnecting...")
				}

				a.Listen()
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// Infinite loop
	forever := make(chan bool)
	<-forever

	return nil
}

// MustListen starts the application, registers the error handler and connects to RabbitMQ
// It panics if an error occurs
func (a *App) MustListen() {
	if err := a.Listen(); err != nil {
		panic(err)
	}
}

// Close closes the connection to RabbitMQ
func (a *App) Close() error {
	err := a.baseConnection.Close()
	if err != nil {
		return err
	}

	return nil
}

// MustClose closes the connection to RabbitMQ and panics if an error occurs
func (a *App) MustClose() {
	if err := a.Close(); err != nil {
		panic(err)
	}
}

func (a *App) Use(middlewares ...Handler) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.middlewares = append(a.middlewares, middlewares...)
}
