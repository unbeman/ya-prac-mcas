package main

import (
	"context"
	"fmt"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/utils"
	"io"
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
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func SendPostRequest(ctx context.Context, client http.Client, url string, body io.Reader) { // TODO: write http connector
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	request.Header.Set("Content-Type", "text/plain")
	if err != nil {
		log.Fatalln(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err) //TODO: retry request?
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Received status code: %v for post request to %v", response.StatusCode, url)
}

func Report(ctx context.Context, client http.Client, metrics map[string]metrics.Metric) {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(len(metrics))
	for _, metric := range metrics {
		url := utils.FormatURL(reportAddr, metric.GetType(), metric.GetName(), metric.GetValue())
		go func() {
			SendPostRequest(ctx2, client, url, nil)
			wg.Done()
		}()
	}
	wg.Wait()
}

func DoWork(ctx context.Context) { // TODO: rename
	log.Println("Agent started")
	reportTicker := time.NewTicker(reportInterval)
	pollTicker := time.NewTicker(pollInterval)
	am := metrics.NewAgentMetrics()
	client := http.Client{}
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped by context")
			return
		case <-reportTicker.C:
			Report(ctx, client, am.GetMetrics())
		case <-pollTicker.C:
			metrics.UpdateMetrics(am)
		}
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	defer func() {
		cancel()
		log.Println("Agent cancelled")
	}()
	DoWork(ctx)
	<-ctx.Done()
}
