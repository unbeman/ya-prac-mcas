package storage

import (
	"sync"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type counterRAMStorage struct {
	sync.RWMutex
	storage map[string]metrics.Counter
}

func NewCounterRAMStorage() *counterRAMStorage {
	return &counterRAMStorage{storage: map[string]metrics.Counter{}}
}

func (cs *counterRAMStorage) Get(id string) metrics.Counter {
	cs.RLock()
	defer cs.RUnlock()
	if value, ok := cs.storage[id]; !ok {
		return nil
	} else {
		return value
	}
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
	counter := cs.Get(id)
	cs.Lock()
	defer cs.Unlock()
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

func (cs *gaugeRAMStorage) Get(id string) metrics.Gauge {
	cs.RLock()
	defer cs.RUnlock()
	if value, ok := cs.storage[id]; !ok {
		return nil
	} else {
		return value
	}
}

func (cs *gaugeRAMStorage) GetAll() []metrics.Gauge {
	cs.RLock()
	defer cs.RUnlock()
	mSlice := make([]metrics.Gauge, 0, len(cs.storage))
	for _, gauge := range cs.storage {
		mSlice = append(mSlice, gauge)
	}
	return mSlice
}

func (cs *gaugeRAMStorage) Set(id string, value float64) {
	gauge := cs.Get(id)
	cs.Lock()
	defer cs.Unlock()
	if gauge == nil {
		gauge = metrics.NewGauge(id)
		cs.storage[id] = gauge
	}
	gauge.Set(value)
}
