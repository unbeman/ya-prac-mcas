package storage

import (
	"errors"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

var (
	ErrInvalidType  = errors.New("invalid type")
	ErrInvalidValue = errors.New("invalid value")
	ErrNotFound     = errors.New("not found")
)

type CounterStorager interface {
	Set(id string, value int64) // error?
	Get(id string) metrics.Counter
	GetAll() []metrics.Counter
}

type GaugeStorager interface {
	Set(id string, value float64) // error?
	Get(id string) metrics.Gauge
	GetAll() []metrics.Gauge
}

type Repository interface {
	SetMetric(typeM, id string, value string) error
	GetMetric(typeM, id string) (metrics.Metric, error)
	GetAll() []metrics.Metric
}
