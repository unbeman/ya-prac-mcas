package parser

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

const (
	PType  = "type"
	PName  = "name"
	PValue = "value"
)

func ParseURI(request *http.Request, keys ...string) (*MetricParams, error) {
	// type и name обязательны
	params := &MetricParams{}
	var valueStr string
	for _, key := range keys {
		value := chi.URLParam(request, key)
		if len(value) == 0 {
			return nil, fmt.Errorf("ParseURI: %v is %w", key, ErrInvalidValue)
		}
		switch key {
		case PType:
			err := checkType(value)
			if err != nil {
				return nil, fmt.Errorf("ParseURI: %v is %w", key, err)
			}
			params.Type = value
		case PName:
			if err := checkName(value); err != nil {
				return nil, fmt.Errorf("ParseURI: %v is %w", PName, err)
			}
			params.Name = value
		case PValue:
			valueStr = value
		default:
			return nil, fmt.Errorf("ParseURI: unknown parameter")
		}
	}
	if len(valueStr) != 0 { //TODO: wrap to func
		switch params.Type {
		case metrics.GaugeType:
			gValue, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				log.Printf("ParamsToMetric: %v \n", err)
				return nil, fmt.Errorf("%w: %v", ErrInvalidValue, valueStr)
			}
			params.ValueGauge = &gValue
		case metrics.CounterType:
			cValue, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				log.Printf("ParamsToMetric: %v \n", err)
				return nil, fmt.Errorf("%w: %v", ErrInvalidValue, valueStr)
			}
			params.ValueCounter = &cValue
		}
	}
	return params, nil
}
