# Architecture Patterns — ML_Elec

**Domain:** Modular IIoT / Predictive Maintenance Platform
**Researched:** 2026-06-30
**Overall confidence:** MEDIUM (synthesized from industry patterns + ML_Elec constraints)

---

## Executive Summary

ML_Elec follows a **microkernel + plugin architecture** running on edge hardware (Raspberry Pi). The system is structured as 6 distinct layers with clear boundaries: Sensors → Edge Collection → Transport (MQTT) → Core Processing (Go + NATS) → Analytics (Python plugins) → Visualization (React dashboard). This aligns with industry-standard IIoT patterns (IIRA 5C architecture, AWS edge patterns) while adapting to the project's constraints: single-device deployment, plugin isolation via child processes, and open-source Apache 2.0 licensing.

The key architectural insight from research: **edge-first predictive maintenance systems must process 90%+ of data locally** and only send aggregated/alerted data northbound. ML_Elec's microkernel design is well-suited for this—plugins run in isolated processes, communicate via JSON-RPC over stdin/stdout, and coordinate through the NATS bus. This provides fault tolerance (plugin crash doesn't kill core) while maintaining low-latency local processing.

---

## Recommended Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        ML_Elec Architecture                      │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   Sensors    │───▶│   Edge       │───▶│   Core       │      │
│  │  (ESP32)     │    │  Collection  │    │  (Go + NATS) │      │
│  │  vibration   │    │  Plugin      │    │              │      │
│  │  temperature │    │  (MQTT)      │    │  ┌────────┐  │      │
│  │  current     │    └──────────────┘    │  │ NATS   │  │      │
│  └──────────────┘                        │  │ Bus    │  │      │
│                                          │  └────────┘  │      │
│  ┌──────────────┐    ┌──────────────┐    │  ┌────────┐  │      │
│  │  Dashboard   │◀──│   REST API   │◀──│  │SQLite  │  │      │
│  │  (React+TS)  │    │   (Core)     │    │  │Storage │  │      │
│  └──────────────┘    └──────────────┘    │  └────────┘  │      │
│                                          └──────────────┘      │
│  ┌──────────────┐    ┌──────────────┐                          │
│  │  Anomaly     │◀──▶│  NATS Bus    │◀──▶  Other Plugins      │
│  │  Detection   │    │  (Internal)  │    │  (Future: RUL,      │
│  │  Plugin      │    └──────────────┘    │   Classification)   │
│  │  (Python)    │                        └──────────────┘      │
│  └──────────────┘                                                │
└─────────────────────────────────────────────────────────────────┘
```

### Component Boundaries

| Component | Responsibility | Language | Communicates With | Isolation Model |
|-----------|---------------|----------|-------------------|-----------------|
| **Core (Go)** | Plugin lifecycle, NATS bus management, REST API, SQLite operations, configuration | Go | All components | Process host (manages plugins) |
| **NATS Bus** | Internal message routing, pub/sub coordination | — | Core, all plugins | Message broker (embedded) |
| **MQTT Plugin** | Sensor data acquisition from ESP32, protocol bridging to NATS | Python | ESP32 sensors, NATS Bus | Child process (stdin/stdout JSON-RPC) |
| **Anomaly Detection Plugin** | Threshold-based + basic ML anomaly detection | Python | NATS Bus | Child process (stdin/stdout JSON-RPC) |
| **Dashboard** | Real-time visualization, metrics display, user interface | React/TS | Core (REST API) | Separate process (HTTP) |
| **SQLite Storage** | Persistent data storage, historical data, configuration | — | Core | Embedded database |
| **ESP32 Sensors** | Physical signal acquisition (vibration, temp, current) | C/C++ | MQTT Plugin | Hardware device |

### Data Flow — Northbound (Sensors → Decisions)

```
ESP32 Sensors
    │
    │  MQTT (WiFi)
    ▼
MQTT Plugin (Python)
    │  • Receives raw sensor packets
    │  • Validates data format
    │  • Publishes to NATS: sensor.{type}.{id}
    ▼
NATS Bus
    │  • Routes messages to subscribers
    │  • Core stores to SQLite
    │  • Anomaly plugin receives for analysis
    ▼
┌─────────────────────────────────────┐
│ Anomaly Detection Plugin (Python)   │
│  • Subscribes to sensor.* topics    │
│  • Applies threshold rules          │
│  • Runs basic ML models             │
│  • Publishes alerts to alert.*      │
│  • Publishes metrics to metrics.*   │
└─────────────────────────────────────┘
    │
    │  NATS: alert.{severity}
    ▼
Core (Go)
    │  • Receives alerts
    │  • Stores to SQLite
    │  • Exposes via REST API
    ▼
Dashboard (React)
    │  • Polls REST API
    │  • Displays real-time metrics
    │  • Shows anomaly alerts
    ▼
User Decision
```

### Data Flow — Southbound (Configuration → Sensors)

```
Dashboard (React)
    │  • User configures thresholds
    │  • User enables/disables plugins
    ▼
Core REST API
    │  • Validates configuration
    │  • Stores to SQLite
    │  • Publishes config changes to NATS
    ▼
