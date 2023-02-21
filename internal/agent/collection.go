package agent

import (
	"log"
	"math/rand"
	"runtime"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type MetricsCollection struct {
	Alloc         metrics.Gauge
	BuckHashSys   metrics.Gauge
	Frees         metrics.Gauge
	GCCPUFraction metrics.Gauge
	GCSys         metrics.Gauge
	HeapAlloc     metrics.Gauge
	HeapIdle      metrics.Gauge
	HeapInuse     metrics.Gauge
	HeapObjects   metrics.Gauge
	HeapReleased  metrics.Gauge
	HeapSys       metrics.Gauge
	LastGC        metrics.Gauge
	Lookups       metrics.Gauge
	MCacheInuse   metrics.Gauge
	MCacheSys     metrics.Gauge
	MSpanInuse    metrics.Gauge
	MSpanSys      metrics.Gauge
	Mallocs       metrics.Gauge
	NextGC        metrics.Gauge
	NumForcedGC   metrics.Gauge
	NumGC         metrics.Gauge
	OtherSys      metrics.Gauge
	PauseTotalNs  metrics.Gauge
	StackInuse    metrics.Gauge
	StackSys      metrics.Gauge
	Sys           metrics.Gauge
	TotalAlloc    metrics.Gauge
	RandomValue   metrics.Gauge
	PollCount     metrics.Counter
}

func NewMetricsCollection() *MetricsCollection {
	log.Println("Metrics CREATED")
	return &MetricsCollection{
		Alloc:         metrics.NewGauge("Alloc"),
		BuckHashSys:   metrics.NewGauge("BuckHashSys"),
		Frees:         metrics.NewGauge("Frees"),
		GCCPUFraction: metrics.NewGauge("GCCPUFraction"),
		GCSys:         metrics.NewGauge("GCSys"),
		HeapAlloc:     metrics.NewGauge("HeapAlloc"),
		HeapIdle:      metrics.NewGauge("HeapIdle"),
		HeapInuse:     metrics.NewGauge("HeapInuse"),
		HeapObjects:   metrics.NewGauge("HeapObjects"),
		HeapReleased:  metrics.NewGauge("HeapReleased"),
		HeapSys:       metrics.NewGauge("HeapSys"),
		LastGC:        metrics.NewGauge("LastGC"),
		Lookups:       metrics.NewGauge("Lookups"),
		MCacheInuse:   metrics.NewGauge("MCacheInuse"),
		MCacheSys:     metrics.NewGauge("MCacheSys"),
		MSpanInuse:    metrics.NewGauge("MSpanInuse"),
		MSpanSys:      metrics.NewGauge("MSpanSys"),
		Mallocs:       metrics.NewGauge("Mallocs"),
		NextGC:        metrics.NewGauge("NextGC"),
		NumForcedGC:   metrics.NewGauge("NumForcedGC"),
		NumGC:         metrics.NewGauge("NumGC"),
		OtherSys:      metrics.NewGauge("OtherSys"),
		PauseTotalNs:  metrics.NewGauge("PauseTotalNs"),
		StackInuse:    metrics.NewGauge("StackInuse"),
		StackSys:      metrics.NewGauge("StackSys"),
		Sys:           metrics.NewGauge("Sys"),
		TotalAlloc:    metrics.NewGauge("TotalAlloc"),
		RandomValue:   metrics.NewGauge("RandomValue"),
		PollCount:     metrics.NewCounter("PollCount"),
	}
}

func UpdateMetrics(am *MetricsCollection) {
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

func (am *MetricsCollection) GetMetrics() map[string]metrics.Metric {
	metricsMap := map[string]metrics.Metric{}
	metricsMap[am.Alloc.GetName()] = am.Alloc
	metricsMap[am.BuckHashSys.GetName()] = am.BuckHashSys
	metricsMap[am.Frees.GetName()] = am.Frees
	metricsMap[am.GCCPUFraction.GetName()] = am.GCCPUFraction
	metricsMap[am.GCSys.GetName()] = am.GCSys
	metricsMap[am.HeapAlloc.GetName()] = am.HeapAlloc
	metricsMap[am.HeapIdle.GetName()] = am.HeapIdle
	metricsMap[am.HeapInuse.GetName()] = am.HeapInuse
	metricsMap[am.HeapObjects.GetName()] = am.HeapObjects
	metricsMap[am.HeapReleased.GetName()] = am.HeapReleased
	metricsMap[am.HeapSys.GetName()] = am.HeapSys
	metricsMap[am.LastGC.GetName()] = am.LastGC
	metricsMap[am.Lookups.GetName()] = am.Lookups
	metricsMap[am.MCacheInuse.GetName()] = am.MCacheInuse
	metricsMap[am.MCacheSys.GetName()] = am.MCacheSys
	metricsMap[am.MSpanInuse.GetName()] = am.MSpanInuse
	metricsMap[am.MSpanSys.GetName()] = am.MSpanSys
	metricsMap[am.Mallocs.GetName()] = am.Mallocs
	metricsMap[am.NextGC.GetName()] = am.NextGC
	metricsMap[am.NumForcedGC.GetName()] = am.NumForcedGC
	metricsMap[am.NumGC.GetName()] = am.NumGC
	metricsMap[am.OtherSys.GetName()] = am.OtherSys
	metricsMap[am.PauseTotalNs.GetName()] = am.PauseTotalNs
	metricsMap[am.StackInuse.GetName()] = am.StackInuse
	metricsMap[am.StackSys.GetName()] = am.StackSys
	metricsMap[am.Sys.GetName()] = am.Sys
	metricsMap[am.TotalAlloc.GetName()] = am.TotalAlloc
	metricsMap[am.RandomValue.GetName()] = am.RandomValue
	metricsMap[am.PollCount.GetName()] = am.PollCount
	return metricsMap
}
