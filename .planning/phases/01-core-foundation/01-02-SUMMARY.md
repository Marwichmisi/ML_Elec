---
phase: 01-core-foundation
plan: 02
subsystem: storage,nats
tags: [go, sqlite, wal, nats, embedded, pubsub]

# Dependency graph
requires:
  - phase: 01-core-foundation
    provides: Go module with all dependencies initialized
provides:
  - SQLite storage with WAL mode, migrations, and CRUD operations
  - Embedded NATS server wrapper with pub/sub support
  - Parameterized queries via squirrel
  - Corruption detection via PRAGMA integrity_check
affects: [01-03, 01-04, 01-05]

# Tech tracking
tech-stack:
  added: [modernc.org/sqlite, github.com/Masterminds/squirrel, github.com/golang-migrate/migrate/v4, github.com/nats-io/nats-server/v2, github.com/nats-io/nats.go]
  patterns: [sqlite-wal, embedded-nats, query-builder, migration-files, tdd-red-green]

key-files:
  created:
    - internal/storage/storage.go
    - internal/storage/migrations/001_init.up.sql
    - internal/storage/migrations/001_init.down.sql
    - internal/storage/storage_test.go
    - internal/nats/nats.go
    - internal/nats/nats_test.go
  modified:
    - go.mod
    - go.sum

key-decisions:
  - "Used modernc.org/sqlite for pure Go SQLite (no CGO)"
  - "Embedded NATS with Port: -1 for random port selection"
  - "squirrel for parameterized queries (no raw SQL strings)"
  - "golang-migrate for schema versioning"

patterns-established:
  - "SQLite WAL: _pragma=journal_mode(WAL) in DSN"
  - "Embedded NATS: server.NewServer(opts) + ReadyForConnections timeout"
  - "TDD: RED (failing tests) → GREEN (implementation) → commit"

requirements-completed: [CORE-01, CORE-04]

coverage:
  - id: D1
    description: "SQLite storage with WAL mode, migrations, CRUD, and corruption check"
    requirement: CORE-04
    verification:
      - kind: unit
        ref: "internal/storage/storage_test.go#TestNewOpensDBInWALMode"
        status: pass
      - kind: unit
        ref: "internal/storage/storage_test.go#TestInsertAndGetSensors"
        status: pass
      - kind: unit
        ref: "internal/storage/storage_test.go#TestCheckCorruptionReturnsNilForHealthyDB"
        status: pass
      - kind: unit
        ref: "internal/storage/storage_test.go#TestCheckCorruptionReturnsErrorForCorruptedDB"
        status: pass
    human_judgment: false
  - id: D2
    description: "Embedded NATS server with lifecycle management and pub/sub"
    requirement: CORE-01
    verification:
      - kind: integration
        ref: "internal/nats/nats_test.go#TestPublishSubscribe"
        status: pass
      - kind: integration
        ref: "internal/nats/nats_test.go#TestConcurrentPubSub"
        status: pass
    human_judgment: false

# Metrics
duration: 49min
completed: 2026-07-01
status: complete
---

# Phase 1 Plan 02: SQLite Storage & NATS Embedded Summary

**SQLite WAL storage with squirrel queries and migrations, plus embedded NATS server wrapper with pub/sub integration — both with full test suites**

## Performance

- **Duration:** 49 min
- **Started:** 2026-07-01T03:34:03Z
- **Completed:** 2026-07-01T04:23:24Z
- **Tasks:** 2
- **Files modified:** 8

## Accomplishments
- SQLite storage package: WAL mode verified via PRAGMA, CRUD via squirrel, corruption check, migrations
- NATS server wrapper: embedded server with random port, pub/sub, concurrent access tested
- All tests pass with race detector enabled
- TDD discipline: RED (failing tests) → GREEN (implementation) for both packages

## Task Commits

Each task was committed atomically:

1. **Task 1: SQLite storage with WAL mode** - `5d3fcf5` (test) + `bdc99aa` (feat)
2. **Task 2: NATS server wrapper** - `5668878` (test) + `f9cf2b6` (feat)

## Files Created/Modified
- `internal/storage/storage.go` - SQLite WAL store with CRUD and corruption check
- `internal/storage/migrations/001_init.up.sql` - sensor_readings table schema
- `internal/storage/migrations/001_init.down.sql` - Rollback migration
- `internal/storage/storage_test.go` - 7 tests covering WAL, CRUD, corruption, limits
- `internal/nats/nats.go` - Embedded NATS server wrapper
- `internal/nats/nats_test.go` - 6 tests covering lifecycle, pub/sub, concurrency
- `go.mod` - Added dependencies (squirrel, migrate, nats-server, nats.go, sqlite)
- `go.sum` - Dependency checksums

## Decisions Made
- Used modernc.org/sqlite for pure Go SQLite (no CGO, cross-compilation)
- Embedded NATS with Port: -1 for random port selection (no network exposure)
- squirrel for parameterized queries (security: no raw SQL strings)
- golang-migrate for schema versioning (file-based migrations)

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed Ready() method returning false after Start()**
- **Found during:** Task 2 (NATS server tests)
- **Issue:** ReadyForConnections(0) returns false when timeout is 0, causing TestStartMakesServerReady to fail
- **Fix:** Changed Ready() to return s.ready flag instead of calling ReadyForConnections(0)
- **Files modified:** internal/nats/nats.go
- **Verification:** All 6 NATS tests pass
- **Committed in:** f9cf2b6 (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (1 bug)
**Impact on plan:** Minor fix to Ready() logic. No scope creep.

## Issues Encountered
- Network timeout issues during dependency download (modernc.org/sqlite) — resolved by retrying with GOFLAGS=-mod=mod
- go.mod file reverted multiple times during dependency installation — resolved by manually editing and running go mod tidy

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Storage package ready for API layer (Plan 03)
- NATS package ready for internal communication (Plan 03)
- Both packages follow Go error handling patterns (wrapped errors, context propagation)

---
*Phase: 01-core-foundation*
*Completed: 2026-07-01*
