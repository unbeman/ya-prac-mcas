package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

func Test_ramRepository_GetCounter(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}

	type want struct {
		metric metrics.Counter
		isErr  bool
	}
	tests := []struct {
		name string
		repo *ramRepository
		args args
		want want
	}{
		{
			name: "exists",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{
					"Dogs": metrics.NewCounter("Dogs", 5),
				},
			},
			args: args{name: "Dogs"},
			want: want{
				metrics.NewCounter("Dogs", 5),
				false,
			},
		},
		{
			name: "not found",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{
					"Cats": metrics.NewCounter("Cats", 3),
				},
			},
			args: args{name: "Dogs"},
			want: want{
				nil,
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.repo
			got, err := rs.GetCounter(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.want.isErr {
				t.Errorf("GetCounter() error = %v, wantErr %v", err, tt.want.isErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.metric) {
				t.Errorf("GetCounter() got = %v, want %v", got, tt.want.metric)
			}
		})
	}
}

func Test_ramRepository_AddCounter(t *testing.T) {
	type args struct {
		ctx   context.Context
		name  string
		value int64
	}

	type want struct {
		metric metrics.Counter
		isErr  bool
	}
	tests := []struct {
		name string
		repo *ramRepository
		args args
		want want
	}{
		{
			name: "good add to exists",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{
					"Dogs": metrics.NewCounter("Dogs", 5),
				},
			},
			args: args{name: "Dogs", value: 2},
			want: want{
				metrics.NewCounter("Dogs", 7),
				false,
			},
		},
		{
			name: "good add to not exists",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{},
			},
			args: args{name: "Dogs", value: 2},
			want: want{
				metrics.NewCounter("Dogs", 2),
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.repo
			got, err := rs.AddCounter(tt.args.ctx, tt.args.name, tt.args.value)
			if (err != nil) != tt.want.isErr {
				t.Errorf("AddCounter() error = %v, wantErr %v", err, tt.want.isErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.metric) {
				t.Errorf("AddCounter() got = %v, want %v", got, tt.want.metric)
			}
		})
	}
}

func Test_ramRepository_GetGauge(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}

	type want struct {
		metric metrics.Gauge
		isErr  bool
	}
	tests := []struct {
		name string
		repo *ramRepository
		args args
		want want
	}{
		{
			name: "exists",
			repo: &ramRepository{
				gaugeStorage: map[string]metrics.Gauge{
					"WaterPercent": metrics.NewGauge("WaterPercent", 0.35),
				},
			},
			args: args{name: "WaterPercent"},
			want: want{
				metrics.NewGauge("WaterPercent", 0.35),
				false,
			},
		},
		{
			name: "not found",
			repo: &ramRepository{
				gaugeStorage: map[string]metrics.Gauge{
					"FoodPercent": metrics.NewGauge("FoodPercent", 0.8),
				},
			},
			args: args{name: "WaterPercent"},
			want: want{
				nil,
				true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.repo
			got, err := rs.GetGauge(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.want.isErr {
				t.Errorf("GetGauge() error = %v, wantErr %v", err, tt.want.isErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.metric) {
				t.Errorf("GetGauge() got = %v, want %v", got, tt.want.metric)
			}
		})
	}
}

func Test_ramRepository_SetGauge(t *testing.T) {
	type args struct {
		ctx   context.Context
		name  string
		value float64
	}

	type want struct {
		metric metrics.Gauge
		isErr  bool
	}
	tests := []struct {
		name string
		repo *ramRepository
		args args
		want want
	}{
		{
			name: "good set for exists",
			repo: &ramRepository{
				gaugeStorage: map[string]metrics.Gauge{
					"WaterPercent": metrics.NewGauge("WaterPercent", 0.35),
				},
			},
			args: args{name: "WaterPercent", value: 0.5},
			want: want{
				metrics.NewGauge("WaterPercent", 0.5),
				false,
			},
		},
		{
			name: "good set for not exists",
			repo: &ramRepository{
				gaugeStorage: map[string]metrics.Gauge{},
			},
			args: args{name: "FoodPercent", value: 2},
			want: want{
				metrics.NewGauge("FoodPercent", 2),
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.repo
			got, err := rs.SetGauge(tt.args.ctx, tt.args.name, tt.args.value)
			if (err != nil) != tt.want.isErr {
				t.Errorf("SetGauge() error = %v, wantErr %v", err, tt.want.isErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.metric) {
				t.Errorf("SetGauge() got = %v, want %v", got, tt.want.metric)
			}
		})
	}
}

func Test_ramRepository_GetAll(t *testing.T) {
	type want struct {
		metricsSlice []metrics.Metric
		checkErr     bool
	}
	tests := []struct {
		name string
		repo *ramRepository
		want want
	}{
		{
			name: "good",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{
					"A": metrics.NewCounter("A", 1),
					"B": metrics.NewCounter("B", 2),
				},
				gaugeStorage: map[string]metrics.Gauge{
					"C": metrics.NewGauge("C", 0.0001),
					"D": metrics.NewGauge("D", 2.34),
				},
			},
			want: want{
				metricsSlice: []metrics.Metric{
					metrics.NewCounter("A", 1),
					metrics.NewCounter("B", 2),
					metrics.NewGauge("C", 0.0001),
					metrics.NewGauge("D", 2.34),
				},
				checkErr: false,
			},
		},
		{
			name: "empty",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{},
				gaugeStorage:   map[string]metrics.Gauge{},
			},
			want: want{
				metricsSlice: []metrics.Metric{},
				checkErr:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.repo
			got, err := rs.GetAll(context.TODO())
			if (err != nil) != tt.want.checkErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.want.checkErr)
				return
			}
			assert.ElementsMatchf(t, got, tt.want.metricsSlice,
				"GetAll() got = %v, want %v", got, tt.want.metricsSlice)
		})
	}
}

