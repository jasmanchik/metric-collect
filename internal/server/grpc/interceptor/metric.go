package interceptor

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"http-metric/internal/service/metric"
)

func Metric(m *metric.Manager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		m.RequestMetric.AddTotalRequests(1)
		method := "GRPC"
		path := info.FullMethod
		m.RequestMetric.PromTotalRequests.WithLabelValues(method, path).Inc()
		timer := prometheus.NewTimer(m.RequestMetric.PromRequestDuration.WithLabelValues(method, path))
		resp, err := handler(ctx, req)
		timer.ObserveDuration()

		if err != nil {
			errStatus, _ := status.FromError(err)
			errCode := errStatus.Code()
			m.RequestMetric.AddTotalErrors(1)
			m.RequestMetric.PromTotalErrors.WithLabelValues(method, path, errCode.String()).Inc()
		}

		return resp, err
	}
}
