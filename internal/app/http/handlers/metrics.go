package handlers

import (
	"errors"
	httpMetric "http-metric/internal/metrics/http"
	"log/slog"
	"net/http"
)

type HttpMetricHandler struct {
	log    *slog.Logger
	metric *httpMetric.RequestsMetric
}

func NewHttpMetric(logger *slog.Logger, m *httpMetric.RequestsMetric) *HttpMetricHandler {
	return &HttpMetricHandler{
		log:    logger,
		metric: m,
	}
}

func (h *HttpMetricHandler) Ping(w http.ResponseWriter, r *http.Request) (any, error) {
	err := h.metric.Ping()
	if err != nil {
		var errMetric *httpMetric.MetricHttpError
		if errors.As(err, &errMetric) {
			return nil, &ErrorHTTP{Code: errMetric.Code, ErrorText: errMetric.ErrorText}
		}
		return nil, err
	}

	return nil, nil
}

func (h *HttpMetricHandler) RequestCounter(w http.ResponseWriter, r *http.Request) (any, error) {
	metricResult, err := h.metric.GetMetricList()
	log := h.log.With(slog.String("op", "router.RequestCounter"))
	if err != nil {
		log.Error("get metric list error: ", slog.String("error", err.Error())) //nolint:govet
		return nil, err
	}

	return metricResult, nil
}
