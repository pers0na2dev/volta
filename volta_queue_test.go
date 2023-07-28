package volta

import "testing"

func TestApp_AddQueue(t *testing.T) {
	app := New(Config{
		DisableLogging: true,
		RabbitMQ:       "amqp://volta:volta@localhost:5672/",
	})

	app.AddQueue(Queue{Name: "test", Exchange: "test"})

	if len(app.queues) == 0 {
		t.Error("No exchanges added")
	}

	if app.queues["test"].Name != "test" {
		t.Errorf("Exchange name is %s, expected test", app.exchanges["test"].Name)
	}

	if app.queues["test"].Exchange != "test" {
		t.Errorf("Exchange type is %s, expected topic", app.exchanges["test"].Type)
	}
}

func TestApp_declareQueue(t *testing.T) {
	app := New(Config{
		DisableLogging: true,
		RabbitMQ:       "amqp://volta:volta@localhost:5672/",
	})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.declareExchange(Exchange{Name: "test", Type: "topic"}); err != nil {
		t.Errorf("App.declareExchange() error = %v", err)
	}

	if err := app.declareQueue(Queue{Name: "test", Exchange: "test", Durable: true}); err != nil {
		t.Errorf("App.declareQueue() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_PurgeQueue(t *testing.T) {
	app := New(Config{
		DisableLogging: true,
		RabbitMQ:       "amqp://volta:volta@localhost:5672/",
	})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.declareExchange(Exchange{Name: "test", Type: "topic"}); err != nil {
		t.Errorf("App.declareExchange() error = %v", err)
	}

	if err := app.declareQueue(Queue{Name: "test", Exchange: "test", Durable: true}); err != nil {
		t.Errorf("App.declareQueue() error = %v", err)
	}

	if err := app.PurgeQueue("test", true); err != nil {
		t.Errorf("App.PurgeQueue() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}
