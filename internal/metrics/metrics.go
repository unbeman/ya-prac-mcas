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

func NewGauge(name string, value float64) *gauge {
	return &gauge{name: name, value: &value}
}

type Counter interface {
	Metric
	Inc()
	Add(value int64)
	Value() int64
	Set(value int64)
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

func (c *counter) Set(value int64) {
	*c.value = value
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

func NewCounter(name string, value int64) *counter {
	return &counter{name: name, value: &value}
}

func NewCounterFromParams(params Params) *counter {
	return &counter{name: params.Name, value: params.ValueCounter}
}

func NewGaugeFromParams(params Params) *gauge {
	return &gauge{name: params.Name, value: params.ValueGauge}
}

func NewMetricFromParams(params Params) Metric {
	var metric Metric

	switch params.Type {
	case CounterType:
		metric = NewCounterFromParams(params)
	case GaugeType:
		metric = NewGaugeFromParams(params)
	}
	return metric
}
