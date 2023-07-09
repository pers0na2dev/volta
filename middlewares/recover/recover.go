package recover

import (
	"fmt"
	"github.com/volta-dev/volta"
)

type Config struct {
	// Next is a function that will be called before executing the middleware.
	Next func(c *volta.Ctx) bool

	// EnableStackTrace is a flag that indicates whether to enable stack trace.
	EnableStackTrace bool

	// StackTraceHandler is a function that will be called when a panic occurs.
	StackTraceHandler func(c *volta.Ctx, e interface{})
}

var ConfigDefault = Config{
	Next:              nil,
	EnableStackTrace:  false,
	StackTraceHandler: defaultStackTraceHandler,
}

func defaultStackTraceHandler(_ *volta.Ctx, e interface{}) {
	fmt.Println("Panic:", e)
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.EnableStackTrace && cfg.StackTraceHandler == nil {
		cfg.StackTraceHandler = defaultStackTraceHandler
	}

	return cfg
}

func New(config ...Config) volta.Handler {
	cfg := configDefault(config...)

	return func(c *volta.Ctx) (err error) {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		defer func() {
			if r := recover(); r != nil {
				if cfg.EnableStackTrace {
					cfg.StackTraceHandler(c, r)
				}

				var ok bool
				if err, ok = r.(error); !ok {
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		return c.Next()
	}
}
