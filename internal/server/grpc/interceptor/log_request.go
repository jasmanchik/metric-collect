package interceptor

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"time"
)

func LogRequests(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime)
		logger.With("interceptor", "logRequest").Info(fmt.Sprintf("Called %s, duration %v.", info.FullMethod, duration))

		return resp, err
	}
}
