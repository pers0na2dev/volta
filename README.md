# 🐰 volta
❤️ A handy library for working with RabbitMQ 🐰 inspired by Express.js and Martini-like code style.

#### Features
- [x] Middlewares
- [x] Automatic Reconnect with retry limit/timeout
- [ ] OnMessage/OnStartup/etc hooks
- [x] JSON Request / JSON Bind
- [x] XML Request / XML Bind
- [ ] Automatic Dead Lettering <on error / timeout>
- [x] Set of ready-made middleware (limitter / request logger)

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
        ConnectRetries:       5,
        ConnectRetryInterval: 10,
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

    if err := app.Listen(); err != nil {
        panic(err)
    }
}

func Handler(ctx *volta.Ctx) error {
    return ctx.Ack(false)
}

```

### 📝 License

This project is licensed under the WTFPL - see the [LICENSE](LICENSE) file for details

### 🤝 Contributing

Feel free to open an issue or create a pull request.