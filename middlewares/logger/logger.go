package logger

import (
	"fmt"
	"github.com/volta-dev/volta"
	"time"
)

func New() volta.Handler {
	return func(ctx *volta.Ctx) error {
		timeStart := time.Now()

		ctx.Next()

		elapsed := time.Since(timeStart)

		fmt.Printf(
			"[%s] correlationId: \"%s\" | exchange: \"%s\" | routingKey: \"%s\" | elapsed: \"%s\"\n",
			time.Now().Format("2006-01-02 15:04:05"),
			ctx.CorrelationId(),
			ctx.Delivery.Exchange,
			ctx.Delivery.RoutingKey,
			elapsed,
		)

		return nil
	}
}
