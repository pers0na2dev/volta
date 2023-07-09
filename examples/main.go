package main

import (
	"encoding/json"
	"github.com/volta-dev/volta"
)

func main() {
	app := volta.New(volta.Config{
		RabbitMQ:  "amqp://volta:volta@localhost:5672/",
		Timeout:   10,
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	})

	app.AddExchanges(
		volta.Exchange{Name: "testing", Type: "fanout"},
	)

	app.AddQueue(
		volta.Queue{Name: "testing.12", RoutingKey: "testing.12", Exchange: "testing"},
	)

	app.Use(GlobalMiddleware)
	app.AddConsumer("testing.12", Handler)
	
	app.Listen()
}

func GlobalMiddleware(ctx *volta.Ctx) error {
	return ctx.Next()
}

func Handler(ctx *volta.Ctx) error {
	return ctx.ReplyJSON(volta.Map{
		"message": "Hello World!",
	})
}
