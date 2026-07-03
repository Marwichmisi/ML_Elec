// Package v1 provides the versioned gRPC plugin SDK for ML-Elec.
// It defines the plugin lifecycle and sensor collection contracts
// that all plugins must implement.
package v1

import (
	"context"
	"fmt"

	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// PluginLifecycle defines the interface for plugin lifecycle management.
// Plugins implement this to handle initialization, startup, and shutdown.
type PluginLifecycle interface {
	// Init initializes the plugin with the provided configuration.
	Init(ctx context.Context, config map[string]string) error

	// Start starts the plugin. Called after successful Init.
	Start(ctx context.Context) error

	// Stop gracefully stops the plugin.
	Stop(ctx context.Context) error
}

// SensorCollector defines the interface for collecting sensor data.
// Plugins implement this to send sensor readings to the core.
type SensorCollector interface {
	// Collect sends a sensor reading to the core for processing.
	Collect(ctx context.Context, req *CollectRequest) (*CollectResponse, error)
}

// GRPCPlugin is the go-plugin implementation that bridges the gRPC transport
// with the PluginLifecycle and SensorCollector interfaces.
type GRPCPlugin struct {
	goplugin.NetRPCUnsupportedPlugin

	// Impl is the plugin lifecycle implementation (server-side).
	Impl PluginLifecycle

	// CollectImpl is the sensor collector implementation (server-side).
	CollectImpl SensorCollector
}

// GRPCServer registers both PluginLifecycle and SensorCollector services on the gRPC server.
func (p *GRPCPlugin) GRPCServer(_ *goplugin.GRPCBroker, s *grpc.Server) error {
	RegisterPluginLifecycleServer(s, &grpcServer{
		impl:        p.Impl,
		collectImpl: p.CollectImpl,
	})
	RegisterSensorCollectorServer(s, &grpcServer{
		impl:        p.Impl,
		collectImpl: p.CollectImpl,
	})
	return nil
}

// GRPCClient returns a client that satisfies both PluginLifecycle and SensorCollector interfaces.
func (p *GRPCPlugin) GRPCClient(_ context.Context, _ *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{
		lifecycle: NewPluginLifecycleClient(c),
		collector: NewSensorCollectorClient(c),
	}, nil
}

// grpcServer is the gRPC server-side adapter that delegates to the actual implementations.
type grpcServer struct {
	UnimplementedPluginLifecycleServer
	UnimplementedSensorCollectorServer

	impl        PluginLifecycle
	collectImpl SensorCollector
}

// Init delegates to the PluginLifecycle implementation.
func (s *grpcServer) Init(ctx context.Context, req *InitRequest) (*InitResponse, error) {
	if err := s.impl.Init(ctx, req.Config); err != nil {
		return nil, fmt.Errorf("plugin init: %w", err)
	}
	return &InitResponse{}, nil
}

// Start delegates to the PluginLifecycle implementation.
func (s *grpcServer) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	if err := s.impl.Start(ctx); err != nil {
		return nil, fmt.Errorf("plugin start: %w", err)
	}
	return &StartResponse{}, nil
}

// Stop delegates to the PluginLifecycle implementation.
func (s *grpcServer) Stop(ctx context.Context, req *StopRequest) (*StopResponse, error) {
	if err := s.impl.Stop(ctx); err != nil {
		return nil, fmt.Errorf("plugin stop: %w", err)
	}
	return &StopResponse{}, nil
}

// Collect delegates to the SensorCollector implementation.
func (s *grpcServer) Collect(ctx context.Context, req *CollectRequest) (*CollectResponse, error) {
	resp, err := s.collectImpl.Collect(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("sensor collect: %w", err)
	}
	return resp, nil
}

// grpcClient is the gRPC client-side adapter that delegates to the gRPC clients.
type grpcClient struct {
	lifecycle PluginLifecycleClient
	collector SensorCollectorClient
}

// Init calls the Init RPC on the remote plugin.
func (c *grpcClient) Init(ctx context.Context, config map[string]string) error {
	_, err := c.lifecycle.Init(ctx, &InitRequest{Config: config})
	if err != nil {
		return fmt.Errorf("rpc init: %w", err)
	}
	return nil
}

// Start calls the Start RPC on the remote plugin.
func (c *grpcClient) Start(ctx context.Context) error {
	_, err := c.lifecycle.Start(ctx, &StartRequest{})
	if err != nil {
		return fmt.Errorf("rpc start: %w", err)
	}
	return nil
}

// Stop calls the Stop RPC on the remote plugin.
func (c *grpcClient) Stop(ctx context.Context) error {
	_, err := c.lifecycle.Stop(ctx, &StopRequest{})
	if err != nil {
		return fmt.Errorf("rpc stop: %w", err)
	}
	return nil
}

// Collect calls the Collect RPC on the remote plugin.
func (c *grpcClient) Collect(ctx context.Context, req *CollectRequest) (*CollectResponse, error) {
	resp, err := c.collector.Collect(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("rpc collect: %w", err)
	}
	return resp, nil
}

// Compile-time interface checks.
var _ goplugin.GRPCPlugin = (*GRPCPlugin)(nil)
var _ PluginLifecycle = (*grpcClient)(nil)
var _ SensorCollector = (*grpcClient)(nil)
