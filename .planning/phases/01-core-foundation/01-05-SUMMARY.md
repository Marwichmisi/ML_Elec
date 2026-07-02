---
phase: 01-core-foundation
plan: 05
subsystem: core
tags: [go, wire, signal-handling, graceful-shutdown, lifo, integration-tests, manual-di]

# Dependency graph
requires:
  - phase: 01-01
    provides: Go module, config package, Makefile
  - phase: 01-03
    provides: REST API server with /health and /api/v1/sensors endpoints
  - phase: 01-04
    provides: Plugin manager with Launch, Kill, ShutdownAll
provides:
  - Main entry point with manual DI, signal handling, LIFO shutdown
  - Binary that compiles, starts, serves /health, shuts down cleanly
  - Integration test suite (7 tests) with race detection
  - LOC verification (1389 < 5000)
  - ML_ELEC_PORT and CONFIG_PATH env var support
affects: [02-01, 04-dashboard]

# Tech tracking
tech-stack:
  added: [github.com/google/wire (replaced by manual DI)]
  patterns: [manual-constructor-injection, signal-notifycontext, lifo-shutdown, embedded-migrations, context-cancellation]

key-files:
  created:
    - cmd/ml-elec/main.go
    - cmd/ml-elec/main_test.go
  modified:
    - internal/nats/nats.go
    - internal/storage/storage.go
    - internal/config/config.go
    - internal/api/server.go
  deleted:
    - cmd/ml-elec/wire.go
    - cmd/ml-elec/wire_gen.go

key-decisions:
  - "Switched from Wire DI to manual constructor injection — simpler for 5 components, no code generation needed"
  - "Embedded migrations via go:embed — avoids relative path issues when running from different directories"
  - "ML_ELEC_PORT env var for test port isolation — enables parallel integration tests"
  - "CONFIG_PATH env var for config file override — enables test-specific configs"
  - "m.Close() skipped in runMigrations — it closes the underlying DB connection (SQLite driver limitation)"

patterns-established:
  - "Manual DI: InitializeApp(cfg) creates all components from config struct"
  - "LIFO shutdown: API → Plugins → NATS → Storage with 30s timeout"
  - "Signal handling: signal.NotifyContext for SIGINT/SIGTERM"
  - "Embedded migrations: go:embed + iofs for portable migration files"
  - "Integration tests: binary builds, starts, health check, graceful shutdown, concurrent load"

requirements-completed: [CORE-01, CORE-02, CORE-03, CORE-04, CORE-05, CORE-07]

coverage:
  - id: D1
    description: "Main entry point with manual DI, signal handling, LIFO shutdown, 30s timeout"
    requirement: CORE-05
    verification:
      - kind: integration
        ref: "cmd/ml-elec/main_test.go#TestFullStartup"
        status: pass
      - kind: integration
        ref: "cmd/ml-elec/main_test.go#TestGracefulShutdown"
        status: pass
      - kind: integration
        ref: "cmd/ml-elec/main_test.go#TestSignalHandlingSIGINT"
        status: pass
    human_judgment: false
  - id: D2
    description: "Binary compiles and starts without error, serves /health endpoint"
    requirement: CORE-03
    verification:
      - kind: integration
        ref: "cmd/ml-elec/main_test.go#TestBinaryCompiles"
        status: pass
      - kind: integration
        ref: "cmd/ml-elec/main_test.go#TestHealthEndpointReturnsJSON"
        status: pass
    human_judgment: false
  - id: D3
    description: "Concurrent API requests handled without data race (20/20 succeed)"
    requirement: CORE-03
    verification:
      - kind: integration
        ref: "cmd/ml-elec/main_test.go#TestConcurrentAPILoad"
        status: pass
      - kind: integration
        ref: "go test -race ./... (no races detected)"
        status: pass
    human_judgment: false
  - id: D4
    description: "Context cancellation propagates correctly to HTTP requests"
    requirement: CORE-03
    verification:
      - kind: integration
        ref: "cmd/ml-elec/main_test.go#TestContextCancellation"
        status: pass
    human_judgment: false
  - id: D5
    description: "Core stays under 5000 LOC (1389 LOC)"
    requirement: CORE-07
    verification:
      - kind: other
        ref: "make check-loc → LOC: 1389 (< 5000)"
        status: pass
    human_judgment: false

# Metrics
duration: 101min
completed: 2026-07-01
status: complete
---

# Phase 1 Plan 05: Wire DI + main.go + Integration Tests Summary

