package cache

import (
	"context"
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	c := NewMemory()
	defer c.Close()
	c.Set("a", 1, 0)
	v, ok := c.Get("a")
	if !ok || v.(int) != 1 {
		t.Fatalf("get mismatch")
	}
}

func TestTTL(t *testing.T) {
	c := NewMemory()
	defer c.Close()
	c.Set("a", 1, 10*time.Millisecond)
	time.Sleep(15 * time.Millisecond)
	if _, ok := c.Get("a"); ok {
		t.Fatalf("value should have expired")
	}
}

func TestCleanupJanitor(t *testing.T) {
	c := NewMemory(WithCleanupInterval(5 * time.Millisecond))
	defer c.Close()
	c.Set("a", 1, 5*time.Millisecond)
	time.Sleep(20 * time.Millisecond) // allow cleanup tick
	if _, ok := c.Get("a"); ok {
		t.Fatalf("value should have been cleaned up")
	}
}

func TestGetOrCompute(t *testing.T) {
	c := NewMemory()
	defer c.Close()
	calls := 0
	compute := func(context.Context) (any, error) {
		calls++
		return 42, nil
	}
	v, err := c.GetOrCompute(context.Background(), "x", 0, compute)
	if err != nil || v.(int) != 42 || calls != 1 {
		t.Fatalf("first compute failed: v=%v err=%v calls=%d", v, err, calls)
	}
	v, err = c.GetOrCompute(context.Background(), "x", 0, compute)
	if err != nil || v.(int) != 42 || calls != 1 {
		t.Fatalf("should have used cached value: v=%v err=%v calls=%d", v, err, calls)
	}
}

func TestLastAccessAndSlidingTTL(t *testing.T) {
	c := NewMemory(WithLastAccessTracking(), WithSlidingTTL())
	defer c.Close()
	c.Set("k", 1, 10*time.Millisecond)
	// Access a few times to update lastAccess and extend expiry
	for i := 0; i < 3; i++ {
		v, ok := c.Get("k")
		if !ok || v.(int) != 1 {
			t.Fatalf("missing value on iteration %d", i)
		}
		time.Sleep(5 * time.Millisecond)
	}
	exp, last, ok := c.Info("k")
	if !ok {
		t.Fatalf("info missing")
	}
	if exp.IsZero() || last.IsZero() {
		t.Fatalf("expected non-zero exp and last access")
	}
	// After sliding accesses, it should not be expired yet
	if _, ok := c.Get("k"); !ok {
		t.Fatalf("value should still be present due to sliding TTL")
	}
}

func TestStats(t *testing.T) {
	c := NewMemory(WithStats())
	defer c.Close()
	_, _ = c.Get("miss")
	c.Set("hit", 1, 0)
	_, _ = c.Get("hit")
	h, m, e, size := c.Stats()
	if h < 1 || m < 1 || e < 0 || size < 1 {
		t.Fatalf("unexpected stats: h=%d m=%d e=%d size=%d", h, m, e, size)
	}
}
