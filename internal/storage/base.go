// Package storage has storage interface Repository, and its implementations.
package storage

import (
	"context"
	"errors"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

// ErrNotFound is returned by Repository methods when entity didn't find.
var ErrNotFound = errors.New("not found")

// Repository describes the storage usage.
type Repository interface {
	AddCounter(ctx context.Context, name string, delta int64) (metrics.Counter, error)
	AddCounters(ctx context.Context, slice []metrics.Counter) ([]metrics.Counter, error)
	GetCounter(ctx context.Context, name string) (metrics.Counter, error)

	SetGauge(ctx context.Context, name string, value float64) (metrics.Gauge, error)
	SetGauges(ctx context.Context, slice []metrics.Gauge) ([]metrics.Gauge, error)
	GetGauge(ctx context.Context, name string) (metrics.Gauge, error)

	GetAll(ctx context.Context) ([]metrics.Metric, error)

	Ping() error
	Shutdown() error
}

// GetRepository return Repository implementation depending on the config.
func GetRepository(cfg configs.RepositoryConfig) (Repository, error) {
	switch {
	case cfg.PG != nil:
		return NewPostgresRepository(*cfg.PG)
	case cfg.RAMWithBackup != nil:
		return NewRAMBackupRepository(cfg.RAMWithBackup)
	default:
		return NewRAMRepository(), nil
	}
}
