package metric

import (
	"log/slog"
)

type Manager struct {
	log           *slog.Logger
	RequestMetric *RequestsMetric
	//db metrics
	//...
}

func New(log *slog.Logger) *Manager {
	return &Manager{
		log:           log,
		RequestMetric: NewRequestCounter(log),
	}
}
