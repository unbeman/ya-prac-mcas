package storage

import (
	"sync"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type ramRepository struct {
	sync.RWMutex
	counterStorage map[string]metrics.Counter
	gaugeStorage   map[string]metrics.Gauge
}

func (rs *ramRepository) Shutdown() error {
	return nil
}

func (rs *ramRepository) Ping() error {
	return nil
}

func NewRAMRepository() *ramRepository {
	return &ramRepository{
		counterStorage: map[string]metrics.Counter{},
		gaugeStorage:   map[string]metrics.Gauge{},
	}
}

func (rs *ramRepository) getCounter(name string) metrics.Counter {
	value, ok := rs.counterStorage[name]
	if !ok {
		return nil
	}
	return value
}

func (rs *ramRepository) GetCounter(name string) metrics.Counter {
	rs.RLock()
	defer rs.RUnlock()
	return rs.getCounter(name)
}

func (rs *ramRepository) AddCounter(name string, value int64) metrics.Counter {
	rs.Lock()
	defer rs.Unlock()
	counter := rs.getCounter(name)
	if counter == nil {
		counter = metrics.NewCounter(name)
		rs.counterStorage[name] = counter
	}
	counter.Add(value)
	return counter
}

func (rs *ramRepository) getGauge(name string) metrics.Gauge {
	value, ok := rs.gaugeStorage[name]
	if !ok {
		return nil
	}
	return value
}

func (rs *ramRepository) GetGauge(name string) metrics.Gauge {
	rs.RLock()
	defer rs.RUnlock()
	return rs.getGauge(name)
}

func (rs *ramRepository) SetGauge(name string, value float64) metrics.Gauge {
	rs.Lock()
	defer rs.Unlock()
	gauge := rs.getGauge(name)
	if gauge == nil {
		gauge = metrics.NewGauge(name)
		rs.gaugeStorage[name] = gauge
	}
	gauge.Set(value)
	return gauge
}

func (rs *ramRepository) GetAll() []metrics.Metric {
	rs.RLock()
	defer rs.RUnlock()
	metricSlice := make([]metrics.Metric, 0, len(rs.counterStorage)+len(rs.gaugeStorage))

	for _, counter := range rs.counterStorage {
		metricSlice = append(metricSlice, counter)
	}
	for _, gauge := range rs.gaugeStorage {
		metricSlice = append(metricSlice, gauge)
	}
	return metricSlice
}
