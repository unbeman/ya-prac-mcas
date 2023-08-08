package controller

import (
	"context"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type Controller struct {
	repository storage.Repository
	hashKey    []byte
}

func NewController(repo storage.Repository, hashKey string) *Controller {
	return &Controller{repository: repo, hashKey: []byte(hashKey)}
}

func (c Controller) GetAll(ctx context.Context) ([]metrics.Metric, error) {
	return c.repository.GetAll(ctx)
}

func (c Controller) Ping(ctx context.Context) error {
	return c.repository.Ping(ctx)
}

func (c Controller) GetMetric(ctx context.Context, params metrics.Params) (metrics.Metric, error) { //TODO: controller layer
	var (
		metric metrics.Metric
		err    error
	)
	switch params.Type {
	case metrics.GaugeType:
		metric, err = c.repository.GetGauge(ctx, params.Name)
	case metrics.CounterType:
		metric, err = c.repository.GetCounter(ctx, params.Name)
	}
	return metric, err
}

func (c Controller) UpdateMetric(ctx context.Context, params metrics.Params) (metrics.Metric, error) { //TODO: controller layer
	var (
		err    error
		metric = metrics.NewMetricFromParams(params)
	)

	if !c.IsValidHash(params.Hash, metric) {
		return nil, ErrInvalidHash
	}
	switch params.Type {
	case metrics.GaugeType:
		metric, err = c.repository.SetGauge(ctx, params.Name, *params.ValueGauge)
	case metrics.CounterType:
		metric, err = c.repository.AddCounter(ctx, params.Name, *params.ValueCounter)
	}
	return metric, err
}

func (c Controller) UpdateMetrics(
	ctx context.Context,
	paramsSlice metrics.ParamsSlice) (metrics.ParamsSlice, error) {

	gauges := make([]metrics.Gauge, 0)
	counters := make([]metrics.Counter, 0)
	for _, params := range paramsSlice {
		metric := metrics.NewMetricFromParams(params)

		if !c.IsValidHash(params.Hash, metric) {
			return nil, ErrInvalidHash
		}

		switch metric.GetType() {
		case metrics.GaugeType:
			gauges = append(gauges, metric.(metrics.Gauge))
		case metrics.CounterType:
			counters = append(counters, metric.(metrics.Counter))
		}
	}

	metricsParams := make(metrics.ParamsSlice, 0, len(gauges)+len(counters))

	if len(gauges) > 0 {
		updatedGauges, err := c.repository.SetGauges(ctx, gauges)
		if err != nil {
			return nil, err
		}

		for _, gauge := range updatedGauges {
			gp := gauge.ToParams()
			gp.Hash = c.GetHash(gauge)
			metricsParams = append(metricsParams, gp)
		}
	}

	if len(counters) > 0 {
		updatedCounters, err := c.repository.AddCounters(ctx, counters)
		if err != nil {
			return nil, err
		}

		for _, counter := range updatedCounters {
			cp := counter.ToParams()
			cp.Hash = c.GetHash(counter)
			metricsParams = append(metricsParams, cp)
		}

	}
	return metricsParams, nil
}

func (c Controller) IsValidHash(hash string, metric metrics.Metric) bool {
	if !c.isKeySet() { // ключа нет, проверять не нужно
		return true
	}
	if !isHashSet(hash) {
		return false
	}

	return hash == metric.Hash(c.hashKey)
}

func (c Controller) GetHash(metric metrics.Metric) string {
	if !c.isKeySet() {
		return ""
	}
	return metric.Hash(c.hashKey)
}

func (c Controller) isKeySet() bool {
	return len(c.hashKey) > 0
}

func isHashSet(hash string) bool {
	return len(hash) > 0
}
