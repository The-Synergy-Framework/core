package ctx

import (
	stdctx "context"
	"errors"
	"strings"
	"time"

	chrono "core/chrono"
)

// contextKey is an unexported type for keys defined in this package to avoid collisions.
// It intentionally does not expose its underlying type.
type contextKey string

const requestContextKey contextKey = "synergy.core.request_context"

// RequestContext carries request-scoped metadata for correlation and tenancy.
//
// It is designed to be lightweight and dependency-free. Do not store large values.
// All fields are optional; use helpers to enrich incrementally through the request path.
type RequestContext struct {
	TraceID   string            // Correlates work across services
	RequestID string            // Identifies this hop/request instance
	UserID    string            // Authenticated user identifier (if any)
	TenantID  string            // Tenant/workspace identifier (if any)
	SessionID string            // Session identifier (if any)
	Labels    map[string]string // Small, sanitized key/value labels
	StartTime time.Time         // Request start time (for duration)
}

// Option mutates a RequestContext during creation.
type Option func(*RequestContext)

// New returns a stdlib context carrying a fresh RequestContext, applying any options.
// If the parent context already has a RequestContext, it is shallow-copied and then options are applied.
func New(parent stdctx.Context, opts ...Option) (stdctx.Context, *RequestContext) {
	rc, _ := From(parent)
	var base RequestContext
	if rc != nil {
		// Shallow copy existing; Labels will be copied below if present
		base = *rc
		if rc.Labels != nil {
			labels := make(map[string]string, len(rc.Labels))
			for k, v := range rc.Labels {
				labels[k] = v
			}
			base.Labels = labels
		}
	}
	if base.StartTime.IsZero() {
		base.StartTime = chrono.Now()
	}
	for _, opt := range opts {
		opt(&base)
	}
	ctx := stdctx.WithValue(parent, requestContextKey, &base)
	return ctx, &base
}

// From extracts the RequestContext from ctx if present.
func From(ctx stdctx.Context) (*RequestContext, bool) {
	v := ctx.Value(requestContextKey)
	if v == nil {
		return nil, false
	}
	rc, ok := v.(*RequestContext)
	return rc, ok && rc != nil
}

// Into attaches the provided RequestContext to ctx.
func Into(ctx stdctx.Context, rc *RequestContext) stdctx.Context {
	if rc == nil {
		return ctx
	}
	return stdctx.WithValue(ctx, requestContextKey, rc)
}

// WithTrace sets the TraceID on the RequestContext inside ctx without overwriting other fields.
func WithTrace(ctx stdctx.Context, traceID string) stdctx.Context {
	if traceID == "" {
		return ctx
	}
	rc, _ := From(ctx)
	if rc == nil {
		ctx, rc = New(ctx)
	}
	rc.TraceID = traceID
	return ctx
}

// WithRequestID sets the RequestID on the RequestContext inside ctx.
func WithRequestID(ctx stdctx.Context, requestID string) stdctx.Context {
	if requestID == "" {
		return ctx
	}
	rc, _ := From(ctx)
	if rc == nil {
		ctx, rc = New(ctx)
	}
	rc.RequestID = requestID
	return ctx
}

// WithUser sets the UserID on the RequestContext inside ctx.
func WithUser(ctx stdctx.Context, userID string) stdctx.Context {
	if userID == "" {
		return ctx
	}
	rc, _ := From(ctx)
	if rc == nil {
		ctx, rc = New(ctx)
	}
	rc.UserID = userID
	return ctx
}

// WithTenant sets the TenantID on the RequestContext inside ctx.
func WithTenant(ctx stdctx.Context, tenantID string) stdctx.Context {
	if tenantID == "" {
		return ctx
	}
	rc, _ := From(ctx)
	if rc == nil {
		ctx, rc = New(ctx)
	}
	rc.TenantID = tenantID
	return ctx
}

// WithSession sets the SessionID on the RequestContext inside ctx.
func WithSession(ctx stdctx.Context, sessionID string) stdctx.Context {
	if sessionID == "" {
		return ctx
	}
	rc, _ := From(ctx)
	if rc == nil {
		ctx, rc = New(ctx)
	}
	rc.SessionID = sessionID
	return ctx
}

// WithLabel sets or replaces a single label on the RequestContext.
func WithLabel(ctx stdctx.Context, key, value string) stdctx.Context {
	key = strings.TrimSpace(key)
	if key == "" {
		return ctx
	}
	rc, _ := From(ctx)
	if rc == nil {
		ctx, rc = New(ctx)
	}
	if rc.Labels == nil {
		rc.Labels = make(map[string]string, 1)
	}
	rc.Labels[key] = value
	return ctx
}

// TenantID returns the tenant identifier from ctx if present.
func TenantID(ctx stdctx.Context) (string, bool) {
	if rc, ok := From(ctx); ok {
		if rc.TenantID != "" {
			return rc.TenantID, true
		}
	}
	return "", false
}

// UserID returns the user identifier from ctx if present.
func UserID(ctx stdctx.Context) (string, bool) {
	if rc, ok := From(ctx); ok {
		if rc.UserID != "" {
			return rc.UserID, true
		}
	}
	return "", false
}

// Duration returns time.Since(StartTime) if present, or 0 if StartTime is zero or ctx has no RequestContext.
func Duration(ctx stdctx.Context) time.Duration {
	if rc, ok := From(ctx); ok && !rc.StartTime.IsZero() {
		return chrono.Default.Since(rc.StartTime)
	}
	return 0
}

// Fields produces a map of safe fields suitable for structured logging.
// Potentially sensitive identifiers (like UserID) are included; callers are responsible for redaction policies.
func Fields(ctx stdctx.Context) map[string]any {
	fields := make(map[string]any, 6)
	rc, ok := From(ctx)
	if !ok || rc == nil {
		return fields
	}
	if rc.TraceID != "" {
		fields["trace_id"] = rc.TraceID
	}
	if rc.RequestID != "" {
		fields["request_id"] = rc.RequestID
	}
	if rc.UserID != "" {
		fields["user_id"] = rc.UserID
	}
	if rc.TenantID != "" {
		fields["tenant_id"] = rc.TenantID
	}
	if rc.SessionID != "" {
		fields["session_id"] = rc.SessionID
	}
	if len(rc.Labels) > 0 {
		fields["labels"] = rc.Labels
	}
	if !rc.StartTime.IsZero() {
		fields["duration_ms"] = chrono.Default.Since(rc.StartTime).Milliseconds()
	}
	return fields
}

// Validate performs basic size/format checks on IDs and labels.
// This function is conservative and avoids external dependencies.
func Validate(rc *RequestContext) error {
	if rc == nil {
		return errors.New("nil RequestContext")
	}
	if len(rc.TraceID) > 128 || len(rc.RequestID) > 128 || len(rc.UserID) > 128 || len(rc.TenantID) > 128 || len(rc.SessionID) > 128 {
		return errors.New("identifier too long")
	}
	if len(rc.Labels) > 32 { // keep labels small
		return errors.New("too many labels")
	}
	for k, v := range rc.Labels {
		if len(k) > 64 || len(v) > 256 {
			return errors.New("label size exceeded")
		}
	}
	return nil
}
