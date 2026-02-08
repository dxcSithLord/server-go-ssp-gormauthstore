# Task Register

**Source of truth** for all implementation tasks across all phases.
Referenced by [PROJECT_PLAN.md](PROJECT_PLAN.md).

> Task IDs are authoritative here. If a stage-level document diverges,
> update it to match this register.

---

## How to read this file

| Column | Meaning |
|--------|---------|
| **ID** | Unique task identifier (TASK-NNN) |
| **Phase.Stage** | Phase and stage the task belongs to |
| **Description** | What the task does |
| **Status** | `done`, `pending`, or `deferred` |
| **Date** | Completion date (blank if pending) |
| **Depends on** | Tasks that must complete before this one |
| **Blocks** | Tasks that cannot start until this one completes |

---

## Phase 1: Critical Foundation

### Stage 1.1 -- GORM v2 Migration

| ID | Description | Status | Date | Depends on | Blocks |
|----|-------------|--------|------|------------|--------|
| TASK-001 | Update `go.mod`: remove `jinzhu/gorm`, add `gorm.io/gorm`, switch to `dxcSithLord/server-go-ssp` | done | 2026-02-06 | -- | TASK-002, TASK-003, TASK-004 |
| TASK-002 | Update `auth_store.go` imports; create `identityRecord` model with GORM v2 tags | done | 2026-02-06 | TASK-001 | TASK-003, TASK-004, TASK-005 |
| TASK-003 | Update `FindIdentity` error handling: `errors.Is(err, gorm.ErrRecordNotFound)` | done | 2026-02-06 | TASK-002 | TASK-004, TASK-005 |
| TASK-004 | Update `auth_store_test.go` for GORM v2 + SQLite in-memory | done | 2026-02-06 | TASK-001, TASK-002, TASK-003 | TASK-005, TASK-006 |
| TASK-005 | Create `auth_store_integration_test.go` (CRUD round-trip, build-tag gated) | done | 2026-02-06 | TASK-004 | TASK-006 |
| TASK-006 | Run full test suite with `-race` | done | 2026-02-06 | TASK-005 | TASK-007, TASK-008 |
| TASK-007 | Security scan (`gosec`, `govulncheck`) | done | 2026-02-06 | TASK-006 | TASK-008 |
| TASK-008 | Code review and PR creation | done | 2026-02-06 | TASK-006, TASK-007 | TASK-009 |
| TASK-009 | Merge to main; tag `v0.2.0-stage1` | done | 2026-02-06 | TASK-008 | TASK-010, TASK-011, TASK-012, TASK-013 |

### Stage 1.2 -- Database Driver Upgrades

| ID | Description | Status | Date | Depends on | Blocks |
|----|-------------|--------|------|------------|--------|
| TASK-010 | Upgrade `github.com/lib/pq` to latest | done | 2026-02-06 | TASK-009 | TASK-014 |
| TASK-011 | Upgrade `github.com/mattn/go-sqlite3` to latest | done | 2026-02-06 | TASK-009 | TASK-014 |
| TASK-012 | Upgrade `github.com/go-sql-driver/mysql` to latest | done | 2026-02-06 | TASK-009 | TASK-014 |
| TASK-013 | Upgrade `github.com/denisenkom/go-mssqldb` to latest | done | 2026-02-06 | TASK-009 | TASK-014 |
| TASK-014 | Update GORM driver packages (`gorm.io/driver/*`) to latest | done | 2026-02-06 | TASK-010, TASK-011, TASK-012, TASK-013 | TASK-015 |
| TASK-015 | Integration tests (all databases) | deferred | | TASK-014 | TASK-016 |
| TASK-016 | Security scan and merge; tag `v0.2.0-stage2` | done | 2026-02-06 | TASK-015 | TASK-017 |

> **TASK-015 note:** SQLite in-memory tests pass. PostgreSQL and MySQL
> integration tests require Docker and are deferred to the CI environment.

### Stage 1.3 -- Transitive Dependencies

| ID | Description | Status | Date | Depends on | Blocks |
|----|-------------|--------|------|------------|--------|
| TASK-017 | Update `golang.org/x/*` packages (crypto v0.47.0, text, image, sync) | done | 2026-02-06 | TASK-016 | TASK-018 |
| TASK-018 | Review and update all indirect dependencies | done | 2026-02-06 | TASK-017 | TASK-019 |
| TASK-019 | Full test suite with race detection | done | 2026-02-06 | TASK-018 | TASK-020 |
| TASK-020 | Merge and tag `v0.2.0-stage3` | done | 2026-02-06 | TASK-019 | TASK-021, TASK-022, TASK-023, TASK-024 |

---

## Phase 2: Security and Testing

### Stage 2.1 -- Security Integration

