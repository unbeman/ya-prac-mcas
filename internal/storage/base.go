package storage

import (
	"context"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type Repository interface {
	AddCounter(name string, delta int64) metrics.Counter
	GetCounter(name string) metrics.Counter

	SetGauge(name string, value float64) metrics.Gauge
	GetGauge(name string) metrics.Gauge

	GetAll() []metrics.Metric

	Load() error
	Save() error
	RunSaver(ctx context.Context)
}

func GetRepository(cfg configs.RepositoryConfig) Repository {
	switch {
	case cfg.FileStorage != nil:
		return NewFileRepository(*cfg.FileStorage)
	default:
		return NewRAMRepository()
	}
}
