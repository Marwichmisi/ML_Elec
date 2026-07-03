# Roadmap: ML_Elec

**Created:** 2026-06-30
**Granularity:** fine
**Total Phases:** 7
**Coverage:** 24/24 v1 requirements mapped ✓

## Phases

- [x] **Phase 1: Core Foundation** - Go microkernel with NATS bus, plugin lifecycle, REST API, SQLite storage, and configuration
- [ ] **Phase 2: Plugin SDK & MQTT Bridge** - Versioned plugin contracts and Python MQTT plugin bridging ESP32 sensors to NATS
- [ ] **Phase 3: Anomaly Detection Engine** - Threshold rules + IsolationForest ML detection with confidence scoring and operator feedback
- [ ] **Phase 4: Real-time & Historical Dashboard** - React dashboard with live metrics, historical charts, and email notifications
- [ ] **Phase 5: Dashboard Configuration** - User-configurable widgets and alert notification management
- [ ] **Phase 6: Documentation & Developer Experience** - Professional README, plugin SDK docs, and getting-started guide
- [ ] **Phase 7: Demo & Multi-Asset Validation** - Live motor demo with real sensors, end-to-end pipeline validation, and multi-asset proof

## Phase Details

### Phase 1: Core Foundation
**Goal**: Users can start a Go binary that boots an embedded NATS bus, manages plugins via child processes, stores data in SQLite, and exposes a REST API
**Depends on**: Nothing (first phase)
**Requirements**: CORE-01, CORE-02, CORE-03, CORE-04, CORE-05, CORE-07
**Success Criteria** (what must be TRUE):
  1. Go binary starts and NATS bus is operational (pub/sub between components works)
  2. Plugin manager launches a mock Python plugin as a child process and communicates via JSON-RPC over stdin/stdout
  3. REST API responds to HTTP requests (health check, sensor data query)
  4. SQLite stores and retrieves sensor readings in WAL mode
  5. Configuration file loads and controls which plugins are enabled/disabled
  6. Core stays under 5000 LOC (no domain logic in core)
**Plans:** 5 plans

Plans:
- [x] 01-01-PLAN.md — Go module + config package (YAML loading, defaults, validation)
- [x] 01-02-PLAN.md — SQLite WAL storage + embedded NATS server
- [x] 01-03-PLAN.md — REST API (health + sensor endpoints)
- [x] 01-04-PLAN.md — Plugin manager (go-plugin, crash isolation, mock plugin)
- [x] 01-05-PLAN.md — Wire DI + main.go + graceful shutdown + integration tests

### Phase 2: Plugin SDK & MQTT Bridge
**Goal**: Developers can write plugins using a versioned SDK, and ESP32 sensor data flows from MQTT through the plugin into the NATS bus with validated timestamps and ranges
**Depends on**: Phase 1
**Requirements**: CORE-06, ACQ-01, ACQ-02, ACQ-03, ACQ-04
**Success Criteria** (what must be TRUE):
  1. Python plugin SDK provides typed client that connects to core via JSON-RPC
  2. MQTT plugin subscribes to ESP32 topics and publishes to NATS `sensor.*` subjects
  3. Data validation rejects readings with out-of-range values, missing fields, or non-monotonic timestamps
  4. QoS levels (0/1/2) are configurable per topic and correctly handled
  5. Sensor readings are persisted to SQLite with correct timestamps
  6. Asset registry supports registering machines with sensor hierarchy
**Plans:** 5 plans

Plans:
- [ ] 02-01-PLAN.md — gRPC Plugin SDK (proto files, generated code, Go interfaces)
- [ ] 02-02-PLAN.md — Config extensions + asset registry migrations & CRUD
- [ ] 02-03-PLAN.md — MQTT plugin with embedded broker, validation, NATS bridge
- [ ] 02-04-PLAN.md — Asset REST API endpoints with pagination
- [ ] 02-05-PLAN.md — Integration tests, crash isolation, performance benchmarks

### Phase 3: Anomaly Detection Engine
**Goal**: System automatically detects anomalous sensor readings using configurable thresholds and ML models, scoring each alert with a confidence level
**Depends on**: Phase 2
**Requirements**: DET-01, DET-02, DET-03
**Success Criteria** (what must be TRUE):
  1. Threshold rules detect out-of-range values for temperature, vibration, and current
  2. IsolationForest ML model identifies multivariate anomalies (combinations of sensors)
  3. Each anomaly alert includes a confidence score (0-100%)
  4. Operators can confirm or dismiss alerts, and that feedback is stored for future model improvement
  5. False positive rate stays below 30% on demo motor data after feedback loop matures
**Plans**: TBD

### Phase 4: Real-time & Historical Dashboard
**Goal**: Users can monitor live sensor data, view historical trends, and receive email notifications when anomalies are detected
**Depends on**: Phase 1
**Requirements**: VIS-01, VIS-02, VIS-05
**Success Criteria** (what must be TRUE):
  1. Dashboard displays live sensor metrics with <2s refresh latency
  2. Historical time-series charts show data over configurable time ranges (1h, 24h, 7d)
  3. Email notifications are sent when alerts reach configurable severity threshold
  4. Dashboard loads in under 3 seconds on Raspberry Pi
  5. Sensor health metrics (completeness %, timestamp gaps) are visible
