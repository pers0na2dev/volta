package volta

import "encoding/json"

type Config struct {
	// RabbitMQ's connection string
	RabbitMQ string

	// Request / RequestJSON timeout in seconds
	Timeout int

	// JSON Marshaler
	Marshal func(interface{}) ([]byte, error)

	// JSON Unmarshaler
	Unmarshal func([]byte, interface{}) error
}

var DefaultConfig = &Config{
	RabbitMQ:  "amqp://guest:guest@localhost:5672/",
	Timeout:   10,
	Marshal:   json.Marshal,
	Unmarshal: json.Unmarshal,
}
