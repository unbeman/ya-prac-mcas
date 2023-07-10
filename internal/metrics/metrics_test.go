package metrics

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func ptrGaugeValue(n float64) *float64 {
	return &n
}

func ptrCounterValue(n int64) *int64 {
	return &n
}

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
				ValueGauge: ptrGaugeValue(0.35),
			},
			want: NewGauge("WaterPercent", 0.35),
		},
		{
			name: "good Counter",
			input: Params{
				Name:         "Dog",
				Type:         "counter",
				ValueCounter: ptrCounterValue(10),
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

func Test_gauge_String(t *testing.T) {
	type fields struct {
		name  string
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "WaterPercent",
				value: ptrGaugeValue(0.8),
			},
			want: "gauge WaterPercent: 0.8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := g.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_String(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "CatCounter",
				value: ptrCounterValue(10),
			},
			want: "counter CatCounter: 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gauge_GetName(t *testing.T) {
	type fields struct {
		name  string
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "WaterPercent",
				value: ptrGaugeValue(0.8),
			},
			want: "WaterPercent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := g.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_GetName(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "DogCounter",
				value: ptrCounterValue(10),
			},
			want: "DogCounter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := c.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gauge_GetValue(t *testing.T) {
	type fields struct {
		name  string
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "WaterPercent",
				value: ptrGaugeValue(0.8),
			},
			want: "0.8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := g.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_GetValue(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "DogCounter",
				value: ptrCounterValue(10),
			},
			want: "10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := c.GetValue(); got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gauge_GetType(t *testing.T) {
	type fields struct {
		name  string
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "WaterPercent",
				value: ptrGaugeValue(0.8),
			},
			want: "gauge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := g.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_GetType(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "good",
			fields: fields{
				name:  "DogCounter",
				value: ptrCounterValue(10),
			},
			want: "counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := c.GetType(); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gauge_ToParams(t *testing.T) {
	type fields struct {
		name  string
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		want   Params
	}{
		{
			name: "good",
			fields: fields{
				name:  "WaterPercent",
				value: ptrGaugeValue(0.8),
			},
			want: Params{
				Name:         "WaterPercent",
				Type:         "gauge",
				ValueCounter: nil,
				ValueGauge:   ptrGaugeValue(0.8),
				Hash:         "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := g.ToParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_counter_ToParams(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}
	tests := []struct {
		name   string
		fields fields
		want   Params
	}{
		{
			name: "good",
			fields: fields{
				name:  "DogCounter",
				value: ptrCounterValue(10),
			},
			want: Params{
				Name:         "DogCounter",
				Type:         "counter",
				ValueCounter: ptrCounterValue(10),
				ValueGauge:   nil,
				Hash:         "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := c.ToParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gauge_Hash(t *testing.T) {
	type fields struct {
		name  string
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		key    []byte
	}{
		{
			name: "good",
			fields: fields{
				name:  "WaterPercent",
				value: ptrGaugeValue(0.8),
			},
			key: []byte("key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := g.Hash(tt.key); len(got) == 0 {
				t.Error("Hash() no hash generated")
			}
		})
	}
}

func Test_counter_Hash(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}
	tests := []struct {
		name   string
		fields fields
		key    []byte
	}{
		{
			name: "good",
			fields: fields{
				name:  "DogCount",
				value: ptrCounterValue(10),
			},
			key: []byte("key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			if got := g.Hash(tt.key); len(got) == 0 {
				t.Error("Hash() no hash generated")
			}
		})
	}
}

func Test_gauge_Set(t *testing.T) {
	type fields struct {
		name  string
		value *float64
	}
	tests := []struct {
		name   string
		fields fields
		value  float64
	}{
		{
			name: "good",
			fields: fields{
				name:  "WaterPercent",
				value: ptrGaugeValue(0.8),
			},
			value: 0.0001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			g.Set(tt.value)
			assert.Equalf(t, tt.value, *g.value, "Set(%v) value not set", tt.value)
		})
	}
}

func Test_counter_Inc(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "good",
			fields: fields{name: "CatCounter", value: ptrCounterValue(10)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			before := *c.value
			c.Inc()
			after := *c.value
			assert.Equal(t, before+1, after, "Inc() value not incremented")
		})
	}
}

func Test_counter_Add(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}

	tests := []struct {
		name   string
		fields fields
		value  int64
	}{
		{
			name:   "good",
			fields: fields{name: "CatCounter", value: ptrCounterValue(10)},
			value:  5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			before := *c.value
			c.Add(tt.value)
			after := *c.value
			assert.Equal(t, before+tt.value, after, "Add() value is not changed")
		})
	}
}

func Test_counter_Set(t *testing.T) {
	type fields struct {
		name  string
		value *int64
	}

	tests := []struct {
		name   string
		fields fields
		value  int64
	}{
		{
			name:   "good",
			fields: fields{name: "CatCounter", value: ptrCounterValue(10)},
			value:  5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				name:  tt.fields.name,
				value: tt.fields.value,
			}
			c.Set(tt.value)
			assert.Equal(t, tt.value, *c.value, "Set() value is not changed")
		})
	}
}
