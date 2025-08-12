package events

import (
	"context"
	"errors"
)

var (
	ErrClosed     = errors.New("events: bus closed")
	ErrNilHandler = errors.New("events: nil handler")
)

// Handler processes an event; returning an error signals failure (may be retried).
type Handler func(ctx context.Context, event any) error

// Subscription represents an active subscription; call Unsubscribe to stop receiving events.
type Subscription interface {
	Unsubscribe()
}

// EventBus is a simple, clean pub/sub interface.
type EventBus interface {
	Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscription, error)
	Publish(ctx context.Context, topic string, event any, opts ...PublishOption) error
	Close() error
}
