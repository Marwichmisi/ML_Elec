package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Structured logging to stderr
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	// Signal handling: SIGINT and SIGTERM trigger graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize all components via Wire DI
	app, err := InitializeApp()
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
