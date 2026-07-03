package main

import (
	"context"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-hclog"

	"ml-elec/internal/plugin"
	sdk "ml-elec/pkg/sdk/v1"
)

// MockSensorPlugin implements both PluginLifecycle and SensorCollector.
type MockSensorPlugin struct{}

// Init initializes the mock plugin (no-op).
func (p *MockSensorPlugin) Init(_ context.Context, _ map[string]string) error { return nil }

// Start starts the mock plugin (no-op).
func (p *MockSensorPlugin) Start(_ context.Context) error { return nil }

// Stop stops the mock plugin (no-op).
func (p *MockSensorPlugin) Stop(_ context.Context) error { return nil }

// Collect handles a sensor collection request (always accepts).
func (p *MockSensorPlugin) Collect(_ context.Context, _ *sdk.CollectRequest) (*sdk.CollectResponse, error) {
	return &sdk.CollectResponse{Accepted: true}, nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "mock-plugin",
		Level:  hclog.Info,
		Output: hclog.DefaultOutput,
	})

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"sensor": &sdk.GRPCPlugin{
				Impl:        &MockSensorPlugin{},
				CollectImpl: &MockSensorPlugin{},
			},
		},
		Logger: logger,
	})
}
