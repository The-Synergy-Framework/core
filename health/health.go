// Package health provides a standardized health checking system for services.
// It defines health status types and interfaces for implementing health checks
// that can be used for monitoring, load balancing, and alerting.
package health

import "context"

// Status represents the health status of a service or component.
// Health statuses follow industry standards and are designed to be
// compatible with container orchestration systems, monitoring tools,
// and alerting systems.
type Status string

const (
	// StatusHealthy indicates the service is fully operational and functioning normally.
	StatusHealthy Status = "healthy"

	// StatusUnhealthy indicates the service is completely down or experiencing critical failures.
	StatusUnhealthy Status = "unhealthy"

	// StatusDegraded indicates the service is partially functional but experiencing issues.
	StatusDegraded Status = "degraded"

	// StatusUnknown indicates the health status cannot be determined.
	StatusUnknown Status = "unknown"
)

// The Result represents the result of a health check operation.
// It provides detailed information about the health status including
// a human-readable message and optional details for debugging.
type Result struct {
	Status  Status
	Message string
	Details map[string]any
}

// Checker defines the interface for implementing health checks.
// Implementations should perform the necessary checks and return
// a Result with appropriate status and details.
//
// Health check implementations should:
//   - Be fast and non-blocking when possible
//   - Provide meaningful error messages
//   - Include relevant metrics in Details
//   - Handle timeouts and cancellation via ctx
//   - Avoid expensive operations
//
// The error return should be reserved for probe failures (e.g., misconfiguration
// or the check itself failing). Health state belongs in Result.Status.
type Checker interface {
	Check(ctx context.Context) (*Result, error)
}

// FuncChecker adapts a function to a Checker.
type FuncChecker func(ctx context.Context) (*Result, error)

// Check implements Checker.
func (f FuncChecker) Check(ctx context.Context) (*Result, error) { return f(ctx) }

// Helper constructors

// OK creates a healthy result with an optional message and details.
func OK(message string, details map[string]any) *Result {
	return &Result{Status: StatusHealthy, Message: message, Details: details}
}

// Degraded creates a degraded result with an optional message and details.
func Degraded(message string, details map[string]any) *Result {
	return &Result{Status: StatusDegraded, Message: message, Details: details}
}

// Unhealthy creates an unhealthy result with an optional message and details.
func Unhealthy(message string, details map[string]any) *Result {
	return &Result{Status: StatusUnhealthy, Message: message, Details: details}
}

// Unknown creates an unknown result with an optional message and details.
func Unknown(message string, details map[string]any) *Result {
	return &Result{Status: StatusUnknown, Message: message, Details: details}
}
