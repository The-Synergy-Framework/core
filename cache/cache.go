package cache

import (
	"context"
	"time"
)

// Cache defines a minimal key/value in-memory cache with TTL support.
// Implementations must be safe for concurrent use.
type Cache interface {
	// Get returns the value for key if present and not expired.
	Get(key string) (value any, ok bool)
	// Set stores value for key with an optional TTL. If ttl<=0, uses default TTL if configured; if no default, stores without expiration.
	Set(key string, value any, ttl time.Duration)
	// Delete removes a key.
	Delete(key string)
	// Clear removes all entries.
	Clear()
	// Size returns the number of live entries (best-effort).
	Size() int
	// Close releases resources (e.g., background janitor). It is safe to call multiple times.
	Close()
	// GetOrCompute returns the cached value for key or invokes compute to produce and store it.
	// If compute returns an error, nothing is cached and the error is returned.
	GetOrCompute(ctx context.Context, key string, ttl time.Duration, compute func(context.Context) (any, error)) (any, error)
	// Info returns the expiration time and last-access time for key, if present and not expired.
	// If last-access tracking is disabled or not yet accessed, lastAccess may be zero.
	Info(key string) (expiresAt time.Time, lastAccess time.Time, ok bool)
	// Stats returns hits, misses, evictions (due to expiry), and current size.
	Stats() (hits, misses, evictions, size int)
}

// Option configures a Memory cache.
type Option func(*memory)

// WithDefaultTTL sets a default TTL applied when Set is called with ttl<=0.
func WithDefaultTTL(ttl time.Duration) Option {
	return func(m *memory) {
		m.defaultTTL = ttl
	}
}

// WithCleanupInterval sets the periodic cleanup tick. If <=0, cleanup runs opportunistically on access only.
func WithCleanupInterval(d time.Duration) Option {
	return func(m *memory) {
		m.cleanupEvery = d
	}
}

// WithLastAccessTracking enables tracking of last-access timestamps.
func WithLastAccessTracking() Option {
	return func(m *memory) {
		m.trackAccess = true
	}
}

// WithSlidingTTL enables expire-after-access behavior. Entries with a TTL will have their expiration extended on access.
func WithSlidingTTL() Option {
	return func(m *memory) {
		m.sliding = true
	}
}

// WithStats enables hit/miss/eviction counters.
func WithStats() Option {
	return func(m *memory) {
		m.trackStats = true
	}
}

// NewMemory returns a new in-memory cache.
func NewMemory(opts ...Option) Cache {
	return newMemory(opts...)
}
