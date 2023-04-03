package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-mcas/configs"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type Sender interface {
	SendMetric(ctx context.Context, mp metrics.Params)
	SendJSONMetrics(ctx context.Context, slice []metrics.Params)
	SendJSONMetric(ctx context.Context, mp metrics.Params)
}

type httpSender struct {
	client      http.Client
	address     string
	timeout     time.Duration
	rateLimiter *rate.Limiter
}

func NewHttpSender(cfg configs.HttConnectionConfig) *httpSender {
	client := http.Client{Timeout: cfg.ClientTimeout}
	rl := rate.NewLimiter(rate.Every(1*time.Second), cfg.RateLimit) //не больше rateLimit запросов в секунду
	return &httpSender{
		client:      client,
		address:     cfg.Address,
		timeout:     cfg.ReportTimeout,
		rateLimiter: rl,
	}
}

func (h *httpSender) SendMetric(ctx context.Context, mp metrics.Params) {
	h.rateLimiter.Wait(ctx)
	ctx2, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()
	url := FormatURL(h.address, mp)
	request, err := http.NewRequestWithContext(ctx2, http.MethodPost, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := h.client.Do(request)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	log.Infof("Metrics send")
	log.Debugf("Received status code: %v for post request to %v\n", response.StatusCode, url)
}

func (h *httpSender) SendJSONMetric(ctx context.Context, mp metrics.Params) {
	h.rateLimiter.Wait(ctx)
	ctx2, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()
	url := fmt.Sprintf("http://%s/update", h.address) //TODO: wrap
	buf, err := json.Marshal(mp)
	if err != nil {
		log.Fatalf("Json marshal failed, %v\n", err)
		return
	}
	request, err := http.NewRequestWithContext(ctx2, http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := h.client.Do(request)
	if err != nil {
		log.Errorln(err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		log.Errorln(err)
	}
	log.Infof("Metrics send")
	log.Debugf("Received status code: %v for post request to %v\n", response.StatusCode, url)
}

func (h *httpSender) SendJSONMetrics(ctx context.Context, slice []metrics.Params) {
	err := h.rateLimiter.Wait(ctx)
	ctx2, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()
	if err != nil {
		log.Error(err)
		return
	}
	url := fmt.Sprintf("http://%s/updates/", h.address) //TODO: wrap
	buf, err := json.Marshal(slice)
	if err != nil {
		log.Errorf("Json marshal failed, %v\n", err)
		return
	}
	request, err := http.NewRequestWithContext(ctx2, http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")
	response, err := h.client.Do(request)
	if err != nil {
		log.Error(err)
		return
	}
	defer response.Body.Close()
	_, err = io.Copy(io.Discard, response.Body)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Metrics send")
	log.Debugf("Received status code: %v for post request to %v\n", response.StatusCode, url)
}

func FormatURL(addr string, mp metrics.Params) string {
	url := fmt.Sprintf("http://%v/update/%v/%v/", addr, mp.Type, mp.Name)
	switch mp.Type {
	case metrics.GaugeType:
		url += fmt.Sprint(*mp.ValueGauge)
	case metrics.CounterType:
		url += fmt.Sprint(*mp.ValueCounter)
	}
	return url
}
