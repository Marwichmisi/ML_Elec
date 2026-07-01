package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ml-elec/internal/config"
	"ml-elec/internal/storage"
)

func TestNewServerCreatesServer(t *testing.T) {
	// Create a temp DB for the store
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	if srv == nil {
		t.Fatal("NewServer returned nil")
	}
	if srv.server == nil {
		t.Fatal("NewServer created server with nil http.Server")
	}
	if srv.store == nil {
		t.Fatal("NewServer created server with nil store")
	}
}

func TestHealthEndpoint(t *testing.T) {
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	// Use httptest to test the handler directly
	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if healthResp.Status != "ok" {
		t.Errorf("expected status \"ok\", got %q", healthResp.Status)
	}
}

func TestGetSensorsNoSensorID(t *testing.T) {
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/sensors")
	if err != nil {
		t.Fatalf("GET /api/v1/sensors failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestGetSensorsEmptyResult(t *testing.T) {
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/sensors?sensor_id=nonexistent")
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

	data, ok := result["data"]
	if !ok {
		t.Fatal("expected \"data\" key in response")
	}

	dataSlice, ok := data.([]interface{})
	if !ok {
		t.Fatal("expected data to be an array")
	}

	if len(dataSlice) != 0 {
		t.Errorf("expected empty data array, got %d items", len(dataSlice))
	}
}

func TestGetSensorsValidRequest(t *testing.T) {
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Insert test data
	ctx := context.Background()
	now := time.Now()
	for i := 0; i < 5; i++ {
		ts := now.Add(time.Duration(i) * time.Second)
		if err := store.InsertSensor(ctx, "sensor-test", float64(i)*1.5, ts); err != nil {
			t.Fatalf("failed to insert sensor reading: %v", err)
		}
	}

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	// Query with limit
	resp, err := http.Get(ts.URL + "/api/v1/sensors?sensor_id=sensor-test&limit=3")
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

	data, ok := result["data"]
	if !ok {
		t.Fatal("expected \"data\" key in response")
	}

	dataSlice, ok := data.([]interface{})
	if !ok {
		t.Fatal("expected data to be an array")
	}

	if len(dataSlice) != 3 {
		t.Errorf("expected 3 readings (limit=3), got %d", len(dataSlice))
	}

	// Verify each reading has expected fields
	for i, item := range dataSlice {
		reading, ok := item.(map[string]interface{})
		if !ok {
			t.Fatalf("reading %d is not a JSON object", i)
		}
		if _, ok := reading["sensor_id"]; !ok {
			t.Errorf("reading %d missing sensor_id", i)
		}
		if _, ok := reading["value"]; !ok {
			t.Errorf("reading %d missing value", i)
		}
		if _, ok := reading["timestamp"]; !ok {
			t.Errorf("reading %d missing timestamp", i)
		}
	}
}

func TestConcurrentHealthRequests(t *testing.T) {
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	const numGoroutines = 10
	done := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			resp, err := http.Get(ts.URL + "/health")
			if err != nil {
				done <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
					done <- fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				return
			}
			done <- nil
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		if err := <-done; err != nil {
			t.Errorf("goroutine %d failed: %v", i, err)
		}
	}
}

func TestJSONResponseFormat(t *testing.T) {
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	// Test error response format
	resp, err := http.Get(ts.URL + "/api/v1/sensors")
	if err != nil {
		t.Fatalf("GET /api/v1/sensors failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check that response has "error" key
	if _, ok := result["error"]; !ok {
		t.Error("expected \"error\" key in error response")
	}

	// Test success response format
	resp2, err := http.Get(ts.URL + "/api/v1/sensors?sensor_id=test")
	if err != nil {
		t.Fatalf("GET /api/v1/sensors failed: %v", err)
	}
	defer resp2.Body.Close()

	var result2 map[string]interface{}
	if err := json.NewDecoder(resp2.Body).Decode(&result2); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check that response has "data" key
	if _, ok := result2["data"]; !ok {
		t.Error("expected \"data\" key in success response")
	}
}
