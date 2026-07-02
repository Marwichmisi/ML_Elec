# Phase 01: Core Foundation - Pattern Map

**Mapped:** 2026-06-30
**Files analyzed:** 18 (all new — greenfield project)
**Analogs found:** 0 / 18 (greenfield — use RESEARCH.md patterns as analogs)

## File Classification

| New/Modified File | Role | Data Flow | Closest Analog | Match Quality |
|-------------------|------|-----------|----------------|---------------|
| `cmd/ml-elec/main.go` | entry-point, orchestrator | lifecycle | RESEARCH Pattern 5: Graceful Shutdown | primary |
| `internal/config/config.go` | config, parser | file-I/O | RESEARCH Pattern 1-5: Config loading | primary |
| `internal/config/config_test.go` | test | unit | RESEARCH Validation Architecture | primary |
| `internal/nats/nats.go` | service, wrapper | pub-sub | RESEARCH Pattern 1: Embedded NATS | primary |
| `internal/nats/nats_test.go` | test | integration | RESEARCH Validation Architecture | primary |
| `internal/storage/storage.go` | repository, service | CRUD | RESEARCH Pattern 3: SQLite WAL | primary |
| `internal/storage/migrations/` | migration | DDL | RESEARCH Pattern 3: SQLite Migrations | primary |
| `internal/storage/storage_test.go` | test | unit/integration | RESEARCH Validation Architecture | primary |
| `internal/plugin/manager.go` | service, manager | event-driven | RESEARCH Pattern 2: Plugin Manager | primary |
| `internal/plugin/manager_test.go` | test | integration | RESEARCH Validation Architecture | primary |
| `internal/api/server.go` | server, middleware | request-response | RESEARCH: net/http standard | primary |
| `internal/api/health.go` | handler, controller | request-response | RESEARCH: REST endpoint pattern | primary |
| `internal/api/sensors.go` | handler, controller | request-response | RESEARCH: REST endpoint pattern | primary |
| `internal/api/api_test.go` | test | integration | RESEARCH Validation Architecture | primary |
| `wire.go` | DI-config | wiring | RESEARCH Pattern 4: Wire DI | primary |
| `wire_gen.go` | generated | wiring | RESEARCH Pattern 4: Wire Generated | primary |
| `config.yaml` | config | static | RESEARCH: YAML config | primary |
| `Makefile` | build | automation | Go standard build patterns | primary |

## Pattern Assignments

### `cmd/ml-elec/main.go` (entry-point, orchestrator, lifecycle)

**Analog:** RESEARCH Pattern 5: Graceful Shutdown with LIFO Order (lines 367-410)

**Imports pattern:**
```go
package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"
    "time"

    "ml-elec/internal/api"
    "ml-elec/internal/config"
    "ml-elec/internal/nats"
    "ml-elec/internal/plugin"
    "ml-elec/internal/storage"
)
```

**Signal handling pattern:**
```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    // Load config
    cfg, err := config.Load()
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    // Initialize components (DI order)
    // Start components
    // Wait for signal
    <-ctx.Done()

    // LIFO shutdown with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    apiServer.Shutdown(shutdownCtx)    // 1st: stop accepting new requests
    pluginMgr.ShutdownAll(shutdownCtx) // 2nd: stop plugins
    natsServer.Shutdown()              // 3rd: stop NATS
    db.Close()                         // 4th: close database
}
```

**Error handling pattern:**
```go
// Exit code 1 for startup failures
if err != nil {
    slog.Error("component failed", "component", "name", "error", err)
    os.Exit(1)
}
```

---

### `internal/config/config.go` (config, parser, file-I/O)

**Analog:** RESEARCH: YAML/TOML config loading (lines 58-59, 100-101)

**Imports pattern:**
```go
package config

import (
    "fmt"
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
    NATS    NATSConfig    `yaml:"nats"`
    Storage StorageConfig `yaml:"storage"`
    API     APIConfig     `yaml:"api"`
    Plugins PluginConfig  `yaml:"plugins"`
}
```

**Config loading pattern (search multiple paths):**
```go
func Load() (*Config, error) {
    paths := []string{
        "./config.yaml",
        filepath.Join(os.Getenv("HOME"), ".config/ml-elec/config.yaml"),
    }

    for _, path := range paths {
        if _, err := os.Stat(path); err == nil {
            return loadFromFile(path)
        }
    }

    // Return defaults
    return DefaultConfig(), nil
}

func loadFromFile(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config %s: %w", path, err)
    }

    cfg := DefaultConfig()
    if err := yaml.Unmarshal(data, cfg); err != nil {
        return nil, fmt.Errorf("parsing config %s: %w", path, err)
    }

    return cfg, nil
}

func DefaultConfig() *Config {
    return &Config{
        NATS: NATSConfig{
            Host: "127.0.0.1",
            Port: -1, // random port
        },
        Storage: StorageConfig{
            Path: "./data/sensors.db",
        },
        API: APIConfig{
            Port: 8080,
        },
        Plugins: PluginConfig{
            Enabled: []string{},
        },
    }
}
```

