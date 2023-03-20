package parser

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type JSONMetric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func (jm *JSONMetric) String() string {
	var delta, value = "nil", "nil"
	if jm.Delta != nil {
		delta = fmt.Sprintf("%v", *jm.Delta)
	}
	if jm.Value != nil {
		value = fmt.Sprintf("%v", *jm.Value)
	}
	return fmt.Sprintf("{ID:%v; MType:%v; Delta:%v; Value:%v};", jm.ID, jm.MType, delta, value)
}

func (jm *JSONMetric) Decode(data io.Reader) error {
	if err := json.NewDecoder(data).Decode(&jm); err != nil {
		return err
	}
	return nil
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

func ParseJSON(jm *JSONMetric, requiredKeys ...string) (*MetricParams, error) {
	//TODO: гарантированной получать type и name, без type value не получать
	params := &MetricParams{}
	for _, key := range requiredKeys {
		switch key {
		case PType:
			err := checkType(jm.MType)
			if err != nil {
				return nil, err
			}
			params.Type = jm.MType
		case PName:
			if err := checkName(jm.ID); err != nil {
				return nil, fmt.Errorf("ParseJSON: ID is %w", err)
			}
			params.Name = jm.ID
		case PValue:
			switch jm.MType {
			case metrics.CounterType:
				if jm.Delta == nil {
					return nil, fmt.Errorf("ParseJSON: Delta is %w", ErrInvalidValue)
				}
				params.ValueCounter = jm.Delta
			case metrics.GaugeType:
				if jm.Value == nil {
					return nil, fmt.Errorf("ParseJSON: Value is %w", ErrInvalidValue)
				}
				params.ValueGauge = jm.Value
			}
		}
	}
	params.Hash = jm.Hash
	return params, nil
}
