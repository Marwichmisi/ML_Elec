package main

import (
	"fmt"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-hclog"

	"ml-elec/internal/plugin"
)

// MockSensorPlugin is a mock implementation of the SensorPlugin interface.
type MockSensorPlugin struct{}

// Echo returns the message back.
func (p *MockSensorPlugin) Echo(msg string) (string, error) {
	return msg, nil
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
			"sensor": &plugin.SensorPluginRPC{},
		},
		Logger: logger,
	})
	_ = fmt.Sprintf // keep import
}
