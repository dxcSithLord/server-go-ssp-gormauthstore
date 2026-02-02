# Dependency Management and Local Build Setup

## Overview

This document tracks all git repository dependencies, their versions, and provides instructions for building from local versioned releases.

---

## Repository References - Development vs Upstream

### Module Path Configuration

**Current Module Path:** `github.com/dxcSithLord/server-go-ssp-gormauthstore`

This module's path is set for dxcSithLord development fork. All import paths for dependencies remain as `github.com/sqrldev` since upstream modules declare themselves with those paths.

**To revert to upstream sqrldev:**
- Change `go.mod` line 3: `module github.com/sqrldev/server-go-ssp-gormauthstore`
- Remove the comment on lines 1-2

### Local Development with Replace Directives

To build against local copies of dependencies (e.g., modified server-go-ssp), add replace directives to `go.mod`:

```go
// Uncomment and adjust paths as needed for local builds
// Replace directives redirect import paths to local directories
replace github.com/sqrldev/server-go-ssp => ../server-go-ssp
```

This approach is used because:
1. Upstream modules declare themselves as `github.com/sqrldev/...` in their go.mod
2. Import paths must match the module's declared path (Go module requirement)
3. Replace directives allow local development without changing import paths

**IMPORTANT:** Remove replace directives before pushing to upstream repositories.

---

## Complete Dependency List

### Direct Dependencies

| Package | Current Version | Latest Available | Notes |
|---------|----------------|------------------|-------|
| **github.com/jinzhu/gorm** | v1.9.16 | **DEPRECATED** → gorm.io/gorm v1.31.1 | Major upgrade required (breaking API changes) |
| **github.com/dxcSithLord/server-go-ssp** | v0.0.0-20241212182118-c8230b16b87d | Latest pseudo-version | SQRL SSP protocol implementation |

### Indirect Dependencies (from go.sum)

| Package | Current Version | Latest Available (Nov 2025) | Update Priority |
|---------|----------------|------------------------------|----------------|
| github.com/davecgh/go-spew | v1.1.1 | v1.1.2 (if available) | LOW |
| github.com/jinzhu/inflection | v1.0.0 | v1.0.0 | Current |
| github.com/jinzhu/now | v1.1.5 | v1.1.5 | Current |
| github.com/lib/pq | v1.1.1 | **v1.10.9** | HIGH (security fixes) |
| github.com/mattn/go-sqlite3 | v1.14.22 | v1.14.24 | MEDIUM |
| github.com/skip2/go-qrcode | v0.0.0-20200617195104 | v0.0.0-20200617195104 | Current (no newer releases) |
| golang.org/x/crypto | **v0.43.0** | v0.43.0 | Current (just updated) |

### Transitive Dependencies (from server-go-ssp)

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/denisenkom/go-mssqldb | v0.0.0-20191124224453 | MSSQL driver (GORM) |
| github.com/erikstmartin/go-testdb | v0.0.0-20160219214506 | Test database |
| github.com/go-sql-driver/mysql | v1.5.0 | MySQL driver |
| github.com/golang-sql/civil | v0.0.0-20190719163853 | SQL civil types |
| github.com/PuerkitoBio/goquery | v1.5.1 | HTML parsing |
| github.com/andybalholm/cascadia | v1.1.0 | CSS selector |

---

## Critical Upgrade Paths

### 1. GORM v1 → v2 (CRITICAL)

**Current:** `github.com/jinzhu/gorm v1.9.16`
**Target:** `gorm.io/gorm v1.31.1`

This is a **breaking change** requiring code modifications:

```go
// Import change
// OLD: import "github.com/jinzhu/gorm"
// NEW: import "gorm.io/gorm"

// Error handling change
// OLD: gorm.IsRecordNotFoundError(err)
// NEW: errors.Is(err, gorm.ErrRecordNotFound)

// Driver import change
// OLD: _ "github.com/jinzhu/gorm/dialects/postgres"
// NEW: import "gorm.io/driver/postgres"

// Connection change
// OLD: gorm.Open("postgres", "connection_string")
// NEW: gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

### 2. PostgreSQL Driver (HIGH)

**Current:** `github.com/lib/pq v1.1.1`
**Target:** `github.com/lib/pq v1.10.9`

Contains security fixes and performance improvements.

### 3. SQLite Driver (MEDIUM)

**Current:** `github.com/mattn/go-sqlite3 v1.14.22`
**Target:** `github.com/mattn/go-sqlite3 v1.14.24`

---

## Local Development Setup

### Option 1: Using Go Module Replace Directives

For local development with local copies of dependencies:

```go
// Add to go.mod
replace (
    github.com/dxcSithLord/server-go-ssp => ../server-go-ssp
    github.com/dxcSithLord/server-go-ssp-gormauthstore => ../server-go-ssp-gormauthstore
)
```

### Option 2: Using Go Workspace (Go 1.18+)

Create `go.work` in parent directory:

```go
go 1.24

