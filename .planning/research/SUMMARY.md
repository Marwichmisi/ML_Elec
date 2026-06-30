# Project Research Summary

**Project:** ML_Elec — Modular Predictive Maintenance Platform
**Domain:** Industrial IoT (IIoT) / Predictive Maintenance
**Researched:** 2026-06-30
**Confidence:** MEDIUM-HIGH

## Executive Summary

ML_Elec is an edge-first, modular predictive maintenance platform targeting Raspberry Pi deployment with ESP32 sensor nodes. The research consensus is clear: the microkernel + plugin architecture (Go core + Python plugins via HashiCorp go-plugin) is the right foundation. This mirrors proven production systems (Terraform, Vault, EdgeX) and aligns with IIoT industry patterns (IIRA 5C model). The stack—Go 1.24+, embedded NATS, SQLite, scikit-learn IsolationForest, React 19 + Vite + shadcn/ui—is battle-tested and appropriate for single-device edge deployment.

The research reveals a critical insight: **60-85% of predictive maintenance projects fail not from model accuracy, but from data quality, operator trust, and integration gaps** (PwC, McKinsey, Siemens). The most dangerous combination for ML_Elec is "pilot purgatory" (never scaling beyond one asset) and "overengineering" (building everything before the core pipeline works). The recommended approach is strict phase discipline: sensor → core → anomaly → alert, proven on ONE motor type first, then expanded. Every phase must ship a working pipeline segment, not just infrastructure.

