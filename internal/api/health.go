package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler returns 200 OK with a JSON health status.
//
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{Status: "ok"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode health response", "error", err)
	}
}
