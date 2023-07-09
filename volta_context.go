package volta

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

type Ctx struct {
	App           *App
	Delivery      amqp091.Delivery
	Channel       *amqp091.Channel
	handlers      []Handler
	handlerCursor int
}

func (ctx *Ctx) Reply(data []byte) error {
	replyCtx, cancel := context.WithTimeout(context.Background(), time.Duration(ctx.App.config.Timeout)*time.Second)
	defer cancel()

	ctx.Channel.PublishWithContext(
		replyCtx,
		"",
		ctx.Delivery.ReplyTo,
		false,
		false,
		amqp091.Publishing{
			CorrelationId: ctx.Delivery.CorrelationId,
			Body:          data,
		},
	)

	return ctx.Ack(false)
}

func (ctx *Ctx) ReplyJSON(data interface{}) error {
	replyCtx, cancel := context.WithTimeout(context.Background(), time.Duration(ctx.App.config.Timeout)*time.Second)
	defer cancel()

	jsonData, err := ctx.App.config.Marshal(data)
	if err != nil {
		return err
	}

	ctx.Channel.PublishWithContext(
		replyCtx,
		"",
		ctx.Delivery.ReplyTo,
		false,
		false,
		amqp091.Publishing{
			ContentType:   "application/json",
			CorrelationId: ctx.Delivery.CorrelationId,
			Body:          jsonData,
		},
	)

	return ctx.Ack(false)
}

func (ctx *Ctx) Next() error {
	ctx.handlerCursor++
	if ctx.handlerCursor < len(ctx.handlers) {
		ctx.handlers[ctx.handlerCursor](ctx)
	}

	return nil
}

func (ctx *Ctx) Bind(data interface{}) error {
	err := ctx.App.config.Unmarshal(ctx.Delivery.Body, data)
	if err != nil {
		return err
	}

	return nil
}
func (ctx *Ctx) Ack(multiple bool) error {
	return ctx.Delivery.Ack(multiple)
}

func (ctx *Ctx) Nack(multiple, requeue bool) error {
	return ctx.Delivery.Nack(multiple, requeue)
}

func (ctx *Ctx) Reject(requeue bool) error {
	return ctx.Delivery.Reject(requeue)
}

func (ctx *Ctx) Body() []byte {
	return ctx.Delivery.Body
}

func (ctx *Ctx) ContentType() string {
	return ctx.Delivery.ContentType
}

func (ctx *Ctx) CorrelationId() string {
	return ctx.Delivery.CorrelationId
}

func (ctx *Ctx) ReplyTo() string {
	return ctx.Delivery.ReplyTo
}

func (ctx *Ctx) MessageId() string {
	return ctx.Delivery.MessageId
}

func (ctx *Ctx) Timestamp() time.Time {
	return ctx.Delivery.Timestamp
}

func (ctx *Ctx) Type() string {
	return ctx.Delivery.Type
}

func (ctx *Ctx) UserId() string {
	return ctx.Delivery.UserId
}

func (ctx *Ctx) AppId() string {
	return ctx.Delivery.AppId
}

func (ctx *Ctx) ConsumerTag() string {
	return ctx.Delivery.ConsumerTag
}

func (ctx *Ctx) MessageCount() uint32 {
	return ctx.Delivery.MessageCount
}

func (ctx *Ctx) DeliveryTag() uint64 {
	return ctx.Delivery.DeliveryTag
}

func (ctx *Ctx) Redelivered() bool {
	return ctx.Delivery.Redelivered
}