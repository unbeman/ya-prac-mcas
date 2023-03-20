package storage

import (
	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type Repository interface {
	AddCounter(name string, delta int64) metrics.Counter
	GetCounter(name string) metrics.Counter

	SetGauge(name string, value float64) metrics.Gauge
	GetGauge(name string) metrics.Gauge

	GetAll() []metrics.Metric

	Ping() error
	Shutdown() error
}

func GetRepository(cfg configs.RepositoryConfig) (Repository, error) { //cfg configs.RepositoryConfig
	switch {
	case cfg.PG != nil:
		return NewPostgresRepository(*cfg.PG)
	default:
		return NewRAMRepository(), nil
	}
}
