package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		return err
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

func (jm *Params) String() string {
	var delta, value = "nil", "nil"
	if jm.ValueCounter != nil {
		delta = fmt.Sprintf("%v", *jm.ValueCounter)
	}
	if jm.ValueGauge != nil {
		value = fmt.Sprintf("%v", *jm.ValueGauge)
	}
	return fmt.Sprintf("{ID:%v; MType:%v; Delta:%v; Value:%v};", jm.Name, jm.Type, delta, value)
}

func ParseURI(request *http.Request, requiredKeys ...string) (Params, error) {
	params := Params{}
	for _, key := range requiredKeys {
		value := chi.URLParam(request, key)
		//if len(value) == 0 {
		//	return nil, fmt.Errorf("ParseURI: %v is %w", key, ErrInvalidValue)
		//}
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
			return params, fmt.Errorf("ParseURI: unknown parameter")
		}
	}
	return params, nil
}

func ParseJSON(data io.Reader, requiredKeys ...string) (Params, error) {
	var params Params
	if err := json.NewDecoder(data).Decode(&params); err != nil {
		return params, err
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
		return fmt.Errorf("CheckType: (%v) - %w", typeStr, ErrInvalidType)
	}
}

func CheckName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("CheckName: (%v) - %w", name, ErrInvalidValue)
	}
	return nil
}

func CheckValues(valueGauge *float64, valueCounter *int64) error {
	if valueGauge == nil && valueCounter == nil {
		return fmt.Errorf("CheckValues: %w", ErrInvalidValue)
	}
	return nil
}
