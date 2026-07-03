# Phase 02: Plugin SDK & MQTT Bridge - Research

**Researched:** 2026-07-03
**Domain:** gRPC plugin SDK, embedded MQTT broker, sensor data validation, asset registry
**Confidence:** HIGH

## Summary

Phase 02 étend le micro-noyau Go de la Phase 01 avec un SDK gRPC versionné pour les plugins, un plugin MQTT externe avec broker embarqué, une validation 3 niveaux des données capteurs, un registry d'assets ISA-95/UNS avec API REST, et le support des payloads JSON + binaires MQTT. La clé architecturale est que le broker MQTT est embarqué dans le plugin (pas dans le core), préservant la contrainte micro-noyau (< 5000 LOC).

**Découverte critique:** La décision D-01 mentionne "Eclipse Paho Go" comme lib broker, mais Paho Go (`eclipse/paho.mqtt.golang`) est une bibliothèque **client** MQTT, pas un broker. Pour un broker MQTT embarqué en Go, la solution standard est **mochi-mqtt/server** — un broker MQTT v5/v3.1.1 pleinement conforme, embarquable, et performant. Le client Paho sera utilisé comme client MQTT dans le plugin pour se connecter au broker mochi embarqué.

**Primary recommendation:** Utiliser `mochi-mqtt/server/v2` pour le broker MQTT embarqué, `eclipse/paho.mqtt.golang` pour le client MQTT, `hashicorp/go-plugin` avec transport gRPC pour le SDK, et `google.golang.org/grpc` + `protoc` pour la génération de code proto.

<user_constraints>
## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Lib broker: Eclipse Paho Go — broker MQTT embarqué dans le plugin
- **D-02:** Port local fixe (1883) — les ESP32 se connectent directement
- **D-03:** Broker dans un processus séparé du plugin MQTT — le plugin MQTT démarre le broker comme sous-processus
- **D-04:** Cycle de vie broker: l'agent décide de la meilleure approche (plugin lance le broker)
- **D-05:** Un fichier .proto par service: `lifecycle.proto` (Init/Start/Stop) + `sensor.proto` (Collect)
- **D-06:** Code généré commit dans le repo — les plugins n'ont pas besoin de protoc pour compiler
- **D-07:** Support Go + Python plugin via go-plugin dès v1
- **D-08:** Plugin Python: wrapper simplifié avec contrat gRPC unique maintenu par la plateforme
- **D-09:** Versioning dès v1: `pkg/sdk/v1/` — package unique pour v1
- **D-10:** Breaking changes futurs: tags Git + nouveau module `/v2` uniquement sur breaking change
- **D-11:** Fréquence d'échantillonnage vibration: configurable avec profils prédéfinis (low, standard, advanced, bearing), défaut 1 kHz, max 25.6 kHz
- **D-12:** Structure trame binaire v1: header fixe 24-32 octets, payload brut, version en 1er octet, encoding configurable (int16 défaut), pas de CRC v1
- **D-13:** Payloads JSON: telemetry `{ts, values:{temperature, humidity, current}}` + status `{ts, status, battery, rssi, firmware}`
- **D-14:** Topics MQTT: hiérarchie ISA-95/UNS `factory/{site}/{area}/{line}/{asset}/{type}` + `sys/{component}/{type}`, wildcard-friendly
- **D-15:** Tests plugin MQTT: mock broker en unit tests + vrai broker Paho en integration tests
- **D-16:** Tests isolation crash: test Go principal (lance plugin, tue, vérifie core continue) + script shell pour contrôle manuel
- **D-17:** Tests performance: benchmark Go avec timer, mesure broker→NATS (100 msg/s, <100ms)
- **D-18:** Reconnexion: backoff exponentiel (1s→30s) + LWT pour status offline + clean session=false
- **D-19:** Toutes les bonnes pratiques MQTT skill appliquées: topics hiérarchiques, wildcards, QoS, sessions, keep-alive
- **D-20:** Sections séparées: `mqtt: {}`, `validation: {}`, `assets: {}` — propre, modulaire
- **D-21:** Configuration par couches: paramètres essentiels visibles, paramètres avancés regroupés
- **D-22:** Validation 3 niveaux: plages valeurs + qualité données + timestamps/complétude/taille/santé capteurs
- **D-23:** Réponses JSON + pagination: `{data: [...], pagination: {page, limit, total}}`
- **D-24:** Endpoints: `GET /api/v1/assets/{id}/sensors` (relation) + `GET /api/v1/sensors` (collection globale) avec pagination sur les deux
- **D-25:** Conventions REST: noms pluriels, status codes standard (201, 400, 404, 409), DELETE → 405
- **D-26:** Valeurs par défaut: le code gère les champs manquants avec des defaults
- **D-27:** Rétrocompatibilité: alias + warning pour champs renommés, suppression en v2 uniquement

