package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"http-metric/internal/service/metric"
	"net/http"
	"strconv"
)

func Metric(m *metric.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		m.RequestMetric.AddTotalRequests(1)
		method := c.Request.Method
		path := c.Request.URL.Path
		m.RequestMetric.PromTotalRequests.WithLabelValues(method, path).Inc()
		timer := prometheus.NewTimer(m.RequestMetric.PromRequestDuration.WithLabelValues(method, path))

		c.Next()

		if c.Writer.Status() >= http.StatusInternalServerError {
			m.RequestMetric.AddTotalErrors(1)
			m.RequestMetric.PromTotalErrors.WithLabelValues(method, path, strconv.Itoa(c.Writer.Status())).Inc()
		} else if c.Writer.Status() >= http.StatusBadRequest {
			m.RequestMetric.AddServerErrors(1)
		}

		timer.ObserveDuration()
	}
}
