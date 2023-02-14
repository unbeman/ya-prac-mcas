package metrics

import (
	"reflect"
	"testing"
)

//TODO: cover methods

func TestNewCounter(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *counter
	}{
		{name: "Positive",
			args: args{name: "Cat"},
			want: &counter{name: "Cat", value: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCounter(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGauge(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *gauge
	}{
		{name: "Positive",
			args: args{name: "Dog"},
			want: &gauge{name: "Dog", value: 0.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGauge(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}
