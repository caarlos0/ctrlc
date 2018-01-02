# ctrlc

CTRL-C is a Go library that provides an easy way to handle
interrups and context timeouts and cancelations on your cli.

## Usage

```go
package main

import "context"
import "github.com/caarlos0/ctrlc"

func main() {
    ctx, cancel := context.WithTimeout(context.Backgroud(), time.Second)
    defer cancel()
    ctrlc.Default.Run(ctx, func() error {
        // do something
        return nil
    })
}
```

