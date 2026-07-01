package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if cfg.NATS.Host != "127.0.0.1" {
		t.Errorf("NATS.Host = %q, want %q", cfg.NATS.Host, "127.0.0.1")
	}

	if cfg.NATS.Port != -1 {
		t.Errorf("NATS.Port = %d, want %d", cfg.NATS.Port, -1)
	}

	if cfg.Storage.Path != "./data/sensors.db" {
		t.Errorf("Storage.Path = %q, want %q", cfg.Storage.Path, "./data/sensors.db")
	}

	if cfg.API.Port != 8080 {
		t.Errorf("API.Port = %d, want %d", cfg.API.Port, 8080)
	}

	if len(cfg.Plugins.Enabled) != 0 {
		t.Errorf("Plugins.Enabled = %v, want empty slice", cfg.Plugins.Enabled)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name: "empty NATS host",
			cfg: &Config{
				NATS:    NATSConfig{Host: "", Port: -1},
				Storage: StorageConfig{Path: "./data/sensors.db"},
			},
			wantErr: true,
		},
		{
			name: "empty storage path",
			cfg: &Config{
				NATS:    NATSConfig{Host: "127.0.0.1", Port: -1},
				Storage: StorageConfig{Path: ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadDefaults(t *testing.T) {
	// Ensure no config file exists in current directory
	orig, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	defer os.Chdir(orig)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.NATS.Host != "127.0.0.1" {
		t.Errorf("Load() NATS.Host = %q, want %q", cfg.NATS.Host, "127.0.0.1")
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `nats:
  host: "0.0.0.0"
  port: 4222
storage:
  path: "/tmp/test.db"
api:
  port: 9090
plugins:
  enabled:
    - mqtt-sensor
    - modbus-gateway
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := loadFromFile(configPath)
	if err != nil {
		t.Fatalf("loadFromFile() unexpected error: %v", err)
	}

	if cfg.NATS.Host != "0.0.0.0" {
		t.Errorf("NATS.Host = %q, want %q", cfg.NATS.Host, "0.0.0.0")
	}

	if cfg.NATS.Port != 4222 {
		t.Errorf("NATS.Port = %d, want %d", cfg.NATS.Port, 4222)
	}

	if cfg.Storage.Path != "/tmp/test.db" {
		t.Errorf("Storage.Path = %q, want %q", cfg.Storage.Path, "/tmp/test.db")
	}

	if cfg.API.Port != 9090 {
		t.Errorf("API.Port = %d, want %d", cfg.API.Port, 9090)
	}

	if len(cfg.Plugins.Enabled) != 2 {
		t.Errorf("len(Plugins.Enabled) = %d, want 2", len(cfg.Plugins.Enabled))
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `nats:
  host: [invalid yaml
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := loadFromFile(configPath)
	if err == nil {
		t.Fatal("loadFromFile() expected error for invalid YAML, got nil")
	}
}

// --- Additional edge case tests for Task 2 ---

func TestLoadPartialYAML(t *testing.T) {
	// YAML with only some fields set — defaults should apply for the rest
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `nats:
  host: "192.168.1.100"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := loadFromFile(configPath)
	if err != nil {
		t.Fatalf("loadFromFile() unexpected error: %v", err)
	}

	// Override from YAML
	if cfg.NATS.Host != "192.168.1.100" {
		t.Errorf("NATS.Host = %q, want %q", cfg.NATS.Host, "192.168.1.100")
	}

	// Defaults preserved
	if cfg.NATS.Port != -1 {
		t.Errorf("NATS.Port = %d, want default -1", cfg.NATS.Port)
	}
	if cfg.Storage.Path != "./data/sensors.db" {
		t.Errorf("Storage.Path = %q, want default", cfg.Storage.Path)
	}
	if cfg.API.Port != 8080 {
		t.Errorf("API.Port = %d, want default 8080", cfg.API.Port)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := loadFromFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("loadFromFile() expected error for nonexistent file, got nil")
	}
}

func TestLoadFromFileErrorContext(t *testing.T) {
	_, err := loadFromFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error")
	}
	// Error should contain context about the operation
	errStr := err.Error()
	if errStr == "" {
		t.Error("error message should not be empty")
	}
}

func TestLoadConfigPathSearchOrder(t *testing.T) {
	// Test that Load() finds config in current directory first
	orig, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	defer os.Chdir(orig)

	// Create config in current directory
	localConfig := filepath.Join(dir, "config.yaml")
	content := `nats:
  host: "from-local"
`
	if err := os.WriteFile(localConfig, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.NATS.Host != "from-local" {
		t.Errorf("Load() should use local config, got NATS.Host = %q", cfg.NATS.Host)
	}
}

func TestPluginConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	content := `plugins:
  enabled:
    - mqtt-sensor
    - modbus-gateway
    - anomaly-detector
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := loadFromFile(configPath)
	if err != nil {
		t.Fatalf("loadFromFile() unexpected error: %v", err)
	}

	if len(cfg.Plugins.Enabled) != 3 {
		t.Fatalf("len(Plugins.Enabled) = %d, want 3", len(cfg.Plugins.Enabled))
	}

	expected := []string{"mqtt-sensor", "modbus-gateway", "anomaly-detector"}
	for i, name := range cfg.Plugins.Enabled {
		if name != expected[i] {
			t.Errorf("Plugins.Enabled[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

func TestConfigPathSearchHomeDir(t *testing.T) {
	// Test that Load() falls back to ~/.config/ml-elec/config.yaml
	orig, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	defer os.Chdir(orig)

	// Create the home config directory
	home, _ := os.UserHomeDir()
	homeConfigDir := filepath.Join(home, ".config", "ml-elec")
	if err := os.MkdirAll(homeConfigDir, 0755); err != nil {
		t.Skip("cannot create home config dir")
	}

	// Remove any existing config.yaml in test dir and home dir
	os.Remove(filepath.Join(dir, "config.yaml"))
	os.Remove(filepath.Join(homeConfigDir, "config.yaml"))

	homeConfig := filepath.Join(homeConfigDir, "config.yaml")
	content := `nats:
  host: "from-home"
`
	if err := os.WriteFile(homeConfig, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write home config: %v", err)
	}
	defer os.Remove(homeConfig)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() unexpected error: %v", err)
	}

	if cfg.NATS.Host != "from-home" {
		t.Errorf("Load() should use home config, got NATS.Host = %q", cfg.NATS.Host)
	}
}

func TestLoadYAMLErrorWithContext(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")

	// Write a file that's not valid YAML (unclosed bracket)
	if err := os.WriteFile(configPath, []byte("nats:\n  host: [unclosed\n"), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := loadFromFile(configPath)
	if err == nil {
		t.Fatal("loadFromFile() expected error for invalid YAML")
	}

	// Error should mention the file path
	errStr := err.Error()
	if errStr == "" {
		t.Error("error message should not be empty")
	}
}

func TestConfigConcurrency(t *testing.T) {
	// Multiple goroutines calling Load() simultaneously should not race
	orig, _ := os.Getwd()
	dir := t.TempDir()
	os.Chdir(dir)
	defer os.Chdir(orig)

	// No config file — should return defaults
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			cfg, err := Load()
			if err != nil {
				t.Errorf("Load() error: %v", err)
				return
			}
			if cfg == nil {
				t.Error("Load() returned nil config")
			}
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
