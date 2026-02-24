# CLAUDE.md - Project Instructions for Claude Code

## Working Preferences

- **Think first, read second, act third.** Always think through the problem,
  read the codebase for relevant files, and understand context before making
  changes. Never speculate about code you have not opened. If a specific file
  is referenced, read it before answering.
- **Check in before major changes.** Before making any major changes, present
  the plan and wait for confirmation before generating code.
- **Keep it simple.** Make every task and code change as simple as possible.
  Comply with good coding standards appropriate for the file type, but avoid
  massive or complex changes. Every change should impact as little code as
  possible. Simplicity, readability, and security are the priorities.
- **Security and defensive coding.** Code with security and defensive coding
  practices. All generated code must meet appropriate standards for the
  language, including security checks from NIST, OWASP, to at least FIPS 140-2
  standard.
- **High-level explanations.** For each step, provide a high-level explanation
  of what changes were made rather than exhaustive detail.
- **Grounded answers only.** Never make claims about code before investigating.
  Give grounded, hallucination-free answers. If unable to determine something,
  say so.
- **Verify "latest" claims.** Verify the current date against any references
  to "latest" and confirm the assumption from a trusted and verified source.
  Ask if unable to determine trust of source.
- **Diagrams in Mermaid.** Any diagrams created should be in Mermaid markdown
  format.
- **Maintain architecture documentation.** Keep `docs/ARCHITECTURE.md` up to
  date so it describes how the application works inside and out.
- **Use available skills.** Consider relevant skills (e.g., PDF skill at
  `/mnt/skills/public/pdf/SKILL.md` for PDF extraction) before proceeding.
  Create a plan first and ask for confirmation before generating code.

## Project Overview

This is a Go library (`gormauthstore`) that implements the `ssp.AuthStore`
interface for persisting SQRL authentication identities via the GORM ORM.

- **Module:** `github.com/dxcSithLord/server-go-ssp-gormauthstore`
  (Note: temporarily set to `dxcSithLord` for development; must revert to
  `github.com/sqrldev/server-go-ssp-gormauthstore` before upstream merge --
  see TASK-043.)
- **Upstream:** `github.com/sqrldev/server-go-ssp-gormauthstore`
- **Go version:** 1.25+ (go.mod: `go 1.25.0`, toolchain `go1.25.7`)
- **Build:** `make ci` (deps + fmt-check + lint + security + test + build)
- **Current release:** `v1.0.0` (tagged 2026-02-08)

## Build and Test Commands

```bash
make ci              # Full CI: deps, fmt-check, lint, security, test, build
make test            # Run tests with race detection and coverage
make test-coverage   # Tests + HTML coverage report
make test-short      # Tests without integration (skip -short tests)
make lint            # golangci-lint
make security        # gosec + govulncheck (informational)
make security-ci     # gosec + govulncheck (strict, SARIF output)
make build           # go build ./...
make build-all       # Cross-compile for linux/darwin/windows (amd64, arm64)
make fmt             # gofmt -s -w + goimports
make fmt-check       # Verify formatting (used in CI)
make bench           # go test -bench=. -benchmem ./...
make pre-commit      # fmt + lint + test-short (before committing)
make tools           # Install dev tools (golangci-lint, gosec, govulncheck, goimports)
make deps            # go mod download + tidy + verify
make clean           # Remove coverage.out, coverage.html, gosec-results.sarif
```

## Project Structure

