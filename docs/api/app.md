# ⚙ App

## AddExchanges

Function to add exchanges to the app.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) AddExchanges(exchange ...Exchange)
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.AddExchanges(
    volta.Exchange{
        Name:       "testing", // Name of the exchange.
        Type:       "fanout", // Type of the exchange.
        Durable:    false, // Durability of the exchange when the server restarts.
        AutoDelete: false, // Auto delete the exchange when there are no more queues bound to it.
        Internal:   false, // If true, messages cannot be published directly to the exchange.
        NoWait:     false, // If true, declare without waiting for a confirmation from the server.
    },
)
```
{% endcode %}

## AddQueue

Function to add queues to the app.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) AddQueue(queue ...Queue)
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.AddQueue(
    volta.Queue{
        Name:       "testing.12", // Name of the queue.
        RoutingKey: "testing.12", // Routing key of the queue.
        Exchange:   "testing", // Exchange to bind the queue to.
        Durable:    false, // Durability of the queue when the server restarts.
        AutoDelete: false, // Auto delete the queue when there are no more consumers subscribed to it.
        Exclusive:  false, // If true, only one consumer can consume from the queue.
        NoWait:     false, // If true, declare without waiting for a confirmation from the server.
    },
)
```
{% endcode %}

## AddConsumer

Function to add consumers to the app.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) AddConsumer(routingKey string, consumer ...Handler)
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.AddConsumer("testing", func(ctx *volta.Ctx) error {
    fmt.Println(ctx.Body())
    return ctx.Ack(false)
})
```
{% endcode %}


## ConsumeNative

ConsumeNative consumes messages from the specified routing key using the AMQP 0.9.1 protocol.
It returns a channel of message deliveries and an error if any occurred.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) ConsumeNative(routingKey string) (<-chan amqp.Delivery, error)
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
messages, err := app.ConsumeNative("testing.12")
if err != nil {
    ...
}

for message := range messages {
    fmt.Println(message.Body)
    message.Ack(false)
}
```
{% endcode %}

## PurgeExchange

Function to purge exchanges from the app.

If force is true, the exchange will be deleted even if it is in use

If force is false, the exchange will be deleted only if it is not in use

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) PurgeExchange(name string, force bool) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if err := app.PurgeExchange("testing", false); err != nil {
    ...
}
```
{% endcode %}

## PurgeQueue

Function to purge queues from the app.

If force is true, the queue will be deleted even if it is in use

If force is false, the queue will be deleted only if it is not in use

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) PurgeQueue(name string, force bool) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if err := app.PurgeQueue("testing.12", false); err != nil {
    ...
}
```
{% endcode %}

## Publish

Function to publish a message to an exchange without response awaiting.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) Publish(name, exchange string, body []byte) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if err := app.Publish("testing.12", "testing", []byte("Hello, World!")); err != nil {
    ...
}
```
{% endcode %}

## PublishJSON

Function to publish a message to an exchange without response awaiting.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) PublishJSON(name, exchange string, body interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if err := app.PublishJSON("testing.12", "testing", volta.Map{"name": "World"}); err != nil {
    ...
}
```
{% endcode %}

## Request

Function to publish a message to an exchange with response awaiting.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) Request(routingKey, exchange string, body []byte) ([]byte, error)
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if response, err := app.Request("testing.12", "testing", []byte("Hello, World!")); err != nil {
    ...
}
```
{% endcode %}

## RequestJSON

Function to publish a message to an exchange with response awaiting.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) RequestJSON(name, exchange string, body interface{}, response interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
var response volta.Map
if response, err := app.RequestJSON(
    "testing.12",
    "testing",
    volta.Map{"name": "World"},
    &response, 
); err != nil {
    ...
}
```
{% endcode %}

## PublishXML

Function to publish a message to an exchange without response awaiting.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) PublishXML(name, exchange string, body interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if err := app.PublishXML("testing.12", "testing", volta.Map{"name": "World"}); err != nil {
    ...
}
```
{% endcode %}

## RequestXML

Function to publish a message to an exchange with response awaiting.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) RequestXML(name, exchange string, body interface{}, response interface{}) error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
var response volta.Map
if response, err := app.RequestXML(
    "testing.12", 
    "testing",
    volta.Map{"name": "World"},
    &response, 
); err != nil {
    ...
}
```
{% endcode %}

## Use

Function to add global middlewares to the app.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) Use(middleware ...Handler)
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.Use(func(ctx *volta.Ctx) error {
    fmt.Println("Before")
    return ctx.Next()
})

app.AddConsumer("testing", func(ctx *volta.Ctx) error {
    fmt.Println("Handler")
    return ctx.Ack(false)
})

if err := app.Listen(); err != nil {
    ...
}

// Output:
// Before
// Handler
```
{% endcode %}

## Listen

Function to initialize all the exchanges and queues and start listening for messages.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) Listen() error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if err := app.Listen(); err != nil {
    ...
}
```
{% endcode %}

## MustListen

Function to initialize all the exchanges and queues and start listening for messages.
Panics if an error occurs.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) MustListen()
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.MustListen()
```
{% endcode %}

## Close    

Function to close the connection to the RabbitMQ server.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) Close() error
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
if err := app.Close(); err != nil {
    ...
}
```
{% endcode %}

## MustClose

Function to close the connection to the RabbitMQ server.
Panics if an error occurs.

{% code title="Signature" lineNumbers="true" %}
```go
func (m *App) MustClose()
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.MustClose()
```
{% endcode %}