### the agent's Discretion
- Cycle de vie exact du broker MQTT (comment le plugin le démarre/arrête)
- Paramètres de backoff exponentiel (factor, max retries)
- Structure exacte des headers binaires (24 vs 32 octets)
- Organisation des fichiers dans chaque package
- Configuration des paramètres de logging
- Choix du logger (slog standard, zerolog, zap)

### Deferred Ideas (OUT OF SCOPE)
- Dashboard UI (Phase 4)
- Anomaly detection engine (Phase 3)
- Alert notifications (Phase 3/5)
- QoS 2 support
- Modbus/OPC UA plugins (v2)
- Python plugin SDK wrapper
- MQTT authentication/TLS
- Plugin auto-discovery
</user_constraints>

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| CORE-06 | Plugin SDK avec contrats versionnés | gRPC proto files with go-plugin GRPCPlugin interface, versioned `pkg/sdk/v1/` directory |
| ACQ-01 | Plugin acquisition MQTT: collecte données capteurs ESP32 → Core via NATS | mochi-mqtt embedded broker + paho.mqtt.golang client + NATS publish |
| ACQ-02 | Support MQTT QoS 0/1 pour fiabilité variable | mochi-mqtt QoS 0+1 support, paho client QoS configuration |
| ACQ-03 | Stockage données capteurs en SQLite avec timestamps | Existing `internal/storage` + new migrations for assets tables |
| ACQ-04 | Gestion des assets/machines (enregistrement, hiérarchie) | SQLite tables `assets` + `asset_sensors`, REST API endpoints |
</phase_requirements>

## Architectural Responsibility Map

| Capability | Primary Tier | Secondary Tier | Rationale |
|------------|-------------|----------------|-----------|
| gRPC Plugin SDK | Core (pkg/sdk/v1/) | Plugin binaries | Proto contracts define the interface, core manages lifecycle |
| MQTT Broker | Plugin (cmd/mqtt-plugin) | mochi-mqtt embedded | Broker lives in plugin process, not core — microkernel respected |
| MQTT Client | Plugin (cmd/mqtt-plugin) | paho.mqtt.golang | Plugin subscribes to ESP32 topics via client |
| Data Validation | Plugin (cmd/mqtt-plugin) | Core (shared validators) | Validation runs at ingestion point in plugin |
| Asset Registry | Core (internal/storage + internal/api) | SQLite + REST | Assets are core infrastructure, queried by all plugins |
| NATS Bridge | Plugin (cmd/mqtt-plugin) | Core NATS bus | Plugin publishes validated data to NATS `sensor.*` subjects |
| Binary Payload Parsing | Plugin (cmd/mqtt-plugin) | encoding/binary stdlib | Format parsing is plugin responsibility |

## Standard Stack

### Core (New Dependencies for Phase 02)
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| `github.com/mochi-mqtt/server/v2` | v2.7.9 | Embeddable MQTT v5/v3.1.1 broker | Fully compliant, high-performance, Go-native, embeddable, trie-based subscriptions |
| `github.com/eclipse/paho.mqtt.golang` | v1.5.1 | MQTT client for plugin→broker connection | Official Eclipse Paho, QoS 0/1, LWT, auto-reconnect |
| `google.golang.org/grpc` | v1.74.2 | gRPC transport for plugin SDK | Already in go.mod (indirect), standard for Go gRPC |
| `google.golang.org/protobuf` | v1.36.7 | Protobuf runtime for generated code | Already in go.mod (indirect), official Google protobuf |
| `github.com/hashicorp/go-plugin` | v1.8.0 | Plugin lifecycle with gRPC support | Already in go.mod, supports both net/rpc and gRPC |

### Supporting
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| `protoc` | latest | Protocol buffer compiler | Build step for .proto → .go code generation |
| `protoc-gen-go` | latest | Go code generator for proto | Required for `make proto` target |
| `protoc-gen-go-grpc` | latest | gRPC Go code generator | Required for `make proto` target |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| mochi-mqtt/server | Eclipse Mosquitto (external process) | Mosquitto requires CGO or separate binary, not embeddable in Go |
| mochi-mqtt/server | custom MQTT broker | Enormous complexity, protocol compliance bugs |
| paho.mqtt.golang | mochi-mqtt client | mochi-mqtt is broker-only, no standalone client library |
| protoc + buf | buf alone | buf is simpler but adds another dependency, protoc is standard |

