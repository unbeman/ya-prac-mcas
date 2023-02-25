package parser

import (
	"fmt"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type JSONMetric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func MetricToJSON(m metrics.Metric) *JSONMetric {
	jM := &JSONMetric{ID: m.GetName(), MType: m.GetType()}
	switch m.GetType() {
	case metrics.CounterType:
		v := m.(metrics.Counter).Value()
		jM.Delta = &v
	case metrics.GaugeType:
		v := m.(metrics.Gauge).Value()
		jM.Value = &v
	}
	return jM
}

func ParseJSON(jm *JSONMetric, needValue bool) (*MetricParams, error) {
	if err := checkName(jm.ID); err != nil {
		return nil, fmt.Errorf("ParseJSON: ID is %w", err)
	}
	err := checkType(jm.MType)
	if err != nil {
		return nil, err
	}
	var mP *MetricParams
	switch jm.MType {
	case metrics.CounterType:
		if needValue && jm.Delta == nil {
			return nil, fmt.Errorf("ParseJSON: Delta is %w", ErrInvalidValue)
		}
		mP = &MetricParams{Name: jm.ID, Type: metrics.CounterType, ValueCounter: jm.Delta}
	case metrics.GaugeType:
		if needValue && jm.Value == nil {
			return nil, fmt.Errorf("ParseJSON: Value is %w", ErrInvalidValue)
		}
		mP = &MetricParams{Name: jm.ID, Type: metrics.GaugeType, ValueGauge: jm.Value}
	}
	return mP, nil
}
