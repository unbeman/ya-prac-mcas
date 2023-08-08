package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	pb "github.com/unbeman/ya-prac-mcas/proto"
)

const (
	PName  string = "name"
	PType  string = "type"
	PValue string = "value"
)

type Params struct { //TODO: make builder .TypeAndName() .Value()
	Name         string   `json:"id"`
	Type         string   `json:"type"`
	ValueCounter *int64   `json:"delta,omitempty"`
	ValueGauge   *float64 `json:"value,omitempty"`
	Hash         string   `json:"hash,omitempty"`
}

type ParamsSlice []Params

func (ps *ParamsSlice) ParseJSON(reader io.Reader) error {
	if err := json.NewDecoder(reader).Decode(ps); err != nil {
		return fmt.Errorf("%w - %v", ErrParseJSON, err)
	}
	for _, params := range *ps {
		if err := CheckType(params.Type); err != nil {
			return err
		}
		if err := CheckName(params.Name); err != nil {
			return err
		}
		if err := CheckValues(params.ValueGauge, params.ValueCounter); err != nil {
			return ErrInvalidValue
		}
	}
	return nil
}

func (ps *ParamsSlice) ParseProto(mps []*pb.Metric) error {

	for _, m := range mps {
		if err := CheckType(m.Type); err != nil {
			return err
		}
		if err := CheckName(m.Name); err != nil {
			return err
		}
		if err := CheckValues(&m.Value, &m.Delta); err != nil {
			return ErrInvalidValue
		}
		p := Params{
			Name:         m.Name,
			Type:         m.Type,
			ValueGauge:   &m.Value,
			ValueCounter: &m.Delta,
			Hash:         m.Hash,
		}
		*ps = append(*ps, p)
	}
	return nil
}

func (ps *ParamsSlice) ToProto() []*pb.Metric {
	protoMetrics := make([]*pb.Metric, 0, len(*ps))
	for _, p := range *ps {
		pm := pb.Metric{
			Name:  p.Name,
			Type:  p.Type,
			Delta: p.GetCounterValue(),
			Value: p.GetGaugeValue(),
			Hash:  p.Hash,
		}
		protoMetrics = append(protoMetrics, &pm)
	}
	return protoMetrics
}

func (p *Params) GetGaugeValue() float64 {
	if p.ValueGauge != nil {
		return *p.ValueGauge
	}
	return 0.0
}

func (p *Params) GetCounterValue() int64 {
	if p.ValueCounter != nil {
		return *p.ValueCounter
	}
	return 0
}

func (p *Params) String() string {
	var delta, value = "nil", "nil"
	if p.ValueCounter != nil {
		delta = fmt.Sprintf("%v", *p.ValueCounter)
	}
	if p.ValueGauge != nil {
		value = fmt.Sprintf("%v", *p.ValueGauge)
	}
	return fmt.Sprintf("{ID:%v; MType:%v; Delta:%v; Value:%v};", p.Name, p.Type, delta, value)
}

func ParseURI(request *http.Request, requiredKeys ...string) (Params, error) {
	params := Params{}
	for _, key := range requiredKeys {
		value := chi.URLParam(request, key)
		switch key {
		case PType:
			err := CheckType(value)
			if err != nil {
				return params, fmt.Errorf("ParseURI: %v - %w", key, err)
			}
			params.Type = value
		case PName:
			if err := CheckName(value); err != nil {
				return params, fmt.Errorf("ParseURI: %v - %w", key, err)
			}
			params.Name = value
		case PValue:
			switch params.Type { // unsafe
			case GaugeType:
				gValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return params, fmt.Errorf("ParseURI: %v = %v - %w", key, value, ErrInvalidValue)
				}
				params.ValueGauge = &gValue
			case CounterType:
				cValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return params, fmt.Errorf("ParseURI: %v = %v - %w", key, value, ErrInvalidValue)
				}
				params.ValueCounter = &cValue
			}
		default:
			return params, fmt.Errorf("ParseURI: %w, (%v) is unknown parameter", ErrParseURI, key)
		}
	}
	return params, nil
}

func ParseJSON(data io.Reader, requiredKeys ...string) (Params, error) {
	var params Params
	if err := json.NewDecoder(data).Decode(&params); err != nil {
		return params, fmt.Errorf("%w - %v", ErrParseJSON, err)
	}
	for _, key := range requiredKeys {
		switch key {
		case PName:
			if err := CheckName(params.Name); err != nil {
				return params, err
			}
		case PType:
			if err := CheckType(params.Type); err != nil {
				return params, err
			}
		case PValue:
			if err := CheckValues(params.ValueGauge, params.ValueCounter); err != nil {
				return params, ErrInvalidValue
			}
		}
	}

	return params, nil
}

func ParseProto(metric *pb.Metric, requiredKeys ...string) (Params, error) {
	params := Params{
		Name:         metric.Name,
		Type:         metric.Type,
		ValueCounter: &metric.Delta,
		ValueGauge:   &metric.Value,
		Hash:         metric.Hash,
	}
	for _, key := range requiredKeys {
		switch key {
		case PName:
			if err := CheckName(params.Name); err != nil {
				return params, err
			}
		case PType:
			if err := CheckType(params.Type); err != nil {
				return params, err
			}
		case PValue:
			if err := CheckValues(params.ValueGauge, params.ValueCounter); err != nil {
				return params, ErrInvalidValue
			}
		}
	}
	return params, nil
}

func CheckType(typeStr string) error {
	switch typeStr {
	case GaugeType, CounterType:
		return nil
	default:
		return fmt.Errorf("checkType: (%v) - %w", typeStr, ErrInvalidType)
	}
}

func CheckName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("checkName: (%v) - %w", name, ErrInvalidValue)
	}
	return nil
}

func CheckValues(valueGauge *float64, valueCounter *int64) error {
	if valueGauge == nil && valueCounter == nil {
		return fmt.Errorf("checkValues: %w", ErrInvalidValue)
	}
	return nil
}
