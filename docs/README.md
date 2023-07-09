---
description: >-
  Volta is a handy library for quickly building asynchronys microservices using
  RabbitMQ in Golang.
---

# ❤ Welcome

### Installation

```bash
go get github.com/volta-dev/volta
```

### Hello, World!

```go
package main

import (
    "encoding/json"
    "github.com/volta-dev/volta"
)

func main() {
    app := volta.New(volta.Config{
        RabbitMQ:  "amqp://guest:guest@localhost:5672/",
    })
    
    app.AddExchanges(volta.Exchange{Name: "test", Type: "topic"})
    app.AddQueue(volta.Queue{Name: "test", RoutingKey: "test", Exchange: "test"})
    
    app.AddConsumer("test", Handler)
    
    app.Listen()
}

func Handler(ctx *volta.Ctx) error {
    return ctx.Reply([]body("Hello, World!"))
}
```
