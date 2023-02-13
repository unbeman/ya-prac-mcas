package parser

import (
	"fmt"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"strconv"
	"strings"
)

func FormatURL(addr string, m metrics.Metric) string {
	return fmt.Sprintf("http://%v/update/%v/%v/%v", addr, m.GetType(), m.GetName(), m.GetValue())
}

func ParseMetric(url string) (metrics.Metric, error) {
	params := strings.Split(url, "/") //TODO: assert len
	typeName := params[2]
	switch typeName {
	case "gauge":
		value, err := strconv.ParseFloat(params[4], 64)
		if err != nil {
			return nil, fmt.Errorf("ParseMetric: cant's parse gauge %v to float64. %v", params[4], err)
		}
		g := metrics.NewGaugeWithValue(params[3], value)
		return g, nil
	case "counter":
		value, err := strconv.ParseInt(params[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ParseMetric: cant's parse counter %v to int64. %v", params[4], err)
		}
		c := metrics.NewCounterWithValue(params[3], value)
		return c, nil
	default:
		return nil, fmt.Errorf("ParseMetric: invalid metric type")
	}
}
