package metrics

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

type Metric interface {
	GetName() string
	GetValue() string
	GetType() string
	Hash(key []byte) string
	ToParams() Params
}

type Gauge interface {
	Metric
	Set(value float64)
	Value() float64
}
type gauge struct {
	name  string
	value *float64
}

func (g *gauge) String() string {
	return fmt.Sprintf("gauge %v: %v", g.GetName(), g.GetValue())
}

func (g *gauge) GetName() string {
	return g.name
}

func (g *gauge) GetValue() string {
	return fmt.Sprintf("%v", g.Value())
}

func (g *gauge) GetType() string {
	return GaugeType
}

func (g *gauge) ToParams() Params {
	v := g.Value()
	return Params{Name: g.name, Type: g.GetType(), ValueGauge: &v}
}

func (g *gauge) Hash(key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(fmt.Sprintf("%s:gauge:%f", g.name, g.Value())))
	return hex.EncodeToString(h.Sum(nil))
}

func (g *gauge) Set(value float64) {
	*g.value = value
}

func (g *gauge) Value() float64 {
	return *g.value
}

func NewGauge(name string) *gauge {
	var v float64
	return &gauge{name: name, value: &v}
}

type Counter interface {
	Metric
	Inc()
	Add(value int64)
	Value() int64
}

type counter struct {
	name  string
	value *int64
}

func (c *counter) String() string {
	return fmt.Sprintf("counter %v: %v", c.GetName(), c.GetValue())
}

func (c *counter) Inc() {
	*c.value++
}

func (c *counter) Add(value int64) {
	*c.value += value
}

func (c *counter) GetName() string {
	return c.name
}

func (c *counter) GetValue() string {
	return fmt.Sprintf("%d", c.Value())
}

func (c *counter) GetType() string {
	return CounterType
}

func (c *counter) ToParams() Params {
	v := c.Value()
	return Params{Name: c.name, Type: c.GetType(), ValueCounter: &v}
}

func (c *counter) Hash(key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(fmt.Sprintf("%s:counter:%d", c.name, c.Value())))
	return hex.EncodeToString(h.Sum(nil))
}

func (c *counter) Value() int64 {
	return *c.value
}

func NewCounter(name string) *counter {
	var v int64
	return &counter{name: name, value: &v}
}

func NewCounterFromParams(params Params) (*counter, error) {
	if err := CheckName(params.Name); err != nil {
		return nil, err
	}
	return &counter{name: params.Name, value: params.ValueCounter}, nil
}

func NewGaugeFromParams(params Params) (*gauge, error) {
	if err := CheckName(params.Name); err != nil {
		return nil, err
	}
	return &gauge{name: params.Name, value: params.ValueGauge}, nil
}

func NewMetricFromParams(params Params) (Metric, error) {
	var (
		metric Metric
		err    error
	)
	if err := CheckType(params.Type); err != nil {
		return nil, err
	}
	switch params.Type {
	case CounterType:
		metric, err = NewCounterFromParams(params)
	case GaugeType:
		metric, err = NewGaugeFromParams(params)
	}
	return metric, err
}
