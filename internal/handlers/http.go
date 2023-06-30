package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type CollectorHandler struct {
	*chi.Mux
	Repository storage.Repository
	HashKey    []byte
}

func NewCollectorHandler(repository storage.Repository, key string) *CollectorHandler {
	ch := &CollectorHandler{
		Mux:        chi.NewMux(),
		Repository: repository,
		HashKey:    []byte(key),
	}

	ch.Use(middleware.RequestID)
	ch.Use(middleware.RealIP)
	ch.Use(logger.Logger("router", log.New()))
	ch.Use(middleware.Recoverer)
	ch.Use(GZipMiddleware)
	ch.Route("/", func(router chi.Router) {
		router.Get("/", ch.GetMetricsHandler)
		router.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", ch.UpdateMetricHandler)
			r.Post("/", ch.UpdateJSONMetricHandler)
		})
		router.Post("/updates/", ch.UpdateJSONMetricsHandler)
		router.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", ch.GetMetricHandler)
			r.Post("/", ch.GetJSONMetricHandler)
		})
		router.Get("/ping", ch.PingHandler)
	})
	return ch
}

func (ch *CollectorHandler) GetMetricHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	params, err := metrics.ParseURI(request, metrics.PType, metrics.PName)
	if errors.Is(err, metrics.ErrInvalidType) {
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	if errors.Is(err, metrics.ErrInvalidValue) {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	metric, err := ch.getMetric(request.Context(), params)
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(writer, "metric not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = writer.Write([]byte(metric.GetValue()))
	if err != nil {
		log.Errorf("Write failed, %v", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) GetMetricsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")

	var b strings.Builder

	metricSlice, err := ch.Repository.GetAll(request.Context())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, metric := range metricSlice {
		_, err := fmt.Fprintf(&b, "%v: %v\n", metric.GetName(), metric.GetValue())
		if err != nil {
			log.Errorf("GetMetricsHandler: can't build metrics list with values %v %v, reason: %v",
				metric.GetName(), metric.GetValue(), err)
		}
	}
	_, err = writer.Write([]byte(b.String()))
	if err != nil {
		log.Errorf("Write failed, %v", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) UpdateMetricHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")
	params, err := metrics.ParseURI(request, metrics.PType, metrics.PName, metrics.PValue)
	if errors.Is(err, metrics.ErrInvalidType) {
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	if errors.Is(err, metrics.ErrInvalidValue) {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	_, err = ch.updateMetric(request.Context(), params)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) GetJSONMetricHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	params, err := metrics.ParseJSON(request.Body, metrics.PName, metrics.PType)
	if errors.Is(err, metrics.ErrInvalidType) {
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	if errors.Is(err, metrics.ErrInvalidValue) {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	metric, err := ch.getMetric(request.Context(), params)
	if errors.Is(err, storage.ErrNotFound) {
		http.Error(writer, "metric not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	params = metric.ToParams()
	params.Hash = ch.getHash(metric)
	if err := json.NewEncoder(writer).Encode(params); err != nil {
		log.Errorf("Write failed, %v", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) UpdateJSONMetricHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	params, err := metrics.ParseJSON(request.Body, metrics.PName, metrics.PType, metrics.PValue)
	if errors.Is(err, metrics.ErrInvalidType) {
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	if errors.Is(err, metrics.ErrInvalidValue) {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	metric := metrics.NewMetricFromParams(params)

	if !ch.isValidHash(params.Hash, metric) {
		http.Error(writer, "invalid hash", http.StatusBadRequest)
		return
	}

	metric, err = ch.updateMetric(request.Context(), params)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	params = metric.ToParams()
	params.Hash = ch.getHash(metric)
	if err := json.NewEncoder(writer).Encode(&params); err != nil {
		log.Errorf("Write failed, %v\n", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) UpdateJSONMetricsHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	paramsSlice := metrics.ParamsSlice{}
	err := paramsSlice.ParseJSON(request.Body)

	if errors.Is(err, metrics.ErrInvalidType) {
		http.Error(writer, err.Error(), http.StatusNotImplemented)
		return
	}
	if errors.Is(err, metrics.ErrInvalidValue) {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	gauges := make([]metrics.Gauge, 0)
	counters := make([]metrics.Counter, 0)
	for _, params := range paramsSlice {
		metric := metrics.NewMetricFromParams(params)

		if !ch.isValidHash(params.Hash, metric) {
			http.Error(writer, "invalid hash", http.StatusBadRequest)
			return
		}

		switch metric.GetType() {
		case metrics.GaugeType:
			gauges = append(gauges, metric.(metrics.Gauge))
		case metrics.CounterType:
			counters = append(counters, metric.(metrics.Counter))
		}
	}

	metricsParams, err := ch.updateMetrics(request.Context(), gauges, counters)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(writer).Encode(&metricsParams); err != nil {
		log.Errorf("Write failed, %v\n", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) PingHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain")

	err := ch.Repository.Ping()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) getMetric(ctx context.Context, params metrics.Params) (metrics.Metric, error) { //TODO: controller layer
	var (
		metric metrics.Metric
		err    error
	)
	switch params.Type {
	case metrics.GaugeType:
		metric, err = ch.Repository.GetGauge(ctx, params.Name)
	case metrics.CounterType:
		metric, err = ch.Repository.GetCounter(ctx, params.Name)
	}
	return metric, err
}

func (ch *CollectorHandler) updateMetric(ctx context.Context, params metrics.Params) (metrics.Metric, error) { //TODO: controller layer
	var (
		err    error
		metric metrics.Metric
	)

	switch params.Type {
	case metrics.GaugeType:
		metric, err = ch.Repository.SetGauge(ctx, params.Name, *params.ValueGauge)
	case metrics.CounterType:
		metric, err = ch.Repository.AddCounter(ctx, params.Name, *params.ValueCounter)
	}
	return metric, err
}

func (ch *CollectorHandler) updateMetrics(
	ctx context.Context,
	gauges []metrics.Gauge,
	counters []metrics.Counter) (metrics.ParamsSlice, error) { //TODO: controller layer
	metricsParams := make(metrics.ParamsSlice, 0, len(gauges)+len(counters))

	if len(gauges) > 0 {
		updatedGauges, err := ch.Repository.SetGauges(ctx, gauges)
		if err != nil {
			return nil, err
		}

		for _, gauge := range updatedGauges {
			gp := gauge.ToParams()
			metricsParams = append(metricsParams, gp)
		}
	}

	if len(counters) > 0 {
		updatedCounters, err := ch.Repository.AddCounters(ctx, counters)
		if err != nil {
			return nil, err
		}

		for _, counter := range updatedCounters {
			cp := counter.ToParams()
			metricsParams = append(metricsParams, cp)
		}

	}
	return metricsParams, nil
}

func (ch *CollectorHandler) isValidHash(hash string, metric metrics.Metric) bool {
	if !ch.isKeySet() { //ключа нет
		return true
	}
	if !isHashSet(hash) {
		return false
	}

	return hash == metric.Hash(ch.HashKey)
}

func (ch *CollectorHandler) getHash(metric metrics.Metric) string {
	if !ch.isKeySet() {
		return ""
	}
	return metric.Hash(ch.HashKey)
}

func (ch *CollectorHandler) isKeySet() bool {
	return len(ch.HashKey) > 0
}

func isHashSet(hash string) bool {
	return len(hash) > 0
}
