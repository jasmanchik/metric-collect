package handlers

import (
	"errors"
	"fmt"
	httpMetric "http-metric/internal/metrics/http"
	"http-metric/internal/platform/web"
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

func (h *HttpMetricHandler) Ping(w http.ResponseWriter, r *http.Request) error {
	l := h.log.With(slog.String("op", "router.ping"))
	defer func() {
		if err := recover(); err != nil {
			h.metric.IncServerErrors()
			l.Error("panic", "error", err) //nolint:govet
			err := web.RespondError(w)
			if err != nil {
				l.Error("defer response error:", slog.String("error", err.Error())) //nolint:govet
			}
		}
	}()
	err := h.metric.Ping()

	h.metric.IncTotalRequests()
	var errMetric *httpMetric.MetricHttpError
	if err != nil && errors.As(err, &errMetric) {
		if errMetric.Code == http.StatusBadRequest {
			h.metric.IncTotalErrors()
		}
		err = web.Response(w, nil, http.StatusBadRequest)
		if err != nil {
			l.Warn("response error:", slog.String("error", err.Error())) //nolint:govet
			return fmt.Errorf("response error: %w", err)
		}
	}
	err = web.Response(w, nil, http.StatusNoContent)
	if err != nil {
		l.Error("response error:", slog.String("error", err.Error())) //nolint:govet
		return fmt.Errorf("response error: %w", err)
	}

	return nil
}

func (h *HttpMetricHandler) RequestCounter(w http.ResponseWriter, r *http.Request) error {
	metricResult, err := h.metric.GetMetricList()
	log := h.log.With(slog.String("op", "router.RequestCounter"))
	if err != nil {
		log.Error("Response error: ", err.Error()) //nolint:govet
	}

	err = web.Response(w, metricResult, http.StatusOK)
	if err != nil {
		log.Error("Response error: ", err.Error()) //nolint:govet
		return fmt.Errorf("response error: %w", err)
	}
	return nil
}