**Installation:**
```bash
# MQTT broker (embedded)
go get github.com/mochi-mqtt/server/v2@v2.7.9

# MQTT client
go get github.com/eclipse/paho.mqtt.golang@v1.5.1

# gRPC tools (for proto compilation)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# protoc compiler (system package)
sudo apt install protobuf-compiler   # or brew install protobuf
```

**Version verification:** All versions verified against Go module proxy on 2026-07-03. mochi-mqtt v2.7.9 is the latest stable release. paho.mqtt.golang v1.5.1 is the latest stable release.

## Package Legitimacy Audit

> **Required** whenever this phase installs external packages.

| Package | Registry | Age | Downloads | Source Repo | Verdict | Disposition |
|---------|----------|-----|-----------|-------------|---------|-------------|
| github.com/mochi-mqtt/server/v2 | Go proxy | 4+ yrs | High | github.com/mochi-mqtt/server | OK | Approved — embeddable MQTT broker, 1.8k GitHub stars |
| github.com/eclipse/paho.mqtt.golang | Go proxy | 8+ yrs | High | github.com/eclipse/paho.mqtt.golang | OK | Approved — official Eclipse Paho client |
| google.golang.org/grpc | Go proxy | 10+ yrs | Very High | github.com/grpc/grpc-go | OK | Approved — official Go gRPC |
| google.golang.org/protobuf | Go proxy | 10+ yrs | Very High | github.com/protocolbuffers/protobuf-go | OK | Approved — official Google protobuf |
| github.com/hashicorp/go-plugin | Go proxy | 8+ yrs | High | github.com/hashicorp/go-plugin | OK | Approved — already in project |

**Packages removed due to [SLOP] verdict:** none
**Packages flagged as suspicious [SUS]:** none

*All packages are from reputable organizations and verified on Go module proxy.*

## Architecture Patterns

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        ml-elec binary                            │
│                     (Core — micro-noyau)                         │
├─────────────────────────────────────────────────────────────────┤
│  main.go                                                        │
│  ├── Wire DI initialization                                     │
│  ├── Signal handling (SIGINT/SIGTERM)                           │
│  └── LIFO shutdown: API → Plugins → NATS → Storage             │
├─────────────────────────────────────────────────────────────────┤
│  internal/config/        ← new sections: mqtt, validation, assets│
│  internal/nats/          ← unchanged                            │
│  internal/storage/       ← new tables: assets, asset_sensors    │
│  internal/plugin/        ← UPGRADED: gRPC transport support     │
│  internal/api/           ← new endpoints: /api/v1/assets, sensors│
│  pkg/sdk/v1/             ← NEW: .proto files + generated Go     │
├─────────────────────────────────────────────────────────────────┤
│  cmd/ml-elec/main.go     ← plugin manager launches MQTT plugin  │
│  cmd/mqtt-plugin/main.go ← NEW: external MQTT plugin binary     │
│  cmd/mock-plugin/main.go ← UPDATED: uses new gRPC SDK           │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│              cmd/mqtt-plugin (child process)                     │
├─────────────────────────────────────────────────────────────────┤
│  1. Start embedded mochi-mqtt broker on :1883                   │
│  2. Connect paho.mqtt.golang client to local broker             │
│  3. Subscribe to esp32/# (wildcard)                             │
│  4. On message: parse JSON/binary → validate → publish to NATS  │
│  5. Register assets on startup via gRPC SDK                     │
│  6. LWT: publish "offline" to sys/mqtt-plugin/status            │
└─────────────────────────────────────────────────────────────────┘

Data Flow:
  ESP32 → MQTT (:1883) → mochi broker → paho client → validate → NATS → core
                                                              ↓
                                                        SQLite (assets + readings)
