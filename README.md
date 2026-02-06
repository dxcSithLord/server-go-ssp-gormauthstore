# sqrl-gormauthstore

SQRL `ssp.AuthStore` implementation using the GORM ORM.

## TL;DR - Project Status

| Area | Status | Detail |
|------|--------|--------|
| **SQRL Protocol Compliance** | **COMPLIANT** | All required storage fields (Idk, Suk, Vuk) plus optional enhancements |
| **GORM Version** | **CURRENT (v2 -- gorm.io/gorm v1.31.1)** | Migrated from deprecated jinzhu/gorm v1.9.16 |
| **Go Version** | 1.24 | Module initialised with Go 1.24.0 toolchain |
| **Test Coverage** | ~35% | Target: 70%+. Unit + integration tests pass with SQLite |
| **CI/CD Pipeline** | Configured | GitHub Actions workflow with lint, security scan, build matrix |
| **Security Hardening** | Partial | Secure memory clearing implemented; not yet integrated into AuthStore |
| **Documentation** | Comprehensive | 10 documents covering requirements, architecture, API, security, and upgrade plan |

### Overall Progress

| Phase | Description | Tasks | Completed | Status |
|-------|-------------|-------|-----------|--------|
| **Phase 1** | Critical Foundation (GORM v2, drivers, deps) | 20 | 16 | In progress |
| **Phase 2** | Security & Testing | 14 | 0 | Not started |
| **Phase 3** | Production Readiness & Release | 10 | 0 | Not started |
| **Docs & Infra** | Documentation, CI/CD, secure memory | -- | Done | Complete |
| **TOTAL** | 44 implementation tasks | 44 | 16 | **36%** |

> **Next milestone:** Phase 2 / Stage 2.1 -- Security Integration.
> See [docs/PROJECT_PLAN.md](docs/PROJECT_PLAN.md) for the full plan.

---

## Overview

This Go library provides database-backed persistence for
[SQRL](https://www.grc.com/sqrl/sqrl.htm) (Secure Quick Reliable Login)
authentication identities. It implements the `ssp.AuthStore` interface defined
by [server-go-ssp](https://github.com/sqrldev/server-go-ssp), allowing any
GORM-supported database (PostgreSQL, MySQL, SQLite, SQL Server) to store SQRL
identity records.

### Features

- **CRUD operations** -- `FindIdentity`, `SaveIdentity`, `DeleteIdentity`
- **Schema management** -- `AutoMigrate` for automatic table creation
- **Multi-database** -- PostgreSQL, MySQL, SQLite, SQL Server via GORM drivers
- **Secure memory** -- Platform-aware clearing of sensitive cryptographic keys
- **Input validation** -- `ValidateIdk()` with length and character-set checks
- **Custom errors** -- `ErrEmptyIdentityKey`, `ErrIdentityKeyTooLong`,
  `ErrInvalidIdentityKeyFormat`

### Quick Start

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    gormauthstore "github.com/sqrldev/server-go-ssp-gormauthstore"
)

dsn := "host=localhost user=postgres dbname=sqrl sslmode=require"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil { panic(err) }

store := gormauthstore.NewAuthStore(db)
store.AutoMigrate()
```

> **Note:** The code uses `gorm.io/gorm` v2. An internal `identityRecord`
> model provides GORM v2 tags for the upstream `SqrlIdentity` struct.

## Documentation

All project documentation lives in the [`docs/`](docs/) directory:

| Document | Purpose |
|----------|---------|
| [PROJECT_PLAN.md](docs/PROJECT_PLAN.md) | Consolidated project plan (44 tasks, 3 phases, 6 stages) |
| [REQUIREMENTS.md](docs/REQUIREMENTS.md) | Reverse-engineered functional and non-functional requirements |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | TOGAF-aligned architecture views |
| [API_SPECIFICATION.md](docs/API_SPECIFICATION.md) | OpenAPI-style Go interface specification |
| [API_TESTS_SPEC.md](docs/API_TESTS_SPEC.md) | 70+ test case specifications |
| [DEPENDENCIES.md](docs/DEPENDENCIES.md) | Dependency management and local build setup |
| [Notice_of_Decisions.md](docs/Notice_of_Decisions.md) | Decision log with SQRL protocol compliance analysis |
| [archive/](docs/archive/) | Superseded planning documents (TODO, UNIFIED_TODO, STAGED_UPGRADE_PLAN, SECURITY_REVIEW) |

## Development

```bash
# Install tools
make tools

# Run full CI locally
make ci

# Run tests with coverage
make test-coverage

# Run security checks
make security

# Format and lint
make fmt lint
```

## License

[MIT](LICENSE)
