package http

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	totalRequests     int
	totalErrors       int
	totalServerErrors int
	log               *slog.Logger
	promTotalRequests prometheus.Counter
	promTotalErrors   prometheus.Counter
	promServerErrors  prometheus.Counter
}

func NewRequestCounter(log *slog.Logger) *RequestsMetric {

	promTotalRequests := promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_requests_count",
		Help: "The total number of requests",
	})
	promTotalErrors := promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_request_errors_count",
		Help: "The total number of error requests",
	})
	promServerErrors := promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_server_errors_count",
		Help: "The total number of fatal error requests",
	})

	return &RequestsMetric{
		log:               log,
		promTotalRequests: promTotalRequests,
		promTotalErrors:   promTotalErrors,
		promServerErrors:  promServerErrors,
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
		TotalRequests:     r.totalRequests,
		TotalErrors:       r.totalErrors,
		TotalServerErrors: r.totalServerErrors,
	}

	return metricCounters, nil
}

func (r *RequestsMetric) IncTotalRequests() {
	r.totalRequests++
	r.promTotalRequests.Inc()
}

func (r *RequestsMetric) IncTotalErrors() {
	r.totalErrors++
	r.promTotalErrors.Inc()
}

func (r *RequestsMetric) IncServerErrors() {
	r.totalServerErrors++
	r.promServerErrors.Inc()
}
