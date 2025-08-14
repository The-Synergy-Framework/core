package metrics

import (
	"context"
	"fmt"
	"sync"
)

var (
	globalMu        sync.RWMutex
	defaultRegistry Registry = &noopRegistry{}
)

// SetDefault sets the global default registry.
// This should typically be called once during application startup.
func SetDefault(r Registry) {
	globalMu.Lock()
	defer globalMu.Unlock()
	if r == nil {
		defaultRegistry = &noopRegistry{}
	} else {
		defaultRegistry = r
	}
}

// Default returns the global default registry.
func Default() Registry {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return defaultRegistry
}

// Multi creates a registry that fans out to multiple registries.
// If any registry fails during metric creation, the error is returned.
func Multi(registries ...Registry) Registry {
	// Filter out nil registries
	var validRegs []Registry
	for _, r := range registries {
		if r != nil {
			validRegs = append(validRegs, r)
		}
	}
	return &multiRegistry{registries: validRegs}
}

type multiRegistry struct {
	registries []Registry
}

func (m *multiRegistry) NewCounter(opts MetricOptions) (Counter, error) {
	var counters []Counter
	for _, r := range m.registries {
		c, err := r.NewCounter(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create counter in registry: %w", err)
		}
		counters = append(counters, c)
	}
	return &multiCounter{counters: counters}, nil
}

func (m *multiRegistry) NewGauge(opts MetricOptions) (Gauge, error) {
	var gauges []Gauge
	for _, r := range m.registries {
		g, err := r.NewGauge(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create gauge in registry: %w", err)
		}
		gauges = append(gauges, g)
	}
	return &multiGauge{gauges: gauges}, nil
}

func (m *multiRegistry) NewHistogram(opts HistogramOptions) (Histogram, error) {
	var histograms []Histogram
	for _, r := range m.registries {
		h, err := r.NewHistogram(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to create histogram in registry: %w", err)
		}
		histograms = append(histograms, h)
	}
	return &multiHistogram{histograms: histograms}, nil
}

type multiCounter struct {
	counters []Counter
}

func (m *multiCounter) Inc(ctx context.Context, labels Labels) {
	for _, c := range m.counters {
		c.Inc(ctx, labels)
	}
}

func (m *multiCounter) Add(ctx context.Context, delta float64, labels Labels) {
	for _, c := range m.counters {
		c.Add(ctx, delta, labels)
	}
}

type multiGauge struct {
	gauges []Gauge
}

func (m *multiGauge) Set(ctx context.Context, value float64, labels Labels) {
	for _, g := range m.gauges {
		g.Set(ctx, value, labels)
	}
}

func (m *multiGauge) Add(ctx context.Context, delta float64, labels Labels) {
	for _, g := range m.gauges {
		g.Add(ctx, delta, labels)
	}
}

func (m *multiGauge) Inc(ctx context.Context, labels Labels) {
	for _, g := range m.gauges {
		g.Inc(ctx, labels)
	}
}

func (m *multiGauge) Dec(ctx context.Context, labels Labels) {
	for _, g := range m.gauges {
		g.Dec(ctx, labels)
	}
}

type multiHistogram struct {
	histograms []Histogram
}

func (m *multiHistogram) Observe(ctx context.Context, value float64, labels Labels) {
	for _, h := range m.histograms {
		h.Observe(ctx, value, labels)
	}
}

// noopRegistry is a no-op implementation for when no registry is configured.
type noopRegistry struct{}

func (n *noopRegistry) NewCounter(opts MetricOptions) (Counter, error) {
	if err := ValidateMetricName(opts.Name); err != nil {
		return nil, err
	}
	if err := ValidateLabels(opts.ConstLabels); err != nil {
		return nil, err
	}
	return &noopCounter{}, nil
}

func (n *noopRegistry) NewGauge(opts MetricOptions) (Gauge, error) {
	if err := ValidateMetricName(opts.Name); err != nil {
		return nil, err
	}
	if err := ValidateLabels(opts.ConstLabels); err != nil {
		return nil, err
	}
	return &noopGauge{}, nil
}

func (n *noopRegistry) NewHistogram(opts HistogramOptions) (Histogram, error) {
	if err := ValidateMetricName(opts.Name); err != nil {
		return nil, err
	}
	if err := ValidateLabels(opts.ConstLabels); err != nil {
		return nil, err
	}
	return &noopHistogram{}, nil
}

type noopCounter struct{}

func (n *noopCounter) Inc(ctx context.Context, labels Labels)                {}
func (n *noopCounter) Add(ctx context.Context, delta float64, labels Labels) {}

type noopGauge struct{}

func (n *noopGauge) Set(ctx context.Context, value float64, labels Labels) {}
func (n *noopGauge) Add(ctx context.Context, delta float64, labels Labels) {}
func (n *noopGauge) Inc(ctx context.Context, labels Labels)                {}
func (n *noopGauge) Dec(ctx context.Context, labels Labels)                {}

type noopHistogram struct{}

func (n *noopHistogram) Observe(ctx context.Context, value float64, labels Labels) {}
