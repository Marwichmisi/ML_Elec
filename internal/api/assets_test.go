package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ml-elec/internal/config"
	"ml-elec/internal/storage"
)

// helper to create a test server with a fresh store
func newTestServer(t *testing.T) (*Server, *httptest.Server, func()) {
	t.Helper()
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)
	ts := httptest.NewServer(srv.server.Handler)
	return srv, ts, func() { ts.Close(); store.Close() }
}

func TestCreateAssetValidRequest(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	body := `{"name":"CNC-01","site":"factory-1","area":"production","line":"L1","type":"cnc"}`
	resp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST /api/v1/assets failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var asset storage.Asset
	if err := json.NewDecoder(resp.Body).Decode(&asset); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if asset.ID == 0 {
		t.Error("expected non-zero asset ID")
	}
	if asset.Name != "CNC-01" {
		t.Errorf("expected name CNC-01, got %q", asset.Name)
	}
	if asset.CreatedAt.IsZero() {
		t.Error("expected non-zero created_at")
	}
}

func TestCreateAssetDuplicateName(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	body := `{"name":"CNC-01","site":"factory-1"}`
	resp1, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("first POST failed: %v", err)
	}
	resp1.Body.Close()

	if resp1.StatusCode != http.StatusCreated {
		t.Fatalf("first POST: expected 201, got %d", resp1.StatusCode)
	}

	// Second POST with same name should return 409
	resp2, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("second POST failed: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("expected status 409 for duplicate name, got %d", resp2.StatusCode)
	}
}

func TestCreateAssetInvalidBody(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	resp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader("{invalid"))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateAssetMissingName(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	body := `{"site":"factory-1","area":"production"}`
	resp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for missing name, got %d", resp.StatusCode)
	}
}

