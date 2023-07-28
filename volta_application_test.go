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
