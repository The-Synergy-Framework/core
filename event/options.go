package events

import "context"

// SubscribeOption configures a subscription.
type SubscribeOption func(*SubscribeConfig)

// SubscribeConfig holds subscription configuration.
type SubscribeConfig struct {
	Retries int
}

// WithRetries sets number of attempts per event for this handler (default 1, i.e., no retry).
func WithRetries(n int) SubscribeOption {
	return func(c *SubscribeConfig) {
		if n > 0 {
			c.Retries = n
		}
	}
}

// PublishOption configures a publish operation.
type PublishOption func(*PublishConfig)

// PublishConfig holds publish configuration.
type PublishConfig struct {
	Headers map[string]string
	Key     string
}

// WithHeaders attaches metadata headers to the event.
func WithHeaders(headers map[string]string) PublishOption {
	return func(c *PublishConfig) {
		if len(headers) == 0 {
			return
		}
		if c.Headers == nil {
			c.Headers = make(map[string]string, len(headers))
		}
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// WithKey sets a partition key for the event (useful for distributed systems).
func WithKey(key string) PublishOption {
	return func(c *PublishConfig) {
		c.Key = key
	}
}

// BusOption configures an EventBus implementation.
type BusOption func(*BusConfig)

// BusConfig holds bus configuration.
type BusConfig struct {
	BufferSize      int
	WorkersPerTopic int
	OnError         func(ctx context.Context, topic string, event any, err error)
}

// WithBuffer sets the per-topic buffer size (default 64).
func WithBuffer(size int) BusOption {
	return func(c *BusConfig) {
		if size > 0 {
			c.BufferSize = size
		}
	}
}

// WithWorkers sets the number of workers per topic (default 1).
func WithWorkers(n int) BusOption {
	return func(c *BusConfig) {
		if n > 0 {
			c.WorkersPerTopic = n
		}
	}
}

// WithOnError sets a hook invoked when a handler returns error after its final retry.
func WithOnError(f func(ctx context.Context, topic string, event any, err error)) BusOption {
	return func(c *BusConfig) {
		c.OnError = f
	}
}
