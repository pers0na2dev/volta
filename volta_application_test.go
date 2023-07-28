package volta

import (
	"testing"
	"time"
)

func TestApp_Listen(t *testing.T) {
	app := New(Config{
		RabbitMQ:             "amqp://volta:volta@localhost:5672/",
		ConnectRetryInterval: 0,
		ConnectRetries:       0,
	})

	// TEST: Listen() should start a new connection to RabbitMQ
	go func() {
		if err := app.Listen(); err != nil {
			t.Errorf("App.Listen() error = %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_Close(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	go func() {
		if err := app.Listen(); err != nil {
			t.Errorf("App.Listen() error = %v", err)
		}
	}()

	time.Sleep(5 * time.Second)

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_connect(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_initExchanges(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	app.AddExchanges(Exchange{Name: "test", Type: "topic"})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.initExchanges(); err != nil {
		t.Errorf("App.initExchanges() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_initQueues(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	app.AddExchanges(Exchange{Name: "test", Type: "topic"})
	app.AddQueue(Queue{Name: "test", Durable: true, Exchange: "test"})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.initExchanges(); err != nil {
		t.Errorf("App.initExchanges() error = %v", err)
	}

	if err := app.initQueues(); err != nil {
		t.Errorf("App.initQueues() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_initConsumers(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	app.AddExchanges(Exchange{Name: "test", Type: "topic"})
	app.AddQueue(Queue{Name: "test", Durable: true, Exchange: "test"})
	app.AddConsumer("test", func(msg *Ctx) error {
		return msg.Ack(false)
	})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.initExchanges(); err != nil {
		t.Errorf("App.initExchanges() error = %v", err)
	}

	if err := app.initQueues(); err != nil {
		t.Errorf("App.initQueues() error = %v", err)
	}

	if err := app.initConsumers(); err != nil {
		t.Errorf("App.initConsumers() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_Use(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	app.Use(func(ctx *Ctx) error {
		return ctx.Ack(false)
	})

	if app.middlewares == nil {
		t.Errorf("App.Use() error = %v", "middlewares is nil")
	}

	if len(app.middlewares) < 1 {
		t.Errorf("App.Use() error = %v", "middlewares is empty")
	}
}
