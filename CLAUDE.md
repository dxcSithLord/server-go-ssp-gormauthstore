# CLAUDE.md - Project Instructions for Claude Code

## Project Overview

This is a Go library (`gormauthstore`) that implements the `ssp.AuthStore`
interface for persisting SQRL authentication identities via the GORM ORM.

- **Module:** `github.com/dxcSithLord/server-go-ssp-gormauthstore`
- **Upstream:** `github.com/sqrldev/server-go-ssp-gormauthstore`
- **Go version:** 1.24+
- **Build:** `make ci` (lint + security + test + build)

## Build and Test Commands

```bash
make ci            # Full CI: lint, security, test, build
make test          # Run tests with race detection and coverage
make test-coverage # Tests + HTML coverage report
make lint          # golangci-lint
make security      # gosec + govulncheck
make build         # go build ./...
make fmt           # gofmt
```

## Project Structure

```
.
├── *.go                          # Go source (auth_store, errors, secure_memory)
├── *_test.go                     # Go tests
├── docs/                         # All planning and design documentation
│   ├── UNIFIED_TODO.md           # Master plan (44 tasks, 3 phases)
│   ├── STAGED_UPGRADE_PLAN.md    # 6-stage upgrade path
│   ├── SECURITY_REVIEW_AND_UPGRADE_PLAN.md
│   ├── REQUIREMENTS.md
│   ├── ARCHITECTURE.md
│   ├── API_SPECIFICATION.md
│   ├── API_TESTS_SPEC.md
│   ├── DEPENDENCIES.md
│   ├── Notice_of_Decisions.md
│   └── TODO.md
├── .github/workflows/ci.yml     # GitHub Actions CI pipeline
├── Makefile                      # Development automation
├── go.mod / go.sum               # Go module definition
└── README.md                     # Project overview with progress summary
```

## Progress Tracking - IMPORTANT

**After completing any task from the project plan, you MUST update the
progress tables in `README.md`.**

Specifically, update the following sections:

1. **TL;DR - Project Status** table at the top of `README.md`:
   - Update the "Status" and "Detail" columns for any area that changed
   - Update "Test Coverage" if new tests were added
   - Update "GORM Version" when the migration is done
   - Update "Security Hardening" as features are integrated

2. **Overall Progress** table in `README.md`:
   - Increment the "Completed" column for the relevant phase
   - Update the "Status" column (e.g., "Not started" -> "In progress" -> "Complete")
   - Recalculate the TOTAL percentage

3. **docs/UNIFIED_TODO.md** task status:
   - Mark completed tasks with their completion date
   - Update the summary counts at the top of that document

### Example: After completing GORM v2 migration

In `README.md`, update:
- "GORM Version" row: change status from "DEPRECATED (v1.9.16)" to "CURRENT (v2.x.x)"
- Phase 1 "Completed" count: increment by the number of tasks finished
- Phase 1 "Status": change from "Not started" to "In progress"
- TOTAL: recalculate

## Current State (as of 2025-02-05)

- **Phase 1 (GORM v2 Migration):** Not started. This blocks all other work.
- **Phase 2 (Security & Testing):** Not started. Secure memory functions exist
  but are not yet integrated into AuthStore methods.
- **Phase 3 (Production Readiness):** Not started.
- **Infrastructure:** CI/CD pipeline, Makefile, golangci-lint, markdownlint,
  secure memory implementation, and comprehensive documentation are all done.

### What exists in the code

- `auth_store.go` - Core AuthStore using deprecated GORM v1
- `errors.go` - Custom validation error types (ErrEmptyIdentityKey, etc.)
- `secure_memory.go` / `secure_memory_common.go` / `secure_memory_windows.go` -
  Platform-aware secure memory clearing (WipeBytes, ClearIdentity, ScrambleBytes)
- `secure_memory_test.go` - Comprehensive tests for secure memory
- `auth_store_test.go` - Basic AuthStore test (needs expansion)

### What needs to be done next

1. **GORM v1 -> v2 migration** in `auth_store.go` (see docs/STAGED_UPGRADE_PLAN.md Stage 1)
2. Integrate `ValidateIdk()` into AuthStore methods
3. Expand test suite to 70%+ coverage
4. Tag v1.0.0 release

## Code Conventions

- Use Go standard formatting (`gofmt`)
- Follow existing error patterns (sentinel errors in `errors.go`)
- Security-sensitive fields (Suk, Vuk) must be cleared with `ClearIdentity()`
- All public functions should validate inputs before database operations
- No logging of sensitive cryptographic material (Idk, Suk, Vuk, Pidk)

## Decision Points

There are 4 open decision points documented in `docs/UNIFIED_TODO.md` that
require human input before proceeding:

- **DP-001:** MemGuard vs custom implementation for secure memory
- **DP-002:** Database driver selection (PostgreSQL-only vs multi-database)
- **DP-003:** Context API design (new methods vs modify existing signatures)
- **DP-004:** Release versioning (v0.x pre-release vs v1.0.0)
