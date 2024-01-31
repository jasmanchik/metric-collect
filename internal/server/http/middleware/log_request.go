package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

func LogRequest(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)
		logger.With("interceptor", "logRequest").Info(fmt.Sprintf("Request %s %s handled %s. Status code: %d", c.Request.Method, c.Request.URL.Path, duration, c.Writer.Status()))
	}
}