**Plans**: TBD
**UI hint**: yes

### Phase 5: Dashboard Configuration
**Goal**: Users can customize their dashboard layout and configure alert notification preferences through the UI
**Depends on**: Phase 4
**Requirements**: VIS-03, VIS-04
**Success Criteria** (what must be TRUE):
  1. Users can add, remove, and rearrange dashboard widgets
  2. Each alert displays context: probable cause, recommended action, and priority level
  3. Widget configuration persists across sessions
  4. Users can configure notification channels (email, in-app) per severity level
**Plans**: TBD
**UI hint**: yes

### Phase 6: Documentation & Developer Experience
**Goal**: New users can get the system running in 5 minutes, and plugin developers have clear contracts to build new plugins
**Depends on**: Phase 1
**Requirements**: DEMO-03, DEMO-04
**Success Criteria** (what must be TRUE):
  1. README contains a 5-minute getting started guide that works on a fresh Raspberry Pi
  2. Plugin SDK documentation includes API reference, versioning policy, and a working example plugin
  3. Documentation covers installation, configuration, and troubleshooting
  4. All code examples in docs are tested and functional
**Plans**: TBD

### Phase 7: Demo & Multi-Asset Validation
**Goal**: Live demo shows end-to-end pipeline (sensor → MQTT → Core → Detection → Dashboard) on a real electric motor with multiple sensor types
**Depends on**: Phase 2, Phase 3, Phase 4
**Requirements**: DEMO-01, DEMO-02
**Success Criteria** (what must be TRUE):
  1. Demo scenario runs with real vibration, temperature, and current sensors on an electric motor
  2. Complete pipeline works: ESP32 → MQTT → Core → Anomaly Detection → Dashboard alerts
  3. System works with at least 2 different motor types (validates multi-asset design)
  4. Demo runs stably for 24+ hours without crashes or data loss
  5. All alert outcomes (confirmed/dismissed) are logged for audit trail
**Plans**: TBD

## Dependency Graph

```
Phase 1: Core Foundation
    ├── Phase 2: Plugin SDK & MQTT Bridge
    │       ├── Phase 3: Anomaly Detection Engine
    │       │       └── Phase 7: Demo & Multi-Asset Validation
    │       └── Phase 7: Demo & Multi-Asset Validation
    ├── Phase 4: Real-time & Historical Dashboard
    │       └── Phase 5: Dashboard Configuration
    │               └── Phase 7: Demo & Multi-Asset Validation
    └── Phase 6: Documentation & Developer Experience
            └── Phase 7: Demo & Multi-Asset Validation
```

## Parallelization Opportunities

With `parallelization: true`, these phases can run simultaneously:

| Wave | Phases | Rationale |
|------|--------|-----------|
| 1 | Phase 1 | Foundation — no dependencies |
| 2 | Phase 2, Phase 4, Phase 6 | All depend only on Phase 1; independent of each other |
| 3 | Phase 3, Phase 5 | Phase 3 depends on Phase 2; Phase 5 depends on Phase 4 |
| 4 | Phase 7 | Depends on Phases 2, 3, 4, 5 |

## Pitfall Mitigations by Phase

| Pitfall | Mitigation Phase | Strategy |
|---------|-----------------|----------|
| Core Bloat (#7) | Phase 1 | Strict core contract, size budget <5000 LOC, plugin-first rule |
| Dirty Sensor Data (#2) | Phase 2 | Validate every data point on ingestion (range, timestamp monotonicity) |
| MQTT Data Tsunami (#10) | Phase 2 | Edge processing, binary payloads, QoS=0 for vibration |
| Alert Fatigue (#3) | Phase 3 | Feedback loop, confidence scoring, dynamic thresholds |
| Pilot Purgatory (#1) | Phase 7 | Multi-asset validation with 2+ motor types |
| SD Card Death (#5) | Phase 1, 4 | WAL mode, log rotation, disk monitoring |
| NTP Time Warp (#14) | Phase 1, 2 | Hardware RTC support, monotonic time for intervals |
| Organizational Resistance (#15) | Phase 4, 5 | Explainable alerts, actionable context, human-in-the-loop |

## Progress

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Core Foundation | 5/5 | ✓ Complete | 2026-07-01 |
| 2. Plugin SDK & MQTT Bridge | 0/5 | Planned | - |
| 3. Anomaly Detection Engine | 0/4 | Not started | - |
| 4. Real-time & Historical Dashboard | 0/3 | Not started | - |
| 5. Dashboard Configuration | 0/2 | Not started | - |
| 6. Documentation & Developer Experience | 0/2 | Not started | - |
| 7. Demo & Multi-Asset Validation | 0/2 | Not started | - |

---
*Roadmap created: 2026-06-30*
*Total v1 requirements: 24*
*Coverage: 24/24 ✓*