**Validation pattern:**
```go
func (c *Config) Validate() error {
    if c.NATS.Host == "" {
        return fmt.Errorf("nats.host is required")
    }
    if c.Storage.Path == "" {
        return fmt.Errorf("storage.path is required")
    }
    return nil
}
```

---

### `internal/nats/nats.go` (service, wrapper, pub-sub)

**Analog:** RESEARCH Pattern 1: Embedded NATS Server (lines 215-253)

**Imports pattern:**
```go
package nats

import (
    "errors"
    "fmt"
    "time"

    "github.com/nats-io/nats-server/v2/server"
)
```

**Server wrapper pattern:**
```go
// Server wraps the NATS server for lifecycle management
type Server struct {
    srv *server.Server
    opts *server.Options
}

func New(cfg *config.NATSConfig) (*Server, error) {
    opts := &server.Options{
        Host:     cfg.Host,
        Port:     cfg.Port,
        MaxConn:  1024,
        MaxPayload: 1048576,
    }

    s, err := server.NewServer(opts)
    if err != nil {
        return nil, fmt.Errorf("creating nats server: %w", err)
    }

    return &Server{srv: s, opts: opts}, nil
}

func (s *Server) Start() error {
    s.srv.ConfigureLogger()
    s.srv.Start()

    if !s.srv.ReadyForConnections(10 * time.Second) {
        return errors.New("nats server not ready")
    }

    return nil
}

func (s *Server) Shutdown() {
    s.srv.Shutdown()
    s.srv.WaitForShutdown()
}

func (s *Server) Client() (*nats.Conn, error) {
    return nats.Connect(s.srv.ClientURL())
}
```

**Error handling pattern:**
```go
// Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("creating nats server: %w", err)
}
```

---

### `internal/storage/storage.go` (repository, service, CRUD)

**Analog:** RESEARCH Pattern 3: SQLite WAL with Migrations (lines 297-332)

**Imports pattern:**
```go
package storage

import (
    "database/sql"
    "fmt"
    "time"

    "github.com/Masterminds/squirrel"
    _ "modernc.org/sqlite"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/sqlite"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)
```

**SQLite WAL connection pattern:**
```go
// Store provides sensor data storage
type Store struct {
    db *sql.DB
}

func New(cfg *config.StorageConfig) (*Store, error) {
    dsn := "file:" + cfg.Path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"
    db, err := sql.Open("sqlite", dsn)
    if err != nil {
        return nil, fmt.Errorf("opening database: %w", err)
    }

    // Verify WAL mode
    var journalMode string
    if err := db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil || journalMode != "wal" {
        db.Close()
        return nil, fmt.Errorf("WAL mode not enabled: %w", err)
    }

    // Run migrations
    if err := runMigrations(cfg.Path); err != nil {
        db.Close()
        return nil, fmt.Errorf("running migrations: %w", err)
    }

    return &Store{db: db}, nil
}
```

**Migrations pattern:**
```go
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
```

**Query builder pattern (squirrel):**
```go
// InsertSensor inserts a sensor reading
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

// GetSensors retrieves sensor readings
func (s *Store) GetSensors(ctx context.Context, sensorID string, limit int) ([]SensorReading, error) {
    query, args, err := squirrel.Select("sensor_id", "value", "timestamp").
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

    // Scan results...
    return readings, nil
}
```

**Corruption check pattern:**
```go
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
```

---

### `internal/plugin/manager.go` (service, manager, event-driven)

**Analog:** RESEARCH Pattern 2: Plugin Manager with Crash Isolation (lines 255-295)

**Imports pattern:**
```go
package plugin

import (
    "fmt"
    "os/exec"
    "sync"

    "github.com/hashicorp/go-plugin"
)
```

