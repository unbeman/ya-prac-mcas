package storage

import (
	"errors"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

var ErrNotFound = errors.New("not found")

type Repository interface { //TODO: rename
	//Get(id string) (metrics.Metric, error)
	Update(id string, metric metrics.Metric) error
	//Add(id string, metric metrics.Metric) error
	//Delete(id string) error
}