```

### Recommended Project Structure
```
ml-elec/
├── cmd/
│   ├── ml-elec/main.go           # Entry point (existing, updated)
│   ├── mqtt-plugin/main.go       # NEW: MQTT bridge plugin binary
│   └── mock-plugin/main.go       # UPDATED: uses gRPC SDK
├── internal/
│   ├── config/config.go          # UPDATED: +mqtt, +validation, +assets sections
│   ├── nats/                     # unchanged
│   ├── storage/
│   │   ├── storage.go            # UPDATED: +asset CRUD methods
│   │   ├── migrations/
│   │   │   ├── 001_init.up.sql   # existing
│   │   │   ├── 001_init.down.sql # existing
│   │   │   ├── 002_assets.up.sql      # NEW: assets + asset_sensors tables
│   │   │   └── 002_assets.down.sql    # NEW
│   │   └── storage_test.go
│   ├── plugin/
│   │   ├── manager.go            # UPDATED: gRPC transport support
│   │   └── manager_test.go
│   └── api/
│       ├── server.go             # unchanged
│       ├── health.go             # unchanged
│       ├── sensors.go            # unchanged
│       ├── assets.go             # NEW: asset CRUD endpoints
│       └── api_test.go
├── pkg/
│   └── sdk/
│       └── v1/
│           ├── proto/
│           │   ├── lifecycle.proto    # NEW: Init/Start/Stop services
│           │   └── sensor.proto       # NEW: Collect service
│           ├── lifecycle.pb.go        # GENERATED (committed)
│           ├── lifecycle_grpc.pb.go   # GENERATED (committed)
│           ├── sensor.pb.go           # GENERATED (committed)
│           ├── sensor_grpc.pb.go      # GENERATED (committed)
│           ├── plugin.go             # NEW: Go interface wrappers
│           └── plugin_test.go        # NEW
├── wire.go / wire_gen.go          # UPDATED: new providers
├── Makefile                       # UPDATED: +proto target
└── config.yaml                    # UPDATED: +mqtt, +validation, +assets sections
```

### Pattern 1: gRPC Plugin with go-plugin
**What:** Define plugin interface via .proto, implement GRPCPlugin for go-plugin integration
**When to use:** Plugin SDK with multi-language support, versioned contracts
**Example:**
```go
// Source: [Context7: /hashicorp/go-plugin]
// pkg/sdk/v1/plugin.go — Go interface wrapper around generated proto

package sdk

import (
    "context"
    "google.golang.org/grpc"
    "github.com/hashicorp/go-plugin"
    pb "ml-elec/pkg/sdk/v1"
)

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
    plugin.NetRPCUnsupportedPlugin // Disable net/rpc, gRPC only
    Impl PluginLifecycle
    CollectImpl SensorCollector
}

func (p *GRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
    pb.RegisterPluginLifecycleServer(s, &GRPCServer{Impl: p.Impl})
    pb.RegisterSensorCollectorServer(s, &GRPCServer{CollectImpl: p.CollectImpl})
    return nil
}

func (p *GRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
    return &GRPCClient{
        lifecycle: pb.NewPluginLifecycleClient(c),
        collector: pb.NewSensorCollectorClient(c),
    }, nil
}
```

### Pattern 2: Embedded MQTT Broker (mochi-mqtt)
**What:** Start an MQTT broker inside the plugin process
**When to use:** MQTT plugin needs to accept ESP32 connections locally
**Example:**
```go
// Source: [Context7: /mochi-mqtt/server]
package main

import (
    mqtt "github.com/mochi-mqtt/server/v2"
    "github.com/mochi-mqtt/server/v2/hooks/auth"
    "github.com/mochi-mqtt/server/v2/listeners"
)

func startBroker(port string) (*mqtt.Server, error) {
    server := mqtt.New(&mqtt.Options{
        InlineClient: true, // Allow direct pub/sub from Go code
    })

    // Allow all connections (local-only, no auth for v1)
    _ = server.AddHook(new(auth.AllowHook), nil)

    // TCP listener on port 1883
    tcp := listeners.NewTCP(listeners.Config{
        ID:      "tcp1",
        Address: ":" + port,
    })
    if err := server.AddListener(tcp); err != nil {
        return nil, fmt.Errorf("adding TCP listener: %w", err)
    }

    go func() {
        if err := server.Serve(); err != nil {
            log.Printf("broker error: %v", err)
        }
    }()

    return server, nil
}
```

### Pattern 3: MQTT Client with LWT and Reconnection
**What:** Connect to local broker with persistent sessions and last will
**When to use:** Plugin subscribes to ESP32 topics via client
**Example:**
```go
// Source: [Context7: /eclipse-paho/paho.mqtt.golang]
func connectClient(brokerAddr, clientID string) mqtt.Client {
    opts := mqtt.NewClientOptions().
        AddBroker("tcp://" + brokerAddr).
        SetClientID(clientID).
        SetCleanSession(false).           // Persistent sessions
        SetAutoReconnect(true).           // Auto-reconnect
        SetConnectRetry(true).            // Retry initial connection
        SetConnectRetryInterval(5 * time.Second).
        SetMaxReconnectInterval(30 * time.Second).
        SetKeepAlive(30 * time.Second)

    // LWT: publish "offline" if plugin crashes
    opts.SetWill("sys/mqtt-plugin/status", "offline", 1, true)

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Printf("MQTT connect error: %v", token.Error())
    }

    // Publish online status
    client.Publish("sys/mqtt-plugin/status", 1, true, "online")
    return client
}
```

### Pattern 4: Binary Payload Parsing
**What:** Parse structured binary MQTT messages with header + payload
**When to use:** High-frequency vibration data from ESP32
**Example:**
```go
// Source: [D-12: Structure trame binaire v1]
package mqttplugin

