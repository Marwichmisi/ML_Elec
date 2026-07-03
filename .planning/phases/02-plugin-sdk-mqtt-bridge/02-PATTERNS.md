# Phase 02: Plugin SDK & MQTT Bridge - Pattern Map

**Mapped:** 2026-07-03
**Files analyzed:** 16
**Analogs found:** 12 / 16

## File Classification

| New/Modified File | Role | Data Flow | Closest Analog | Match Quality |
|-------------------|------|-----------|----------------|---------------|
| `cmd/mqtt-plugin/main.go` | main binary | event-driven | `cmd/ml-elec/main.go` | exact |
| `cmd/mqtt-plugin/main_test.go` | test | integration | `internal/api/api_test.go` | role-match |
| `pkg/sdk/v1/proto/lifecycle.proto` | config/contract | request-response | *(no analog — new pattern)* | — |
| `pkg/sdk/v1/proto/sensor.proto` | config/contract | request-response | *(no analog — new pattern)* | — |
| `pkg/sdk/v1/plugin.go` | service/middleware | request-response | `internal/plugin/manager.go` | role-match |
| `pkg/sdk/v1/plugin_test.go` | test | unit | `internal/plugin/manager_test.go` | exact |
| `internal/storage/migrations/002_assets.up.sql` | migration | CRUD | `internal/storage/migrations/001_init.up.sql` | exact |
| `internal/storage/migrations/002_assets.down.sql` | migration | CRUD | `internal/storage/migrations/001_init.down.sql` | exact |
| `internal/api/assets.go` | controller | CRUD | `internal/api/sensors.go` | exact |
| `internal/api/assets_test.go` | test | unit | `internal/api/api_test.go` | exact |
| `internal/config/config.go` *(modified)* | config | transform | — (self) | — |
| `internal/storage/storage.go` *(modified)* | service | CRUD | — (self) | — |
| `internal/plugin/manager.go` *(modified)* | service | request-response | — (self) | — |
| `cmd/mock-plugin/main.go` *(modified)* | main binary | request-response | — (self) | — |
| `wire.go / wire_gen.go` *(modified)* | config/wiring | transform | — (self) | — |
| `Makefile` *(modified)* | config | transform | — (self) | — |

## Pattern Assignments

### `cmd/mqtt-plugin/main.go` (main binary, event-driven)

**Analog:** `cmd/ml-elec/main.go`

**Imports pattern** (lines 1-17):
```go
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-hclog"

	"ml-elec/internal/config"
	"ml-elec/internal/plugin"
)
```

**App struct + InitializeApp pattern** (lines 20-50):
```go
// App holds all MQTT plugin components.
type App struct {
	Config *config.Config
	Broker *mqtt.Server   // mochi-mqtt embedded broker
	Client mqtt.Client     // paho.mqtt.golang client
}

func InitializeApp(cfg *config.Config) (*App, error) {
	// 1. Start embedded mochi-mqtt broker on :1883
	// 2. Connect paho.mqtt.golang client to local broker
	// 3. Subscribe to esp32/# (wildcard)
	// 4. On message: parse JSON/binary → validate → publish to NATS
	return &App{...}, nil
}
```

**Signal handling + LIFO shutdown pattern** (lines 52-139):
```go
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	app, err := InitializeApp(cfg)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	// Start broker + client + NATS bridge
	// ...

	<-ctx.Done()
	slog.Info("shutdown signal received")

	// LIFO: MQTT client → Broker → NATS
	// ...
}
```

**go-plugin serve pattern** — for gRPC registration from `cmd/mock-plugin/main.go` (lines 25-31):
```go
goplugin.Serve(&goplugin.ServeConfig{
	HandshakeConfig: plugin.HandshakeConfig,
	Plugins: map[string]goplugin.Plugin{
		"sensor": &plugin.SensorPluginGRPC{Impl: &MQTTPlugin{}},
	},
	Logger: logger,
})
```

---

### `pkg/sdk/v1/plugin.go` (service/middleware, request-response)

**Analog:** `internal/plugin/manager.go`

**Imports pattern** (lines 3-10):
```go
package sdk

import (
	"context"
	"fmt"

	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	pb "ml-elec/pkg/sdk/v1"
)
```