use (
    ./server-go-ssp
    ./server-go-ssp-gormauthstore
)
```

### Option 3: Vendoring Dependencies

```bash
# Vendor all dependencies
go mod vendor

# Build with vendor
go build -mod=vendor ./...
```

---

## Version Pinning Strategy

### Recommended go.mod for Production

```go
module github.com/sqrldev/server-go-ssp-gormauthstore

go 1.23

require (
    // Pin to specific versions for reproducible builds
    gorm.io/gorm v1.31.1
    gorm.io/driver/postgres v1.5.9
    github.com/sqrldev/server-go-ssp v0.0.0-20241212182118-c8230b16b87d
)

require (
    github.com/jackc/pgpassfile v1.0.0 // indirect
    github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
    github.com/jackc/pgx/v5 v5.7.1 // indirect
    github.com/jackc/puddle/v2 v2.2.2 // indirect
    github.com/jinzhu/inflection v1.0.0 // indirect
    github.com/jinzhu/now v1.1.5 // indirect
    golang.org/x/crypto v0.43.0 // indirect
    golang.org/x/sync v0.9.0 // indirect
    golang.org/x/text v0.20.0 // indirect
)
```

---

## Security Advisories

### Checked Dependencies (No Known Vulnerabilities as of Nov 2025)

- ✅ golang.org/x/crypto v0.43.0 - Updated with security fixes
- ✅ github.com/jinzhu/gorm v1.9.16 - No critical CVEs (but deprecated)

### Potential Concerns

- ⚠️ github.com/lib/pq v1.1.1 - Old version, upgrade recommended
- ⚠️ github.com/denisenkom/go-mssqldb - Old snapshot version
- ⚠️ golang.org/x/net - Ensure latest for security

---

## Release Version References

### Go Standard Library Extensions

| Package | Latest Version | Release Date | Notes |
|---------|---------------|--------------|-------|
| golang.org/x/crypto | v0.43.0 | Nov 2025 | Security-critical |
| golang.org/x/sync | v0.17.0 | 2025 | Concurrency primitives |
| golang.org/x/text | v0.30.0 | 2025 | Text processing |
| golang.org/x/net | v0.45.0 | 2025 | Network utilities |
| golang.org/x/sys | v0.37.0 | 2025 | System calls |

### Third-Party Libraries

| Package | Latest Tagged Release | GitHub Stars | Maintenance Status |
|---------|----------------------|--------------|-------------------|
| gorm.io/gorm | v1.31.1 | 37k+ | Active |
| github.com/lib/pq | v1.10.9 | 8k+ | Maintenance mode |
| github.com/mattn/go-sqlite3 | v1.14.24 | 7k+ | Active |
| github.com/skip2/go-qrcode | N/A (pseudo-versions only) | 1k+ | Minimal maintenance |

---

## Verification Commands

```bash
# List all module dependencies
go list -m all

# Check for outdated dependencies
go list -u -m all

# Verify module checksums
go mod verify

# Check for security vulnerabilities (requires internet)
govulncheck ./...

# Show why a dependency is needed
go mod why -m github.com/lib/pq

# Graph of dependencies
go mod graph
```

---

## Files Requiring Reversion Before Upstream Push

Execute this script to revert the module declaration to upstream sqrldev:

```bash
#!/bin/bash
# revert_to_upstream.sh

# Remove the development comment and revert module declaration
head -n 2 go.mod | grep -q 'dxcSithLord' && {
    # Remove first two comment lines and update module path
    tail -n +3 go.mod > go.mod.tmp
    sed -i 's/^module github\.com\/dxcSithLord\/server-go-ssp-gormauthstore$/module github.com\/sqrldev\/server-go-ssp-gormauthstore/' go.mod.tmp
    mv go.mod.tmp go.mod
}

# Ensure no replace directives remain (they should already be commented out)
sed -i '/^replace.*=>.*$/d' go.mod

# Regenerate go.sum
rm -f go.sum
go mod tidy

echo "Module path reverted to upstream sqrldev"
echo "Import paths for dependencies remain as github.com/sqrldev (no changes needed)"
```

**Note:** All import paths for dependencies (server-go-ssp, etc.) already use `github.com/sqrldev` paths and do not need to be changed. Only the module declaration in go.mod requires updating.

---

**Last Updated:** November 17, 2025
