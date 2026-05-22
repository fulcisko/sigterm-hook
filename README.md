# sigterm-hook

Library for graceful shutdown orchestration with dependency-aware teardown ordering.

## Installation

```bash
go get github.com/yourusername/sigterm-hook
```

## Usage

```go
package main

import (
    "context"
    "log"

    "github.com/yourusername/sigterm-hook"
)

func main() {
    h := sighook.New()

    // Register shutdown hooks with optional dependencies
    db := h.Register("database", func(ctx context.Context) error {
        log.Println("closing database connection...")
        return db.Close()
    })

    h.Register("server", func(ctx context.Context) error {
        log.Println("shutting down HTTP server...")
        return server.Shutdown(ctx)
    }, sighook.DependsOn(db)) // server shuts down before database

    // Block until SIGTERM or SIGINT is received, then run hooks in order
    if err := h.Wait(context.Background()); err != nil {
        log.Fatalf("shutdown error: %v", err)
    }

    log.Println("shutdown complete")
}
```

Hooks are executed in reverse dependency order — dependents are torn down before their dependencies, ensuring a clean and predictable shutdown sequence.

## Features

- Dependency-aware teardown ordering
- Context propagation with timeout support
- Listens for `SIGTERM` and `SIGINT` signals
- Zero external dependencies

## License

MIT © [yourusername](https://github.com/yourusername)