package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/rs/cors"
	"ml-elec/internal/config"
	"ml-elec/internal/storage"
)

// Server holds the HTTP server and dependencies.
type Server struct {
	server *http.Server
	store  *storage.Store
}

// NewServer creates a new API server with routes and CORS middleware.
func NewServer(cfg *config.APIConfig, store *storage.Store) *Server {
	mux := http.NewServeMux()

	s := &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		store: store,
	}

	// Register routes
	mux.HandleFunc("GET /health", s.HealthHandler)
	mux.HandleFunc("GET /api/v1/sensors", s.SensorsHandler)
	mux.HandleFunc("POST /api/v1/assets", s.CreateAssetHandler)
	mux.HandleFunc("GET /api/v1/assets", s.ListAssetsHandler)
	mux.HandleFunc("GET /api/v1/assets/{id}", s.GetAssetHandler)
	mux.HandleFunc("GET /api/v1/assets/{id}/sensors", s.ListAssetSensorsHandler)
	mux.HandleFunc("POST /api/v1/assets/{id}/sensors", s.CreateAssetSensorHandler)
	mux.HandleFunc("DELETE /api/v1/assets", s.DeleteAssetsHandler)

	// Apply CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	s.server.Handler = c.Handler(mux)

	return s
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	slog.Info("starting API server", "addr", s.server.Addr)
	return s.server.ListenAndServe()
}

// ServerAddr returns the server's listen address.
func (s *Server) ServerAddr() string {
	return s.server.Addr
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
