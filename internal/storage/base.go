package storage

import (
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

//var ErrNotFound = errors.New("not found")
//var ErrStorage = errors.New("storage error")

type Repository interface { //TODO: rename
	//Get(id string) (metrics.Metric, error)
	//Update(id string, metric metrics.Metric) error
	//Add(id string, metric metrics.Metric) error
	//Delete(id string) error
	UpdateCounterRepo(metric metrics.Counter)
	UpdateGaugeRepo(metric metrics.Gauge)
	GetAll() (map[string]metrics.Metric, bool)
	GetCounter(id string) (metrics.Counter, bool)
	GetGauge(id string) (metrics.Gauge, bool)
	GetAllCounter() (map[string]metrics.Counter, bool)
	GetAllGauge() (map[string]metrics.Gauge, bool)
}
