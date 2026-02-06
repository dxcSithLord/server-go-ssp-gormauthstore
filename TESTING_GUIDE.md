# Testing Guide

## Overview

This project includes comprehensive test suites covering both code and documentation.

## Test Categories

### 1. Code Tests

#### Secure Memory Tests
**File:** `secure_memory_test.go`

Tests for secure memory operations:
- `TestWipeBytes` - Memory wiping functionality
- `TestScrambleBytes` - Byte scrambling
- `TestWipeString` - String wiping
- `TestClearIdentity` - Identity clearing
- `TestSecureIdentityWrapper` - Secure wrapper functionality
- `TestValidateIdk` - Identity key validation

**Coverage:** ~90-95% of secure memory functions

#### AuthStore Tests
**File:** `auth_store_test.go`

Tests for authentication store operations:
- `TestSave` - Basic save functionality

**Note:** This test suite needs expansion as part of Phase 2 (see PROJECT_PLAN.md)

### 2. Documentation Tests

**File:** `docs_test.go`

Comprehensive documentation validation:
- **Existence Tests** - Verify all required documentation files exist
- **Link Integrity Tests** - Validate internal markdown links
- **Content Structure Tests** - Ensure required sections are present
- **Consistency Tests** - Check cross-document consistency
- **Security Tests** - Verify no hardcoded credentials
- **Diagram Tests** - Validate mermaid diagram syntax

**Coverage:** 14 documents, 105+ assertions

See [docs/DOCUMENTATION_TESTS.md](docs/DOCUMENTATION_TESTS.md) for detailed documentation.

## Running Tests

### Run All Tests
```bash
make test
# or
go test ./... -v
```

### Run Specific Test Suites

#### Run Only Documentation Tests
```bash
go test -v -run "^Test(Documentation|Archive|Markdown|README|CLAUDE|ProjectPlan|Architecture|API|Dependencies|Notice|Phase1|Requirements|NoHardcoded|Mermaid)" ./...
```

#### Run Only Secure Memory Tests
```bash
go test -v -run "^Test.*Memory" ./...
```

#### Run Only AuthStore Tests
```bash
go test -v -run "^TestSave" ./...
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
# Detect race conditions
go test -race ./...
```

### Run Short Tests (Skip Integration Tests)

```bash
go test -short ./...
```

## Test Output

### Successful Test Run
```
=== RUN   TestDocumentationExists
=== RUN   TestDocumentationExists/README.md
--- PASS: TestDocumentationExists/README.md (0.00s)
...
PASS
ok      github.com/dxcSithLord/server-go-ssp-gormauthstore    0.023s
```

### Failed Test Run
```
=== RUN   TestDocumentationExists
--- FAIL: TestDocumentationExists (0.00s)
    docs_test.go:25: Required documentation file does not exist: README.md
FAIL
```

## CI/CD Integration

Tests are automatically run in the CI/CD pipeline:

```yaml
# .github/workflows/ci.yml
- name: Run tests
  run: go test -v -race -coverprofile=coverage.out ./...

- name: Check coverage
  run: |
    coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo "Coverage: $coverage"
```

## Writing New Tests

### Adding Code Tests

1. Create or update `*_test.go` file
2. Follow Go testing conventions:
   ```go
   func TestFeatureName(t *testing.T) {
       // Arrange
       input := "test data"

       // Act
       result := YourFunction(input)

       // Assert
       if result != expected {
           t.Errorf("Expected %v, got %v", expected, result)
       }
   }
   ```

3. Use sub-tests for multiple scenarios:
   ```go
   func TestFeature(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected string
       }{
           {"case1", "input1", "output1"},
           {"case2", "input2", "output2"},
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               result := YourFunction(tt.input)
               if result != tt.expected {
                   t.Errorf("Expected %v, got %v", tt.expected, result)
               }
           })
       }
   }
   ```

### Adding Documentation Tests

