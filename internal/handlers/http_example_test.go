package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"

	"github.com/unbeman/ya-prac-mcas/internal/controller"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	"github.com/unbeman/ya-prac-mcas/internal/utils"
)

func ExampleCollectorHandler_GetMetricHandler() {
	ch := NewCollectorHandler(controller.NewController(storage.NewRAMRepository(), ""), nil, nil)
	ch.controller.UpdateMetric(context.TODO(), newCounterParams("Dog", 10, ""))

	request := utils.NewGetMetricTestRequest("counter", "Dog")
	w := httptest.NewRecorder()

	ch.GetMetricHandler(w, request)

	result := w.Result()
	defer result.Body.Close()

	fmt.Println(result.StatusCode)

	value, _ := io.ReadAll(result.Body)
	fmt.Println(string(value))

	// Output:
	// 200
	// 10
}

func ExampleCollectorHandler_GetMetricsHandler() {
	ch := NewCollectorHandler(controller.NewController(storage.NewRAMRepository(), ""), nil, nil)
	ch.controller.UpdateMetric(context.TODO(), newCounterParams("Dog", 10, ""))
	ch.controller.UpdateMetric(context.TODO(), newGaugeParams("WaterPercent", 0.8, ""))

	request := newGetMetricsTestRequest()
	w := httptest.NewRecorder()

	ch.GetMetricsHandler(w, request)

	result := w.Result()
	defer result.Body.Close()

	fmt.Println(result.StatusCode)

	value, _ := io.ReadAll(result.Body)
	fmt.Println(string(value))

	// Output:
	// 200
	// Dog: 10
	// WaterPercent: 0.8
}

func ExampleCollectorHandler_GetJSONMetricHandler() {
	ch := NewCollectorHandler(controller.NewController(storage.NewRAMRepository(), ""), nil, nil)
	ch.controller.UpdateMetric(context.TODO(), newCounterParams("Dog", 10, ""))

	request := newGetMetricJSONTestRequest(metrics.Params{Name: "Dog", Type: "counter"})
	w := httptest.NewRecorder()

	ch.GetJSONMetricHandler(w, request)

	result := w.Result()
	defer result.Body.Close()

	fmt.Println(result.StatusCode)

	value, _ := io.ReadAll(result.Body)
	fmt.Println(string(value))

	// Output:
	// 200
	// {"id":"Dog","type":"counter","delta":10}
}

func ExampleCollectorHandler_UpdateJSONMetricHandler() {
	ch := NewCollectorHandler(controller.NewController(storage.NewRAMRepository(), ""), nil, nil)
	ch.controller.UpdateMetric(context.TODO(), newCounterParams("Dog", 10, ""))

	request := newUpdateMetricJSONTestRequest(metrics.Params{
		Name:         "Dog",
		Type:         "counter",
		ValueCounter: func(n int64) *int64 { return &n }(5),
	})
	w := httptest.NewRecorder()

	ch.UpdateJSONMetricHandler(w, request)

	result := w.Result()
	defer result.Body.Close()

	fmt.Println(result.StatusCode)

	value, _ := io.ReadAll(result.Body)
	fmt.Println(string(value))

	// Output:
	// 200
	// {"id":"Dog","type":"counter","delta":15}
}

func ExampleCollectorHandler_UpdateJSONMetricsHandler() {
	ch := NewCollectorHandler(controller.NewController(storage.NewRAMRepository(), ""), nil, nil)

	request := newUpdatesMetricsJSONTestRequest([]metrics.Params{
		{
			Name:         "Dog",
			Type:         "counter",
			ValueCounter: func(n int64) *int64 { return &n }(5),
		},
		{
			Name:       "WaterPercent",
			Type:       "gauge",
			ValueGauge: func(n float64) *float64 { return &n }(0.8),
		},
	})
	w := httptest.NewRecorder()

	ch.UpdateJSONMetricsHandler(w, request)

	result := w.Result()
	defer result.Body.Close()

	fmt.Println(result.StatusCode)

	value, _ := io.ReadAll(result.Body)
	fmt.Println(string(value))

	// Output:
	// 200
	// [{"id":"WaterPercent","type":"gauge","value":0.8},{"id":"Dog","type":"counter","delta":5}]
}

func ExampleCollectorHandler_UpdateMetricHandler() {
	ch := NewCollectorHandler(controller.NewController(storage.NewRAMRepository(), ""), nil, nil)
	ch.controller.UpdateMetric(context.TODO(), newCounterParams("Dog", 10, ""))

	request := utils.NewUpdateMetricTestRequest("counter", "Dog", "5")
	w := httptest.NewRecorder()

	ch.UpdateMetricHandler(w, request)

	result := w.Result()
	defer result.Body.Close()

	fmt.Println(result.StatusCode)

	// Output:
	// 200
}
