# ⚡ Make Volta Faster

### Custom JSON Encoder/Decoder <a href="#custom-json-encoderdecoder" id="custom-json-encoderdecoder"></a>

Since we use **encoding/json** as default json library due to stability and producibility. However, the standard library is a bit slow compared to 3rd party libraries. If you're not happy with the performance of **encoding/json**, we recommend you to use these libraries:

* [goccy/go-json](https://github.com/goccy/go-json)
* [bytedance/sonic](https://github.com/bytedance/sonic)
* [segmentio/encoding](https://github.com/segmentio/encoding)
* [mailru/easyjson](https://github.com/mailru/easyjson)
* [minio/simdjson-go](https://github.com/minio/simdjson-go)
* [wI2L/jettison](https://github.com/wI2L/jettison)

Example

```
package main

import "github.com/volta-dev/volta"
import "github.com/bytedance/sonic"

func main() {
    app := volta.New(volta.Config{
        ...
        Marshal: sonic.Marshal,
        Unmarshal: sonic.Unmarshal,
        ...
    })

    # ...
}
```
