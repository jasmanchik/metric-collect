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
	TotalRequests       int
	TotalErrors         int
	TotalServerErrors   int
	log                 *slog.Logger
	PromREDMetric       *prometheus.CounterVec
	PromTotalRequests   *prometheus.CounterVec
	PromTotalErrors     *prometheus.CounterVec
	PromRequestDuration *prometheus.HistogramVec
}

func NewRequestCounter(log *slog.Logger) *RequestsMetric {
	promTotalErrors := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "The total number of error requests",
	}, []string{"method", "endpoint", "status"})

	promTotalRequests := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of requests",
	}, []string{"method", "endpoint"})

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	err := prometheus.Register(requestDuration)
	if err != nil {
		log.Error("failed to register request duration metric", slog.String("error", err.Error()))
	}

	return &RequestsMetric{
		log:                 log,
		PromTotalRequests:   promTotalRequests,
		PromTotalErrors:     promTotalErrors,
		PromRequestDuration: requestDuration,
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
