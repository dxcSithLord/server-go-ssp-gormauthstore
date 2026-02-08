# Enterprise Architecture Documentation
## SQRL GORM Authentication Store
## TOGAF-Aligned Architecture Views

**Document Version:** 2.0
**Date:** November 18, 2025 (updated February 8, 2026)
**Architecture Framework:** TOGAF 9.2
**Modeling Language:** ArchiMate 3.1 (represented via Mermaid)

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Business Architecture](#business-architecture)
3. [Application Architecture](#application-architecture)
4. [Data Architecture](#data-architecture)
5. [Technology Architecture](#technology-architecture)
6. [Objectives and Requirements Mapping](#objectives-and-requirements-mapping)
7. [Component Interaction Diagrams](#component-interaction-diagrams)
8. [Deployment Architecture](#deployment-architecture)

---

## Architecture Overview

This document provides TOGAF-aligned architectural views of the SQRL GORM Authentication Store, showing the relationship between business objectives, functional requirements, and technical implementation.

---

## Business Architecture

### Business Capability Model

```mermaid
graph TB
    subgraph "Business Capabilities"
        BC1[Passwordless Authentication]
        BC2[Identity Management]
        BC3[Privacy-Preserving Auth]
        BC4[Cross-Platform Support]
    end

    subgraph "Supporting Capabilities"
        SC1[Data Persistence]
        SC2[Schema Management]
        SC3[Identity Lifecycle]
        SC4[Security Controls]
    end

    BC1 --> SC1
    BC1 --> SC4
    BC2 --> SC3
    BC2 --> SC1
    BC3 --> SC4
    BC4 --> SC2

    subgraph "Enabling Services"
        ES1[gormauthstore Library]
        ES2[GORM ORM]
        ES3[Database Drivers]
    end

    SC1 --> ES1
    SC2 --> ES1
    SC3 --> ES1
    SC4 --> ES1
    ES1 --> ES2
    ES2 --> ES3

    style BC1 fill:#e1f5ff
    style BC2 fill:#e1f5ff
    style BC3 fill:#e1f5ff
    style BC4 fill:#e1f5ff
    style ES1 fill:#fff9c4
```

### Business Objectives to Functional Requirements

```mermaid
graph LR
    subgraph "Business Objectives"
        BO1[BO-001: Enable SQRL<br/>Authentication in Go]
        BO2[BO-002: Database<br/>Agnostic Storage]
        BO3[BO-003: Cryptographic<br/>Key Security]
        BO4[BO-004: Simple<br/>Integration]
        BO5[BO-005: Production<br/>Grade Deployment]
    end

    subgraph "Functional Requirements"
        FR1[FR-001: Identity Storage]
        FR2[FR-002: Identity Retrieval]
        FR3[FR-003: Identity Deletion]
        FR4[FR-004: Schema Management]
        FR5[FR-005: Interface Compliance]
    end

    subgraph "Non-Functional Requirements"
        NFR1[NFR-001: Performance]
        NFR2[NFR-002: Reliability]
        NFR3[NFR-003: Security]
        NFR4[NFR-004: Maintainability]
        NFR5[NFR-005: Compatibility]
    end

    BO1 --> FR1
    BO1 --> FR2
    BO1 --> FR3
    BO1 --> FR4
    BO1 --> FR5

    BO2 --> FR4
    BO2 --> FR5
    BO2 --> NFR5

    BO3 --> FR2
    BO3 --> FR3
    BO3 --> NFR3

    BO4 --> FR5

    BO5 --> NFR1
    BO5 --> NFR2
    BO5 --> NFR3

    style BO1 fill:#ffccbc
    style BO2 fill:#ffccbc
    style BO3 fill:#ffccbc
    style BO4 fill:#ffccbc
    style BO5 fill:#ffccbc

    style FR1 fill:#c8e6c9
    style FR2 fill:#c8e6c9
    style FR3 fill:#c8e6c9
    style FR4 fill:#c8e6c9
    style FR5 fill:#c8e6c9

    style NFR1 fill:#fff9c4
    style NFR2 fill:#fff9c4
    style NFR3 fill:#fff9c4
    style NFR4 fill:#fff9c4
    style NFR5 fill:#fff9c4
```

---

## Application Architecture

### Component Architecture

```mermaid
graph TB
    subgraph "Application Layer: gormauthstore"
        AS[AuthStore]
        SM[Secure Memory Manager]
        VAL[Input Validator]
        ERR[Error Handler]
    end

    subgraph "Integration Layer"
        INT[ssp.AuthStore Interface]
    end

    subgraph "Data Access Layer"
        GORM[GORM ORM v2]
    end

    subgraph "Database Drivers"
        PG[PostgreSQL Driver]
        MY[MySQL Driver]
        SQ[SQLite Driver]
        MS[MSSQL Driver]
    end

    subgraph "Database Layer"
        DB[(Database)]
    end

    INT -.implements.- AS
    AS --> SM
    AS --> VAL
    AS --> ERR
    AS --> GORM
    SM -.used by.- AS

    GORM --> PG
    GORM --> MY
    GORM --> SQ
    GORM --> MS

    PG --> DB
    MY --> DB
    SQ --> DB
    MS --> DB

    style AS fill:#4CAF50
    style INT fill:#2196F3
    style GORM fill:#FF9800
    style DB fill:#9C27B0
```

### Application Component Details

```mermaid
graph LR
    subgraph "AuthStore Component"
        direction TB
        AS_NEW[NewAuthStore]
        AS_FIND[FindIdentity<br/>FindIdentityWithContext]
        AS_SAVE[SaveIdentity<br/>SaveIdentityWithContext]
        AS_DEL[DeleteIdentity<br/>DeleteIdentityWithContext]
        AS_MIG[AutoMigrate<br/>AutoMigrateWithContext]
        AS_SEC[FindIdentitySecure<br/>FindIdentitySecureWithContext]
    end

    subgraph "Secure Memory Component"
        direction TB
        SM_WIPE[WipeBytes]
        SM_CLEAR[ClearIdentity]
        SM_WRAP[SecureIdentityWrapper]
        SM_SCRAM[ScrambleBytes]
    end

    subgraph "Validation Component"
        direction TB
        VAL_IDK[ValidateIdk]
        VAL_CHAR[isValidIdkChar]
    end

    subgraph "Error Component"
        direction TB
        ERR_EMPTY[ErrEmptyIdentityKey]
        ERR_LONG[ErrIdentityKeyTooLong]
        ERR_FMT[ErrInvalidIdentityKeyFormat]
    end

    AS_FIND --> VAL_IDK
    AS_DEL --> VAL_IDK
    AS_SAVE --> VAL_IDK

    AS_FIND -.returns.- SM_WRAP
    SM_WRAP --> SM_CLEAR
    SM_CLEAR --> SM_WIPE
    SM_WIPE --> SM_SCRAM

    VAL_IDK --> VAL_CHAR
    VAL_IDK --> ERR_EMPTY
    VAL_IDK --> ERR_LONG
    VAL_IDK --> ERR_FMT

    style AS_FIND fill:#4CAF50
    style AS_SAVE fill:#4CAF50
    style AS_DEL fill:#4CAF50
```

---

## Data Architecture

### Conceptual Data Model

```mermaid
erDiagram
    SQRL_IDENTITY {
        string idk PK "Identity Key (43-44 chars)"
        string suk "Server Unlock Key"
        string vuk "Verify Unlock Key"
        string pidk "Previous Identity Key"
        bool sqrl_only "SQRL-only flag"
        bool hardlock "Hard lock status"
        bool disabled "Disabled flag"
        string rekeyed FK "New identity link"
        int btn "Button response (0-3)"
    }

    SQRL_IDENTITY ||--o| SQRL_IDENTITY : "rekeyed to"
```

### Data Sensitivity Classification

```mermaid
graph TB
    subgraph "Data Classification"
        direction TB
        CRIT[CRITICAL: Cryptographic Keys]
        SENS[SENSITIVE: Identity Keys]
        CONF[CONFIDENTIAL: Status Flags]
        PUB[PUBLIC: None]
    end

    subgraph "Fields by Classification"
        direction TB
        F_CRIT[Suk, Vuk]
        F_SENS[Idk, Pidk, Rekeyed]
        F_CONF["`SQRLOnly, Hardlock,
Disabled, Btn`"]
    end

    subgraph "Security Controls"
        direction TB
        SC_MEM[Memory Clearing]
        SC_TLS[TLS Encryption]
        SC_LOG[No Logging]
        SC_ERR[Error Sanitization]
        SC_ENC["DB Encryption (TDE)"]
    end

    CRIT --> F_CRIT
    SENS --> F_SENS
    CONF --> F_CONF

    F_CRIT --> SC_MEM
    F_CRIT --> SC_TLS
    F_CRIT --> SC_LOG
    F_CRIT --> SC_ERR
    F_CRIT --> SC_ENC

    F_SENS --> SC_TLS
    F_SENS --> SC_ENC

    style CRIT fill:#f44336
    style SENS fill:#ff9800
    style CONF fill:#ffc107
    style PUB fill:#4caf50
```

### Data Flow Diagram

```mermaid
flowchart TB
    START([Client Application])

    START -->|SaveIdentity| SAVE{SaveIdentity}
    START -->|FindIdentity| FIND{FindIdentity}
    START -->|DeleteIdentity| DEL{DeleteIdentity}

    SAVE -->|Validate| VAL1[ValidateIdk]
    FIND -->|Validate| VAL2[ValidateIdk]
    DEL -->|Validate| VAL3[ValidateIdk]

    VAL1 -->|Valid| GORM_S[GORM.Save]
    VAL2 -->|Valid| GORM_F[GORM.Where.First]
    VAL3 -->|Valid| GORM_D[GORM.Where.Delete]

    VAL1 -->|Invalid| ERR1[Return Error]
    VAL2 -->|Invalid| ERR2[Return Error]
    VAL3 -->|Invalid| ERR3[Return Error]

    GORM_S -->|SQL| DB[(Database)]
    GORM_F -->|SQL| DB
    GORM_D -->|SQL| DB

    DB -->|Record| GORM_F
    DB -->|Success| GORM_S
    DB -->|Success| GORM_D

    GORM_F -->|Identity| WRAP[SecureIdentityWrapper]
    WRAP -->|Protected| RET_FIND([Return to Client])

    GORM_S -->|Success| RET_SAVE([Return to Client])
    GORM_D -->|Success| RET_DEL([Return to Client])

    WRAP -.Cleanup.- CLEAR[ClearIdentity]
    CLEAR -.Wipe.- WIPE[WipeBytes]

    style DB fill:#9C27B0
    style SAVE fill:#4CAF50
    style FIND fill:#4CAF50
    style DEL fill:#4CAF50
    style CLEAR fill:#f44336
    style WIPE fill:#f44336
```

---

## Technology Architecture

### Technology Stack

```mermaid
graph TB
    subgraph "Runtime Layer"
        GO[Go Runtime 1.25+]
    end

    subgraph "Application Framework"
        LIB[gormauthstore Library]
        SSP[server-go-ssp]
    end

    subgraph "ORM Layer"
        GORM2[GORM v2<br/>gorm.io/gorm v1.31.1]
    end

    subgraph "Database Drivers"
        DRV_PG[gorm.io/driver/postgres<br/>v1.5.9]
        DRV_MY[gorm.io/driver/mysql<br/>v1.5.7]
        DRV_SQ[gorm.io/driver/sqlite<br/>v1.6.0]
        DRV_MS[gorm.io/driver/sqlserver<br/>v1.5.3]
    end

    subgraph "Native Drivers"
        PQ[github.com/lib/pq<br/>v1.10.9]
        MY[github.com/go-sql-driver/mysql<br/>v1.9.3]
        SQ3[github.com/mattn/go-sqlite3<br/>v1.14.33]
        MSSQL[github.com/denisenkom/go-mssqldb<br/>v0.12.3]
    end

    subgraph "Database Management Systems"
        PG_DB[(PostgreSQL 12+)]
        MY_DB[(MySQL 8+)]
        SQ_DB[(SQLite 3.35+)]
        MS_DB[(SQL Server 2019+)]
    end

    GO --> LIB
    GO --> SSP
    SSP --> LIB

    LIB --> GORM2

    GORM2 --> DRV_PG
    GORM2 --> DRV_MY
    GORM2 --> DRV_SQ
    GORM2 --> DRV_MS

    DRV_PG --> PQ
    DRV_MY --> MY
    DRV_SQ --> SQ3
    DRV_MS --> MSSQL

    PQ --> PG_DB
    MY --> MY_DB
    SQ3 --> SQ_DB
    MSSQL --> MS_DB

    style GO fill:#00ADD8
    style LIB fill:#4CAF50
    style GORM2 fill:#FF9800
```

### Dependency Upgrade Path

```mermaid
graph LR
    subgraph "Current State (DEPRECATED)"
        OLD_GORM[jinzhu/gorm<br/>v1.9.16<br/>⚠️ DEPRECATED]
        OLD_PQ[lib/pq<br/>v1.1.1<br/>🔴 CRITICAL]
        OLD_SQ[go-sqlite3<br/>v1.14.22<br/>🟡 OUTDATED]
    end

    subgraph "Target State (CURRENT)"
        NEW_GORM[gorm.io/gorm<br/>v1.31.1<br/>✅ LATEST]
        NEW_PQ[lib/pq<br/>v1.10.9<br/>✅ LATEST]
        NEW_SQ[go-sqlite3<br/>v1.14.33<br/>✅ LATEST]
    end

    OLD_GORM -.BREAKING<br/>CHANGE.-> NEW_GORM
    OLD_PQ -.SECURITY<br/>FIXES.-> NEW_PQ
    OLD_SQ -.BUG<br/>FIXES.-> NEW_SQ

    style OLD_GORM fill:#ffcdd2
    style OLD_PQ fill:#ffcdd2
    style OLD_SQ fill:#fff9c4

    style NEW_GORM fill:#c8e6c9
    style NEW_PQ fill:#c8e6c9
    style NEW_SQ fill:#c8e6c9
```

---

## Objectives and Requirements Mapping

### High-Level Objectives Traceability

```mermaid
graph TB
    subgraph "Strategic Objectives"
        SO1[Secure Passwordless<br/>Authentication]
        SO2[Privacy-Preserving<br/>Identity]
        SO3[Developer<br/>Productivity]
        SO4[Enterprise<br/>Readiness]
    end

    subgraph "Business Objectives"
        BO1[Enable SQRL<br/>in Go Apps]
        BO2[Database<br/>Agnostic]
        BO3[Key<br/>Security]
        BO4[Simple<br/>Integration]
        BO5[Production<br/>Grade]
    end

    subgraph "Functional Requirements"
        FR1[Identity<br/>Storage]
        FR2[Identity<br/>Retrieval]
        FR3[Identity<br/>Deletion]
        FR4[Schema<br/>Management]
        FR5[Interface<br/>Compliance]
    end

    subgraph "Implementation"
        I1[SaveIdentity]
        I2[FindIdentity]
        I3[DeleteIdentity]
        I4[AutoMigrate]
        I5[AuthStore<br/>Type]
    end

    subgraph "Security Controls"
        SEC1[Memory<br/>Clearing]
        SEC2[SQL Injection<br/>Prevention]
        SEC3[Input<br/>Validation]
        SEC4[Error<br/>Sanitization]
    end

    SO1 --> BO1
    SO1 --> BO3
    SO2 --> BO3
    SO3 --> BO4
    SO4 --> BO5
    SO4 --> BO2

    BO1 --> FR1
    BO1 --> FR2
    BO1 --> FR3
    BO2 --> FR4
    BO3 --> SEC1
    BO3 --> SEC2
    BO4 --> FR5
    BO5 --> SEC3

    FR1 --> I1
    FR2 --> I2
    FR3 --> I3
    FR4 --> I4
    FR5 --> I5

    SEC1 --> I2
    SEC2 --> I1
    SEC2 --> I2
    SEC2 --> I3
    SEC3 --> I1
    SEC3 --> I2
    SEC3 --> I3
    SEC4 --> I1
    SEC4 --> I2
    SEC4 --> I3

    style SO1 fill:#e1bee7
    style SO2 fill:#e1bee7
    style SO3 fill:#e1bee7
    style SO4 fill:#e1bee7

    style BO1 fill:#ffccbc
    style BO3 fill:#ffccbc

    style SEC1 fill:#f44336,color:#fff
    style SEC2 fill:#f44336,color:#fff
    style SEC3 fill:#f44336,color:#fff
    style SEC4 fill:#f44336,color:#fff
```

---

## Component Interaction Diagrams

### Authentication Flow Sequence

```mermaid
sequenceDiagram
    participant App as Application
    participant SSP as SQRL SSP API
    participant AS as AuthStore
    participant GORM as GORM ORM
    participant DB as Database

    Note over App,DB: Authentication Request Flow

    App->>SSP: HandleCliRequest(req)
    activate SSP

    SSP->>AS: FindIdentity(idk)
    activate AS

    AS->>AS: ValidateIdk(idk)

    AS->>GORM: Where("idk = ?").First()
    activate GORM

    GORM->>DB: SELECT * FROM sqrl_identities<br/>WHERE idk = $1
    activate DB

    DB-->>GORM: SqrlIdentity record
    deactivate DB

    GORM-->>AS: *SqrlIdentity
    deactivate GORM

    AS->>AS: NewSecureIdentityWrapper(identity)

    AS-->>SSP: *SecureIdentityWrapper
    deactivate AS

    SSP->>SSP: Verify cryptographic signature

    SSP-->>App: Authentication result
    deactivate SSP

    Note over App,DB: Cleanup Phase

    App->>AS: wrapper.Destroy()
    activate AS

    AS->>AS: ClearIdentity()
    AS->>AS: WipeBytes(sensitive fields)

    deactivate AS
```

### Identity Lifecycle State Machine

```mermaid
stateDiagram-v2
    [*] --> Created: SaveIdentity(new)

    Created --> Active: Authentication successful
    Created --> Disabled: Admin action
    Created --> Deleted: DeleteIdentity

    Active --> Disabled: User disables
    Active --> Hardlocked: User enables hardlock
    Active --> Rekeyed: User rekeys identity
    Active --> Deleted: User removes account

    Hardlocked --> Active: User disables hardlock
    Hardlocked --> Deleted: User removes account

    Disabled --> Active: User re-enables
    Disabled --> Deleted: User removes account

    Rekeyed --> Deleted: Old identity cleanup

    Deleted --> [*]

    note right of Created
        Idk, Suk, Vuk stored
        SQRLOnly = false
        Hardlock = false
        Disabled = false
    end note

    note right of Rekeyed
        Rekeyed field points to
        new identity's Idk
    end note
```

### Error Handling Flow

```mermaid
flowchart TD
    START([API Method Called])

    START --> VAL{Input Validation}

    VAL -->|Empty Idk| ERR_EMPTY[Return ErrEmptyIdentityKey]
    VAL -->|Too Long| ERR_LONG[Return ErrIdentityKeyTooLong]
    VAL -->|Invalid Chars| ERR_FMT[Return ErrInvalidIdentityKeyFormat]
    VAL -->|Valid| DB_OP[Execute Database Operation]

    DB_OP --> DB_RESULT{Database Result}

    DB_RESULT -->|Not Found| MAP_ERR[Map to ssp.ErrNotFound]
    DB_RESULT -->|Connection Error| PROP_ERR[Propagate DB Error]
    DB_RESULT -->|Success| SUCCESS[Return Success]

    ERR_EMPTY --> LOG1[Log Error Context]
    ERR_LONG --> LOG2[Log Error Context]
    ERR_FMT --> LOG3[Log Error Context]
    MAP_ERR --> LOG4[Log Error Context]
    PROP_ERR --> LOG5[Log Error Context]

    LOG1 --> SANITIZE1[Sanitize Error Message]
    LOG2 --> SANITIZE2[Sanitize Error Message]
    LOG3 --> SANITIZE3[Sanitize Error Message]
    LOG4 --> SANITIZE4[Sanitize Error Message]
    LOG5 --> SANITIZE5[Sanitize Error Message]

    SANITIZE1 --> RET_ERR1([Return Error])
    SANITIZE2 --> RET_ERR2([Return Error])
    SANITIZE3 --> RET_ERR3([Return Error])
    SANITIZE4 --> RET_ERR4([Return Error])
    SANITIZE5 --> RET_ERR5([Return Error])

    SUCCESS --> RET_OK([Return Success])

    style ERR_EMPTY fill:#ffcdd2
    style ERR_LONG fill:#ffcdd2
    style ERR_FMT fill:#ffcdd2
    style MAP_ERR fill:#ffcdd2
    style PROP_ERR fill:#ffcdd2
    style SUCCESS fill:#c8e6c9
```

---

## Deployment Architecture

### Deployment Options

```mermaid
graph TB
    subgraph "Application Tier"
        APP1[App Instance 1]
        APP2[App Instance 2]
        APP3[App Instance N]
    end

    subgraph "gormauthstore Library"
        LIB1[AuthStore 1]
        LIB2[AuthStore 2]
        LIB3[AuthStore N]
    end

    subgraph "Connection Pool"
        POOL[GORM Connection Pool<br/>Max: Configurable<br/>Idle: Configurable]
    end

    subgraph "Database Cluster"
        DB_PRIMARY[(Primary DB)]
        DB_REPLICA1[(Replica 1)]
        DB_REPLICA2[(Replica 2)]
    end

    APP1 --> LIB1
    APP2 --> LIB2
    APP3 --> LIB3

    LIB1 --> POOL
    LIB2 --> POOL
    LIB3 --> POOL

    POOL -->|Writes| DB_PRIMARY
    POOL -->|Reads| DB_REPLICA1
    POOL -->|Reads| DB_REPLICA2

    DB_PRIMARY -.Replication.-> DB_REPLICA1
    DB_PRIMARY -.Replication.-> DB_REPLICA2

    style DB_PRIMARY fill:#4CAF50
    style DB_REPLICA1 fill:#8BC34A
    style DB_REPLICA2 fill:#8BC34A
```

### Network Security Architecture

```mermaid
graph TB
    subgraph "DMZ - Application Tier"
        APP[Application Server<br/>TLS 1.3]
    end

    subgraph "Internal Network - Database Tier"
        DB[(Database<br/>TLS 1.3<br/>TDE Enabled)]
    end

    subgraph "Security Controls"
        FW1[Firewall<br/>Port 5432/3306 Only]
        TLS[TLS Encryption<br/>In Transit]
        TDE[Transparent Data<br/>Encryption at Rest]
        AUTH[Certificate-Based<br/>Authentication]
    end

    APP -->|Encrypted Connection| FW1
    FW1 -->|Allowed| TLS
    TLS --> DB
    DB --> TDE
    TLS --> AUTH

    style FW1 fill:#f44336,color:#fff
    style TLS fill:#f44336,color:#fff
    style TDE fill:#f44336,color:#fff
    style AUTH fill:#f44336,color:#fff
```

### Multi-Database Deployment

```mermaid
graph LR
    subgraph "Application Layer"
        APP[Application]
        AS[AuthStore Factory]
    end

    subgraph "Database Connections"
        CONN_PG[PostgreSQL Connection]
        CONN_MY[MySQL Connection]
        CONN_SQ[SQLite Connection]
    end

    subgraph "Databases"
        PG[(PostgreSQL<br/>Production)]
        MY[(MySQL<br/>Staging)]
        SQ[(SQLite<br/>Dev/Test)]
    end

    APP --> AS
    AS -.Environment: PROD.-> CONN_PG
    AS -.Environment: STAGING.-> CONN_MY
    AS -.Environment: DEV.-> CONN_SQ

    CONN_PG --> PG
    CONN_MY --> MY
    CONN_SQ --> SQ

    style PG fill:#336791
    style MY fill:#4479A1
    style SQ fill:#003B57
```

---

## Architecture Decision Records

### ADR-001: Use GORM v2 ORM

**Status:** Accepted

**Context:** Need database abstraction layer supporting multiple databases

**Decision:** Use GORM v2 (gorm.io/gorm)

**Consequences:**
- ✅ Multi-database support
- ✅ Active maintenance
- ✅ Automatic migrations
- ❌ Learning curve for GORM-specific patterns

---

### ADR-002: Implement Secure Memory Clearing

**Status:** Accepted

**Context:** SQRL keys are sensitive cryptographic material

**Decision:** Implement platform-aware memory wiping (WipeBytes)

**Consequences:**
- ✅ Defense-in-depth security
- ✅ Compliance with CWE-226 mitigation
- ❌ Cannot guarantee complete clearing (Go string immutability)
- ℹ️ Documented limitations for users

---

### ADR-003: Input Validation Before Database Operations

**Status:** Accepted

**Context:** Need to prevent invalid data and DoS attacks

**Decision:** Validate all inputs (ValidateIdk) before database calls

**Consequences:**
- ✅ Early error detection
- ✅ DoS prevention (length limits)
- ✅ Clear error messages
- ❌ Slight performance overhead

---

## Appendix: Notation Guide

### Diagram Symbols

| Symbol | Meaning |
|--------|---------|
| Rectangle | Component/Service |
| Cylinder | Database |
| Diamond | Decision Point |
| Arrow | Data Flow |
| Dashed Arrow | Dependency |
| Subgraph | Logical Grouping |

### Color Coding

| Color | Meaning |
|-------|---------|
| 🔴 Red (#f44336) | Security Control / Critical |
| 🟠 Orange (#ff9800) | Framework / Infrastructure |
| 🟢 Green (#4CAF50) | Application Component |
| 🔵 Blue (#2196F3) | Interface / Contract |
| 🟣 Purple (#9C27B0) | Data Store |
| 🟡 Yellow (#fff9c4) | Non-Functional Requirement |

---

**Document Control:**
- Version: 2.0
- Last Updated: 2026-02-08
- Next Review: Before v1.0.0 release

**END OF ARCHITECTURE DOCUMENTATION**
