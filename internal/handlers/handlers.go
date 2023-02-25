package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/parser"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

const (
	ParamType  = "type"
	ParamName  = "name"
	ParamValue = "value"
)

type CollectorHandler struct {
	*chi.Mux
	Repository storage.Repository
}

func NewCollectorHandler(repository storage.Repository) *CollectorHandler {
	ch := &CollectorHandler{
		Mux:        chi.NewMux(),
		Repository: repository,
	}
	ch.Use(middleware.RequestID)
	ch.Use(middleware.RealIP)
	ch.Use(middleware.Logger)
	ch.Use(middleware.Recoverer)
	ch.Route("/", func(router chi.Router) {
		router.Get("/", ch.GetMetricsHandler())
		router.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", ch.UpdateMetricHandler())
			r.Post("/", ch.UpdateJsonMetricHandler())
		})
		router.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", ch.GetMetricHandler())
			r.Post("/", ch.GetJsonMetricHandler())
		})
	})
	return ch
}

func (ch *CollectorHandler) GetMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		params, err := parser.ParseURI(request, ParamType, ParamName)
		if errors.Is(err, parser.ErrInvalidType) {
			http.Error(writer, err.Error(), http.StatusNotImplemented)
			return
		}
		if errors.Is(err, parser.ErrInvalidValue) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		var metric metrics.Metric
		switch params.Type {
		case metrics.GaugeType:
			metric = ch.Repository.GetGauge(params.Name)
		case metrics.CounterType:
			metric = ch.Repository.GetCounter(params.Name)
		}
		if metric == nil {
			http.Error(writer, "metric not found", http.StatusNotFound)
			return
		}

		_, err = writer.Write([]byte(metric.GetValue()))
		if err != nil {
			log.Printf("Write failed, %v", err)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (ch *CollectorHandler) GetMetricsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
		var b strings.Builder
		for _, metric := range ch.Repository.GetAll() {
			_, err := fmt.Fprintf(&b, "%v: %v\n", metric.GetName(), metric.GetValue())
			if err != nil {
				log.Printf("GetMetricsHandler: can't build metrics list with values %v %v, reason: %v",
					metric.GetName(), metric.GetValue(), err)
			}
		}
		_, err := writer.Write([]byte(b.String()))
		if err != nil {
			log.Printf("Write failed, %v", err)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (ch *CollectorHandler) UpdateMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		params, err := parser.ParseURI(request, ParamType, ParamName, ParamValue)
		if errors.Is(err, parser.ErrInvalidType) {
			http.Error(writer, err.Error(), http.StatusNotImplemented)
			return
		}
		if errors.Is(err, parser.ErrInvalidValue) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		switch params.Type {
		case metrics.GaugeType:
			_ = ch.Repository.SetGauge(params.Name, *params.ValueGauge)
		case metrics.CounterType:
			_ = ch.Repository.AddCounter(params.Name, *params.ValueCounter)
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (ch *CollectorHandler) GetJsonMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		var jsonMetric *parser.JSONMetric
		if err := json.NewDecoder(request.Body).Decode(&jsonMetric); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		params, err := parser.ParseJSON(jsonMetric, false)
		if errors.Is(err, parser.ErrInvalidType) {
			http.Error(writer, err.Error(), http.StatusNotImplemented)
			return
		}
		if errors.Is(err, parser.ErrInvalidValue) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		var metric metrics.Metric
		switch params.Type {
		case metrics.GaugeType:
			metric = ch.Repository.GetGauge(params.Name)
		case metrics.CounterType:
			metric = ch.Repository.GetCounter(params.Name)
		}

		if metric == nil {
			http.Error(writer, "metric not found", http.StatusNotFound)
			return
		}

		jsonMetric = parser.MetricToJSON(metric)
		if err := json.NewEncoder(writer).Encode(jsonMetric); err != nil {
			log.Printf("Write failed, %v", err)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (ch *CollectorHandler) UpdateJsonMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		var jsonMetric *parser.JSONMetric
		if err := json.NewDecoder(request.Body).Decode(&jsonMetric); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		params, err := parser.ParseJSON(jsonMetric, true)
		if errors.Is(err, parser.ErrInvalidType) {
			http.Error(writer, err.Error(), http.StatusNotImplemented)
			return
		}
		if errors.Is(err, parser.ErrInvalidValue) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		var metric metrics.Metric
		switch params.Type {
		case metrics.GaugeType:
			metric = ch.Repository.SetGauge(params.Name, *params.ValueGauge)
		case metrics.CounterType:
			metric = ch.Repository.AddCounter(params.Name, *params.ValueCounter)
		}
		log.Printf("JSON UPD %#v", metric)
		jsonMetric = parser.MetricToJSON(metric)
		if err := json.NewEncoder(writer).Encode(&jsonMetric); err != nil {
			log.Printf("Write failed, %v\n", err)
			return
		}
		log.Printf("UpdateJsonMetricHandler Output: %v\n", jsonMetric)
		writer.WriteHeader(http.StatusOK)
	}
}
