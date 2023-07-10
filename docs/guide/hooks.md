---
description: >-
  Hooks - actions that are performed in parallel to certain events. Hooks can
  help you with collecting metrics or other things :)
---

# 🪝 Hooks

Hooks don't interact with middleware in any way, so for example recover middleware won't work here. Be careful not to shoot yourself in the foot.

### OnMessage

Runs in parallel to the start of the message handler from the queue. Do not start Nack/Ack/Reject from a hook.

{% code title="Signature" %}
```go
func (m *App) OnMessage(handler ...OnMessage)
```
{% endcode %}

{% code title="Example" %}
```go
app.OnMessage(func(message amqp091.Delivery) {
    fmt.Println("Message received!")
    
    // DONT RUN Nack/Ack/Reject HERE
})
```
{% endcode %}
