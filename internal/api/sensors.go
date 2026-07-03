package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ml-elec/internal/storage"
)

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SensorsHandler handles GET /api/v1/sensors.
// With ?sensor_id= param: returns sensor readings for that sensor.
// Without ?sensor_id=: returns all sensors globally with pagination.
func (s *Server) SensorsHandler(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get("sensor_id")
	if sensorID != "" {
		s.GetSensorsHandler(w, r)
		return
	}
	s.ListAllSensorsHandler(w, r)
}

// GetSensorsHandler returns sensor readings for a given sensor_id.
//
// @Summary Get sensor readings
// @Tags sensors
// @Produce json
// @Param sensor_id query string true "Sensor ID"
// @Param limit query int false "Limit" default(100)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/sensors [get]
func (s *Server) GetSensorsHandler(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get("sensor_id")
	if sensorID == "" {
		writeError(w, http.StatusBadRequest, "sensor_id is required")
		return
	}

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	readings, err := s.store.GetSensors(r.Context(), sensorID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query sensors")
		return
	}

	// Ensure data is always an array, never null
	if readings == nil {
		readings = []storage.SensorReading{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": readings,
	})
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Response already started, can't change status — log only
		_ = err
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
