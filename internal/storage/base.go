package storage

import (
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
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

type Repository struct {
	Gauge   GaugeStorager
	Counter CounterStorager
}

func NewRepository(gaugeRepo GaugeStorager, counterRepo CounterStorager) *Repository {
	return &Repository{Gauge: gaugeRepo, Counter: counterRepo}
}
