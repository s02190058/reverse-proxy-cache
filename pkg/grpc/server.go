package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

var (
	ErrBadStart = errors.New("cannot start a gRPC server")
)

// Server is a wrapper over grpc.Server.
type Server struct {
	grpcServer *grpc.Server
	notify     chan error
}

// NewServer creates a gRPC server wrapper.
func NewServer(l *slog.Logger) *Server {
	server := &Server{
		grpcServer: grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				recoveryUnaryServerInterceptor(l),
				logging.UnaryServerInterceptor(
					interceptorLogger(l),
					logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
				),
			),
		),
	}

	return server
}

// RegisterHandlers allows to add handlers to the gRPC server that is encapsulated.
func (s *Server) RegisterHandlers(register func(s *grpc.Server)) {
	register(s.grpcServer)
}

// Start launches a gRPC server in a separate goroutine.
func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", addr(port))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrBadStart, err)
	}
	go func() {
		s.notify <- s.grpcServer.Serve(lis)
	}()

	return nil
}

// Notify is a wrapper function that converts chan error to <-chan error.
// We use it to prevent the caller from writing to the channel.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Stop shut down a gRPC server gracefully.
func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
}

// interceptorLogger adapts slog logger to interceptor logger.
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// addr returns a :port gRPC server address.
func addr(port string) string {
	return ":" + port
}
