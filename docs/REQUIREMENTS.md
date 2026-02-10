# Requirements Specification
## SQRL GORM Authentication Store (Reverse Engineered)

**Document Version:** 2.0
**Analysis Date:** November 18, 2025 (updated February 8, 2026)
**Project:** github.com/dxcSithLord/server-go-ssp-gormauthstore
**Analysis Method:** Reverse engineering from existing implementation

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Business Context and Objectives](#business-context-and-objectives)
3. [Functional Requirements](#functional-requirements)
4. [Non-Functional Requirements](#non-functional-requirements)
5. [Interface Requirements](#interface-requirements)
6. [Data Requirements](#data-requirements)
7. [Security Requirements](#security-requirements)
8. [Quality Attributes](#quality-attributes)
9. [Constraints and Assumptions](#constraints-and-assumptions)
10. [Traceability Matrix](#traceability-matrix)

---

## Executive Summary

### Purpose
This library provides persistent storage capabilities for SQRL (Secure QR Login) authentication identities using GORM ORM, enabling Go applications to implement SQRL authentication with database-backed identity management.

### Scope
- **In Scope:** Database persistence for SQRL identities, CRUD operations, schema migration
- **Out of Scope:** SQRL protocol implementation, HTTP endpoints, cryptographic operations, user interface

### Key Stakeholders
- **Primary Users:** Go developers implementing SQRL authentication
- **Consumers:** `github.com/dxcSithLord/server-go-ssp` (SQRL protocol implementation)
- **End Users:** Web applications requiring passwordless authentication

---

## Business Context and Objectives

### Business Problem
Organizations require secure, passwordless authentication systems that:
- Eliminate password-related vulnerabilities
- Provide privacy-preserving authentication
- Support user sovereignty over identity data
- Enable cross-device authentication without data synchronization

### Solution Approach
Implement a database abstraction layer that stores SQRL cryptographic identities, allowing the SQRL Server-Side Protocol implementation to persist user identities across sessions while maintaining the security properties of the SQRL protocol.

### Business Objectives

| Objective ID | Description | Success Criteria | Priority |
|--------------|-------------|------------------|----------|
| BO-001 | Enable SQRL authentication in Go applications | Library used in production SQRL implementations | HIGH |
| BO-002 | Provide database-agnostic storage | Support PostgreSQL, MySQL, SQLite via GORM | HIGH |
| BO-003 | Ensure cryptographic key security | No keys leaked in logs or error messages | CRITICAL |
| BO-004 | Minimize integration complexity | < 10 lines of code to integrate | MEDIUM |
| BO-005 | Support production-grade deployments | Concurrent access, transaction safety | HIGH |

---

## Functional Requirements

### FR-001: Identity Storage
**Objective:** Persist SQRL identity records to database

**Requirements:**
- **FR-001.1:** Store SqrlIdentity structures with all fields
- **FR-001.2:** Use database transactions for data consistency
- **FR-001.3:** Support upsert operations (insert or update)
- **FR-001.4:** Return error on database failure

**Implementation:** `SaveIdentity(*SqrlIdentity) error`

**Test Criteria:**
- Identity successfully saved to database
- Duplicate Idk results in update, not insert
- Transaction rolled back on error
- Error returned with context on failure

**Trace:** BO-001, BO-002

---

### FR-002: Identity Retrieval
**Objective:** Retrieve SQRL identity by Identity Key

**Requirements:**
- **FR-002.1:** Query by Identity Key (Idk) field
- **FR-002.2:** Return `ssp.ErrNotFound` when identity doesn't exist
- **FR-002.3:** Return database errors with context
- **FR-002.4:** Use parameterized queries (SQL injection prevention)

**Implementation:** `FindIdentity(idk string) (*SqrlIdentity, error)`

**Test Criteria:**
- Existing identity retrieved successfully
- Non-existent Idk returns `ssp.ErrNotFound`
- SQL injection attempts safely handled
- Database errors propagated with context

**Trace:** BO-001, BO-003

---

### FR-003: Identity Deletion
**Objective:** Remove SQRL identity from database

**Requirements:**
- **FR-003.1:** Delete identity by Identity Key (Idk)
- **FR-003.2:** No error if identity doesn't exist (idempotent)
- **FR-003.3:** Use parameterized queries
- **FR-003.4:** Physically delete (not soft delete)

**Implementation:** `DeleteIdentity(idk string) error`

**Test Criteria:**
- Identity deleted successfully
- Subsequent FindIdentity returns `ssp.ErrNotFound`
- Multiple deletions don't cause errors
- Deletion is permanent

**Trace:** BO-001, BO-003

---

### FR-004: Schema Management
**Objective:** Automatically create/update database schema

**Requirements:**
- **FR-004.1:** Create `sqrl_identities` table if not exists
- **FR-004.2:** Add missing columns on schema evolution
- **FR-004.3:** Preserve existing data during migration
- **FR-004.4:** Support GORM AutoMigrate functionality

**Implementation:** `AutoMigrate() error`

**Test Criteria:**
- Table created on first run
- Columns match SqrlIdentity struct
- Existing data preserved on migration
- Indexes created for performance

**Trace:** BO-002, BO-004

---

### FR-005: AuthStore Interface Compliance
**Objective:** Implement ssp.AuthStore interface contract

**Requirements:**
- **FR-005.1:** Implement all methods of `ssp.AuthStore`
- **FR-005.2:** Maintain interface compatibility with `server-go-ssp`
- **FR-005.3:** Follow semantic versioning for API changes

**Implementation:** `type AuthStore struct`

**Test Criteria:**
- `var _ ssp.AuthStore = (*AuthStore)(nil)` compiles
- Integration tests with `server-go-ssp` pass
- No breaking changes without major version bump

**Trace:** BO-001, BO-004

---

### FR-006: Context Support
**Objective:** Enable timeout and cancellation control for database operations

**Requirements:**
- **FR-006.1:** All CRUD methods have `*WithContext(ctx)` variants
- **FR-006.2:** Original methods delegate to WithContext with `context.Background()`
- **FR-006.3:** Cancelled context returns error before database operation

**Implementation:**
- `FindIdentityWithContext(ctx, idk)`
- `SaveIdentityWithContext(ctx, identity)`
- `DeleteIdentityWithContext(ctx, idk)`
- `FindIdentitySecureWithContext(ctx, idk)`
- `AutoMigrateWithContext(ctx)`

**Trace:** BO-005

---

### FR-007: Secure Identity Retrieval
**Objective:** Provide RAII-style automatic cleanup of sensitive cryptographic material

**Requirements:**
- **FR-007.1:** Return identity wrapped in `SecureIdentityWrapper`
- **FR-007.2:** `Destroy()` clears all sensitive fields
- **FR-007.3:** Accessing destroyed wrapper returns `ErrWrappedIdentityDestroyed`

**Implementation:** `FindIdentitySecure(idk)`, `FindIdentitySecureWithContext(ctx, idk)`

**Trace:** BO-003

---

### FR-008: Input Validation
**Objective:** Validate inputs before database operations

**Requirements:**
- **FR-008.1:** Reject empty identity keys (`ErrEmptyIdentityKey`)
- **FR-008.2:** Reject keys exceeding 256 characters (`ErrIdentityKeyTooLong`)
- **FR-008.3:** Reject keys with invalid characters (`ErrInvalidIdentityKeyFormat`)
- **FR-008.4:** Reject nil identity on save (`ErrNilIdentity`)

**Implementation:** `ValidateIdk()` called by all CRUD methods

**Trace:** BO-003, SEC-002, SEC-003

---

## Non-Functional Requirements

### NFR-001: Performance
**Objective:** Minimize latency impact on authentication flow

| Metric | Requirement | Measurement |
|--------|-------------|-------------|
| FindIdentity latency | < 50ms (p95) | Benchmark tests |
| SaveIdentity latency | < 100ms (p95) | Benchmark tests |
| Concurrent requests | 100 req/s per instance | Load tests |
| Database connections | Configurable pool size | Connection monitoring |

**Trace:** BO-005

---

### NFR-002: Reliability
**Objective:** Ensure consistent operation under failure conditions

| Attribute | Requirement | Verification |
|-----------|-------------|--------------|
| Data consistency | ACID transaction compliance | Integration tests |
| Error handling | All errors returned with context | Code review |
| Graceful degradation | No panics on database failure | Fault injection tests |
| Idempotency | SaveIdentity, DeleteIdentity idempotent | Idempotency tests |

**Trace:** BO-005

---

### NFR-003: Security
**Objective:** Protect cryptographic material and prevent attacks

| Control | Requirement | Implementation |
|---------|-------------|----------------|
| SQL Injection | All queries parameterized | GORM ORM usage |
| Sensitive data logging | No Suk/Vuk in logs/errors | Code audit |
| Memory clearing | Clear sensitive data after use | secure_memory.go |
| TLS enforcement | Database connections over TLS (configurable) | Connection string |

**Trace:** BO-003

**Note:** See SECURITY_REQUIREMENTS.md for detailed security controls

---

### NFR-004: Maintainability
**Objective:** Enable long-term maintenance and evolution

| Attribute | Requirement | Measurement |
|-----------|-------------|-------------|
| Code coverage | ≥ 70% line coverage | `go test -cover` |
| Cyclomatic complexity | < 10 per function | `gocyclo` |
| Documentation | All exported symbols documented | `golint` |
| Dependencies | Minimize external dependencies | Dependency graph |

**Trace:** BO-001

---

### NFR-005: Compatibility
**Objective:** Support diverse deployment environments

| Aspect | Requirement | Verification |
|--------|-------------|--------------|
| Go version | Go 1.25+ | CI/CD matrix |
| Databases | PostgreSQL 12+, MySQL 8+, SQLite 3.35+ | Integration tests |
| Operating systems | Linux, Windows, macOS | CI/CD matrix |
| Architectures | amd64, arm64 | Build verification |

**Trace:** BO-002

---

## Interface Requirements

### INT-001: Constructor Interface
```go
func NewAuthStore(db *gorm.DB) *AuthStore
```

**Purpose:** Initialize AuthStore with existing GORM database connection

**Preconditions:**
- `db` is non-nil and connected
- Database driver loaded (postgres, mysql, sqlite)

**Postconditions:**
- Returns configured AuthStore instance
- No database operations performed
- AuthStore ready for use

---

### INT-002: ssp.AuthStore Interface Compliance

**Interface Definition:**
```go
type AuthStore interface {
    FindIdentity(idk string) (*SqrlIdentity, error)
    SaveIdentity(identity *SqrlIdentity) error
    DeleteIdentity(idk string) error
}
```

**Contract:**
- Methods must be safe for concurrent use
- Errors must implement `error` interface
- `ssp.ErrNotFound` must be returned for missing identities
- All database errors propagated

---

## Data Requirements

### SqrlIdentity Data Model

**Structure:**
```go
type SqrlIdentity struct {
    Idk      string  `gorm:"primary_key"`  // Identity Key (43-44 chars base64url)
    Suk      string  // Server Unlock Key (cryptographic material)
    Vuk      string  // Verify Unlock Key (cryptographic material)
    Pidk     string  // Previous Identity Key (for rekeying)
    SQRLOnly bool    // SQRL-only authentication flag
    Hardlock bool    // Hard lock status
    Disabled bool    // Account disabled status
    Rekeyed  string  // Link to new identity (Idk)
    Btn      int     // User button response (0-3)
}
```

### Data Constraints

| Field | Constraint | Rationale |
|-------|------------|-----------|
| Idk | PRIMARY KEY, NOT NULL, UNIQUE | Identity uniqueness |
| Idk | Length: 43-44 characters | Base64url(32 bytes) |
| Idk | Character set: [A-Za-z0-9_-] | URL-safe base64 |
| Suk | NOT NULL | Required for authentication |
| Vuk | NOT NULL | Required for authentication |
| Pidk | Nullable | Only set during rekeying |
| Rekeyed | Nullable | Only set when identity replaced |
| Btn | Range: 0-3 | SQRL protocol constraint |

### Database Schema

**Table Name:** `sqrl_identities`

**Indexes:**
- Primary: `idk` (clustered)
- Optional: `rekeyed` (for identity chain navigation)

**Size Estimates:**
- Row size: ~500 bytes (with cryptographic keys)
- 10,000 users: ~5 MB
- 1,000,000 users: ~500 MB

---

## Security Requirements

### SEC-001: Cryptographic Key Protection (CWE-226, CWE-200, CWE-244)

**Threat:** Sensitive cryptographic keys (Suk, Vuk) leaked via memory dumps, logs, or error messages

**Controls:**
1. **SEC-001.1:** No sensitive fields in error messages
2. **SEC-001.2:** No sensitive fields in debug/trace logs
3. **SEC-001.3:** Memory clearing after identity use (defense-in-depth)
4. **SEC-001.4:** TLS for database connections (encryption in transit)

**Implementation:**
- `secure_memory.go`: WipeBytes, ClearIdentity functions
- `auth_store.go`: Error messages exclude sensitive data
- Documentation: Recommend TLS connection strings

**Verification:**
- Security code review
- Error message audit
- Memory analysis (manual testing)

**Priority:** CRITICAL
**Trace:** BO-003

---

### SEC-002: SQL Injection Prevention (CWE-89)

**Threat:** Attacker manipulates Idk parameter to execute arbitrary SQL

**Controls:**
1. **SEC-002.1:** All queries use parameterized statements
2. **SEC-002.2:** GORM ORM prevents direct SQL concatenation
3. **SEC-002.3:** Input validation (character whitelist)

**Implementation:**
- GORM's `.Where("idk = ?", idk)` syntax
- `ValidateIdk()` function validates character set

**Verification:**
- SQL injection test suite
- Static analysis (gosec)

**Priority:** CRITICAL
**Trace:** BO-003

---

### SEC-003: Denial of Service Prevention (CWE-400)

**Threat:** Resource exhaustion via database connection pool or unbounded queries

**Controls:**
1. **SEC-003.1:** Connection pool limits configured
2. **SEC-003.2:** Query timeouts enforced
3. **SEC-003.3:** Input length validation (Idk max 256 chars)

**Implementation:**
- GORM connection pool configuration
- `MaxIdkLength` constant (256)

**Verification:**
- Load testing
- Resource monitoring

**Priority:** HIGH
**Trace:** BO-005

---

### SEC-004: Data at Rest Protection

**Threat:** Database compromise exposes SQRL keys

**Controls:**
1. **SEC-004.1:** Recommend database encryption (TDE)
2. **SEC-004.2:** Document encryption-at-rest best practices
3. **SEC-004.3:** Consider field-level encryption (future enhancement)

**Implementation:**
- Documentation in README
- Production deployment guide

**Verification:**
- Documentation review

**Priority:** MEDIUM
**Note:** Database encryption is deployment-specific, not library responsibility

---

## Quality Attributes

### Testability
- All functions unit testable
- Database mockable via GORM interfaces
- Test coverage measured and enforced (≥70%)

### Observability
- All operations return errors with context
- Compatible with logging wrappers
- Performance metrics collectable

### Extensibility
- Interface-based design
- Additional databases supported via GORM drivers
- Schema evolution via AutoMigrate

---

## Constraints and Assumptions

### Technical Constraints
1. **CON-001:** Must use GORM v2+ (gorm.io/gorm)
2. **CON-002:** Requires Go 1.25 or later
3. **CON-003:** Database must support transactions
4. **CON-004:** Identity Keys assumed URL-safe base64 encoded

### Business Constraints
1. **CON-005:** Open source (MIT License)
2. **CON-006:** No external service dependencies
3. **CON-007:** Must remain compatible with `server-go-ssp` interface

### Assumptions
1. **ASM-001:** Database driver loaded by calling application
2. **ASM-002:** GORM connection configured by calling application
3. **ASM-003:** Calling code handles concurrent request serialization if needed
4. **ASM-004:** Identity Keys generated by SQRL client (not validated cryptographically)
5. **ASM-005:** Database provides persistence guarantees (not this library's responsibility)

---

## Traceability Matrix

### Business Objectives → Functional Requirements

| Business Objective | Functional Requirements |
|-------------------|------------------------|
| BO-001 (Enable SQRL) | FR-001, FR-002, FR-003, FR-004, FR-005 |
| BO-002 (Database agnostic) | FR-004, FR-005 |
| BO-003 (Key security) | FR-002, FR-003, SEC-001, SEC-002 |
| BO-004 (Simple integration) | FR-005, INT-001 |
| BO-005 (Production grade) | NFR-001, NFR-002, SEC-003, FR-006 |

### Functional Requirements → Code Implementation

| Functional Requirement | Implementation | Location |
|------------------------|----------------|----------|
| FR-001 (Storage) | `SaveIdentity()` | auth_store.go:37-39 |
| FR-002 (Retrieval) | `FindIdentity()` | auth_store.go:24-34 |
| FR-003 (Deletion) | `DeleteIdentity()` | auth_store.go:42-44 |
| FR-004 (Schema) | `AutoMigrate()` | auth_store.go:19-21 |
| FR-005 (Interface) | `type AuthStore` | auth_store.go:72-74 |
| FR-006 (Context) | `*WithContext()` methods | auth_store.go (5 methods) |
| FR-007 (Secure) | `FindIdentitySecure()` | auth_store.go:150-171 |
| FR-008 (Validation) | `ValidateIdk()` | secure_memory_common.go |

### Security Controls → Implementation

| Security Control | Implementation | Location |
|------------------|----------------|----------|
| SEC-001 (Key protection) | `ClearIdentity()`, `WipeBytes()` | secure_memory_common.go:86-105 |
| SEC-002 (SQL injection) | GORM parameterized queries | auth_store.go (all methods) |
| SEC-003 (DoS) | `ValidateIdk()`, `MaxIdkLength` | secure_memory_common.go:165-197 |

---

## Document Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Technical Lead | [To be assigned] | | |
| Security Architect | [To be assigned] | | |
| Product Owner | [To be assigned] | | |

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2025-11-18 | Automated Analysis | Initial reverse-engineered requirements |
| 2.0 | 2026-02-08 | Review | Added FR-006/007/008, updated versions and paths |

---

## References

1. SQRL Protocol Specification: https://www.grc.com/sqrl/sqrl.htm
2. GORM Documentation: https://gorm.io/docs/
3. CWE-226: https://cwe.mitre.org/data/definitions/226.html
4. CWE-200: https://cwe.mitre.org/data/definitions/200.html
5. CWE-89: https://cwe.mitre.org/data/definitions/89.html
6. Go Secure Coding Practices: https://go.dev/doc/security/

---

**END OF REQUIREMENTS SPECIFICATION**
