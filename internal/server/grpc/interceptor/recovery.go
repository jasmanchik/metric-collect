package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"http-metric/internal/service/metric"
	"log/slog"
	"strconv"
)

func Recovery(logger *slog.Logger, m *metric.Manager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if err := recover(); err != nil {
				logger.With("interceptor", "recovery").Error("panic", "error", err)
				if err != nil {
					method := "GRPC"
					path := info.FullMethod

					m.RequestMetric.AddServerErrors(1)
					m.RequestMetric.PromTotalErrors.WithLabelValues(method, path, strconv.Itoa(int(codes.Internal))).Inc()
					//как отдать InternalError статус?
				}
			}
		}()
		resp, err := handler(ctx, req)

		return resp, err
	}
}
