package events

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMemoryBus_PublishSubscribe(t *testing.T) {
	bus := NewMemoryBus(WithBuffer(16), WithWorkers(1))
	defer bus.Close()

	var got atomic.Int64
	var wg sync.WaitGroup
	wg.Add(2)

	_, err := bus.Subscribe("topic", func(ctx context.Context, evt any) error {
		defer wg.Done()
		if v, ok := evt.(int); ok {
			got.Add(int64(v))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.Publish(context.Background(), "topic", 3); err != nil {
		t.Fatalf("publish: %v", err)
	}
	if err := bus.Publish(context.Background(), "topic", 4); err != nil {
		t.Fatalf("publish: %v", err)
	}

	waitDone(t, &wg)
	if got.Load() != 7 {
		t.Fatalf("got %d, want 7", got.Load())
	}
}

func TestMemoryBus_Retries(t *testing.T) {
	bus := NewMemoryBus(WithBuffer(8), WithWorkers(1))
	defer bus.Close()

	var tries atomic.Int32
	var wg sync.WaitGroup
	wg.Add(3)

	_, err := bus.Subscribe("t", func(ctx context.Context, evt any) error {
		tries.Add(1)
		wg.Done()
		return errors.New("fail")
	}, WithRetries(3))
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.Publish(context.Background(), "t", 1); err != nil {
		t.Fatalf("publish: %v", err)
	}

	waitDone(t, &wg)
	if tries.Load() != 3 {
		t.Fatalf("got tries=%d, want 3", tries.Load())
	}
}

func TestMemoryBus_ClosePreventsPublishAndSubscribe(t *testing.T) {
	bus := NewMemoryBus()
	bus.Close()

	if _, err := bus.Subscribe("x", func(context.Context, any) error { return nil }); !errors.Is(err, ErrClosed) {
		t.Fatalf("subscribe err=%v, want ErrClosed", err)
	}
	if err := bus.Publish(context.Background(), "x", 1); !errors.Is(err, ErrClosed) {
		t.Fatalf("publish err=%v, want ErrClosed", err)
	}
}

func TestTypedHelpers(t *testing.T) {
	bus := NewMemoryBus(WithBuffer(1), WithWorkers(1))
	defer bus.Close()

	var got int32
	var wg sync.WaitGroup
	wg.Add(1)

	_, err := SubscribeTyped[int](bus, "ints", func(ctx context.Context, v int) error {
		atomic.AddInt32(&got, int32(v))
		wg.Done()
		return nil
	})
	if err != nil {
		t.Fatalf("subscribe typed: %v", err)
	}

	if err := bus.Publish(context.Background(), "ints", 9); err != nil {
		t.Fatalf("publish: %v", err)
	}

	waitDone(t, &wg)
	if atomic.LoadInt32(&got) != 9 {
		t.Fatalf("got %d, want 9", got)
	}
}

func TestPublishContextCancel(t *testing.T) {
	bus := NewMemoryBus(WithBuffer(0)) // unbuffered forces blocking
	defer bus.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		_ = bus.Publish(ctx, "no_subs", 1)
		close(done)
	}()

	select {
	case <-done:
		// ok: publish returned due to context deadline
	case <-time.After(200 * time.Millisecond):
		t.Fatal("publish did not return on context deadline")
	}
}

func TestPublishWithHeaders(t *testing.T) {
	bus := NewMemoryBus(WithBuffer(1), WithWorkers(1))
	defer bus.Close()

	var receivedHeaders map[string]string
	var wg sync.WaitGroup
	wg.Add(1)

	_, err := bus.Subscribe("test", func(ctx context.Context, evt any) error {
		defer wg.Done()
		headers, _ := HeadersFrom(ctx)
		receivedHeaders = headers
		return nil
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	testHeaders := map[string]string{"key1": "value1", "key2": "value2"}
	if err := bus.Publish(context.Background(), "test", "data", WithHeaders(testHeaders)); err != nil {
		t.Fatalf("publish: %v", err)
	}

	waitDone(t, &wg)
	if len(receivedHeaders) != 2 || receivedHeaders["key1"] != "value1" || receivedHeaders["key2"] != "value2" {
		t.Fatalf("got headers %v, want %v", receivedHeaders, testHeaders)
	}
}

func waitDone(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
		return
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for handlers")
	}
}
