package chrono

import "time"

// Clock provides time-related functions behind an interface for testability.
// Implementations must be safe for concurrent use.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
	// Since returns the duration since t.
	Since(t time.Time) time.Duration
}

// SystemClock is the default Clock using the standard library time package.
// It has no internal state and is safe for concurrent use.
type SystemClock struct{}

// Now implements Clock.
func (SystemClock) Now() time.Time { return time.Now() }

// Since implements Clock.
func (SystemClock) Since(t time.Time) time.Duration { return time.Since(t) }

// Default is the package-level clock used by helper functions.
var Default Clock = SystemClock{}

// SetDefault sets the package-level default clock.
// This is primarily intended for tests.
func SetDefault(c Clock) {
	if c == nil {
		Default = SystemClock{}
		return
	}
	Default = c
}
