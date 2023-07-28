package volta

import "testing"

func TestApp_AddExchanges(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	app.AddExchanges(Exchange{Name: "test", Type: "topic"})

	if len(app.exchanges) == 0 {
		t.Error("No exchanges added")
	}

	if app.exchanges["test"].Name != "test" {
		t.Errorf("Exchange name is %s, expected test", app.exchanges["test"].Name)
	}

	if app.exchanges["test"].Type != "topic" {
		t.Errorf("Exchange type is %s, expected topic", app.exchanges["test"].Type)
	}
}

func TestApp_declareExchange(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.declareExchange(Exchange{Name: "test", Type: "topic"}); err != nil {
		t.Errorf("App.declareExchange() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}

func TestApp_PurgeExchange(t *testing.T) {
	app := New(Config{
		RabbitMQ: "amqp://volta:volta@localhost:5672/",
	})

	if err := app.connect(); err != nil {
		t.Errorf("App.connect() error = %v", err)
	}

	if err := app.declareExchange(Exchange{Name: "test", Type: "topic"}); err != nil {
		t.Errorf("App.declareExchange() error = %v", err)
	}

	if err := app.PurgeExchange("test", true); err != nil {
		t.Errorf("App.PurgeExchange() error = %v", err)
	}

	if err := app.Close(); err != nil {
		t.Errorf("App.Close() error = %v", err)
	}
}