```
.
├── auth_store.go                       # Core AuthStore: CRUD + FindIdentitySecure + *WithContext() variants
├── errors.go                           # Sentinel errors (ErrEmptyIdentityKey, ErrNilIdentity, etc.)
├── secure_memory.go                    # WipeBytes -- Unix implementation (build tag: !windows)
├── secure_memory_common.go             # WipeString, ClearIdentity, SecureIdentityWrapper, ValidateIdk, ScrambleBytes
├── secure_memory_windows.go            # WipeBytes -- Windows implementation
├── auth_store_test.go                  # Basic CRUD test (TestSave)
├── auth_store_comprehensive_test.go    # 27 unit tests (TC-001 to TC-027)
├── auth_store_context_test.go          # 13 context support tests (CTX-001 to CTX-013)
├── auth_store_security_test.go         # 14 security tests (SEC-001 to SEC-014)
├── auth_store_integration_test.go      # 10 integration tests (build-tag: integration)
├── auth_store_bench_test.go            # 6 benchmarks (PERF-001 to PERF-006)
├── secure_memory_test.go               # 18 secure memory/validation tests + 4 benchmarks
├── test_helpers_test.go                # testIdentityBuilder, openTestDB, newTestStore, seedIdentity
├── docs_test.go                        # 17 documentation integrity tests
├── docs/                               # All planning and design documentation
│   ├── PROJECT_PLAN.md                 # Consolidated plan (54 tasks, 3 phases, 7 stages)
│   ├── TASKS.md                        # Authoritative task register (source of truth)
│   ├── REQUIREMENTS.md                 # Functional and non-functional requirements
│   ├── ARCHITECTURE.md                 # TOGAF-aligned architecture views (Mermaid)
│   ├── API_SPECIFICATION.md            # Go interface specification with examples
│   ├── API_TESTS_SPEC.md               # 70+ test case specifications
│   ├── DEPENDENCIES.md                 # Dependency management and local build setup
│   ├── DOCUMENTATION_TESTS.md          # Documentation test descriptions
│   ├── PRODUCTION.md                   # Production deployment guide
│   ├── UPGRADE_FROM_V0.md              # Migration guide (v0.x → v1.0.0)
│   ├── Notice_of_Decisions.md          # Decision log and SQRL protocol compliance analysis
│   └── archive/                        # Superseded planning documents
├── .github/workflows/ci.yml            # GitHub Actions CI pipeline
├── .golangci.yml                       # golangci-lint v2 configuration
├── .markdownlintrc                     # Markdown lint configuration
├── .whitesource                        # WhiteSource/Mend dependency scanning configuration
├── Makefile                            # Development automation
├── go.mod / go.sum                     # Go module definition
├── CHANGELOG.md                        # Changelog (Keep a Changelog format)
├── TESTING_GUIDE.md                    # Test suite overview, patterns, and running instructions
├── TEST_RESULTS_SUMMARY.md             # Test results and coverage summary
├── LICENSE                             # MIT License
└── README.md                           # Project overview, status badges, quick start
```

## Test Suite Summary

**100 tests total, 100% coverage across all source files.**

| File | Tests | Description |
|------|-------|-------------|
| `auth_store_test.go` | 1 | Basic CRUD round-trip |
| `auth_store_comprehensive_test.go` | 27 | TC-001 to TC-027: CRUD edge cases, validation |
| `auth_store_context_test.go` | 13 | CTX-001 to CTX-013: `*WithContext()` methods |
| `auth_store_security_test.go` | 14 | SEC-001 to SEC-014: injection, DoS, unicode |
| `auth_store_integration_test.go` | 10 | Build-tag `integration`; requires `TEST_DATABASE_URL` |
| `auth_store_bench_test.go` | 6 benchmarks | PERF-001 to PERF-006 |
| `secure_memory_test.go` | 18 + 4 benchmarks | Memory clearing and validation |
| `docs_test.go` | 17 | Documentation integrity (files, links, Mermaid, no secrets) |

**Running integration tests:**

```bash
TEST_DATABASE_URL="host=localhost user=postgres dbname=sqrl sslmode=require" \
  go test -tags=integration -v ./...
```

## Progress Tracking - IMPORTANT

**After completing any task, you MUST update the authoritative task register
at `docs/TASKS.md`.** This is the single source of truth for task status.

Specifically:

1. **`docs/TASKS.md`** (source of truth):
   - Mark the task as `done` with the completion date
   - Update the Summary table counts at the bottom

2. **`README.md`** badge and summary line:
   - Update the tasks badge count (e.g., `50/54`)
   - Update the one-line summary under "Overall Progress"

3. **`README.md` TL;DR table** (if applicable):
   - Update "Status" and "Detail" columns for any area that changed
   - Update "Test Coverage" if new tests were added

Do **not** duplicate per-phase breakdowns in README.md. The README links
to `docs/TASKS.md` and `docs/PROJECT_PLAN.md` for detailed status.

## Current State

