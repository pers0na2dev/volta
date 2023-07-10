package limiter

import (
	"github.com/volta-dev/volta"
	"time"
)

type Config struct {
	// Limit is the maximum number of requests allowed in a minute
	Limits int

	// Next is a function to skip middleware based on some condition
	Next func(c *volta.Ctx) bool

	// OnLimitReached is a function that will be called when the limit is reached
	OnLimitReached volta.Handler

	// Requests is a map of routingKey addresses and their request count
	requests map[string]int
}

func New(config Config) volta.Handler {
	config.requests = make(map[string]int)

	return func(c *volta.Ctx) (err error) {
		if config.Next != nil && config.Next(c) {
			return c.Next()
		}

		if count, ok := config.requests[c.Delivery.RoutingKey]; ok {
			if count >= config.Limits {
				if config.OnLimitReached != nil {
					return config.OnLimitReached(c)
				}

				return c.ReplyJSON(volta.Map{
					"message": "Limit reached",
				})
			}

			config.requests[c.Delivery.RoutingKey] = count + 1
		} else {
			config.requests[c.Delivery.RoutingKey] = 1
		}

		// Wait for 1 minute before resetting the request count
		time.AfterFunc(time.Minute, func() {
			delete(config.requests, c.Delivery.RoutingKey)
		})

		// Continue processing the next middleware or route handler
		return c.Next()
	}
}
