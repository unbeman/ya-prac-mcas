package metrics

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestParamsSlice_ParseJSON(t *testing.T) {
	type want struct {
		slice      ParamsSlice
		checkSlice bool
		checkErr   bool
	}

	tests := []struct {
		name        string
		inputReader io.Reader
		want        want
	}{
		{
			name: "good",
			inputReader: func() io.Reader {
				gaugeA := 0.0001
				gaugeB := 0.34
				var counterC int64 = 1
				var counterD int64 = 123
				slice := ParamsSlice{
					Params{
						Name:       "GaugeA",
						Type:       GaugeType,
						ValueGauge: &gaugeA,
					},
					Params{
						Name:       "GaugeB",
						Type:       GaugeType,
						ValueGauge: &gaugeB,
					},
					Params{
						Name:         "CounterC",
						Type:         CounterType,
						ValueCounter: &counterC,
					},
					Params{
						Name:         "CounterD",
						Type:         CounterType,
						ValueCounter: &counterD,
					},
				}
				b, _ := json.Marshal(slice)
				return bytes.NewReader(b)
			}(),
			want: want{slice: func() ParamsSlice {
				gaugeA := 0.0001
				gaugeB := 0.34
				var counterC int64 = 1
				var counterD int64 = 123
				slice := ParamsSlice{
					Params{
						Name:       "GaugeA",
						Type:       GaugeType,
						ValueGauge: &gaugeA,
					},
					Params{
						Name:       "GaugeB",
						Type:       GaugeType,
						ValueGauge: &gaugeB,
					},
					Params{
						Name:         "CounterC",
						Type:         CounterType,
						ValueCounter: &counterC,
					},
					Params{
						Name:         "CounterD",
						Type:         CounterType,
						ValueCounter: &counterD,
					},
				}
				return slice
			}(),
			},
		},
		{
			name: "empty json list",
			inputReader: func() io.Reader {
				js := "[]"
				return strings.NewReader(js)
			}(),
			want: want{
				checkSlice: true,
				slice:      ParamsSlice{},
				checkErr:   false,
			},
		},
		{
			name: "invalid json",
			inputReader: func() io.Reader {
				js := `[{"id": "Dog", "type": "counter", ""]`
				return strings.NewReader(js)
			}(),
			want: want{
				checkSlice: false,
				slice:      nil,
				checkErr:   true,
			},
		},
		{
			name: "invalid type",
			inputReader: func() io.Reader {
				slice := ParamsSlice{Params{Name: "Good", Type: "invalid", ValueCounter: nil}}
				js, _ := json.Marshal(slice)
				return bytes.NewReader(js)
			}(),
			want: want{
				checkSlice: false,
				slice:      nil,
				checkErr:   true,
			},
		},
		{
			name: "invalid name",
			inputReader: func() io.Reader {
				slice := ParamsSlice{Params{Name: "", Type: CounterType, ValueCounter: nil}}
				js, _ := json.Marshal(slice)
				return bytes.NewReader(js)
			}(),
			want: want{
				checkSlice: false,
				slice:      nil,
				checkErr:   true,
			},
		},
		{
			name: "invalid value",
			inputReader: func() io.Reader {
				raw := `[{"id": "Dog", "type": "counter"},
{"id": "WaterPercent", "type": "gauge"}]`
				return strings.NewReader(raw)
			}(),
			want: want{
				checkSlice: false,
				slice:      nil,
				checkErr:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ps ParamsSlice

			err := ps.ParseJSON(tt.inputReader)
			if (err != nil) != tt.want.checkErr {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.want.checkErr)
			}

			if tt.want.checkSlice && !reflect.DeepEqual(ps, tt.want.slice) {
				t.Errorf("ParseJSON() got = %v, want %v", ps, tt.want.slice)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	type want struct {
		params      Params
		checkMetric bool
		checkErr    bool
	}

	tests := []struct {
		name        string
		inputReader io.Reader
		inputArgs   []string
		want        want
	}{
		{
			name:      "good",
			inputArgs: []string{PType, PName, PValue},
			inputReader: func() io.Reader {
				gaugeA := 0.0001
				p := Params{
					Name:       "GaugeA",
					Type:       GaugeType,
					ValueGauge: &gaugeA,
				}
				b, _ := json.Marshal(p)
				return bytes.NewReader(b)
			}(),
			want: want{params: func() Params {
				gaugeA := 0.0001
				p := Params{
					Name:       "GaugeA",
					Type:       GaugeType,
					ValueGauge: &gaugeA,
				}
				return p
			}(),
			},
		},
		{
			name:      "empty json",
			inputArgs: []string{},
			inputReader: func() io.Reader {
				js := "{}"
				return strings.NewReader(js)
			}(),
			want: want{
				checkMetric: true,
				checkErr:    false,
			},
		},
		{
			name:      "invalid json",
			inputArgs: []string{PType, PName, PValue},
			inputReader: func() io.Reader {
				js := `{"id": "Dog", "type": "counter", "value": 3,`
				return strings.NewReader(js)
			}(),
			want: want{
				checkMetric: false,
				checkErr:    true,
			},
		},
		{
			name:      "invalid type",
			inputArgs: []string{PType, PName, PValue},
			inputReader: func() io.Reader {
				p := Params{Name: "Good", Type: "invalid", ValueCounter: nil}
				js, _ := json.Marshal(p)
				return bytes.NewReader(js)
			}(),
			want: want{
				checkMetric: false,
				checkErr:    true,
			},
		},
		{
			name:      "invalid name",
			inputArgs: []string{PType, PName, PValue},
			inputReader: func() io.Reader {
				p := Params{Name: "", Type: CounterType, ValueCounter: nil}
				js, _ := json.Marshal(p)
				return bytes.NewReader(js)
			}(),
			want: want{
				checkMetric: false,
				checkErr:    true,
			},
		},
		{
			name:      "invalid value",
			inputArgs: []string{PType, PName, PValue},
			inputReader: func() io.Reader {
				raw := `{"id": "Dog", "type": "counter"}`
				return strings.NewReader(raw)
			}(),
			want: want{
				checkMetric: false,
				checkErr:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			params, err := ParseJSON(tt.inputReader, tt.inputArgs...)
			if (err != nil) != tt.want.checkErr {
				t.Errorf("ParseJSON() error = %v, wantErr %v", err, tt.want.checkErr)
			}

			if tt.want.checkMetric && !reflect.DeepEqual(params, tt.want.params) {
				t.Errorf("ParseJSON() got = %v, want %v", params, tt.want.params)
			}
		})
	}
}
