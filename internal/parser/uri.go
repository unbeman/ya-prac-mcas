package parser

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

func ParseURI(request *http.Request, requiredKeys ...string) (*MetricParams, error) {
	//TODO: гарантированной получать type и name, без type value не получать
	params := &MetricParams{}
	for _, key := range requiredKeys {
		value := chi.URLParam(request, key)
		//if len(value) == 0 {
		//	return nil, fmt.Errorf("ParseURI: %v is %w", key, ErrInvalidValue)
		//}
		switch key {
		case PType:
			err := checkType(value)
			if err != nil {
				return nil, fmt.Errorf("ParseURI: %v - %w", key, err)
			}
			params.Type = value
		case PName:
			if err := checkName(value); err != nil {
				return nil, fmt.Errorf("ParseURI: %v - %w", key, err)
			}
			params.Name = value
		case PValue:
			switch params.Type { // unsafe
			case metrics.GaugeType:
				gValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Printf("ParseURI: %v \n", err)
					return nil, fmt.Errorf("ParseURI: %v = %v - %w", key, value, ErrInvalidValue)
				}
				params.ValueGauge = &gValue
			case metrics.CounterType:
				cValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					log.Printf("ParseURI: %v \n", err)
					return nil, fmt.Errorf("ParseURI: %v = %v - %w", key, value, ErrInvalidValue)
				}
				params.ValueCounter = &cValue
			}
		default:
			return nil, fmt.Errorf("ParseURI: unknown parameter")
		}
	}
	return params, nil
}
