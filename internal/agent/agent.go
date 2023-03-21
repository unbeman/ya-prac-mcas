package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type Sender interface {
	SendMetric(ctx context.Context, m metrics.Metric)
}

type agentMetrics struct {
	address        string
	client         http.Client
	collection     *MetricsCollection
	hashKey        []byte
	pollInterval   time.Duration
	reportInterval time.Duration
	reportTimeout  time.Duration
}

func NewAgentMetrics(cfg *configs.AgentConfig) *agentMetrics {
	return &agentMetrics{address: cfg.Address,
		client:         http.Client{Timeout: cfg.Connection.ClientTimeout},
		collection:     NewMetricsCollection(),
		hashKey:        []byte(cfg.Key),
		pollInterval:   cfg.PollInterval,
		reportInterval: cfg.ReportInterval,
		reportTimeout:  cfg.ReportTimeout,
	}
}

func (am *agentMetrics) Report(ctx context.Context, ms map[string]metrics.Metric) {
	ctx2, cancel := context.WithTimeout(ctx, am.reportTimeout)
	defer cancel()
	paramSlice := am.prepareMetrics(ms)
	am.SendJSONMetrics(ctx2, paramSlice)
}

func (am *agentMetrics) DoWork(ctx context.Context) {
	log.Infoln("Agent started")
	reportTicker := time.NewTicker(am.reportInterval)
	pollTicker := time.NewTicker(am.pollInterval)
	for {
		select {
		case <-ctx.Done():
			log.Infoln("Worker stopped by context")
			return
		case <-reportTicker.C:
			am.Report(ctx, am.collection.GetMetrics())
			am.collection.PollCount = metrics.NewCounter("PollCount")
		case <-pollTicker.C:
			UpdateMetrics(am.collection)
		}
	}
}

func (am agentMetrics) SendMetric(ctx context.Context, m metrics.Metric) { //TODO: write http connector
	url := FormatURL(am.address, m)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := am.client.Do(request)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Debugf("Received status code: %v for post request to %v\n", response.StatusCode, url)
}

func (am agentMetrics) SendJSONMetrics(ctx context.Context, slice []metrics.Params) { //TODO: write http connector
	url := fmt.Sprintf("http://%s/updates", am.address) //TODO: wrap

	buf, err := json.Marshal(slice)
	if err != nil {
		log.Fatalf("Json marshal failed, %v\n", err)
		return
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := am.client.Do(request)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		log.Errorln(err)
	}
	log.Debugf("Received status code: %v for post request to %v\n", response.StatusCode, url)
}

func (am agentMetrics) SendJSONMetric(ctx context.Context, m metrics.Metric) { //TODO: write http connector
	url := fmt.Sprintf("http://%s/update", am.address) //TODO: wrap
	params := m.ToParams()
	params.Hash = am.getHash(m)
	buf, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Json marshal failed, %v\n", err)
		return
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := am.client.Do(request)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		log.Errorln(err)
	}
	log.Debugf("Received status code: %v for post request to %v\n", response.StatusCode, url)
}

func (am agentMetrics) getHash(metric metrics.Metric) string {
	if len(am.hashKey) == 0 {
		return ""
	}
	return metric.Hash(am.hashKey)
}

func FormatURL(addr string, m metrics.Metric) string {
	return fmt.Sprintf("http://%v/update/%v/%v/%v", addr, m.GetType(), m.GetName(), m.GetValue())
}

func (am agentMetrics) prepareMetrics(ms map[string]metrics.Metric) []metrics.Params {
	paramSlice := make([]metrics.Params, 0, len(ms))
	for _, metric := range ms {
		params := metric.ToParams()
		params.Hash = am.getHash(metric)
		paramSlice = append(paramSlice, params)
	}
	return paramSlice
}
