---
description: >-
  Recover middleware for Volta that recovers from panics anywhere in the stack
  chain and handles the control to the centralized ErrorHandler.
---

# Recover

### Signature

```go
func New(config ...Config) volta.Handler
```

### Examples

Import the middleware package that is part of the Fiber web framework

```go
import (
  "github.com/volta-dev/volta"
  "github.com/volta-dev/volta/middlewares/recover"
)
```

After you initiate your Volta app, you can use the following possibilities:

```go
// Initialize default config
app.Use(recover.New())

// This panic will be caught by the middleware
app.AddConsumer("test", func(ctx *volta.Ctx) error {
    panic("Hardcoded panic")	
})
```

### Config

```go
// Config defines the config for middleware.
type Config struct {
    // Next defines a function to skip this middleware when returned true.
    //
    // Optional. Default: nil
    Next func(c *volta.Ctx) bool

    // EnableStackTrace enables handling stack trace
    //
    // Optional. Default: false
    EnableStackTrace bool

    // StackTraceHandler defines a function to handle stack trace
    //
    // Optional. Default: defaultStackTraceHandler
    StackTraceHandler func(c *volta.Ctx, e interface{})
}
```

### Default Config

```go
func defaultStackTraceHandler(_ *volta.Ctx, e interface{}) {
    fmt.Println("Panic:", e)
}

var ConfigDefault = Config{
    Next:              nil,
    EnableStackTrace:  false,
    StackTraceHandler: defaultStackTraceHandler,
}
```