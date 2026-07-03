package storage

import (
	"context"
	"testing"
)

func TestCreateAsset(t *testing.T) {
	store, err := NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Create asset_sensors table for foreign key test
	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			site TEXT NOT NULL DEFAULT 'factory-1',
			area TEXT NOT NULL DEFAULT '',
			line TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'machine',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create assets table: %v", err)
	}

	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS asset_sensors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			asset_id INTEGER NOT NULL,
			sensor_id TEXT NOT NULL,
			sensor_type TEXT NOT NULL DEFAULT 'generic',
			topic TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
			UNIQUE(asset_id, sensor_id)
		)
	`)
	if err != nil {
		t.Fatalf("failed to create asset_sensors table: %v", err)
	}

	ctx := context.Background()

	t.Run("create asset with defaults", func(t *testing.T) {
		asset := &Asset{
			Name: "CNC-01",
		}

		result, err := store.CreateAsset(ctx, asset)
		if err != nil {
			t.Fatalf("CreateAsset() error = %v", err)
		}

		if result.ID == 0 {
			t.Error("CreateAsset() ID = 0, want non-zero")
		}
		if result.Name != "CNC-01" {
			t.Errorf("CreateAsset() Name = %q, want %q", result.Name, "CNC-01")
		}
		if result.Site != "factory-1" {
			t.Errorf("CreateAsset() Site = %q, want %q", result.Site, "factory-1")
		}
		if result.Type != "machine" {
			t.Errorf("CreateAsset() Type = %q, want %q", result.Type, "machine")
		}
		if result.CreatedAt.IsZero() {
			t.Error("CreateAsset() CreatedAt is zero")
		}
		if result.UpdatedAt.IsZero() {
			t.Error("CreateAsset() UpdatedAt is zero")
		}
	})

	t.Run("create asset with all fields", func(t *testing.T) {
		asset := &Asset{
			Name: "CNC-02",
			Site: "factory-2",
			Area: "production",
			Line: "L1",
			Type: "cnc",
		}

		result, err := store.CreateAsset(ctx, asset)
		if err != nil {
			t.Fatalf("CreateAsset() error = %v", err)
		}

		if result.Site != "factory-2" {
			t.Errorf("CreateAsset() Site = %q, want %q", result.Site, "factory-2")
		}
		if result.Area != "production" {
			t.Errorf("CreateAsset() Area = %q, want %q", result.Area, "production")
		}
		if result.Line != "L1" {
			t.Errorf("CreateAsset() Line = %q, want %q", result.Line, "L1")
		}
		if result.Type != "cnc" {
			t.Errorf("CreateAsset() Type = %q, want %q", result.Type, "cnc")
		}
	})

	t.Run("create duplicate asset name returns error", func(t *testing.T) {
		asset := &Asset{
			Name: "CNC-01", // duplicate
		}

		_, err := store.CreateAsset(ctx, asset)
		if err == nil {
			t.Error("CreateAsset() with duplicate name should return error")
		}
	})
}

func TestGetAsset(t *testing.T) {
	store, err := NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Create assets table
	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			site TEXT NOT NULL DEFAULT 'factory-1',
			area TEXT NOT NULL DEFAULT '',
			line TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'machine',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create assets table: %v", err)
	}

	ctx := context.Background()

	// Create test asset
	created, err := store.CreateAsset(ctx, &Asset{Name: "TestAsset"})
	if err != nil {
		t.Fatalf("CreateAsset() error = %v", err)
	}

	t.Run("get existing asset", func(t *testing.T) {
		asset, err := store.GetAsset(ctx, created.ID)
		if err != nil {
			t.Fatalf("GetAsset() error = %v", err)
		}

		if asset.ID != created.ID {
			t.Errorf("GetAsset() ID = %d, want %d", asset.ID, created.ID)
		}
		if asset.Name != "TestAsset" {
			t.Errorf("GetAsset() Name = %q, want %q", asset.Name, "TestAsset")
		}
	})

	t.Run("get non-existent asset returns error", func(t *testing.T) {
		_, err := store.GetAsset(ctx, 99999)
		if err == nil {
			t.Error("GetAsset() for non-existent asset should return error")
		}
	})
}

func TestListAssets(t *testing.T) {
	store, err := NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Create assets table
	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			site TEXT NOT NULL DEFAULT 'factory-1',
			area TEXT NOT NULL DEFAULT '',
			line TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'machine',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create assets table: %v", err)
	}

	ctx := context.Background()

	// Create test assets
	for i := 0; i < 5; i++ {
		_, err := store.CreateAsset(ctx, &Asset{Name: "Asset-" + string(rune('A'+i))})
		if err != nil {
			t.Fatalf("CreateAsset() error = %v", err)
		}
	}

	t.Run("list with pagination", func(t *testing.T) {
		assets, total, err := store.ListAssets(ctx, 1, 2)
		if err != nil {
			t.Fatalf("ListAssets() error = %v", err)
		}

		if total != 5 {
			t.Errorf("ListAssets() total = %d, want 5", total)
		}
		if len(assets) != 2 {
			t.Errorf("ListAssets() returned %d assets, want 2", len(assets))
		}
	})

	t.Run("list second page", func(t *testing.T) {
		assets, total, err := store.ListAssets(ctx, 2, 2)
		if err != nil {
			t.Fatalf("ListAssets() error = %v", err)
		}

		if total != 5 {
			t.Errorf("ListAssets() total = %d, want 5", total)
		}
		if len(assets) != 2 {
			t.Errorf("ListAssets() returned %d assets, want 2", len(assets))
		}
	})

	t.Run("list empty page", func(t *testing.T) {
		assets, total, err := store.ListAssets(ctx, 10, 2)
		if err != nil {
			t.Fatalf("ListAssets() error = %v", err)
		}

		if total != 5 {
			t.Errorf("ListAssets() total = %d, want 5", total)
		}
		if len(assets) != 0 {
			t.Errorf("ListAssets() returned %d assets, want 0", len(assets))
		}
	})
}

func TestCreateAssetSensor(t *testing.T) {
	store, err := NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Create tables
	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			site TEXT NOT NULL DEFAULT 'factory-1',
			area TEXT NOT NULL DEFAULT '',
			line TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'machine',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create assets table: %v", err)
	}

	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS asset_sensors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			asset_id INTEGER NOT NULL,
			sensor_id TEXT NOT NULL,
			sensor_type TEXT NOT NULL DEFAULT 'generic',
			topic TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
			UNIQUE(asset_id, sensor_id)
		)
	`)
	if err != nil {
		t.Fatalf("failed to create asset_sensors table: %v", err)
	}

	ctx := context.Background()

	// Create test asset
	asset, err := store.CreateAsset(ctx, &Asset{Name: "TestAsset"})
	if err != nil {
		t.Fatalf("CreateAsset() error = %v", err)
	}

	t.Run("create sensor for asset", func(t *testing.T) {
		sensor := &AssetSensor{
			SensorID:   "temp-01",
			SensorType: "temperature",
			Topic:      "factory-1/test/temperature",
		}

		result, err := store.CreateAssetSensor(ctx, asset.ID, sensor)
		if err != nil {
			t.Fatalf("CreateAssetSensor() error = %v", err)
		}

		if result.ID == 0 {
			t.Error("CreateAssetSensor() ID = 0, want non-zero")
		}
		if result.AssetID != asset.ID {
			t.Errorf("CreateAssetSensor() AssetID = %d, want %d", result.AssetID, asset.ID)
		}
		if result.SensorID != "temp-01" {
			t.Errorf("CreateAssetSensor() SensorID = %q, want %q", result.SensorID, "temp-01")
		}
	})

	t.Run("create duplicate sensor returns error", func(t *testing.T) {
		sensor := &AssetSensor{
			SensorID:   "temp-01", // duplicate
			SensorType: "temperature",
		}

		_, err := store.CreateAssetSensor(ctx, asset.ID, sensor)
		if err == nil {
			t.Error("CreateAssetSensor() with duplicate sensor_id should return error")
		}
	})
}

