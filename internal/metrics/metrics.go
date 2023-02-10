package metrics

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
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
	return fmt.Sprintf("%v", g.value)
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
	return fmt.Sprintf("%v", c.value)
}

func (c *counter) GetType() string {
	return "counter"
}

func NewCounter(name string) *counter {
	return &counter{name: name}
}

type AgentMetrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	RandomValue   Gauge
	PollCount     Counter
}

func NewAgentMetrics() *AgentMetrics {
	log.Println("Metrics CREATED")
	return &AgentMetrics{
		Alloc:         NewGauge("Alloc"),
		BuckHashSys:   NewGauge("BuckHashSys"),
		Frees:         NewGauge("Frees"),
		GCCPUFraction: NewGauge("GCCPUFraction"),
		GCSys:         NewGauge("GCSys"),
		HeapAlloc:     NewGauge("HeapAlloc"),
		HeapIdle:      NewGauge("HeapIdle"),
		HeapInuse:     NewGauge("HeapInuse"),
		HeapObjects:   NewGauge("HeapObjects"),
		HeapReleased:  NewGauge("HeapReleased"),
		HeapSys:       NewGauge("HeapSys"),
		LastGC:        NewGauge("LastGC"),
		Lookups:       NewGauge("Lookups"),
		MCacheInuse:   NewGauge("MCacheInuse"),
		MCacheSys:     NewGauge("MCacheSys"),
		MSpanInuse:    NewGauge("MSpanInuse"),
		MSpanSys:      NewGauge("MSpanSys"),
		Mallocs:       NewGauge("Mallocs"),
		NextGC:        NewGauge("NextGC"),
		NumForcedGC:   NewGauge("NumForcedGC"),
		NumGC:         NewGauge("NumGC"),
		OtherSys:      NewGauge("OtherSys"),
		PauseTotalNs:  NewGauge("PauseTotalNs"),
		StackInuse:    NewGauge("StackInuse"),
		StackSys:      NewGauge("StackSys"),
		Sys:           NewGauge("Sys"),
		TotalAlloc:    NewGauge("TotalAlloc"),
		RandomValue:   NewGauge("RandomValue"),
		PollCount:     NewCounter("PollCount"),
	}
}

func UpdateMetrics(am *AgentMetrics) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	am.Alloc.Set(float64(memStats.Alloc))
	am.BuckHashSys.Set(float64(memStats.BuckHashSys))
	am.Frees.Set(float64(memStats.Frees))
	am.GCCPUFraction.Set(memStats.GCCPUFraction)
	am.GCSys.Set(float64(memStats.GCSys))
	am.HeapAlloc.Set(float64(memStats.HeapAlloc))
	am.HeapIdle.Set(float64(memStats.HeapIdle))
	am.HeapInuse.Set(float64(memStats.HeapInuse))
	am.HeapObjects.Set(float64(memStats.HeapObjects))
	am.HeapReleased.Set(float64(memStats.HeapReleased))
	am.LastGC.Set(float64(memStats.LastGC))
	am.Lookups.Set(float64(memStats.Lookups))
	am.MCacheInuse.Set(float64(memStats.MCacheInuse))
	am.MCacheSys.Set(float64(memStats.MCacheSys))
	am.MSpanInuse.Set(float64(memStats.MSpanInuse))
	am.MSpanSys.Set(float64(memStats.MSpanSys))
	am.Mallocs.Set(float64(memStats.Mallocs))
	am.NextGC.Set(float64(memStats.NextGC))
	am.NumForcedGC.Set(float64(memStats.NumForcedGC))
	am.NumGC.Set(float64(memStats.NumGC))
	am.OtherSys.Set(float64(memStats.OtherSys))
	am.PauseTotalNs.Set(float64(memStats.PauseTotalNs))
	am.StackInuse.Set(float64(memStats.StackInuse))
	am.StackSys.Set(float64(memStats.StackSys))
	am.Sys.Set(float64(memStats.Sys))
	am.TotalAlloc.Set(float64(memStats.TotalAlloc))

	am.PollCount.Inc()
	am.RandomValue.Set(rand.Float64())
	log.Println("Metrics updated")
}

func (am *AgentMetrics) GetMetrics() map[string]Metric {
	metrics := map[string]Metric{}
	metrics[am.Alloc.GetName()] = am.Alloc
	metrics[am.BuckHashSys.GetName()] = am.BuckHashSys
	metrics[am.Frees.GetName()] = am.Frees
	metrics[am.GCCPUFraction.GetName()] = am.GCCPUFraction
	metrics[am.GCSys.GetName()] = am.GCSys
	metrics[am.HeapAlloc.GetName()] = am.HeapAlloc
	metrics[am.HeapIdle.GetName()] = am.HeapIdle
	metrics[am.HeapInuse.GetName()] = am.HeapInuse
	metrics[am.HeapObjects.GetName()] = am.HeapObjects
	metrics[am.HeapReleased.GetName()] = am.HeapReleased
	metrics[am.HeapSys.GetName()] = am.HeapSys
	metrics[am.LastGC.GetName()] = am.LastGC
	metrics[am.Lookups.GetName()] = am.Lookups
	metrics[am.MCacheInuse.GetName()] = am.MCacheInuse
	metrics[am.MCacheSys.GetName()] = am.MCacheSys
	metrics[am.MSpanInuse.GetName()] = am.MSpanInuse
	metrics[am.MSpanSys.GetName()] = am.MSpanSys
	metrics[am.Mallocs.GetName()] = am.Mallocs
	metrics[am.NextGC.GetName()] = am.NextGC
	metrics[am.NumForcedGC.GetName()] = am.NumForcedGC
	metrics[am.NumGC.GetName()] = am.NumGC
	metrics[am.OtherSys.GetName()] = am.OtherSys
	metrics[am.PauseTotalNs.GetName()] = am.PauseTotalNs
	metrics[am.StackInuse.GetName()] = am.StackInuse
	metrics[am.StackSys.GetName()] = am.StackSys
	metrics[am.Sys.GetName()] = am.Sys
	metrics[am.TotalAlloc.GetName()] = am.TotalAlloc
	metrics[am.RandomValue.GetName()] = am.RandomValue
	metrics[am.PollCount.GetName()] = am.PollCount
	return metrics
}
