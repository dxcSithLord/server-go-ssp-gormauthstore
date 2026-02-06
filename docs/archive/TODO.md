# Security Review TODO Checklist

## Priority: CRITICAL (Week 1-2)

### Dependency Upgrades

- [ ] **Migrate from deprecated `github.com/jinzhu/gorm` to `gorm.io/gorm` v1.31.1**
  - Update import path in `auth_store.go` (line 4)
  - Change error handling from `gorm.IsRecordNotFoundError(err)` to `errors.Is(err, gorm.ErrRecordNotFound)` (line 28)
  - Update database driver import from `github.com/jinzhu/gorm/dialects/postgres` to `gorm.io/driver/postgres`
  - Update connection syntax in test files

- [ ] **Initialize Go Module**
  - Review generated `go.mod` file
  - Run `go mod tidy` to populate `go.sum`
  - Test compilation with new module structure

- [ ] **Update Go Runtime Support**
  - Ensure compatibility with Go 1.23+ (current supported versions: Go 1.23, Go 1.24)
  - Test builds on both supported Go versions

### Breaking Code Changes Required

```go
// auth_store.go - OLD CODE (lines 3-6):
import (
    "github.com/jinzhu/gorm"
    ssp "github.com/sqrldev/server-go-ssp"
)

// auth_store.go - NEW CODE:
import (
    "errors"

    ssp "github.com/sqrldev/server-go-ssp"
    "gorm.io/gorm"
)

// auth_store.go - OLD ERROR CHECK (line 28):
if gorm.IsRecordNotFoundError(err) {

// auth_store.go - NEW ERROR CHECK:
if errors.Is(err, gorm.ErrRecordNotFound) {
```

## Priority: HIGH (Week 3-4)

### Security Vulnerability Fixes

- [ ] **CWE-226: Sensitive Information Not Removed Before Reuse**
  - [x] Created `secure_memory.go` with platform-aware memory clearing (Unix/Linux)
  - [x] Created `secure_memory_windows.go` with Windows-specific implementation
  - [x] Implemented `ClearIdentity()` function
  - [x] Implemented `SecureIdentityWrapper` for RAII-style cleanup
  - [ ] Integrate secure clearing into `FindIdentity()` return path
  - [ ] Document proper usage patterns for callers

- [ ] **CWE-200: Exposure of Sensitive Information**
  - [ ] Update test database connection to use `sslmode=prefer` or `sslmode=require`
  - [ ] Implement environment variable configuration for database credentials
  - [ ] Add warning in documentation about plain-text storage implications

- [ ] **CWE-20: Improper Input Validation**
  - [x] Created `ValidateIdk()` function for identity key validation
  - [x] Created custom error types in `errors.go`
  - [ ] Integrate validation into `FindIdentity()`, `DeleteIdentity()`
  - [ ] Add validation for `SaveIdentity()` (check for nil, empty Idk)

### Test Coverage Enhancement

- [x] Created comprehensive test suite for secure memory functions (`secure_memory_test.go`)
- [ ] Achieve 70%+ code coverage (currently ~25%)
- [ ] Add integration tests for database operations
- [ ] Add concurrent access tests
- [ ] Add SQL injection prevention tests
- [ ] Add benchmark tests for performance baseline

## Priority: MEDIUM (Week 5-8)

### CI/CD Pipeline

- [x] Created GitHub Actions workflow (`.github/workflows/ci.yml`)
- [x] Created golangci-lint configuration (`.golangci.yml`)
- [x] Created Makefile for development tasks
- [ ] Set up Codecov for coverage reporting
- [ ] Configure branch protection rules on main branch
- [ ] Enable CodeQL security scanning
- [ ] Add Dependabot for automatic dependency updates

### Production Readiness

- [ ] Add `context.Context` support to all database operations
- [ ] Implement structured logging (avoid logging sensitive data)
- [ ] Add metrics/observability hooks
- [ ] Create performance benchmarks
- [ ] Document deployment requirements
- [ ] Add health check functionality

### Advanced Security Features

- [ ] Consider integrating `github.com/awnumar/memguard` v0.23.0 for protected memory
- [ ] Implement field-level encryption hooks for sensitive data at rest
- [ ] Add audit logging for all identity operations
- [ ] Create rate limiting interfaces for abuse prevention
- [ ] Document threat model and security assumptions

## Dependency Version Reference

### Current Versions (as of November 17, 2025)

| Package | Version | Status |
|---------|---------|--------|
| Go Runtime | 1.24.10 / 1.23.8 | Supported (1.24, 1.23) |
| gorm.io/gorm | v1.31.1 | Latest (Nov 2025) |
| gorm.io/driver/postgres | v1.5.9 | Latest |
| github.com/sqrldev/server-go-ssp | v0.0.0-20241212182118 | Latest (Dec 2024) |
| github.com/awnumar/memguard | v0.23.0 | Latest (Aug 2025) |
| golangci-lint | v1.61.0 | Latest |
| gosec | v2.x | Latest |

### Upgrade Path Summary

1. **github.com/jinzhu/gorm** (v1.x) → **gorm.io/gorm** (v1.31.1)
   - Major API changes required
   - Error handling changes from method-based to errors.Is()
   - Driver imports change completely

2. **Go Runtime** (unknown) → **Go 1.23+**
   - Required for latest GORM compatibility
   - End of life for Go 1.22: Feb 2025
   - End of life for Go 1.23: Expected Aug 2025 (with Go 1.25 release)

## Quick Start Commands

```bash
# Install development tools
make tools

# Run full CI pipeline locally
make ci

# Run tests with coverage
make test-coverage

# Run security checks
make security

# Format and lint code
make fmt lint

# Build for all platforms
make build-all
```

## Files Created in This Review

1. `SECURITY_REVIEW_AND_UPGRADE_PLAN.md` - Comprehensive security analysis
2. `TODO.md` - This checklist
3. `go.mod` - Go module definition with target dependencies
4. `.github/workflows/ci.yml` - CI/CD pipeline configuration
5. `.golangci.yml` - Linting configuration
6. `Makefile` - Development automation
7. `secure_memory.go` - Unix/Linux secure memory clearing
8. `secure_memory_windows.go` - Windows secure memory clearing
9. `secure_memory_test.go` - Tests for secure memory functions
10. `errors.go` - Custom error definitions

## Estimated Effort

- **Phase 1 (Critical):** 26 hours
- **Phase 2 (High):** 40 hours
- **Phase 3 (Medium):** 36 hours
- **Phase 4 (Advanced):** 44 hours
- **Total:** ~146 hours (4 weeks full-time)

## Key Contacts

- **Repository:** github.com/sqrldev/server-go-ssp-gormauthstore
- **Original Author:** Scott White
- **SQRL Protocol:** https://www.grc.com/sqrl/sqrl.htm
- **GORM Documentation:** https://gorm.io/docs/

## Notes

- All sensitive data (Idk, Suk, Vuk, Pidk) contains cryptographic material
- Current code properly uses parameterized queries (SQL injection protected)
- No existing go.mod means this is a legacy GOPATH-based project
- Test coverage is minimal - only single integration test exists
- No CI/CD infrastructure currently in place
- License: MIT (changed from Apache 2.0 in Dec 2019)
