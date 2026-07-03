package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name     string
		check    func(t *testing.T, cfg *Config)
	}{
		{
			name: "MQTT port default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.Port != 1883 {
					t.Errorf("MQTT.Port = %d, want 1883", cfg.MQTT.Port)
				}
			},
		},
		{
			name: "MQTT client_id default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.ClientID != "mqtt-plugin" {
					t.Errorf("MQTT.ClientID = %q, want %q", cfg.MQTT.ClientID, "mqtt-plugin")
				}
			},
		},
		{
			name: "MQTT clean_session default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.CleanSession != false {
					t.Errorf("MQTT.CleanSession = %v, want false", cfg.MQTT.CleanSession)
				}
			},
		},
		{
			name: "MQTT keep_alive default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.KeepAlive != 30 {
					t.Errorf("MQTT.KeepAlive = %d, want 30", cfg.MQTT.KeepAlive)
				}
			},
		},
		{
			name: "MQTT topics defaults",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.Topics.Subscribe != "esp32/#" {
					t.Errorf("MQTT.Topics.Subscribe = %q, want %q", cfg.MQTT.Topics.Subscribe, "esp32/#")
				}
				if cfg.MQTT.Topics.Status != "sys/mqtt-plugin/status" {
					t.Errorf("MQTT.Topics.Status = %q, want %q", cfg.MQTT.Topics.Status, "sys/mqtt-plugin/status")
				}
			},
		},
		{
			name: "Validation temperature range default",
			check: func(t *testing.T, cfg *Config) {
				temp, ok := cfg.Validation.Ranges["temperature"]
				if !ok {
					t.Error("Validation.Ranges[temperature] not found")
					return
				}
				if temp.Min != -40.0 {
					t.Errorf("temperature.Min = %f, want -40.0", temp.Min)
				}
				if temp.Max != 150.0 {
					t.Errorf("temperature.Max = %f, want 150.0", temp.Max)
				}
			},
		},
		{
			name: "Validation humidity range default",
			check: func(t *testing.T, cfg *Config) {
				humidity, ok := cfg.Validation.Ranges["humidity"]
				if !ok {
					t.Error("Validation.Ranges[humidity] not found")
					return
				}
				if humidity.Min != 0.0 {
					t.Errorf("humidity.Min = %f, want 0.0", humidity.Min)
				}
				if humidity.Max != 100.0 {
					t.Errorf("humidity.Max = %f, want 100.0", humidity.Max)
				}
			},
		},
		{
			name: "Validation current range default",
			check: func(t *testing.T, cfg *Config) {
				current, ok := cfg.Validation.Ranges["current"]
				if !ok {
					t.Error("Validation.Ranges[current] not found")
					return
				}
				if current.Min != 0.0 {
					t.Errorf("current.Min = %f, want 0.0", current.Min)
				}
				if current.Max != 1000.0 {
					t.Errorf("current.Max = %f, want 1000.0", current.Max)
				}
			},
		},
		{
			name: "Validation timestamp defaults",
			check: func(t *testing.T, cfg *Config) {
				if cfg.Validation.Timestamp.MaxFutureDrift != "5s" {
					t.Errorf("Timestamp.MaxFutureDrift = %q, want %q", cfg.Validation.Timestamp.MaxFutureDrift, "5s")
				}
				if cfg.Validation.Timestamp.MaxPastDrift != "24h" {
					t.Errorf("Timestamp.MaxPastDrift = %q, want %q", cfg.Validation.Timestamp.MaxPastDrift, "24h")
				}
			},
		},
		{
			name: "Validation health defaults",
			check: func(t *testing.T, cfg *Config) {
				if cfg.Validation.Health.MaxPayloadSize != 1024 {
					t.Errorf("Health.MaxPayloadSize = %d, want 1024", cfg.Validation.Health.MaxPayloadSize)
				}
				if cfg.Validation.Health.MinPayloadSize != 1 {
					t.Errorf("Health.MinPayloadSize = %d, want 1", cfg.Validation.Health.MinPayloadSize)
				}
			},
		},
		{
			name: "Assets auto_register default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.Assets.AutoRegister != true {
					t.Errorf("Assets.AutoRegister = %v, want true", cfg.Assets.AutoRegister)
				}
			},
		},
		{
			name: "Assets default_site default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.Assets.DefaultSite != "factory-1" {
					t.Errorf("Assets.DefaultSite = %q, want %q", cfg.Assets.DefaultSite, "factory-1")
				}
			},
		},
		{
			name: "MQTT QoS default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.QoS != 1 {
					t.Errorf("MQTT.QoS = %d, want 1", cfg.MQTT.QoS)
				}
			},
		},
		{
			name: "MQTT QoSPerTopic default empty",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.QoSPerTopic == nil {
					t.Error("MQTT.QoSPerTopic should not be nil")
				}
				if len(cfg.MQTT.QoSPerTopic) != 0 {
					t.Errorf("MQTT.QoSPerTopic should be empty, got %v", cfg.MQTT.QoSPerTopic)
				}
			},
		},
		{
			name: "MQTT NATSHost default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.NATSHost != "127.0.0.1" {
					t.Errorf("MQTT.NATSHost = %q, want %q", cfg.MQTT.NATSHost, "127.0.0.1")
				}
			},
		},
		{
			name: "MQTT NATSPort default",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.NATSPort != 4222 {
					t.Errorf("MQTT.NATSPort = %d, want 4222", cfg.MQTT.NATSPort)
				}
			},
		},
		{
			name: "MQTT ReconnectBackoff defaults",
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.ReconnectBackoff.InitialInterval != "1s" {
					t.Errorf("ReconnectBackoff.InitialInterval = %q, want %q", cfg.MQTT.ReconnectBackoff.InitialInterval, "1s")
				}
				if cfg.MQTT.ReconnectBackoff.MaxInterval != "30s" {
					t.Errorf("ReconnectBackoff.MaxInterval = %q, want %q", cfg.MQTT.ReconnectBackoff.MaxInterval, "30s")
				}
				if cfg.MQTT.ReconnectBackoff.Multiplier != 2.0 {
					t.Errorf("ReconnectBackoff.Multiplier = %f, want 2.0", cfg.MQTT.ReconnectBackoff.Multiplier)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			tt.check(t, cfg)
		})
	}
}

func TestLoadFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		check    func(t *testing.T, cfg *Config)
	}{
		{
			name: "all sections present",
			yaml: `
mqtt:
  port: 1884
  client_id: "test-client"
  clean_session: true
  keep_alive: 60
  qos: 0
  qos_per_topic:
    esp32/vibration: 0
    esp32/temperature: 1
  nats_host: "192.168.1.100"
  nats_port: 4223
  reconnect_backoff:
    initial_interval: "2s"
    max_interval: "60s"
    multiplier: 3.0
  topics:
    subscribe: "test/#"
    status: "sys/test/status"
validation:
  ranges:
    temperature:
      min: -20.0
      max: 100.0
  timestamp:
    max_future_drift: "10s"
    max_past_drift: "12h"
  health:
    max_payload_size: 2048
    min_payload_size: 2
assets:
  auto_register: false
  default_site: "test-site"
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.Port != 1884 {
					t.Errorf("MQTT.Port = %d, want 1884", cfg.MQTT.Port)
				}
				if cfg.MQTT.ClientID != "test-client" {
					t.Errorf("MQTT.ClientID = %q, want %q", cfg.MQTT.ClientID, "test-client")
				}
				if cfg.MQTT.CleanSession != true {
					t.Errorf("MQTT.CleanSession = %v, want true", cfg.MQTT.CleanSession)
				}
				if cfg.MQTT.KeepAlive != 60 {
					t.Errorf("MQTT.KeepAlive = %d, want 60", cfg.MQTT.KeepAlive)
				}
				if cfg.MQTT.QoS != 0 {
					t.Errorf("MQTT.QoS = %d, want 0", cfg.MQTT.QoS)
				}
				if cfg.MQTT.QoSPerTopic["esp32/vibration"] != 0 {
					t.Errorf("QoSPerTopic[esp32/vibration] = %d, want 0", cfg.MQTT.QoSPerTopic["esp32/vibration"])
				}
				if cfg.MQTT.QoSPerTopic["esp32/temperature"] != 1 {
					t.Errorf("QoSPerTopic[esp32/temperature] = %d, want 1", cfg.MQTT.QoSPerTopic["esp32/temperature"])
				}
				if cfg.MQTT.NATSHost != "192.168.1.100" {
					t.Errorf("MQTT.NATSHost = %q, want %q", cfg.MQTT.NATSHost, "192.168.1.100")
				}
				if cfg.MQTT.NATSPort != 4223 {
					t.Errorf("MQTT.NATSPort = %d, want 4223", cfg.MQTT.NATSPort)
				}
				if cfg.MQTT.ReconnectBackoff.InitialInterval != "2s" {
					t.Errorf("ReconnectBackoff.InitialInterval = %q, want %q", cfg.MQTT.ReconnectBackoff.InitialInterval, "2s")
				}
				if cfg.MQTT.ReconnectBackoff.MaxInterval != "60s" {
					t.Errorf("ReconnectBackoff.MaxInterval = %q, want %q", cfg.MQTT.ReconnectBackoff.MaxInterval, "60s")
				}
				if cfg.MQTT.ReconnectBackoff.Multiplier != 3.0 {
					t.Errorf("ReconnectBackoff.Multiplier = %f, want 3.0", cfg.MQTT.ReconnectBackoff.Multiplier)
				}
				if cfg.MQTT.Topics.Subscribe != "test/#" {
					t.Errorf("MQTT.Topics.Subscribe = %q, want %q", cfg.MQTT.Topics.Subscribe, "test/#")
				}
				if cfg.MQTT.Topics.Status != "sys/test/status" {
					t.Errorf("MQTT.Topics.Status = %q, want %q", cfg.MQTT.Topics.Status, "sys/test/status")
				}
				if cfg.Validation.Timestamp.MaxFutureDrift != "10s" {
					t.Errorf("Timestamp.MaxFutureDrift = %q, want %q", cfg.Validation.Timestamp.MaxFutureDrift, "10s")
				}
				if cfg.Validation.Timestamp.MaxPastDrift != "12h" {
					t.Errorf("Timestamp.MaxPastDrift = %q, want %q", cfg.Validation.Timestamp.MaxPastDrift, "12h")
				}
				if cfg.Validation.Health.MaxPayloadSize != 2048 {
					t.Errorf("Health.MaxPayloadSize = %d, want 2048", cfg.Validation.Health.MaxPayloadSize)
				}
				if cfg.Validation.Health.MinPayloadSize != 2 {
					t.Errorf("Health.MinPayloadSize = %d, want 2", cfg.Validation.Health.MinPayloadSize)
				}
				if cfg.Assets.AutoRegister != false {
					t.Errorf("Assets.AutoRegister = %v, want false", cfg.Assets.AutoRegister)
				}
				if cfg.Assets.DefaultSite != "test-site" {
					t.Errorf("Assets.DefaultSite = %q, want %q", cfg.Assets.DefaultSite, "test-site")
				}
			},
		},
		{
			name: "missing mqtt section uses defaults",
			yaml: `
storage:
  path: ./data/test.db
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.Port != 1883 {
					t.Errorf("MQTT.Port = %d, want 1883 (default)", cfg.MQTT.Port)
				}
				if cfg.MQTT.ClientID != "mqtt-plugin" {
					t.Errorf("MQTT.ClientID = %q, want %q (default)", cfg.MQTT.ClientID, "mqtt-plugin")
				}
			},
		},
		{
			name: "missing validation section uses defaults",
			yaml: `
storage:
  path: ./data/test.db
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Validation.Timestamp.MaxFutureDrift != "5s" {
					t.Errorf("Timestamp.MaxFutureDrift = %q, want %q (default)", cfg.Validation.Timestamp.MaxFutureDrift, "5s")
				}
				if cfg.Validation.Health.MaxPayloadSize != 1024 {
					t.Errorf("Health.MaxPayloadSize = %d, want 1024 (default)", cfg.Validation.Health.MaxPayloadSize)
				}
			},
		},
		{
			name: "missing assets section uses defaults",
			yaml: `
storage:
  path: ./data/test.db
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Assets.AutoRegister != true {
					t.Errorf("Assets.AutoRegister = %v, want true (default)", cfg.Assets.AutoRegister)
				}
				if cfg.Assets.DefaultSite != "factory-1" {
					t.Errorf("Assets.DefaultSite = %q, want %q (default)", cfg.Assets.DefaultSite, "factory-1")
				}
			},
		},
		{
			name: "partial mqtt section merges with defaults",
			yaml: `
mqtt:
  port: 1884
  client_id: "custom-client"
`,
			check: func(t *testing.T, cfg *Config) {
				if cfg.MQTT.Port != 1884 {
					t.Errorf("MQTT.Port = %d, want 1884", cfg.MQTT.Port)
				}
				if cfg.MQTT.ClientID != "custom-client" {
					t.Errorf("MQTT.ClientID = %q, want %q", cfg.MQTT.ClientID, "custom-client")
				}
				// Defaults should be preserved
				if cfg.MQTT.CleanSession != false {
					t.Errorf("MQTT.CleanSession = %v, want false (default)", cfg.MQTT.CleanSession)
				}
				if cfg.MQTT.KeepAlive != 30 {
					t.Errorf("MQTT.KeepAlive = %d, want 30 (default)", cfg.MQTT.KeepAlive)
				}
				if cfg.MQTT.Topics.Subscribe != "esp32/#" {
					t.Errorf("MQTT.Topics.Subscribe = %q, want %q (default)", cfg.MQTT.Topics.Subscribe, "esp32/#")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.yaml), 0o644); err != nil {
				t.Fatalf("failed to write config file: %v", err)
			}

			cfg, err := loadFromFile(configPath)
			if err != nil {
				t.Fatalf("loadFromFile() error = %v", err)
			}
			tt.check(t, cfg)
		})
	}
}
