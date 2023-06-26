package agent

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-mcas/internal/storage"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type MetricsCollection struct {
	storage storage.Repository
}

func NewMetricsCollection() *MetricsCollection {
	log.Infoln("Metrics CREATED")
	return &MetricsCollection{storage: storage.NewRAMRepository()}
}

func (am *MetricsCollection) UpdateRuntimeMetrics(ctx context.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	gaugeSlice := []metrics.Gauge{
		metrics.NewGauge("Alloc", float64(memStats.Alloc)),
		metrics.NewGauge("BuckHashSys", float64(memStats.BuckHashSys)),
		metrics.NewGauge("Frees", float64(memStats.Frees)),
		metrics.NewGauge("GCCPUFraction", memStats.GCCPUFraction),
		metrics.NewGauge("GCSys", float64(memStats.GCSys)),
		metrics.NewGauge("HeapAlloc", float64(memStats.HeapAlloc)),
		metrics.NewGauge("HeapIdle", float64(memStats.HeapIdle)),
		metrics.NewGauge("HeapInuse", float64(memStats.HeapInuse)),
		metrics.NewGauge("HeapObjects", float64(memStats.HeapObjects)),
		metrics.NewGauge("HeapReleased", float64(memStats.HeapReleased)),
		metrics.NewGauge("HeapSys", float64(memStats.HeapReleased)),
		metrics.NewGauge("LastGC", float64(memStats.LastGC)),
		metrics.NewGauge("Lookups", float64(memStats.Lookups)),
		metrics.NewGauge("MCacheInuse", float64(memStats.MCacheInuse)),
		metrics.NewGauge("MCacheSys", float64(memStats.MCacheSys)),
		metrics.NewGauge("MSpanInuse", float64(memStats.MSpanInuse)),
		metrics.NewGauge("MSpanSys", float64(memStats.MSpanSys)),
		metrics.NewGauge("Mallocs", float64(memStats.Mallocs)),
		metrics.NewGauge("NextGC", float64(memStats.NextGC)),
		metrics.NewGauge("NumForcedGC", float64(memStats.NumForcedGC)),
		metrics.NewGauge("NumGC", float64(memStats.NumGC)),
		metrics.NewGauge("OtherSys", float64(memStats.OtherSys)),
		metrics.NewGauge("PauseTotalNs", float64(memStats.PauseTotalNs)),
		metrics.NewGauge("StackInuse", float64(memStats.StackInuse)),
		metrics.NewGauge("StackSys", float64(memStats.StackSys)),
		metrics.NewGauge("Sys", float64(memStats.Sys)),
		metrics.NewGauge("TotalAlloc", float64(memStats.TotalAlloc)),
		metrics.NewGauge("RandomValue", rand.Float64()),
	}
	err := am.storage.SetGauges(ctx, gaugeSlice)
	if err != nil {
		log.Error(err)
	}
	_, err = am.storage.AddCounter(ctx, "PollCount", 1)
	if err != nil {
		log.Error(err)
	}
	log.Infoln("Metrics updated")
}

func (am *MetricsCollection) UpdateMemCPUMetrics(ctx context.Context) {
	memStats, err := mem.VirtualMemory()
	if err != nil {
		log.Error(err)
	}
	gaugeSlice := []metrics.Gauge{
		metrics.NewGauge("TotalMemory", float64(memStats.Total)),
		metrics.NewGauge("FreeMemory", float64(memStats.Free)),
	}

	cpusUsed, err := cpu.Percent(0, false)
	if err != nil {
		log.Error(err)
	}
	for idx, percent := range cpusUsed {
		gaugeSlice = append(gaugeSlice, metrics.NewGauge(fmt.Sprintf("CPUutilization%d", idx+1), percent))
	}

	err = am.storage.SetGauges(ctx, gaugeSlice)
	if err != nil {
		log.Error(err)
	}
}

func (am *MetricsCollection) GetMetrics(ctx context.Context) []metrics.Metric {
	ms, _ := am.storage.GetAll(ctx)
	return ms
}

func (am *MetricsCollection) ResetPollCount(ctx context.Context) {
	oldPC, _ := am.storage.GetCounter(ctx, "PollCount")
	am.storage.AddCounter(ctx, "PollCount", -1*oldPC.Value())
}
