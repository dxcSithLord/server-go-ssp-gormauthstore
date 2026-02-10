# Project Plan

## Consolidated Implementation Plan for gormauthstore

**Version:** 3.3 (Updated February 8, 2026)
**Date:** February 8, 2026
**Status:** Phase 1 complete, Phase 2 complete, Phase 3 Stage 3.1 complete

---

## Current State

| Metric | Value | Notes |
|--------|-------|-------|
| **GORM Version** | v2 (gorm.io/gorm v1.31.1) | Migrated from deprecated jinzhu/gorm v1.9.16 |
| **Go Version** | 1.25.7 | go.mod min 1.25.0, toolchain go1.25.7 |
| **Test Coverage** | 100% | 100 tests (90 default + 10 integration), 10 benchmarks. Target: 70%+ |
| **Security Scans** | CI/CD configured | gosec clean, golangci-lint v2.8.0, 14 security tests |
| **Secure Memory** | Implemented | ClearIdentity, WipeBytes, SecureIdentityWrapper available |
| **Input Validation** | Integrated | ValidateIdk() called by FindIdentity, SaveIdentity, DeleteIdentity |

### Completed (not repeated in tasks below)

- Go module initialization (`go.mod`, `go.sum`)
- Secure memory functions (`secure_memory*.go`) -- safe copy-then-wipe implementation
- Custom error types (`errors.go`)
- CI/CD pipeline (`.github/workflows/ci.yml`)
- Makefile, `.golangci.yml`, `.markdownlintrc`
- Comprehensive documentation (10 documents)
- CodeRabbit unit tests (`docs_test.go`, `TESTING_GUIDE.md`, `TEST_RESULTS_SUMMARY.md`)

> **Note:** PR #5 (`claude/update-crypto-refactor-memory`) is superseded.
> It proposed an alternative secure memory implementation with issues
> (deprecated APIs, wrong Windows function name, vulnerable crypto version).
> Master's implementation is safer and should be used.

---

## Phase 1: Critical Foundation

**Priority:** P0 -- CRITICAL
**Blocks:** All subsequent work
**Stages:** 3 (GORM migration, driver upgrades, transitive deps)
**Task register:** [TASKS.md](TASKS.md) (authoritative task IDs and dependencies)

### Stage 1.1: GORM v2 Migration

**Estimated effort:** 6-8 hours
**Risk:** MEDIUM (breaking API changes)
**Branch:** `stage-1-gorm-v2`

| Task | File(s) | Action | Est. |
|------|---------|--------|------|
| **TASK-001** Update go.mod | `go.mod` | Remove `jinzhu/gorm`, add `gorm.io/gorm` | 0.5h |
| **TASK-002** Update auth_store.go imports | `auth_store.go` | `jinzhu/gorm` -> `gorm.io/gorm`, add `"errors"` | 0.25h |
| **TASK-003** Update FindIdentity error handling | `auth_store.go` | `gorm.IsRecordNotFoundError()` -> `errors.Is(err, gorm.ErrRecordNotFound)` | 0.25h |
| **TASK-004** Update auth_store_test.go | `auth_store_test.go` | Import path + connection API | 0.5h |
| **TASK-005** Create GORM v2 integration test | `auth_store_integration_test.go` (new) | CRUD round-trip on PostgreSQL | 2h |
| **TASK-006** Run full test suite | -- | `go test ./... -v -race` | 0.5h |
| **TASK-007** Security scan | -- | `gosec ./...` + `govulncheck ./...` | 0.5h |
| **TASK-008** Code review and PR | -- | Branch PR, 2+ reviewers | 2h |
| **TASK-009** Merge and tag | -- | Tag `v0.2.0-stage1` | 0.5h |

**Key breaking changes:**

