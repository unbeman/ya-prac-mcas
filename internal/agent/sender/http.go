// Package sender describes connection with metrics server.
package sender

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/utils"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type httpSender struct {
	client      http.Client
	address     string
	timeout     time.Duration
	rateLimiter *rate.Limiter
	publicKey   *rsa.PublicKey
}

func NewHTTPSender(cfg configs.ConnectionConfig, pubKey *rsa.PublicKey) (*httpSender, error) {
	client := http.Client{Timeout: cfg.ClientTimeout}
	rl := rate.NewLimiter(rate.Every(defaultRate), cfg.RateTokensCount) // не больше RateTokensCount запросов в секунду.
	return &httpSender{
		client:      client,
		address:     cfg.Address,
		timeout:     cfg.ReportTimeout,
		rateLimiter: rl,
		publicKey:   pubKey,
	}, nil
}

func (h *httpSender) SendMetric(ctx context.Context, mp metrics.Params) {
	h.rateLimiter.Wait(ctx)
	ctx2, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	url := FormatURL(h.address, mp)

	request, err := http.NewRequestWithContext(ctx2, http.MethodPost, url, nil)
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")

	ip, err := utils.GetOutboundIP()
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("X-Real-IP", ip)

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
		log.Errorf("json marshal failed, %v", err)
		return
	}

	var encryptedKey string

	if h.publicKey != nil {
		buf, encryptedKey, err = utils.GetEncryptedMessage(h.publicKey, buf)
		if err != nil {
			log.Errorf("encryption err, %v", err)
			return
		}
	}

	request, err := http.NewRequestWithContext(ctx2, http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")

	request.Header.Set("Encrypted-Key", encryptedKey)

	ip, err := utils.GetOutboundIP()
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("X-Real-IP", ip)

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

func (h *httpSender) SendMetrics(ctx context.Context, slice metrics.ParamsSlice) {
	err := h.rateLimiter.Wait(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	ctx2, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	url := fmt.Sprintf("http://%s/updates/", h.address) //TODO: wrap
	buf, err := json.Marshal(slice)
	if err != nil {
		log.Errorf("json marshal failed, %v", err)
		return
	}

	var encryptedKey string

	if h.publicKey != nil {
		buf, encryptedKey, err = utils.GetEncryptedMessage(h.publicKey, buf)
		if err != nil {
			log.Errorf("encryption err, %v", err)
			return
		}
	}

	request, err := http.NewRequestWithContext(ctx2, http.MethodPost, url, bytes.NewBuffer(buf))
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("Content-Type", "text/plain")

	request.Header.Set("Encrypted-Key", encryptedKey)

	ip, err := utils.GetOutboundIP()
	if err != nil {
		log.Error(err)
		return
	}
	request.Header.Set("X-Real-IP", ip)

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

	log.Info("Metrics send")
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
