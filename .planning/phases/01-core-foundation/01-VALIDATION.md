---
phase: 01
slug: core-foundation
status: draft
nyquist_compliant: false
wave_0_complete: false
created: 2026-07-01
---

# Phase 01 — Validation Strategy

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
| 01-01-01 | 01 | 1 | CORE-01 | — | NATS lifecycle isolation | integration | `go test ./internal/nats/...` | ❌ Wave 0 | ⬜ pending |
| 01-02-01 | 01 | 1 | CORE-02 | T-01-01 | Plugin crash isolation via child process | integration | `go test ./internal/plugin/...` | ❌ Wave 0 | ⬜ pending |
| 01-03-01 | 01 | 1 | CORE-03 | — | Race-free concurrent requests | unit | `go test ./internal/api/...` | ❌ Wave 0 | ⬜ pending |
| 01-04-01 | 01 | 1 | CORE-04 | T-01-02 | Parameterized queries (SQL injection prevention) | unit | `go test ./internal/storage/...` | ❌ Wave 0 | ⬜ pending |
| 01-05-01 | 01 | 1 | CORE-05 | T-01-03 | Config file integrity | unit | `go test ./internal/config/...` | ❌ Wave 0 | ⬜ pending |
| 01-06-01 | 01 | 1 | CORE-07 | — | N/A | manual | `find . -name '*.go' ! -name '*_test.go' ! -path './generated/*' \| xargs wc -l` | ❌ Wave 0 | ⬜ pending |

*Status: ⬜ pending · ✅ green · ❌ red · ⚠️ flaky*

---

## Wave 0 Requirements

- [ ] `internal/nats/nats_test.go` — stubs for CORE-01
- [ ] `internal/plugin/manager_test.go` — stubs for CORE-02
- [ ] `internal/api/api_test.go` — stubs for CORE-03
- [ ] `internal/storage/storage_test.go` — stubs for CORE-04
- [ ] `internal/config/config_test.go` — stubs for CORE-05
- [ ] `cmd/ml-elec/main_test.go` — covers CORE-07 (LOC check)

*If none: "Existing infrastructure covers all phase requirements."*

---

## Manual-Only Verifications

| Behavior | Requirement | Why Manual | Test Instructions |
|----------|-------------|------------|-------------------|
| Core < 5000 LOC | CORE-07 | Size constraint verification | `find . -name '*.go' ! -name '*_test.go' ! -path './generated/*' \| xargs wc -l` |

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
