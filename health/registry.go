package health

import (
	"context"
	"time"
)

// Entry represents a registered checker and timing info.
type Entry struct {
	Name     string
	Result   *Result
	Error    error
	Duration time.Duration
}

// Summary aggregates results across all checkers.
type Summary struct {
	Overall Status
	Entries []Entry
}

// Registry holds named checkers.
type Registry struct {
	checks map[string]Checker
}

// New creates a registry.
func New() *Registry { return &Registry{checks: map[string]Checker{}} }

// Register adds or replaces a checker.
func (r *Registry) Register(name string, c Checker) { r.checks[name] = c }

// RunAll executes all registered checks with the given context.
// Overall status is the worst among results (Unhealthy > Degraded > Unknown > Healthy).
func (r *Registry) RunAll(ctx context.Context) Summary {
	entries := make([]Entry, 0, len(r.checks))
	overall := StatusHealthy
	for name, c := range r.checks {
		start := time.Now()
		res, err := c.Check(ctx)
		dur := time.Since(start)
		entries = append(entries, Entry{Name: name, Result: res, Error: err, Duration: dur})
		st := statusFrom(res, err)
		if worse(st, overall) {
			overall = st
		}
	}
	return Summary{Overall: overall, Entries: entries}
}

func statusFrom(res *Result, err error) Status {
	if err != nil {
		return StatusUnknown
	}
	if res == nil {
		return StatusUnknown
	}
	return res.Status
}

func worse(a, b Status) bool {
	rank := func(s Status) int {
		switch s {
		case StatusUnhealthy:
			return 4
		case StatusDegraded:
			return 3
		case StatusUnknown:
			return 2
		case StatusHealthy:
			fallthrough
		default:
			return 1
		}
	}
	return rank(a) > rank(b)
}