1. Update `docs_test.go`
2. Add document to appropriate test function:
   ```go
   func TestDocumentationExists(t *testing.T) {
       requiredDocs := []string{
           "README.md",
           "CLAUDE.md",
           "docs/YOUR_NEW_DOC.md",  // Add here
       }
       // ...
   }
   ```

3. Add content validation if needed:
   ```go
   func TestYourNewDocument(t *testing.T) {
       content, err := os.ReadFile("docs/YOUR_NEW_DOC.md")
       if err != nil {
           t.Fatalf("Cannot read YOUR_NEW_DOC.md: %v", err)
       }

       // Validate required sections
       if !strings.Contains(string(content), "## Required Section") {
           t.Error("Missing required section")
       }
   }
   ```

## Coverage Goals

| Component | Current | Target | Status |
|-----------|---------|--------|--------|
| **secure_memory.go** | ~95% | 95% | ✅ Met |
| **secure_memory_common.go** | ~90% | 90% | ✅ Met |
| **secure_memory_windows.go** | ~95% | 95% | ✅ Met |
| **errors.go** | 100% | 100% | ✅ Met |
| **auth_store.go** | ~25% | 80% | ⚠️ Phase 2 |
| **Documentation** | 100% | 100% | ✅ Met |
| **Overall** | ~30% | 70% | ⚠️ Phase 2 |

## Test Best Practices

### DO
✅ Write clear, descriptive test names
✅ Use table-driven tests for multiple scenarios
✅ Test both success and failure cases
✅ Test edge cases and boundary conditions
✅ Use sub-tests for better organization
✅ Clean up resources (use `defer`)
✅ Keep tests independent and isolated
✅ Run tests with `-race` flag

### DON'T
❌ Test implementation details
❌ Write flaky tests
❌ Ignore test failures
❌ Skip tests without good reason
❌ Use sleep for timing (use proper synchronization)
❌ Depend on external services in unit tests
❌ Leave commented-out tests

## Debugging Tests

### Run Single Test
```bash
go test -v -run TestSpecificFunction
```

### Run with Verbose Output
```bash
go test -v ./...
```

### Run with CPU Profiling
```bash
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

### Run with Memory Profiling
```bash
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

## Benchmarking

Run performance benchmarks:

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkWipeBytes ./...

# With memory stats
go test -bench=. -benchmem ./...
```

## Test Utilities

### Test Helpers

The project includes test helper functions:

**Current GORM v1 syntax (pre-migration):**
```go
// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }
    return db
}
```

**GORM v2 syntax (post-migration -- Phase 1 complete):**
```go
// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open test database: %v", err)
    }
    return db
}
```

Use `t.Helper()` to mark functions as test helpers - this improves error reporting.

## Continuous Testing

### Watch Mode (Using External Tool)
```bash
# Install gotestsum
go install gotest.tools/gotestsum@latest

# Run in watch mode
gotestsum --watch
```

### Pre-commit Hook
```bash
# Add to .git/hooks/pre-commit
#!/bin/bash
make test
```

## Troubleshooting

### Tests Pass Locally but Fail in CI

1. Check Go version consistency
2. Verify environment variables
3. Check for race conditions (`-race` flag)
4. Verify file paths are absolute
5. Check for timezone dependencies

### Flaky Tests

1. Identify with `go test -count=100`
2. Check for race conditions
3. Look for timing dependencies
4. Verify proper test isolation

### Slow Tests

1. Profile with `-cpuprofile`
2. Use `-short` flag for quick tests
3. Consider parallel execution
4. Move integration tests to separate suite

## See Also

- [docs/DOCUMENTATION_TESTS.md](docs/DOCUMENTATION_TESTS.md) - Documentation test specification
- [TEST_RESULTS_SUMMARY.md](TEST_RESULTS_SUMMARY.md) - Latest test results
- [docs/API_TESTS_SPEC.md](docs/API_TESTS_SPEC.md) - API test specifications
- [docs/PROJECT_PLAN.md](docs/PROJECT_PLAN.md) - Phase 2 test expansion plan

---

**Last Updated:** 2026-02-05
**Maintainer:** Testing Team

**END OF TESTING GUIDE**