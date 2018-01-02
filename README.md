# ctrlc

CTRL-C is a Go library that provides an easy way of having a task that
is context-aware and deals with SIGINT and SIGTERM signals.

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/caarlos0/ctrlc"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Backgroud(), time.Second)
    defer cancel()
    err := ctrlc.Default.Run(ctx, func() error {
        // this is a task that doe something
        return nil
    })
    // will err if context times out, if the task returns an error or if
    // a SIGTERM or SIGINT is received (CTRL-C for example).
    if err != nil {
        log.Fatalln(err)
    }
}
```

