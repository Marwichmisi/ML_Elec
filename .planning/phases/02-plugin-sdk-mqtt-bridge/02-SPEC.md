# Phase 2: Plugin SDK & MQTT Bridge — Specification

**Created:** 2026-07-03
**Ambiguity score:** 0.14 (gate: ≤ 0.20)
**Requirements:** 5 locked

## Goal

Developers can write plugins using a versioned gRPC SDK (proto files), and ESP32 sensor data flows from an embedded MQTT broker through an external plugin into the NATS bus with validated timestamps, configurable range validation, and an asset registry linking machines to sensors.

## Background

Phase 1 delivered a Go microkernel with an embedded NATS bus, a plugin manager (HashiCorp go-plugin, `net/rpc` transport), REST API, SQLite storage (WAL), and YAML configuration. The current plugin contract is a single `Echo(msg) (string, error)` method — a placeholder. No MQTT support, no versioned SDK, no asset registry, and no data validation exist today. Phase 2 builds the real plugin ecosystem: a gRPC-based SDK with proto files for multi-language support, an MQTT plugin as an external child process with an embedded broker, configurable data validation, and an asset/machine registry.

**Current codebase state:**
- `internal/plugin/manager.go` — plugin manager with go-plugin, `SensorPlugin` interface (Echo only)
- `internal/storage/storage.go` — SQLite CRUD, `sensor_readings` table
- `internal/config/config.go` — YAML config with `plugins.enabled` list
- `cmd/mock-plugin/main.go` — mock Go plugin (Echo implementation)
- No `.proto` files, no `pkg/` directory, no MQTT code, no MQTT dependency in go.mod
- Core LOC: ~1400 (well under 5000 limit)

## Requirements

1. **Plugin SDK (gRPC proto)**: Define versioned plugin contracts via `.proto` files with gRPC transport.
   - Current: Single `Echo` method via `net/rpc` — no versioning, no typed contracts
   - Target: `pkg/sdk/` directory with `.proto` files defining `PluginLifecycle` service (Init, Start, Stop) and `SensorCollector` service (Collect). Go code generated from protos. Plugin manager upgraded to support gRPC transport.
   - Acceptance: A proto file compiles without error; generated Go code implements the interfaces; existing `cmd/mock-plugin` is updated to use the new gRPC contract and all existing tests still pass.

2. **MQTT Plugin (external process)**: Plugin MQTT bridge runs as a child process via go-plugin, with an embedded MQTT broker.
   - Current: No MQTT code, no MQTT dependency in go.mod
   - Target: `cmd/mqtt-plugin/main.go` — external plugin binary that runs an embedded MQTT broker (QoS 0+1), subscribes to ESP32 topics via dynamic discovery (wildcard `esp32/#`), validates incoming readings, and publishes to NATS `sensor.*` subjects. Launched by the plugin manager as a child process.
   - Acceptance: Plugin starts, connects to embedded broker, receives MQTT messages, and publishes to NATS; crash isolation test passes (killing MQTT plugin process does not crash core).

3. **Data validation**: Sensor readings are validated on ingestion with configurable thresholds.
   - Current: No validation — raw data is stored as-is
   - Target: Readings are rejected if: missing required fields, non-monotonic timestamps, or values outside configurable min/max ranges (configurable per sensor type in `config.yaml`). Rejected readings are logged with reason.
   - Acceptance: A reading with out-of-range temperature is rejected and logged; a reading with a future timestamp is rejected; a valid reading passes; all validation rules are configurable in config.yaml.

4. **Asset registry**: Machines and sensors are registered with a hierarchy.
   - Current: No asset concept — only raw `sensor_readings` table
   - Target: SQLite tables `assets` (id, name, type, parent_id) and `asset_sensors` (asset_id, sensor_id, sensor_type). REST API endpoints: `POST /api/v1/assets`, `GET /api/v1/assets`, `POST /api/v1/assets/{id}/sensors`. Plugins register assets on startup.
   - Acceptance: Creating a machine asset with 3 sensors (temperature, vibration, current) succeeds; hierarchy is queryable; API returns correct JSON structure.