**Plugin manager pattern:**
```go
// Manager handles plugin lifecycle
type Manager struct {
    clients map[string]*plugin.Client
    mu      sync.RWMutex
}

func NewManager() *Manager {
    return &Manager{
        clients: make(map[string]*plugin.Client),
    }
}

// Launch starts a plugin as a child process
func (m *Manager) Launch(name, path string, handshake plugin.HandshakeConfig, plugins map[string]plugin.Plugin) error {
    client := plugin.NewClient(&plugin.ClientConfig{
        HandshakeConfig: handshake,
        Plugins:         plugins,
        Cmd:             exec.Command(path),
        Managed:         true,
    })

    rpcClient, err := client.Client()
    if err != nil {
        client.Kill()
        return fmt.Errorf("connecting to plugin %s: %w", name, err)
    }

    m.mu.Lock()
    m.clients[name] = client
    m.mu.Unlock()

    return nil
}

// ShutdownAll terminates all running plugins
func (m *Manager) ShutdownAll() {
    m.mu.RLock()
    defer m.mu.RUnlock()

    for name, client := range m.clients {
        client.Kill()
        delete(m.clients, name)
    }
}
```

**Handshake configuration pattern:**
```go
var HandshakeConfig = plugin.HandshakeConfig{
    ProtocolVersion:  1,
    MagicCookieKey:   "ML_ELEC_PLUGIN",
    MagicCookieValue: "ml-elec-v1",
}
```

**Crash isolation pattern:**
```go
// Plugins are isolated via child processes
// A plugin crash does NOT crash the core
// client.Kill() cleans up the process
```

---

### `internal/api/server.go` (server, middleware, request-response)

**Analog:** RESEARCH: net/http standard + swaggo (lines 82-101)

**Imports pattern:**
```go
package api

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "time"

    "github.com/rs/cors"
    "ml-elec/internal/storage"
)
```

**Server pattern:**
```go
// Server holds the HTTP server and dependencies
type Server struct {
    server *http.Server
    store  *storage.Store
}

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
    mux.HandleFunc("GET /api/v1/sensors", s.GetSensorsHandler)

    return s
}

func (s *Server) Start() error {
    slog.Info("starting API server", "addr", s.server.Addr)
    return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.server.Shutdown(ctx)
}
```

**CORS middleware pattern:**
```go
// Apply CORS middleware
c := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders: []string{"Content-Type", "Authorization"},
})
s.server.Handler = c.Handler(mux)
```

---

### `internal/api/health.go` (handler, controller, request-response)

**Analog:** RESEARCH: REST endpoint pattern (lines 174-178)

**Imports pattern:**
```go
package api

import (
    "encoding/json"
    "net/http"
)

// HealthResponse represents the health check response
type HealthResponse struct {
    Status string `json:"status"`
}
```

**Health handler pattern:**
```go
// HealthHandler returns 200 OK
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
```

---

### `internal/api/sensors.go` (handler, controller, request-response)

**Analog:** RESEARCH: REST endpoint pattern (lines 174-178)

**Imports pattern:**
```go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"

    "ml-elec/internal/storage"
)
```

**Sensor query handler pattern:**
```go
// GetSensorsHandler returns sensor readings
// @Summary Get sensor readings
// @Tags sensors
// @Produce json
// @Param sensor_id query string true "Sensor ID"
// @Param limit query int false "Limit" default(100)
// @Success 200 {object} storage.SensorReading
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
        if parsed, err := strconv.Atoi(l); err == nil {
            limit = parsed
        }
    }

    readings, err := s.store.GetSensors(r.Context(), sensorID, limit)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to query sensors")
        return
    }

    writeJSON(w, http.StatusOK, map[string]interface{}{
        "data": readings,
    })
}
```

**Response formatting pattern (D-10):**
```go
// Standard JSON response structure: {"data": ..., "error": ...}
type ErrorResponse struct {
    Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, ErrorResponse{Error: message})
}
```

---

### `wire.go` (DI-config, wiring)

**Analog:** RESEARCH Pattern 4: Wire Dependency Injection (lines 334-365)

**Imports pattern:**
```go
//go:build wireinject

package main

import (
    "github.com/google/wire"
    "ml-elec/internal/api"
    "ml-elec/internal/config"
    "ml-elec/internal/nats"
    "ml-elec/internal/plugin"
    "ml-elec/internal/storage"
)
```

**Wire injector pattern:**
```go
type App struct {
    Config     *config.Config
    NATSServer *nats.Server
    Store      *storage.Store
    PluginMgr  *plugin.Manager
    APIServer  *api.Server
}

func InitializeApp() (*App, error) {
    wire.Build(
        config.Load,
        nats.New,
        storage.New,
        plugin.NewManager,
        api.NewServer,
        wire.Struct(new(App), "*"),
    )
    return nil, nil
}
```

---

### `config.yaml` (config, static)

