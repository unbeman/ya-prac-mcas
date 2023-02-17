package agent

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type Sender interface {
	SendMetric(ctx context.Context, m metrics.Metric)
}

type agentMetrics struct {
	address        string
	client         http.Client
	collection     *MetricsCollection
	pollInterval   time.Duration
	reportInterval time.Duration
}

func NewAgentMetrics(addr string, pI, rI time.Duration) *agentMetrics {
	return &agentMetrics{address: addr,
		client:         http.Client{},
		collection:     NewMetricsCollection(),
		pollInterval:   pI,
		reportInterval: rI}
}

func (am *agentMetrics) Report(ctx context.Context, ms map[string]metrics.Metric) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	for _, metric := range ms {
		wg.Add(1)
		go func(m metrics.Metric) {
			am.SendMetric(ctx2, m)
			wg.Done()
		}(metric)
	}
	wg.Wait()
}

func (am *agentMetrics) DoWork(ctx context.Context) {
	log.Println("Agent started")
	reportTicker := time.NewTicker(am.reportInterval)
	pollTicker := time.NewTicker(am.pollInterval)
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped by context")
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
		log.Println(err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Received status code: %v for post request to %v", response.StatusCode, url)
}
