# core/cache

Minimal in-memory cache with TTL and GetOrCompute helper.

## Features
- Concurrent-safe in-memory cache
- Optional default TTL and periodic cleanup
- `GetOrCompute` to populate on demand
- No external dependencies

## Install
```bash
go get core/cache
```

## Quick start
```go
import (
	"context"
	"time"
	"core/cache"
)

func main() {
	c := cache.NewMemory(cache.WithDefaultTTL(5*time.Minute), cache.WithCleanupInterval(time.Minute))
	defer c.Close()

	// Set / Get
	c.Set("k", 123, 0)
	v, ok := c.Get("k")
	_ = v; _ = ok

	// GetOrCompute
	val, err := c.GetOrCompute(context.Background(), "user:42", 10*time.Minute, func(ctx context.Context) (any, error) {
		return fetchUser(ctx, "42")
	})
	_ = val; _ = err
}
``` 