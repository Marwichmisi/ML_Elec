---
phase: 02-plugin-sdk-mqtt-bridge
plan: 04
subsystem: api
tags: [rest, assets, sensors, pagination, http]

# Dependency graph
requires:
  - phase: 02-plugin-sdk-mqtt-bridge/02-02
    provides: [Asset CRUD storage methods, SQLite schema for assets/asset_sensors]
provides:
  - "REST endpoints for asset CRUD with pagination"
  - "Sensor registration and listing endpoints"
  - "Global sensor listing with pagination"
affects: [03-anomaly-detection, 04-dashboard]

# Tech tracking
tech-stack:
  added: []
  patterns: [pagination response {data, pagination}, method-not-allowed safety, unified handler routing]

key-files:
  created:
    - internal/api/assets.go
    - internal/api/assets_test.go
  modified:
    - internal/api/server.go
    - internal/api/sensors.go
    - internal/storage/storage.go

key-decisions:
  - "Unified SensorsHandler routes to GetSensorsHandler (with sensor_id) or ListAllSensorsHandler (global)"
  - "DELETE /api/v1/assets returns 405 per SPEC prohibition for v1 safety"
  - "CreateAsset returns 409 Conflict on duplicate name via UNIQUE constraint detection"

patterns-established:
  - "Pagination response pattern: {data: [...], pagination: {page, limit, total}}"
  - "Unified handler pattern: single route with conditional behavior based on query params"

requirements-completed: [ACQ-04]

coverage:
  - id: D1
    description: "POST /api/v1/assets creates asset with 201, duplicate returns 409, invalid returns 400"
    requirement: ACQ-04
    verification:
      - kind: unit
        ref: "internal/api/assets_test.go#TestCreateAssetValidRequest"
        status: pass
      - kind: unit
        ref: "internal/api/assets_test.go#TestCreateAssetDuplicateName"
        status: pass
      - kind: unit
        ref: "internal/api/assets_test.go#TestCreateAssetInvalidBody"
        status: pass
    human_judgment: false
  - id: D2
    description: "GET /api/v1/assets returns paginated list with {data, pagination} structure"
    requirement: ACQ-04
    verification:
      - kind: unit
        ref: "internal/api/assets_test.go#TestListAssetsPagination"
        status: pass
    human_judgment: false
  - id: D3
    description: "GET /api/v1/assets/{id} returns single asset, 404 for nonexistent"
    requirement: ACQ-04
    verification:
      - kind: unit
        ref: "internal/api/assets_test.go#TestGetAssetByID"
        status: pass
      - kind: unit
        ref: "internal/api/assets_test.go#TestGetAssetNotFound"
        status: pass
    human_judgment: false
  - id: D4
    description: "Sensor CRUD: POST/GET for asset sensors, global sensor listing with pagination"
    requirement: ACQ-04
    verification:
      - kind: unit
        ref: "internal/api/assets_test.go#TestCreateAssetSensor"
        status: pass
      - kind: unit
        ref: "internal/api/assets_test.go#TestListAssetSensors"
        status: pass
      - kind: unit
        ref: "internal/api/assets_test.go#TestListAllSensors"
        status: pass
    human_judgment: false
  - id: D5
    description: "DELETE /api/v1/assets returns 405 Method Not Allowed per SPEC prohibition"
    requirement: ACQ-04
    verification:
      - kind: unit
        ref: "internal/api/assets_test.go#TestDeleteAssetsReturns405"
        status: pass
    human_judgment: false

# Metrics
duration: 6min
completed: 2026-07-03
status: complete
---

# Phase 02 Plan 04: Asset Registry REST API Summary

**REST endpoints for asset CRUD with pagination (201/409/400/404/405), sensor registration, and global sensor listing using unified SensorsHandler routing**

## Performance

- **Duration:** 6 min
- **Started:** 2026-07-03T06:25:59Z
- **Completed:** 2026-07-03T06:32:47Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments
- Created 7 REST endpoints for asset and sensor management with proper HTTP status codes
- POST /api/v1/assets with 201 Created, 409 Conflict (duplicate), 400 Bad Request (invalid/missing)
- GET /api/v1/assets with pagination returning {data, pagination} structure per D-23
- GET /api/v1/assets/{id} with 404 for nonexistent assets
- Sensor CRUD: POST/GET for asset sensors, global sensor listing
- DELETE /api/v1/assets returns 405 Method Not Allowed per SPEC prohibition
- Unified SensorsHandler routes by sensor_id query param (backward compatible)
- Added ListAllSensors storage method with pagination
- 19 API tests passing with race detection

## Task Commits

Each task was committed atomically:

1. **Task 1: Create asset REST endpoints with pagination (TDD RED)** - `f5c817e` (test)
2. **Task 1: Create asset REST endpoints with pagination (TDD GREEN)** - `997f6af` (feat)

## Files Created/Modified
- `internal/api/assets.go` - Asset and sensor REST handlers with pagination, error handling, and 405 safety
- `internal/api/assets_test.go` - 13 comprehensive tests covering all endpoints and edge cases
- `internal/api/server.go` - Registered 7 new routes for asset/sensor endpoints
- `internal/api/sensors.go` - Added unified SensorsHandler routing by sensor_id param
- `internal/storage/storage.go` - Added ListAllSensors method; updated NewForTest with asset tables

## Decisions Made
- Unified SensorsHandler routes to GetSensorsHandler (with sensor_id) or ListAllSensorsHandler (global) for backward compatibility
- DELETE /api/v1/assets returns 405 per SPEC prohibition — safety for v1
- CreateAsset detects UNIQUE constraint violations for 409 Conflict response
- Updated existing tests to match new global sensor listing behavior (GET /api/v1/sensors without sensor_id returns all sensors)

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Updated old tests for new API contract**
- **Found during:** Task 1 (GREEN phase)
- **Issue:** Old tests expected 400 for GET /api/v1/sensors without sensor_id, but new API returns all sensors globally
- **Fix:** Updated TestGetSensorsNoSensorID and TestJSONResponseFormat to test new behavior
- **Files modified:** internal/api/api_test.go
- **Verification:** All 19 tests pass with race detection
- **Committed in:** 997f6af

**2. [Rule 2 - Missing Critical] Added asset tables to NewForTest**
- **Found during:** Task 1 (GREEN phase)
- **Issue:** Test database only had sensor_readings table, asset tests failed with 500
- **Fix:** Added assets and asset_sensors table creation to NewForTest
- **Files modified:** internal/storage/storage.go
- **Verification:** All asset API tests pass
- **Committed in:** 997f6af

---

**Total deviations:** 2 auto-fixed (1 bug, 1 missing critical)
**Impact on plan:** Both deviations necessary for correctness. Old test updates are expected API evolution. No scope creep.

## Issues Encountered
None

## User Setup Required
None - no external service configuration required.

## Next Phase Readiness
- Asset registry REST API complete with full CRUD and pagination
- Ready for MQTT plugin to register assets on startup (Plan 3)
- Dashboard can query assets and sensors via REST API (Phase 4)

---
*Phase: 02-plugin-sdk-mqtt-bridge*
*Completed: 2026-07-03*

## Self-Check: PASSED

All created files exist on disk. All commits verified in git log. All 19 API tests pass with race detection.
