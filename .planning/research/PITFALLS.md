# Domain Pitfalls: Predictive Maintenance / IIoT

**Domain:** Maintenance prédictive & intelligence industrielle
**Researched:** 2026-06-30
**Overall confidence:** MEDIUM-HIGH (multiple industry sources cross-referenced)

---

## Executive Summary

60-85% of predictive maintenance projects fail to deliver ROI within 18 months (PwC, McKinsey, Siemens). The failures are almost never about model accuracy — they're about data quality, integration gaps, operator trust, and scaling from pilot to production. For ML_Elec specifically, the combination of edge-first constraints (Raspberry Pi/ESP32), microkernel plugin architecture, and open-source scope creates a unique intersection of pitfalls that must be addressed from Phase 1.

---

## Critical Pitfalls

Mistakes that cause project failure or major rewrites.

### Pitfall 1: "Pilot Purgatory" — Never Scaling Beyond First Asset

**What goes wrong:** A working demo on one sensor/motor gives false confidence. The team tries to scale to 10+ assets and discovers the data pipeline, model, and alerting were never generalized. The pilot dies.

**Why it happens:** Pilots are scoped to one machine with clean data. Real deployments have different machine types, different operating conditions, different sensor calibrations, and no unified failure taxonomy.

**Consequences:** The demo works perfectly at soutenance but is useless in production. The project becomes a proof-of-concept that never ships.

**Prevention:**
- Design the data pipeline for multi-asset from day one (asset ID in every message, configurable per-asset thresholds)
- Use a **failure-mode taxonomy** — map work order labels to sensor signatures before training
- Test with at least 2 different motor types during development
- Keep the model simple enough to explain to a non-technical maintenance team

**Detection:** If the model only works on data from the exact same motor used during training, you're in pilot purgatory.

**Phase:** Phase 2-3 (data pipeline design, first plugin architecture)

**Sources:** IIoT World AI Manufacturing Day 2026, Algoscale PdM pilot analysis, Siemens Senseye framework

---

### Pitfall 2: "Garbage In, Garbage Out" — Dirty Sensor Data

**What goes wrong:** Sensor data arrives with missing values, timestamps drift, calibration drifts silently, sensors produce zero-value readings during dropouts. The model trains on this garbage and produces unreliable predictions.

**Why it happens:** Industrial sensors are not lab instruments. They experience:
- **Intermittent dropout** (not complete failure — harder to detect)
- **Calibration drift** over months
- **Clock skew** on Raspberry Pi without RTC (timestamps jump to 2023 after power cycle)
- **Electrical noise** from VFDs and nearby equipment

**Consequences:** False positive rate spikes. Operators lose trust. "The system that cries wolf."

**Prevention:**
- **Validate every data point** on ingestion: range checks, timestamp monotonicity, staleness detection
- **Use hardware RTC** (DS3231, ~3€) on Raspberry Pi — default system clock is unreliable
- **Never mix monotonic and wall-clock time** for intervals
- Log sensor health metrics: uptime, last-seen timestamp, drift rate
- **Start with QoS=0 for vibration monitoring** — missed window > delayed window (30s-old vibration data is useless)

**Detection:** Set up a sensor health dashboard showing: data completeness %, timestamp gaps, zero-value spike counts.

**Phase:** Phase 1 (core data ingestion), Phase 2 (MQTT plugin)

**Sources:** Raspberry Pi production deployment guides, vibration monitoring architecture case studies, MQTT QoS analysis

---

### Pitfall 3: "Alert Fatigue" — False Positives Destroy Operator Trust

**What goes wrong:** The system generates too many false alerts. Operators start ignoring them. Within months, unread alerts pile up. The system earns a nickname: "the system that cries wolf."

