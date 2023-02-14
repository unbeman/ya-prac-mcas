package parser

import (
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
			args: args{addr: "localhost:8080", m: func() metrics.Metric {
				c := metrics.NewCounter("CatCount")
				c.Add(3)
				return c
			}()},
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
