package metric

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"math/rand"
	"sync/atomic"
	"time"
)

type RequestsMetric struct {
	log *slog.Logger

	totalRequests       int32
	totalErrors         int32
	totalServerErrors   int32
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

var (
	BadRequestError = errors.New("bad request")
)

type RequestMetricResult struct {
	TotalRequests     int32 `json:"total_requests_count"`
	TotalErrors       int32 `json:"total_request_errors_count"`
	TotalServerErrors int32 `json:"total_server_errors_count"`
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
		return BadRequestError
	}

	return nil
}

func (r *RequestsMetric) GetMetricList() *RequestMetricResult {

	metricCounters := &RequestMetricResult{
		TotalRequests:     r.totalRequests,
		TotalErrors:       r.totalErrors,
		TotalServerErrors: r.totalServerErrors,
	}

	return metricCounters
}

func (r *RequestsMetric) AddTotalRequests(i int32) {
	atomic.AddInt32(&r.totalRequests, i)
}

func (r *RequestsMetric) AddTotalErrors(i int32) {
	atomic.AddInt32(&r.totalErrors, i)
}

func (r *RequestsMetric) AddServerErrors(i int32) {
	atomic.AddInt32(&r.totalServerErrors, i)
}
