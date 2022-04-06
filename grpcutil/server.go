package grpcutil

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/multierr"
	"google.golang.org/grpc"

	"github.com/bryanl/lilutil/log"
)

// RegisterFn is a function that takes a GRPC server as input.
type RegisterFn func(s *grpc.Server) error

// ServerConfig is configuration for GRPCServer.
type ServerConfig struct {
	// Listener is where the server will listen.
	Listener net.Listener

	// RegisterFunc is a function that allows you to register servers.
	RegisterFunc RegisterFn
}

// Validate validates the ServerConfig. If not valid, an error is returned.
func (config *ServerConfig) Validate() error {
	var err error

	if config.Listener == nil {
		err = multierr.Append(err, errors.New("listener is required"))
	}

	if config.RegisterFunc == nil {
		err = multierr.Append(err, errors.New("register function is required"))
	}

	return err
}

// Server provides a GRPC server
type Server struct {
	config ServerConfig
	name   string
}

// NewServer creates an instance of Server.
func NewServer(name string, config ServerConfig) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	server := &Server{
		name:   name,
		config: config,
	}

	return server, nil
}

// Start starts the GRPC server.
func (server *Server) Start(ctx context.Context) (<-chan struct{}, error) {
	logger := log.From(ctx).WithName(server.name)

	s := grpc.NewServer()

	if err := server.config.RegisterFunc(s); err != nil {
		return nil, fmt.Errorf("grpc register: %w", err)
	}

	lis := server.config.Listener
	logger.Info("Starting GRPC server", "addr", lis.Addr().String())

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error(err, "Unable to stop GRPC server cleanly")
		}
		logger.Info("GRPC server has stopped")
	}()

	ch := make(chan struct{}, 1)

	go func() {
		<-ctx.Done()
		logger.Info("Stopping GRPC server gracefully")
		s.GracefulStop()
		close(ch)
	}()

	return ch, nil
}
