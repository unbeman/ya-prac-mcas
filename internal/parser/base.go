package parser

import (
	"fmt"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

const (
	PName  string = "name"
	PType  string = "type"
	PValue string = "value"
)

type MetricParams struct { //TODO: make builder .TypeAndName() .Value()
	Name         string
	Type         string
	ValueCounter *int64
	ValueGauge   *float64
	Hash         string
}

func checkType(typeStr string) error {
	switch typeStr {
	case metrics.GaugeType, metrics.CounterType:
		return nil
	default:
		return fmt.Errorf("checkType: %v - %w", typeStr, ErrInvalidType)
	}
}

func checkName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("checkType: %v - %w", name, ErrInvalidValue)
	}
	return nil
}
