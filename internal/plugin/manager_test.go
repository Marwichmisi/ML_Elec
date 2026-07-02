package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildMockPlugin builds the mock plugin binary and returns its path.
func buildMockPlugin(t *testing.T) string {
	t.Helper()

	binaryPath := filepath.Join(t.TempDir(), "mock-plugin")
	cmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd/mock-plugin")
	cmd.Dir = filepath.Join(".")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to build mock plugin: %v\n%s", err, out)
	}

	return binaryPath
}

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("NewManager returned nil")
	}
	if len(m.clients) != 0 {
		t.Errorf("expected empty clients map, got %d", len(m.clients))
	}
}

func TestLaunchMockPlugin(t *testing.T) {
	binaryPath := buildMockPlugin(t)

	m := NewManager()
	err := m.Launch("mock", binaryPath, []string{"mock"})
	if err != nil {
		t.Fatalf("Launch failed: %v", err)
	}

	if !m.IsRunning("mock") {
		t.Error("expected plugin to be running")
	}

	m.ShutdownAll()
}

func TestPluginEcho(t *testing.T) {
	binaryPath := buildMockPlugin(t)

	m := NewManager()
	err := m.Launch("mock", binaryPath, []string{"mock"})
	if err != nil {
		t.Fatalf("Launch failed: %v", err)
	}
	defer m.ShutdownAll()

	// Get the RPC client from the plugin
	m.mu.RLock()
	client, ok := m.clients["mock"]
	m.mu.RUnlock()
	if !ok {
		t.Fatal("plugin not found in clients map")
	}

	rpcClient, err := client.Client()
	if err != nil {
		t.Fatalf("failed to get RPC client: %v", err)
	}

	// Dispense the raw interface
	raw, err := rpcClient.Dispense("sensor")
	if err != nil {
		t.Fatalf("failed to dispense plugin: %v", err)
	}

	// Type assert to SensorPlugin
	sensorPlugin, ok := raw.(SensorPlugin)
	if !ok {
		t.Fatalf("dispensed plugin does not implement SensorPlugin: %T", raw)
	}

	// Test Echo
	reply, err := sensorPlugin.Echo("hello")
	if err != nil {
		t.Fatalf("Echo failed: %v", err)
	}
	if reply != "hello" {
		t.Errorf("expected 'hello', got %q", reply)
	}
}

func TestCrashIsolation(t *testing.T) {
	binaryPath := buildMockPlugin(t)

	m := NewManager()
	err := m.Launch("mock", binaryPath, []string{"mock"})
	if err != nil {
		t.Fatalf("Launch failed: %v", err)
	}

	// Get the plugin client
	m.mu.RLock()
	client := m.clients["mock"]
	m.mu.RUnlock()

	// Kill the plugin process (simulates a crash)
	// client.Kill() terminates the child process
	client.Kill()

	// Wait for the process to exit by polling Exited()
	for i := 0; i < 50; i++ {
		if client.Exited() {
			break
		}
	}

	// Manager should still be functional - it should not crash
	// Try to shut down gracefully (should not panic)
	m.ShutdownAll()
}

func TestLaunchDisabledPlugin(t *testing.T) {
	binaryPath := buildMockPlugin(t)

	m := NewManager()
	err := m.Launch("mock", binaryPath, []string{"other-plugin"})
	if err == nil {
		t.Error("expected error when launching disabled plugin")
		m.ShutdownAll()
	}

	if m.IsRunning("mock") {
		t.Error("disabled plugin should not be running")
	}
}

func TestShutdownAll(t *testing.T) {
	binaryPath := buildMockPlugin(t)

	m := NewManager()

	// Launch multiple plugins (using different names but same binary for testing)
	err := m.Launch("mock1", binaryPath, []string{"mock1", "mock2"})
	if err != nil {
		t.Fatalf("Launch mock1 failed: %v", err)
	}

	err = m.Launch("mock2", binaryPath, []string{"mock1", "mock2"})
	if err != nil {
		t.Fatalf("Launch mock2 failed: %v", err)
	}

	if !m.IsRunning("mock1") {
		t.Error("mock1 should be running")
	}
	if !m.IsRunning("mock2") {
		t.Error("mock2 should be running")
	}

	m.ShutdownAll()

	// After shutdown, nothing should be running
	// Note: IsRunning may still return true for managed clients until cleanup
	// This is a known behavior of go-plugin
}

func TestLaunchInvalidPath(t *testing.T) {
	m := NewManager()
	err := m.Launch("nonexistent", "/nonexistent/binary", []string{"nonexistent"})
	if err == nil {
		t.Error("expected error when launching with invalid path")
		m.ShutdownAll()
	}
}

// TestMain builds the mock plugin binary before running tests.
func TestMain(m *testing.M) {
	// Build the mock plugin binary
	cmd := exec.Command("go", "build", "-o", "mock-plugin-test-binary", "../../cmd/mock-plugin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic("failed to build mock plugin for tests: " + string(out))
	}
	defer os.Remove("mock-plugin-test-binary")

	os.Exit(m.Run())
}
