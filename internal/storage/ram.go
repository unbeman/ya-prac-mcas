package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

// ramRepository is implementation of Repository,
// describes storage based on maps.
type ramRepository struct {
	sync.RWMutex
	counterStorage map[string]metrics.Counter
	gaugeStorage   map[string]metrics.Gauge
}

// NewRAMRepository creates ramRepository.
func NewRAMRepository() *ramRepository {
	return &ramRepository{
		counterStorage: map[string]metrics.Counter{},
		gaugeStorage:   map[string]metrics.Gauge{},
	}
}

// getCounter returns metrics.Counter by name.
func (rs *ramRepository) getCounter(name string) (metrics.Counter, error) {
	value, ok := rs.counterStorage[name]
	if !ok {
		return nil, fmt.Errorf("counter (%v) %w", name, ErrNotFound)
	}
	return value, nil
}

// GetCounter returns metrics.Counter by name, calling getCounter.
func (rs *ramRepository) GetCounter(ctx context.Context, name string) (metrics.Counter, error) {
	rs.RLock()
	defer rs.RUnlock()
	return rs.getCounter(name)
}

// addCounter increases by delta counter and return metrics.Counter.
func (rs *ramRepository) addCounter(ctx context.Context, name string, value int64) (metrics.Counter, error) {
	counter, err := rs.getCounter(name)
	if errors.Is(err, ErrNotFound) {
		counter = metrics.NewCounter(name, 0)
		rs.counterStorage[name] = counter
	}
	counter.Add(value)
	return counter, nil
}

// AddCounter increases by delta counter and return metrics.Counter,
func (rs *ramRepository) AddCounter(ctx context.Context, name string, value int64) (metrics.Counter, error) {
	rs.Lock()
	defer rs.Unlock()
	return rs.addCounter(ctx, name, value)
}

// getGauge returns metrics.Gauge by name.
func (rs *ramRepository) getGauge(name string) (metrics.Gauge, error) {
	value, ok := rs.gaugeStorage[name]
	if !ok {
		return nil, fmt.Errorf("gauge (%v) %w", name, ErrNotFound)
	}
	return value, nil
}

// GetGauge returns metrics.Gauge by name.
func (rs *ramRepository) GetGauge(ctx context.Context, name string) (metrics.Gauge, error) {
	rs.RLock()
	defer rs.RUnlock()
	return rs.getGauge(name)
}

// setGauge sets new value gauge and returns metrics.Gauge.
func (rs *ramRepository) setGauge(ctx context.Context, name string, value float64) (metrics.Gauge, error) {
	gauge, err := rs.getGauge(name)
	if errors.Is(err, ErrNotFound) {
		gauge = metrics.NewGauge(name, 0)
		rs.gaugeStorage[name] = gauge
	}
	gauge.Set(value)
	return gauge, nil
}

// SetGauge sets new value gauge and returns metrics.Gauge.
func (rs *ramRepository) SetGauge(ctx context.Context, name string, value float64) (metrics.Gauge, error) {
	rs.Lock()
	defer rs.Unlock()
	return rs.setGauge(ctx, name, value)
}

// GetAll returns all saved metrics.
func (rs *ramRepository) GetAll(ctx context.Context) ([]metrics.Metric, error) {
	rs.RLock()
	defer rs.RUnlock()
	metricSlice := make([]metrics.Metric, 0, len(rs.counterStorage)+len(rs.gaugeStorage))

	for _, counter := range rs.counterStorage {
		metricSlice = append(metricSlice, counter)
	}
	for _, gauge := range rs.gaugeStorage {
		metricSlice = append(metricSlice, gauge)
	}
	return metricSlice, nil
}

// AddCounters increase each metrics.Counter on value in slice and return slice of result.
func (rs *ramRepository) AddCounters(ctx context.Context, slice []metrics.Counter) ([]metrics.Counter, error) {
	rs.Lock()
	defer rs.Unlock()
	for idx, counter := range slice {
		updatedCounter, err := rs.addCounter(ctx, counter.GetName(), counter.Value())
		if err != nil {
			return nil, err
		}
		slice[idx] = updatedCounter
	}
	return slice, nil
}

// SetGauges set new value for each metrics.Gauge in slice and return the result slice.
func (rs *ramRepository) SetGauges(ctx context.Context, slice []metrics.Gauge) ([]metrics.Gauge, error) {
	rs.Lock()
	defer rs.Unlock()
	for idx, gauge := range slice {
		updatedGauge, err := rs.setGauge(ctx, gauge.GetName(), gauge.Value())
		if err != nil {
			return nil, err
		}
		slice[idx] = updatedGauge
	}
	return slice, nil
}

// Shutdown .
func (rs *ramRepository) Shutdown() error {
	return nil
}

// Ping .
func (rs *ramRepository) Ping() error {
	return nil
}
