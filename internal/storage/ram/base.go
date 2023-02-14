package ram

import (
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"sync"
)

type ramStorage struct { //TODO: split into two repos for each metric type?
	sync.RWMutex
	counterRepo map[string]metrics.Counter
	gaugeRepo   map[string]metrics.Gauge
}

func (rs *ramStorage) GetCounter(id string) (metrics.Counter, bool) {
	rs.RLock()
	defer rs.RUnlock()
	if val, ok := rs.counterRepo[id]; ok {
		return val, ok
	}
	return nil, false
}

func (rs *ramStorage) GetGauge(id string) (metrics.Gauge, bool) {
	rs.RLock()
	defer rs.RUnlock()
	if val, ok := rs.gaugeRepo[id]; ok {
		return val, ok
	}
	return nil, false
}

func (rs *ramStorage) GetAllGauge() (map[string]metrics.Gauge, bool) {
	rs.RLock()
	defer rs.RUnlock()
	return rs.gaugeRepo, true
}

func (rs *ramStorage) GetAllCounter() (map[string]metrics.Counter, bool) {
	rs.RLock()
	defer rs.RUnlock()
	return rs.counterRepo, true
}

func (rs *ramStorage) GetAll() (map[string]metrics.Metric, bool) {
	rs.RLock()
	defer rs.RUnlock()
	metricMap := map[string]metrics.Metric{}
	for k, v := range rs.gaugeRepo {
		metricMap[k] = v
	}
	for k, v := range rs.counterRepo {
		metricMap[k] = v
	}
	return metricMap, true
}

func (rs *ramStorage) UpdateCounterRepo(counter metrics.Counter) {
	rs.Lock()
	defer rs.Unlock()
	rs.counterRepo[counter.GetName()] = counter
}

func (rs *ramStorage) UpdateGaugeRepo(gauge metrics.Gauge) {
	rs.Lock()
	defer rs.Unlock()
	rs.gaugeRepo[gauge.GetName()] = gauge
}

//func (rs *ramStorage) Update(id string, metric metrics.Metric) error {
//	rs.Lock()
//	defer rs.Unlock()
//	rs.repo[id] = metric
//	return nil
//}

func NewRAMStorage() *ramStorage {
	return &ramStorage{gaugeRepo: map[string]metrics.Gauge{}, counterRepo: map[string]metrics.Counter{}}
}
