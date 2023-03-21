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

func (jm *Params) ParseURI(request *http.Request, requiredKeys ...string) error {
	params := &Params{}
	for _, key := range requiredKeys {
		value := chi.URLParam(request, key)
		//if len(value) == 0 {
		//	return nil, fmt.Errorf("ParseURI: %v is %w", key, ErrInvalidValue)
		//}
		switch key {
		case PType:
			err := CheckType(value)
			if err != nil {
				return fmt.Errorf("ParseURI: %v - %w", key, err)
			}
			params.Type = value
		case PName:
			if err := CheckName(value); err != nil {
				return fmt.Errorf("ParseURI: %v - %w", key, err)
			}
			params.Name = value
		case PValue:
			switch params.Type { // unsafe
			case GaugeType:
				gValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("ParseURI: %v = %v - %w", key, value, ErrInvalidValue)
				}
				params.ValueGauge = &gValue
			case CounterType:
				cValue, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("ParseURI: %v = %v - %w", key, value, ErrInvalidValue)
				}
				params.ValueCounter = &cValue
			}
		default:
			return fmt.Errorf("ParseURI: unknown parameter")
		}
	}
	return nil
}

func (jm *Params) ParseJson(data io.Reader) error {
	if err := json.NewDecoder(data).Decode(&jm); err != nil {
		return err
	}
	return nil
}

func CheckType(typeStr string) error {
	switch typeStr {
	case GaugeType, CounterType:
		return nil
	default:
		return fmt.Errorf("CheckType: %v - %w", typeStr, ErrInvalidType)
	}
}

func CheckName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("CheckType: %v - %w", name, ErrInvalidValue)
	}
	return nil
}
