# Feature Landscape

**Domain:** Predictive Maintenance / Industrial IoT (IIoT)
**Researched:** 2026-06-30
**Project:** ML_Elec — Modular open-source predictive maintenance platform

## Executive Summary

Predictive maintenance platforms have evolved from simple threshold-based monitoring to sophisticated AI-driven systems. The market分为三个层次: table stakes (must-have), differentiators (competitive advantage), and anti-features (things to deliberately NOT build). For ML_Elec's v1, focus on core sensor acquisition, basic anomaly detection, and a simple dashboard. Advanced AI, CMMS integration, and enterprise features belong in v2+.

## Table Stakes

Features users expect. Missing = product feels incomplete or unusable.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| **Real-time sensor data acquisition** | Core function — without data, nothing works | Medium | Vibration, temperature, current sensors via MQTT |
| **MQTT/IoT protocol support** | Industry standard for IoT communication | Low | ESP32 → Core data transport |
| **Basic anomaly detection** | Primary value proposition | Medium | Threshold-based + simple ML (scikit-learn) |
| **Asset management/registry** | Users need to know what they're monitoring | Low | Simple asset hierarchy (plant → line → machine) |
| **Real-time dashboard** | Visualization is table stakes for monitoring | Medium | React + TypeScript, live metrics |
| **Alerting/notifications** | Users must know when something is wrong | Low | Email, webhook, or in-app alerts |
| **Historical data storage** | Trend analysis requires history | Low | SQLite for v1, sufficient for <100 assets |
| **Configuration management** | Users need to configure sensors/thresholds | Medium | Centralized config for core + plugins |
| **Basic user authentication** | Security requirement | Low | Simple JWT or session-based auth |
| **Plugin architecture** | Core value prop — extensibility | High | Child process + JSON-RPC (HashiCorp model) |

## Differentiators

Features that set product apart. Not expected, but valued.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| **Edge-first deployment** | Runs on Raspberry Pi — no cloud required | Medium | Unique positioning vs cloud-only competitors |
| **Multi-sensor fusion** | 30-50% higher accuracy than single-sensor | High | Vibration + temperature + current correlation |
| **RUL estimation** | Predict remaining useful life, not just alerts | High | Requires ML models (LSTM, transformers) |
| **Failure mode classification** | Identify specific fault type (bearing, imbalance) | High | Trained on failure libraries |
| **Open source + plugin ecosystem** | Community contributions, extensibility | Medium | Apache 2.0, plugin marketplace |
| **No-code sensor setup** | Deploy without data scientists | Medium | Pre-configured templates for common assets |
| **Prescriptive recommendations** | "What to do" not just "what's wrong" | High | Requires FMEA knowledge base |
| **Multi-protocol support** | Modbus, OPC-UA, MQTT, REST | Medium | Connects to existing industrial equipment |
| **Offline operation** | Works without internet | Medium | Critical for industrial environments |
| **Modular architecture** | Pick only what you need | Medium | Micro-kernel + plugins |

## Anti-Features

Features to explicitly NOT build for v1.

| Anti-Feature | Why Avoid | What to Do Instead |
|--------------|-----------|-------------------|
| **Full 3D digital twin** | Too complex, no ROI for v1 | 2D dashboard with health scores |
| **ERP/GMAO integration** | Enterprise feature, premature | Plugins for v2+ (premium) |
| **Cloud/SaaS deployment** | v1 is on-premise only | Local SQLite, optional cloud sync later |
| **Mobile native app** | Web-first, responsive design | Mobile-friendly web dashboard |
| **Advanced NLP/AI assistants** | Gimmicky for industrial users | Simple rule-based recommendations |
| **Massive ML model training** | Requires data scientists | Pre-trained models + simple threshold |
| **Complex FMEA libraries** | Domain expertise needed | Basic failure mode templates |
| **Real-time streaming analytics** | Overkill for v1 use cases | Batch processing with near-real-time alerts |
| **Multi-tenancy** | Enterprise feature | Single-tenant on-premise |

## Feature Dependencies

```
Sensor Acquisition → Data Storage → Anomaly Detection → Alerting
                                ↓
                        Historical Analysis → Dashboard
                                ↓
                        ML Models → RUL Estimation → Prescriptive Recommendations
                                ↓
                        Plugin System → Community Plugins → Advanced Features
```

## MVP Recommendation (Phase 1)

Prioritize for minimum viable product:
1. **Real-time sensor acquisition via MQTT** — without this, nothing works
2. **Basic anomaly detection** — core value proposition
3. **Simple dashboard** — users need to see what's happening
4. **Asset registry** — organize what you're monitoring
5. **Alerting** — notify users of issues

**Defer to Phase 2:**
- Advanced ML (RUL, failure classification)
- Multi-sensor fusion
- Plugin architecture
- Mobile optimization

**Defer to v2+:**
- ERP/GMAO integration
- Cloud deployment
- Advanced prescriptive analytics
- Community plugin marketplace

## Sources

- iFactory AI: "Predictive Maintenance Software 2026: 13 Expert Reviews & Comparison Guide"
- Tractian: "2026 Complete Guide to Predictive Maintenance Analytics"
- Siemens: "Senseye Predictive Maintenance" — asset health, failure prediction
- AspenTech: "Aspen Mtell" — prescriptive maintenance, FMEA integration
- EMQ: "The Unified MQTT Platform for Edge Computing" — MQTT, OPC-UA protocols
- Ifm: "moneo IIoT Core" — essential IIoT functions
- MachineCDN: "Predictive Maintenance Software Comparison 2026"
- Multiple vendor comparisons and feature matrices

---

*Last updated: 2026-06-30 after research phase*