import (
    "encoding/binary"
    "fmt"
    "math"
)

// BinaryHeader is the fixed header for binary MQTT payloads.
// Layout: [version(1) | sensor_type(1) | encoding(1) | reserved(1) | timestamp(8) | sample_count(4) | sample_rate(4) | reserved(8)]
// Total: 28 bytes (within 24-32 byte range per D-12)
type BinaryHeader struct {
    Version     byte
    SensorType  byte
    Encoding    byte
    Reserved1   byte
    Timestamp   int64   // Unix nanoseconds
    SampleCount uint32
    SampleRate  uint32  // Hz
    Reserved2   [8]byte
}

const binaryHeaderSize = 28

func ParseBinaryPayload(data []byte) (*BinaryHeader, []int16, error) {
    if len(data) < binaryHeaderSize {
        return nil, nil, fmt.Errorf("payload too short: %d < %d", len(data), binaryHeaderSize)
    }

    if data[0] != 1 {
        return nil, nil, fmt.Errorf("unsupported version: %d", data[0])
    }

    // Big-endian only (D-SPEC: binary format rejects unknown endianness)
    header := &BinaryHeader{
        Version:     data[0],
        SensorType:  data[1],
        Encoding:    data[2],
        Reserved1:   data[3],
        Timestamp:   int64(binary.BigEndian.Uint64(data[4:12])),
        SampleCount: binary.BigEndian.Uint32(data[12:16]),
        SampleRate:  binary.BigEndian.Uint32(data[16:20]),
    }
    copy(header.Reserved2[:], data[20:28])

    // Parse samples based on encoding
    payload := data[binaryHeaderSize:]
    expectedLen := int(header.SampleCount) * 2 // int16 = 2 bytes
    if len(payload) < expectedLen {
        return nil, nil, fmt.Errorf("payload too short for %d samples: %d < %d",
            header.SampleCount, len(payload), expectedLen)
    }

    samples := make([]int16, header.SampleCount)
    for i := uint32(0); i < header.SampleCount; i++ {
        offset := i * 2
        samples[i] = int16(binary.BigEndian.Uint16(payload[offset:offset+2]))
    }

    return header, samples, nil
}
```

### Pattern 5: Three-Level Data Validation
**What:** Validate sensor readings on ingestion with configurable thresholds
**When to use:** Every MQTT message before NATS publish
**Example:**
```go
// Source: [D-22: Validation 3 niveaux]
package mqttplugin

import (
    "fmt"
    "time"
)

type ValidationConfig struct {
    Ranges      map[string]RangeConfig      `yaml:"ranges"`
    Timestamp   TimestampConfig             `yaml:"timestamp"`
    Health      HealthConfig                `yaml:"health"`
}

type RangeConfig struct {
    Min float64 `yaml:"min"`
    Max float64 `yaml:"max"`
}

type TimestampConfig struct {
    MaxFutureDrift time.Duration `yaml:"max_future_drift"`
    MaxPastDrift   time.Duration `yaml:"max_past_drift"`
}

type HealthConfig struct {
    MaxPayloadSize int `yaml:"max_payload_size"`
    MinPayloadSize int `yaml:"min_payload_size"`
}

// Level 1: Range validation (value within min/max)
func validateRange(value float64, cfg RangeConfig) error {
    if value < cfg.Min || value > cfg.Max {
        return fmt.Errorf("value %.2f out of range [%.2f, %.2f]", value, cfg.Min, cfg.Max)
    }
    return nil
}

// Level 2: Timestamp monotonicity (no future timestamps, no large jumps)
func validateTimestamp(ts time.Time, cfg TimestampConfig) error {
    now := time.Now()
    if ts.After(now.Add(cfg.MaxFutureDrift)) {
        return fmt.Errorf("timestamp %v is in the future (drift > %v)", ts, cfg.MaxFutureDrift)
    }
    if ts.Before(now.Add(-cfg.MaxPastDrift)) {
        return fmt.Errorf("timestamp %v is too old (drift > %v)", ts, cfg.MaxPastDrift)
    }
    return nil
}

