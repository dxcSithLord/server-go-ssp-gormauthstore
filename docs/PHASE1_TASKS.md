# Phase 1: Critical Foundation -- Task List

**Status:** In progress (Stages 1.1-1.3 code complete, pending PR/merge/tag)
**Blocks:** All Phase 2 and Phase 3 work
**Reference:** [PROJECT_PLAN.md](PROJECT_PLAN.md) | [TASKS.md](TASKS.md) (authoritative task register)

---

## Pre-flight

- [x] Verify Go version: `go version` (1.24.7) -- 2026-02-06
- [x] Verify current tests pass: `go test ./... -v` -- 2026-02-06
- [x] Verify latest dependency versions at pkg.go.dev -- 2026-02-06
- [x] Check `server-go-ssp` compatibility with GORM v2 -- 2026-02-06
  - Resolved: upstream uses GORM v1 `sql:""` tags; created `identityRecord`
    wrapper with GORM v2 `gorm:""` tags to bridge the gap.

---

## Stage 1.1: GORM v2 Migration

**Branch:** `claude/implement-phase1-tasks-P5Ja3`
**Risk:** MEDIUM (breaking API changes)

- [x] **TASK-001** Update `go.mod` with GORM v2 dependencies -- 2026-02-06
  - Removed `github.com/jinzhu/gorm`
  - Added `gorm.io/gorm` v1.31.1
  - Switched ssp dependency to `github.com/dxcSithLord/server-go-ssp`

- [x] **TASK-002** Update `auth_store.go` imports -- 2026-02-06
  - Changed `"github.com/jinzhu/gorm"` to `"gorm.io/gorm"`
  - Added `"errors"` import
  - Created `identityRecord` model with GORM v2 tags

- [x] **TASK-003** Update error handling in `FindIdentity()` -- 2026-02-06
  - Replaced `gorm.IsRecordNotFoundError(err)` with `errors.Is(err, gorm.ErrRecordNotFound)`

- [x] **TASK-004** Update `auth_store_test.go` for GORM v2 -- 2026-02-06
  - Updated import paths to GORM v2 + dxcSithLord fork
  - Switched from PostgreSQL to SQLite in-memory for unit tests
  - Added `openTestDB()` helper

- [x] **TASK-005** Create GORM v2 integration test -- 2026-02-06
  - New file: `auth_store_integration_test.go`
  - CRUD round-trip, multi-record, all-fields, not-found, delete tests
  - Build-tag gated: `//go:build integration`

- [x] **TASK-006** Run full test suite -- 2026-02-06
  - `go test -tags=integration ./... -v -race` -- all pass

- [x] **TASK-007** Security scan after GORM upgrade -- 2026-02-06
  - `gosec ./...` -- zero issues
  - `govulncheck` -- blocked by network (vuln.go.dev inaccessible)
  - **TODO:** Re-run `govulncheck ./...` once network access to vuln.go.dev is
    available to verify that the `golang.org/x/crypto` upgrade from v0.43.0
    (CVE-2025-47914, CVE-2025-58181) to v0.47.0 resolves all known
    vulnerabilities. Track via CI pipeline or manual verification.

- [ ] **TASK-008** Code review and PR creation
  - Branch PR to main
  - 2+ reviewers

- [ ] **TASK-009** Merge to main and tag `v0.2.0-stage1`

**Stage 1.1 gate:** All unit + integration tests pass, security scan clean,
PR approved.

---

## Stage 1.2: Database Driver Upgrades

**Branch:** `claude/implement-phase1-tasks-P5Ja3`
**Depends on:** Stage 1.1
**Risk:** LOW

- [x] **TASK-010** Upgrade `github.com/lib/pq` to latest -- 2026-02-06
  - Upgraded to v1.11.1

- [x] **TASK-011** Upgrade `github.com/mattn/go-sqlite3` to latest -- 2026-02-06
  - Upgraded to v1.14.33

- [x] **TASK-012** Upgrade `github.com/go-sql-driver/mysql` to latest -- 2026-02-06
  - Available via `gorm.io/driver/mysql` v1.6.0 (consumer-facing dependency)

- [x] **TASK-013** Upgrade `github.com/denisenkom/go-mssqldb` to latest -- 2026-02-06
  - Available via `gorm.io/driver/sqlserver` v1.6.3 (consumer-facing dependency)

- [x] **TASK-014** Update GORM driver packages (`gorm.io/driver/*`) to latest -- 2026-02-06
  - `gorm.io/driver/postgres` v1.6.0, `gorm.io/driver/mysql` v1.6.0,
    `gorm.io/driver/sqlite` v1.6.0, `gorm.io/driver/sqlserver` v1.6.3
  - NOTE: Only sqlite is a direct dependency; others are consumer-facing

- [ ] **TASK-015** Integration tests (all databases) -- deferred
  - SQLite in-memory `:memory:` -- all pass
  - PostgreSQL and MySQL require Docker; deferred to CI environment

- [ ] **TASK-016** Security scan and merge; tag `v0.2.0-stage2`
  - Part of overall PR

**Stage 1.2 gate:** Integration tests pass on SQLite. No new vulnerabilities.

---

## Stage 1.3: Transitive Dependencies

**Branch:** `claude/implement-phase1-tasks-P5Ja3`
**Depends on:** Stage 1.2
**Risk:** LOW

- [x] **TASK-017** Update all golang.org/x packages -- 2026-02-06
  - `golang.org/x/crypto` v0.47.0 (was v0.43.0 with known CVEs)
  - `golang.org/x/text` v0.33.0
  - `golang.org/x/image` v0.35.0
  - `golang.org/x/sync` v0.19.0

- [x] **TASK-018** Review and update all indirect dependencies -- 2026-02-06
  - `go get -u ./... && go mod tidy` -- clean

- [x] **TASK-019** Full test suite with race detection -- 2026-02-06
  - `go test -tags=integration ./... -v -race` -- all pass

- [ ] **TASK-020** Merge and tag `v0.2.0-stage3`

**Stage 1.3 gate:** All deps current, race detector clean, all tests pass.

---

## Phase 1 Completion Criteria

- [x] GORM v2 (`gorm.io/gorm`) replacing deprecated `jinzhu/gorm`
- [x] Database drivers at latest secure versions
- [x] All transitive dependencies current
- [x] Zero gosec issues
- [x] All tests pass with `-race` flag
- [ ] PR reviewed and merged (TASK-008, TASK-009)
- [ ] Stage tags created: `v0.2.0-stage1`, `v0.2.0-stage2`, `v0.2.0-stage3`

When Phase 1 is complete, update:
1. `README.md` -- Phase 1 status to "Complete", GORM Version to current
2. `docs/PROJECT_PLAN.md` -- mark tasks with completion dates
