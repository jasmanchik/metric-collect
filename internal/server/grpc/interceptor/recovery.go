package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"log/slog"
)

func Recovery(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if err := recover(); err != nil {
				logger.With("interceptor", "recovery").Error("panic", "error", err)
			}
		}()
		resp, err := handler(ctx, req)

		return resp, err
	}
}