func Test_ramRepository_AddCounters(t *testing.T) {
	type want struct {
		counters []metrics.Counter
		isErr    bool
	}
	tests := []struct {
		name     string
		repo     *ramRepository
		argSlice []metrics.Counter
		want     want
	}{
		{
			name: "good add to exists",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{
					"Dogs": metrics.NewCounter("Dogs", 5),
				},
			},
			argSlice: []metrics.Counter{metrics.NewCounter("Dogs", 2)},
			want: want{
				[]metrics.Counter{metrics.NewCounter("Dogs", 7)},
				false,
			},
		},
		{
			name: "good add to not exists",
			repo: &ramRepository{
				counterStorage: map[string]metrics.Counter{
					"Dogs": metrics.NewCounter("Dogs", 2),
				},
			},
			argSlice: []metrics.Counter{
				metrics.NewCounter("Dogs", 2),
				metrics.NewCounter("Cats", 3),
			},
			want: want{
				[]metrics.Counter{
					metrics.NewCounter("Dogs", 4),
					metrics.NewCounter("Cats", 3),
				},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.repo
			got, err := rs.AddCounters(context.TODO(), tt.argSlice)
			if (err != nil) != tt.want.isErr {
				t.Errorf("AddCounters() error = %v, wantErr %v", err, tt.want.isErr)
				return
			}
			assert.ElementsMatchf(t, got, tt.want.counters,
				"AddCounters() got = %v, want %v", got, tt.want.counters)
		})
	}
}

func Test_ramRepository_SetGauges(t *testing.T) {
	type want struct {
		gauges []metrics.Gauge
		isErr  bool
	}
	tests := []struct {
		name     string
		repo     *ramRepository
		argSlice []metrics.Gauge
		want     want
	}{
		{
			name: "good set for exists",
			repo: &ramRepository{
				gaugeStorage: map[string]metrics.Gauge{
					"WaterPercent": metrics.NewGauge("WaterPercent", 0.35),
				},
			},
			argSlice: []metrics.Gauge{metrics.NewGauge("WaterPercent", 1)},
			want: want{
				[]metrics.Gauge{metrics.NewGauge("WaterPercent", 1)},
				false,
			},
		},
		{
			name: "good set for not exists",
			repo: &ramRepository{
				gaugeStorage: map[string]metrics.Gauge{
					"WaterPercent": metrics.NewGauge("WaterPercent", 0.62),
				},
			},
			argSlice: []metrics.Gauge{
				metrics.NewGauge("WaterPercent", 0.2),
				metrics.NewGauge("FoodPercent", 1),
			},
			want: want{
				[]metrics.Gauge{
					metrics.NewGauge("WaterPercent", 0.2),
					metrics.NewGauge("FoodPercent", 1),
				},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := tt.repo
			got, err := rs.SetGauges(context.TODO(), tt.argSlice)
			if (err != nil) != tt.want.isErr {
				t.Errorf("SetGauges() error = %v, wantErr %v", err, tt.want.isErr)
				return
			}
			assert.ElementsMatchf(t, got, tt.want.gauges,
				"SetGauges() got = %v, want %v", got, tt.want.gauges)
		})
	}
}
