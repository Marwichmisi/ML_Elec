---
phase: 02
slug: plugin-sdk-mqtt-bridge
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-07-03
---

# Phase 02 — Validation Strategy

> Per-phase validation contract for feedback sampling during execution.

---

## Test Infrastructure

| Property | Value |
|----------|-------|
| **Framework** | Go testing (standard library) |
| **Config file** | go.mod |
| **Quick run command** | `go test ./...` |
| **Full suite command** | `go test -race -cover ./...` |
| **Estimated runtime** | ~30 seconds |

---

## Sampling Rate

- **After every task commit:** Run `go test ./...`
- **After every plan wave:** Run `go test -race -cover ./...`
- **Before `/gsd-verify-work`:** Full suite must be green
- **Max feedback latency:** 30 seconds

---

## Per-Task Verification Map

| Task ID | Plan | Wave | Requirement | Threat Ref | Secure Behavior | Test Type | Automated Command | File Exists | Status |
|---------|------|------|-------------|------------|-----------------|-----------|-------------------|-------------|--------|
| 02-01-01 | 01 | 1 | CORE-06 | — | Proto compilation, gRPC lifecycle | build | `make proto && go build ./pkg/sdk/v1/...` | ❌ W0 | ⬜ pending |
| 02-01-02 | 01 | 1 | CORE-06 | — | gRPC plugin Init/Start/Stop | integration | `go test ./pkg/sdk/v1/... -run TestLifecycle` | ❌ W0 | ⬜ pending |
| 02-02-01 | 02 | 1 | ACQ-01 | T-02-01 | MQTT plugin subscribe esp32/# | integration | `go test ./cmd/mqtt-plugin/... -run TestSubscribe` | ❌ W0 | ⬜ pending |
| 02-02-02 | 02 | 1 | ACQ-01 | T-02-02 | MQTT→NATS bridge | integration | `go test ./cmd/mqtt-plugin/... -run TestBridge` | ❌ W0 | ⬜ pending |
| 02-02-03 | 02 | 2 | ACQ-02 | — | QoS 0+1 handling | unit | `go test ./cmd/mqtt-plugin/... -run TestQoS` | ❌ W0 | ⬜ pending |
| 02-03-01 | 03 | 2 | ACQ-03 | T-02-03 | Three-level validation | unit | `go test ./internal/validation/... -run TestValidate` | ❌ W0 | ⬜ pending |
| 02-03-02 | 03 | 2 | ACQ-03 | — | Rejected readings logged | unit | `go test ./internal/validation/... -run TestRejection` | ❌ W0 | ⬜ pending |
| 02-04-01 | 04 | 3 | ACQ-04 | — | Asset creation with sensors | unit | `go test ./internal/api/... -run TestCreateAsset` | ❌ W0 | ⬜ pending |
| 02-04-02 | 04 | 3 | ACQ-04 | — | Asset hierarchy query | unit | `go test ./internal/api/... -run TestGetAssets` | ❌ W0 | ⬜ pending |
| 02-05-01 | 05 | 3 | ACQ-01 | — | JSON payload parsing | unit | `go test ./cmd/mqtt-plugin/... -run TestJSON` | ❌ W0 | ⬜ pending |
| 02-05-02 | 05 | 3 | ACQ-01 | T-02-04 | Binary payload parsing | unit | `go test ./cmd/mqtt-plugin/... -run TestBinary` | ❌ W0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `pkg/sdk/v1/plugin_test.go` — stubs for CORE-06
- [ ] `cmd/mqtt-plugin/main_test.go` — stubs for ACQ-01, ACQ-02
- [ ] `internal/validation/validate_test.go` — stubs for ACQ-03
- [ ] `internal/storage/asset_test.go` — stubs for ACQ-04
- [ ] `internal/api/assets_test.go` — stubs for ACQ-04 REST
- [ ] Makefile proto target — covers CORE-06 code generation
- [ ] `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- [ ] `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Crash isolation: killing MQTT plugin process does not crash core | ACQ-01 | Requires process kill + core health check | Start core → start MQTT plugin → kill MQTT plugin PID → verify core responds to HTTP |
| MQTT plugin auto-reconnects on embedded broker crash | ACQ-01 | Requires broker restart simulation | Start broker → connect plugin → stop broker → restart broker → verify plugin reconnects |

*If none: "All phase behaviors have automated verification."*

---

## Validation Sign-Off

- [ ] All tasks have `<automated>` verify or Wave 0 dependencies
- [ ] Sampling continuity: no 3 consecutive tasks without automated verify
- [ ] Wave 0 covers all MISSING references
- [ ] No watch-mode flags
- [ ] Feedback latency < 30s
- [ ] `nyquist_compliant: true` set in frontmatter

**Approval:** pending