**Analog:** RESEARCH: YAML config (lines 58-59)

**Config file pattern:**
```yaml
# ML-Elec Core Configuration
# Docs: https://github.com/ml-elec/ml-elec

nats:
  host: 127.0.0.1
  port: -1  # random port for embedded mode

storage:
  path: ./data/sensors.db

api:
  port: 8080

plugins:
  enabled: []
  # - mqtt-sensor
  # - modbus-gateway
```

---

### `Makefile` (build, automation)

**Analog:** Go standard build patterns

**Makefile pattern:**
```makefile
.PHONY: build test run clean wire lint

# Build the binary
build:
	go build -o bin/ml-elec ./cmd/ml-elec

# Run tests
test:
	go test ./...

# Run tests with race detector
test-race:
	go test -race ./...

# Run tests with coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Generate Wire code
wire:
	wire ./...

# Run linter
lint:
	golangci-lint run

# Run the application
run: build
	./bin/ml-elec

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out

# Check LOC constraint (must be < 5000)
check-loc:
	@echo "Checking LOC constraint (< 5000)..."
	@LOC=$$(find . -name '*.go' ! -name '*_test.go' ! -path './generated/*' | xargs wc -l | tail -1 | awk '{print $$1}'); \
	if [ $$LOC -lt 5000 ]; then \
		echo "✓ LOC: $$LOC (< 5000)"; \
	else \
		echo "✗ LOC: $$LOC (exceeds 5000 limit)"; \
		exit 1; \
	fi
```

---

## Shared Patterns

### Error Handling (Go idiomatic)
**Source:** RESEARCH: Anti-Patterns to Avoid (lines 412-418)
**Apply to:** All packages

```go
// Wrap errors with context using fmt.Errorf("context: %w", err)
if err != nil {
    return fmt.Errorf("creating nats server: %w", err)
}

// Use slog for structured logging
slog.Error("component failed", "component", "name", "error", err)

// Exit with code 1 for startup failures
os.Exit(1)
```

### Context Propagation
**Source:** RESEARCH: Graceful Shutdown (lines 367-410)
**Apply to:** All handlers and services

```go
// Pass context to all functions
func (s *Store) GetSensors(ctx context.Context, sensorID string, limit int) ([]SensorReading, error) {
    // Use ctx for database queries
    rows, err := s.db.QueryContext(ctx, query, args...)
}

// HTTP handlers use r.Context()
func (s *Server) GetSensorsHandler(w http.ResponseWriter, r *http.Request) {
    readings, err := s.store.GetSensors(r.Context(), sensorID, limit)
}
```

### Testing Patterns
**Source:** RESEARCH: Validation Architecture (lines 612-644)
**Apply to:** All test files

```go
// Unit tests - test individual functions
func TestConfigLoad(t *testing.T) {
    // Arrange
    // Act
    // Assert
}

// Integration tests - test component interaction
func TestNATSPubSub(t *testing.T) {
    // Setup NATS server
    // Test publish/subscribe
    // Cleanup
}

// Use table-driven tests
func TestSensorReading(t *testing.T) {
    tests := []struct {
        name     string
        sensorID string
        value    float64
        wantErr  bool
    }{
        {"valid reading", "sensor-1", 23.5, false},
        {"empty sensor ID", "", 23.5, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

### Shutdown Pattern (LIFO)
**Source:** RESEARCH: Graceful Shutdown (lines 367-410)
**Apply to:** `cmd/ml-elec/main.go`

```go
// Initialize components in order
cfg := config.Load()
natsServer := nats.New(cfg.NATS)
db := storage.New(cfg.Storage)
pluginMgr := plugin.NewManager()
apiServer := api.NewServer(cfg.API, db)

// Start components
natsServer.Start()
pluginMgr.StartAll()
apiServer.Start()

// Wait for signal
<-ctx.Done()

// Shutdown in LIFO order with timeout
shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

apiServer.Shutdown(shutdownCtx)    // 1st
pluginMgr.ShutdownAll()            // 2nd
natsServer.Shutdown()              // 3rd
db.Close()                         // 4th
```

---

## No Analog Found

Files with no close match in the codebase (planner should use RESEARCH.md patterns instead):

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| All files | — | — | Greenfield project — no existing codebase |

## Metadata

**Analog search scope:** Entire workspace
**Files scanned:** 0 (greenfield)
**Pattern extraction date:** 2026-06-30

**Note:** This is a greenfield project. All patterns are sourced from RESEARCH.md which contains verified patterns from official documentation and Context7 fetches. The planner should use these patterns as the primary implementation guide.
