# Phase 1: Critical Foundation -- Task List

**Status:** Not started
**Blocks:** All Phase 2 and Phase 3 work
**Reference:** [PROJECT_PLAN.md](PROJECT_PLAN.md)

---

## Pre-flight

- [ ] Tag current state: `git tag -a v0.1.0-pre-upgrade -m "Pre-upgrade baseline"`
- [ ] Verify Go version: `go version` (expect 1.24+)
- [ ] Verify current tests pass: `go test ./... -v`
- [ ] Verify latest dependency versions at pkg.go.dev (Nov 2025 references may be stale)
- [ ] Check `server-go-ssp` compatibility with GORM v2 (import conflict risk -- see Gaps in PROJECT_PLAN.md)

---

## Stage 1.1: GORM v2 Migration

**Branch:** `stage-1-gorm-v2`
**Risk:** MEDIUM (breaking API changes)

- [ ] **TASK-001** Update `go.mod` with GORM v2 dependencies
  - Remove `github.com/jinzhu/gorm`
  - Add `gorm.io/gorm@latest`
  - Run `go mod tidy`
  - Verify: `go mod graph | grep gorm`

- [ ] **TASK-002** Update `auth_store.go` imports
  - Change `"github.com/jinzhu/gorm"` to `"gorm.io/gorm"`
  - Add `"errors"` import

- [ ] **TASK-003** Update error handling in `FindIdentity()`
  - Replace `gorm.IsRecordNotFoundError(err)` with `errors.Is(err, gorm.ErrRecordNotFound)`

- [ ] **TASK-004** Update `auth_store_test.go` for GORM v2
  - Update import path
  - Change `gorm.Open("postgres", dsn)` to `gorm.Open(postgres.Open(dsn), &gorm.Config{})`

- [ ] **TASK-005** Create GORM v2 integration test
  - New file: `auth_store_integration_test.go`
  - CRUD round-trip test on PostgreSQL
  - Build-tag gated: `//go:build integration`

- [ ] **TASK-006** Run full test suite
  - `go test ./... -v -race`
  - All tests must pass with race detection

- [ ] **TASK-007** Security scan after GORM upgrade
  - `gosec ./...` -- zero new issues
  - `govulncheck ./...` -- zero new vulnerabilities

- [ ] **TASK-008** Code review and PR creation
  - Branch PR to main
  - 2+ reviewers

- [ ] **TASK-009** Merge to main and tag `v0.2.0-stage1`

**Stage 1.1 gate:** All unit + integration tests pass, security scan clean,
PR approved.

---

## Stage 1.2: Database Driver Upgrades

**Branch:** `stage-2-db-drivers`
**Depends on:** Stage 1.1 merged
**Risk:** LOW

- [ ] **TASK-010** Upgrade `github.com/lib/pq` to latest
  - Fixes CVE-2021-3121 (HIGH)
  - `go get github.com/lib/pq@latest`

- [ ] **TASK-011** Upgrade `github.com/mattn/go-sqlite3` to latest
  - `go get github.com/mattn/go-sqlite3@latest`

- [ ] **TASK-012** Upgrade `github.com/go-sql-driver/mysql` to latest
  - `go get github.com/go-sql-driver/mysql@latest`

- [ ] **TASK-013** Upgrade `github.com/denisenkom/go-mssqldb` to latest
  - `go get github.com/denisenkom/go-mssqldb@latest`

- [ ] **TASK-014** Update GORM driver packages
  - `go get -u gorm.io/driver/postgres gorm.io/driver/mysql gorm.io/driver/sqlite gorm.io/driver/sqlserver`
  - `go mod tidy`

- [ ] **TASK-015** Integration tests (all databases)
  - PostgreSQL: Docker `postgres:16`
  - MySQL: Docker `mysql:8.4`
  - SQLite: in-memory `:memory:`

- [ ] **TASK-016** Security scan and merge
  - Tag `v0.2.0-stage2`

**Stage 1.2 gate:** Integration tests pass on PostgreSQL, MySQL, SQLite. No
new vulnerabilities.

---

## Stage 1.3: Transitive Dependencies

**Branch:** `stage-3-transitive`
**Depends on:** Stage 1.2 merged
**Risk:** LOW

- [ ] **TASK-017** Update all golang.org/x packages
  - `go get -u golang.org/x/sys golang.org/x/text golang.org/x/crypto`

- [ ] **TASK-018** Review and update all indirect dependencies
  - `go get -u ./...`
  - `go mod tidy`
  - Review `git diff go.mod` for unexpected changes

- [ ] **TASK-019** Full test suite with race detection
  - `go test ./... -v -race`

- [ ] **TASK-020** Merge and tag `v0.2.0-stage3`

**Stage 1.3 gate:** All deps current, race detector clean, all tests pass.

---

## Phase 1 Completion Criteria

- [ ] All 20 tasks (TASK-001 to TASK-020) complete
- [ ] GORM v2 (`gorm.io/gorm`) replacing deprecated `jinzhu/gorm`
- [ ] All database drivers at latest secure versions
- [ ] All transitive dependencies current
- [ ] Zero high/critical security vulnerabilities
- [ ] All tests pass with `-race` flag
- [ ] Three stage tags created: `v0.2.0-stage1`, `v0.2.0-stage2`, `v0.2.0-stage3`

When Phase 1 is complete, update:
1. `README.md` -- Phase 1 status to "Complete", GORM Version to current
2. `docs/PROJECT_PLAN.md` -- mark tasks with completion dates