NATS Bus
    │  • Routes config to relevant plugins
    ▼
Plugins
    │  • MQTT Plugin: adjusts collection rates
    │  • Anomaly Plugin: updates thresholds
    ▼
ESP32 Sensors (via MQTT)
    │  • Receives config updates
    │  • Adjusts sampling parameters
```

---

## Patterns to Follow

### Pattern 1: Microkernel + Isolated Plugins

**What:** Core handles lifecycle, plugins run in separate processes communicating via JSON-RPC over stdin/stdout.

**When:** Always — this is ML_Elec's fundamental architecture.

**Why:** Plugin crash doesn't kill core. Multi-language support (Go core, Python plugins). Independent development and testing.

**Example (from HashiCorp go-plugin, which ML_Elec mirrors):**

```go
// Core: Launches plugin as child process
cmd := exec.Command("python3", "mqtt_plugin.py")
stdin, _ := cmd.StdinPipe()
stdout, _ := cmd.StdoutPipe()

// JSON-RPC communication
client := jsonrpc2.NewClient(stdout, stdin)
```

**Boundaries:**
- Core owns: lifecycle, NATS bus, SQLite, REST API
- Plugins own: domain logic, external protocol handling
- Communication: only via JSON-RPC (stdin/stdout) and NATS pub/sub

### Pattern 2: NATS as Internal Message Bus

**What:** All inter-component communication flows through NATS subjects.

**When:** Always — core ↔ plugin, plugin ↔ plugin coordination.

**Why:** Decouples components, supports multi-language (Go + Python clients), lightweight for edge.

**Subject hierarchy:**
```
sensor.vibration.motor1     # Raw sensor data
sensor.temperature.motor1   # Temperature readings
alert.critical.bearing      # Anomaly alerts
alert.warning.temperature   # Warning-level alerts
config.anomaly.threshold    # Configuration updates
metrics.anomaly.score       # Health scores
```

**Key pattern:** Core subscribes to `sensor.*` for storage. Anomaly plugin subscribes to `sensor.*` for analysis. Dashboard reads from REST API (not directly from NATS).

### Pattern 3: Edge-First Processing

**What:** 90%+ of data processing happens locally on Raspberry Pi.

**When:** Always — ML_Elec is edge-first by design.

**Why:** Latency, bandwidth, reliability. Edge handles real-time decisions; cloud (if added later) handles heavy analytics.

**Implementation:**
- MQTT Plugin: receives and buffers sensor data locally
- Anomaly Detection: runs entirely on device
- SQLite: stores all historical data locally
- Dashboard: serves from local REST API
- NATS: embedded mode (no external NATS server needed)

### Pattern 4: Layered Data Processing

**What:** Data flows through distinct layers, each adding value.

**When:** Always — follows standard IIoT architecture (IIRA 5C model).

**Layers:**
1. **Connection Layer** (ESP32 + MQTT): Physical signal acquisition
2. **Conversion Layer** (MQTT Plugin): Protocol bridging, data normalization
3. **Cyber Layer** (Core + NATS): Data storage, routing, configuration
4. **Cognition Layer** (Anomaly Plugin): Analysis, anomaly detection, scoring
5. **Configuration Layer** (Dashboard): User interaction, visualization

**Why:** Each layer has clear responsibility. Easy to replace/upgrade individual layers.

---

## Anti-Patterns to Avoid

### Anti-Pattern 1: Plugin Interdependencies

**What:** Plugin A depends on Plugin B running.

**Why bad:** If B crashes, A breaks. Defeats modularity.

**Instead:** Plugins communicate only via NATS pub/sub. If B is down, A continues with available data.

### Anti-Pattern 2: Core as "Dumping Ground"

**What:** Core accumulates domain logic that should be in plugins.

**Why bad:** Core becomes complex, hard to maintain, violates microkernel principle.

**Instead:** Core stays minimal: lifecycle management, bus coordination, API exposure. All domain logic lives in plugins.

### Anti-Pattern 3: Synchronous Plugin Communication

**What:** Plugin A calls Plugin B directly and waits for response.

**Why bad:** Tight coupling, single point of failure, latency.

**Instead:** Use NATS request-reply pattern or publish/subscribe. Plugins don't know about each other.

### Anti-Pattern 4: Cloud-First Design

**What:** Design assumes cloud connectivity for core functionality.

**Why bad:** ML_Elec is edge-first. Cloud may not be available.

**Instead:** All core functionality works offline. Cloud (if added) is optional enhancement for backup/sync/heavy analytics.

---

## Component Details

### Core (Go)

**Responsibilities:**
- Plugin lifecycle management (start, stop, restart, health check)
- NATS bus initialization and management
- REST API endpoints for dashboard
- SQLite operations (read/write/query)
- Configuration management (load, validate, apply)
- Alert routing and notification

**Key interfaces:**
```go
type PluginManager interface {
    StartPlugin(name string) error
    StopPlugin(name string) error
    RestartPlugin(name string) error
    HealthCheck(name string) (PluginStatus, error)
}

