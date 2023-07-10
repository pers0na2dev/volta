package volta

import "github.com/rabbitmq/amqp091-go"

type Handler func(*Ctx) error

type OnMessage func(delivery amqp091.Delivery)

type Map map[string]interface{}

type Exchange struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type Queue struct {
	Name       string
	RoutingKey string
	Exchange   string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}
