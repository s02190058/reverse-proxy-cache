package grpc

import (
	"errors"
	"fmt"
	"net"

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
func NewServer() *Server {
	server := &Server{
		grpcServer: grpc.NewServer(),
	}

	return server
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

// addr returns a :port gRPC server address.
func addr(port string) string {
	return ":" + port
}
