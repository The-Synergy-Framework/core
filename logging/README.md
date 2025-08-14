# core/logging

Enterprise-grade structured logging built on `log/slog` with optimized context integration, performance, and reliability.

## Features

- **Built on `log/slog`**: Full compatibility with Go's standard logging
- **Context integration**: Automatic injection of fields from `core/context`
- **Performance optimized**: Efficient memory allocation and context field handling
- **Enterprise-ready defaults**: Proper error handling and sensible configurations
- **Async GELF support**: Production-ready Graylog integration with buffering
- **Thread-safe**: Proper synchronization for concurrent usage
- **Structured logging**: Rich attribute support with grouping
- **Package-level convenience**: Direct logging functions for common use cases

## Quick Start

```go
import (
	"context"
	"core/logging"
)

func main() {
	// Set up logger with production defaults
	logger := logging.NewJSON(nil, logging.DefaultConfig()) // writes to stderr
	logging.SetDefault(logger)

	// Use package-level functions
	logging.Info("service starting", "version", "1.0.0", "port", 8080)
	
	// Or use logger instance
	logger.Info("service ready")
}
```

## Context Integration

```go
import (
	"context"
	ctxpkg "core/context"
	"core/logging"
)

func handleRequest(ctx context.Context) {
	// Context fields are automatically included
	ctx = ctxpkg.WithRequestID(ctx, "req-123")
	ctx = ctxpkg.WithTenant(ctx, "tenant-1")
	
	// All context fields will be included in the log
	logging.InfoContext(ctx, "processing request", "endpoint", "/api/users")
	// Output: {"time":"...","level":"INFO","msg":"processing request","request_id":"req-123","tenant":"tenant-1","endpoint":"/api/users"}
}
```

## Configuration

```go
import (
	"log/slog"
	"os"
	"core/logging"
)

func configuredLogger() *logging.Logger {
	config := &logging.Config{
		Level:     slog.LevelDebug,
		AddSource: true, // Include source file:line in logs
	}
	
	return logging.NewJSON(os.Stdout, config)
}
```

## Structured Logging

```go
logger := logging.Default()

// Key-value pairs
logger.Info("user created", 
	"user_id", "123",
	"email", "user@example.com",
	"created_at", time.Now())

// With permanent attributes
userLogger := logger.With("user_id", "123", "session", "abc")
userLogger.Info("login successful")
userLogger.Info("password changed")

// With groups
apiLogger := logger.WithGroup("api")
apiLogger.Info("request", "method", "POST", "path", "/users")
// Output: {"api.method":"POST","api.path":"/users",...}
```

## Production GELF Logging

```go
import (
	"core/logging"
	"core/logging/gelf"
)

func setupGelfLogging() *logging.Logger {
	config := &gelf.Config{
		Level:      slog.LevelInfo,
		Async:      true,              // Non-blocking
		BufferSize: 1000,              // Message buffer size
		Timeout:    5 * time.Second,   // Connection timeout
	}
	
	handler, err := gelf.New("graylog.company.com:12201", config)
	if err != nil {
		panic(err)
	}
	
	return logging.New(handler)
}
```

## Performance Features

### Efficient Context Handling
- Context fields are only processed when logging at enabled levels
- No unnecessary allocations for disabled log levels
- Optimized attribute conversion and grouping

### Async GELF Logging
- Non-blocking UDP writes prevent application slowdown
- Configurable buffer sizes for high-throughput scenarios  
- Graceful degradation when GELF endpoint is unavailable

### Memory Optimization
- Pre-allocated attribute slices reduce GC pressure
- Efficient string formatting and key sanitization
- Minimal allocations in hot logging paths

## API Reference

### Core Types
```go
type Logger struct { ... }
type Config struct {
    Level      slog.Level
    AddSource  bool
    TimeFormat string
}
```

### Constructors
```go
New(handler slog.Handler) *Logger
NewJSON(w io.Writer, config *Config) *Logger
NewText(w io.Writer, config *Config) *Logger
DefaultConfig() *Config
```

### Global Functions
```go
SetDefault(logger *Logger)
Default() *Logger

// Package-level logging (uses default logger)
Debug(msg string, args ...any)
DebugContext(ctx context.Context, msg string, args ...any)
Info(msg string, args ...any)
InfoContext(ctx context.Context, msg string, args ...any)
Warn(msg string, args ...any)
WarnContext(ctx context.Context, msg string, args ...any)
Error(msg string, args ...any)
ErrorContext(ctx context.Context, msg string, args ...any)
```

### Instance Methods
```go
// Enrichment
(*Logger).With(args ...any) *Logger
(*Logger).WithGroup(name string) *Logger

// Logging with context
(*Logger).DebugContext(ctx context.Context, msg string, args ...any)
(*Logger).InfoContext(ctx context.Context, msg string, args ...any)
(*Logger).WarnContext(ctx context.Context, msg string, args ...any)
(*Logger).ErrorContext(ctx context.Context, msg string, args ...any)

// Logging without context (uses context.Background())
(*Logger).Debug(msg string, args ...any)
(*Logger).Info(msg string, args ...any)
(*Logger).Warn(msg string, args ...any)
(*Logger).Error(msg string, args ...any)

// Advanced
(*Logger).LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr)
(*Logger).Enabled(ctx context.Context, level slog.Level) bool
(*Logger).ToSlog() *slog.Logger
```

## Breaking Changes from v1

- **Better defaults**: Default logger writes to stderr instead of being discarded
- **Configuration objects**: Use `Config` structs instead of direct `slog.HandlerOptions`
- **Context methods**: Added `*Context` variants for all logging methods
- **Thread safety**: Proper synchronization for global default logger
- **GELF improvements**: Async processing and better error handling
- **Removed methods**: `WithContext()` method removed in favor of `*Context` methods

## Migration Guide

**Before:**
```go
logger := logging.NewJSON(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
logger.WithContext(ctx).Info("message")
```

**After:**
```go
logger := logging.NewJSON(os.Stderr, &logging.Config{Level: slog.LevelInfo})
logger.InfoContext(ctx, "message")
``` 