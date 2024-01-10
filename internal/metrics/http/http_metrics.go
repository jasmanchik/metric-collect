package http

import (
	"errors"
	"log/slog"
	"math/rand"
	"net/http"
	"time"
)

type RequestMetricResult struct {
	TotalRequests     int `json:"total_requests_count"`
	TotalErrors       int `json:"total_request_errors_count"`
	TotalServerErrors int `json:"total_server_errors_count"`
}

type RequestsMetric struct {
	TotalRequests     int
	TotalErrors       int
	TotalServerErrors int
	log               *slog.Logger
}

func NewRequestCounter(log *slog.Logger) *RequestsMetric {
	return &RequestsMetric{
		log: log,
	}
}

func (r *RequestsMetric) Ping() error {
	const op = "router.ping"
	logger := r.log.With("op", op)

	sleepTime := rand.Intn(1000-100) + 100
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	errorRand := rand.Float64()
	if errorRand <= 0.05 {
		logger.Error("internal server error with 0.5% chance")
		panic(errors.New("internal server error"))
	}
	if errorRand <= 0.15 {
		logger.Warn("bad Request error")
		return &MetricHttpError{
			ErrorText: "Bad Request error",
			Code:      http.StatusBadRequest,
		}
	}

	return nil
}

func (r *RequestsMetric) GetMetricList() (*RequestMetricResult, error) {

	metricCounters := &RequestMetricResult{
		TotalRequests:     r.TotalRequests,
		TotalErrors:       r.TotalErrors,
		TotalServerErrors: r.TotalServerErrors,
	}

	return metricCounters, nil
}
