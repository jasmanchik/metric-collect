package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"http-metric/internal/service/metric"
	"log/slog"
	"net/http"
)

type HttpMetricHandler struct {
	log    *slog.Logger
	metric *metric.Manager
}

func NewHttpMetricHandler(logger *slog.Logger, m *metric.Manager) *HttpMetricHandler {
	return &HttpMetricHandler{
		log:    logger,
		metric: m,
	}
}

func (m *HttpMetricHandler) Ping(c *gin.Context) {
	err := m.metric.RequestMetric.Ping()
	if err != nil {
		if errors.Is(err, metric.BadRequestError) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"details": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "pong",
	})
}

func (m *HttpMetricHandler) RequestCounter(c *gin.Context) {
	metricResult := m.metric.RequestMetric.GetMetricList()

	c.JSON(http.StatusOK, metricResult)
}
