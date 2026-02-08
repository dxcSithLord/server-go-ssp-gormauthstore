# Dependency Management and Local Build Setup

## Overview

This document tracks all module dependencies, their versions, and provides
instructions for building from local versioned releases.

**Last reviewed:** February 8, 2026

---

## Repository References - Development vs Upstream

### Module Path Configuration

**Current Module Path:** `github.com/dxcSithLord/server-go-ssp-gormauthstore`

This module's path is set for dxcSithLord development fork. The upstream
`ssp.AuthStore` interface is imported via
`github.com/dxcSithLord/server-go-ssp`.

**To revert to upstream sqrldev:**
- Change `go.mod` module line to `module github.com/sqrldev/server-go-ssp-gormauthstore`
- Update the ssp import to `github.com/sqrldev/server-go-ssp`
- Remove the development comment on lines 1-2

### Local Development with Replace Directives

To build against local copies of dependencies (e.g., modified server-go-ssp),
add replace directives to `go.mod`:

```go
// Uncomment and adjust paths as needed for local builds
replace github.com/dxcSithLord/server-go-ssp => ../server-go-ssp
```

**IMPORTANT:** Remove replace directives before pushing to upstream repositories.

---

## Complete Dependency List

### Direct Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/dxcSithLord/server-go-ssp` | v0.0.0-20260202110616 | SQRL SSP protocol (AuthStore interface) |
| `gorm.io/driver/sqlite` | v1.6.0 | SQLite database driver (test dependency) |
| `gorm.io/gorm` | v1.31.1 | GORM v2 ORM framework |

### Indirect Dependencies

| Package | Version | Source | Purpose |
|---------|---------|--------|---------|
| `github.com/fogleman/gg` | v1.3.0 | server-go-ssp | Graphics (QR code) |
| `github.com/golang/freetype` | v0.0.0 | server-go-ssp | Font rendering |
| `github.com/jinzhu/inflection` | v1.0.0 | gorm | Pluralization |
| `github.com/jinzhu/now` | v1.1.5 | gorm | Time helpers |
| `github.com/mattn/go-sqlite3` | v1.14.33 | sqlite driver | SQLite C bindings |
| `github.com/pkg/errors` | v0.9.1 | server-go-ssp | Error wrapping |
| `github.com/stretchr/testify` | v1.9.0 | server-go-ssp | Test assertions |
| `github.com/yeqown/go-qrcode/v2` | v2.2.5 | server-go-ssp | QR code generation |
| `github.com/yeqown/go-qrcode/writer/standard` | v1.3.0 | server-go-ssp | QR writer |
| `github.com/yeqown/reedsolomon` | v1.0.0 | server-go-ssp | Error correction |
| `golang.org/x/image` | v0.35.0 | server-go-ssp | Image processing |
| `golang.org/x/text` | v0.33.0 | gorm | Text processing |

### Production Database Drivers (Optional)

These drivers are not in go.mod but are required for production deployments
with databases other than SQLite:

| Package | Recommended Version | Database |
|---------|-------------------|----------|
| `gorm.io/driver/postgres` | v1.5.9 | PostgreSQL 12+ |
| `gorm.io/driver/mysql` | v1.5.7 | MySQL 8+ |
| `gorm.io/driver/sqlserver` | v1.5.3 | SQL Server 2019+ |

Install with:

```bash
go get gorm.io/driver/postgres   # or mysql, sqlserver
```

---

## Security Advisories

### Current Status (February 2026)

- gosec: 0 issues
- govulncheck: clean (network-blocked in CI; verified manually)
- All direct dependencies at latest stable versions

### Dependency Provenance

| Package | Maintenance Status | Notes |
|---------|-------------------|-------|
| `gorm.io/gorm` | Active (37k+ stars) | Major ORM framework |
| `gorm.io/driver/sqlite` | Active | Official GORM driver |
| `github.com/mattn/go-sqlite3` | Active (7k+ stars) | CGo SQLite bindings |
| `golang.org/x/image` | Active (Go team) | Standard library extension |
| `golang.org/x/text` | Active (Go team) | Standard library extension |

---

## Local Development Setup

### Option 1: Using Go Module Replace Directives

For local development with local copies of dependencies:

```go
// Add to go.mod
replace github.com/dxcSithLord/server-go-ssp => ../server-go-ssp
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
go mod vendor
go build -mod=vendor ./...
```

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
go mod why -m github.com/mattn/go-sqlite3

# Graph of dependencies
go mod graph
```

---

## Files Requiring Reversion Before Upstream Push

Execute this script to revert the module declaration to upstream sqrldev:

```bash
#!/bin/bash
# revert_to_upstream.sh

# Update module path
sed -i 's|github.com/dxcSithLord/server-go-ssp-gormauthstore|github.com/sqrldev/server-go-ssp-gormauthstore|' go.mod

# Update ssp import in source files
sed -i 's|github.com/dxcSithLord/server-go-ssp|github.com/sqrldev/server-go-ssp|' auth_store.go

# Remove replace directives
sed -i '/^replace.*=>.*$/d' go.mod

# Regenerate go.sum
rm -f go.sum
go mod tidy

echo "Module path reverted to upstream sqrldev"
```

---

**Document Control:**
- Version: 2.0
- Last Updated: February 8, 2026
- Next Review: Before v1.0.0 release
