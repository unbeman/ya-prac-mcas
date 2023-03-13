package storage

import (
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type Repository interface {
	AddCounter(name string, delta int64) metrics.Counter
	GetCounter(name string) metrics.Counter

	SetGauge(name string, value float64) metrics.Gauge
	GetGauge(name string) metrics.Gauge

	GetAll() []metrics.Metric
}

func GetRepository() Repository { //cfg configs.RepositoryConfig
	switch {
	default:
		return NewRAMRepository()
	}
}