**Manual DI entry point with signal handling, LIFO shutdown (API→Plugins→NATS→Storage), 30s timeout, embedded migrations, and 7 integration tests passing with race detector — 1389 LOC**

## Performance

- **Duration:** 101 min
- **Started:** 2026-07-01T04:57:22Z
- **Completed:** 2026-07-01T06:38:48Z
- **Tasks:** 2
- **Files modified:** 6

## Accomplishments
- Main entry point with manual DI orchestrating component lifecycle (Config → NATS → Storage → Plugins → API)
- Graceful shutdown in LIFO order with 30-second timeout enforced
- Signal handling via signal.NotifyContext for SIGINT and SIGTERM
- Embedded migrations via go:embed (avoids relative path issues)
- ML_ELEC_PORT and CONFIG_PATH env vars for test isolation
- 7 integration tests: binary compiles, full startup, graceful shutdown (SIGTERM + SIGINT), concurrent API load (20 requests), JSON format verification, context cancellation
- All tests pass with race detector enabled
- LOC: 1389 (well under 5000 limit)
- go vet clean

## Task Commits

Each task was committed atomically:

1. **Task 1: Wire DI + main.go with graceful shutdown** - `305c4a9` (feat)
2. **Task 3: Integration tests + LOC verification** - `dc50c91` (feat)

## Files Created/Modified
- `cmd/ml-elec/main.go` - Entry point with manual DI, signal handling, LIFO shutdown, 30s timeout
- `cmd/ml-elec/main_test.go` - 7 integration tests with parallel execution and port isolation
- `internal/nats/nats.go` - Updated constructor to accept *config.NATSConfig
- `internal/storage/storage.go` - Embedded migrations, updated constructor to accept *config.StorageConfig
- `internal/config/config.go` - Added CONFIG_PATH env var support
- `internal/api/server.go` - Added ServerAddr() method

## Decisions Made
- Switched from Wire DI to manual constructor injection — simpler for 5 components, no code generation needed
- Embedded migrations via go:embed — avoids relative path issues when running from different directories
- ML_ELEC_PORT env var for test port isolation — enables parallel integration tests without port conflicts
- Skipped m.Close() in runMigrations — SQLite driver closes the underlying connection (known limitation)

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed migration path relative to working directory**
- **Found during:** Task 1 (binary startup test)
- **Issue:** `file://migrations` in runMigrations was relative to CWD, failing when binary runs from project root
- **Fix:** Embedded migrations via `//go:embed migrations/*.sql` and `iofs.New()`
- **Files modified:** internal/storage/storage.go
- **Verification:** Binary starts, storage verified, all tests pass
- **Committed in:** 305c4e9 (Task 1 commit)

**2. [Rule 1 - Bug] Fixed m.Close() closing underlying database**
- **Found during:** Task 1 (storage test failure)
- **Issue:** `m.Close()` in runMigrations closed the SQLite connection via the driver wrapper
- **Fix:** Removed m.Close() — Store.Close() handles database cleanup
- **Files modified:** internal/storage/storage.go
- **Verification:** TestNewOpensDBInWALMode passes, binary starts successfully
- **Committed in:** 305c4e9 (Task 1 commit)

**3. [Rule 1 - Bug] Fixed port override after Wire initialization**
- **Found during:** Task 3 (integration tests)
- **Issue:** ML_ELEC_PORT env var was set after Wire created the API server, so port override had no effect
- **Fix:** Replaced Wire DI with manual constructor injection — port override happens before API server creation
- **Files modified:** cmd/ml-elec/main.go (removed wire.go, wire_gen.go)
- **Verification:** All 7 integration tests pass with unique ports
- **Committed in:** dc50c91 (Task 3 commit)

---

**Total deviations:** 3 auto-fixed (3 bugs)
**Impact on plan:** All fixes essential for correctness. Wire→manual DI switch simplified the codebase without changing behavior. Embedded migrations make the binary portable.

## Issues Encountered
None beyond the auto-fixed deviations above.

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Core foundation complete — binary orchestrates all components
- Plugin manager ready for Phase 2 MQTT plugin integration
- API server ready for dashboard integration (Phase 4)
- All acceptance criteria from SPEC.md verified
- Phase 1 is ready for completion audit

---
*Phase: 01-core-foundation*
*Completed: 2026-07-01*

## Self-Check: PASSED

1. Key files exist: cmd/ml-elec/main.go, cmd/ml-elec/main_test.go
2. Commits exist: 305c4e9 (Task 1), dc50c91 (Task 3)
3. All 7 integration tests pass with race detector
4. go vet clean
5. LOC: 1389 (< 5000)
