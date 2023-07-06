package grpc

import (
	"context"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

// recoveryUnaryServerInterceptor catches the panic.
func recoveryUnaryServerInterceptor(l *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		defer func() {
			if r := recover(); r != nil {
				l.Error("panic occurred", slog.Attr{
					Key:   "error",
					Value: slog.AnyValue(r),
				})
			}
		}()

		return handler(ctx, req)
	}
}
