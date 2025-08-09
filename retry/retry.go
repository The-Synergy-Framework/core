package retry

import (
	"context"
	"errors"
	"time"
)

// Func is the function to retry.
type Func func(ctx context.Context) error

// ResultFunc runs a function that returns a value and error.
type ResultFunc[T any] func(ctx context.Context) (T, error)

// Policy returns the base backoff for a given 1-based attempt number.
// Implementations should be pure and fast.
type Policy func(attempt int) time.Duration

// Jitter tweaks the computed delay for an attempt.
type Jitter func(base time.Duration, attempt int) time.Duration

// RetryIf determines whether an error is retryable.
type RetryIf func(err error) bool

// OnRetry is called after a failed attempt, before sleeping.
type OnRetry func(ctx context.Context, attempt int, err error, nextDelay time.Duration)

// Options configures retry behavior.
type Options struct {
	MaxAttempts int
	Policy      Policy
	Jitter      Jitter
	MaxDelay    time.Duration
	RetryIf     RetryIf
	OnRetry     OnRetry
}

// Option applies a mutation to Options.
type Option func(*Options)

// WithMaxAttempts sets the maximum number of attempts (>= 1). Default 3.
func WithMaxAttempts(n int) Option { return func(o *Options) { o.MaxAttempts = n } }

// WithPolicy sets the backoff policy. Default: Exponential(100ms, 2).
func WithPolicy(p Policy) Option { return func(o *Options) { o.Policy = p } }

// WithJitter sets a jitter function to randomize delays. Default: none.
func WithJitter(j Jitter) Option { return func(o *Options) { o.Jitter = j } }

// WithMaxDelay caps the computed delay. Default: 30s.
func WithMaxDelay(d time.Duration) Option { return func(o *Options) { o.MaxDelay = d } }

// WithRetryIf sets a predicate to determine retryable errors. Default: retry any non-nil error.
func WithRetryIf(p RetryIf) Option { return func(o *Options) { o.RetryIf = p } }

// WithOnRetry sets a callback invoked after each failed attempt.
func WithOnRetry(cb OnRetry) Option { return func(o *Options) { o.OnRetry = cb } }

func defaults() Options {
	return Options{
		MaxAttempts: 3,
		Policy:      Exponential(100*time.Millisecond, 2.0),
		MaxDelay:    30 * time.Second,
		RetryIf:     func(err error) bool { return err != nil },
	}
}

// Do executes fn with retries according to options.
// Returns nil on success or the last error encountered.
func Do(ctx context.Context, fn Func, opts ...Option) error {
	if fn == nil {
		return errors.New("retry: nil function")
	}
	cfg := defaults()
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		// Respect context cancellation before attempt begins
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := fn(ctx); err == nil {
			return nil
		} else {
			lastErr = err
			if !cfg.RetryIf(err) || attempt == cfg.MaxAttempts {
				return lastErr
			}
			// Compute next delay
			d := cfg.Policy(attempt)
			if cfg.Jitter != nil {
				d = cfg.Jitter(d, attempt)
			}
			if cfg.MaxDelay > 0 && d > cfg.MaxDelay {
				d = cfg.MaxDelay
			}
			if d < 0 {
				d = 0
			}
			if cfg.OnRetry != nil {
				cfg.OnRetry(ctx, attempt, err, d)
			}
			// Sleep respecting context
			if err := sleep(ctx, d); err != nil {
				return err
			}
		}
	}
	return lastErr
}

// DoWithResult executes fn with retries and returns its result or the last error.
func DoWithResult[T any](ctx context.Context, fn ResultFunc[T], opts ...Option) (T, error) {
	var zero T
	if fn == nil {
		return zero, errors.New("retry: nil function")
	}
	var out T
	err := Do(ctx, func(c context.Context) error {
		v, err := fn(c)
		if err == nil {
			out = v
		}
		return err
	}, opts...)
	if err != nil {
		return zero, err
	}
	return out, nil
}

func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
