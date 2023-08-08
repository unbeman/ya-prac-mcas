package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/controller"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	mock_storage "github.com/unbeman/ya-prac-mcas/internal/storage/mock"
	"github.com/unbeman/ya-prac-mcas/internal/utils"
)

func TestCollectorHandler_GetMetricHandler(t *testing.T) {
	textContentType := "text/plain"
	type want struct {
		code          int
		checkResponse bool
		response      string
		contentType   string
	}

	type metric struct {
		name  string
		mType string
	}

	tests := []struct {
		name   string
		metric metric
		want   want
		setup  func(*mock_storage.MockRepository)
	}{
		{
			name:   "OK gauge",
			metric: metric{name: "OK", mType: "gauge"},
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				response:      "1.35",
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetGauge(gomock.Any(), gomock.Any()).
					Return(metrics.NewGauge("OK", 1.35), nil)
			},
		},
		{
			name:   "OK counter",
			metric: metric{name: "OK", mType: "counter"},
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				response:      "1",
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetCounter(gomock.Any(), gomock.Any()).
					Return(metrics.NewCounter("OK", 1), nil)
			},
		},
		{
			name:   "invalid type",
			metric: metric{name: "OK", mType: "fruit"},
			want: want{
				code:          http.StatusNotImplemented,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "invalid name",
			metric: metric{name: "", mType: "gauge"},
			want: want{
				code:          http.StatusBadRequest,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "not found gauge",
			metric: metric{name: "NotRegistered", mType: "gauge"},
			want: want{
				code:          http.StatusNotFound,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetGauge(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrNotFound)
			},
		},
		{
			name:   "not found counter",
			metric: metric{name: "NotRegistered", mType: "counter"},
			want: want{
				code:          http.StatusNotFound,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetCounter(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrNotFound)
			},
		},
		{
			name:   "internal error",
			metric: metric{name: "DogCount", mType: "counter"},
			want: want{
				code:          http.StatusInternalServerError,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetCounter(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := mock_storage.NewMockRepository(ctrl)
			tt.setup(mockRepository)

			ch := NewCollectorHandler(controller.NewController(mockRepository, ""), nil)

			request := utils.NewGetMetricTestRequest(tt.metric.mType, tt.metric.name)

			w := httptest.NewRecorder()
			ch.GetMetricHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			answer, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			if tt.want.checkResponse {
				assert.Equal(t, tt.want.response, string(answer))
			}
		})
	}
}

func TestCollectorHandler_GetMetricsHandler(t *testing.T) {
	textContentType := "text/plain"
	htmlContentType := "text/html"
	type want struct {
		code          int
		checkResponse bool
		response      string
		contentType   string
	}

	tests := []struct {
		name  string
		want  want
		setup func(*mock_storage.MockRepository)
	}{
		{
			name: "OK",
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				response: `GaugeA: 1.35
GaugeD: 0.001
CounterB: 10
CounterC: 12345
`,
				contentType: htmlContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetAll(gomock.Any()).Return([]metrics.Metric{
					metrics.NewGauge("GaugeA", 1.35),
					metrics.NewGauge("GaugeD", 0.001),
					metrics.NewCounter("CounterB", 10),
					metrics.NewCounter("CounterC", 12345),
				}, nil)
			},
		},
		{
			name: "internal error",
			want: want{
				code:          http.StatusInternalServerError,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetAll(gomock.Any()).Return(
					[]metrics.Metric{},
					errors.New("repository error"),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := mock_storage.NewMockRepository(ctrl)
			tt.setup(mockRepository)

			ch := NewCollectorHandler(controller.NewController(mockRepository, ""), nil)

			request := newGetMetricsTestRequest()

			w := httptest.NewRecorder()
			ch.GetMetricsHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			answer, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			if tt.want.checkResponse {
				assert.Equal(t, tt.want.response, string(answer))
			}
		})
	}
}

func newGetMetricsTestRequest() *http.Request {
	r := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)
	return r
}

