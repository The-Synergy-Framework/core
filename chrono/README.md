# core/chrono

Minimal time utilities with a testable Clock interface.

## Features
- `Clock` interface (Now, Since) + `SystemClock`
- Package-level `Default` clock and `SetDefault`
- Helpers: `Now`, `Since`, `IsExpired`, `FormatApprox`

## Install
```bash
go get core/chrono
```

## Quick start
```go
import (
	chrono "core/chrono"
)

func work() {
	start := chrono.Now()
	// ... do things ...
	elapsed := chrono.Since(start)
	fmt.Println("took", chrono.FormatApprox(elapsed))
}
```

## Testability
```go
type fakeClock struct{ t time.Time }
func (f fakeClock) Now() time.Time { return f.t }
func (f fakeClock) Since(t time.Time) time.Duration { return f.t.Sub(t) }

func test() {
	chrono.SetDefault(fakeClock{t: time.Date(2024,1,1,0,0,0,0,time.UTC)})
	defer chrono.SetDefault(nil)
}
```

## API
```go
// Clock
Now() time.Time
Since(time.Time) time.Duration

// Package helpers
Now() time.Time
Since(t time.Time) time.Duration
IsExpired(t time.Time) bool
FormatApprox(d time.Duration) string
``` 