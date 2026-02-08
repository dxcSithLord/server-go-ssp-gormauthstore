# Documentation Test Suite

## Overview

This document describes the comprehensive test suite for validating the project documentation. These tests ensure documentation integrity, consistency, and quality.

## Test File

**Location:** `/docs_test.go`

**Purpose:** Validate all project documentation files for completeness, consistency, and integrity

## Test Categories

### 1. Existence Tests

**Test:** `TestDocumentationExists`

Verifies that all critical documentation files exist:
- README.md
- CLAUDE.md
- docs/PROJECT_PLAN.md
- docs/REQUIREMENTS.md
- docs/ARCHITECTURE.md
- docs/API_SPECIFICATION.md
- docs/API_TESTS_SPEC.md
- docs/DEPENDENCIES.md
- docs/Notice_of_Decisions.md
- docs/archive/PHASE1_TASKS.md

**Test:** `TestArchiveDocumentationExists`

Verifies that archived documentation files exist:
- docs/archive/SECURITY_REVIEW_AND_UPGRADE_PLAN.md
- docs/archive/STAGED_UPGRADE_PLAN.md
- docs/archive/TODO.md
- docs/archive/UNIFIED_TODO.md

### 2. Link Integrity Tests

**Test:** `TestMarkdownLinkIntegrity`

Validates that all internal markdown links point to existing files:
- Parses all markdown documents
- Extracts markdown links: `[text](path)`
- Resolves relative paths
- Verifies target files exist
- Skips external links (http, https, mailto)
- Allows anchor-only links

### 3. Content Structure Tests

**Test:** `TestREADMEContent`

Verifies README.md has essential sections:
- TL;DR - Project Status
- Overview
- Documentation
- Development
- License
- Overall Progress table
- GORM Version status
- Test Coverage status

**Test:** `TestCLAUDEInstructions`

Verifies CLAUDE.md has key sections:
- Working Preferences
- Project Overview
- Build and Test Commands
- Project Structure
- Progress Tracking
- Current State
- Code Conventions
- Decision Points
- Build commands (make ci, make test, make lint, make security)

**Test:** `TestProjectPlanStructure`

Verifies PROJECT_PLAN.md structure:
- All 3 phases present (Phase 1, 2, 3)
- All 6 stages present (1.1, 1.2, 1.3, 2.1, 2.2, 3.1)
- At least 44 tasks documented
- Decision Points section present

**Test:** `TestArchitectureDocument`

Verifies ARCHITECTURE.md has required content:
- Business Architecture
- Application Architecture
- Data Architecture
- Technology Architecture
- Deployment Architecture
- At least 5 mermaid diagrams

**Test:** `TestAPISpecification`

Verifies API_SPECIFICATION.md structure:
- Interface Contract
- Data Models
- Operations
- Error Responses
- Security
- Examples
- All operations documented (NewAuthStore, AutoMigrate, FindIdentity, SaveIdentity, DeleteIdentity)
- SqrlIdentity data model
- SQL Injection discussion in security section

**Test:** `TestAPITestsSpec`

Verifies API_TESTS_SPEC.md completeness:
- Unit Tests section
- Integration Tests section
- Security Tests section
- Performance Tests section
- Test case identifiers (TC-, IT-, SEC-, PERF-)
- Coverage target (70%) mentioned

**Test:** `TestDependenciesDocument`

Verifies DEPENDENCIES.md content:
- Key dependencies documented (gorm.io/gorm, github.com/lib/pq, github.com/mattn/go-sqlite3, golang.org/x/crypto)
- GORM migration discussed
- Security considerations

**Test:** `TestNoticeOfDecisions`

Verifies Notice_of_Decisions.md content:
- Decision points documented
- Protocol compliance discussed
- Decision response form present

**Test:** `TestPhase1Tasks`

Verifies PHASE1_TASKS.md completeness:
- All 3 stages present (1.1, 1.2, 1.3)
- Task identifiers present (TASK-001, TASK-010, TASK-020)
- Completion criteria defined

**Test:** `TestRequirementsDocument`

