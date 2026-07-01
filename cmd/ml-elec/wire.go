//go:build wireinject

package main

import (
	"github.com/google/wire"
	"ml-elec/internal/api"
	"ml-elec/internal/config"
	"ml-elec/internal/nats"
	"ml-elec/internal/plugin"
	"ml-elec/internal/storage"
)

// ProvideNATSConfig extracts NATS config from the main config.
func ProvideNATSConfig(cfg *config.Config) *config.NATSConfig {
	return &cfg.NATS
}

// ProvideStorageConfig extracts storage config from the main config.
func ProvideStorageConfig(cfg *config.Config) *config.StorageConfig {
	return &cfg.Storage
}

// ProvideAPIConfig extracts API config from the main config.
func ProvideAPIConfig(cfg *config.Config) *config.APIConfig {
	return &cfg.API
}

// App holds all application components wired together.
type App struct {
	Config     *config.Config
	NATSServer *nats.Server
	Store      *storage.Store
	PluginMgr  *plugin.Manager
	APIServer  *api.Server
}

// InitializeApp wires all components together using Google Wire.
func InitializeApp() (*App, error) {
	wire.Build(
		config.Load,
		ProvideNATSConfig,
		ProvideStorageConfig,
		ProvideAPIConfig,
		nats.New,
		storage.New,
		plugin.NewManager,
		api.NewServer,
		wire.Struct(new(App), "*"),
	)
	return nil, nil
}
