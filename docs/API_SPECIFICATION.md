# API Specification
## SQRL GORM Authentication Store
## OpenAPI 3.1-Style Documentation for Go Interface

**Specification Version:** 2.0.0
**API Version:** v1.0.0-rc
**Format:** OpenAPI 3.1 concepts applied to Go interfaces
**Date:** February 8, 2026

---

## Table of Contents

1. [Overview](#overview)
2. [Interface Contract](#interface-contract)
3. [Data Models](#data-models)
4. [Operations](#operations)
5. [Error Responses](#error-responses)
6. [Security](#security)
7. [Examples](#examples)
8. [Testing Specification](#testing-specification)

---

## Overview

### API Description

This library provides a Go interface (`ssp.AuthStore`) for persisting SQRL authentication identities in SQL databases via GORM ORM. While not a REST/HTTP API, this specification documents the programmatic interface using OpenAPI 3.1 concepts.

### Base Information

```yaml
openapi: 3.1.0
info:
  title: SQRL GORM Authentication Store API
  description: Database persistence layer for SQRL identities
  version: 1.0.0
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT
  contact:
    name: SQRL Development Team
    url: https://github.com/sqrldev/server-go-ssp-gormauthstore

externalDocs:
  description: SQRL Protocol Specification
  url: https://www.grc.com/sqrl/sqrl.htm
```

### Supported Databases

| Database | Minimum Version | Driver Package |
|----------|----------------|----------------|
| PostgreSQL | 12.0 | gorm.io/driver/postgres v1.5.9 |
| MySQL | 8.0 | gorm.io/driver/mysql v1.5.7 |
| SQLite | 3.35.0 | gorm.io/driver/sqlite v1.6.0 |
| SQL Server | 2019 | gorm.io/driver/sqlserver v1.5.3 |

---

## Interface Contract

### ssp.AuthStore Interface

```yaml
components:
  interfaces:
    AuthStore:
      description: Interface for SQRL identity persistence
      methods:
        - FindIdentity
        - SaveIdentity
        - DeleteIdentity
      extended_methods:
        - FindIdentityWithContext
        - SaveIdentityWithContext
        - DeleteIdentityWithContext
        - FindIdentitySecure
        - FindIdentitySecureWithContext
        - AutoMigrateWithContext

  types:
    AuthStore:
      type: object
      description: GORM-based implementation of ssp.AuthStore
      properties:
        db:
          $ref: '#/components/schemas/GormDB'
      methods:
        NewAuthStore:
          type: constructor
          parameters:
            - name: db
              required: true
              schema:
                $ref: '#/components/schemas/GormDB'
          returns:
            - type: '*AuthStore'

        AutoMigrate:
          type: method
          description: Create or update database schema
          parameters: []
          returns:
            - type: error
              description: Database error if migration fails

        FindIdentity:
          type: method
          description: Retrieve SQRL identity by Identity Key
          parameters:
            - name: idk
              in: parameter
              required: true
              schema:
                type: string
                minLength: 1
                maxLength: 256
                pattern: '^[A-Za-z0-9\+/=\-_\.]+$'
              description: Identity Key (URL-safe base64)
          returns:
            - type: '*ssp.SqrlIdentity'
              description: Found identity
            - type: error
              description: ssp.ErrNotFound if not exists, or database error
          security:
            - sensitive_data: [Suk, Vuk]

        SaveIdentity:
          type: method
          description: Insert or update SQRL identity
          parameters:
            - name: identity
              in: parameter
              required: true
              schema:
                $ref: '#/components/schemas/SqrlIdentity'
          returns:
            - type: error
              description: Database error if save fails
          security:
            - input_validation: [Idk, Suk, Vuk]

        DeleteIdentity:
          type: method
          description: Remove SQRL identity from database
          parameters:
            - name: idk
              in: parameter
              required: true
              schema:
                type: string
                minLength: 1
                maxLength: 256
          returns:
            - type: error
              description: Database error if delete fails (no error if not found)
          idempotent: true
```

---

## Data Models

### SqrlIdentity Schema

```yaml
components:
  schemas:
    SqrlIdentity:
      type: object
      description: SQRL cryptographic identity record
      required:
        - Idk
        - Suk
        - Vuk
      properties:
        Idk:
          type: string
          format: base64url
          minLength: 43
          maxLength: 44
          pattern: '^[A-Za-z0-9\-_]{43,44}$'
          description: Identity Key (unique identifier)
          example: "k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E"
          x-sensitive: true
          x-primary-key: true

        Suk:
          type: string
          format: base64url
          description: Server Unlock Key (cryptographic material)
          example: "M2n5P8r1T5v9Y3z7C1e5G9k3N7q1S5w9A3d7F1h5J9"
          x-sensitive: true
          x-security-level: critical

        Vuk:
          type: string
          format: base64url
          description: Verify Unlock Key (cryptographic material)
          example: "Q5t9W3a7D1g5K9n3R7u1X5b9E3h7L1p5T9y3C7f1I5"
          x-sensitive: true
          x-security-level: critical

        Pidk:
          type: string
          format: base64url
          nullable: true
          description: Previous Identity Key (for rekeying)
          example: "B9f3J7m1P5s9V3y7A1d5G9k3N7q1T5w9C3e7H1l5O9"
          x-sensitive: true

        SQRLOnly:
          type: boolean
          default: false
          description: SQRL-only authentication flag
          example: false

        Hardlock:
          type: boolean
          default: false
          description: Hard lock status (prevents non-SQRL auth)
          example: false

        Disabled:
          type: boolean
          default: false
          description: Account disabled status
          example: false

        Rekeyed:
          type: string
          format: base64url
          nullable: true
          description: New identity Idk if this identity was rekeyed
          example: null
          x-foreign-key: SqrlIdentity.Idk

        Btn:
          type: integer
          format: int32
          minimum: 0
          maximum: 3
          default: 0
          description: User button response (0-3)
          example: 0

      example:
        Idk: "k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E"
        Suk: "M2n5P8r1T5v9Y3z7C1e5G9k3N7q1S5w9A3d7F1h5J9"
        Vuk: "Q5t9W3a7D1g5K9n3R7u1X5b9E3h7L1p5T9y3C7f1I5"
        Pidk: null
        SQRLOnly: false
        Hardlock: false
        Disabled: false
        Rekeyed: null
        Btn: 0
```

### Error Models

```yaml
components:
  schemas:
    Error:
      type: object
      description: Standard Go error interface
      properties:
        message:
          type: string
          description: Error message (sanitized)
        code:
          type: string
          enum:
            - ErrNotFound
            - ErrEmptyIdentityKey
            - ErrIdentityKeyTooLong
            - ErrInvalidIdentityKeyFormat
            - ErrNilIdentity
            - ErrNilDatabase
            - ErrWrappedIdentityDestroyed

    ValidationError:
      allOf:
        - $ref: '#/components/schemas/Error'
        - type: object
          properties:
            field:
              type: string
              description: Field that failed validation
            value:
              type: string
              description: Sanitized representation of invalid value
              x-sensitive: never-include-actual-value
```

---

## Operations

### Operation: NewAuthStore

**Constructor for AuthStore**

```yaml
paths:
  /authstore/new:  # Conceptual path
    post:
      summary: Create new AuthStore instance
      operationId: NewAuthStore
      description: |
        Initializes a new AuthStore with the provided GORM database connection.
        The database connection must be established before calling this constructor.
      parameters:
        - name: db
          in: parameter
          required: true
          schema:
            type: object
            description: '*gorm.DB instance'
      responses:
        '200':
          description: AuthStore created successfully
          content:
            application/go:
              schema:
                type: object
                properties:
                  instance:
                    type: '*AuthStore'

      x-code-example: |
        import (
            "gorm.io/driver/postgres"
            "gorm.io/gorm"
            "github.com/sqrldev/server-go-ssp-gormauthstore"
        )

        dsn := "host=localhost user=postgres dbname=sqrl_db sslmode=require"
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            panic(err)
        }

        store := gormauthstore.NewAuthStore(db)
```

---

### Operation: AutoMigrate

**Create or update database schema**

```yaml
paths:
  /authstore/migrate:  # Conceptual path
    post:
      summary: Create or update database schema
      operationId: AutoMigrate
      description: |
        Uses GORM's AutoMigrate to create the sqrl_identities table or add missing columns.
        Safe to call on every application start.
      responses:
        '200':
          description: Schema migrated successfully
        '500':
          description: Database error
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/Error'

      x-code-example: |
        err := store.AutoMigrate()
        if err != nil {
            log.Fatalf("Migration failed: %v", err)
        }
```

---

### Operation: FindIdentity

**Retrieve SQRL identity by Idk**

```yaml
paths:
  /identity/{idk}:  # Conceptual path
    get:
      summary: Find SQRL identity by Identity Key
      operationId: FindIdentity
      description: |
        Queries the database for a SQRL identity matching the provided Idk.
        Returns ssp.ErrNotFound if identity doesn't exist.
      parameters:
        - name: idk
          in: path
          required: true
          schema:
            type: string
            minLength: 1
            maxLength: 256
            pattern: '^[A-Za-z0-9\+/=\-_\.]+$'
          description: Identity Key (URL-safe base64)
          examples:
            valid:
              value: "k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E"
              summary: Valid 43-character Idk

      responses:
        '200':
          description: Identity found
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/SqrlIdentity'
        '404':
          description: Identity not found
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                code: "ErrNotFound"
                message: "Not Found"
        '400':
          description: Invalid Idk format
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/ValidationError'
        '500':
          description: Database error
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/Error'

      security:
        - database_auth: []

      x-code-example: |
        identity, err := store.FindIdentity("k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E")
        if err != nil {
            if err == ssp.ErrNotFound {
                // Handle not found
                return
            }
            // Handle database error
            log.Printf("Database error: %v", err)
            return
        }

        // Use identity (remember to clear sensitive data)
        defer ClearIdentity(identity)

        // ... use identity.Suk, identity.Vuk ...
```

---

### Operation: SaveIdentity

**Insert or update SQRL identity**

```yaml
paths:
  /identity:  # Conceptual path
    put:
      summary: Save SQRL identity (upsert)
      operationId: SaveIdentity
      description: |
        Inserts a new SQRL identity or updates an existing one (based on Idk).
        Uses GORM's Save method which performs upsert based on primary key.
      requestBody:
        required: true
        content:
          application/go:
            schema:
              $ref: '#/components/schemas/SqrlIdentity'

      responses:
        '200':
          description: Identity saved successfully
        '400':
          description: Validation error
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/ValidationError'
        '500':
          description: Database error
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/Error'

      security:
        - database_auth: []

      x-code-example: |
        identity := &ssp.SqrlIdentity{
            Idk:      "k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E",
            Suk:      "M2n5P8r1T5v9Y3z7C1e5G9k3N7q1S5w9A3d7F1h5J9",
            Vuk:      "Q5t9W3a7D1g5K9n3R7u1X5b9E3h7L1p5T9y3C7f1I5",
            SQRLOnly: false,
            Hardlock: false,
            Disabled: false,
        }

        err := store.SaveIdentity(identity)
        if err != nil {
            log.Printf("Save failed: %v", err)
            return
        }
```

---

### Operation: DeleteIdentity

**Remove SQRL identity from database**

```yaml
paths:
  /identity/{idk}:  # Conceptual path
    delete:
      summary: Delete SQRL identity
      operationId: DeleteIdentity
      description: |
        Permanently removes the SQRL identity with the specified Idk.
        Idempotent - no error if identity doesn't exist.
      parameters:
        - name: idk
          in: path
          required: true
          schema:
            type: string
            minLength: 1
            maxLength: 256
          description: Identity Key to delete

      responses:
        '200':
          description: Identity deleted successfully (or didn't exist)
        '400':
          description: Invalid Idk format
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/ValidationError'
        '500':
          description: Database error
          content:
            application/go:
              schema:
                $ref: '#/components/schemas/Error'

      security:
        - database_auth: []

      x-code-example: |
        err := store.DeleteIdentity("k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E")
        if err != nil {
            log.Printf("Delete failed: %v", err)
            return
        }
        // Success - identity deleted (or didn't exist)
```

---

## Error Responses

### Standard Errors

| Error Code | Go Constant | HTTP Equivalent | Description |
|------------|-------------|----------------|-------------|
| `ErrNotFound` | `ssp.ErrNotFound` | 404 | Identity not found |
| `ErrEmptyIdentityKey` | `gormauthstore.ErrEmptyIdentityKey` | 400 | Idk is empty string |
| `ErrIdentityKeyTooLong` | `gormauthstore.ErrIdentityKeyTooLong` | 400 | Idk exceeds 256 chars |
| `ErrInvalidIdentityKeyFormat` | `gormauthstore.ErrInvalidIdentityKeyFormat` | 400 | Idk contains invalid characters |
| `ErrNilIdentity` | `gormauthstore.ErrNilIdentity` | 400 | Nil identity passed to SaveIdentity |
| `ErrNilDatabase` | `gormauthstore.ErrNilDatabase` | 500 | Database connection is nil |
| `ErrWrappedIdentityDestroyed` | `gormauthstore.ErrWrappedIdentityDestroyed` | 500 | SecureIdentityWrapper already destroyed |
| `gorm.ErrRecordNotFound` | `gorm.ErrRecordNotFound` | 404 | Mapped to ssp.ErrNotFound internally |

### Error Handling Pattern

```go
identity, err := store.FindIdentity(idk)
if err != nil {
    switch {
    case errors.Is(err, ssp.ErrNotFound):
        // Handle not found (404)
        return nil, fmt.Errorf("identity not found")

    case errors.Is(err, gormauthstore.ErrEmptyIdentityKey):
        // Handle validation error (400)
        return nil, fmt.Errorf("invalid input: %w", err)

    case errors.Is(err, gormauthstore.ErrIdentityKeyTooLong):
        // Handle validation error (400)
        return nil, fmt.Errorf("invalid input: %w", err)

    case errors.Is(err, gormauthstore.ErrInvalidIdentityKeyFormat):
        // Handle validation error (400)
        return nil, fmt.Errorf("invalid input: %w", err)

    default:
        // Handle database error (500)
        return nil, fmt.Errorf("database error: %w", err)
    }
}
```

---

## Security

### Security Schemes

```yaml
components:
  securitySchemes:
    database_auth:
      type: database
      description: Database connection authentication (PostgreSQL, MySQL, etc.)
      x-credentials:
        - username
        - password
        - ssl_cert
        - ssl_key
```

### Security Considerations

| Threat | Mitigation | Implementation |
|--------|-----------|----------------|
| **SQL Injection** | Parameterized queries | GORM ORM (`.Where("idk = ?", idk)`) |
| **Sensitive Data Exposure** | Memory clearing | `ClearIdentity()`, `WipeBytes()` |
| **DoS (Length Attacks)** | Input validation | `MaxIdkLength` (256 chars) |
| **Character Injection** | Character whitelist | `isValidIdkChar()` validation |
| **Database Compromise** | Encryption at rest | TDE (deployment-specific) |
| **Man-in-the-Middle** | TLS for DB connections | `sslmode=require` in DSN |

### Data Sensitivity Levels

```yaml
x-data-classification:
  critical:
    - Suk  # Server Unlock Key
    - Vuk  # Verify Unlock Key
    description: Cryptographic keys - never log, clear after use

  sensitive:
    - Idk   # Identity Key
    - Pidk  # Previous Identity Key
    - Rekeyed
    description: Identifiers - sanitize in errors, clear after use

  confidential:
    - SQRLOnly
    - Hardlock
    - Disabled
    - Btn
    description: Account status - safe to log
```

---

## Examples

### Complete Usage Example

```go
package main

import (
    "log"

    ssp "github.com/sqrldev/server-go-ssp"
    "github.com/sqrldev/server-go-ssp-gormauthstore"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    // 1. Connect to database
    dsn := "host=localhost user=postgres password=secret dbname=sqrl_prod sslmode=require"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Connection failed:", err)
    }

    // 2. Create AuthStore
    store := gormauthstore.NewAuthStore(db)

    // 3. Run migrations
    if err := store.AutoMigrate(); err != nil {
        log.Fatal("Migration failed:", err)
    }

    // 4. Save an identity
    identity := &ssp.SqrlIdentity{
        Idk: "k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E",
        Suk: "M2n5P8r1T5v9Y3z7C1e5G9k3N7q1S5w9A3d7F1h5J9",
        Vuk: "Q5t9W3a7D1g5K9n3R7u1X5b9E3h7L1p5T9y3C7f1I5",
    }

    if err := store.SaveIdentity(identity); err != nil {
        log.Fatal("Save failed:", err)
    }

    // 5. Find an identity
    found, err := store.FindIdentity("k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E")
    if err != nil {
        if err == ssp.ErrNotFound {
            log.Println("Identity not found")
            return
        }
        log.Fatal("Find failed:", err)
    }

    // 6. Use identity (with secure cleanup)
    defer gormauthstore.ClearIdentity(found)
    log.Printf("Found identity: %s", found.Idk)  // Safe to log Idk
    // Use found.Suk and found.Vuk...

    // 7. Delete identity
    if err := store.DeleteIdentity("k1vMZ8C9B2Q8h5K3x7N9m4P6w8R1t5Y2u9Z3v7C1d4E"); err != nil {
        log.Fatal("Delete failed:", err)
    }
}
```

### Secure Wrapper Example

```go
// Use SecureIdentityWrapper for automatic cleanup
wrapper, err := store.FindIdentitySecure(idk)
if err != nil {
    return err
}
defer wrapper.Destroy()

// Access via wrapper
identity := wrapper.GetIdentity()
if identity == nil {
    return errors.New("wrapper invalidated")
}

// Use identity.Suk, identity.Vuk...
// Automatic cleanup on function return
```

### Context-Aware Example

```go
// Use context for timeout and cancellation control
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()

identity, err := store.FindIdentityWithContext(ctx, idk)
if err != nil {
    // Handle timeout, cancellation, or database errors
    return err
}
defer gormauthstore.ClearIdentity(identity)
```

---

## Testing Specification

### Test Coverage Requirements

| Component | Target Coverage | Current Coverage |
|-----------|----------------|------------------|
| auth_store.go | ≥80% | 100% |
| secure_memory_common.go | ≥90% | 100% |
| secure_memory.go | ≥95% | 100% |
| errors.go | 100% | 100% |
| **Overall** | **≥70%** | **100%** |

### Test Categories

```yaml
x-test-suite:
  unit_tests:
    - TestNewAuthStore
    - TestAutoMigrate
    - TestFindIdentity_Found
    - TestFindIdentity_NotFound
    - TestFindIdentity_InvalidIdk
    - TestSaveIdentity_Insert
    - TestSaveIdentity_Update
    - TestSaveIdentity_ValidationError
    - TestDeleteIdentity_Exists
    - TestDeleteIdentity_NotExists
    - TestDeleteIdentity_InvalidIdk

  integration_tests:
    - TestPostgreSQLIntegration
    - TestMySQLIntegration
    - TestSQLiteIntegration
    - TestConcurrentAccess
    - TestTransactionRollback

  security_tests:
    - TestSQLInjectionPrevention
    - TestInputValidation_LengthLimits
    - TestInputValidation_CharacterSet
    - TestSensitiveDataNotInErrors
    - TestMemoryClearing

  performance_tests:
    - BenchmarkFindIdentity
    - BenchmarkSaveIdentity
    - BenchmarkDeleteIdentity
    - BenchmarkConcurrentReads
```

### API Test Template

See the separate API test file: `API_TESTS.md`

---

## Versioning

### API Versions

| Version | Status | Go Import Path | Breaking Changes |
|---------|--------|---------------|------------------|
| v0.1.0 | Released | github.com/sqrldev/... | N/A (initial) |
| v0.2.0 | Released | github.com/dxcSithLord/... | GORM v1 → v2, Go 1.24+ |
| v0.3.0-rc1 | Released | github.com/dxcSithLord/... | None (test + security) |
| v1.0.0 | Pending | github.com/sqrldev/... | Module path revert |

### Deprecation Policy

- Methods marked deprecated in v1.x will be removed in v2.0
- Minimum 6 months notice for deprecations
- Migration guide provided for breaking changes

---

## Appendix: OpenAPI 3.1 Mapping

### Go Concepts → OpenAPI Mappings

| Go Concept | OpenAPI 3.1 Equivalent |
|------------|----------------------|
| Interface | `components.interfaces` (custom) |
| Method | `paths.{operation}` |
| Struct | `components.schemas` |
| Error | `responses.{code}.content` |
| Constructor | `POST /resource/new` |
| Property | `schema.properties` |

---

**Document Control:**
- Version: 2.0.0
- Last Updated: 2026-02-08
- Next Review: Before v1.0.0 release

**END OF API SPECIFICATION**
