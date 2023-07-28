package volta

import "encoding/json"

type Config struct {
	// RabbitMQ's connection string
	RabbitMQ string

	// Request / RequestJSON timeout in seconds
	Timeout int

	// Retry count for connecting to RabbitMQ
	ConnectRetries int

	// Retry interval for connecting to RabbitMQ
	ConnectRetryInterval int

	// JSON Marshaler
	Marshal func(interface{}) ([]byte, error)

	// JSON Unmarshaler
	Unmarshal func([]byte, interface{}) error

	// Disable logging
	DisableLogging bool
}

var DefaultConfig = &Config{
	RabbitMQ:             "amqp://guest:guest@localhost:5672/",
	Timeout:              10,
	ConnectRetries:       5,
	ConnectRetryInterval: 10,
	Marshal:              json.Marshal,
	Unmarshal:            json.Unmarshal,
	DisableLogging:       false,
}
