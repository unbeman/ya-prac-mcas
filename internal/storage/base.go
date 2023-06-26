package storage

import (
	"context"
	"errors"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	AddCounter(ctx context.Context, name string, delta int64) (metrics.Counter, error)
	AddCounters(ctx context.Context, slice []metrics.Counter) error
	GetCounter(ctx context.Context, name string) (metrics.Counter, error)

	SetGauge(ctx context.Context, name string, value float64) (metrics.Gauge, error)
	SetGauges(ctx context.Context, slice []metrics.Gauge) error
	GetGauge(ctx context.Context, name string) (metrics.Gauge, error)

	GetAll(ctx context.Context) ([]metrics.Metric, error)

	Ping() error
	Shutdown() error
}

func GetRepository(cfg configs.RepositoryConfig) (Repository, error) { //cfg configs.RepositoryConfig
	switch {
	case cfg.PG != nil:
		return NewPostgresRepository(*cfg.PG)
	case cfg.RAMWithBackup != nil:
		return NewRAMBackupRepository(cfg.RAMWithBackup)
	default:
		return NewRAMRepository(), nil
	}
}