type NATSBroker interface {
    Publish(subject string, data []byte) error
    Subscribe(subject string, handler MessageHandler) error
}

type Storage interface {
    StoreSensorData(data SensorReading) error
    QuerySensorData(query DataQuery) ([]SensorReading, error)
    StoreAlert(alert Alert) error
}
```

### MQTT Plugin (Python)

**Responsibilities:**
- Connect to MQTT broker (Mosquitto on Raspberry Pi)
- Subscribe to ESP32 sensor topics
- Validate and normalize incoming data
- Bridge to NATS: publish sensor readings
- Handle connection errors and reconnection

**Key pattern:**
```python
# MQTT → NATS bridge
def on_mqtt_message(client, userdata, msg):
    sensor_data = parse_sensor_payload(msg.payload)
    nats_publish(f"sensor.{sensor_data.type}.{sensor_data.id}", sensor_data)
```

### Anomaly Detection Plugin (Python)

**Responsibilities:**
- Subscribe to sensor.* NATS topics
- Apply threshold rules (configurable)
- Run basic ML models (scikit-learn)
- Publish alerts to alert.* topics
- Publish health scores to metrics.* topics

**Key pattern:**
```python
# Anomaly detection loop
async def detect_anomalies():
    async for msg in nats.subscribe("sensor.*"):
        score = calculate_anomaly_score(msg.data)
        if score > threshold:
            await nats.publish(f"alert.warning.{msg.subject}", alert)
        await nats.publish("metrics.anomaly.score", score)
```

### Dashboard (React + TypeScript)

**Responsibilities:**
- Display real-time sensor data (polling REST API)
- Show anomaly alerts with severity levels
- Allow configuration of thresholds and plugins
- Provide historical data visualization

**Key pattern:**
- Polls REST API every N seconds (configurable)
- WebSocket for real-time updates (future enhancement)
- State management for sensor readings and alerts

---

## Build Order Implications

### Phase 1: Core Foundation (Critical Path)

**Build first:** Core Go binary with NATS bus, plugin manager, SQLite, REST API.

**Why:** Everything depends on the core. Plugins need the bus. Dashboard needs the API.

**Dependencies:** None — this is the foundation.

### Phase 2: Data Acquisition (Depends on Phase 1)

**Build second:** MQTT Plugin (Python) for sensor data collection.

**Why:** Without data, nothing else works. This is the "Connection" layer.

**Dependencies:** Core NATS bus must be operational.

### Phase 3: Basic Analytics (Depends on Phase 1 + 2)

**Build third:** Anomaly Detection Plugin (Python).

**Why:** Core value proposition — detect issues before failure.

**Dependencies:** Core must be running. MQTT plugin must be feeding data.

### Phase 4: User Interface (Depends on Phase 1)

**Build fourth:** Dashboard (React + TypeScript).

**Why:** Allows users to see data and configure system.

**Dependencies:** Core REST API must be operational.

### Phase 5: Integration + Demo (Depends on All)

**Build last:** End-to-end integration, demo preparation.

**Why:** Validate everything works together.

**Dependencies:** All previous phases complete.

---

## Scalability Considerations

| Concern | At 1 Device | At 10 Devices | At 100 Devices |
|---------|-------------|---------------|----------------|
| NATS Throughput | Embedded mode sufficient | Embedded mode sufficient | Consider NATS cluster |
| SQLite Performance | Fine for 100K readings/day | Fine for 1M readings/day | Migrate to TimescaleDB |
| Plugin Isolation | Single process each | Process per device | Container per device |
| Dashboard Polling | 1s interval OK | 5s interval for batch | WebSocket streaming |
| MQTT Broker | Mosquitto single instance | Mosquitto clustering | EMQX or HiveMQ |

**Note:** ML_Elec v1 targets single-device (1 Raspberry Pi + ESP32s). Scalability is future concern.

---

## Key Architectural Decisions

| Decision | Choice | Rationale | Trade-off |
|----------|--------|-----------|-----------|
| Plugin Isolation | Child process + JSON-RPC | HashiCorp pattern, proven, simple | Overhead of process creation |
| Internal Bus | NATS embedded | Lightweight, multi-language, no external dependency | Less features than Kafka |
| Storage | SQLite | Zero-config, sufficient for v1, edge-friendly | Limited concurrency, no time-series optimization |
| Protocol | MQTT | Standard IoT, ESP32 native support | Requires MQTT broker (Mosquitto) |
| Dashboard | React + TypeScript | Type safety, ecosystem, maintainability | Heavier than simple HTML |

---

## Sources

- IIRA v1.10 (Industrial Internet Reference Architecture) — layered patterns
- AWS Industrial IoT Architecture Patterns — edge data flows
- MARTIN: End-to-end Microservice Architecture for PdM — component boundaries
- TIP4.0: IIoT Platform for Predictive Maintenance — modular platform design
- Superkind: PdM Architecture 6-Layer Model — industry best practices
- HashiCorp go-plugin — plugin isolation pattern reference
- NATS documentation — messaging patterns
- Edge Computing for IIoT (Springer, 2025) — edge architecture patterns