func TestListAssetsPagination(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	// Insert 3 assets
	for i := 1; i <= 3; i++ {
		body := `{"name":"asset-` + string(rune('0'+i)) + `","site":"factory-1"}`
		resp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatalf("POST asset %d failed: %v", i, err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("POST asset %d: expected 201, got %d", i, resp.StatusCode)
		}
	}

	// GET page 1, limit 2
	resp, err := http.Get(ts.URL + "/api/v1/assets?page=1&limit=2")
	if err != nil {
		t.Fatalf("GET /api/v1/assets failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("expected data to be an array")
	}
	if len(data) != 2 {
		t.Errorf("expected 2 items on page 1, got %d", len(data))
	}

	pagination, ok := result["pagination"].(map[string]interface{})
	if !ok {
		t.Fatal("expected pagination object")
	}
	total, ok := pagination["total"].(float64)
	if !ok || int(total) != 3 {
		t.Errorf("expected pagination.total=3, got %v", pagination["total"])
	}

	// GET page 2, limit 2
	resp2, err := http.Get(ts.URL + "/api/v1/assets?page=2&limit=2")
	if err != nil {
		t.Fatalf("GET page 2 failed: %v", err)
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&result2); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data2, ok := result2["data"].([]interface{})
	if !ok {
		t.Fatal("expected data to be an array")
	}
	if len(data2) != 1 {
		t.Errorf("expected 1 item on page 2, got %d", len(data2))
	}
}

func TestGetAssetByID(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	// Create an asset
	body := `{"name":"CNC-01","site":"factory-1","area":"production","line":"L1","type":"cnc"}`
	createResp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer createResp.Body.Close()

	var created storage.Asset
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode created asset: %v", err)
	}

	// GET by ID
	resp, err := http.Get(ts.URL + "/api/v1/assets/1")
	if err != nil {
		t.Fatalf("GET /api/v1/assets/1 failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var asset storage.Asset
	if err := json.NewDecoder(resp.Body).Decode(&asset); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if asset.Name != "CNC-01" {
		t.Errorf("expected name CNC-01, got %q", asset.Name)
	}
}

func TestGetAssetNotFound(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/api/v1/assets/99999")
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestCreateAssetSensor(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	// Create asset
	assetBody := `{"name":"CNC-01","site":"factory-1"}`
	assetResp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(assetBody))
	if err != nil {
		t.Fatalf("POST asset failed: %v", err)
	}
	assetResp.Body.Close()

	// Create sensor for asset
	sensorBody := `{"sensor_id":"vibration-01","sensor_type":"vibration","topic":"esp32/vibration"}`
	resp, err := http.Post(ts.URL+"/api/v1/assets/1/sensors", "application/json", strings.NewReader(sensorBody))
	if err != nil {
		t.Fatalf("POST sensor failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	var sensor storage.AssetSensor
	if err := json.NewDecoder(resp.Body).Decode(&sensor); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if sensor.SensorID != "vibration-01" {
		t.Errorf("expected sensor_id vibration-01, got %q", sensor.SensorID)
	}
	if sensor.AssetID != 1 {
		t.Errorf("expected asset_id 1, got %d", sensor.AssetID)
	}
}

func TestListAssetSensors(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	// Create asset
	assetBody := `{"name":"CNC-01","site":"factory-1"}`
	assetResp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(assetBody))
	if err != nil {
		t.Fatalf("POST asset failed: %v", err)
	}
	assetResp.Body.Close()

	// Add two sensors
	sensorBody1 := `{"sensor_id":"vibration-01","sensor_type":"vibration","topic":"esp32/vibration"}`
	resp1, err := http.Post(ts.URL+"/api/v1/assets/1/sensors", "application/json", strings.NewReader(sensorBody1))
	if err != nil {
		t.Fatalf("POST sensor 1 failed: %v", err)
	}
	resp1.Body.Close()

	sensorBody2 := `{"sensor_id":"temp-01","sensor_type":"temperature","topic":"esp32/temp"}`
	resp2, err := http.Post(ts.URL+"/api/v1/assets/1/sensors", "application/json", strings.NewReader(sensorBody2))
	if err != nil {
		t.Fatalf("POST sensor 2 failed: %v", err)
	}
	resp2.Body.Close()

	// List sensors for asset
	resp, err := http.Get(ts.URL + "/api/v1/assets/1/sensors")
	if err != nil {
		t.Fatalf("GET sensors failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var sensors []storage.AssetSensor
	if err := json.NewDecoder(resp.Body).Decode(&sensors); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(sensors) != 2 {
		t.Errorf("expected 2 sensors, got %d", len(sensors))
	}
}

func TestListAllSensors(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	// Create two assets with sensors
	assetBody1 := `{"name":"CNC-01","site":"factory-1"}`
	resp1, _ := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(assetBody1))
	resp1.Body.Close()

	assetBody2 := `{"name":"CNC-02","site":"factory-1"}`
	resp2, _ := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(assetBody2))
	resp2.Body.Close()

	sensorBody1 := `{"sensor_id":"vibration-01","sensor_type":"vibration","topic":"esp32/vibration"}`
	sResp1, _ := http.Post(ts.URL+"/api/v1/assets/1/sensors", "application/json", strings.NewReader(sensorBody1))
	sResp1.Body.Close()

	sensorBody2 := `{"sensor_id":"temp-01","sensor_type":"temperature","topic":"esp32/temp"}`
	sResp2, _ := http.Post(ts.URL+"/api/v1/assets/2/sensors", "application/json", strings.NewReader(sensorBody2))
	sResp2.Body.Close()

	// List all sensors
	resp, err := http.Get(ts.URL + "/api/v1/sensors")
	if err != nil {
		t.Fatalf("GET /api/v1/sensors failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := result["data"].([]interface{})
	if !ok {
		t.Fatal("expected data to be an array")
	}
	if len(data) != 2 {
		t.Errorf("expected 2 sensors globally, got %d", len(data))
	}

	pagination, ok := result["pagination"].(map[string]interface{})
	if !ok {
		t.Fatal("expected pagination object")
	}
	total, ok := pagination["total"].(float64)
	if !ok || int(total) != 2 {
		t.Errorf("expected pagination.total=2, got %v", pagination["total"])
	}
}

func TestDeleteAssetsReturns405(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/api/v1/assets", nil)
	if err != nil {
		t.Fatalf("creating DELETE request failed: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /api/v1/assets failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", resp.StatusCode)
	}
}

func TestCreateAssetMissingRequiredFields(t *testing.T) {
	_, ts, cleanup := newTestServer(t)
	defer cleanup()

	// Empty name
	body := `{"name":"","site":"factory-1"}`
	resp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty name, got %d", resp.StatusCode)
	}
}