**GRPCPlugin struct + interface pattern** — adapted from `internal/plugin/manager.go` lines 12-79:
```go
// PluginLifecycle is the Go interface for plugin lifecycle management.
type PluginLifecycle interface {
	Init(ctx context.Context, config map[string]string) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// SensorCollector is the Go interface for sensor data collection.
type SensorCollector interface {
	Collect(ctx context.Context, request *pb.CollectRequest) (*pb.CollectResponse, error)
}

// GRPCPlugin implements plugin.GRPCPlugin for go-plugin integration.
type GRPCPlugin struct {
	goplugin.NetRPCUnsupportedPlugin // Disable net/rpc, gRPC only
	Impl        PluginLifecycle
	CollectImpl SensorCollector
}

func (p *GRPCPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterPluginLifecycleServer(s, &GRPCServer{Impl: p.Impl})
	pb.RegisterSensorCollectorServer(s, &GRPCServer{CollectImpl: p.CollectImpl})
	return nil
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		lifecycle: pb.NewPluginLifecycleClient(c),
		collector: pb.NewSensorCollectorClient(c),
	}, nil
}
```

**Error handling pattern** — wrapping errors with context (from `internal/plugin/manager.go` lines 110-121):
```go
func (c *GRPCClient) Init(ctx context.Context, config map[string]string) error {
	_, err := c.lifecycle.Init(ctx, &pb.InitRequest{Config: config})
	if err != nil {
		return fmt.Errorf("grpc init: %w", err)
	}
	return nil
}
```

---

### `internal/api/assets.go` (controller, CRUD)

**Analog:** `internal/api/sensors.go`

**Imports pattern** (lines 1-9):
```go
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"ml-elec/internal/storage"
)
```

**Pagination response pattern** (from D-23, D-24, D-25):
```go
// PaginationResponse wraps paginated results.
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
		Total int `json:"total"`
	} `json:"pagination"`
}
```

**CRUD handler pattern** — adapted from `internal/api/sensors.go` lines 26-54:
```go
// CreateAssetHandler creates a new asset.
//
// @Summary Create asset
// @Tags assets
// @Accept json
// @Produce json
// @Param asset body AssetRequest true "Asset data"
// @Success 201 {object} Asset
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/assets [post]
func (s *Server) CreateAssetHandler(w http.ResponseWriter, r *http.Request) {
	var req AssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	asset, err := s.store.CreateAsset(r.Context(), &req)
	if err != nil {
		if isConflict(err) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create asset")
		return
	}

	writeJSON(w, http.StatusCreated, asset)
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
```

**Error helper pattern** — from `internal/api/sensors.go` lines 57-69:
```go
// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		_ = err
	}
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}
```

---

### `internal/storage/migrations/002_assets.up.sql` (migration, CRUD)

**Analog:** `internal/storage/migrations/001_init.up.sql`

**Migration pattern** (from `001_init.up.sql` lines 1-9):
```sql
CREATE TABLE IF NOT EXISTS assets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    site TEXT NOT NULL DEFAULT 'factory-1',
    area TEXT NOT NULL DEFAULT '',
    line TEXT NOT NULL DEFAULT '',
    type TEXT NOT NULL DEFAULT 'machine',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_assets_site ON assets(site);
CREATE INDEX IF NOT EXISTS idx_assets_name ON assets(name);

CREATE TABLE IF NOT EXISTS asset_sensors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    asset_id INTEGER NOT NULL,
    sensor_id TEXT NOT NULL,
    sensor_type TEXT NOT NULL DEFAULT 'generic',
    topic TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(asset_id, sensor_id)
);

CREATE INDEX IF NOT EXISTS idx_asset_sensors_asset_id ON asset_sensors(asset_id);
CREATE INDEX IF NOT EXISTS idx_asset_sensors_sensor_id ON asset_sensors(sensor_id);
```

---

### `internal/storage/migrations/002_assets.down.sql` (migration, CRUD)

**Analog:** `internal/storage/migrations/001_init.down.sql`

**Down migration pattern:**
```sql
DROP TABLE IF EXISTS asset_sensors;
DROP TABLE IF EXISTS assets;
```

---

### `internal/storage/storage.go` *(modified — add asset CRUD)* (service, CRUD)

**Existing analog:** `internal/storage/storage.go` (self)

**InsertAsset method pattern** — adapted from lines 96-111 (InsertSensor):
```go
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

// CreateAsset inserts a new asset.
func (s *Store) CreateAsset(ctx context.Context, asset *Asset) (*Asset, error) {
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
```