Key risks are: dirty sensor data breaking ML reliability (Pitfall #2), alert fatigue destroying operator trust (#3), hardware failures on edge devices (SD card death #5, thermal throttling #6), and core bloat collapsing the microkernel architecture (#7). Each has concrete prevention strategies documented in the pitfalls research. The architecture is deliberately minimal—Core handles only lifecycle, NATS, SQLite, and API. All domain logic lives in plugins. This boundary must be enforced from Phase 1.

## Key Findings

### Recommended Stack

The stack is built around edge-first constraints: single Raspberry Pi, no cloud dependency, plugin isolation, and open-source Apache 2.0 licensing. NATS replaces MQTT as the internal bus because it's faster (sub-ms latency), lighter (single binary), and supports request/reply patterns essential for plugin orchestration. MQTT (Mosquitto) handles only the ESP32 → Core bridge where it excels—constrained device communication.

**Core technologies:**
- **Go 1.24+**: Core runtime — static compilation, single binary, native concurrency, industry standard for IoT gateways
- **HashiCorp go-plugin v1.8.0**: Plugin isolation — battle-tested (Terraform/Vault), process isolation, Python support via gRPC
- **NATS 2.10+ (embedded)**: Internal message bus — single binary, sub-ms latency, JetStream persistence, leaf nodes for future scaling
- **SQLite 3.45+ (WAL mode)**: Configuration + state — zero-config, sub-1ms queries, sufficient for v1 workload (<100 assets)
- **scikit-learn 1.5+ (IsolationForest)**: Anomaly detection — gold standard for unsupervised PdM, 84% precision / 81% recall on NASA data
- **React 19 + Vite 6 + shadcn/ui**: Dashboard — 2025/2026 default stack, TanStack Query for real-time updates
- **Mosquitto 2.0+**: ESP32 MQTT broker — 200KB footprint, QoS levels, runs on Raspberry Pi

### Expected Features

**Must have (table stakes):**
- Real-time sensor data acquisition via MQTT — without this, nothing works
- Basic anomaly detection (threshold + scikit-learn) — core value proposition
- Real-time dashboard with live metrics — visualization is table stakes
- Asset management/registry — users need to know what they're monitoring
- Alerting/notifications — users must know when something is wrong
- Historical data storage — trend analysis requires history
- Configuration management — users configure sensors/thresholds
- Basic authentication — security requirement
- Plugin architecture — core value prop, extensibility

**Should have (competitive):**
- Edge-first deployment on Raspberry Pi — unique vs cloud-only competitors
- Multi-sensor fusion — 30-50% higher accuracy than single-sensor
- Open source + plugin ecosystem — community contributions
- Offline operation — critical for industrial environments

**Defer (v2+):**
- RUL estimation, failure mode classification, prescriptive recommendations
- ERP/GMAO integration, cloud deployment, multi-tenancy
- Community plugin marketplace, mobile native app
- Advanced ML (LSTM, transformers), FMEA libraries

### Architecture Approach

ML_Elec follows a microkernel + isolated plugin architecture. The Go Core handles only lifecycle management, NATS bus coordination, SQLite operations, and REST API exposure—all domain logic lives in Python plugins communicating via JSON-RPC over stdin/stdout (HashiCorp go-plugin pattern). Data flows northbound: ESP32 → MQTT (Mosquitto) → MQTT Plugin → NATS Bus → Core (SQLite storage) + Anomaly Plugin (analysis) → Dashboard (React). Configuration flows southbound: Dashboard → REST API → NATS → Plugins. The 6-layer model (Connection → Conversion → Cyber → Cognition → Configuration) aligns with IIRA 5C architecture.

**Major components:**
1. **Core (Go)** — Plugin lifecycle, NATS bus, REST API, SQLite, configuration
2. **NATS Bus (embedded)** — Internal message routing, pub/sub coordination
3. **MQTT Plugin (Python)** — Sensor data acquisition, protocol bridging to NATS
4. **Anomaly Detection Plugin (Python)** — Threshold + ML anomaly detection
5. **Dashboard (React + TypeScript)** — Real-time visualization, configuration UI

### Critical Pitfalls

1. **Pilot Purgatory (#1)** — Never scaling beyond first asset. Prevent: design data pipeline for multi-asset from day one (asset ID in every message, configurable per-asset thresholds). Test with ≥2 motor types during development.
2. **Dirty Sensor Data (#2)** — Missing values, clock skew, calibration drift. Prevent: validate every data point on ingestion (range checks, timestamp monotonicity). Use hardware RTC (DS3231). Log sensor health metrics.
3. **Alert Fatigue (#3)** — False positives destroy operator trust. Prevent: implement feedback loop from day one (every alert confirmed/dismissed by operator). Use confidence scoring, dynamic thresholds, severity levels.
4. **SD Card Death (#5)** — Edge hardware reliability. Prevent: use Compute Module 4 eMMC or USB SSD. Read-only root filesystem. SQLite WAL mode with `PRAGMA synchronous=NORMAL`. Monitor disk usage at 70%.
5. **Core Bloat (#7)** — Microkernel becomes monolith. Prevent: strict core contract (lifecycle, NATS, config, API only). Core size budget <5000 LOC. Plugin-first rule for all new features.

## Implications for Roadmap

Based on combined research, suggested 5-phase structure:

### Phase 1: Core Foundation
**Rationale:** Everything depends on the core. Plugins need the bus. Dashboard needs the API. This is the foundation with zero external dependencies.
**Delivers:** Go binary with embedded NATS, plugin manager (go-plugin), SQLite storage, REST API, configuration system.
**Addresses:** Plugin architecture (table stakes), configuration management, historical data storage.
**Avoids:** Core bloat (#7) — strict core contract from day one, size budget <5000 LOC.
**Stack:** Go 1.24+, go-plugin v1.8.0, NATS 2.10+ embedded, SQLite 3.45+ WAL, cobra, viper, zerolog.
**Pattern:** Microkernel + isolated plugins, NATS subject hierarchy design.
**Research flag:** LOW — well-documented patterns (HashiCorp, NATS, SQLite).

### Phase 2: Data Acquisition (MQTT Plugin)
**Rationale:** Without data, nothing else works. This is the "Connection" layer. Depends on Core NATS bus being operational.
**Delivers:** MQTT Plugin (Python) bridging ESP32 sensors to NATS, sensor data validation, protocol handling.
**Addresses:** Real-time sensor data acquisition (table stakes), MQTT/IoT protocol support.
**Avoids:** Data tsunami (#10) — edge processing, binary payloads, QoS=0 for vibration. Dirty data (#2) — validate every data point on ingestion.
**Stack:** Python 3.11+, paho-mqtt 1.6+, nats-py 2.7+, scikit-learn (for future).
**Pattern:** MQTT → NATS bridge, structured topic hierarchy.
**Research flag:** MEDIUM — MQTT QoS tuning, ESP32 firmware constraints need validation.

### Phase 3: Anomaly Detection Plugin
**Rationale:** Core value proposition—detect issues before failure. Depends on Phase 1 (Core) + Phase 2 (MQTT data feeding).
**Delivers:** Python plugin with threshold-based + IsolationForest ML anomaly detection, alert generation, health scoring.
**Addresses:** Basic anomaly detection (table stakes), alerting/notifications.
**Avoids:** Alert fatigue (#3) — feedback loop, confidence scoring, dynamic thresholds. Pilot purgatory (#1) — design for multi-asset from day one.
**Stack:** scikit-learn 1.5+ IsolationForest, NumPy 1.26+, pandas 2.2+.
**Pattern:** NATS subscribe sensor.*, publish alert.* and metrics.*.
**Research flag:** MEDIUM — ML model validation, threshold tuning methodology.

### Phase 4: Dashboard & User Interface
**Rationale:** Users need to see data and configure the system. Depends on Phase 1 (REST API) + Phase 3 (alerts to display).
**Delivers:** React + TypeScript dashboard with real-time metrics, anomaly alerts, asset management, configuration UI.
**Addresses:** Real-time dashboard (table stakes), asset management, configuration management.
**Avoids:** No action path (#4) — actionable alerts with context, log every alert outcome. Overengineering (#12) — strict scope discipline.
**Stack:** React 19, Vite 6, TypeScript 5.5+, shadcn/ui, TanStack Query 5+, Zustand 5+, Recharts 2+, Tailwind CSS 4+.
**Pattern:** REST API polling (upgrade to WebSocket later), TanStack Query cache invalidation.
**Research flag:** LOW — well-documented dashboard patterns.

### Phase 5: Integration & Demo
**Rationale:** Validate everything works together end-to-end. Prepare for soutenance demo.
**Delivers:** End-to-end pipeline validation, demo scenario with real sensor data, documentation, deployment guide.
**Addresses:** All table stakes features integrated and working.
**Avoids:** Pilot purgatory (#1) — demo must show multi-asset capability. Organizational resistance (#15) — alerts with context, explainability.
**Pattern:** Full pipeline: sensor → core → plugin → anomaly → alert → dashboard → user decision.
**Research flag:** LOW — integration testing, demo preparation.

### Phase Ordering Rationale

- **Phase 1 before everything** — Core is the foundation; plugins and dashboard depend on it
- **Phase 2 before Phase 3** — Anomaly detection needs sensor data flowing
- **Phase 3 before Phase 4** — Dashboard needs alerts and health scores to display
- **Phase 5 last** — Integration validates all prior phases work together
- **Each phase delivers a working pipeline segment** — avoids "pilot purgatory" by proving incrementally
- **Core bloat prevention enforced at every phase** — strict contract from Phase 1

### Research Flags

Phases likely needing deeper research during planning:
- **Phase 2:** MQTT QoS tuning for vibration data, ESP32 firmware constraints, binary payload formats
- **Phase 3:** ML model validation methodology, threshold tuning, IsolationForest hyperparameter optimization
- **Phase 5:** Demo scenario design, multi-asset pipeline validation

Phases with standard patterns (skip research-phase):
- **Phase 1:** Well-documented patterns (HashiCorp go-plugin, NATS, SQLite, Go REST APIs)
- **Phase 4:** Established dashboard patterns (React + Vite + shadcn/ui + TanStack Query)

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | **HIGH** | All technologies are industry-standard, battle-tested, well-documented |
| Features | **HIGH** | Clear table stakes from multiple vendor comparisons and PdM literature |
| Architecture | **MEDIUM-HIGH** | Microkernel + plugin pattern proven (HashiCorp, VS Code), IIRA 5C alignment, but IIoT-specific details need validation |
| Pitfalls | **HIGH** | Multiple industry reports (PwC, McKinsey, Siemens) cross-referenced, quantitative data available |

**Overall confidence:** MEDIUM-HIGH

### Gaps to Address

- **ESP32 firmware specifics:** Research covered MQTT protocol but not ESP32 firmware implementation details (memory constraints, sensor libraries). Handle during Phase 2 planning.
- **TimescaleDB migration path:** SQLite is sufficient for v1 but research notes it's not ideal for high-frequency time-series. Document migration criteria for v2 during Phase 1.
- **Plugin contract versioning:** Pitfall #8 (version skew) identified the need for semantic versioning and deprecation cycles, but specific contract design needs Phase 2 planning.
- **NATS subject hierarchy optimization:** Research identified the pattern but specific topic naming conventions for multi-asset scenarios need validation.
- **ML model explainability:** Pitfall #3 and #15 emphasize operator trust and explainability (SHAP values), but integration with IsolationForest scoring needs technical design in Phase 3.

## Sources

### Primary (HIGH confidence)
- HashiCorp go-plugin documentation (v1.8.0) — plugin isolation patterns
- NATS.io official documentation + MachineMetrics case study — embedded messaging
- scikit-learn IsolationForest docs + production benchmarks — anomaly detection
- PwC, McKinsey, Siemens PdM failure reports — pitfalls validation
- IIRA v1.10 (Industrial Internet Reference Architecture) — architecture patterns

### Secondary (MEDIUM confidence)
- React dashboard guide 2026 (usedatabrain.com) — dashboard stack patterns
- SQLite edge computing patterns (sqliteforum.com) — WAL optimization
- MQTT broker benchmarks 2026 (arxiv.org) — protocol comparison
- arc42 plugin architecture quality model — plugin pitfalls
- Nature Scientific Reports adaptive ML for PdM — model drift patterns

### Tertiary (LOW confidence)
- "Scalable Micro-Kernel with Go 2025" (gitconnected.com) — architectural patterns
- ESP32 production deployment guides — hardware constraints
- ClarityPoint industrial AI alignment analysis — organizational pitfalls

---
*Research completed: 2026-06-30*
*Ready for roadmap: yes*
