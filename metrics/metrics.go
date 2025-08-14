package metrics

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Labels represent key-value pairs for metric dimensions.
type Labels map[string]string

// MetricOptions configure a metric instrument.
type MetricOptions struct {
	Name        string // Required: metric name (must match [a-zA-Z_:][a-zA-Z0-9_:]*
	Help        string // Optional: human-readable description
	Unit        string // Optional: unit of measurement (e.g., "seconds", "bytes")
	ConstLabels Labels // Optional: labels attached to all samples
}

// HistogramOptions configure a histogram instrument.
type HistogramOptions struct {
	MetricOptions
	Buckets []float64 // Optional: custom buckets (default: DefaultBuckets)
}

// Counter represents a monotonically increasing metric.
type Counter interface {
	// Inc increments the counter by 1
	Inc(ctx context.Context, labels Labels)
	// Add increments the counter by delta (must be >= 0)
	Add(ctx context.Context, delta float64, labels Labels)
}

// Gauge represents a metric that can go up and down.
type Gauge interface {
	// Set sets the gauge to value
	Set(ctx context.Context, value float64, labels Labels)
	// Add adds delta to the gauge value
	Add(ctx context.Context, delta float64, labels Labels)
	// Inc increments the gauge by 1
	Inc(ctx context.Context, labels Labels)
	// Dec decrements the gauge by 1
	Dec(ctx context.Context, labels Labels)
}

// Histogram samples observations into buckets.
type Histogram interface {
	// Observe adds a single observation to the histogram
	Observe(ctx context.Context, value float64, labels Labels)
}

// Registry creates and manages metric instruments.
type Registry interface {
	// NewCounter creates a new counter with the given options
	NewCounter(opts MetricOptions) (Counter, error)
	// NewGauge creates a new gauge with the given options
	NewGauge(opts MetricOptions) (Gauge, error)
	// NewHistogram creates a new histogram with the given options
	NewHistogram(opts HistogramOptions) (Histogram, error)
}

// Timer measures elapsed time and records it to a histogram.
type Timer struct {
	start  time.Time
	hist   Histogram
	ctx    context.Context
	labels Labels
}

// DefaultBuckets for histograms (in seconds, suitable for request durations).
var DefaultBuckets = []float64{
	0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
}

var metricNameRegex = regexp.MustCompile(`^[a-zA-Z_:][a-zA-Z0-9_:]*$`)

// ValidateMetricName validates a metric name according to Prometheus conventions.
func ValidateMetricName(name string) error {
	if name == "" {
		return fmt.Errorf("metric name cannot be empty")
	}
	if !metricNameRegex.MatchString(name) {
		return fmt.Errorf("invalid metric name %q: must match [a-zA-Z_:][a-zA-Z0-9_:]*", name)
	}
	return nil
}

// ValidateLabels validates label keys and values.
func ValidateLabels(labels Labels) error {
	for key, value := range labels {
		if key == "" {
			return fmt.Errorf("label key cannot be empty")
		}
		if strings.HasPrefix(key, "__") {
			return fmt.Errorf("label key %q cannot start with '__' (reserved prefix)", key)
		}
		if !metricNameRegex.MatchString(key) {
			return fmt.Errorf("invalid label key %q: must match [a-zA-Z_:][a-zA-Z0-9_:]*", key)
		}
		// Value validation is more lenient - just check for null bytes
		if strings.Contains(value, "\x00") {
			return fmt.Errorf("label value for key %q cannot contain null bytes", key)
		}
	}
	return nil
}

// NewTimer creates a timer that will record elapsed time to the given histogram.
func NewTimer(ctx context.Context, hist Histogram, labels Labels) *Timer {
	return &Timer{
		start:  time.Now(),
		hist:   hist,
		ctx:    ctx,
		labels: labels,
	}
}

// Stop stops the timer and records the elapsed time to the histogram.
func (t *Timer) Stop() {
	if t.hist != nil {
		elapsed := time.Since(t.start).Seconds()
		t.hist.Observe(t.ctx, elapsed, t.labels)
	}
}

// ObserveDuration is a convenience function to time a function call.
func ObserveDuration(ctx context.Context, hist Histogram, labels Labels, fn func()) {
	timer := NewTimer(ctx, hist, labels)
	defer timer.Stop()
	fn()
}