5. **MQTT data format support**: Support both JSON and binary MQTT payloads.
   - Current: No MQTT support
   - Target: JSON payloads parsed from standard ESP32 format (`{"temperature": 25.3, "vibration": 0.2, "timestamp": 1719900000}`). Binary payloads supported for high-frequency vibration data (configurable sensor type mapping).
   - Acceptance: A JSON MQTT message is parsed and stored correctly; a binary MQTT message (known format) is parsed and stored correctly; unknown format is rejected with log message.

## Boundaries

**In scope:**
- `pkg/sdk/` — gRPC plugin SDK with `.proto` files and generated Go code
- `cmd/mqtt-plugin/` — MQTT bridge plugin (external child process)
- Embedded MQTT broker (QoS 0+1) within the MQTT plugin
- Data validation (configurable ranges, timestamp monotonicity)
- Asset registry (SQLite tables + REST API)
- JSON + binary MQTT payload parsing
- Dynamic MQTT topic discovery (wildcard subscriptions)
- Plugin manager upgrade to support gRPC transport
- Config.yaml additions for MQTT, validation rules, asset configuration

**Out of scope:**
- Dashboard UI (Phase 4) — no frontend work in this phase
- Anomaly detection engine (Phase 3) — no detection logic
- Alert notifications (Phase 3/5) — no email/in-app alerts
- QoS 2 support — QoS 0+1 sufficient for industrial sensors
- Modbus/OPC UA plugins (v2 requirements) — MQTT only for now
- Python plugin SDK — gRPC proto supports it but no Python SDK wrapper in this phase
- MQTT authentication/TLS — plaintext for v1, security deferred
- Plugin auto-discovery — plugins listed in config.yaml

## Constraints

- Core must stay under 5000 LOC (MQTT + validation logic lives in plugin, not core)
- Broker embedded in MQTT plugin (not in core) — respects microkernel architecture
- Plugin manager must remain backward-compatible with existing `net/rpc` plugins (go-plugin supports both)
- Validation thresholds configurable via `config.yaml` (not hardcoded)
- MQTT plugin must handle 100 msg/s with <100ms latency from broker to NATS
- Go skills (golang-*) are mandatory for the executor agent — all code must follow Go best practices (error handling, naming, testing, lint)

## Acceptance Criteria

- [ ] `.proto` files compile and generate valid Go code
- [ ] `pkg/sdk/` contains versioned plugin interfaces (Init, Start, Stop, Collect)
- [ ] `cmd/mqtt-plugin` starts as child process and passes plugin manager lifecycle tests
- [ ] MQTT plugin subscribes to `esp32/#` and receives messages from embedded broker
- [ ] JSON MQTT payload is parsed and stored in SQLite with correct timestamp
- [ ] Binary MQTT payload is parsed and stored in SQLite (known format)
- [ ] Readings with out-of-range values are rejected and logged
- [ ] Readings with non-monotonic timestamps are rejected and logged
- [ ] All validation thresholds are configurable in config.yaml
- [ ] `POST /api/v1/assets` creates a machine with sensors
- [ ] `GET /api/v1/assets` returns asset hierarchy as JSON
- [ ] Existing Phase 1 tests still pass (no regressions)
- [ ] MQTT plugin handles 100 msg/s with <100ms broker-to-NATS latency
- [ ] Crash isolation: killing MQTT plugin process does not crash core
- [ ] MQTT plugin auto-reconnects on embedded broker crash
- [ ] SDK exposes no filesystem access methods (plugin isolation)
- [ ] No PII stored in MQTT logs or messages
- [ ] Parsing errors are always logged (never silently ignored)
- [ ] Binary format rejects data with unknown endianness
- [ ] DELETE endpoint for assets returns 405 Method Not Allowed

## Edge Coverage

**Coverage:** 11/11 applicable edges resolved · 0 unresolved

