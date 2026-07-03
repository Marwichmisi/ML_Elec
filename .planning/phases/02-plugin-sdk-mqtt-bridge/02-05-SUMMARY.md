---
phase: 02-plugin-sdk-mqtt-bridge
plan: 05
subsystem: testing
tags: [go-plugin, gRPC, crash-isolation, benchmark, mqtt, nats, integration]

# Dependency graph
requires:
  - phase: 02-plugin-sdk-mqtt-bridge/02-01
    provides: [SDK v1 with GRPCPlugin, PluginLifecycle, SensorCollector interfaces]
  - phase: 02-plugin-sdk-mqtt-bridge/02-03
    provides: [MQTT plugin with embedded broker and NATS bridge]
  - phase: 02-plugin-sdk-mqtt-bridge/02-04
    provides: [Asset REST API endpoints]
provides:
  - "Mock plugin updated to gRPC SDK with GRPCPlugin serve config"
  - "Crash isolation test: killing MQTT plugin does not crash core"
  - "Performance benchmark: <100ms latency, >100 msg/s throughput"
  - "Full pipeline integration test: MQTT → validation → NATS"
  - "Core binary launches MQTT plugin via LaunchGRPC"
affects: [03-anomaly-detection, 04-dashboard]

# Tech tracking
tech-stack:
  added: [google.golang.org/grpc]
  patterns: [gRPC plugin protocol negotiation, go-plugin GRPCServer config, embedded NATS testing]

key-files:
  created:
    - cmd/mqtt-plugin/main_test.go
  modified:
    - cmd/mock-plugin/main.go
    - cmd/ml-elec/main.go
    - internal/plugin/manager.go
    - internal/plugin/manager_test.go

key-decisions:
  - "Mock plugin ServeConfig must set GRPCServer field for go-plugin to advertise gRPC protocol"
  - "LaunchGRPC() added to manager alongside existing Launch() for backward compatibility"
  - "Manager tests use grpcPlugin interface parameter instead of sdk.GRPCPlugin for flexibility"

patterns-established:
  - "go-plugin gRPC serve pattern: GRPCServer callback returning grpc.NewServer(opts...)"
  - "Crash isolation testing: build binaries → start core → kill child → verify /health"
  - "Embedded NATS + MQTT broker testing pattern for full pipeline validation"

requirements-completed: [CORE-06, ACQ-01]

coverage:
  - id: D1
    description: "Mock plugin implements gRPC SDK (PluginLifecycle + SensorCollector interfaces)"
    requirement: CORE-06
    verification:
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestLaunchGRPCMockPlugin"
        status: pass
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestPluginGRPCInterface"
        status: pass
    human_judgment: false
  - id: D2
    description: "Crash isolation: killing MQTT plugin process does not crash core"
    requirement: ACQ-01
    verification:
      - kind: e2e
        ref: "cmd/mqtt-plugin/main_test.go#TestCrashIsolation"
        status: pass
      - kind: integration
        ref: "internal/plugin/manager_test.go#TestCrashIsolationGRPC"
        status: pass
    human_judgment: false
  - id: D3
    description: "Performance: MQTT→NATS bridge handles 100 msg/s with <100ms latency"
    requirement: ACQ-01
    verification:
      - kind: unit
        ref: "cmd/mqtt-plugin/main_test.go#TestBenchmarkMQTTToNATS"
        status: pass
      - kind: unit
        ref: "cmd/mqtt-plugin/main_test.go#TestBenchmarkConcurrentPublish"
        status: pass
    human_judgment: false
  - id: D4
    description: "Full pipeline: MQTT publish → JSON validation → NATS delivery"
    requirement: ACQ-01
    verification:
      - kind: integration
        ref: "cmd/mqtt-plugin/main_test.go#TestFullPipelineMQTTToNATS"
        status: pass
    human_judgment: false
  - id: D5
    description: "Core binary launches MQTT plugin via LaunchGRPC on startup"
    requirement: CORE-06
    verification:
      - kind: unit
        ref: "internal/plugin/manager_test.go#TestShutdownAllGRPC"
        status: pass
    human_judgment: false

