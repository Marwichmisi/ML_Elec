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
