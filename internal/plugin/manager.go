package plugin

import (
	"fmt"
	"net/rpc"
	"sync"

	goplugin "github.com/hashicorp/go-plugin"
)

// SensorPlugin is the interface plugins must implement.
type SensorPlugin interface {
	Echo(msg string) (string, error)
}

// HandshakeConfig is the shared handshake configuration for all ML-Elec plugins.
var HandshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ML_ELEC_PLUGIN",
	MagicCookieValue: "ml-elec-v1",
}

// SensorPluginRPC is the go-plugin Plugin implementation for net/rpc.
type SensorPluginRPC struct{}

// Server returns the RPC server for the plugin.
func (p *SensorPluginRPC) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &SensorPluginRPCServer{}, nil
}

// Client returns an RPC client that implements SensorPlugin.
func (p *SensorPluginRPC) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &SensorPluginRPCClient{client: c}, nil
}

// SensorPluginRPCServer is the RPC server side.
type SensorPluginRPCServer struct {
	impl SensorPlugin
}

// EchoArgs is the RPC argument for Echo.
type EchoArgs struct {
	Msg string
}

// EchoReply is the RPC reply for Echo.
type EchoReply struct {
	Msg string
}

// Echo handles the Echo RPC call.
func (s *SensorPluginRPCServer) Echo(args *EchoArgs, reply *EchoReply) error {
	msg, err := s.impl.Echo(args.Msg)
	if err != nil {
		return err
	}
	reply.Msg = msg
	return nil
}

// SensorPluginRPCClient is the RPC client side.
type SensorPluginRPCClient struct {
	client *rpc.Client
}

// Echo calls the Echo RPC on the plugin.
func (c *SensorPluginRPCClient) Echo(msg string) (string, error) {
	var reply EchoReply
	err := c.client.Call("Plugin.Echo", &EchoArgs{Msg: msg}, &reply)
	if err != nil {
		return "", fmt.Errorf("rpc echo: %w", err)
	}
	return reply.Msg, nil
}

// Manager handles plugin lifecycle with child process isolation.
type Manager struct {
	clients map[string]*goplugin.Client
	mu      sync.RWMutex
}

// NewManager creates a new plugin manager with an empty client map.
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*goplugin.Client),
	}
}

// Launch starts a plugin as a child process if it is in the enabled list.
// STUB: returns error until GREEN phase implementation.
func (m *Manager) Launch(name, path string, enabledPlugins []string) error {
	return fmt.Errorf("not implemented")
}

// Kill stops a running plugin and removes it from the manager.
// STUB: returns error until GREEN phase implementation.
func (m *Manager) Kill(name string) error {
	return fmt.Errorf("not implemented")
}

// ShutdownAll stops all running plugins.
// STUB: no-op until GREEN phase implementation.
func (m *Manager) ShutdownAll() {
}

// IsRunning checks if a plugin is currently running.
// STUB: returns false until GREEN phase implementation.
func (m *Manager) IsRunning(name string) bool {
	return false
}
