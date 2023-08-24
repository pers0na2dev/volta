package main

import (
	"encoding/json"
	"fmt"
	"github.com/volta-dev/volta"
	"time"
)

func main() {
	app := volta.New(volta.Config{
		RabbitMQ:             "amqp://auka:auka@localhost:5672/",
		Timeout:              10,
		Marshal:              json.Marshal,
		Unmarshal:            json.Unmarshal,
		ConnectRetries:       5,
		ConnectRetryInterval: 10,
	})

	app.AddExchanges(
		volta.Exchange{Name: "testing", Type: "topic"},
	)

	app.AddQueue(
		volta.Queue{Name: "testing.1", RoutingKey: "testing.1", Exchange: "testing"},
		volta.Queue{Name: "testing.*", RoutingKey: "testing.*", Exchange: "testing"},
	)

	app.AddConsumer("testing.1", Handler)
	app.AddConsumer("testing.*", Handler3)

	go func() {
		if err := app.Listen(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(1 * time.Second)

	err := app.Publish("testing.1", "testing", []byte("123"))
	if err != nil {
		panic(err)
	}

	select {}
}

func Handler(c *volta.Ctx) error {
	fmt.Println("FromHandler1:", c.RoutingKey())
	return c.Reply([]byte("123"))
}

func Handler3(c *volta.Ctx) error {
	fmt.Println("FromHandler3:", c.RoutingKey())
	return c.Ack(true)
}
