# Testing Guide

## Overview

This project has 90 tests achieving 100% code coverage, plus 10 benchmarks,
across 8 test files.

## Test Suites

### 1. AuthStore Core Tests

**File:** `auth_store_test.go`

- `TestSave` - Basic CRUD functionality

### 2. Comprehensive Unit Tests

**File:** `auth_store_comprehensive_test.go` -- 27 tests (TC-001 to TC-027)

Covers FindIdentity, SaveIdentity, DeleteIdentity with all edge cases
including validation errors, not-found, upsert, concurrent access,
and field persistence.

### 3. Context Support Tests

**File:** `auth_store_context_test.go` -- 13 tests (CTX-001 to CTX-013)

Covers all `*WithContext()` methods with valid context, cancelled context,
validation, backward compatibility, and full CRUD round-trip.

### 4. Security Tests

**File:** `auth_store_security_test.go` -- 13 tests (SEC-001 to SEC-013)

Covers SQL injection prevention, DoS protection, unicode handling,
memory clearing, error message sanitization, and concurrent access safety.

### 5. Integration Tests

**File:** `auth_store_integration_test.go` -- 10 tests (build-tag gated)

Run with: `go test -tags=integration ./...`

Requires `TEST_DATABASE_URL` environment variable for PostgreSQL.

### 6. Secure Memory Tests

**File:** `secure_memory_test.go` -- 14 tests + 4 benchmarks

Covers WipeBytes, ScrambleBytes, WipeString, ClearIdentity,
SecureIdentityWrapper, ValidateIdk.

### 7. Documentation Tests

**File:** `docs_test.go` -- 14 test functions

Validates documentation file existence, link integrity, content structure,
cross-document consistency, security (no hardcoded credentials), and
mermaid diagram syntax.

### 8. Benchmarks

**File:** `auth_store_bench_test.go` -- 6 benchmarks (PERF-001 to PERF-006)

Covers FindIdentity, SaveIdentity, DeleteIdentity, FindIdentitySecure,
ValidateIdk, and ClearIdentity performance.

## Test Infrastructure

**File:** `test_helpers_test.go`

Provides reusable test helpers:
- `testIdentityBuilder` -- fluent builder pattern for test identities
- `openTestDB(t)` -- creates in-memory SQLite database
- `newTestStore(t)` -- creates AuthStore with auto-migration
- `seedIdentity(t, store, identity)` -- inserts test data

```go
// Example usage
store := newTestStore(t)
identity := newTestIdentity().
    withIdk("test-idk").
    withSuk("test-suk").
    withVuk("test-vuk").
    build()
seedIdentity(t, store, identity)
```

## Running Tests

### Run All Tests

```bash
make test
# or
go test -v -race ./...
```

### Run Specific Test Suite

```bash
# Comprehensive unit tests
go test -v -run "^TestTC" ./...

# Context support tests
go test -v -run "^Test.*WithContext\|^TestWithContext\|^TestOriginalMethods" ./...

# Security tests
go test -v -run "^TestSEC" ./...

# Documentation tests
go test -v -run "^Test(Documentation|Archive|Markdown|README|CLAUDE|ProjectPlan|Architecture|API|Dependencies|Notice|Phase1|Requirements|NoHardcoded|Mermaid)" ./...

# Secure memory tests
go test -v -run "^Test(Wipe|Scramble|Clear|SecureIdentity|Validate)" ./...
```

### Run Tests with Coverage

```bash
# All tests with coverage
make test-coverage

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Run Tests with Race Detection

```bash
go test -race ./...
```

### Run Benchmarks

```bash
# All benchmarks
go test -bench=. -benchmem ./...

# Specific benchmark
go test -bench=BenchmarkFindIdentity -benchmem ./...
```

### Run Short Tests (Skip Integration)

```bash
go test -short ./...
```

## Coverage Summary

| Component | Coverage |
|-----------|----------|
| `auth_store.go` | 100% |
| `secure_memory.go` | 100% |
| `secure_memory_common.go` | 100% |
| `errors.go` | 100% |
| **Overall** | **100%** |

**Coverage threshold:** 70% minimum enforced in CI.

## CI/CD Integration

Tests run automatically in GitHub Actions CI pipeline:

```yaml
- name: Run tests
  run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

- name: Check coverage threshold (70%)
  run: |
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo "Coverage: $coverage"
```

## Test Best Practices

### DO

- Write clear, descriptive test names
- Use table-driven tests for multiple scenarios
- Test both success and failure cases
- Test edge cases and boundary conditions
- Clean up resources (use `defer ClearIdentity(result)`)
- Run tests with `-race` flag
- Use `testIdentityBuilder` pattern for test data

### DON'T

- Test implementation details
- Write flaky tests
- Log sensitive fields (Idk, Suk, Vuk) in test output
- Use sleep for timing
- Depend on external services in unit tests

## See Also

- [docs/API_TESTS_SPEC.md](docs/API_TESTS_SPEC.md) - Test case specifications
- [TEST_RESULTS_SUMMARY.md](TEST_RESULTS_SUMMARY.md) - Test results summary
- [docs/DOCUMENTATION_TESTS.md](docs/DOCUMENTATION_TESTS.md) - Documentation test details

---

**Last Updated:** 2026-02-08
