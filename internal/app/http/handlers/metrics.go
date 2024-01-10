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
	log     *slog.Logger
	counter *httpMetric.RequestsMetric
}

func NewHttpMetric(logger *slog.Logger, m *httpMetric.RequestsMetric) *HttpMetricHandler {
	return &HttpMetricHandler{
		log:     logger,
		counter: m,
	}
}

func (m *HttpMetricHandler) Ping(w http.ResponseWriter, r *http.Request) error {
	l := m.log.With(slog.String("op", "router.ping"))
	defer func() {
		if err := recover(); err != nil {
			m.counter.TotalServerErrors++
			l.Error("panic", "error", err) //nolint:govet
			err := web.RespondError(w)
			if err != nil {
				l.Error("defer response error:", slog.String("error", err.Error())) //nolint:govet
			}
		}
	}()
	err := m.counter.Ping()
	m.counter.TotalRequests++
	var errMetric *httpMetric.MetricHttpError
	if err != nil && errors.As(err, &errMetric) {
		if errMetric.Code == http.StatusBadRequest {
			m.counter.TotalErrors++
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

func (m *HttpMetricHandler) RequestCounter(w http.ResponseWriter, r *http.Request) error {
	metricResult, err := m.counter.GetMetricList()
	log := m.log.With(slog.String("op", "router.RequestCounter"))
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

func (m *HttpMetricHandler) Metrics(w http.ResponseWriter, r *http.Request) error {

	return nil
}
