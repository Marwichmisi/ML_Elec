package main

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"ml-elec/internal/api"
	"ml-elec/internal/config"
	"ml-elec/internal/nats"
	"ml-elec/internal/plugin"
	"ml-elec/internal/storage"
	sdk "ml-elec/pkg/sdk/v1"
)

// App holds all application components.
type App struct {
	Config     *config.Config
	NATSServer *nats.Server
	Store      *storage.Store
	PluginMgr  *plugin.Manager
	APIServer  *api.Server
}

// InitializeApp creates and wires all components manually.
func InitializeApp(cfg *config.Config) (*App, error) {
	natsServer, err := nats.New(&cfg.NATS)
	if err != nil {
		return nil, err
	}

	store, err := storage.New(&cfg.Storage)
	if err != nil {
		return nil, err
	}

	pluginMgr := plugin.NewManager()
	apiServer := api.NewServer(&cfg.API, store)

	return &App{
		Config:     cfg,
		NATSServer: natsServer,
		Store:      store,
		PluginMgr:  pluginMgr,
		APIServer:  apiServer,
	}, nil
}

// launchPlugins starts enabled plugins (e.g. MQTT plugin) as child processes.
func launchPlugins(app *App) {
	for _, name := range app.Config.Plugins.Enabled {
		// Look for plugin binary in ./bin/ directory
		binPath := filepath.Join("bin", name)
		if _, err := exec.LookPath(binPath); err != nil {
			slog.Warn("plugin binary not found, skipping", "name", name, "path", binPath)
			continue
		}

		// All Phase 2 plugins use gRPC transport
		if err := app.PluginMgr.LaunchGRPC(name, binPath, app.Config.Plugins.Enabled, &sdk.GRPCPlugin{}); err != nil {
			slog.Error("failed to launch plugin", "name", name, "error", err)
			continue
		}
		slog.Info("plugin launched", "name", name, "path", binPath)
	}
}

func main() {
	// Structured logging to stderr
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	// Signal handling: SIGINT and SIGTERM trigger graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Allow port override via ML_ELEC_PORT environment variable
	if portStr := os.Getenv("ML_ELEC_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.API.Port = port
		}
	}

	// Initialize all components
	app, err := InitializeApp(cfg)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	slog.Info("all components initialized",
		"nats_host", app.Config.NATS.Host,
		"nats_port", app.Config.NATS.Port,
		"storage_path", app.Config.Storage.Path,
		"api_port", app.Config.API.Port,
	)

	// Start components in order: NATS → Storage (verify) → Plugins → API
	if err := app.NATSServer.Start(); err != nil {
		slog.Error("failed to start NATS server", "error", err)
		os.Exit(1)
	}
	slog.Info("NATS server started")

	if err := app.Store.CheckCorruption(); err != nil {
		slog.Error("storage corruption detected", "error", err)
		os.Exit(1)
	}
	slog.Info("storage verified", "path", app.Config.Storage.Path)

	// Launch enabled plugins (MQTT, etc.) as child processes
	launchPlugins(app)

	// Start API server in goroutine (ListenAndServe blocks)
	go func() {
		if err := app.APIServer.Start(); err != nil && err.Error() != "http: Server closed" {
			slog.Error("API server error", "error", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("shutdown signal received")

	// LIFO shutdown with 30-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1st: Stop accepting new requests (API)
	if err := app.APIServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("failed to shutdown API server", "error", err)
	} else {
		slog.Info("API server stopped")
	}

	// 2nd: Kill all plugins
	app.PluginMgr.ShutdownAll()
	slog.Info("plugins stopped")

	// 3rd: Stop NATS
	app.NATSServer.Shutdown()
	slog.Info("NATS server stopped")

	// 4th: Close database
	if err := app.Store.Close(); err != nil {
		slog.Error("failed to close storage", "error", err)
	} else {
		slog.Info("storage closed")
	}

	slog.Info("shutdown complete")
}