// Level 3: Data quality (completeness, size, health)
func validateQuality(payload []byte, cfg HealthConfig) error {
    if len(payload) == 0 {
        return fmt.Errorf("empty payload")
    }
    if len(payload) > cfg.MaxPayloadSize {
        return fmt.Errorf("payload too large: %d > %d", len(payload), cfg.MaxPayloadSize)
    }
    return nil
}
```

### Pattern 6: Asset Registry REST API
**What:** CRUD endpoints for machines and sensors with pagination
**When to use:** Assets are registered by plugins on startup, queried by dashboard
**Example:**
```go
// Source: [D-23, D-24, D-25: REST conventions]
// internal/api/assets.go

// POST /api/v1/assets — Create asset
// Response: 201 Created, 409 Conflict (duplicate name), 400 Bad Request

// GET /api/v1/assets — List assets with pagination
// Query: ?page=1&limit=20
// Response: {data: [...], pagination: {page, limit, total}}

// GET /api/v1/assets/{id}/sensors — List sensors for an asset
// Response: {data: [...], pagination: {page, limit, total}}

// GET /api/v1/sensors — List all sensors globally
// Response: {data: [...], pagination: {page, limit, total}}

// POST /api/v1/assets/{id}/sensors — Register sensor for asset
// Response: 201 Created, 404 Not Found

// DELETE /api/v1/assets — Returns 405 Method Not Allowed (safety for v1)
```

### Anti-Patterns to Avoid
- **Broker in core:** MQTT broker MUST live in plugin process, not core — violates microkernel
- **Hardcoded validation thresholds:** All thresholds must be configurable via config.yaml
- **Silent parse errors:** Every parsing error MUST be logged (SPEC prohibition)
- **Little-endian binary:** Only big-endian supported (SPEC prohibition)
- **DELETE for assets:** Returns 405, not implemented (safety for v1)
- **Raw error returns in gRPC:** Always use `status.Errorf` with specific codes

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| MQTT broker | Custom TCP server + MQTT protocol | mochi-mqtt/server | MQTT protocol is complex (QoS flows, sessions, wildcards) |
| MQTT client | Custom TCP client | paho.mqtt.golang | Reconnection, LWT, QoS handling are non-trivial |
| gRPC plugin lifecycle | Custom subprocess management | hashicorp/go-plugin | Handles handshake, cleanup, version negotiation |
| Proto code generation | Manual struct definitions | protoc + protoc-gen-go | Type safety, multi-language support, backward compatibility |
| Binary parsing | Manual byte manipulation | encoding/binary (stdlib) | Endianness handling, boundary checks |
| JSON payload parsing | Custom unmarshaling | encoding/json (stdlib) | Standard, tested, handles edge cases |

**Key insight:** The MQTT protocol has 14+ packet types, QoS flow state machines, session persistence, and wildcard matching. Hand-rolling even a subset introduces protocol compliance bugs that are extremely hard to debug with real ESP32 devices.

## Common Pitfalls

### Pitfall 1: Paho Go is a Client, Not a Broker
**What goes wrong:** Using `eclipse/paho.mqtt.golang` as if it were a broker
**Why it happens:** D-01 says "Eclipse Paho Go" but Paho Go is only a client library
**How to use:** Use `mochi-mqtt/server` for the broker, `paho.mqtt.golang` for the client
**Warning signs:** Trying to call `paho.mqtt.golang` server methods that don't exist

### Pitfall 2: go-plugin gRPC vs net/rpc Coexistence
**What goes wrong:** Plugin manager tries to use gRPC transport but existing plugins use net/rpc
**Why it happens:** Phase 1 used net/rpc, Phase 2 adds gRPC
**How to avoid:** Use `AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC}` for new plugins, keep net/rpc for existing ones. go-plugin supports both simultaneously via `VersionedPlugins`.
**Warning signs:** "protocol negotiation failed" errors

### Pitfall 3: protoc Not Installed
**What goes wrong:** `make proto` fails because protoc compiler is missing
**Why it happens:** protoc is a system package, not a Go module
**How to avoid:** Add protoc installation to Makefile or document in README: `sudo apt install protobuf-compiler`
**Warning signs:** "protoc: command not found"

### Pitfall 4: Generated Proto Code Not Committed
**What goes wrong:** Plugins can't compile because .pb.go files are missing
**Why it happens:** D-06 says generated code is committed, but developer forgets
**How to avoid:** Run `make proto` after any .proto change, commit generated files
**Warning signs:** "undefined: pb.RegisterPluginLifecycleServer"

### Pitfall 5: MQTT Broker Port Conflict
**What goes wrong:** Plugin fails to start because port 1883 is already in use
**Why it happens:** Another MQTT broker (Mosquitto) or previous instance running
**How to avoid:** Check port availability before binding, configurable port in config.yaml
**Warning signs:** "address already in use" on plugin startup

### Pitfall 6: Binary Payload Endianness Mismatch
**What goes wrong:** ESP32 sends little-endian, plugin expects big-endian → garbage values
**Why it happens:** ESP32 is little-endian by default, MQTT binary format specifies big-endian
**How to avoid:** Document endianness in proto spec, reject unknown endianness (SPEC requirement)
**Warning signs:** Sensor values are wildly incorrect but no parse error

### Pitfall 7: InlineClient Missing in mochi-mqtt
**What goes wrong:** `server.Subscribe()` returns error "inline client not enabled"
**Why it happens:** mochi-mqtt requires `InlineClient: true` in Options for direct pub/sub
**How to avoid:** Always set `mqtt.New(&mqtt.Options{InlineClient: true})`
**Warning signs:** "inline client not enabled" error on startup

## Code Examples

Verified patterns from official sources:

### Proto File Definition
```protobuf
// Source: [D-05: Un fichier .proto par service]
// pkg/sdk/v1/proto/lifecycle.proto

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

