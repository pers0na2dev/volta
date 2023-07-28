package volta

import (
	"github.com/rabbitmq/amqp091-go"
	"testing"
)

func TestApp_OnMessage(t *testing.T) {
	app := New(Config{
		RabbitMQ:       "amqp://volta:volta@localhost:5672/",
		DisableLogging: true,
	})

	app.OnMessage(func(delivery amqp091.Delivery) {
		t.Log("OnMessage")
	})

	if len(app.onMessage) != 1 {
		t.Errorf("App.OnMessage() error = %v", app.onMessage)
	}

	app.OnMessage(func(delivery amqp091.Delivery) {
		t.Log("OnMessage")
	})

	if len(app.onMessage) != 2 {
		t.Errorf("App.OnMessage() error = %v", app.onMessage)
	}
}