**ListAssets method pattern** — adapted from lines 114-145 (GetSensors):
```go
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
```

---

### `internal/plugin/manager.go` *(modified — gRPC transport)* (service, request-response)

**Existing analog:** `internal/plugin/manager.go` (self)

**gRPC transport addition pattern** — add alongside existing net/rpc support:
```go
// GRPCPlugin wraps a go-plugin GRPCPlugin implementation.
type GRPCPlugin struct {
	goplugin.NetRPCUnsupportedPlugin
	Impl interface{} // PluginLifecycle or SensorCollector
}

// LaunchGRPC starts a plugin with gRPC transport.
func (m *Manager) LaunchGRPC(name, path string, enabledPlugins []string, grpcPlugin goplugin.Plugin) error {
	if !isPluginEnabled(name, enabledPlugins) {
		return fmt.Errorf("plugin %q is not enabled", name)
	}

	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig:  HandshakeConfig,
		Plugins:          map[string]goplugin.Plugin{"sensor": grpcPlugin},
		Cmd:              exec.Command(path),
		Managed:          true,
		AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolGRPC},
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return fmt.Errorf("connecting to grpc plugin %q: %w", name, err)
	}

	_, err = rpcClient.Dispense("sensor")
	if err != nil {
		client.Kill()
		return fmt.Errorf("dispensing grpc plugin %q: %w", name, err)
	}

	m.mu.Lock()
	m.clients[name] = client
	m.mu.Unlock()

	return nil
}
```

---

### `internal/config/config.go` *(modified — add mqtt, validation, assets)* (config)

**Existing analog:** `internal/config/config.go` (self)

**Config struct extension pattern** — from lines 12-38:
```go
// Config holds all configuration for the application.
type Config struct {
	NATS       NATSConfig       `yaml:"nats"`
	Storage    StorageConfig    `yaml:"storage"`
	API        APIConfig        `yaml:"api"`
	Plugins    PluginsConfig    `yaml:"plugins"`
	MQTT       MQTTConfig       `yaml:"mqtt"`
	Validation ValidationConfig `yaml:"validation"`
	Assets     AssetsConfig     `yaml:"assets"`
}

// MQTTConfig holds MQTT broker and client configuration.
type MQTTConfig struct {
	Port         int    `yaml:"port"`
	ClientID     string `yaml:"client_id"`
	CleanSession bool   `yaml:"clean_session"`
	KeepAlive    int    `yaml:"keep_alive"`
	Topics       struct {
		Subscribe string `yaml:"subscribe"`
		Status    string `yaml:"status"`
	} `yaml:"topics"`
}

// ValidationConfig holds 3-level data validation configuration.
type ValidationConfig struct {
	Ranges    map[string]RangeConfig `yaml:"ranges"`
	Timestamp TimestampConfig        `yaml:"timestamp"`
	Health    HealthConfig           `yaml:"health"`
}

type RangeConfig struct {
	Min float64 `yaml:"min"`
	Max float64 `yaml:"max"`
}

type TimestampConfig struct {
	MaxFutureDrift string `yaml:"max_future_drift"`
	MaxPastDrift   string `yaml:"max_past_drift"`
}

type HealthConfig struct {
	MaxPayloadSize int `yaml:"max_payload_size"`
	MinPayloadSize int `yaml:"min_payload_size"`
}

// AssetsConfig holds asset registry configuration.
type AssetsConfig struct {
	AutoRegister bool   `yaml:"auto_register"`
	DefaultSite  string `yaml:"default_site"`
}
```

**Default config pattern** — from lines 41-57:
```go
func DefaultConfig() *Config {
	return &Config{
		// ... existing defaults ...
		MQTT: MQTTConfig{
			Port:         1883,
			ClientID:     "mqtt-plugin",
			CleanSession: false,
			KeepAlive:    30,
		},
		Validation: ValidationConfig{
			Timestamp: TimestampConfig{
				MaxFutureDrift: "5s",
				MaxPastDrift:   "24h",
			},
			Health: HealthConfig{
				MaxPayloadSize: 1024,
				MinPayloadSize: 1,
			},
		},
		Assets: AssetsConfig{
			AutoRegister: true,
			DefaultSite:  "factory-1",
		},
	}
}
```

---

### `cmd/mock-plugin/main.go` *(modified — use new gRPC SDK)* (main binary, request-response)

**Existing analog:** `cmd/mock-plugin/main.go` (self)

