package metrics

import (
	httpMetric "http-metric/internal/metrics/http"
	"log/slog"
)

type MetricRecorder struct {
	log        *slog.Logger
	HttpMetric *httpMetric.RequestsMetric
	//db metrics
	//...
}

func NewMetric(log *slog.Logger) *MetricRecorder {
	return &MetricRecorder{
		log:        log,
		HttpMetric: httpMetric.NewRequestCounter(log),
	}
}
