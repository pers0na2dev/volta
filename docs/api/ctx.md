# 🧠 Ctx

## Ack

Function to acknowledge a message.

multiple: bool - If true, the server will acknowledge all messages up to and including the message delivered to the client. If false, only the supplied message will be acknowledged.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Ack(multiple bool) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    return ctx.Ack(false)
}
```
{% endcode %}

## Nack

Function to reject a message.

multiple: bool - If true, the server will reject all messages up to and including the message delivered to the client. If false, only the supplied message will be rejected.

requeue: bool - If true, the server will attempt to requeue the message. If false or the requeue fails the messages will be discarded or dead-lettered.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Nack(multiple, requeue bool) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    return ctx.Nack(false, false)
}
```
{% endcode %}

## Reject 

Function to reject a message.

requeue: bool - If true, the server will attempt to requeue the message. If false or the requeue fails the messages will be discarded or dead-lettered.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Reject(requeue bool) error 
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    return ctx.Reject(false)
}
```
{% endcode %}

## BindJSON

Function to bind a message body (JSON-type) to a struct.

body: interface{} - The struct to bind the body to.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) BindJSON(body interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
type User struct {
    Name string `json:"name"`
}

func Handler(ctx *volta.Ctx) error {
    var user User
    if err := ctx.BindJSON(&user); err != nil {
        ...
    }

    return ctx.Reply([]byte("Hello, " + user.Name))
}
```
{% endcode %}

## BindXML

Function to bind a message body (XML-type) to a struct.

body: interface{} - The struct to bind the body to.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) BindXML(body interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
type User struct {
    Name string `xml:"name"`
}

func Handler(ctx *volta.Ctx) error {
    var user User
    if err := ctx.BindXML(&user); err != nil {
        ...
    }

    return ctx.Reply([]byte("Hello, " + user.Name))
}
```
{% endcode %}

## Body

Function to get the message body.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Body() []byte
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    body := ctx.Body()
    ...
}
```
{% endcode %}

## Locals

Function to get the context local storage.

{% code title="Signature"  %}
```go
func (ctx *Ctx) Locals(key string, value ...interface{}) interface{}
```
{% endcode %}

{% code title="Example"  %}
```go
func Middleware(ctx *volta.Ctx) error {
    ctx.Locals("key", "value")
    ...
}

func Handler(ctx *volta.Ctx) error {
    value := ctx.Locals("key").(string)
    ...
}
```

{% endcode %}

## Reply

Function to reply to a message.

body: []byte - The message body to reply with.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Reply(body []byte) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    return ctx.Reply([]byte("Hello, World!"))
}
```
{% endcode %}

## ReplyJSON

Function to reply to a message with automatically json marshal.

body: interface{} - The message body to reply with.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) ReplyJSON(body interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
type User struct {
    Name string `json:"name"`
}

func Handler(ctx *volta.Ctx) error {
    return ctx.ReplyJSON(&User{Name: "John"})
}
```
{% endcode %}

## ReplyXML 

Function to reply to a message with automatically xml marshal.

body: interface{} - The message body to reply with.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) ReplyXML(body interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
type User struct {
    Name string `xml:"name"`
}

func Handler(ctx *volta.Ctx) error {
    return ctx.ReplyXML(&User{Name: "John"})
}
```
{% endcode %}

## ContentType

Function to get the message content type.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) ContentType() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    contentType := ctx.ContentType()
    ...
}
```
{% endcode %}

## CorrelationId

Function to get the message correlation id.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) CorrelationId() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    correlationId := ctx.CorrelationId()
    ...
}
```
{% endcode %}

## ReplyTo

Function to get the message reply to.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) ReplyTo() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    replyTo := ctx.ReplyTo()
    ...
}
```
{% endcode %}

## MessageId

Function to get the message id.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) MessageId() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    messageId := ctx.MessageId()
    ...
}
```
{% endcode %}

## Timestamp

Function to get the message timestamp.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Timestamp() time.Time
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    timestamp := ctx.Timestamp()
    ...
}
```
{% endcode %}

## Type

Function to get the message type.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Type() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    messageType := ctx.Type()
    ...
}
```
{% endcode %}

## UserId

Function to get the message user id.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) UserId() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    userId := ctx.UserId()
    ...
}
```
{% endcode %}

## AppId

Function to get the message app id.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) AppId() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error {
    appId := ctx.AppId()
    ...
}
```
{% endcode %}

## ConsumerTag

Function to get the message consumer tag.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) ConsumerTag() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error { 
    consumerTag := ctx.ConsumerTag()
    ...
}
```
{% endcode %}

## MessageCount

Function to get the message count.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) MessageCount() uint32
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error { 
    messageCount := ctx.MessageCount()
    ...
}
```
{% endcode %}

## DeliveryTag

Function to get the message delivery tag.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) DeliveryTag() uint64
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error { 
    deliveryTag := ctx.DeliveryTag()
    ...
}
```
{% endcode %}

## Redelivered

Function to get the message redelivered.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Redelivered() bool
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error { 
    redelivered := ctx.Redelivered()
    ...
}
```
{% endcode %}

## Exchange

Function to get the message exchange.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) Exchange() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error { 
    exchange := ctx.Exchange()
    ...
}
```
{% endcode %}

## RoutingKey

Function to get the message routing key.

{% code title="Signature" lineNumbers="true" %}
```go
func (ctx *Ctx) RoutingKey() string
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
func Handler(ctx *volta.Ctx) error { 
    routingKey := ctx.RoutingKey()
    ...
}
```
{% endcode %}