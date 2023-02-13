package ram

import (
	"fmt"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	"sync"
)

type ramStorage struct { //TODO: threadsafe
	sync.RWMutex
	repo map[string]metrics.Metric
}

func (rs *ramStorage) Get(id string) (metrics.Metric, error) {
	rs.RLock()
	defer rs.RUnlock()
	if val, ok := rs.repo[id]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("%q: %w", id, storage.ErrNotFound)
}

func (rs *ramStorage) Update(id string, metric metrics.Metric) error {
	rs.Lock()
	defer rs.Unlock()
	rs.repo[id] = metric
	return nil
}

func NewRAMStorage() *ramStorage {
	return &ramStorage{repo: map[string]metrics.Metric{}}
}