**Updated serve pattern** — from lines 18-31:
```go
package main

import (
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-hclog"

	"ml-elec/internal/plugin"
	sdk "ml-elec/pkg/sdk/v1"
)

// MockSensorPlugin implements both PluginLifecycle and SensorCollector.
type MockSensorPlugin struct{}

func (p *MockSensorPlugin) Init(ctx context.Context, config map[string]string) error { return nil }
func (p *MockSensorPlugin) Start(ctx context.Context) error { return nil }
func (p *MockSensorPlugin) Stop(ctx context.Context) error  { return nil }
func (p *MockSensorPlugin) Collect(ctx context.Context, req *sdk.CollectRequest) (*sdk.CollectResponse, error) {
	return &sdk.CollectResponse{Accepted: true}, nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "mock-plugin",
		Level:  hclog.Info,
		Output: hclog.DefaultOutput,
	})

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			"sensor": &sdk.GRPCPlugin{Impl: &MockSensorPlugin{}, CollectImpl: &MockSensorPlugin{}},
		},
		Logger: logger,
	})
}
```

---

### `pkg/sdk/v1/proto/lifecycle.proto` (config/contract)

**No existing analog** — first proto file in project.

**Proto pattern** from RESEARCH.md (lines 593-617):
```protobuf
syntax = "proto3";
package sdk.v1;
option go_package = "ml-elec/pkg/sdk/v1";

service PluginLifecycle {
    rpc Init(InitRequest) returns (InitResponse);
    rpc Start(StartRequest) returns (StartResponse);
    rpc Stop(StopRequest) returns (StopResponse);
}

message InitRequest {
    map<string, string> config = 1;
}
message InitResponse {}
message StartRequest {}
message StartResponse {}
message StopRequest {}
message StopResponse {}
```

---

### `pkg/sdk/v1/proto/sensor.proto` (config/contract)

**No existing analog** — first proto file in project.

**Proto pattern** from RESEARCH.md (lines 619-641):
```protobuf
syntax = "proto3";
package sdk.v1;
option go_package = "ml-elec/pkg/sdk/v1";

service SensorCollector {
    rpc Collect(CollectRequest) returns (CollectResponse);
}

message CollectRequest {
    string sensor_id = 1;
    string topic = 2;
    bytes payload = 3;
    string payload_format = 4; // "json" or "binary"
}

message CollectResponse {
    bool accepted = 1;
    string rejection_reason = 2;
}
```

---

### `wire.go / wire_gen.go` *(modified — new providers)* (config/wiring)

**Existing analog:** `cmd/ml-elec/main.go` lines 29-49 (manual DI pattern)

Note: Le projet utilise actuellement du wiring manuel (pas wire.Generate). L'extension suit le même pattern.

```go
// New providers for Phase 02:
// - config.MQTTConfig → MQTT plugin config
// - config.ValidationConfig → validation config
// - config.AssetsConfig → assets config
```

---

### `Makefile` *(modified — add proto target)* (config)

**Existing analog:** `Makefile`

**Proto target pattern** from RESEARCH.md (lines 710-718):
```makefile
# Add after existing targets:
.PHONY: proto
proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       pkg/sdk/v1/proto/lifecycle.proto \
	       pkg/sdk/v1/proto/sensor.proto

# Update build target to include mqtt-plugin:
build: proto
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) ./cmd/ml-elec
	go build -o bin/mqtt-plugin ./cmd/mqtt-plugin
```

---

### `config.yaml` *(modified — add mqtt, validation, assets)* (config)

**Existing analog:** `config.yaml`

**Config extension pattern** from RESEARCH.md (lines 675-707):
```yaml
# Add after existing sections:
mqtt:
  port: 1883
  client_id: "mqtt-plugin"
  clean_session: false
  keep_alive: 30
  topics:
    subscribe: "esp32/#"
    status: "sys/mqtt-plugin/status"

validation:
  ranges:
    temperature:
      min: -40.0
      max: 150.0
    humidity:
      min: 0.0
      max: 100.0
    current:
      min: 0.0
      max: 1000.0
  timestamp:
    max_future_drift: 5s
    max_past_drift: 24h
  health:
    max_payload_size: 1024
    min_payload_size: 1

assets:
  auto_register: true
  default_site: "factory-1"
```

---

### Test Files

#### `cmd/mqtt-plugin/main_test.go` (test, integration)

**Analog:** `internal/api/api_test.go`

