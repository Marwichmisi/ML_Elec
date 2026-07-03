---
phase: 02
slug: plugin-sdk-mqtt-bridge
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-07-03
---

# Phase 02 — Validation Strategy

> Contrat de validation par phase pour l'échantillonnage de feedback pendant l'exécution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go testing (standard) |
| **Config file** | none — see Wave 0 |
| **Quick run command** | `go test ./...` |
| **Full suite command** | `go test -v -race ./...` |
| **Estimated runtime** | ~30 seconds |

---

## Sampling Rate

- **Après chaque commit de tâche:** Run `go test ./...`
- **Après chaque wave de plan:** Run `go test -v -race ./...`
- **Avant `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 02-01-01 | 01 | 1 | CORE-06 | — | Proto compilation safe | unit | `protoc --go_out=. --go-grpc_out=. pkg/sdk/v1/*.proto` | ❌ W0 | ⬜ pending |
| 02-01-02 | 01 | 1 | CORE-06 | — | gRPC plugin lifecycle | integration | `go test ./internal/plugin/ -run TestGRPCPlugin` | ❌ W0 | ⬜ pending |
| 02-02-01 | 02 | 2 | ACQ-01 | T-02-01 | MQTT message validation | unit | `go test ./cmd/mqtt-plugin/ -run TestMQTTBridge` | ❌ W0 | ⬜ pending |
| 02-02-02 | 02 | 2 | ACQ-02 | T-02-02 | QoS level handling | unit | `go test ./cmd/mqtt-plugin/ -run TestQoS` | ❌ W0 | ⬜ pending |
| 02-03-01 | 03 | 2 | ACQ-03 | — | Data validation rules | unit | `go test ./internal/validation/ -run TestValidate` | ❌ W0 | ⬜ pending |
| 02-03-02 | 03 | 2 | ACQ-03 | — | Timestamp monotonicity | unit | `go test ./internal/validation/ -run TestTimestamp` | ❌ W0 | ⬜ pending |
| 02-04-01 | 04 | 3 | ACQ-04 | T-02-03 | Asset CRUD operations | unit | `go test ./internal/storage/ -run TestAsset` | ❌ W0 | ⬜ pending |
| 02-04-02 | 04 | 3 | ACQ-04 | — | Asset REST endpoints | unit | `go test ./internal/api/ -run TestAsset` | ❌ W0 | ⬜ pending |
| 02-05-01 | 05 | 3 | CORE-06 | — | Binary payload parsing | unit | `go test ./cmd/mqtt-plugin/ -run TestBinary` | ❌ W0 | ⬜ pending |
| 02-05-02 | 05 | 3 | ACQ-01 | — | JSON payload parsing | unit | `go test ./cmd/mqtt-plugin/ -run TestJSON` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `pkg/sdk/v1/lifecycle.proto` — proto file definitions
- [ ] `pkg/sdk/v1/sensor.proto` — proto file definitions
- [ ] `pkg/sdk/v1/*.pb.go` — generated Go code
- [ ] `cmd/mqtt-plugin/main.go` — MQTT bridge plugin
- [ ] `internal/validation/validator.go` — data validation
- [ ] `internal/storage/migrations/002_assets.up.sql` — asset tables
- [ ] `internal/storage/migrations/002_assets.down.sql` — rollback
- [ ] `internal/api/assets.go` — asset REST endpoints
- [ ] `internal/plugin/grpc.go` — gRPC plugin support

*Existing infrastructure covers Phase 1 requirements. Wave 0 adds proto compilation tools.*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| MQTT broker lifecycle (start/stop) | ACQ-01 | Requires process management verification | Start plugin, verify broker listens on :1883, stop plugin, verify port released |
| Crash isolation (kill MQTT plugin) | ACQ-01 | Requires OS process kill signal | `kill -9 <mqtt-plugin-pid>`, verify core continues running |
| 100 msg/s throughput | ACQ-01 | Requires timing measurement | Send 100 messages, measure broker→NATS latency |
| Auto-reconnect on broker crash | ACQ-01 | Requires broker restart sequence | Crash broker, verify plugin reconnects within 30s |

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
