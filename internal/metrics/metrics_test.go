package metrics

import (
	"reflect"
	"testing"
)

func TestNewMetricFromParams(t *testing.T) {

	tests := []struct {
		name  string
		input Params
		want  Metric
	}{
		{
			name: "good Gauge",
			input: Params{
				Name:       "WaterPercent",
				Type:       "gauge",
				ValueGauge: func(n float64) *float64 { return &n }(0.35),
			},
			want: NewGauge("WaterPercent", 0.35),
		},
		{
			name: "good Counter",
			input: Params{
				Name:         "Dog",
				Type:         "counter",
				ValueCounter: func(n int64) *int64 { return &n }(10),
			},
			want: NewCounter("Dog", 10),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMetricFromParams(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetricFromParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
