package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

// SensorReading represents a sensor data point.
type SensorReading struct {
	ID        int64     `json:"id"`
	SensorID  string    `json:"sensor_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// Store provides sensor data storage with SQLite WAL mode.
type Store struct {
	db *sql.DB
}

// New opens a SQLite database in WAL mode, runs migrations, and returns a Store.
func New(dbPath string) (*Store, error) {
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

	// Run migrations
	if err := runMigrations(dbPath); err != nil {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return &Store{db: db}, nil
}

func runMigrations(dbPath string) error {
	m, err := migrate.New("file://migrations", "sqlite://"+dbPath)
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}
	defer m.Close()

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

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
