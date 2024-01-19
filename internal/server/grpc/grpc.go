package grpc

import (
	"http-metric/internal/service/metric"
	"log/slog"
)

type ServerGrpc struct {
	log     *slog.Logger
	port    int
	metrics *metric.Manager
}

func (s *ServerGrpc) Start() error {
	//todo start grpc
	return nil
}

func (s *ServerGrpc) Stop() error {
	//todo stop grpc
	return nil
}

func New(logger *slog.Logger, port int, metric *metric.Manager) ServerGrpc {
	return ServerGrpc{
		log:     logger,
		port:    port,
		metrics: metric,
	}
}
