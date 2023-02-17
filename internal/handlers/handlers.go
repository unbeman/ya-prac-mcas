package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

const (
	ParamType  = "type"
	ParamName  = "name"
	ParamValue = "value"
)

type CollectorHandler struct {
	*chi.Mux
	Repository *storage.Repository
}

func NewCollectorHandler(repository *storage.Repository) (*CollectorHandler, error) {
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
		})
		router.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", ch.GetMetricHandler())
		})
	})
	return ch, nil
}

func getParams(request *http.Request, keys ...string) (map[string]string, error) {
	params := make(map[string]string)
	for _, key := range keys {
		value := chi.URLParam(request, key)
		if len(value) == 0 {
			return nil, fmt.Errorf("empty %v", key)
		}
		params[key] = value
	}
	return params, nil
}

//TODO: split http handler and business logic into layers

func (ch *CollectorHandler) GetMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		params, err := getParams(request, ParamType, ParamName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
			return
		}
		var metric metrics.Metric
		switch params[ParamType] {
		case metrics.GaugeTypeName:
			metric = ch.Repository.Gauge.Get(params[ParamName])
		case metrics.CounterTypeName:
			metric = ch.Repository.Counter.Get(params[ParamName])
		default:
			http.Error(writer, fmt.Sprintf("invalid type %v", params[ParamType]), http.StatusMethodNotAllowed)
			return
		}
		if metric == nil {
			http.Error(writer, fmt.Sprintf("%v %v not found", params[ParamType], params[ParamName]), http.StatusNotFound)
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

// TODO: create html template with header and body, move to file

func (ch *CollectorHandler) GetMetricsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var b strings.Builder
		for _, metric := range ch.Repository.Gauge.GetAll() {
			fmt.Fprintf(&b, "%v: %v\n", metric.GetName(), metric.GetValue())
		}
		for _, metric := range ch.Repository.Counter.GetAll() {
			fmt.Fprintf(&b, "%v: %v\n", metric.GetName(), metric.GetValue())
		}
		_, err := writer.Write([]byte(b.String()))
		if err != nil {
			log.Printf("Write failed, %v", err)
			return
		}
		writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
		writer.WriteHeader(http.StatusOK)
	}
}

func (ch *CollectorHandler) UpdateMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		params, err := getParams(request, ParamType, ParamName, ParamValue)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
			return
		}
		switch params[ParamType] {
		case metrics.GaugeTypeName:
			err = ch.UpdateGauge(params)
		case metrics.CounterTypeName:
			err = ch.UpdateCounter(params)
		default:
			http.Error(writer, fmt.Sprintf("invalid type %v", params[ParamType]), http.StatusNotImplemented)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func (ch *CollectorHandler) UpdateCounter(params map[string]string) error {
	mValue, err := strconv.ParseInt(params[ParamValue], 10, 64) //TODO: wrap and move to parser
	if err != nil {
		return fmt.Errorf("invalid value %v: %w", params[ParamValue], err)
	}

	ch.Repository.Counter.Set(params[ParamName], mValue)
	return nil
}

func (ch *CollectorHandler) UpdateGauge(params map[string]string) error {
	mValue, err := strconv.ParseFloat(params[ParamValue], 64) //TODO: wrap and move to parser
	if err != nil {
		return fmt.Errorf("invalid value %v: %w", params[ParamValue], err)
	}

	ch.Repository.Gauge.Set(params[ParamName], mValue)
	return nil
}
