package agent

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/internal/agent/sender"
	"github.com/unbeman/ya-prac-mcas/internal/utils"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type agentMetrics struct {
	reporter       sender.Sender
	collection     *MetricsCollection
	tickerPool     *utils.TickerPool
	hashKey        []byte
	pollInterval   time.Duration
	reportInterval time.Duration
	reportTimeout  time.Duration
}

func NewAgentMetrics(cfg *configs.AgentConfig) *agentMetrics {
	reporter := sender.NewHTTPSender(cfg.Connection)
	collector := NewMetricsCollection()
	tickerPool := utils.NewTickerPool()
	return &agentMetrics{
		reporter:       reporter,
		collection:     collector,
		tickerPool:     tickerPool,
		hashKey:        []byte(cfg.Key),
		pollInterval:   cfg.PollInterval,
		reportInterval: cfg.ReportInterval,
	}
}

func (am *agentMetrics) Report(ctx context.Context) {
	paramSlice := am.prepareMetrics(am.collection.GetMetrics(ctx))
	am.reporter.SendJSONMetrics(ctx, paramSlice)
}

func (am *agentMetrics) Run(ctx context.Context) {
	log.Infoln("Agent started")

	am.tickerPool.AddTask(ctx, "UpdateRuntimeMetrics", am.collection.UpdateRuntimeMetrics, am.pollInterval)
	am.tickerPool.AddTask(ctx, "UpdateMemCPUMetrics", am.collection.UpdateMemCPUMetrics, am.pollInterval)
	am.tickerPool.AddTask(ctx, "Report", am.Report, am.reportInterval)

	am.tickerPool.Wait()
}

func (am agentMetrics) getHash(metric metrics.Metric) string {
	if len(am.hashKey) == 0 {
		return ""
	}
	return metric.Hash(am.hashKey)
}

func (am agentMetrics) prepareMetrics(ms []metrics.Metric) []metrics.Params {
	paramSlice := make([]metrics.Params, 0, len(ms))
	for _, metric := range ms {
		params := metric.ToParams()
		params.Hash = am.getHash(metric)
		paramSlice = append(paramSlice, params)
	}
	return paramSlice
}
