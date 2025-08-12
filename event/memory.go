package events

import (
	"context"
	"sync"
)

type memoryBus struct {
	cfg    BusConfig
	mu     sync.RWMutex
	topics map[string]*topic
	closed bool
}

type topic struct {
	ch      chan item
	workers int

	mu     sync.RWMutex
	subs   map[int64]subscription
	nextID int64
}

type subscription struct {
	handler Handler
	config  SubscribeConfig
}

type memorySub struct {
	bus   *memoryBus
	topic string
	id    int64
}

type item struct {
	ctx   context.Context
	event any
}

// NewMemoryBus creates an in-memory EventBus.
func NewMemoryBus(opts ...BusOption) EventBus {
	cfg := BusConfig{BufferSize: 64, WorkersPerTopic: 1}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &memoryBus{
		cfg:    cfg,
		topics: make(map[string]*topic),
	}
}

func (b *memoryBus) ensureTopic(name string) *topic {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	t := b.topics[name]
	if t == nil {
		t = &topic{
			ch:      make(chan item, b.cfg.BufferSize),
			workers: b.cfg.WorkersPerTopic,
			subs:    make(map[int64]subscription),
		}
		b.topics[name] = t

		// Start worker goroutines for this topic
		for i := 0; i < t.workers; i++ {
			go b.worker(name, t)
		}
	}
	return t
}

func (b *memoryBus) worker(topicName string, t *topic) {
	for item := range t.ch {
		// Snapshot current subscriptions to avoid holding locks during handler execution
		t.mu.RLock()
		subs := make([]subscription, 0, len(t.subs))
		for _, sub := range t.subs {
			subs = append(subs, sub)
		}
		t.mu.RUnlock()

		// Process each subscription
		for _, sub := range subs {
			retries := sub.config.Retries
			if retries <= 0 {
				retries = 1
			}

			var lastErr error
			for attempt := 1; attempt <= retries; attempt++ {
				if err := sub.handler(item.ctx, item.event); err != nil {
					lastErr = err
					continue
				}
				lastErr = nil
				break
			}

			// Call error handler if all retries failed
			if lastErr != nil && b.cfg.OnError != nil {
				b.cfg.OnError(item.ctx, topicName, item.event, lastErr)
			}
		}
	}
}

func (b *memoryBus) Subscribe(topicName string, handler Handler, opts ...SubscribeOption) (Subscription, error) {
	if handler == nil {
		return nil, ErrNilHandler
	}

	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return nil, ErrClosed
	}
	b.mu.RUnlock()

	topic := b.ensureTopic(topicName)
	if topic == nil {
		return nil, ErrClosed
	}

	// Build subscription config
	cfg := SubscribeConfig{Retries: 1}
	for _, opt := range opts {
		opt(&cfg)
	}

	// Register subscription
	topic.mu.Lock()
	id := topic.nextID + 1
	topic.nextID = id
	topic.subs[id] = subscription{
		handler: handler,
		config:  cfg,
	}
	topic.mu.Unlock()

	return &memorySub{bus: b, topic: topicName, id: id}, nil
}

func (s *memorySub) Unsubscribe() {
	s.bus.mu.RLock()
	topic := s.bus.topics[s.topic]
	s.bus.mu.RUnlock()

	if topic == nil {
		return
	}

	topic.mu.Lock()
	delete(topic.subs, s.id)
	topic.mu.Unlock()
}

func (b *memoryBus) Publish(ctx context.Context, topicName string, event any, opts ...PublishOption) error {
	if ctx == nil {
		ctx = context.Background()
	}

	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return ErrClosed
	}
	b.mu.RUnlock()

	// Process publish options
	var cfg PublishConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	// Attach headers to context if provided
	if len(cfg.Headers) > 0 {
		ctx = ContextWithHeaders(ctx, cfg.Headers)
	}

	topic := b.ensureTopic(topicName)
	if topic == nil {
		return ErrClosed
	}

	item := item{ctx: ctx, event: event}

	// Send to topic channel, respecting context cancellation
	select {
	case topic.ch <- item:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (b *memoryBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	b.closed = true
	for _, topic := range b.topics {
		close(topic.ch)
	}

	return nil
}
