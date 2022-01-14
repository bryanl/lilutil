package grpcutil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/bryanl/lilutil/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

var (
	// defaultGRPCDialOptions is the default set of GRPC dial options.
	defaultGRPCDialOptions = []grpc.DialOption{
		// defaulting to insecure connections because this is a POC.
		grpc.WithInsecure(),
	}

	// defaultGatewayShutdownTimeout is the default gateway shutdown timeout.
	defaultGatewayShutdownTimeout = time.Second * 3
)

// Endpoint is an endpoint function.
type Endpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// GatewayConfig is configuration Gateway.
type GatewayConfig struct {
	// ServerAddr is the GRPC address.
	ServerAddr string
	// HTTPAddr is the address the HTTP server listens on.
	HTTPAddr string
	// Endpoints is a slice of endpoint functions.
	Endpoints []Endpoint
}

// Gateway is a GRPC HTTP gateway.
type Gateway struct {
	name            string
	config          GatewayConfig
	grpcDialOptions []grpc.DialOption
	shutdownTimeout time.Duration
}

// NewGateway creates an instance of Gateway.
func NewGateway(name string, config GatewayConfig) *Gateway {
	g := &Gateway{
		name:            name,
		config:          config,
		grpcDialOptions: defaultGRPCDialOptions,
		shutdownTimeout: defaultGatewayShutdownTimeout,
	}

	return g
}

// Start starts the gateway. It returns a channel that, when closed, stops the gateway.
func (g *Gateway) Start(ctx context.Context) (<-chan struct{}, error) {
	logger := log.From(ctx).WithName(g.name)

	mux := runtime.NewServeMux()

	for _, endpoint := range g.config.Endpoints {
		if err := g.registerEndpoint(ctx, mux, endpoint); err != nil {
			return nil, fmt.Errorf("register endpoint: %w", err)
		}
	}

	ch := make(chan struct{}, 1)

	handler := cors.AllowAll().Handler(mux)

	httpServer := &http.Server{
		Addr:        g.config.HTTPAddr,
		Handler:     handler,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	go func() {
		logger.Info("start up",
			"addr", g.config.HTTPAddr,
			"server", g.config.ServerAddr)

		err := httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error(err, "failed to stop HTTP server cleanly")
		}
		logger.Info("gateway has stopped")
	}()

	go func() {
		<-ctx.Done()
		logger.Info("stopping gracefully")

		gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), g.shutdownTimeout)
		defer cancelShutdown()

		if err := httpServer.Shutdown(gracefulCtx); err != nil {
			logger.Error(err, "attempting to shut server down")
		}

		close(ch)
	}()

	return ch, nil
}

func (g *Gateway) registerEndpoint(ctx context.Context, mux *runtime.ServeMux, fn Endpoint) error {
	return fn(ctx, mux, g.config.ServerAddr, g.grpcDialOptions)
}
