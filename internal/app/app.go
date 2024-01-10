package app

import (
	httpapp "http-metric/internal/app/http"
	"http-metric/internal/metrics"
	"log/slog"
	"time"
)

type App struct {
	HTTPServer *httpapp.App
}

func New(logger *slog.Logger, port int, timeout time.Duration) *App {
	metric := metrics.NewMetric(logger)
	server := httpapp.New(logger, port, timeout, metric)

	return &App{
		HTTPServer: server,
	}
}