| ID | Description | Status | Date | Depends on | Blocks |
|----|-------------|--------|------|------------|--------|
| TASK-021 | Integrate `ValidateIdk` into `FindIdentity` | done | 2026-02-06 | TASK-020 | TASK-024, TASK-025 |
| TASK-022 | Integrate `ValidateIdk` into `SaveIdentity` (+ nil check) | done | 2026-02-06 | TASK-020 | TASK-025 |
| TASK-023 | Integrate `ValidateIdk` into `DeleteIdentity` | done | 2026-02-06 | TASK-020 | TASK-025 |
| TASK-024 | Implement `FindIdentitySecure` helper (returns `SecureIdentityWrapper`) | done | 2026-02-06 | TASK-021 | TASK-025 |
| TASK-025 | Security test suite (`auth_store_security_test.go`) | done | 2026-02-06 | TASK-021, TASK-022, TASK-023, TASK-024 | TASK-026 |
| TASK-026 | Security scan (`gosec` + `govulncheck`) | done | 2026-02-07 | TASK-025 | TASK-027 |
| TASK-027 | Merge and tag `v0.2.0-stage4` | done | 2026-02-07 | TASK-026 | TASK-028, TASK-029, TASK-030, TASK-031 |

> **TASK-021/022/023 note:** Completed early as part of Phase 1 work.
> ValidateIdk is now called before every database operation.

### Stage 2.2 -- Comprehensive Test Suite

| ID | Description | Status | Date | Depends on | Blocks |
|----|-------------|--------|------|------------|--------|
| TASK-028 | Unit tests (27 cases per `API_TESTS_SPEC`) | done | 2026-02-07 | TASK-027 | TASK-032 |
| TASK-029 | Integration tests (10 cases, multi-DB + concurrency) | done | 2026-02-07 | TASK-027 | TASK-032 |
| TASK-030 | Benchmarks (6 cases, PERF-001 to PERF-006) | done | 2026-02-07 | TASK-027 | TASK-032 |
| TASK-031 | Test data helpers (`testIdentityBuilder` pattern) | done | 2026-02-07 | TASK-027 | TASK-028, TASK-029 |
| TASK-032 | Measure and verify coverage (98.8% overall, 90.9%+ `auth_store.go`) | done | 2026-02-07 | TASK-028, TASK-029, TASK-030 | TASK-033 |
| TASK-033 | CI test workflow (coverage gate, Go 1.24 update) | done | 2026-02-07 | TASK-032 | TASK-034 |
| TASK-034 | Merge and tag `v0.3.0-rc1` | done | 2026-02-08 | TASK-032, TASK-033 | TASK-035, TASK-036, TASK-037 |

---

## Phase 3: Production Readiness and Release

### Stage 3.1 -- Production Hardening

| ID | Description | Status | Date | Depends on | Blocks |
|----|-------------|--------|------|------------|--------|
| TASK-035 | Context support (DP-003 resolved: Option A) | done | 2026-02-07 | TASK-034 | TASK-038 |
| TASK-036 | Production deployment documentation (`docs/PRODUCTION.md`) | done | 2026-02-07 | TASK-034 | TASK-038 |
| TASK-037 | Migration guide (`docs/UPGRADE_FROM_V0.md`) | done | 2026-02-07 | TASK-034 | TASK-038 |
| TASK-038 | Final security audit (gosec clean, govulncheck network-blocked) | done | 2026-02-07 | TASK-035, TASK-036, TASK-037 | TASK-039, TASK-040 |

### Stage 3.2 -- Release v1.0.0

| ID | Description | Status | Date | Depends on | Blocks |
|----|-------------|--------|------|------------|--------|
| TASK-039 | Update `README.md` (context support, status tables, docs links) | done | 2026-02-07 | TASK-038 | TASK-041 |
| TASK-040 | Create `CHANGELOG.md` | done | 2026-02-07 | TASK-038 | TASK-041 |
| TASK-041 | Tag `v1.0.0` | pending | | TASK-039, TASK-040 | TASK-042, TASK-043, TASK-044 |
| TASK-042 | GitHub Release (release page with changelog) | pending | | TASK-041 | -- |
| TASK-043 | Revert module path to `sqrldev` | pending | | TASK-041 | TASK-044 |
| TASK-044 | Submit to pkg.go.dev | pending | | TASK-041, TASK-043 | -- |

---

## Summary

| Phase | Stage | Tasks | Done | Pending | Deferred |
|-------|-------|-------|------|---------|----------|
| 1 | 1.1 GORM v2 Migration | 9 | 9 | 0 | 0 |
| 1 | 1.2 Database Drivers | 7 | 6 | 0 | 1 |
| 1 | 1.3 Transitive Deps | 4 | 4 | 0 | 0 |
| 2 | 2.1 Security Integration | 7 | 7 | 0 | 0 |
| 2 | 2.2 Comprehensive Tests | 7 | 7 | 0 | 0 |
| 3 | 3.1 Production Hardening | 4 | 4 | 0 | 0 |
| 3 | 3.2 Release v1.0.0 | 6 | 2 | 4 | 0 |
| **Total** | | **44** | **39** | **4** | **1** |
