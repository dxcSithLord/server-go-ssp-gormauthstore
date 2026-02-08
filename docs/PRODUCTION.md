# Production Deployment Guide

## gormauthstore Production Configuration

**Version:** 1.0
**Date:** February 7, 2026
**Audience:** DevOps engineers, backend developers deploying SQRL authentication

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Database Configuration](#database-configuration)
3. [Connection Pool Settings](#connection-pool-settings)
4. [TLS Configuration](#tls-configuration)
5. [Schema Migration](#schema-migration)
6. [Context and Timeout Configuration](#context-and-timeout-configuration)
7. [Security Checklist](#security-checklist)
8. [Monitoring and Health Checks](#monitoring-and-health-checks)
9. [Backup and Recovery](#backup-and-recovery)

---

```mermaid
flowchart LR
  App["App (SQRL Server)"] --> AuthStore["gormauthstore"]
  AuthStore --> DB["`Database
(PostgreSQL/MySQL/
SQL Server)`"]
  App --> Monitoring["Metrics/Alerts"]
```
## Prerequisites

- **Go:** 1.24 or later
- **Database:** One of the following:
  - PostgreSQL 12+ (recommended for production)
  - MySQL 8+
  - SQLite 3.35+ (development/testing only)
  - SQL Server 2019+
- **GORM:** v2 (`gorm.io/gorm` v1.31.1+)

---

## Database Configuration

### PostgreSQL (Recommended)

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    gormauthstore "github.com/sqrldev/server-go-ssp-gormauthstore"
)

dsn := "host=db.example.com user=sqrl_app password=REDACTED " +
    "dbname=sqrl_auth port=5432 sslmode=verify-full " +
    "sslrootcert=/etc/ssl/certs/ca.pem " +
    "sslcert=/etc/ssl/certs/client.pem " +
    "sslkey=/etc/ssl/private/client-key.pem"

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Warn),
})
if err != nil {
    log.Fatalf("failed to connect to database: %v", err)
}

store := gormauthstore.NewAuthStore(db)
```

### MySQL

```go
import "gorm.io/driver/mysql"

dsn := "sqrl_app:REDACTED@tcp(db.example.com:3306)/sqrl_auth" +
    "?charset=utf8mb4&parseTime=True&loc=Local&tls=custom"

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

### SQL Server

```go
import "gorm.io/driver/sqlserver"

dsn := "sqlserver://sqrl_app:REDACTED@db.example.com:1433" +
    "?database=sqrl_auth&encrypt=true&trustServerCertificate=false"

db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
```

### SQLite (Development Only)

```go
import "gorm.io/driver/sqlite"

db, err := gorm.Open(sqlite.Open("/var/lib/sqrl/auth.db"), &gorm.Config{})
```

> **WARNING:** SQLite is not recommended for production multi-process
> deployments due to write locking. Use PostgreSQL or MySQL for production.

---

## Connection Pool Settings

Configure the underlying `database/sql` connection pool for production
workloads:

```go
sqlDB, err := db.DB()
if err != nil {
    log.Fatalf("failed to get sql.DB: %v", err)
}

// Maximum number of open connections to the database.
// Set based on expected concurrent request volume.
// Rule of thumb: 2-4x the number of CPU cores.
sqlDB.SetMaxOpenConns(25)

// Maximum number of idle connections in the pool.
// Keep some connections warm to avoid reconnect latency.
sqlDB.SetMaxIdleConns(10)

// Maximum lifetime of a connection before it is closed and replaced.
// Prevents stale connections after database failover.
sqlDB.SetConnMaxLifetime(5 * time.Minute)

// Maximum time a connection can sit idle before being closed.
sqlDB.SetConnMaxIdleTime(3 * time.Minute)
```

### Sizing Guidelines

| Deployment | MaxOpenConns | MaxIdleConns | Rationale |
|------------|-------------|-------------|-----------|
| Small (< 100 req/s) | 10 | 5 | Low overhead |
| Medium (100-1000 req/s) | 25 | 10 | Balanced |
| Large (> 1000 req/s) | 50-100 | 25 | High throughput |

> **Note:** Each SQRL authentication request involves 1-2 database operations
> (FindIdentity + optional SaveIdentity). Size the pool accordingly.

---

## TLS Configuration

### Database TLS (In Transit)

**PostgreSQL:**

```text
sslmode=verify-full
sslrootcert=/etc/ssl/certs/ca.pem
sslcert=/etc/ssl/certs/client.pem
sslkey=/etc/ssl/private/client-key.pem
```

**MySQL:**

```go
import (
    "crypto/tls"
    "crypto/x509"
    "os"
    mysqldriver "github.com/go-sql-driver/mysql"
)

rootCertPool := x509.NewCertPool()
pem, _ := os.ReadFile("/etc/ssl/certs/ca.pem")
rootCertPool.AppendCertsFromPEM(pem)

mysqldriver.RegisterTLSConfig("custom", &tls.Config{
    RootCAs:    rootCertPool,
    MinVersion: tls.VersionTLS12,
})
```

### Transparent Data Encryption (At Rest)

The SQRL protocol does not mandate field-level encryption. Use database-level
Transparent Data Encryption (TDE) for data-at-rest protection:

- **PostgreSQL:** Use `pgcrypto` extension or full-disk encryption
- **MySQL:** Enable InnoDB tablespace encryption
- **SQL Server:** Enable TDE via `ALTER DATABASE ... SET ENCRYPTION ON`

### TLS Versions

- **Minimum:** TLS 1.2
- **Recommended:** TLS 1.3

---

## Schema Migration

Run `AutoMigrate` during application startup. In production, use the
context-aware variant with a timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := store.AutoMigrateWithContext(ctx); err != nil {
    log.Fatalf("schema migration failed: %v", err)
}
```

### Migration Safety

- `AutoMigrate` only creates tables and adds missing columns
- It does **not** drop columns, delete data, or change column types
- For destructive schema changes, use a dedicated migration tool
  (e.g., [golang-migrate](https://github.com/golang-migrate/migrate),
  [goose](https://github.com/pressly/goose))

---

## Context and Timeout Configuration

All `*WithContext()` methods accept `context.Context` for timeout and
cancellation control. Use timeouts for all database operations in production:

```go
// Per-request timeout (recommended: 5-10 seconds for SQRL auth)
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()

identity, err := store.FindIdentityWithContext(ctx, idk)
if err != nil {
    // Handle timeout, cancellation, or database errors
}
```

### Recommended Timeouts

| Operation | Timeout | Rationale |
|-----------|---------|-----------|
| `FindIdentityWithContext` | 5s | Single SELECT by primary key |
| `SaveIdentityWithContext` | 5s | Single INSERT/UPDATE |
| `DeleteIdentityWithContext` | 5s | Single DELETE |
| `AutoMigrateWithContext` | 30s | DDL operations may be slower |
| `FindIdentitySecureWithContext` | 5s | Same as FindIdentity |

### Context Propagation

Pass the HTTP request context through to database operations for automatic
cancellation when clients disconnect:

```go
func handleSQRLAuth(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    identity, err := store.FindIdentityWithContext(ctx, idk)
    // Context cancellation propagates from HTTP to DB
}
```

---

## Security Checklist

### Pre-Deployment

- [ ] Database credentials stored in secrets manager (not in code or env vars)
- [ ] TLS enabled for all database connections (`sslmode=verify-full`)
- [ ] TDE enabled on the database for data-at-rest encryption
- [ ] Database user has minimum required privileges (SELECT, INSERT, UPDATE,
      DELETE on `sqrl_identities` table only)
- [ ] Network firewall restricts database access to application servers only
- [ ] Connection string does not appear in logs or error messages

### Application Security

- [ ] `ClearIdentity()` or `SecureIdentityWrapper.Destroy()` called after
      every identity retrieval
- [ ] No logging of sensitive fields (Idk, Suk, Vuk, Pidk)
- [ ] Context timeouts configured for all database operations
- [ ] Error messages do not leak identity key values
- [ ] Race detector clean (`go test -race`)
- [ ] Security scan clean (`gosec ./...`, `govulncheck ./...`)

### Infrastructure Security

- [ ] Database backups encrypted
- [ ] Backup access restricted and audited
- [ ] Database audit logging enabled
- [ ] Connection pool limits configured to prevent resource exhaustion
- [ ] Health checks configured (see below)

---

## Monitoring and Health Checks

### Database Health Check

```go
func healthCheck(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get sql.DB: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    return sqlDB.PingContext(ctx)
}
```

### Connection Pool Metrics

```go
sqlDB, _ := db.DB()
stats := sqlDB.Stats()

// Key metrics to monitor:
// stats.OpenConnections  - current open connections
// stats.InUse            - connections currently in use
// stats.Idle             - idle connections
// stats.WaitCount        - total number of connections waited for
// stats.WaitDuration     - total time blocked waiting for connections
// stats.MaxIdleClosed    - connections closed due to SetMaxIdleConns
// stats.MaxLifetimeClosed - connections closed due to SetConnMaxLifetime
```

### Alerting Thresholds

| Metric | Warning | Critical |
|--------|---------|----------|
| `OpenConnections / MaxOpenConns` | > 80% | > 95% |
| `WaitDuration / WaitCount` (avg wait) | > 100ms | > 500ms |
| Health check ping | > 1s | > 3s or fail |

---

## Backup and Recovery

### Backup Strategy

The `sqrl_identities` table contains cryptographic public keys. While not
secret (public keys by definition), they are critical for user authentication.

**Recommended backup schedule:**

| Frequency | Type | Retention |
|-----------|------|-----------|
| Every 15 minutes | WAL/binlog | 24 hours |
| Daily | Full dump | 30 days |
| Weekly | Full dump | 1 year |

### Recovery Procedure

```bash
# PostgreSQL point-in-time recovery
pg_restore -d sqrl_auth backup_YYYYMMDD.dump

# Verify recovery
psql -d sqrl_auth -c "SELECT COUNT(*) FROM sqrl_identities;"
```

### Data Sensitivity

| Field | Classification | Backup Handling |
|-------|---------------|-----------------|
| Idk | Public key | Standard backup |
| Suk | Public key (unlock) | Standard backup |
| Vuk | Public key (verify) | Standard backup |
| Pidk | Public key (previous) | Standard backup |
| Status flags | Non-sensitive | Standard backup |

> **Note:** SQRL identity keys are Ed25519 **public** keys. They do not
> require the same protection level as private keys, but their loss would
> prevent users from authenticating.

---

## Example: Complete Production Setup

```go
package main

import (
    "context"
    "log"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    gormauthstore "github.com/sqrldev/server-go-ssp-gormauthstore"
)

func main() {
    // 1. Connect with TLS
    dsn := "host=db.example.com user=sqrl_app " +
        "password=REDACTED dbname=sqrl_auth " +
        "sslmode=verify-full sslrootcert=/etc/ssl/certs/ca.pem"

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Warn),
    })
    if err != nil {
        log.Fatalf("database connection failed: %v", err)
    }

    // 2. Configure connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    // 3. Create store and migrate
    store := gormauthstore.NewAuthStore(db)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := store.AutoMigrateWithContext(ctx); err != nil {
        log.Fatalf("migration failed: %v", err)
    }

    // 4. Use in request handlers with context
    // identity, err := store.FindIdentityWithContext(r.Context(), idk)
    // defer gormauthstore.ClearIdentity(identity)

    log.Println("gormauthstore ready for production")
}
```

---

**Document Control:**

- Version: 1.0
- Date: February 7, 2026
- Next Review: Quarterly