Verifies REQUIREMENTS.md structure:
- Functional Requirements
- Non-Functional Requirements
- Security Requirements
- Interface Requirements
- Data Requirements
- Requirement identifiers (FR-001, NFR-001, SEC-001, INT-001)
- Traceability Matrix

### 4. Consistency Tests

**Test:** `TestDocumentationConsistency`

Verifies cross-document consistency:
- **Go Version:** Consistent version references across README.md and CLAUDE.md
- **Phase Count:** Consistent phase references (Phase 1, 2, 3) across documents
- **GORM Status:** Consistent GORM deprecation status

### 5. Security Tests

**Test:** `TestNoHardcodedPasswords`

Ensures no passwords or secrets in documentation:
- Checks for suspicious patterns (password=, passwd=, pwd=, secret=, token=, key=)
- Allows obvious placeholders and examples
- Logs warnings for manual review

**Test:** `TestArchiveDocumentationMarked`

Verifies archived docs have deprecation notices:
- Checks for markers: "superseded", "archived", "obsolete", "deprecated"
- Logs informational notes (non-failing)

### 6. Diagram Validation Tests

**Test:** `TestMermaidDiagramSyntax`

Performs basic validation of mermaid diagrams:
- Finds all mermaid code blocks
- Verifies diagram type is present (graph, sequenceDiagram, flowchart, stateDiagram, erDiagram)
- Checks for unclosed blocks
- Documents with diagrams:
  - docs/ARCHITECTURE.md
  - docs/Notice_of_Decisions.md
  - docs/API_TESTS_SPEC.md

## Running the Tests

### Run all documentation tests
```bash
go test -v ./...
```

### Run specific test
```bash
go test -v -run TestREADMEContent ./...
```

### Run with coverage
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run as part of CI
```bash
make test  # Runs all tests including documentation tests
```

## Test Coverage

The documentation test suite provides:

- **14 test functions** covering different aspects of documentation quality
- **100+ sub-tests** for individual documents and sections
- **Link integrity validation** for all markdown links
- **Content structure validation** for all major documents
- **Cross-document consistency checks**
- **Security validation** (no hardcoded secrets)
- **Diagram syntax validation** for mermaid diagrams

## Maintenance

### Adding New Documentation

When adding new documentation files:

1. Add the file path to `TestDocumentationExists` or `TestArchiveDocumentationExists`
2. If the document has required sections, add a new test function
3. If the document contains links, it will be automatically validated by `TestMarkdownLinkIntegrity`
4. Update this DOCUMENTATION_TESTS.md file

### Modifying Existing Documentation

When modifying documentation:

1. Ensure all links are valid (tests will catch broken links)
2. Maintain required sections (tests will catch missing sections)
3. Keep cross-document references consistent
4. Run `go test ./docs_test.go` to verify changes

## Benefits

This documentation test suite provides:

1. **Early Detection:** Catches broken links and missing sections during development
2. **Consistency:** Ensures all documents follow the same structure
3. **Quality:** Validates that documentation meets minimum standards
4. **Automation:** Can be integrated into CI/CD pipelines
5. **Regression Prevention:** Prevents documentation degradation over time
6. **Confidence:** Provides confidence that documentation is complete and accurate

## Future Enhancements

Potential future improvements:

1. Spell checking integration
2. Grammar checking
3. Markdown linting (beyond basic syntax)
4. External link validation (HTTP requests)
5. Documentation coverage metrics
6. Auto-generated documentation from code
7. Mermaid diagram rendering validation
8. Cross-reference validation between code and docs

## See Also

- [PROJECT_PLAN.md](PROJECT_PLAN.md) - Overall project plan
- [REQUIREMENTS.md](REQUIREMENTS.md) - Project requirements
- [API_TESTS_SPEC.md](API_TESTS_SPEC.md) - API test specifications
- [ARCHITECTURE.md](ARCHITECTURE.md) - Architecture documentation

---

**Document Control:**
- Version: 1.1
- Created: 2026-02-05
- Last Updated: 2026-02-08

**END OF DOCUMENTATION TEST SPECIFICATION**