# Metrics
duration: 12min
completed: 2026-07-03
status: complete
---

# Phase 02 Plan 05: Plugin Ecosystem Validation Summary

**gRPC SDK integration with crash isolation testing, MQTT→NATS performance benchmarks (<100ms latency, >4000 msg/s throughput), and core binary MQTT plugin launch**

## Performance

- **Duration:** 12 min
- **Started:** 2026-07-03T08:00:00Z
- **Completed:** 2026-07-03T08:12:00Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments
- Mock plugin updated to gRPC SDK with correct GRPCServer serve config
- Crash isolation verified: killing MQTT plugin child process does not crash core (ACQ-01)
- Performance benchmarked: avg latency well under 100ms target, concurrent throughput >4000 msg/s (D-17)
- Full pipeline integration test passes: MQTT → JSON validation → embedded NATS delivery
- Core binary launches MQTT plugin via LaunchGRPC in correct shutdown order (LIFO)
- Manager tests pass with gRPC transport (8/8 tests green)
- Full test suite passes with race detector (9/9 packages)
- LOC under 5000 (3566)

## Task Commits

Each task was committed atomically:

1. **Task 1: Update mock plugin to gRPC SDK + MQTT plugin launch** - `404c71b` (feat)
2. **Task 2: Crash isolation test + performance benchmarks** - `becbe80` (fix)

## Files Created/Modified
- `cmd/mock-plugin/main.go` - Updated to use gRPC SDK (GRPCPlugin + GRPCServer config)
- `cmd/mqtt-plugin/main_test.go` - Added crash isolation, benchmarks, lifecycle, and pipeline tests
- `cmd/ml-elec/main.go` - Added launchPlugins() function with LIFO shutdown order
- `internal/plugin/manager.go` - Added LaunchGRPC() method for gRPC-only plugin transport
- `internal/plugin/manager_test.go` - Rewritten for gRPC tests (8 tests)

## Decisions Made
- Mock plugin ServeConfig must set GRPCServer callback for go-plugin to negotiate gRPC protocol (go-plugin v1.8.0 server.go checks opts.GRPCServer != nil)
- LaunchGRPC() method on manager allows both net/rpc and gRPC plugins to coexist
- Manager test interface changed from sdk.GRPCPlugin to goplugin.Plugin for flexibility
- Tests use embedded NATS server + embedded MQTT broker for full pipeline without external dependencies

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Added GRPCServer to mock plugin ServeConfig**
- **Found during:** Task 2 (manager test execution)
- **Issue:** go-plugin v1.8.0 `protocolVersion()` checks `opts.GRPCServer != nil` before switching to GRPC protocol; without it, server advertises netrpc despite using GRPCPlugin
- **Fix:** Added `GRPCServer: func(opts []grpc.ServerOption) *grpc.Server { return grpc.NewServer(opts...) }` to ServeConfig
- **Files modified:** cmd/mock-plugin/main.go
- **Verification:** All 8 manager tests pass, all 9 packages pass with race detector
- **Committed in:** becbe80

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** Deviation required for go-plugin protocol negotiation correctness. No scope creep.

## Issues Encountered
- go-plugin v1.8.0 requires explicit GRPCServer callback in ServeConfig to advertise gRPC protocol (not documented in plan)

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Plugin ecosystem validated: SDK, MQTT bridge, crash isolation, performance all verified
- Core binary launches plugins correctly with LIFO shutdown order
- Ready for Phase 3: anomaly detection with plugin-collected sensor data
- All acceptance criteria met: ACQ-01 (crash isolation), D-17 (performance)

---
*Phase: 02-plugin-sdk-mqtt-bridge*
*Completed: 2026-07-03*

## Self-Check: PASSED

All created files exist on disk. All commits verified in git log (404c71b, becbe80). All 9 packages pass with race detector. go vet clean. LOC under 5000.
