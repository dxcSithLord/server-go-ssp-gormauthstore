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
├── auth_store.go                       # Core AuthStore (CRUD + FindIdentitySecure)
├── errors.go                           # Sentinel errors
├── secure_memory.go                    # WipeBytes (Unix)
├── secure_memory_common.go             # WipeString, ClearIdentity, SecureIdentityWrapper, ValidateIdk
├── secure_memory_windows.go            # WipeBytes (Windows)
├── auth_store_test.go                  # Basic CRUD test
├── auth_store_comprehensive_test.go    # 27 unit tests (TC-001 to TC-027)
├── auth_store_context_test.go          # 13 context support tests (CTX-001 to CTX-013)
├── auth_store_security_test.go         # 13 security tests (SEC-001 to SEC-013)
├── auth_store_integration_test.go      # 10 integration tests (build-tag: integration)
├── auth_store_bench_test.go            # 6 benchmarks (PERF-001 to PERF-006)
├── secure_memory_test.go               # Secure memory + validation tests + benchmarks
├── test_helpers_test.go                # testIdentityBuilder, newTestStore, seedIdentity
├── docs_test.go                        # Documentation integrity tests
├── docs/                               # All planning and design documentation
│   ├── PROJECT_PLAN.md                 # Consolidated plan (44 tasks, 3 phases, 7 stages)
│   ├── TASKS.md                        # Authoritative task register
│   ├── REQUIREMENTS.md
│   ├── ARCHITECTURE.md
│   ├── API_SPECIFICATION.md
│   ├── API_TESTS_SPEC.md
│   ├── DEPENDENCIES.md
│   ├── PRODUCTION.md                   # Production deployment guide
│   ├── UPGRADE_FROM_V0.md              # Migration guide (v0.x to v1.0.0)
│   ├── Notice_of_Decisions.md
│   └── archive/                        # Superseded planning documents
├── .github/workflows/ci.yml           # GitHub Actions CI pipeline
├── Makefile                            # Development automation
├── go.mod / go.sum                     # Go module definition
└── README.md                           # Project overview with progress summary
```

## Progress Tracking - IMPORTANT

**After completing any task, you MUST update the authoritative task register
at `docs/TASKS.md`.** This is the single source of truth for task status.

Specifically:

1. **`docs/TASKS.md`** (source of truth):
   - Mark the task as `done` with the completion date
   - Update the Summary table counts at the bottom

2. **`README.md`** badge and summary line:
   - Update the tasks badge count (e.g., `39/44`)
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

### Key documents to consult

| When you need | Read |
|---------------|------|
| Task status (source of truth) | [`docs/TASKS.md`](docs/TASKS.md) |
| Plan, phases, decision points | [`docs/PROJECT_PLAN.md`](docs/PROJECT_PLAN.md) |
| Requirements | [`docs/REQUIREMENTS.md`](docs/REQUIREMENTS.md) |
| Architecture | [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) |
| Decision log | [`docs/Notice_of_Decisions.md`](docs/Notice_of_Decisions.md) |

### What exists in the code

- `auth_store.go` - Core AuthStore using GORM v2 with `FindIdentitySecure`
  and `*WithContext()` variants for all methods
- `errors.go` - Sentinel errors (ErrEmptyIdentityKey, ErrNilIdentity, etc.)
- `secure_memory.go` / `secure_memory_common.go` / `secure_memory_windows.go` -
  Platform-aware secure memory clearing (WipeBytes, ClearIdentity, ScrambleBytes)

## Code Conventions

- Use Go standard formatting (`gofmt`)
- Follow existing error patterns (sentinel errors in `errors.go`)
- Security-sensitive fields (Suk, Vuk) must be cleared with `ClearIdentity()`
- All public functions should validate inputs before database operations
- No logging of sensitive cryptographic material (Idk, Suk, Vuk, Pidk)

## Decision Points

There are 2 open decision points documented in `docs/PROJECT_PLAN.md` that
require human input before proceeding:

- **DP-001:** MemGuard vs custom implementation for secure memory
- **DP-002:** Database driver selection (PostgreSQL-only vs multi-database)
- **DP-003:** RESOLVED -- Context API: Option A implemented as `*WithContext()`
  methods alongside originals for backward compatibility
