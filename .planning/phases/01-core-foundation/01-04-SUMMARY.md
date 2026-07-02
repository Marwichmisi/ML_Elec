---
phase: 01-core-foundation
plan: 04
subsystem: plugin
tags: [go, plugin, hashicorp-go-plugin, json-rpc, child-process, crash-isolation]

# Dependency graph
requires:
  - phase: 01-01
    provides: Go module with dependencies, config package with PluginsConfig
provides:
  - Plugin manager with Launch, Kill, ShutdownAll, IsRunning
  - Mock plugin binary implementing Echo via JSON-RPC
  - Crash isolation: plugin crash doesn't crash manager
  - Config-based plugin enable/disable
affects: [01-05, 02-01]

# Tech tracking
tech-stack:
  added: [github.com/hashicorp/go-plugin, github.com/hashicorp/go-hclog]
  patterns: [plugin-lifecycle, child-process-isolation, json-rpc-stdin-stdout, handshake-config]

key-files:
  created:
    - internal/plugin/manager.go
    - internal/plugin/manager_test.go
    - cmd/mock-plugin/main.go
  modified: []

key-decisions:
  - "Used net/rpc (not gRPC) for plugin communication - simpler for JSON-RPC stdin/stdout"
  - "SensorPluginRPC type has Impl field for server-side implementation"
  - "HandshakeConfig with ML_ELEC_PLUGIN magic cookie for plugin identification"

patterns-established:
  - "Plugin lifecycle: Launch checks enabled list, creates go-plugin client, verifies connection"
  - "Crash isolation: go-plugin with Managed=true ensures plugin crash doesn't affect core"
  - "RPC pattern: SensorPluginRPCServer for server-side, SensorPluginRPCClient for client-side"

requirements-completed: [CORE-02]

coverage:
  - id: D1
    description: "Plugin manager with Launch, Kill, ShutdownAll, IsRunning methods"
    requirement: CORE-02
    verification:
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestNewManager"
        status: pass
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestLaunchMockPlugin"
        status: pass
    human_judgment: false
  - id: D2
    description: "Mock plugin binary implementing Echo via JSON-RPC"
    requirement: CORE-02
    verification:
      - kind: integration
        ref: "internal/plugin/manager_test.go#TestPluginEcho"
        status: pass
    human_judgment: false
  - id: D3
    description: "Crash isolation: killing plugin process doesn't crash manager"
    requirement: CORE-02
    verification:
      - kind: integration
        ref: "internal/plugin/manager_test.go#TestCrashIsolation"
        status: pass
    human_judgment: false
  - id: D4
    description: "Config-based plugin enable/disable - disabled plugin not started"
    requirement: CORE-02
    verification:
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestLaunchDisabledPlugin"
        status: pass
    human_judgment: false
  - id: D5
    description: "ShutdownAll kills all running plugins"
    requirement: CORE-02
    verification:
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestShutdownAll"
        status: pass
    human_judgment: false
  - id: D6
    description: "Invalid binary path returns error"
    requirement: CORE-02
    verification:
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestLaunchInvalidPath"
        status: pass
    human_judgment: false

# Metrics
duration: 13min
completed: 2026-07-01
status: complete
---

# Phase 1 Plan 04: Plugin Manager with Crash Isolation Summary

**Plugin manager using HashiCorp go-plugin with child process isolation, JSON-RPC stdin/stdout communication, and crash recovery tests**

## Performance

- **Duration:** 13 min
- **Started:** 2026-07-01T04:27:15Z
- **Completed:** 2026-07-01T04:41:04Z
- **Tasks:** 1
- **Files modified:** 3

## Accomplishments
- Implemented plugin manager with Launch, Kill, ShutdownAll, IsRunning methods
- Integrated HashiCorp go-plugin for child process isolation and JSON-RPC communication
- Created mock plugin binary implementing Echo via JSON-RPC for testing
- Verified crash isolation: plugin crash doesn't crash the manager
- Config-based plugin enable/disable working correctly

## Task Commits

Each task was committed atomically:

1. **Task 1: Plugin manager with go-plugin and crash isolation** - `a8f111e` (test - RED phase)
2. **Task 1: Plugin manager with go-plugin and crash isolation** - `71c959f` (feat - GREEN phase)

**Plan metadata:** `docs(01-04): complete plugin manager plan` (pending)

## Files Created/Modified
- `internal/plugin/manager.go` - Plugin manager with go-plugin integration, RPC types, lifecycle management
- `internal/plugin/manager_test.go` - 7 comprehensive tests including crash isolation and race detection
- `cmd/mock-plugin/main.go` - Mock plugin binary implementing Echo via JSON-RPC

## Decisions Made
- Used net/rpc (not gRPC) for plugin communication - simpler for JSON-RPC stdin/stdout
- SensorPluginRPC type has Impl field for server-side implementation
- HandshakeConfig with ML_ELEC_PLUGIN magic cookie for plugin identification
- TDD approach: RED phase (stub) → GREEN phase (implementation) → tests pass

## Deviations from Plan

None - plan executed exactly as written.

---

**Total deviations:** 0
**Impact on plan:** No deviations. Plan executed with TDD discipline.

## Issues Encountered
None - plan executed smoothly.

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Plugin manager complete and tested - ready for Plan 05 (main.go integration)
- Mock plugin binary available for testing
- Crash isolation verified - core will be safe from plugin failures
- Config-based enable/disable ready for Phase 02 plugin SDK

---
*Phase: 01-core-foundation*
*Completed: 2026-07-01*

## Self-Check: PASSED

1. Key files exist: internal/plugin/manager.go, internal/plugin/manager_test.go, cmd/mock-plugin/main.go
2. Commits exist: a8f111e (RED), 71c959f (GREEN)
3. Tests pass: all 7 tests pass with race detector clean
