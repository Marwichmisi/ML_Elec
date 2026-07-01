package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application.
type Config struct {
	NATS    NATSConfig    `yaml:"nats"`
	Storage StorageConfig `yaml:"storage"`
	API     APIConfig     `yaml:"api"`
	Plugins PluginsConfig `yaml:"plugins"`
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
	}
}

// Load searches for a config file in multiple locations and returns the parsed config.
// Search order: ./config.yaml, ~/.config/ml-elec/config.yaml, then defaults.
func Load() (*Config, error) {
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