For authoritative task status, read [`docs/TASKS.md`](docs/TASKS.md).
For the full plan (phases, stages, decisions), read
[`docs/PROJECT_PLAN.md`](docs/PROJECT_PLAN.md).

**Overall: 50 of 54 tasks complete (93%). v1.0.0 tagged 2026-02-08.**

| Status | Count | Tasks |
|--------|-------|-------|
| Done | 50 | All Phase 1, Phase 2, Phase 3.1, most of Phase 3.2, all doc/maintenance |
| Pending | 1 | TASK-042 (GitHub Release page) |
| On-hold | 2 | TASK-043, TASK-044 (pending server-go-ssp upstream updates) |
| Deferred | 1 | TASK-015 (PostgreSQL/MySQL integration tests -- requires Docker CI) |

### Key documents to consult

| When you need | Read |
|---------------|------|
| Task status (source of truth) | [`docs/TASKS.md`](docs/TASKS.md) |
| Plan, phases, decision points | [`docs/PROJECT_PLAN.md`](docs/PROJECT_PLAN.md) |
| Requirements | [`docs/REQUIREMENTS.md`](docs/REQUIREMENTS.md) |
| Architecture | [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) |
| Decision log | [`docs/Notice_of_Decisions.md`](docs/Notice_of_Decisions.md) |
| Test patterns and running | [`TESTING_GUIDE.md`](TESTING_GUIDE.md) |
| API details | [`docs/API_SPECIFICATION.md`](docs/API_SPECIFICATION.md) |

### What exists in the code

- `auth_store.go` -- Core `AuthStore` using GORM v2 with all CRUD methods,
  `FindIdentitySecure`, and `*WithContext()` variants for every method.
  Internal `identityRecord` model provides GORM v2 tags (upstream
  `SqrlIdentity` uses legacy GORM v1 `sql:""` tags).
- `errors.go` -- Sentinel errors: `ErrEmptyIdentityKey`, `ErrIdentityKeyTooLong`,
  `ErrInvalidIdentityKeyFormat`, `ErrNilIdentity`, `ErrNilDatabase`,
  `ErrWrappedIdentityDestroyed`.
- `secure_memory.go` / `secure_memory_common.go` / `secure_memory_windows.go` --
  Platform-aware secure memory clearing: `WipeBytes` (Unix/Windows),
  `WipeString`, `ClearIdentity`, `ScrambleBytes`, `SecureIdentityWrapper`,
  `ValidateIdk`.

## Code Conventions

- Use Go standard formatting (`gofmt -s`); import ordering via `goimports`
- Follow existing error patterns -- sentinel errors declared in `errors.go`
- Security-sensitive fields (`Suk`, `Vuk`) must be cleared via `ClearIdentity()`
  or `SecureIdentityWrapper.Destroy()` when no longer needed
- All public functions must validate inputs before database operations
  (`ValidateIdk()` called on every method that takes an `idk` parameter)
- Never log sensitive cryptographic material (`Idk`, `Suk`, `Vuk`, `Pidk`)
- Prefer `*WithContext()` variants in new code; keep original signatures for
  backward compatibility
- Line length limit: 140 characters (enforced by `lll` linter); comments exempt
- Locale: UK English (enforced by `misspell` linter)
- Linter config: `.golangci.yml` (golangci-lint v2 format); test files are
  exempt from `gosec`, `gocyclo`, `dupl`, `errcheck`, and related linters

## Decision Points

Open decision points documented in `docs/PROJECT_PLAN.md` and
`docs/Notice_of_Decisions.md`:

- **DP-001:** MemGuard vs custom implementation for secure memory -- **OPEN.**
  Current implementation uses custom `WipeBytes`/`ScrambleBytes`. True secure
  memory requires `[]byte` throughout (not `string`); see comments in
  `secure_memory_common.go:WipeString` for documented limitations.
- **DP-002:** Database driver selection (PostgreSQL-only vs multi-database) --
  **OPEN.** Current `go.mod` only includes `gorm.io/driver/sqlite` for tests.
  PostgreSQL/MySQL/SQLServer drivers not listed as direct dependencies.
- **DP-003:** Context API -- **RESOLVED.** Option A implemented: all methods
  have `*WithContext()` companions alongside original signatures for backward
  compatibility.
