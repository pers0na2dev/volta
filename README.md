# 🐰 volta
❤️ A handy library for working with RabbitMQ 🐰 inspired by Express.js and Martini-like code style.

#### Features
- [x] Middlewares
- [x] Automatic Reconnect with retry limit/timeout
- [ ] OnMessage/OnStartup/etc hooks
- [ ] XML Request / XML Bind
- [ ] Graceful shutdown
- [ ] Automatic requeue on error
- [ ] Automatic Dead Lettering <on error / timeout>
- [ ] Custom Logger support
- [ ] Set of ready-made middleware (limitter / request logger)

### 📥 Installation
```bash
go get github.com/volta-dev/volta
```

### 👷 Usage
```go
package main

import (
    "encoding/json"
    "github.com/volta-dev/volta"
)

func main() {
    app := volta.New(volta.Config{
        RabbitMQ:  "amqp://guest:guest@localhost:5672/",
        Timeout:   10,
        Marshal:   json.Marshal,
        Unmarshal: json.Unmarshal,
    })
    
    // Register a exchange "test" with type "topic"
    app.AddExchanges(
        volta.Exchange{Name: "test", Type: "topic"},
    )
    
    // Register a queue "test" with routing key "test" and exchange "test"
    app.AddQueue(
        volta.Queue{Name: "test", RoutingKey: "test", Exchange: "test"},
    )
    
    // Register a handler for the "test" queue
    app.AddConsumer("test", Handler)
    
    app.Listen()
}

func Handler(ctx *volta.Ctx) error {
    return ctx.Ack(false)
}

```