package handlers

//TODO: rename file
import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	"html/template"
	"net/http"
	"strconv"
)

type CollectorHandler struct {
	*chi.Mux
	Storage storage.Repository
}

func NewCollectorHandler(stor storage.Repository) *CollectorHandler {
	ch := &CollectorHandler{
		Mux:     chi.NewMux(),
		Storage: stor,
	}
	ch.Use(middleware.RequestID)
	ch.Use(middleware.RealIP)
	ch.Use(middleware.Logger)
	ch.Use(middleware.Recoverer)
	ch.Route("/", func(rout chi.Router) {
		rout.Get("/", ch.GetMetricsHandler())
		rout.Route("/update", func(r chi.Router) {
			r.Post("/{type}/{name}/{value}", ch.UpdateMetricHandler())
		})
		rout.Route("/value", func(r chi.Router) {
			r.Get("/{type}/{name}", ch.GetMetricHandler())
		})
	})
	return ch
}

//TODO: split http handler and business logic into layers

func (ch *CollectorHandler) GetMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mType := chi.URLParam(request, "type")
		mName := chi.URLParam(request, "name")
		var metric metrics.Metric
		var ok bool
		switch mType {
		case metrics.GaugeTypeName:
			metric, ok = ch.Storage.GetGauge(mName)
		case metrics.CounterTypeName:
			metric, ok = ch.Storage.GetCounter(mName)
		default:
			writer.Header().Set("Content-Type", "text/plain")
			http.Error(writer, fmt.Sprintf("invalid type %v", mType), http.StatusMethodNotAllowed)
			return
		}
		if !ok {
			writer.Header().Set("Content-Type", "text/plain")
			http.Error(writer, fmt.Sprintf("metric %v not found", mName), http.StatusNotFound)
			return
		}
		writer.Write([]byte(metric.GetValue())) //TODO: check error
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
	}
}

//TODO: create html template with header and body, move to file

var ListTemplate = "{{ range $key, $value := . }}{{ $value.GetName }}: {{ $value.GetValue }}\n{{ end }}"

func (ch *CollectorHandler) GetMetricsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		metricsMap, _ := ch.Storage.GetAll()
		t, _ := template.New("metrics list").Parse(ListTemplate)
		_ = t.Execute(writer, metricsMap)
		writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
		writer.WriteHeader(http.StatusOK)
	}
}

func (ch *CollectorHandler) UpdateMetricHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mType := chi.URLParam(request, "type")
		switch mType {
		case metrics.GaugeTypeName:
			ch.UpdateGauge(writer, request)
		case metrics.CounterTypeName:
			ch.UpdateCounter(writer, request)
		default:
			writer.Header().Set("Content-Type", "text/plain")
			http.Error(writer, fmt.Sprintf("invalid type %v", mType), http.StatusNotImplemented)
			return
		}
	}
}

func (ch *CollectorHandler) UpdateCounter(writer http.ResponseWriter, request *http.Request) {
	mName := chi.URLParam(request, "name")
	mValue, err := strconv.ParseInt(chi.URLParam(request, "value"), 10, 64) //TODO: wrap and move to parser
	if err != nil {
		writer.Header().Set("Content-Type", "text/plain")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	counter, ok := ch.Storage.GetCounter(mName)
	if !ok {
		counter = metrics.NewCounter(mName)
		ch.Storage.UpdateCounterRepo(counter)
	}
	counter.Add(mValue)

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}

func (ch *CollectorHandler) UpdateGauge(writer http.ResponseWriter, request *http.Request) {
	mName := chi.URLParam(request, "name")
	mValue, err := strconv.ParseFloat(chi.URLParam(request, "value"), 64) //TODO: wrap and move to parser
	if err != nil {
		writer.Header().Set("Content-Type", "text/plain")
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	gauge, ok := ch.Storage.GetGauge(mName)
	if !ok {
		gauge = metrics.NewGauge(mName)
		ch.Storage.UpdateGaugeRepo(gauge)
	}
	gauge.Set(mValue)
	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
}