| Category | Requirement | Status | Resolution / Reason |
|----------|-------------|--------|---------------------|
| concurrency | R1 (SDK gRPC) | ✅ covered | SDK supports concurrent calls; Init/Start/Stop are sequential (one goroutine per plugin) |
| unclassified | R2 (MQTT Plugin) | ✅ covered | Auto-reconnect on embedded broker crash; plugin manager handles restart |
| boundary | R3 (Validation) | ✅ covered | Inclusive bounds — value equal to min/max is ACCEPTED (e.g. 150°C accepted if max=150) |
| adjacency | R3 (Validation) | ✅ covered | Readings are independent — identical values at same timestamp are both stored |
| empty | R3 (Validation) | ✅ covered | Empty MQTT payload (0 bytes) → rejected + log; missing JSON field → rejected + log |
| ordering | R3 (Validation) | ✅ covered | FIFO stable — readings with same timestamp stored in arrival order |
| precision | R3 (Validation) | ✅ covered | No rounding — float64 values stored natively (SQLite handles float precision) |
| idempotency | R4 (Asset Registry) | ✅ covered | POST with duplicate asset name → 409 Conflict error |
| concurrency | R4 (Asset Registry) | ✅ covered | Applicative lock (sync.Mutex) + SQLite WAL serializes concurrent writes |
| empty | R5 (MQTT Format) | ✅ covered | Empty MQTT payload → rejected + log |
| encoding | R5 (MQTT Format) | ✅ covered | Binary format: big-endian float32 for vibration data |

## Prohibitions (must-NOT)

**Coverage:** 5/5 applicable prohibitions resolved · 0 unresolved

| Prohibition (must-NOT statement) | Requirement | Status | Verification / Reason |
|----------------------------------|-------------|--------|------------------------|
| MUST NOT expose filesystem access methods from SDK to core (plugin isolation) | R1 (SDK gRPC) | resolved | verification: judgment — architectural review confirms no fs APIs in proto |
| MUST NOT store PII data in MQTT logs or messages | R2 (MQTT Plugin) | resolved | verification: judgment — code review for PII in log statements |
| MUST NOT silently ignore parsing errors — every error must be logged | R3 (Validation) | resolved | verification: test — negative test with malformed payload verifies log output |
| MUST NOT accept binary data with unknown endianness — big-endian only | R5 (MQTT Format) | resolved | verification: test — negative test with little-endian payload verifies rejection |
| MUST NOT expose DELETE endpoint for assets (safety for v1) | R4 (Asset Registry) | resolved | verification: test — HTTP test confirms DELETE returns 405 Method Not Allowed |

## Ambiguity Report

| Dimension          | Score | Min  | Status | Notes                                        |
|--------------------|-------|------|--------|----------------------------------------------|
| Goal Clarity       | 0.90  | 0.75 | ✓      | SDK gRPC + MQTT external plugin + asset registry |
| Boundary Clarity   | 0.90  | 0.70 | ✓      | Explicit out-of-scope (dashboard, detection, alerts) |
| Constraint Clarity | 0.80  | 0.65 | ✓      | 100 msg/s, QoS 0+1, core < 5000 LOC         |
| Acceptance Criteria| 0.80  | 0.70 | ✓      | 14 pass/fail criteria                        |
| **Ambiguity**      | 0.14  | ≤0.20| ✓      |                                              |

## Interview Log

| Round | Perspective     | Question summary                              | Decision locked                                        |
|-------|-----------------|-----------------------------------------------|--------------------------------------------------------|
| 1     | Researcher      | SDK Go pur, gRPC, ou Go+Python ?              | gRPC + proto files (multi-langage)                     |
| 1     | Researcher      | MQTT plugin externe ou intégré au core ?       | Plugin externe via go-plugin                           |
| 2     | Simplifier      | Broker MQTT externe ou embarqué ?              | Broker embarqué dans le plugin MQTT                    |
| 2     | Simplifier      | Registry d'assets dans Phase 2 ?              | Oui, tables SQLite + REST API                          |
| 2     | Simplifier      | Découverte des topics MQTT ?                   | Dynamique via wildcard esp32/#                         |
| 3     | Boundary Keeper | Dashboard/hors scope ?                        | Dashboard + Détection hors scope                       |
| 3     | Boundary Keeper | Format données MQTT ?                         | JSON + binaire                                         |
| 4     | Failure Analyst | Débit messages MQTT ?                         | 100 msg/s, <100ms                                      |
| 4     | Failure Analyst | QoS supporté ?                                | QoS 0+1 uniquement                                     |
| 5     | Seed Closer     | Méthodes cycle de vie SDK ?                   | Init/Start/Stop/Collect (4 méthodes)                   |
| 5     | Seed Closer     | Seuils validation configurables ?             | Oui, via config.yaml                                   |

---

*Phase: 02-plugin-sdk-mqtt-bridge*
*Spec created: 2026-07-03*
*Next step: /gsd-discuss-phase 2 — implementation decisions (how to build what's specified above)*