func TestCollectorHandler_UpdateMetricHandler(t *testing.T) {
	textContentType := "text/plain"
	type want struct {
		code          int
		checkResponse bool
		response      string
		contentType   string
	}

	type metric struct {
		name  string
		mType string
		value string
	}

	tests := []struct {
		name   string
		metric metric
		want   want
		setup  func(*mock_storage.MockRepository)
	}{
		{
			name:   "OK gauge",
			metric: metric{name: "OK", mType: "gauge", value: "1.35"},
			want: want{
				code:          http.StatusOK,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().SetGauge(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(metrics.NewGauge("OK", 1.35), nil)
			},
		},
		{
			name:   "OK counter",
			metric: metric{name: "OK", mType: "counter", value: "1"},
			want: want{
				code:          http.StatusOK,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().AddCounter(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(metrics.NewCounter("OK", 2), nil)
			},
		},
		{
			name:   "invalid type",
			metric: metric{name: "OK", mType: "fruit"},
			want: want{
				code:          http.StatusNotImplemented,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "invalid name",
			metric: metric{name: "", mType: "gauge"},
			want: want{
				code:          http.StatusBadRequest,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "invalid value",
			metric: metric{name: "", mType: "counter", value: "a10.5"},
			want: want{
				code:          http.StatusBadRequest,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "internal error",
			metric: metric{name: "DogCount", mType: "counter", value: "4"},
			want: want{
				code:          http.StatusInternalServerError,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().AddCounter(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := mock_storage.NewMockRepository(ctrl)
			tt.setup(mockRepository)

			ch := NewCollectorHandler(controller.NewController(mockRepository, ""), nil)

			request := utils.NewUpdateMetricTestRequest(tt.metric.mType, tt.metric.name, tt.metric.value)

			w := httptest.NewRecorder()
			ch.UpdateMetricHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			answer, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			if tt.want.checkResponse {
				assert.Equal(t, tt.want.response, string(answer))
			}
		})
	}
}

func TestCollectorHandler_GetJSONMetricHandler(t *testing.T) {
	jsonContentType := "application/json"
	textContentType := "text/plain"

	type want struct {
		code          int
		checkResponse bool
		contentType   string
		metric        metrics.Params
	}

	tests := []struct {
		name   string
		metric metrics.Params
		want   want
		setup  func(*mock_storage.MockRepository)
	}{
		{
			name:   "OK gauge",
			metric: metrics.Params{Name: "OK", Type: "gauge"},
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				metric:        newGaugeParams("OK", 1.35, ""),
				contentType:   jsonContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetGauge(gomock.Any(), gomock.Any()).
					Return(metrics.NewGauge("OK", 1.35), nil)
			},
		},
		{
			name:   "OK counter",
			metric: metrics.Params{Name: "OK", Type: "counter"},
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				metric:        newCounterParams("OK", 1, ""),
				contentType:   jsonContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetCounter(gomock.Any(), gomock.Any()).
					Return(metrics.NewCounter("OK", 1), nil)
			},
		},
		{
			name:   "invalid type",
			metric: metrics.Params{Name: "OK", Type: "fruit"},
			want: want{
				code:          http.StatusNotImplemented,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "invalid name",
			metric: metrics.Params{Name: "", Type: "gauge"},
			want: want{
				code:          http.StatusBadRequest,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "not found gauge",
			metric: metrics.Params{Name: "NotRegistered", Type: "gauge"},
			want: want{
				code:          http.StatusNotFound,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetGauge(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrNotFound)
			},
		},
		{
			name:   "not found counter",
			metric: metrics.Params{Name: "NotRegistered", Type: "counter"},
			want: want{
				code:          http.StatusNotFound,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetCounter(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrNotFound)
			},
		},
		{
			name:   "internal error",
			metric: metrics.Params{Name: "DogCount", Type: "counter"},
			want: want{
				code:          http.StatusInternalServerError,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().GetCounter(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := mock_storage.NewMockRepository(ctrl)
			tt.setup(mockRepository)

			ch := NewCollectorHandler(controller.NewController(mockRepository, ""), nil)

			request := newGetMetricJSONTestRequest(tt.metric)

			w := httptest.NewRecorder()
			ch.GetJSONMetricHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			answer, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			if tt.want.checkResponse {
				var resultMetric metrics.Params
				json.Unmarshal(answer, &resultMetric)
				assert.Equal(t, tt.want.metric, resultMetric)
			}
		})
	}
}

func newGetMetricJSONTestRequest(metric metrics.Params) *http.Request {
	body, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}
	r := httptest.NewRequest(
		http.MethodPost,
		"/value",
		bytes.NewBuffer(body),
	)
	return r
}

func newGaugeParams(name string, valueGauge float64, hash string) metrics.Params {
	return metrics.Params{Name: name, Type: metrics.GaugeType, ValueGauge: &valueGauge, Hash: hash}
}

func newCounterParams(name string, valueCounter int64, hash string) metrics.Params {
	return metrics.Params{Name: name, Type: metrics.CounterType, ValueCounter: &valueCounter, Hash: hash}
}

func newInvalidTypeParams(name string, valueCounter int64, hash string) metrics.Params {
	return metrics.Params{Name: name, Type: "fruit", ValueCounter: &valueCounter, Hash: hash}
}

func TestCollectorHandler_UpdateJSONMetricHandler(t *testing.T) {
	jsonContentType := "application/json"
	textContentType := "text/plain"

	type want struct {
		code          int
		checkResponse bool
		contentType   string
		metric        metrics.Params
	}

	tests := []struct {
		name   string
		metric metrics.Params
		want   want
		setup  func(*mock_storage.MockRepository)
	}{
		{
			name:   "OK gauge",
			metric: newGaugeParams("OK", 1.35, ""),
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				metric:        newGaugeParams("OK", 1.35, ""),
				contentType:   jsonContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().SetGauge(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(metrics.NewGauge("OK", 1.35), nil)
			},
		},
		{
			name:   "OK counter",
			metric: newCounterParams("OK", 1, ""),
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				metric:        newCounterParams("OK", 2, ""),
				contentType:   jsonContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().AddCounter(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(metrics.NewCounter("OK", 2), nil)
			},
		},
		{
			name:   "invalid type",
			metric: newInvalidTypeParams("invalid", 3, ""),
			want: want{
				code:          http.StatusNotImplemented,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:   "invalid name",
			metric: newGaugeParams("", 0.0001, ""),
			want: want{
				code:          http.StatusBadRequest,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		//todo: invalid value
		{
			name:   "internal error",
			metric: newCounterParams("DogCount", 10, ""),
			want: want{
				code:          http.StatusInternalServerError,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().AddCounter(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := mock_storage.NewMockRepository(ctrl)
			tt.setup(mockRepository)

			ch := NewCollectorHandler(controller.NewController(mockRepository, ""), nil)

			request := newUpdateMetricJSONTestRequest(tt.metric)

			w := httptest.NewRecorder()
			ch.UpdateJSONMetricHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			answer, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			if tt.want.checkResponse {
				var resultMetric metrics.Params
				json.Unmarshal(answer, &resultMetric)
				assert.Equal(t, tt.want.metric, resultMetric)
			}
		})
	}
}

func newUpdateMetricJSONTestRequest(metric metrics.Params) *http.Request {
	body, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}
	r := httptest.NewRequest(
		http.MethodPost,
		"/update",
		bytes.NewBuffer(body),
	)
	return r
}

func TestCollectorHandler_UpdateJSONMetricsHandler(t *testing.T) {
	jsonContentType := "application/json"
	textContentType := "text/plain"

	type want struct {
		code          int
		checkResponse bool
		contentType   string
		metricsList   metrics.ParamsSlice
	}

	tests := []struct {
		name        string
		metricsList metrics.ParamsSlice
		want        want
		setup       func(*mock_storage.MockRepository)
	}{
		{
			name: "OK",
			metricsList: metrics.ParamsSlice{
				newGaugeParams("WaterPercentage", 0.35, ""),
				newGaugeParams("FoodPercentage", 0.8, ""),
				newCounterParams("DogsCount", 10, ""),
				newCounterParams("CatsCount", 3, ""),
			},
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				metricsList: metrics.ParamsSlice{
					newGaugeParams("WaterPercentage", 0.35, ""),
					newGaugeParams("FoodPercentage", 0.8, ""),
					newCounterParams("DogsCount", 12, ""),
					newCounterParams("CatsCount", 5, ""),
				},
				contentType: jsonContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().SetGauges(gomock.Any(), gomock.Any()).
					Return([]metrics.Gauge{
						metrics.NewGauge("WaterPercentage", 0.35),
						metrics.NewGauge("FoodPercentage", 0.8),
					}, nil)
				mR.EXPECT().AddCounters(gomock.Any(), gomock.Any()).
					Return([]metrics.Counter{
						metrics.NewCounter("DogsCount", 12),
						metrics.NewCounter("CatsCount", 5),
					}, nil)
			},
		},
		{
			name:        "OK only gauge",
			metricsList: metrics.ParamsSlice{newGaugeParams("OK", 1.35, "")},
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				metricsList:   metrics.ParamsSlice{newGaugeParams("OK", 1.35, "")},
				contentType:   jsonContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().SetGauges(gomock.Any(), gomock.Any()).
					Return([]metrics.Gauge{metrics.NewGauge("OK", 1.35)}, nil)
			},
		},
		{
			name:        "OK only counter",
			metricsList: metrics.ParamsSlice{newCounterParams("OK", 1, "")},
			want: want{
				code:          http.StatusOK,
				checkResponse: true,
				metricsList:   metrics.ParamsSlice{newCounterParams("OK", 2, "")},
				contentType:   jsonContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().AddCounters(gomock.Any(), gomock.Any()).
					Return([]metrics.Counter{metrics.NewCounter("OK", 2)}, nil)
			},
		},
		{
			name:        "invalid type",
			metricsList: metrics.ParamsSlice{newInvalidTypeParams("invalid", 3, "")},
			want: want{
				code:          http.StatusNotImplemented,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		{
			name:        "invalid name",
			metricsList: metrics.ParamsSlice{newGaugeParams("", 0.0001, "")},
			want: want{
				code:          http.StatusBadRequest,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {},
		},
		//todo: invalid value
		{
			name:        "internal error",
			metricsList: metrics.ParamsSlice{newCounterParams("DogCount", 10, "")},
			want: want{
				code:          http.StatusInternalServerError,
				checkResponse: false,
				contentType:   textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().AddCounters(gomock.Any(), gomock.Any()).
					Return([]metrics.Counter{}, errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := mock_storage.NewMockRepository(ctrl)
			tt.setup(mockRepository)

			ch := NewCollectorHandler(controller.NewController(mockRepository, ""), nil)

			request := newUpdatesMetricsJSONTestRequest(tt.metricsList)

			w := httptest.NewRecorder()
			ch.UpdateJSONMetricsHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			answer, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			if tt.want.checkResponse {
				var resultMetrics metrics.ParamsSlice
				json.Unmarshal(answer, &resultMetrics)
				assert.ElementsMatchf(t, resultMetrics, tt.want.metricsList,
					"SetGauges() got = %v, want %v", resultMetrics, tt.want.metricsList)
			}
		})
	}
}

func newUpdatesMetricsJSONTestRequest(metric metrics.ParamsSlice) *http.Request {
	body, err := json.Marshal(metric)
	if err != nil {
		panic(err)
	}
	r := httptest.NewRequest(
		http.MethodPost,
		"/update",
		bytes.NewBuffer(body),
	)
	return r
}

func TestPingHandler(t *testing.T) {
	textContentType := "text/plain"
	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name  string
		want  want
		setup func(*mock_storage.MockRepository)
	}{
		{
			name: "OK",
			want: want{
				code:        http.StatusOK,
				contentType: textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().Ping(gomock.Any()).Return(nil)
			},
		},
		{
			name: "internal error",
			want: want{
				code:        http.StatusInternalServerError,
				contentType: textContentType,
			},
			setup: func(mR *mock_storage.MockRepository) {
				mR.EXPECT().Ping(gomock.Any()).Return(errors.New("repository error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepository := mock_storage.NewMockRepository(ctrl)
			tt.setup(mockRepository)

			ch := NewCollectorHandler(controller.NewController(mockRepository, ""), nil)

			request := newPingTestRequest()

			w := httptest.NewRecorder()
			ch.PingHandler(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Contains(t, result.Header.Get("Content-Type"), tt.want.contentType)
			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func newPingTestRequest() *http.Request {
	return httptest.NewRequest(http.MethodGet, "/ping", nil)
}

func newBenchmarkHandler(repo storage.Repository) *CollectorHandler {
	ramRepo := storage.NewRAMRepository()
	ch := NewCollectorHandler(controller.NewController(ramRepo, ""), nil)
	return ch
}

func BenchmarkUpdateHandlers(b *testing.B) {
	ramRepo := storage.NewRAMRepository()
	handlerWithRAM := newBenchmarkHandler(ramRepo)

	pgRepo, _ := storage.NewPostgresRepository(configs.PostgresConfig{
		DSN:          "postgresql://postgres:1211@localhost:5432/mcas",
		MigrationDir: configs.PGMigrationDirDefault},
	)
	handlerWithPG := newBenchmarkHandler(pgRepo)
	b.Run("RAM UpdateJSONMetricsHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := newUpdatesMetricsJSONTestRequest(metrics.ParamsSlice{
				newGaugeParams("WaterPercentage", 0.35, ""),
				newGaugeParams("FoodPercentage", 0.8, ""),
				newCounterParams("DogsCount", 10, ""),
				newCounterParams("CatsCount", 3, ""),
			})
			w := httptest.NewRecorder()
			b.StartTimer()

			handlerWithRAM.UpdateJSONMetricsHandler(w, request)
		}
	})

	b.Run("PG UpdateJSONMetricsHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := newUpdatesMetricsJSONTestRequest(metrics.ParamsSlice{
				newGaugeParams("WaterPercentage", 0.35, ""),
				newGaugeParams("FoodPercentage", 0.8, ""),
				newCounterParams("DogsCount", 10, ""),
				newCounterParams("CatsCount", 3, ""),
			})
			w := httptest.NewRecorder()
			b.StartTimer()

			handlerWithPG.UpdateJSONMetricsHandler(w, request)
		}
	})

	b.Run("RAM UpdateJSONMetricHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := newUpdateMetricJSONTestRequest(newGaugeParams("WaterPercentage", 0.35, ""))
			w := httptest.NewRecorder()
			b.StartTimer()

			handlerWithRAM.UpdateJSONMetricHandler(w, request)
		}
	})

	b.Run("PG UpdateJSONMetricHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := newUpdateMetricJSONTestRequest(newGaugeParams("WaterPercentage", 0.35, ""))
			w := httptest.NewRecorder()
			b.StartTimer()

			handlerWithPG.UpdateJSONMetricHandler(w, request)
		}
	})

	b.Run("RAM UpdateMetricHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := utils.NewUpdateMetricTestRequest("gauge", "WaterPercentage", "0.35")
			w := httptest.NewRecorder()
			b.StartTimer()

			handlerWithRAM.UpdateJSONMetricHandler(w, request)
		}
	})

	b.Run("PG UpdateMetricHandler", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			request := utils.NewUpdateMetricTestRequest("gauge", "WaterPercentage", "0.35")
			w := httptest.NewRecorder()
			b.StartTimer()

			handlerWithPG.UpdateJSONMetricHandler(w, request)
		}
	})
}
