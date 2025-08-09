package ctx

import (
	"context"
	"testing"
	"time"
)

func TestNewAndFrom(t *testing.T) {
	parent := context.Background()
	ctx, rc := New(parent)
	if rc == nil {
		t.Fatalf("expected rc, got nil")
	}
	if _, ok := From(parent); ok {
		t.Fatalf("parent should not have rc")
	}
	if got, ok := From(ctx); !ok || got == nil {
		t.Fatalf("expected rc in ctx")
	}
}

func TestEnrichers(t *testing.T) {
	ctx, _ := New(context.Background())
	ctx = WithTrace(ctx, "t-1")
	ctx = WithRequestID(ctx, "r-1")
	ctx = WithUser(ctx, "u-1")
	ctx = WithTenant(ctx, "ten-1")
	ctx = WithSession(ctx, "s-1")
	ctx = WithLabel(ctx, "k", "v")

	rc, ok := From(ctx)
	if !ok {
		t.Fatalf("missing rc")
	}
	if rc.TraceID != "t-1" || rc.RequestID != "r-1" || rc.UserID != "u-1" || rc.TenantID != "ten-1" || rc.SessionID != "s-1" {
		t.Fatalf("unexpected rc values: %+v", rc)
	}
	if rc.Labels["k"] != "v" {
		t.Fatalf("label not set")
	}
}

func TestAccessorsAndDuration(t *testing.T) {
	ctx, rc := New(context.Background())
	rc.StartTime = time.Now().Add(-10 * time.Millisecond)
	ctx = Into(ctx, rc)

	if id, ok := TenantID(ctx); ok || id != "" {
		t.Fatalf("unexpected tenant: %v %v", id, ok)
	}
	ctx = WithTenant(ctx, "ten")
	if id, ok := TenantID(ctx); !ok || id != "ten" {
		t.Fatalf("tenant accessor failed")
	}
	if d := Duration(ctx); d <= 0 {
		t.Fatalf("expected positive duration, got %v", d)
	}
}

func TestValidate(t *testing.T) {
	_, rc := New(context.Background())
	// valid
	rc.Labels = map[string]string{"a": "1"}
	if err := Validate(rc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// too many labels (33 distinct keys)
	rc.Labels = map[string]string{}
	for i := 0; i < 33; i++ {
		key := "k" + string(rune('a'+(i%26))) + string(rune('A'+(i%26))) + string(rune('0'+(i%10)))
		rc.Labels[key] = "v"
	}
	if err := Validate(rc); err == nil {
		t.Fatalf("expected error for too many labels")
	}
}