```go
// Import path
"github.com/jinzhu/gorm"  ->  "gorm.io/gorm"

// Error check
gorm.IsRecordNotFoundError(err)  ->  errors.Is(err, gorm.ErrRecordNotFound)

// Connection
gorm.Open("postgres", dsn)  ->  gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

**Exit criteria:**
All unit + integration tests pass, security scan clean, PR approved.

### Stage 1.2: Database Driver Upgrades

**Estimated effort:** 2-3 hours
**Risk:** LOW (no API changes)
**Depends on:** Stage 1.1
**Branch:** `stage-2-db-drivers`

| Task | Package | Action | Est. |
|------|---------|--------|------|
| **TASK-010** Upgrade lib/pq | `github.com/lib/pq` | v1.1.1 -> latest (security and stability fixes) | 0.25h |
| **TASK-011** Upgrade go-sqlite3 | `github.com/mattn/go-sqlite3` | -> latest | 0.25h |
| **TASK-012** Upgrade go-sql-driver/mysql | `github.com/go-sql-driver/mysql` | -> latest (TLS 1.3, MySQL 8.4) | 0.25h |
| **TASK-013** Upgrade go-mssqldb | `github.com/denisenkom/go-mssqldb` | -> latest (SQL Server 2022) | 0.25h |
| **TASK-014** Update GORM driver packages | `gorm.io/driver/*` | -> latest | 0.25h |
| **TASK-015** Integration tests (all DBs) | -- | PostgreSQL, MySQL, SQLite via Docker | 1.5h |
| **TASK-016** Security scan and merge | -- | Tag `v0.2.0-stage2` | 0.25h |

> **NOTE:** Version numbers from Nov 2025 plans are stale. Verify latest
> versions at pkg.go.dev before executing each `go get`.

**Exit criteria:**
Integration tests pass on PostgreSQL, MySQL, SQLite. No new vulnerabilities.

### Stage 1.3: Transitive Dependencies

**Estimated effort:** 1-2 hours
**Risk:** LOW
**Depends on:** Stage 1.2
**Branch:** `stage-3-transitive`

| Task | Action | Est. |
|------|--------|------|
| **TASK-017** Update golang.org/x packages | `go get golang.org/x/crypto@v0.45.0` or later (v0.43.0 has CVEs), then `go get -u golang.org/x/...` | 0.5h |
| **TASK-018** Update all indirect deps | `go get -u ./... && go mod tidy`, review diff | 0.5h |
| **TASK-019** Full test suite with race detection | `go test ./... -v -race` | 0.5h |
| **TASK-020** Merge and tag | Tag `v0.2.0-stage3` | 0.5h |

**Exit criteria:**
All deps at latest compatible versions, race detector clean.

---

## Phase 2: Security and Testing

**Priority:** P1 -- HIGH
**Depends on:** Phase 1 complete
**Stages:** 2 (security integration, comprehensive tests)

### Stage 2.1: Security Integration

**Estimated effort:** 5 hours
**Branch:** `stage-4-security`

| Task | File(s) | Action | Est. |
|------|---------|--------|------|
| **TASK-021** Integrate ValidateIdk into FindIdentity | `auth_store.go` | Add validation before DB query | 0.5h |
| **TASK-022** Integrate ValidateIdk into SaveIdentity | `auth_store.go` | Validate Idk + nil check | 0.5h |
| **TASK-023** Integrate ValidateIdk into DeleteIdentity | `auth_store.go` | Add validation | 0.25h |
| **TASK-024** Implement FindIdentitySecure helper | `auth_store.go` | Returns SecureIdentityWrapper | 1h |
| **TASK-025** Security test suite | `auth_store_security_test.go` (new) | SQL injection, DoS, char injection | 2h |
| **TASK-026** Security scan | -- | gosec + govulncheck | 0.25h |
| **TASK-027** Merge and tag | -- | Tag `v0.2.0-stage4` | 0.5h |

**Exit criteria:**
All AuthStore methods validate inputs. 10+ security tests pass.

### Stage 2.2: Comprehensive Test Suite

**Estimated effort:** 12 hours
**Branch:** `stage-5-testing`

| Task | File(s) | Action | Est. |
|------|---------|--------|------|
| **TASK-028** Unit tests (20+ cases) | `auth_store_comprehensive_test.go` (new) | TC-001 to TC-020 per API_TESTS_SPEC | 4h |
| **TASK-029** Integration tests (15 cases) | `auth_store_integration_test.go` | Multi-DB + concurrency | 3h |
| **TASK-030** Benchmarks (6 cases) | `auth_store_bench_test.go` (new) | PERF-001 to PERF-006 | 2h |
| **TASK-031** Test data helpers | `test_helpers.go` (new) | TestIdentityBuilder pattern | 1h |
| **TASK-032** Measure and verify coverage | -- | Target: 70%+ overall, 80%+ auth_store.go | 1h |
| **TASK-033** CI test workflow | `.github/workflows/ci.yml` | Ensure coverage gate enforced | 0.5h |
| **TASK-034** Merge and tag | -- | Tag `v0.3.0-rc1` | 0.5h |

**Exit criteria:**
70+ tests, 70%+ coverage, benchmarks baseline established.

---

## Phase 3: Production Readiness and Release

**Priority:** P1 -- HIGH
**Depends on:** Phase 2 complete
**Stages:** 2 (Stage 3.1: Production Hardening, Stage 3.2: Release v1.0.0)

### Stage 3.1: Production Hardening

| Task | Action | Est. |
|------|--------|------|
| **TASK-035** Context support (DECISION REQUIRED -- DP-003) | Add `context.Context` to methods or defer | 2h |
| **TASK-036** Production deployment documentation | `PRODUCTION.md`: DB config, pools, TLS | 1h |
| **TASK-037** Migration guide | `UPGRADE_FROM_V0.md`: breaking changes, steps | 0.5h |
| **TASK-038** Final security audit | gosec + govulncheck + manual review | 0.5h |

### Stage 3.2: Release v1.0.0

| Task | Action | Est. |
|------|--------|------|
| **TASK-039** Update README.md | Installation, quick start, API examples | 0.5h |
| **TASK-040** Create CHANGELOG.md | All changes from v0.x to v1.0.0 | 0.25h |
| **TASK-041** Tag v1.0.0 | `git tag -a v1.0.0` | 0.25h |
| **TASK-042** GitHub Release | Release page with changelog | 0.5h |
| **TASK-043** Revert module path to sqrldev | go.mod module path change | 0.25h |
| **TASK-044** Submit to pkg.go.dev | `GOPROXY=proxy.golang.org go list -m ...@v1.0.0` | 0.25h |

**Exit criteria:**
v1.0.0 tagged, GitHub release created, available on pkg.go.dev.

---

## Decision Points (2 open, 2 resolved)

### DP-001: MemGuard vs Custom Secure Memory

**Context:** Custom `WipeBytes`/`ClearIdentity` already implemented. MemGuard
(`github.com/awnumar/memguard`) provides kernel-backed locked pages.

**Options:**
- **A.** Keep custom implementation (simpler, no new dependency)
- **B.** Integrate MemGuard (stronger guarantees, adds dependency)

**Recommendation:** A for v1.0.0, evaluate B for v1.1.0.
**Blocks:** Nothing critical. Can be decided at any time.

### DP-002: Database Driver Selection

**Context:** Should the library ship with all 4 database drivers or
PostgreSQL-only?

**Options:**
- **A.** Multi-database (PostgreSQL, MySQL, SQLite, SQL Server)
- **B.** PostgreSQL-only (simplest, smallest dependency footprint)

**Recommendation:** A (multi-database) -- matches upstream GORM philosophy.
**Blocks:** Stage 1.2 scope (TASK-010 to TASK-014).

### DP-003: Context API Design (RESOLVED)

**Context:** Modern Go expects `context.Context` on I/O methods.

**Decision:** Option A (add context to all methods), implemented as
`*WithContext()` variants alongside originals for backward compatibility.
Upstream interface update (Option 2A) to be coordinated separately.

**Status:** RESOLVED -- Implemented 2026-02-07.
**Blocks:** None (TASK-035 complete).

### DP-004: Goauthentik Integration (RESOLVED)

**Status:** RESOLVED -- No action required.

**Context:** Question raised whether this library should integrate with
goauthentik.io or other IAM systems.

**Resolution:** SQRL is an independent authentication protocol, separate from
OAuth/OIDC/SAML. This library is SQRL-specific and does not need IAM
integration. Applications can use both SQRL and goauthentik as separate auth
methods if desired.

---

## Risk Register

| ID | Risk | Probability | Impact | Mitigation |
|----|------|-------------|--------|------------|
| R-001 | GORM v2 migration breaks functionality | MEDIUM | CRITICAL | Integration tests before merge; rollback to v0.1.0 tag |
| R-002 | Test coverage below 70% target | MEDIUM | MEDIUM | 12h allocated; CI gate enforces minimum |
| R-003 | Breaking changes impact downstream users | LOW | MEDIUM | Migration guide; maintain v0.x branch for fixes |
| R-004 | Security vulnerability discovered | LOW | CRITICAL | Multiple scan tools; hotfix process < 24h |
| R-005 | Upstream ssp incompatible with GORM v2 | MEDIUM | HIGH | Test early in Stage 1.1; fork if needed |

### Rollback Procedures

```bash
# Roll back to any previous stage tag
git checkout v0.2.0-stageN
git checkout -b rollback-to-stageN

# Roll back to pre-upgrade baseline
git checkout v0.1.0-pre-upgrade

# Database rollback (PostgreSQL)
psql -U postgres sqrl_db < backup_YYYYMMDD_HHMMSS.sql
```

---

## Identified Gaps

1. **Stale version references:** Plan originally specified `gorm.io/gorm
   v1.31.1`, `lib/pq v1.10.9`, etc. from Nov 2025. Verify latest versions at
   pkg.go.dev before executing each upgrade.

2. **Upstream ssp compatibility:** `server-go-ssp` uses GORM v1 internally.
   Confirm that importing both v1 (via ssp) and v2 (via gormauthstore) does
   not cause module conflicts. Test early in Stage 1.1.

3. **CI/CD gaps not in task list:**
   - Dependabot configuration for automatic dependency updates
   - Branch protection rules on main
   - Codecov integration for coverage reporting
   These should be addressed as part of TASK-033 or as follow-up.

4. **Go version divergence:** Plan was written targeting Go 1.23+. The
   `go.mod` now specifies Go 1.25 (upgraded from 1.24 on 2026-02-08).
   Dependency compatibility verified against 1.25.

5. **Missing pre-upgrade tag:** The plan assumes a `v0.1.0-pre-upgrade` tag
   exists as rollback baseline. This should be created before Stage 1.1
   begins.

---

## Archived Documents

The following documents were merged into this plan and moved to
`docs/archive/`:

| Document | Reason |
|----------|--------|
| `UNIFIED_TODO.md` | Superseded by this plan |
| `STAGED_UPGRADE_PLAN.md` | Step-by-step detail merged into stages above |
| `SECURITY_REVIEW_AND_UPGRADE_PLAN.md` | Findings captured; implementation sections done |
| `TODO.md` | Original checklist fully superseded |

---

**Document Control:**
- Consolidated from: UNIFIED_TODO.md, STAGED_UPGRADE_PLAN.md,
  SECURITY_REVIEW_AND_UPGRADE_PLAN.md, TODO.md
- Total tasks: 44 (39 done, 4 pending, 1 deferred)
- Authoritative task status: [TASKS.md](TASKS.md)
- Decision points: 2 open (DP-001, DP-002), 2 resolved (DP-003, DP-004)
