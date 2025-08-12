# core/event

Clean, simple pub/sub API with an in-memory implementation.

## Features
- Topic-based publish/subscribe
- Multiple subscribers per topic
- Per-subscription retries
- Error handling hooks
- Context-aware publishing with cancellation
- Type-safe handler helpers
- Header support for metadata
- Concurrency-safe

## Install
```bash
go get core/event
```

## Quick start
```go
import (
	"context"
	"core/event"
)

func main() {
	bus := events.NewMemoryBus(
		events.WithBuffer(128),
		events.WithWorkers(2),
	)
	defer bus.Close()

	// Basic subscription
	sub, _ := bus.Subscribe("user.created", func(ctx context.Context, evt any) error {
		// Access headers with events.HeadersFrom(ctx)
		return nil
	}, events.WithRetries(3))
	defer sub.Unsubscribe()

	// Type-safe subscription
	events.SubscribeTyped[string](bus, "user.created", func(ctx context.Context, userID string) error {
		// userID is automatically type-cast from any
		return nil
	})

	// Publish with headers
	_ = bus.Publish(context.Background(), "user.created", "user-123", 
		events.WithHeaders(map[string]string{"source": "auth-service"}))
}
```

## Architecture

**Single Interface**: `EventBus` provides clean publish/subscribe operations.

**Memory Implementation**: Each topic uses a buffered channel with configurable worker goroutines. Workers snapshot subscribers to avoid lock contention during handler execution.

**Headers**: Metadata is passed through context, accessible via `HeadersFrom(ctx)`.

## Configuration

**Bus options**:
- `WithBuffer(n)`: per-topic buffer size (default 64)
- `WithWorkers(n)`: workers per topic (default 1)  
- `WithOnError(func(...))`: hook for handler failures after retries

**Subscribe options**:
- `WithRetries(n)`: retry attempts per handler (default 1)

**Publish options**:
- `WithHeaders(map[string]string)`: attach metadata headers
- `WithKey(string)`: partition key (for future distributed adapters)

## Guarantees

- **Concurrency**: Handlers run concurrently via topic workers
- **Ordering**: Per-topic FIFO ordering; not per-subscriber
- **Cancellation**: Publish respects context cancellation
- **Clean shutdown**: `Close()` stops all workers and prevents new operations

## Errors

- `events.ErrClosed`: bus has been closed
- `events.ErrNilHandler`: handler cannot be nil

## Testing

Use the in-memory bus for unit tests. Small buffers help test cancellation behavior.

```go
func TestExample(t *testing.T) {
	bus := events.NewMemoryBus(events.WithBuffer(1))
	defer bus.Close()
	// ... test code
}
``` 