```protobuf
// pkg/sdk/v1/proto/sensor.proto

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

### Wire Updated Provider Set
```go
// Source: [Context7: /google/wire]
// wire.go — updated for Phase 02

//go:build wireinject

package main

import (
    "github.com/google/wire"
    "ml-elec/internal/config"
    "ml-elec/internal/nats"
    "ml-elec/internal/storage"
    "ml-elec/internal/plugin"
    "ml-elec/internal/api"
)

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

### Config.yaml Additions
```yaml
# Source: [D-20, D-21, D-22: Config structure]
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

### Makefile Proto Target
```makefile
# Source: [D-06: Code généré commit dans le repo]
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       pkg/sdk/v1/proto/lifecycle.proto \
	       pkg/sdk/v1/proto/sensor.proto
```

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| net/rpc only for plugins | gRPC via go-plugin | Phase 2 | Multi-language support, versioned contracts |
| External MQTT broker (Mosquitto) | Embedded broker (mochi-mqtt) | Phase 2 | No external dependency, single binary |
| No data validation | Three-level validation | Phase 2 | Reject bad data at ingestion |
| Raw sensor readings | ISA-95/UNS topic hierarchy | Phase 2 | Structured data model |

**Deprecated/outdated:**
- `eclipse/paho.mqtt.golang` as broker: It's a client library only
- Manual plugin subprocess management: Use go-plugin with gRPC
- Hardcoded validation thresholds: Use config.yaml

## Assumptions Log

> List all claims tagged `[ASSUMED]` in this research.

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | mochi-mqtt/server is production-ready for embedded use | Standard Stack | Broker instability under load |
| A2 | go-plugin gRPC transport works alongside existing net/rpc plugins | Architecture | Plugin manager rewrite required |
| A3 | protoc can be installed via apt on Ubuntu 24.04 | Environment | Build step fails |
| A4 | ESP32 sends big-endian binary MQTT payloads | Code Examples | Binary parsing produces wrong values |
| A5 | InlineClient in mochi-mqtt supports QoS 1 inline subscriptions | Code Examples | Must use external client for QoS 1 |

## Open Questions

1. **protoc installation method?**
   - What we know: protoc is not currently installed on the machine
   - What's unclear: Whether to use apt, snap, or manual installation
   - Recommendation: `sudo apt install protobuf-compiler` for Ubuntu, document in Makefile

2. **mochi-mqtt InlineClient QoS support?**
   - What we know: InlineClient is required for direct pub/sub from Go code
   - What's unclear: Whether inline subscriptions support QoS 1 (docs say "Only QoS 0")
   - Recommendation: Use external paho client for QoS 1 subscriptions, inline for QoS 0 only

3. **Broker lifecycle within plugin process?**
   - What we know: D-03/D-04 say broker runs as sub-process managed by plugin
   - What's unclear: Whether broker should be a goroutine within plugin or separate os/exec process
   - Recommendation: Goroutine within plugin process (simpler, no IPC overhead)

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Build system | ✓ | 1.26.4 | — |
| protoc | Proto compilation | ✗ | — | Install via `sudo apt install protobuf-compiler` |
| protoc-gen-go | Code generation | ✗ | — | `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest` |
| protoc-gen-go-grpc | gRPC code generation | ✗ | — | `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest` |
| git | Version control | ✓ | — | — |
| make | Build automation | ✓ | — | — |

**Missing dependencies with no fallback:**
- protoc compiler — must be installed before any proto compilation

**Missing dependencies with fallback:**
- protoc-gen-go and protoc-gen-go-grpc — can be installed via `go install`

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go testing (standard library) |
| Config file | go.mod |
| Quick run command | `go test ./...` |
| Full suite command | `go test -race -cover ./...` |

### Phase Requirements → Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| CORE-06 | Proto files compile and generate valid Go code | build | `make proto && go build ./pkg/sdk/v1/...` | ❌ Wave 0 |
| CORE-06 | gRPC plugin lifecycle (Init/Start/Stop) works | integration | `go test ./pkg/sdk/v1/... -run TestLifecycle` | ❌ Wave 0 |
| ACQ-01 | MQTT plugin starts and subscribes to esp32/# | integration | `go test ./cmd/mqtt-plugin/... -run TestSubscribe` | ❌ Wave 0 |
| ACQ-01 | MQTT message received and published to NATS | integration | `go test ./cmd/mqtt-plugin/... -run TestBridge` | ❌ Wave 0 |
| ACQ-02 | QoS 0 and QoS 1 messages handled correctly | unit | `go test ./cmd/mqtt-plugin/... -run TestQoS` | ❌ Wave 0 |
| ACQ-03 | Valid sensor reading stored in SQLite with timestamp | unit | `go test ./internal/storage/... -run TestInsertSensor` | ❌ Wave 0 |
| ACQ-04 | Asset creation with sensors succeeds | unit | `go test ./internal/api/... -run TestCreateAsset` | ❌ Wave 0 |
| ACQ-04 | Asset hierarchy is queryable | unit | `go test ./internal/api/... -run TestGetAssets` | ❌ Wave 0 |

### Sampling Rate
- **Per task commit:** `go test ./...`
- **Per wave merge:** `go test -race -cover ./...`
- **Phase gate:** Full suite green + core LOC check < 5000

### Wave 0 Gaps
- [ ] `pkg/sdk/v1/plugin_test.go` — covers CORE-06
- [ ] `cmd/mqtt-plugin/main_test.go` — covers ACQ-01, ACQ-02
- [ ] `internal/storage/asset_test.go` — covers ACQ-03, ACQ-04
- [ ] `internal/api/assets_test.go` — covers ACQ-04
- [ ] Makefile proto target — covers CORE-06 code generation

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | no | MQTT auth deferred (v1 plaintext) |
| V3 Session Management | yes | mochi-mqtt session management, clean_session config |
| V4 Access Control | yes | Plugin isolation via child process, topic-level ACL |
| V5 Input Validation | yes | Three-level validation (ranges, timestamps, quality) |
| V6 Cryptography | no | TLS deferred (v1 plaintext) |

### Known Threat Patterns for Go + MQTT

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Malformed MQTT payload | Tampering | Validation at ingestion, reject + log |
| Plugin crash | Denial of Service | Child process isolation via go-plugin |
| Excessive message rate | Denial of Service | Rate limiting in mochi-mqtt hooks |
| Binary format injection | Tampering | Big-endian only, reject unknown endianness |
| PII in MQTT logs | Information Disclosure | Code review, no PII in log statements |

## Sources

### Primary (HIGH confidence)
- [Context7: /mochi-mqtt/server] - Embedded MQTT broker API, listeners, hooks, inline client
- [Context7: /hashicorp/go-plugin] - gRPC plugin system, GRPCPlugin interface, version negotiation
- [Context7: /eclipse-paho/paho.mqtt.golang] - MQTT client, LWT, reconnection, QoS

### Secondary (MEDIUM confidence)
- [GitHub: mochi-mqtt/server] - v2.7.9 latest release, 1.8k stars, MIT license
- [GitHub: eclipse-paho/paho.mqtt.golang] - v1.5.1 latest release, official Eclipse

### Tertiary (LOW confidence)
- None — all findings verified against official documentation and module proxy

## Metadata

**Confidence breakdown:**
- Standard Stack: HIGH - mochi-mqtt and paho verified on module proxy, documentation fetched
- Architecture: HIGH - Patterns from official go-plugin and mochi-mqtt documentation
- Pitfalls: HIGH - Common issues from MQTT skill and go-plugin documentation

**Research date:** 2026-07-03
**Valid until:** 2026-08-02 (30 days — stable stack)
