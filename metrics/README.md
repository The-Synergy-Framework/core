# core/metrics

Enterprise-grade, thread-safe metrics API with comprehensive validation and observability support.

## Features

- **Thread-safe metric interfaces**: Counter, Gauge, Histogram with proper concurrency
- **Context-aware operations**: All metric operations accept `context.Context`
- **Comprehensive validation**: Metric names and label validation following Prometheus conventions
- **Flexible labeling**: Dynamic labels per operation + constant labels per metric
- **Enterprise timer utilities**: Built-in timing with proper cleanup and error handling
- **No-op default**: Safe to use without configuration
- **Multi-registry support**: Fan-out to multiple backends
- **Clean interface design**: Core package contains only interfaces and validation

## Quick Start

```go
import (
	"context"
	"core/metrics"
)

func main() {
	// Set up your registry implementation
	// reg := prometheus.New() // or opentelemetry.New(), etc.
	// metrics.SetDefault(reg)

	// Create metrics with validation
	counter, err := metrics.Default().NewCounter(metrics.MetricOptions{
		Name: "requests_total",
		Help: "Total number of requests",
		Unit: "requests",
	})
	if err != nil {
		panic(err)
	}

	// Use with context and labels
	ctx := context.Background()
	counter.Inc(ctx, metrics.Labels{"method": "GET", "route": "/"})
	counter.Add(ctx, 5, metrics.Labels{"method": "POST", "route": "/api"})
}
```

## Timer Usage

```go
// Method 1: Manual timer
hist, _ := metrics.Default().NewHistogram(metrics.HistogramOptions{
	MetricOptions: metrics.MetricOptions{
		Name: "request_duration_seconds",
		Help: "Request duration in seconds",
		Unit: "seconds",
	},
})

timer := metrics.NewTimer(ctx, hist, metrics.Labels{"route": "/"})
// ... do work ...
timer.Stop() // Records elapsed time

// Method 2: Function wrapper
metrics.ObserveDuration(ctx, hist, metrics.Labels{"route": "/"}, func() {
	// ... do work ...
})
```

## Multi-Registry Setup
```go
// Fan out to multiple registries
prometheus := /* your prometheus registry */
cloudwatch := /* your cloudwatch registry */

multi := metrics.Multi(prometheus, cloudwatch)
metrics.SetDefault(multi)
```

## Validation
```go
// Names must follow [a-zA-Z_:][a-zA-Z0-9_:]*
err := metrics.ValidateMetricName("invalid-name") // Returns error

// Label keys follow same rules, values are more flexible
err := metrics.ValidateLabels(metrics.Labels{
	"__reserved": "value", // Error: reserved prefix
	"valid_key": "any value is ok",
})
```

## Production Adapters

The core package provides only interfaces. For production use, you'll need adapter implementations:

- **Prometheus**: `core/metrics/prometheus` (recommended for most use cases)
- **OpenTelemetry**: `core/metrics/opentelemetry` (modern observability)
- **StatsD**: `core/metrics/statsd` (legacy systems)
- **CloudWatch**: `core/metrics/cloudwatch` (AWS environments)

## No-Op Behavior

Without a configured registry, all operations are no-ops:

```go
// Safe to call even with no registry configured
counter, _ := metrics.Default().NewCounter(metrics.MetricOptions{Name: "test"})
counter.Inc(context.Background(), nil) // No-op, no errors
```

## Breaking Changes from v1

- **Context required**: All metric operations now require `context.Context`
- **Error handling**: Registry methods return `(Metric, error)` instead of just `Metric`
- **Separate options**: Histograms use `HistogramOptions` instead of generic `Options`
- **Dynamic labels**: Labels are passed per operation, not just at creation time
- **Timer API**: Replaced `Start()/Stop()` with proper timer objects
- **No built-in implementations**: Core package is interface-only, use adapter packages
- **Thread safety**: Improved validation and error handling 