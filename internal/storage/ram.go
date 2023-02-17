package storage

import (
	"sync"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type counterRamStorage struct {
	sync.RWMutex
	storage map[string]metrics.Counter
}

func NewCounterRamStorage() *counterRamStorage {
	return &counterRamStorage{storage: map[string]metrics.Counter{}}
}

func (cs *counterRamStorage) Get(id string) metrics.Counter {
	cs.RLock()
	defer cs.RUnlock()
	if value, ok := cs.storage[id]; !ok {
		return nil
	} else {
		return value
	}
}

func (cs *counterRamStorage) GetAll() []metrics.Counter {
	cs.RLock()
	defer cs.RUnlock()
	mSlice := make([]metrics.Counter, 0, len(cs.storage))
	for _, counter := range cs.storage {
		mSlice = append(mSlice, counter)
	}
	return mSlice
}

func (cs *counterRamStorage) Set(id string, value int64) { //TODO: return metrics.Counter?
	counter := cs.Get(id)
	cs.Lock()
	defer cs.Unlock()
	if counter == nil {
		counter = metrics.NewCounter(id)
		cs.storage[id] = counter
	}
	counter.Add(value)
}

type gaugeRamStorage struct {
	sync.RWMutex
	storage map[string]metrics.Gauge
}

func NewGaugeRamStorage() *gaugeRamStorage {
	return &gaugeRamStorage{storage: map[string]metrics.Gauge{}}
}

func (cs *gaugeRamStorage) Get(id string) metrics.Gauge {
	cs.RLock()
	defer cs.RUnlock()
	if value, ok := cs.storage[id]; !ok {
		return nil
	} else {
		return value
	}
}

func (cs *gaugeRamStorage) GetAll() []metrics.Gauge {
	cs.RLock()
	defer cs.RUnlock()
	mSlice := make([]metrics.Gauge, 0, len(cs.storage))
	for _, gauge := range cs.storage {
		mSlice = append(mSlice, gauge)
	}
	return mSlice
}

func (cs *gaugeRamStorage) Set(id string, value float64) { //TODO: return metrics.Counter?
	gauge := cs.Get(id)
	cs.Lock()
	defer cs.Unlock()
	if gauge == nil {
		gauge = metrics.NewGauge(id)
		cs.storage[id] = gauge
	}
	gauge.Set(value)
}
