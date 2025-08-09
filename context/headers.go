package ctx

// Canonical header keys for simple request correlation and tenancy.
// These are optional conveniences and may be mapped to project-specific names.
const (
	HeaderRequestID = "X-Request-Id"
	HeaderTraceID   = "X-Trace-Id" // fallback when not using W3C traceparent
	HeaderUserID    = "X-User-Id"
	HeaderTenantID  = "X-Tenant-Id"
	HeaderSessionID = "X-Session-Id"
)
