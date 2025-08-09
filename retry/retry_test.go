package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_SucceedsFirstTry(t *testing.T) {
	ctx := context.Background()
	calls := 0
	err := Do(ctx, func(context.Context) error {
		calls++
		return nil
	})
	if err != nil || calls != 1 {
		t.Fatalf("want success in 1 call, got err=%v calls=%d", err, calls)
	}
}

func TestDo_MaxAttempts(t *testing.T) {
	ctx := context.Background()
	calls := 0
	wantErr := errors.New("boom")
	err := Do(ctx, func(context.Context) error {
		calls++
		return wantErr
	}, WithMaxAttempts(3), WithPolicy(Constant(1*time.Millisecond)))
	if !errors.Is(err, wantErr) {
		t.Fatalf("want err %v, got %v", wantErr, err)
	}
	if calls != 3 {
		t.Fatalf("want 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := Do(ctx, func(context.Context) error {
		calls++
		return errors.New("boom")
	}, WithMaxAttempts(5))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("function should not have executed, calls=%d", calls)
	}
}

func TestDoWithResult(t *testing.T) {
	ctx := context.Background()
	calls := 0
	res, err := DoWithResult[int](ctx, func(context.Context) (int, error) {
		calls++
		if calls < 2 {
			return 0, errors.New("retry")
		}
		return 42, nil
	}, WithPolicy(Constant(1*time.Millisecond)))
	if err != nil || res != 42 || calls != 2 {
		t.Fatalf("want res=42, calls=2, err=nil ; got res=%d calls=%d err=%v", res, calls, err)
	}
}

func TestJitterApplied(t *testing.T) {
	ctx := context.Background()
	calls := 0
	var slept time.Duration
	onRetry := func(_ context.Context, _ int, _ error, next time.Duration) { slept = next }
	_ = Do(ctx, func(context.Context) error {
		calls++
		if calls == 1 {
			return errors.New("boom")
		}
		return nil
	}, WithMaxAttempts(2), WithPolicy(Constant(10*time.Millisecond)), WithJitter(EqualJitter(nil)), WithOnRetry(onRetry))
	if slept <= 0 || slept > 10*time.Millisecond {
		t.Fatalf("expected jittered sleep in (0, 10ms], got %v", slept)
	}
}
