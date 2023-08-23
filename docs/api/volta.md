# ⚡ Volta

### New&#x20;

This method creates a new **App** named instance. You need to pass a **Config** to create instance of applicaton.

{% code title="Signature" lineNumbers="true" %}
```go
func New(config Config) *App
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app := volta.New(volta.Config{
	RabbitMQ:  "amqp://volta:volta@localhost:5672/",
	Timeout:   10,
	Marshal:   json.Marshal,
	Unmarshal: json.Unmarshal,
	ConnectRetries:       5,
	ConnectRetryInterval: 10,
})
```
{% endcode %}

### Config fields

<table><thead><tr><th>Property</th><th>Type</th><th>Description</th><th data-hidden></th></tr></thead><tbody><tr><td><pre><code>RabbitMQ
</code></pre></td><td>string</td><td>URL to connect to RabbitMQ</td><td></td></tr><tr><td><pre><code>Timeout
</code></pre></td><td>int</td><td>Timeout - the time to wait for a response from app.Request / app.RequestJSON.</td><td></td></tr><tr><td><pre><code>Marshal
</code></pre></td><td>func(interface{}) ([]byte, error)</td><td>The function responsible for JSON Marshalling. Defaults: json.Marshal</td><td></td></tr><tr><td><pre><code>Unmarshal
</code></pre></td><td>func([]byte, interface{}) error</td><td>The function responsible for JSON Unmarshalling. Defaults: json.Unmarshal</td><td></td></tr><tr><td><pre><code>ConnectRetries
</code></pre></td><td>int</td><td>Number of reconnection attempts</td><td></td></tr><tr><td><pre><code>ConnectRetryInterval
</code></pre></td><td>int</td><td>Interval between reconnections</td><td></td></tr></tbody></table>


### JSONConsumer&#x20;

**JSONConsumer** is a helper function that creates a handler that will unmarshal the request body to the given type.

{% code title="Signature" lineNumbers="true" %}
```go
func JSONConsumer[Data any](callback func(ctx *Ctx, body Data) error) Handler
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.AddConsumer("testing.12", volta.JSONConsumer[SomeDto](JsonConsumer))

type SomeDto struct {
	Message string `json:"message"`
}

func JsonConsumer(ctx *volta.Ctx, dto SomeDto) error {
	return ctx.ReplyJSON(volta.Map{
		"message": dto.Message,
	})
}
```
{% endcode %}

### XMLConsumer&#x20;

**XMLConsumer** is a helper function that creates a handler that will unmarshal the request body to the given type.

{% code title="Signature" lineNumbers="true" %}
```go
func XMLConsumer[Data any](callback func(ctx *Ctx, body Data) error) Handler
```
{% endcode %}

{% code title="Example" lineNumbers="true" %}
```go
app.AddConsumer("testing.12", volta.XMLConsumer[SomeDto](XmlConsumer))

type SomeDto struct {
    Message string `xml:"message"`
}

func XmlConsumer(ctx *volta.Ctx, dto SomeDto) error {
    return ctx.ReplyXML(volta.Map{
        "message": dto.Message,
    })
}
```
{% endcode %}