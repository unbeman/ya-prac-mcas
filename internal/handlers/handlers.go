package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
		})
		router.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", ch.GetMetricHandler())
		})
	})
	return ch
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

func (ch *CollectorHandler) GetMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		params, err := getParams(request, ParamType, ParamName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		metric, err := ch.Repository.GetMetric(params[ParamType], params[ParamName])
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, storage.ErrInvalidType) {
			http.Error(writer, err.Error(), http.StatusNotImplemented)
			return
		}
		if errors.Is(err, storage.ErrInvalidValue) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
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
		params, err := getParams(request, ParamType, ParamName, ParamValue)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
			return
		}
		err = ch.Repository.SetMetric(params[ParamType], params[ParamName], params[ParamValue])
		if errors.Is(err, storage.ErrInvalidType) {
			http.Error(writer, err.Error(), http.StatusNotImplemented)
			return
		}
		if errors.Is(err, storage.ErrInvalidValue) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
