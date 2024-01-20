package grpc

import (
	"fmt"
	grpclib "google.golang.org/grpc"
	metricsv1 "http-metric/api/grpc/go"
	"http-metric/internal/server/grpc/handlers"
	"http-metric/internal/server/grpc/interceptor"
	"http-metric/internal/service/metric"
	"log/slog"
	"net"
)

type ServerGrpc struct {
	log     *slog.Logger
	port    int
	metrics *metric.Manager
	server  *grpclib.Server
}

func (s *ServerGrpc) Start() error {
	const op = "grpc.Start"
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	s.log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))

	if err := s.server.Serve(l); err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	return nil
}

func (s *ServerGrpc) Stop() error {
	const op = "grpc.Stop"

	s.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", s.port))

	s.server.GracefulStop()

	return nil
}

func New(logger *slog.Logger, port int, metric *metric.Manager) *ServerGrpc {

	gRPCServer := grpclib.NewServer(grpclib.ChainUnaryInterceptor(
		interceptor.Recovery(logger, metric),
		interceptor.LogRequests(logger),
		interceptor.Metric(metric),
	))

	metricsv1.RegisterMetricServer(gRPCServer, handlers.NewGrpcMetricHandler(logger, metric))

	return &ServerGrpc{
		log:     logger,
		port:    port,
		metrics: metric,
		server:  gRPCServer,
	}
}