**Integration test pattern** — from `api_test.go` lines 16-36 + 38-70:
```go
package main

import (
	"context"
	"testing"
	"time"

	"ml-elec/internal/config"
)

func TestMQTTPluginStarts(t *testing.T) {
	cfg := &config.MQTTConfig{
		Port:         1883,
		ClientID:     "test-plugin",
		CleanSession: true,
	}

	// Start broker, connect client, verify subscription
	// ...
}

func TestMQTTMessageBridge(t *testing.T) {
	// Start broker, publish test message, verify NATS receive
	// ...
}
```

#### `internal/api/assets_test.go` (test, unit)

**Analog:** `internal/api/api_test.go`

**Test pattern** — from `api_test.go` lines 148-217:
```go
func TestCreateAssetValidRequest(t *testing.T) {
	store, err := storage.NewForTest(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	cfg := &config.APIConfig{Port: 0}
	srv := NewServer(cfg, store)

	ts := httptest.NewServer(srv.server.Handler)
	defer ts.Close()

	body := `{"name":"CNC-01","site":"factory-1","area":"production","line":"L1","type":"cnc"}`
	resp, err := http.Post(ts.URL+"/api/v1/assets", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST /api/v1/assets failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}
}

func TestListAssetsPagination(t *testing.T) {
	// Insert multiple assets, verify pagination response
	// ...
}
```

#### `pkg/sdk/v1/plugin_test.go` (test, unit)

**Analog:** `internal/plugin/manager_test.go`

**gRPC test pattern:**
```go
package sdk

import (
	"context"
	"testing"
)

func TestGRPCPluginLifecycle(t *testing.T) {
	// Create in-memory gRPC server/client
	// Test Init/Start/Stop cycle
	// ...
}
```

---

## Shared Patterns

### Structured Logging
**Source:** `cmd/ml-elec/main.go` line 54
**Apply to:** All new files
```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))
// Usage: slog.Info("message", "key", value)
//        slog.Error("message", "error", err)
```

### Error Handling
**Source:** `internal/storage/storage.go` (throughout)
**Apply to:** All service and storage files
```go
// Wrap errors with context
return fmt.Errorf("building insert query: %w", err)
return fmt.Errorf("inserting sensor reading: %w", err)
```

### REST Response Format
**Source:** `internal/api/sensors.go` lines 12-14, 57-69
**Apply to:** All controller files
```go
type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) { ... }
func writeError(w http.ResponseWriter, status int, message string) { ... }
```

### SQLite Migration Pattern
**Source:** `internal/storage/storage.go` lines 70-93
**Apply to:** All migration files
```go
//go:embed migrations/*.sql
var migrationsFS embed.FS

func runMigrations(db *sql.DB, dbPath string) error {
	d, _ := iofs.New(migrationsFS, "migrations")
	driver, _ := sqlite.WithInstance(db, &sqlite.Config{})
	m, _ := migrate.NewWithInstance("iofs", d, "sqlite", driver)
	m.Up()
	return nil
}
```

### go-plugin Handshake
**Source:** `internal/plugin/manager.go` lines 18-22
**Apply to:** All plugin binaries
```go
var HandshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ML_ELEC_PLUGIN",
	MagicCookieValue: "ml-elec-v1",
}
```

### Test Helpers
**Source:** `internal/storage/storage.go` lines 162-195 (NewForTest)
**Apply to:** All test files
```go
// Use t.TempDir() for test databases
store, err := storage.NewForTest(t.TempDir() + "/test.db")
defer store.Close()

// Use httptest.NewServer for API tests
ts := httptest.NewServer(srv.server.Handler)
defer ts.Close()
```

---

## No Analog Found

Files with no close match in the codebase (planner should use RESEARCH.md patterns instead):

| File | Role | Data Flow | Reason |
|------|------|-----------|--------|
| `pkg/sdk/v1/proto/lifecycle.proto` | config/contract | request-response | First .proto file in project — no existing proto pattern |
| `pkg/sdk/v1/proto/sensor.proto` | config/contract | request-response | First .proto file in project — no existing proto pattern |
| `internal/storage/migrations/002_assets.down.sql` | migration | CRUD | No existing down migration to reference (only 001_init exists) |

## Metadata

**Analog search scope:** `cmd/`, `internal/`, `pkg/`, root config files
**Files scanned:** 16 source files + 5 test files
**Pattern extraction date:** 2026-07-03
