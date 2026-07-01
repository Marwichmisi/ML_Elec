package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewOpensDBInWALMode(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer store.Close()

	// Verify WAL mode
	var journalMode string
	err = store.db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("expected WAL mode, got %s", journalMode)
	}
}

func TestInsertAndGetSensors(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	sensorID := "sensor-1"
	value := 23.5
	ts := time.Now().UTC()

	// Insert a reading
	err = store.InsertSensor(ctx, sensorID, value, ts)
	if err != nil {
		t.Fatalf("InsertSensor() failed: %v", err)
	}

	// Get readings
	readings, err := store.GetSensors(ctx, sensorID, 10)
	if err != nil {
		t.Fatalf("GetSensors() failed: %v", err)
	}

	if len(readings) != 1 {
		t.Fatalf("expected 1 reading, got %d", len(readings))
	}

	if readings[0].SensorID != sensorID {
		t.Errorf("expected sensor_id %s, got %s", sensorID, readings[0].SensorID)
	}
	if readings[0].Value != value {
		t.Errorf("expected value %f, got %f", value, readings[0].Value)
	}
}

func TestGetSensorsReturnsEmptyForNonExistentSensor(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	readings, err := store.GetSensors(ctx, "non-existent", 10)
	if err != nil {
		t.Fatalf("GetSensors() failed: %v", err)
	}

	if len(readings) != 0 {
		t.Errorf("expected 0 readings, got %d", len(readings))
	}
}

func TestGetSensorsRespectsLimit(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	sensorID := "sensor-1"

	// Insert 5 readings
	for i := 0; i < 5; i++ {
		ts := time.Now().UTC().Add(time.Duration(i) * time.Second)
		err = store.InsertSensor(ctx, sensorID, float64(i), ts)
		if err != nil {
			t.Fatalf("InsertSensor() failed: %v", err)
		}
	}

	// Get with limit 3
	readings, err := store.GetSensors(ctx, sensorID, 3)
	if err != nil {
		t.Fatalf("GetSensors() failed: %v", err)
	}

	if len(readings) != 3 {
		t.Errorf("expected 3 readings, got %d", len(readings))
	}
}

func TestCheckCorruptionReturnsNilForHealthyDB(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	defer store.Close()

	err = store.CheckCorruption()
	if err != nil {
		t.Errorf("CheckCorruption() returned error for healthy DB: %v", err)
	}
}

func TestCheckCorruptionReturnsErrorForCorruptedDB(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	// Create a valid DB first
	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	store.Close()

	// Corrupt the DB file by writing garbage
	f, err := os.OpenFile(dbPath, os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("failed to open DB for corruption: %v", err)
	}
	_, err = f.Write([]byte("THIS IS CORRUPTED DATA"))
	if err != nil {
		t.Fatalf("failed to write corrupted data: %v", err)
	}
	f.Close()

	// Try to open and check corruption
	store2, err := New(dbPath)
	if err != nil {
		// If New() fails due to corruption, that's acceptable
		t.Logf("New() failed for corrupted DB (acceptable): %v", err)
		return
	}
	defer store2.Close()

	err = store2.CheckCorruption()
	if err == nil {
		t.Error("CheckCorruption() should return error for corrupted DB")
	}
}

func TestClose(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := New(dbPath)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	err = store.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}
