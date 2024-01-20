package handlers

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metricsv1 "http-metric/api/grpc/go"
	"http-metric/internal/service/metric"
	"log/slog"
)

type GrpcMetricHandler struct {
	log    *slog.Logger
	metric *metric.Manager
	metricsv1.UnimplementedMetricServer
}

func NewGrpcMetricHandler(logger *slog.Logger, m *metric.Manager) *GrpcMetricHandler {
	return &GrpcMetricHandler{
		log:    logger,
		metric: m,
	}
}

func (m *GrpcMetricHandler) Ping(context.Context, *metricsv1.PingRequest) (*metricsv1.PingResponse, error) {

	err := m.metric.RequestMetric.Ping()
	if err != nil {
		if errors.Is(err, metric.BadRequestError) {
			return &metricsv1.PingResponse{}, status.Error(codes.InvalidArgument, "bad request error")
		}
	}

	return &metricsv1.PingResponse{}, nil
}

func (m *GrpcMetricHandler) RequestsCounter(context.Context, *metricsv1.RequestsCounterRequest) (*metricsv1.RequestsCounterResponse, error) {

	metricResult := m.metric.RequestMetric.GetMetricList()

	return &metricsv1.RequestsCounterResponse{
		TotalRequestsCount:      metricResult.TotalRequests,
		TotalRequestErrorsCount: metricResult.TotalErrors,
		TotalServerErrorsCount:  metricResult.TotalServerErrors,
	}, nil
}
