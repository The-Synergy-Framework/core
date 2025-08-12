package events

import "context"

// TypedHandler is a generic handler that accepts a specific event type T.
type TypedHandler[T any] func(ctx context.Context, event T) error

// AsHandler wraps a TypedHandler[T] into an untyped Handler that asserts the value at runtime.
// If the assertion fails, the handler is a no-op and returns nil.
func AsHandler[T any](h TypedHandler[T]) Handler {
	return func(ctx context.Context, event any) error {
		v, ok := event.(T)
		if !ok {
			return nil
		}
		return h(ctx, v)
	}
}

// SubscribeTyped is a helper that subscribes a typed handler to an EventBus.
func SubscribeTyped[T any](bus EventBus, topic string, handler TypedHandler[T], opts ...SubscribeOption) (Subscription, error) {
	return bus.Subscribe(topic, AsHandler(handler), opts...)
}
