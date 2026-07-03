package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application.
type Config struct {
	NATS       NATSConfig       `yaml:"nats"`
	Storage    StorageConfig    `yaml:"storage"`
	API        APIConfig        `yaml:"api"`
	Plugins    PluginsConfig    `yaml:"plugins"`
	MQTT       MQTTConfig       `yaml:"mqtt"`
	Validation ValidationConfig `yaml:"validation"`
	Assets     AssetsConfig     `yaml:"assets"`
}

// MQTTConfig holds MQTT broker and client configuration.
type MQTTConfig struct {
	Port         int               `yaml:"port"`
	ClientID     string            `yaml:"client_id"`
	CleanSession bool              `yaml:"clean_session"`
	KeepAlive    int               `yaml:"keep_alive"`
	QoS          byte              `yaml:"qos"`            // Default QoS for subscriptions (0 or 1)
	QoSPerTopic  map[string]byte   `yaml:"qos_per_topic"` // Per-topic QoS overrides
	NATSHost     string            `yaml:"nats_host"`     // NATS server host to connect to
	NATSPort     int               `yaml:"nats_port"`     // NATS server port to connect to
	ReconnectBackoff ReconnectBackoffConfig `yaml:"reconnect_backoff"`
	Topics       struct {
		Subscribe string `yaml:"subscribe"`
		Status    string `yaml:"status"`
	} `yaml:"topics"`
}

// ReconnectBackoffConfig defines exponential backoff parameters for MQTT reconnection.
type ReconnectBackoffConfig struct {
	InitialInterval string  `yaml:"initial_interval"` // e.g. "1s"
	MaxInterval     string  `yaml:"max_interval"`     // e.g. "30s"
	Multiplier      float64 `yaml:"multiplier"`       // e.g. 2.0
}

// ValidationConfig holds 3-level data validation configuration.
type ValidationConfig struct {
	Ranges    map[string]RangeConfig `yaml:"ranges"`
	Timestamp TimestampConfig        `yaml:"timestamp"`
	Health    HealthConfig           `yaml:"health"`
}

// RangeConfig defines min/max bounds for a sensor type.
type RangeConfig struct {
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
}

// TimestampConfig defines acceptable timestamp drift bounds.
type TimestampConfig struct {
	MaxFutureDrift string `yaml:"max_future_drift"`
	MaxPastDrift   string `yaml:"max_past_drift"`
}

// HealthConfig defines payload health validation bounds.
type HealthConfig struct {
	MaxPayloadSize int `yaml:"max_payload_size"`
	MinPayloadSize int `yaml:"min_payload_size"`
}

// AssetsConfig holds asset registry configuration.
type AssetsConfig struct {
	AutoRegister bool   `yaml:"auto_register"`
	DefaultSite  string `yaml:"default_site"`
}

// NATSConfig holds NATS server configuration.
type NATSConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// StorageConfig holds SQLite storage configuration.
type StorageConfig struct {
	Path string `yaml:"path"`
}

// APIConfig holds REST API configuration.
type APIConfig struct {
	Port int `yaml:"port"`
}

// PluginsConfig holds plugin configuration.
type PluginsConfig struct {
	Enabled []string `yaml:"enabled"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		NATS: NATSConfig{
			Host: "127.0.0.1",
			Port: -1,
		},
		Storage: StorageConfig{
			Path: "./data/sensors.db",
		},
		API: APIConfig{
			Port: 8080,
		},
		Plugins: PluginsConfig{
			Enabled: []string{},
		},
		MQTT: MQTTConfig{
			Port:         1883,
			ClientID:     "mqtt-plugin",
			CleanSession: false,
			KeepAlive:    30,
			QoS:          1,
			QoSPerTopic:  make(map[string]byte),
			NATSHost:     "127.0.0.1",
			NATSPort:     4222,
			ReconnectBackoff: ReconnectBackoffConfig{
				InitialInterval: "1s",
				MaxInterval:     "30s",
				Multiplier:      2.0,
			},
			Topics: struct {
				Subscribe string `yaml:"subscribe"`
				Status    string `yaml:"status"`
			}{
				Subscribe: "esp32/#",
				Status:    "sys/mqtt-plugin/status",
			},
		},
		Validation: ValidationConfig{
			Ranges: map[string]RangeConfig{
				"temperature": {Min: -40.0, Max: 150.0},
				"humidity":    {Min: 0.0, Max: 100.0},
				"current":     {Min: 0.0, Max: 1000.0},
			},
			Timestamp: TimestampConfig{
				MaxFutureDrift: "5s",
				MaxPastDrift:   "24h",
			},
			Health: HealthConfig{
				MaxPayloadSize: 1024,
				MinPayloadSize: 1,
			},
		},
		Assets: AssetsConfig{
			AutoRegister: true,
			DefaultSite:  "factory-1",
		},
	}
}

// Load searches for a config file in multiple locations and returns the parsed config.
// Search order: CONFIG_PATH env var, ./config.yaml, ~/.config/ml-elec/config.yaml, then defaults.
func Load() (*Config, error) {
	// Check CONFIG_PATH environment variable first
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			return loadFromFile(configPath)
		}
	}

	home, _ := os.UserHomeDir()

	paths := []string{
		"./config.yaml",
		filepath.Join(home, ".config", "ml-elec", "config.yaml"),
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return loadFromFile(path)
		}
	}

	return DefaultConfig(), nil
}

func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}

	return cfg, nil
}

// Validate checks that required fields are present.
func (c *Config) Validate() error {
	if c.NATS.Host == "" {
		return fmt.Errorf("nats.host is required")
	}
	if c.Storage.Path == "" {
		return fmt.Errorf("storage.path is required")
	}
	return nil
}
