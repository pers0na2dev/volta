package volta

import (
	"testing"
	"time"
)

func TestApp_Publish(t *testing.T) {
	app := New(Config{
		RabbitMQ:       "amqp://volta:volta@localhost:5672/",
		DisableLogging: true,
	})

	app.AddExchanges(Exchange{Name: "test", Type: "topic"})
	app.AddQueue(Queue{Name: "test", Exchange: "test", Durable: true})
	app.AddConsumer("test", func(ctx *Ctx) error {
		if string(ctx.Body()) != "test" {
			t.Errorf("Body is %s, expected test", ctx.Body())
		}

		ctx.Ack(false)

		return nil
	})

	go func(app *App) {
		if err := app.Listen(); err != nil {
			t.Errorf("App.Listen() error = %v", err)
		}
	}(app)

	time.Sleep(1 * time.Second)

	err := app.Publish("test", "test", []byte("test"))
	if err != nil {
		t.Errorf("App.Publish() error = %v", err)
	}

	time.Sleep(5 * time.Second)

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}
