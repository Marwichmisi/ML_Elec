package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"ml-elec/internal/config"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// SensorReading represents a sensor data point.
type SensorReading struct {
	ID        int64     `json:"id"`
	SensorID  string    `json:"sensor_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// Asset represents a machine or equipment in the ISA-95 hierarchy.
type Asset struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Site      string    `json:"site"`
	Area      string    `json:"area"`
	Line      string    `json:"line"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AssetSensor represents a sensor linked to an asset.
type AssetSensor struct {
	ID         int64     `json:"id"`
	AssetID    int64     `json:"asset_id"`
	SensorID   string    `json:"sensor_id"`
	SensorType string    `json:"sensor_type"`
	Topic      string    `json:"topic"`
	CreatedAt  time.Time `json:"created_at"`
}

// Store provides sensor data storage with SQLite WAL mode.
type Store struct {
	db *sql.DB
}

// New opens a SQLite database in WAL mode, runs migrations, and returns a Store.
func New(cfg *config.StorageConfig) (*Store, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating storage directory: %w", err)
	}

	dsn := "file:" + cfg.Path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Verify WAL mode
	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		db.Close()
		return nil, fmt.Errorf("verifying WAL mode: %w", err)
	}
	if journalMode != "wal" {
		db.Close()
		return nil, fmt.Errorf("WAL mode not enabled: got %s", journalMode)
	}

	// Run migrations
	if err := runMigrations(db, cfg.Path); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &Store{db: db}, nil
}

func runMigrations(db *sql.DB, dbPath string) error {
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("creating migration source: %w", err)
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("creating sqlite driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}
	// Note: m.Close() would close the underlying database connection.
	// We skip it here — the Store.Close() method handles database cleanup.

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}

// InsertSensor inserts a sensor reading.
func (s *Store) InsertSensor(ctx context.Context, sensorID string, value float64, ts time.Time) error {
	query, args, err := squirrel.Insert("sensor_readings").
		Columns("sensor_id", "value", "timestamp").
		Values(sensorID, value, ts).
		ToSql()
	if err != nil {
		return fmt.Errorf("building insert query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("inserting sensor reading: %w", err)
	}

	return nil
}

// GetSensors retrieves sensor readings for a given sensor ID, limited by the limit parameter.
func (s *Store) GetSensors(ctx context.Context, sensorID string, limit int) ([]SensorReading, error) {
	query, args, err := squirrel.Select("id", "sensor_id", "value", "timestamp").
		From("sensor_readings").
		Where(squirrel.Eq{"sensor_id": sensorID}).
		OrderBy("timestamp DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select query: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying sensor readings: %w", err)
	}
	defer rows.Close()

	var readings []SensorReading
	for rows.Next() {
		var r SensorReading
		if err := rows.Scan(&r.ID, &r.SensorID, &r.Value, &r.Timestamp); err != nil {
			return nil, fmt.Errorf("scanning sensor reading: %w", err)
		}
		readings = append(readings, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating sensor readings: %w", err)
	}

	return readings, nil
}

// CheckCorruption runs PRAGMA integrity_check and returns an error if the database is corrupted.
func (s *Store) CheckCorruption() error {
	var result string
	err := s.db.QueryRow("PRAGMA integrity_check").Scan(&result)
	if err != nil {
		return fmt.Errorf("running integrity check: %w", err)
	}
	if result != "ok" {
		return fmt.Errorf("database corruption detected: %s", result)
	}
	return nil
}

// NewForTest opens a SQLite database in WAL mode without running migrations.
// Creates the sensor_readings table directly for test use.
func NewForTest(dbPath string) (*Store, error) {
	dsn := "file:" + dbPath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Verify WAL mode
	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		db.Close()
		return nil, fmt.Errorf("verifying WAL mode: %w", err)
	}
	if journalMode != "wal" {
		db.Close()
		return nil, fmt.Errorf("WAL mode not enabled: got %s", journalMode)
	}

	// Create table directly for tests (no migration files needed)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sensor_readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sensor_id TEXT NOT NULL,
			value REAL NOT NULL,
			timestamp DATETIME NOT NULL
		)
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("creating sensor_readings table: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// CreateAsset inserts a new asset and returns it with auto-generated ID and timestamps.
func (s *Store) CreateAsset(ctx context.Context, asset *Asset) (*Asset, error) {
	// Apply defaults for empty fields
	if asset.Site == "" {
		asset.Site = "factory-1"
	}
	if asset.Type == "" {
		asset.Type = "machine"
	}

	query, args, err := squirrel.Insert("assets").
		Columns("name", "site", "area", "line", "type").
		Values(asset.Name, asset.Site, asset.Area, asset.Line, asset.Type).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building insert asset query: %w", err)
	}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting asset: %w", err)
	}
	return asset, nil
}

// GetAsset retrieves an asset by ID.
func (s *Store) GetAsset(ctx context.Context, id int64) (*Asset, error) {
	query, args, err := squirrel.Select("id", "name", "site", "area", "line", "type", "created_at", "updated_at").
		From("assets").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select asset query: %w", err)
	}

	var asset Asset
	err = s.db.QueryRowContext(ctx, query, args...).Scan(
		&asset.ID, &asset.Name, &asset.Site, &asset.Area, &asset.Line, &asset.Type,
		&asset.CreatedAt, &asset.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("querying asset: %w", err)
	}
	return &asset, nil
}

// ListAssets returns assets with pagination.
func (s *Store) ListAssets(ctx context.Context, page, limit int) ([]Asset, int, error) {
	offset := (page - 1) * limit

	// Count total
	var total int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM assets").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("counting assets: %w", err)
	}

	// Fetch page
	query, args, err := squirrel.Select("id", "name", "site", "area", "line", "type", "created_at", "updated_at").
		From("assets").
		OrderBy("id ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("building select assets query: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("querying assets: %w", err)
	}
	defer rows.Close()

	var assets []Asset
	for rows.Next() {
		var a Asset
		if err := rows.Scan(&a.ID, &a.Name, &a.Site, &a.Area, &a.Line, &a.Type, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scanning asset: %w", err)
		}
		assets = append(assets, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating assets: %w", err)
	}

	return assets, total, nil
}

// CreateAssetSensor inserts a sensor linked to an asset.
func (s *Store) CreateAssetSensor(ctx context.Context, assetID int64, sensor *AssetSensor) (*AssetSensor, error) {
	query, args, err := squirrel.Insert("asset_sensors").
		Columns("asset_id", "sensor_id", "sensor_type", "topic").
		Values(assetID, sensor.SensorID, sensor.SensorType, sensor.Topic).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building insert asset sensor query: %w", err)
	}

	sensor.AssetID = assetID
	err = s.db.QueryRowContext(ctx, query, args...).Scan(&sensor.ID, &sensor.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting asset sensor: %w", err)
	}
	return sensor, nil
}

// GetAssetSensors returns sensors for a given asset.
func (s *Store) GetAssetSensors(ctx context.Context, assetID int64) ([]AssetSensor, error) {
	query, args, err := squirrel.Select("id", "asset_id", "sensor_id", "sensor_type", "topic", "created_at").
		From("asset_sensors").
		Where(squirrel.Eq{"asset_id": assetID}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("building select asset sensors query: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying asset sensors: %w", err)
	}
	defer rows.Close()

	var sensors []AssetSensor
	for rows.Next() {
		var s AssetSensor
		if err := rows.Scan(&s.ID, &s.AssetID, &s.SensorID, &s.SensorType, &s.Topic, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning asset sensor: %w", err)
		}
		sensors = append(sensors, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating asset sensors: %w", err)
	}

	return sensors, nil
}