func TestGetAssetSensors(t *testing.T) {
	store, err := NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Create tables
	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			site TEXT NOT NULL DEFAULT 'factory-1',
			area TEXT NOT NULL DEFAULT '',
			line TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'machine',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create assets table: %v", err)
	}

	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS asset_sensors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			asset_id INTEGER NOT NULL,
			sensor_id TEXT NOT NULL,
			sensor_type TEXT NOT NULL DEFAULT 'generic',
			topic TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
			UNIQUE(asset_id, sensor_id)
		)
	`)
	if err != nil {
		t.Fatalf("failed to create asset_sensors table: %v", err)
	}

	ctx := context.Background()

	// Create test asset
	asset, err := store.CreateAsset(ctx, &Asset{Name: "TestAsset"})
	if err != nil {
		t.Fatalf("CreateAsset() error = %v", err)
	}

	// Create sensors
	_, err = store.CreateAssetSensor(ctx, asset.ID, &AssetSensor{
		SensorID:   "temp-01",
		SensorType: "temperature",
		Topic:      "factory-1/test/temperature",
	})
	if err != nil {
		t.Fatalf("CreateAssetSensor() error = %v", err)
	}

	_, err = store.CreateAssetSensor(ctx, asset.ID, &AssetSensor{
		SensorID:   "vib-01",
		SensorType: "vibration",
		Topic:      "factory-1/test/vibration",
	})
	if err != nil {
		t.Fatalf("CreateAssetSensor() error = %v", err)
	}

	t.Run("get sensors for asset", func(t *testing.T) {
		sensors, err := store.GetAssetSensors(ctx, asset.ID)
		if err != nil {
			t.Fatalf("GetAssetSensors() error = %v", err)
		}

		if len(sensors) != 2 {
			t.Errorf("GetAssetSensors() returned %d sensors, want 2", len(sensors))
		}
	})

	t.Run("get sensors for non-existent asset returns empty", func(t *testing.T) {
		sensors, err := store.GetAssetSensors(ctx, 99999)
		if err != nil {
			t.Fatalf("GetAssetSensors() error = %v", err)
		}

		if len(sensors) != 0 {
			t.Errorf("GetAssetSensors() returned %d sensors, want 0", len(sensors))
		}
	})
}

func TestCascadeDeleteAsset(t *testing.T) {
	store, err := NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Create tables
	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS assets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			site TEXT NOT NULL DEFAULT 'factory-1',
			area TEXT NOT NULL DEFAULT '',
			line TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'machine',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create assets table: %v", err)
	}

	_, err = store.db.Exec(`
		CREATE TABLE IF NOT EXISTS asset_sensors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			asset_id INTEGER NOT NULL,
			sensor_id TEXT NOT NULL,
			sensor_type TEXT NOT NULL DEFAULT 'generic',
			topic TEXT NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
			UNIQUE(asset_id, sensor_id)
		)
	`)
	if err != nil {
		t.Fatalf("failed to create asset_sensors table: %v", err)
	}

	ctx := context.Background()

	// Create asset with sensors
	asset, err := store.CreateAsset(ctx, &Asset{Name: "DeleteMe"})
	if err != nil {
		t.Fatalf("CreateAsset() error = %v", err)
	}

	_, err = store.CreateAssetSensor(ctx, asset.ID, &AssetSensor{
		SensorID:   "temp-01",
		SensorType: "temperature",
	})
	if err != nil {
		t.Fatalf("CreateAssetSensor() error = %v", err)
	}

	// Delete asset
	_, err = store.db.ExecContext(ctx, "DELETE FROM assets WHERE id = ?", asset.ID)
	if err != nil {
		t.Fatalf("delete asset error = %v", err)
	}

	// Verify sensors are also deleted
	sensors, err := store.GetAssetSensors(ctx, asset.ID)
	if err != nil {
		t.Fatalf("GetAssetSensors() error = %v", err)
	}

	if len(sensors) != 0 {
		t.Errorf("GetAssetSensors() returned %d sensors after cascade delete, want 0", len(sensors))
	}
}
