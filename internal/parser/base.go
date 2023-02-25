package parser

import (
	"fmt"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type MetricParams struct {
	Name         string
	Type         string
	ValueCounter *int64
	ValueGauge   *float64
}

//func getType(typeStr string) (metrics.TypeMetric, error) {
//	var tM metrics.TypeMetric
//	switch typeStr {
//	case metrics.CounterTypeStr:
//		tM = metrics.CounterType
//		return tM, nil
//	case metrics.GaugeTypeStr:
//		tM = metrics.GaugeType
//		return tM, nil
//	default:
//		return tM, fmt.Errorf("getType: %v is %w", typeStr, ErrInvalidType)
//	}
//}

func checkType(typeStr string) error {
	switch typeStr {
	case metrics.GaugeType, metrics.CounterType:
		return nil
	default:
		return fmt.Errorf("checkType: %v is %w", typeStr, ErrInvalidType)
	}
}

func checkName(name string) error {
	if len(name) == 0 {
		return ErrInvalidValue
	}
	return nil
}
