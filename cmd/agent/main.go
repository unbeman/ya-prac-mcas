package main

import (
	"context"
	"github.com/unbeman/ya-prac-mcas/internal/agent"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	reportAddr     = "127.0.0.1:8080"
	pollInterval   = 1 * time.Second
	reportInterval = 2 * time.Second
)

func Report(ctx context.Context, cm agent.ClientMetric, ms map[string]metrics.Metric) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(len(ms))
	for _, metric := range ms {
		go func(m metrics.Metric) {
			cm.SendMetric(ctx2, m)
			wg.Done()
		}(metric)
	}
	wg.Wait()
}

func DoWork(ctx context.Context, clientMetic agent.ClientMetric) {
	log.Println("Agent started")
	reportTicker := time.NewTicker(reportInterval)
	pollTicker := time.NewTicker(pollInterval)
	am := agent.NewAgentMetrics()
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped by context")
			return
		case <-reportTicker.C:
			Report(ctx, clientMetic, am.GetMetrics())
			am.PollCount = metrics.NewCounter("PollCount")
		case <-pollTicker.C:
			agent.UpdateMetrics(am)
		}
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	defer func() {
		cancel()
		log.Println("Agent cancelled")
	}()
	cm := agent.NewClientMetric(reportAddr, http.Client{})

	DoWork(ctx, cm)
	<-ctx.Done()
}
