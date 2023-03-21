package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type ramRepository struct {
	sync.RWMutex
	counterStorage map[string]metrics.Counter
	gaugeStorage   map[string]metrics.Gauge
}

func NewRAMRepository() *ramRepository {
	return &ramRepository{
		counterStorage: map[string]metrics.Counter{},
		gaugeStorage:   map[string]metrics.Gauge{},
	}
}

func (rs *ramRepository) getCounter(name string) (metrics.Counter, error) {
	value, ok := rs.counterStorage[name]
	if !ok {
		return nil, fmt.Errorf("counter (%v) %w", name, ErrNotFound)
	}
	return value, nil
}

func (rs *ramRepository) GetCounter(ctx context.Context, name string) (metrics.Counter, error) {
	rs.RLock()
	defer rs.RUnlock()
	return rs.getCounter(name)
}

func (rs *ramRepository) AddCounter(ctx context.Context, name string, value int64) (metrics.Counter, error) {
	rs.Lock()
	defer rs.Unlock()
	counter, err := rs.getCounter(name)
	if errors.Is(err, ErrNotFound) {
		counter = metrics.NewCounter(name)
		rs.counterStorage[name] = counter
	}
	counter.Add(value)
	return counter, nil
}

func (rs *ramRepository) getGauge(name string) (metrics.Gauge, error) {
	value, ok := rs.gaugeStorage[name]
	if !ok {
		return nil, fmt.Errorf("gauge (%v) %w", name, ErrNotFound)
	}
	return value, nil
}

func (rs *ramRepository) GetGauge(ctx context.Context, name string) (metrics.Gauge, error) {
	rs.RLock()
	defer rs.RUnlock()
	return rs.getGauge(name)
}

func (rs *ramRepository) SetGauge(ctx context.Context, name string, value float64) (metrics.Gauge, error) {
	rs.Lock()
	defer rs.Unlock()
	gauge, err := rs.getGauge(name)
	if errors.Is(err, ErrNotFound) {
		gauge = metrics.NewGauge(name)
		rs.gaugeStorage[name] = gauge
	}
	gauge.Set(value)
	return gauge, nil
}

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

func (rs *ramRepository) AddCounters(ctx context.Context, slice []metrics.Counter) error {
	var err error
	for _, counter := range slice {
		_, err = rs.AddCounter(ctx, counter.GetName(), counter.Value())
		if err != nil {
			return err
		}
	}
	return nil
}

func (rs *ramRepository) SetGauges(ctx context.Context, slice []metrics.Gauge) error {
	var err error
	for _, gauge := range slice {
		_, err = rs.SetGauge(ctx, gauge.GetName(), gauge.Value())
		if err != nil {
			return err
		}
	}
	return nil
}

func (rs *ramRepository) Shutdown() error {
	return nil
}

func (rs *ramRepository) Ping() error {
	return nil
}
