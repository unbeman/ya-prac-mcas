package parser

import (
	"errors"
	"fmt"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"strconv"
	"strings"
)

func FormatURL(addr string, m metrics.Metric) string {
	return fmt.Sprintf("http://%v/update/%v/%v/%v", addr, m.GetType(), m.GetName(), m.GetValue())
}

var ErrNotEnoughParams = errors.New("not enough parameters")
var ErrParse = errors.New("can't parse")
var ErrUnknownType = errors.New("unknown type")

func ParseMetric(url string) (metrics.Metric, error) {
	params := strings.Split(url, "/")
	if len(params) != 5 {
		return nil, ErrNotEnoughParams
	}
	typeName := params[2]
	switch typeName {
	case "gauge":
		value, err := strconv.ParseFloat(params[4], 64)
		if err != nil {
			return nil, ErrParse
		}
		g := metrics.NewGaugeWithValue(params[3], value)
		return g, nil
	case "counter":
		value, err := strconv.ParseInt(params[4], 10, 64)
		if err != nil {
			return nil, ErrParse
		}
		c := metrics.NewCounterWithValue(params[3], value)
		return c, nil
	default:
		return nil, ErrUnknownType
	}
}
