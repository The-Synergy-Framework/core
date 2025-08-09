package retry

import (
	"math"
	"math/rand"
	"time"
)

// Constant returns a policy that always returns d.
func Constant(d time.Duration) Policy { return func(int) time.Duration { return d } }

// Linear returns a policy that grows linearly: base*attempt.
func Linear(base time.Duration) Policy {
	return func(attempt int) time.Duration { return time.Duration(int64(base) * int64(attempt)) }
}

// Exponential returns a policy that grows exponentially: base * factor^(attempt-1).
func Exponential(base time.Duration, factor float64) Policy {
	if factor <= 1 {
		factor = 2
	}
	return func(attempt int) time.Duration {
		pow := math.Pow(factor, float64(attempt-1))
		return time.Duration(float64(base) * pow)
	}
}

// FullJitter implements AWS full jitter: rand(0, base).
func FullJitter(r *rand.Rand) Jitter {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return func(base time.Duration, _ int) time.Duration {
		if base <= 0 {
			return 0
		}
		return time.Duration(r.Int63n(int64(base)))
	}
}

// EqualJitter returns (base/2) + rand(0, base/2).
func EqualJitter(r *rand.Rand) Jitter {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return func(base time.Duration, _ int) time.Duration {
		if base <= 0 {
			return 0
		}
		half := base / 2
		return half + time.Duration(r.Int63n(int64(half)))
	}
}
