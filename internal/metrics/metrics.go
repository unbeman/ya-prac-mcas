package metrics

import (
	"fmt"
)

type Metric interface {
	GetName() string
	GetValue() string
	GetType() string
}

type Gauge interface {
	Metric
	Set(value float64)
}
type gauge struct {
	name  string
	value float64
}

func (g *gauge) GetName() string {
	return g.name
}

func (g *gauge) GetValue() string {
	return fmt.Sprintf("%f", g.value)
}

func (g *gauge) GetType() string {
	return "gauge"
}

func (g *gauge) Set(value float64) {
	g.value = value
}

func NewGauge(name string) *gauge {
	return &gauge{name: name}
}

func NewGaugeWithValue(name string, value float64) *gauge {
	return &gauge{name: name, value: value}
}

type Counter interface {
	Metric
	Inc()
	//Add()
}

type counter struct {
	name  string
	value int64
}

func (c *counter) Inc() {
	c.value++
}

func (c *counter) GetName() string {
	return c.name
}

func (c *counter) GetValue() string {
	return fmt.Sprintf("%d", c.value)
}

func (c *counter) GetType() string {
	return "counter"
}

func NewCounter(name string) *counter {
	return &counter{name: name}
}

func NewCounterWithValue(name string, value int64) *counter {
	return &counter{name: name, value: value}
}
