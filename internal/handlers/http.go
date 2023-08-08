// Package handlers describes server router and handlers methods.
package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/internal/controller"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type CollectorHandler struct {
	*chi.Mux
	controller *controller.Controller
}

func NewCollectorHandler(controller *controller.Controller, privateRSAKey *rsa.PrivateKey) *CollectorHandler {
	ch := &CollectorHandler{
		Mux:        chi.NewMux(),
		controller: controller,
	}

	ch.Use(middleware.RequestID)
	ch.Use(middleware.RealIP)
	ch.Use(logger.Logger("router", log.New()))
	ch.Use(middleware.Recoverer)
	ch.Use(GZipMiddleware)
	ch.Route("/", func(router chi.Router) {
		router.Get("/", ch.GetMetricsHandler)

		router.Post("/update/{type}/{name}/{value}", ch.UpdateMetricHandler)

		router.Group(func(r chi.Router) {
			r.Use(DecryptMiddleware(privateRSAKey))
			r.Post("/updates/", ch.UpdateJSONMetricsHandler)
			r.Post("/update/", ch.UpdateJSONMetricHandler)
		})

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
	if err != nil {
		ch.processError(writer, err)
		return
	}

	metric, err := ch.controller.GetMetric(request.Context(), params)
	if err != nil {
		ch.processError(writer, err)
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

	metricSlice, err := ch.controller.GetAll(request.Context())
	if err != nil {
		ch.processError(writer, err)
		return
	}

	for _, metric := range metricSlice {
		_, err = fmt.Fprintf(&b, "%v: %v\n", metric.GetName(), metric.GetValue())
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
	if err != nil {
		ch.processError(writer, err)
		return
	}

	_, err = ch.controller.UpdateMetric(request.Context(), params)
	if err != nil {
		ch.processError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) GetJSONMetricHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	params, err := metrics.ParseJSON(request.Body, metrics.PName, metrics.PType)
	if err != nil {
		ch.processError(writer, err)
		return
	}

	metric, err := ch.controller.GetMetric(request.Context(), params)
	if err != nil {
		ch.processError(writer, err)
		return
	}

	params = metric.ToParams()
	params.Hash = ch.controller.GetHash(metric)
	if err := json.NewEncoder(writer).Encode(params); err != nil {
		log.Errorf("Write failed, %v", err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) UpdateJSONMetricHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params, err := metrics.ParseJSON(request.Body, metrics.PName, metrics.PType, metrics.PValue)
	if err != nil {
		ch.processError(writer, err)
		return
	}

	metric, err := ch.controller.UpdateMetric(request.Context(), params)
	if err != nil {
		ch.processError(writer, err)
		return
	}

	params = metric.ToParams()
	params.Hash = ch.controller.GetHash(metric)
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

	if err != nil {
		ch.processError(writer, err)
		return
	}

	metricsParams, err := ch.controller.UpdateMetrics(request.Context(), paramsSlice)
	if errors.Is(err, controller.ErrInvalidHash) {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	if err != nil {
		ch.processError(writer, err)
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

	err := ch.controller.Ping(request.Context())
	if err != nil {
		ch.processError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) processError(w http.ResponseWriter, err error) {
	var httpCode int
	switch {
	case errors.Is(err, controller.ErrInvalidHash):
		httpCode = http.StatusBadRequest
	case errors.Is(err, metrics.ErrInvalidType):
		httpCode = http.StatusNotImplemented
	case errors.Is(err, metrics.ErrInvalidValue):
		httpCode = http.StatusBadRequest
	case errors.Is(err, metrics.ErrParseURI):
		httpCode = http.StatusMethodNotAllowed
	case errors.Is(err, metrics.ErrParseJSON):
		httpCode = http.StatusBadRequest
	case errors.Is(err, storage.ErrNotFound):
		httpCode = http.StatusNotFound
	default:
		httpCode = http.StatusInternalServerError
	}
	http.Error(w, err.Error(), httpCode)
}
