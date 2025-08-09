# core/context

Minimal request context for correlation and tenancy. Dependency-light, transport-agnostic, and safe to adopt anywhere.

## Features
- TraceID and RequestID for correlation
- UserID, TenantID, SessionID helpers
- Labels map for small, sanitized tags
- StartTime + Duration helper
- Safe logging fields via `Fields(ctx)`
- Canonical header constants for simple propagation

## Install
```bash
go get core/context
```

## Quick start
```go
import (
	"context"
	ctxpkg "core/context"
)

func handler(parent context.Context) {
	ctx, rc := ctxpkg.New(parent)
	ctx = ctxpkg.WithTenant(ctx, "tenant-123")
	ctx = ctxpkg.WithUser(ctx, "user-42")
	ctx = ctxpkg.WithRequestID(ctx, "req-abc")
	ctx = ctxpkg.WithTrace(ctx, "trace-xyz")
	ctx = ctxpkg.WithLabel(ctx, "feature", "on")

	_ = ctxpkg.Validate(rc)
	logFields := ctxpkg.Fields(ctx)
	_ = logFields // pass to your logger
}
```

## API
```go
// New/From/Into
New(parent context.Context, opts ...Option) (context.Context, *RequestContext)
From(ctx context.Context) (*RequestContext, bool)
Into(ctx context.Context, rc *RequestContext) context.Context

// Enrichers
WithTrace(ctx context.Context, traceID string) context.Context
WithRequestID(ctx context.Context, requestID string) context.Context
WithUser(ctx context.Context, userID string) context.Context
WithTenant(ctx context.Context, tenantID string) context.Context
WithSession(ctx context.Context, sessionID string) context.Context
WithLabel(ctx context.Context, key, value string) context.Context

// Accessors
TenantID(ctx context.Context) (string, bool)
UserID(ctx context.Context) (string, bool)

// Utilities
Duration(ctx context.Context) time.Duration
Fields(ctx context.Context) map[string]any
Validate(rc *RequestContext) error
```

## Headers
```go
const (
	HeaderRequestID = "X-Request-Id"
	HeaderTraceID   = "X-Trace-Id"
	HeaderUserID    = "X-User-Id"
	HeaderTenantID  = "X-Tenant-Id"
	HeaderSessionID = "X-Session-Id"
)
```

## Design
- Unexported typed key prevents collisions
- No external deps; ID generation is out-of-scope for v1
- Transport-agnostic; HTTP/gRPC helpers can be added as subpackages later

## Testing
- Table-driven tests recommended; construct contexts with `New` and assert via `From` 