**Why it happens:**
- Threshold set once during pilot, never revisited
- No operational context in the model (doesn't know what job is running, what material, what tooling)
- Model trained during unusual stability; production has variability
- No feedback loop from operators to improve the model

**Consequences:** 47 unread alert emails. System becomes expensive decoration. Phase 2 never happens.

**Prevention:**
- **Implement a feedback loop from day one**: every alert must be confirmed/dismissed by an operator, and that feedback feeds back into the model
- **Use confidence scoring** — only surface high-confidence alerts to operators, log low-confidence ones separately
- **Dynamic thresholds** — don't set a fixed threshold; use percentile-based or adaptive thresholding
- **Alert deduplication** — suppress repeated alerts for the same anomaly within a configurable window
- **Severity levels** — CRITICAL (immediate action), WARNING (investigate), INFO (log only)

**Detection:** Track alert-to-action ratio. If <60% of alerts trigger a work order within 24h, the loop is broken.

**Phase:** Phase 3 (anomaly detection plugin), Phase 5 (dashboard + alerts)

**Sources:** ClarityPoint industrial AI failure analysis, KGT Solutions PdM audit framework, Nebulaworks 18-month production lessons

---

### Pitfall 4: "The Integration Gap" — No Path from Alert to Action

**What goes wrong:** The system detects anomalies and displays them on a dashboard, but nothing happens. No work order is created. No technician is dispatched. The detection-to-action pipeline has a manual gap.

**Why it happens:**
- Dashboard is designed for engineers, not operators
- No API to write work orders into CMMS
- No defined workflow: who owns the alert? What's the SLA? What's the escalation path?
- "Predictive alerts that don't trigger action are just alarms"

**Consequences:** The system is technically correct but operationally useless. Maintenance teams continue firefighting.

**Prevention:**
- **For ML_Elec v1**: Since we're not integrating with GMAO/ERP (out of scope), the "action" path must be the dashboard itself — make alerts actionable with context (asset name, location, sensor readings, recommended action)
- **Design the alert payload** to contain enough information for a technician to act without further investigation
- **Log every alert outcome**: actioned / ignored / missed — this data is gold for model improvement
- **API-first design** so GMAO integration is trivial in v2

**Detection:** Count alerts that triggered a response within 24h. If <60%, the pipeline is broken.

**Phase:** Phase 5 (dashboard), Phase 7 (API design for future GMAO integration)

**Sources:** KGT Solutions PdM audit, McKinsey reference cohort (30% downtime reduction when loop is closed)

---

### Pitfall 5: "SD Card Death" — Hardware Reliability on Edge Devices

**What goes wrong:** Raspberry Pi SD card silently corrupts after 6-24 months of continuous logging. The system crashes. Nobody notices until the next breakdown happens.

**Why it happens:**
- Consumer SD cards: 100-3,000 write cycles (TLC/QLC NAND)
- Continuous logging + swap + SQLite writes accelerate wear
- Power loss during write = filesystem corruption
- **No warning** — Pi doesn't report SD card degradation until it fails completely

**Consequences:** Unexplained crashes. Data loss. System bricked until someone physically re-flashes.

**Prevention:**
- **Use Compute Module 4 eMMC or USB SSD** instead of SD card (if budget allows)
- **Read-only root filesystem** with tmpfs overlay for writable dirs
- **Log rotation** configured aggressively (the `/var/log` fills up and blocks SQLite WAL writes)
- **SQLite WAL mode** with `PRAGMA synchronous=NORMAL` to reduce write amplification
- **UPS hat** or at minimum, graceful shutdown on power loss detection
- **Monitor disk usage** — set alert at 70% capacity

**Detection:** Run `dmesg | grep -i "mmc\|sd\|error"` periodically. Monitor SMART data if using SSD.

**Phase:** Phase 1 (infrastructure hardening), Phase 4 (deployment/production)

**Sources:** Industrial Monitor Direct Pi deployment guides, SiliconWit edge gateway architecture, DEV.to production Pi lessons

---

### Pitfall 6: "Thermal Throttling Ghost" — Invisible Performance Degradation

**What goes wrong:** Raspberry Pi throttles CPU under sustained load without logging a warning. MQTT reconnects slow down. Dashboard latency increases. The system appears to have "software bugs" but the root cause is hardware.

**Why it happens:**
- BCM2711 generates 3-5W under IIoT workloads
- In sealed enclosure + industrial ambient (35-45°C), junction temp exceeds 80°C
- DVFS silently reduces clock from 1.5 GHz to 1.0 GHz or lower
- **No log entry** — you must actively poll `vcgencmd get_throttled`

**Consequences:** Intermittent performance issues. MQTT lag. Dashboard "real-time" graphs are 45 seconds old. False diagnoses waste development time.

**Prevention:**
- **Active cooling** (heatsink + fan) mandatory for continuous operation
- **Monitor throttle status** in health endpoint: `vcgencmd get_throttled`
- **Benchmark under sustained load** (not just burst) — run for 24h, not 24 minutes
- **Don't trust development machine benchmarks** — Pi has 10-50x less memory headroom than dev laptop

**Detection:** `vcgencmd get_throttled` returns `0x0` when healthy, `0x50005` when throttled. Monitor continuously.

**Phase:** Phase 1 (infrastructure), Phase 4 (deployment validation)

**Sources:** Industrial Monitor Direct thermal analysis, LinkedIn Raspberry Pi edge workload study

---

### Pitfall 7: "Core Bloat" — Microkernel Becomes a Monolith

**What goes wrong:** The Go core accumulates "just one more feature" that should have been a plugin. The core becomes hard to maintain, test, and extend. Plugins can't be developed independently.

**Why it happens:**
- Pressure to ship fast for soutenance
- "It's easier to add to core than build a plugin interface"
- No clear boundary between core responsibility and plugin responsibility
- Core starts handling business logic that belongs in plugins

**Consequences:** Every core change risks breaking all plugins. Testing becomes impossible. The microkernel architecture collapses into a monolith with extra steps.

**Prevention:**
- **Strict core contract**: Core provides ONLY lifecycle, messaging (NATS), configuration, and API. Nothing else.
- **Plugin-first rule**: Any new feature goes through the plugin interface first. If it can't be a plugin, revisit the core contract.
- **Core size budget**: Core should be <5000 lines of Go. If it's growing, something is wrong.
- **Test core in isolation**: Core must work with mocked plugins. If tests require real plugins, the boundary is wrong.

**Detection:** Count lines of code in core vs plugins. If core > plugins, you've bloated.

**Phase:** Phase 1 (core architecture), reinforced at every phase

**Sources:** ArchMan microkernel patterns, VS Code architecture analysis, Go plugin system design patterns

---

### Pitfall 8: "Plugin Version Skew" — Breaking Changes Kill Ecosystem

**What goes wrong:** Core updates its API. Existing plugins break. Plugin developers (or you, 6 months later) don't know why. No deprecation path.

**Why it happens:**
- No semantic versioning for plugin contracts
- No compatibility testing between core versions and plugins
- "It's just a small API change" — until 5 plugins break

**Consequences:** Plugin ecosystem collapses. Nobody writes plugins because they break on every core update.

**Prevention:**
- **Semantic versioning** for the plugin contract (not just the core)
- **Deprecation cycle**: old API stays for 2 minor versions before removal
- **Compatibility matrix**: which plugin versions work with which core versions
- **Plugin contract tests**: CI runs plugin compatibility tests on every core change
- **Version negotiation**: plugin declares expected core version on load; core rejects incompatible plugins with clear error

**Detection:** Any time a core update requires changing plugin code, you have a versioning problem.

**Phase:** Phase 2 (plugin interface design), ongoing

**Sources:** arc42 plugin architecture quality model, Extism plugin security analysis, Go plugin system patterns

---

### Pitfall 9: "Model Drift" — Silent Degradation After Deployment

**What goes wrong:** Model works great in month 1. By month 6, it's catching 40% of what it should. Nobody notices because there's no monitoring.

**Why it happens:**
- Equipment ages → vibration patterns change
- Overhauls change the baseline signature (new bearings ≠ old bearings)
- Production mix changes (different materials, tooling, feed rates)
- Sensor calibration drifts
- **Nobody owns the model after go-live** — data science demonstrated on Jupyter, then moved on

**Consequences:** False positive rate spikes. Real failures go undetected. System quietly degrades to run-to-failure with extra steps.

**Prevention:**
- **Monitor model performance metrics** in production: alert precision, false positive rate, detection latency
- **Retraining schedule**: plan for periodic retraining (monthly/quarterly) even if data hasn't changed
- **Baseline snapshots**: save model performance baselines; alert when metrics degrade >10%
- **Feedback loop**: operator confirmations become training data for next model version
- **Asset lifecycle awareness**: flag when an asset has been overhauled and needs baseline recalibration

**Detection:** Track alert-to-confirmed-anomaly ratio. If it drops below 70%, model is drifting.

**Phase:** Phase 3 (anomaly detection), Phase 6 (monitoring/observability)

**Sources:** Nebulaworks 18-month production lessons, Algoscale PdM data foundation analysis, Nature Scientific Reports adaptive ML models

---

### Pitfall 10: "MQTT Data Tsunami" — Flooding the Broker with Raw Samples

**What goes wrong:** ESP32 sends raw vibration samples (10kHz) as JSON over MQTT. Broker chokes. Memory spikes. Dashboard shows data from 45 seconds ago. "Real-time" is a lie.

**Why it happens:**
- `{"value": 0.123}` = 20 bytes of JSON to transmit 4 bytes of float
- Each MQTT message has 40-60 bytes of overhead
- At 10kHz, that's 600KB/s of overhead alone
- QoS=1 causes backlog floods after network hiccups

**Consequences:** Broker becomes the bottleneck. Data is delayed or lost. System doesn't scale beyond 1-2 sensors.

**Prevention:**
- **Edge processing**: FFT, feature extraction, and anomaly detection happen on ESP32/RPi, not in the cloud
- **Binary payloads**: send processed features as binary blobs, not JSON
- **Windowed transmission**: buffer N samples, compute features, send one message per window
- **QoS=0 for vibration monitoring**: missed window > 30s-old window
- **Structured topic hierarchy**: `company/site/zone/machine/sensor` — enables efficient subscription

**Detection:** Monitor broker message rate and latency. If latency >5s, you're flooding.

**Phase:** Phase 2 (MQTT plugin), Phase 4 (edge processing optimization)

**Sources:** Guatu Labs vibration monitoring architecture, MQTTfy industrial monitoring case study, MQTT performance evaluation papers

---

## Moderate Pitfalls

### Pitfall 11: "No Failure History" — Training Without Labels

**What goes wrong:** Supervised models need labeled failure data. Most facilities haven't documented historical failures. Models trained on "normal" data can't predict specific failure modes.

**Prevention:** Start with unsupervised anomaly detection (no labels needed). Collect failure events as they happen. Build labeled dataset over time. Don't attempt supervised learning until you have 50+ confirmed failure events.

**Phase:** Phase 3 (anomaly detection plugin design)

---

### Pitfall 12: "Overengineering the First Version"

**What goes wrong:** Building digital twins, cloud integration, GMAO sync, mobile app, and 3D visualization before the basic sensor→alert pipeline works.

**Prevention:** ML_Elec's out-of-scope list is correct. Resist the urge. The demo must show: sensor data → anomaly detected → operator notified. Everything else is v2+.

**Phase:** All phases (scope discipline)

---

### Pitfall 13: "Plugin Resource Leaks" — Crashing the Host Over Time

**What goes wrong:** Plugins that don't release memory, file handles, or goroutines on unload. After days/weeks, the core process runs out of resources.

**Prevention:** Implement resource budgets per plugin (CPU, memory, runtime). Monitor per-plugin resource usage. Kill plugins that exceed limits. Test plugins under sustained load (24h+), not just burst.

**Phase:** Phase 2 (plugin isolation design)

**Sources:** arc42 plugin architecture, Go plugin system patterns

---

### Pitfall 14: "NTP Time Warp" — Clock Skew Breaking Everything

**What goes wrong:** Raspberry Pi boots without internet. System clock defaults to 2023. All timestamps are wrong. MQTT messages get timestamps in 2023. When NTP syncs, retention policy deletes them as "too old."

**Prevention:** Hardware RTC (DS3231). Use monotonic time for intervals. Configure `chrony` for aggressive sync on boot. Never use `datetime.now()` for critical timestamps without NTP verification.

**Phase:** Phase 1 (infrastructure)

**Sources:** DEV.to Raspberry Pi production deployment guide

---

### Pitfall 15: "Organizational Resistance" — Technicians Don't Trust the System

**What goes wrong:** Maintenance teams weren't consulted during design. Dashboards are built for engineers, not operators. Alerts lack context. "Will this replace me?" fear goes unaddressed.

**Prevention:** Involve operators from the design phase. Build dashboards around their workflow. Make alerts explain WHY (SHAP values, contributing features). Keep humans in the loop (operator approves every work order). Show wins early.

**Phase:** Phase 5 (dashboard design)

**Sources:** Springer PdM barriers study, ClarityPoint industrial AI alignment analysis

---

## Minor Pitfalls

### Pitfall 16: "JSON-RPC Overhead" — IPC Bottleneck for High-Frequency Data

**What goes wrong:** Child process communication via JSON-RPC adds serialization overhead. At high sensor frequencies, the IPC channel becomes a bottleneck.

**Prevention:** Batch messages. Use shared memory or mmap for high-frequency data. JSON-RPC is fine for control messages; data should use a faster channel.

**Phase:** Phase 2 (core-plugin communication)

---

### Pitfall 17: "NATS as Single Point of Failure" — No Broker Redundancy

**What goes wrong:** NATS embedded mode has no clustering. If the NATS process crashes, all inter-plugin communication stops.

**Prevention:** For v1 (single-node), this is acceptable. For v2+, consider NATS clustering or a fallback mechanism. Monitor NATS health.

**Phase:** Phase 2 (NATS setup)

---

### Pitfall 18: "SQLite Under Concurrent Write Load"

**What goes wrong:** Multiple plugins writing to SQLite simultaneously causes locking and slow queries. WAL mode helps but doesn't eliminate the issue.

**Prevention:** Use WAL mode. Serialize writes through the core. Consider a write-ahead queue. For v1, SQLite is fine if writes are moderate.

**Phase:** Phase 1 (storage design)

---

### Pitfall 19: "ESP32 Memory Constraints" — Running Out of RAM on Microcontrollers

**What goes wrong:** ESP32 has ~320KB RAM. Running MQTT client + sensor reading + basic processing can exhaust memory. Crashes are hard to debug.

**Prevention:** Profile ESP32 memory usage early. Use lightweight MQTT libraries. Offload heavy processing to RPi. Keep ESP32 firmware minimal (collect + transmit, nothing else).

**Phase:** Phase 2 (ESP32 firmware)

---

### Pitfall 20: "Forgetting the Human Factor" — No SOP for Alert Response

**What goes wrong:** Alert fires. Nobody knows what to do. Alert gets ignored. Next alert gets ignored. Pattern becomes habit.

**Prevention:** Define a Standard Operating Procedure for every alert severity level. Who gets notified? What's the response time? What's the escalation path? Even for v1, document this in the README.

**Phase:** Phase 5 (dashboard + alerting)

---

## Phase-Specific Warnings

| Phase Topic | Likely Pitfall | Mitigation |
|-------------|---------------|------------|
| Core architecture (Go) | Core bloat (#7) | Strict core contract, size budget <5000 LOC |
| MQTT plugin | Data tsunami (#10) | Edge processing, binary payloads, QoS=0 |
| Anomaly detection | False positive fatigue (#3) | Feedback loop, confidence scoring, dynamic thresholds |
| Dashboard | No action path (#4) | Actionable alerts with context, log outcomes |
| Edge deployment | SD card death (#5) | SSD/eMMC, read-only root, log rotation |
| Edge deployment | Thermal throttling (#6) | Active cooling, health endpoint monitoring |
| Plugin ecosystem | Version skew (#8) | SemVer, deprecation cycle, compatibility tests |
| Production monitoring | Model drift (#9) | Retraining schedule, performance baselines, feedback loop |
| Soutenance demo | Pilot purgatory (#1) | Multi-asset pipeline from day one, failure taxonomy |
| All phases | Overengineering (#12) | Strict scope discipline, out-of-scope list |

---

## Key Insight for ML_Elec

The most dangerous pitfall for this project is the **combination** of #1 (pilot purgatory) and #12 (overengineering). As a soutenance project, there's temptation to build everything at once. The research is clear: **scope in 4-6 week increments around specific assets**. Start with ONE motor, ONE sensor type, ONE anomaly detection method. Prove the full pipeline works (sensor → core → plugin → anomaly → alert). Then expand.

The second most dangerous is **operator trust** (#3, #4, #15). If the demo shows false positives or alerts without context, the jury won't trust the system. Build explainability into the anomaly detection from day one.

---

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| PdM project failures | HIGH | Multiple industry reports (PwC, McKinsey, Siemens, IIoT World) cross-referenced |
| Raspberry Pi edge pitfalls | HIGH | Multiple production deployment guides with specific failure modes |
| MQTT/IIoT data pipeline | HIGH | Case studies with specific metrics (latency, bandwidth, QoS analysis) |
| Plugin architecture | MEDIUM | Based on general microkernel patterns + VS Code analysis, not IIoT-specific |
| Anomaly detection pitfalls | MEDIUM-HIGH | Academic papers + production case studies with quantitative data |
| Organizational pitfalls | MEDIUM | Based on consulting reports, less specific to our context |

---

## Sources

- IIoT World AI Manufacturing Day 2026 panel (Viorel, 2026-06-24)
- Algoscale PdM pilot analysis (2026-04-20)
- Nebulaworks "Predictive Maintenance at Scale" (2025-11-08)
- Siemens Senseye "Avoiding Failure in Projects" framework
- KGT Solutions PdM audit framework (2026-04-30)
- ClarityPoint industrial AI alignment analysis (2026-03-25)
- Springer Schmalenbach Journal PdM barriers study (2025-01-21)
- Fierce Sensors AI-IoT failure analysis (2026-02-20)
- Industrial Monitor Direct Pi deployment guides (2026-04)
- Guatu Labs vibration monitoring architecture (2026-05-08)
- arc42 plugin architecture quality model (2026-06-09)
- ArchMan microkernel patterns (2026)
- Nature Scientific Reports adaptive ML for PdM (2026-03-07)
- ACM SIGKDD thresholding techniques for PdM (Giannoulidis et al., 2022)
- MDPI SVM-based false positive classification (2023)
