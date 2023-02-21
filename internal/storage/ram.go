package storage

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type RAMRepository struct {
	counter CounterStorager
	gauge   GaugeStorager
}

func NewRAMRepository() *RAMRepository {
	return &RAMRepository{counter: NewCounterRAMStorage(), gauge: NewGaugeRAMStorage()}
}

func (rs *RAMRepository) GetMetric(typeM, id string) (metrics.Metric, error) {
	switch typeM {
	case metrics.CounterType:
		counter := rs.counter.Get(id)
		if counter == nil {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, id)
		}
		return counter, nil
	case metrics.GaugeType:
		gauge := rs.gauge.Get(id)
		if gauge == nil {
			return nil, fmt.Errorf("%w: %v", ErrNotFound, id)
		}
		return gauge, nil
	default:
		return nil, fmt.Errorf("%w: %v", ErrInvalidType, typeM)
	}
}

func (rs *RAMRepository) SetMetric(typeM, id, value string) error {
	switch typeM {
	case metrics.CounterType:
		cValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidValue, value)
		}
		rs.counter.Set(id, cValue)
	case metrics.GaugeType:
		gValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidValue, value)
		}
		rs.gauge.Set(id, gValue)
	default:
		return fmt.Errorf("%w: %v", ErrInvalidType, typeM)
	}
	return nil
}

// не уверена что так правильно делать, в итоге поучается по три прохода по элементам слайсов
func (rs *RAMRepository) GetAll() []metrics.Metric {
	counterSlice := rs.counter.GetAll()
	gaugeSlice := rs.gauge.GetAll()
	metricSlice := make([]metrics.Metric, 0, len(counterSlice)+len(gaugeSlice))
	for _, counter := range counterSlice {
		metricSlice = append(metricSlice, counter)
	}
	for _, gauge := range gaugeSlice {
		metricSlice = append(metricSlice, gauge)
	}
	return metricSlice
}

type counterRAMStorage struct {
	sync.RWMutex
	storage map[string]metrics.Counter
}

func NewCounterRAMStorage() *counterRAMStorage {
	return &counterRAMStorage{storage: map[string]metrics.Counter{}}
}

func (cs *counterRAMStorage) get(id string) metrics.Counter {
	value, ok := cs.storage[id]
	if !ok {
		return nil
	}
	return value
}

func (cs *counterRAMStorage) Get(id string) metrics.Counter {
	cs.RLock()
	defer cs.RUnlock()
	return cs.get(id)
}

func (cs *counterRAMStorage) GetAll() []metrics.Counter {
	cs.RLock()
	defer cs.RUnlock()
	mSlice := make([]metrics.Counter, 0, len(cs.storage))
	for _, counter := range cs.storage {
		mSlice = append(mSlice, counter)
	}
	return mSlice
}

func (cs *counterRAMStorage) Set(id string, value int64) {
	cs.Lock()
	defer cs.Unlock()
	counter := cs.get(id)
	if counter == nil {
		counter = metrics.NewCounter(id)
		cs.storage[id] = counter
	}
	counter.Add(value)
}

type gaugeRAMStorage struct {
	sync.RWMutex
	storage map[string]metrics.Gauge
}

func NewGaugeRAMStorage() *gaugeRAMStorage {
	return &gaugeRAMStorage{storage: map[string]metrics.Gauge{}}
}

func (gs *gaugeRAMStorage) get(id string) metrics.Gauge {
	value, ok := gs.storage[id]
	if !ok {
		return nil
	}
	return value
}

func (gs *gaugeRAMStorage) Get(id string) metrics.Gauge {
	gs.RLock()
	defer gs.RUnlock()
	return gs.get(id)
}

func (gs *gaugeRAMStorage) GetAll() []metrics.Gauge {
	gs.RLock()
	defer gs.RUnlock()
	mSlice := make([]metrics.Gauge, 0, len(gs.storage))
	for _, gauge := range gs.storage {
		mSlice = append(mSlice, gauge)
	}
	return mSlice
}

func (gs *gaugeRAMStorage) Set(id string, value float64) {
	gs.Lock()
	defer gs.Unlock()
	gauge := gs.get(id)
	if gauge == nil {
		gauge = metrics.NewGauge(id)
		gs.storage[id] = gauge
	}
	gauge.Set(value)
}
