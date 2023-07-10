package utils

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
)

func NewGetMetricTestRequest(mType, name string) *http.Request {
	r := httptest.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/values/%v/%v", mType, name),
		nil,
	)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("type", mType)
	ctx.URLParams.Add("name", name)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func NewUpdateMetricTestRequest(mType, name, value string) *http.Request {
	r := httptest.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/update/%v/%v/%v", mType, name, value),
		nil,
	)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("type", mType)
	ctx.URLParams.Add("name", name)
	ctx.URLParams.Add("value", value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}
