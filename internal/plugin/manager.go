package plugin

import (
	"fmt"
	"net/rpc"
	"os/exec"
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
// On the server side (plugin binary), Impl must be set to the actual implementation.
// On the client side (manager), Impl is not used.
type SensorPluginRPC struct {
	Impl SensorPlugin
}

// Server returns the RPC server for the plugin.
func (p *SensorPluginRPC) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &SensorPluginRPCServer{impl: p.Impl}, nil
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

// LaunchGRPC starts a plugin with gRPC transport (net/rpc disabled).
// The grpcPlugin parameter should be a *sdk.GRPCPlugin or equivalent go-plugin.GRPCPlugin.
func (m *Manager) LaunchGRPC(name, path string, enabledPlugins []string, grpcPlugin goplugin.Plugin) error {
	if !isPluginEnabled(name, enabledPlugins) {
		return fmt.Errorf("plugin %q is not enabled", name)
	}

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig:  HandshakeConfig,
		Plugins:          map[string]goplugin.Plugin{"sensor": grpcPlugin},
		Cmd:              exec.Command(path),
		Managed:          true,
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
	})

	// Verify the plugin starts and connects
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("connecting to grpc plugin %q: %w", name, err)
	}

	// Dispense the plugin to verify it works
	_, err = rpcClient.Dispense("sensor")
	if err != nil {
		client.Kill()
		return fmt.Errorf("dispensing grpc plugin %q: %w", name, err)
	}

	m.mu.Lock()
	m.clients[name] = client
	m.mu.Unlock()

	return nil
}

// Launch starts a plugin as a child process if it is in the enabled list.
func (m *Manager) Launch(name, path string, enabledPlugins []string) error {
	if !isPluginEnabled(name, enabledPlugins) {
		return fmt.Errorf("plugin %q is not enabled", name)
	}

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"sensor": &SensorPluginRPC{},
		},
		Cmd:     exec.Command(path),
		Managed: true,
	})

	// Verify the plugin starts and connects
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("connecting to plugin %q: %w", name, err)
	}

	// Dispense the plugin to verify it works
	_, err = rpcClient.Dispense("sensor")
	if err != nil {
		client.Kill()
		return fmt.Errorf("dispensing plugin %q: %w", name, err)
	}

	m.mu.Lock()
	m.clients[name] = client
	m.mu.Unlock()

	return nil
}

// Kill stops a running plugin and removes it from the manager.
func (m *Manager) Kill(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, ok := m.clients[name]
	if !ok {
		return fmt.Errorf("plugin %q not found", name)
	}

	client.Kill()
	delete(m.clients, name)
	return nil
}

// ShutdownAll stops all running plugins.
func (m *Manager) ShutdownAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, client := range m.clients {
		client.Kill()
		delete(m.clients, name)
	}
}

// IsRunning checks if a plugin is currently running.
func (m *Manager) IsRunning(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.clients[name]
	return ok
}

// isPluginEnabled checks if a plugin name is in the enabled list.
func isPluginEnabled(name string, enabled []string) bool {
	for _, e := range enabled {
		if e == name {
			return true
		}
	}
	return false
}
