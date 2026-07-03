package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ml-elec/internal/storage"
)

// AssetRequest represents the JSON body for creating an asset.
type AssetRequest struct {
	Name string `json:"name"`
	Site string `json:"site"`
	Area string `json:"area"`
	Line string `json:"line"`
	Type string `json:"type"`
}

// SensorRequest represents the JSON body for creating an asset sensor.
type SensorRequest struct {
	SensorID   string `json:"sensor_id"`
	SensorType string `json:"sensor_type"`
	Topic      string `json:"topic"`
}

// PaginationResponse wraps paginated results.
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	} `json:"pagination"`
}

// parsePagination extracts page and limit from query parameters.
func parsePagination(r *http.Request) (page, limit int) {
	page = 1
	limit = 20
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	return
}

// CreateAssetHandler creates a new asset.
//
// @Summary Create asset
// @Tags assets
// @Accept json
// @Produce json
// @Param asset body AssetRequest true "Asset data"
// @Success 201 {object} storage.Asset
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/assets [post]
func (s *Server) CreateAssetHandler(w http.ResponseWriter, r *http.Request) {
	var req AssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	asset := &storage.Asset{
		Name: req.Name,
		Site: req.Site,
		Area: req.Area,
		Line: req.Line,
		Type: req.Type,
	}

	created, err := s.store.CreateAsset(r.Context(), asset)
	if err != nil {
		// Check for unique constraint violation (duplicate name)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "duplicate") {
			writeError(w, http.StatusConflict, "asset with this name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create asset")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// ListAssetsHandler returns assets with pagination.
//
// @Summary List assets
// @Tags assets
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} PaginationResponse
// @Router /api/v1/assets [get]
func (s *Server) ListAssetsHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

	assets, total, err := s.store.ListAssets(r.Context(), page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query assets")
		return
	}

	resp := PaginationResponse{Data: assets}
	resp.Pagination.Page = page
	resp.Pagination.Limit = limit
	resp.Pagination.Total = total

	writeJSON(w, http.StatusOK, resp)
}

// GetAssetHandler returns a single asset by ID.
//
// @Summary Get asset
// @Tags assets
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {object} storage.Asset
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/assets/{id} [get]
func (s *Server) GetAssetHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid asset ID")
		return
	}

	asset, err := s.store.GetAsset(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			writeError(w, http.StatusNotFound, "asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get asset")
		return
	}

	writeJSON(w, http.StatusOK, asset)
}

// ListAssetSensorsHandler returns sensors for a given asset.
//
// @Summary List asset sensors
// @Tags assets
// @Produce json
// @Param id path int true "Asset ID"
// @Success 200 {array} storage.AssetSensor
// @Router /api/v1/assets/{id}/sensors [get]
func (s *Server) ListAssetSensorsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid asset ID")
		return
	}

	sensors, err := s.store.GetAssetSensors(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query sensors")
		return
	}

	// Ensure null slice becomes empty array in JSON
	if sensors == nil {
		sensors = []storage.AssetSensor{}
	}

	writeJSON(w, http.StatusOK, sensors)
}

// CreateAssetSensorHandler creates a sensor for an asset.
//
// @Summary Create asset sensor
// @Tags assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param sensor body SensorRequest true "Sensor data"
// @Success 201 {object} storage.AssetSensor
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/assets/{id}/sensors [post]
func (s *Server) CreateAssetSensorHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	assetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid asset ID")
		return
	}

	var req SensorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.SensorID) == "" {
		writeError(w, http.StatusBadRequest, "sensor_id is required")
		return
	}

	sensor := &storage.AssetSensor{
		SensorID:   req.SensorID,
		SensorType: req.SensorType,
		Topic:      req.Topic,
	}

	created, err := s.store.CreateAssetSensor(r.Context(), assetID, sensor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create sensor")
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// ListAllSensorsHandler returns all sensors globally with pagination.
//
// @Summary List all sensors
// @Tags sensors
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(20)
// @Success 200 {object} PaginationResponse
// @Router /api/v1/sensors [get]
func (s *Server) ListAllSensorsHandler(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePagination(r)

	sensors, total, err := s.store.ListAllSensors(r.Context(), page, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query sensors")
		return
	}

	// Ensure null slice becomes empty array in JSON
	if sensors == nil {
		sensors = []storage.AssetSensor{}
	}

	resp := PaginationResponse{Data: sensors}
	resp.Pagination.Page = page
	resp.Pagination.Limit = limit
	resp.Pagination.Total = total

	writeJSON(w, http.StatusOK, resp)
}

// DeleteAssetsHandler returns 405 Method Not Allowed (v1 prohibition).
//
// @Summary Delete assets (not allowed)
// @Tags assets
// @Failure 405 {object} ErrorResponse
// @Router /api/v1/assets [delete]
func (s *Server) DeleteAssetsHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusMethodNotAllowed, "delete assets is not allowed")
}
