---
description: >-
  Limiter is middleware which allows you to limit the number of messages passed
  in the consumer.
---

# Limiter

### Signature

```go
func New(config Config) volta.Handler
```

### Examples

Import the middleware package that is part of the Fiber web framework

```go
import (
  "github.com/volta-dev/volta"
  "github.com/volta-dev/volta/middlewares/limiter"
)
```

After you initiate your Volta app, you can use the following possibilities:

```go
// Initialize
app.Use(limiter.New(limiter.Config{
	Limits: 5, // 5 consumings per minute
	Next: func(c *volta.Ctx) bool {
		if c.Delivery.RoutingKey == "bypass.limit" {
			return true
		}
		
		return false
	},
	OnLimitReached: func(ctx *volta.Ctx) error {
		return ctx.Nack(false, true)
	},
}))
```

### Config

<pre class="language-go"><code class="lang-go">// Config defines the config for middleware.
type Config struct {
<strong>    // Limit is the maximum number of requests allowed in a minute
</strong>    Limits int
    
    // Next is a function to skip middleware based on some condition
    Next func(c *volta.Ctx) bool
    
    // OnLimitReached is a function that will be called when the limit is reached
    OnLimitReached volta.Handler
}
</code></pre>

### Default Config

```go
var ConfigDefault = Config{
    Limits:            0,
    Next:              nil,
    OnLimitReached:    nil,
}
```

