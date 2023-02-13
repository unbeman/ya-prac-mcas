package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"testing"
)

func TestFormatURL(t *testing.T) {
	type args struct {
		addr string
		m    metrics.Metric
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Good format",
			args: args{addr: "localhost:8080", m: metrics.NewCounterWithValue("CatCount", 3)},
			want: "http://localhost:8080/update/counter/CatCount/3",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatURL(tt.args.addr, tt.args.m))
		})
	}
}

func TestParseMetric(t *testing.T) {

	tests := []struct {
		name    string
		url     string
		want    metrics.Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{name: "Good parse counter",
			url:     "localhost:8080/update/counter/CatCount/3",
			want:    metrics.NewCounterWithValue("CatCount", 3),
			wantErr: assert.NoError,
		},
		{name: "Good parse gauge",
			url:     "localhost:8080/update/gauge/CatGauge/3.0",
			want:    metrics.NewGaugeWithValue("CatGauge", 3.0),
			wantErr: assert.NoError,
		},
		{name: "Err parse value count",
			url:     "localhost:8080/update/counter/CatCount/3.abba",
			want:    nil,
			wantErr: assert.Error,
		},
		{name: "Err parse value gauge",
			url:     "localhost:8080/update/gauge/CatGauge/3.abba",
			want:    nil,
			wantErr: assert.Error,
		},
		{name: "Err parse unknown metric type",
			url:     "localhost:8080/update/dog/CatGauge/3.0",
			want:    nil,
			wantErr: assert.Error,
		},
		{name: "Err invalid url format",
			url:     "localhost:8080/update/0/counter/Cat",
			want:    nil,
			wantErr: assert.Error,
		},
		{name: "Err not enough params",
			url:     "localhost:8080/update/counter/",
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMetric(tt.url)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseMetric(%v)", tt.url)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseMetric(%v)", tt.url)
		})
	}
}
