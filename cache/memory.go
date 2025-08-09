package cache

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	val        any
	exp        time.Time // zero means no expiration
	lastAccess time.Time
	ttl        time.Duration // original TTL used for sliding TTL
}

type memory struct {
	mu           sync.RWMutex
	items        map[string]entry
	defaultTTL   time.Duration
	cleanupEvery time.Duration
	stop         chan struct{}

	// features
	trackAccess bool
	sliding     bool
	trackStats  bool

	// stats
	hits      int
	misses    int
	evictions int
}

func newMemory(opts ...Option) *memory {
	m := &memory{
		items: make(map[string]entry),
		stop:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	if m.cleanupEvery > 0 {
		go m.janitor()
	}
	return m
}

func (m *memory) janitor() {
	t := time.NewTicker(m.cleanupEvery)
	defer t.Stop()
	for {
		select {
		case <-m.stop:
			return
		case <-t.C:
			m.cleanup()
		}
	}
}

func (m *memory) cleanup() {
	now := time.Now()
	m.mu.Lock()
	for k, e := range m.items {
		if !e.exp.IsZero() && now.After(e.exp) {
			delete(m.items, k)
			if m.trackStats {
				m.evictions++
			}
		}
	}
	m.mu.Unlock()
}

func (m *memory) Get(key string) (any, bool) {
	now := time.Now()
	m.mu.Lock()
	e, ok := m.items[key]
	if !ok {
		if m.trackStats {
			m.misses++
		}
		m.mu.Unlock()
		return nil, false
	}
	// expired?
	if !e.exp.IsZero() && now.After(e.exp) {
		delete(m.items, key)
		if m.trackStats {
			m.evictions++
			m.misses++
		}
		m.mu.Unlock()
		return nil, false
	}
	// last access
	if m.trackAccess {
		e.lastAccess = now
		// sliding TTL
		if m.sliding && e.ttl > 0 {
			e.exp = now.Add(e.ttl)
		}
		m.items[key] = e
	}
	if m.trackStats {
		m.hits++
	}
	v := e.val
	m.mu.Unlock()
	return v, true
}

func (m *memory) Set(key string, value any, ttl time.Duration) {
	exp := time.Time{}
	origTTL := time.Duration(0)
	if ttl <= 0 {
		if m.defaultTTL > 0 {
			exp = time.Now().Add(m.defaultTTL)
			origTTL = m.defaultTTL
		}
	} else {
		exp = time.Now().Add(ttl)
		origTTL = ttl
	}
	e := entry{val: value, exp: exp, ttl: origTTL}
	if m.trackAccess {
		e.lastAccess = time.Now()
	}
	m.mu.Lock()
	m.items[key] = e
	m.mu.Unlock()
}

func (m *memory) Delete(key string) {
	m.mu.Lock()
	delete(m.items, key)
	m.mu.Unlock()
}

func (m *memory) Clear() {
	m.mu.Lock()
	m.items = make(map[string]entry)
	m.mu.Unlock()
}

func (m *memory) Size() int {
	m.mu.RLock()
	n := len(m.items)
	m.mu.RUnlock()
	return n
}

func (m *memory) Close() {
	select {
	case <-m.stop:
		return
	default:
		close(m.stop)
	}
}

func (m *memory) GetOrCompute(ctx context.Context, key string, ttl time.Duration, compute func(context.Context) (any, error)) (any, error) {
	if v, ok := m.Get(key); ok {
		return v, nil
	}
	if compute == nil {
		return nil, nil
	}
	v, err := compute(ctx)
	if err != nil {
		return nil, err
	}
	m.Set(key, v, ttl)
	return v, nil
}

func (m *memory) Info(key string) (expiresAt time.Time, lastAccess time.Time, ok bool) {
	now := time.Now()
	m.mu.RLock()
	e, ok := m.items[key]
	m.mu.RUnlock()
	if !ok {
		return time.Time{}, time.Time{}, false
	}
	if !e.exp.IsZero() && now.After(e.exp) {
		return time.Time{}, time.Time{}, false
	}
	return e.exp, e.lastAccess, true
}

func (m *memory) Stats() (hits, misses, evictions, size int) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hits, m.misses, m.evictions, len(m.items)
}

var _ Cache = (*memory)(nil)
