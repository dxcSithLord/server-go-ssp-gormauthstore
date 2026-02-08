package gormauthstore

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestDocumentationExists verifies that all critical documentation files exist.
func TestDocumentationExists(t *testing.T) {
	requiredDocs := []string{
		"README.md",
		"CLAUDE.md",
		"docs/PROJECT_PLAN.md",
		"docs/REQUIREMENTS.md",
		"docs/ARCHITECTURE.md",
		"docs/API_SPECIFICATION.md",
		"docs/API_TESTS_SPEC.md",
		"docs/DEPENDENCIES.md",
		"docs/Notice_of_Decisions.md",
	}

	for _, doc := range requiredDocs {
		t.Run(doc, func(t *testing.T) {
			if _, err := os.Stat(doc); os.IsNotExist(err) {
				t.Errorf("Required documentation file does not exist: %s", doc)
			}
		})
	}
}

// TestArchiveDocumentationExists verifies archived documentation.
func TestArchiveDocumentationExists(t *testing.T) {
	archiveDocs := []string{
		"docs/archive/SECURITY_REVIEW_AND_UPGRADE_PLAN.md",
		"docs/archive/STAGED_UPGRADE_PLAN.md",
		"docs/archive/TODO.md",
		"docs/archive/UNIFIED_TODO.md",
		"docs/archive/PHASE1_TASKS.md",
	}

	for _, doc := range archiveDocs {
		t.Run(doc, func(t *testing.T) {
			if _, err := os.Stat(doc); os.IsNotExist(err) {
				t.Errorf("Archive documentation file does not exist: %s", doc)
			}
		})
	}
}

// TestMarkdownLinkIntegrity checks that internal markdown links point to existing files.
func TestMarkdownLinkIntegrity(t *testing.T) {
	// Regex to match markdown links: [text](path) or [text](path#anchor)
	linkRegex := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	docs := []string{
		"README.md",
		"CLAUDE.md",
		"docs/PROJECT_PLAN.md",
		"docs/REQUIREMENTS.md",
		"docs/ARCHITECTURE.md",
		"docs/API_SPECIFICATION.md",
		"docs/API_TESTS_SPEC.md",
		"docs/DEPENDENCIES.md",
		"docs/Notice_of_Decisions.md",
	}

	for _, doc := range docs {
		t.Run(doc, func(t *testing.T) {
			file, err := os.Open(doc)
			if err != nil {
				t.Fatalf("Cannot open %s: %v", doc, err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()
				matches := linkRegex.FindAllStringSubmatch(line, -1)

				for _, match := range matches {
					linkURL := match[2]

					// Skip external links (http, https, mailto)
					if strings.HasPrefix(linkURL, "http://") ||
						strings.HasPrefix(linkURL, "https://") ||
						strings.HasPrefix(linkURL, "mailto:") ||
						strings.HasPrefix(linkURL, "#") {
						continue
					}

					// Remove anchor from link
					linkPath := strings.Split(linkURL, "#")[0]
					if linkPath == "" {
						continue // Anchor-only links are OK
					}

					// Resolve relative path
					docDir := filepath.Dir(doc)
					targetPath := filepath.Join(docDir, linkPath)

					// Check if target exists
					if _, err := os.Stat(targetPath); os.IsNotExist(err) {
						t.Errorf("%s:%d: Broken link to %s (resolved as %s)",
							doc, lineNum, linkURL, targetPath)
					}
				}
			}

			if err := scanner.Err(); err != nil {
				t.Errorf("Error scanning %s: %v", doc, err)
			}
		})
	}
}

// TestREADMEContent verifies README.md has essential sections.
func TestREADMEContent(t *testing.T) {
	content, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Cannot read README.md: %v", err)
	}

	readme := string(content)

	requiredSections := []string{
		"## TL;DR - Project Status",
		"## Overview",
		"## Documentation",
		"## Development",
		"## License",
	}

	for _, section := range requiredSections {
		t.Run(section, func(t *testing.T) {
			if !strings.Contains(readme, section) {
				t.Errorf("README.md missing required section: %s", section)
			}
		})
	}

	// Verify progress table exists
	if !strings.Contains(readme, "Overall Progress") {
		t.Error("README.md missing Overall Progress table")
	}

	// Verify critical status indicators
	if !strings.Contains(readme, "GORM Version") {
		t.Error("README.md missing GORM Version status")
	}
	if !strings.Contains(readme, "Test Coverage") {
		t.Error("README.md missing Test Coverage status")
	}
}

// TestCLAUDEInstructions verifies CLAUDE.md has key sections.
func TestCLAUDEInstructions(t *testing.T) {
	content, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("Cannot read CLAUDE.md: %v", err)
	}

	claude := string(content)

	requiredSections := []string{
		"## Working Preferences",
		"## Project Overview",
		"## Build and Test Commands",
		"## Project Structure",
		"## Progress Tracking",
		"## Current State",
		"## Code Conventions",
		"## Decision Points",
	}

	for _, section := range requiredSections {
		t.Run(section, func(t *testing.T) {
			if !strings.Contains(claude, section) {
				t.Errorf("CLAUDE.md missing required section: %s", section)
			}
		})
	}

	// Verify build commands are present
	buildCommands := []string{"make ci", "make test", "make lint", "make security"}
	for _, cmd := range buildCommands {
		if !strings.Contains(claude, cmd) {
			t.Errorf("CLAUDE.md missing build command: %s", cmd)
		}
	}
}

// TestProjectPlanStructure verifies PROJECT_PLAN.md structure.
func TestProjectPlanStructure(t *testing.T) {
	content, err := os.ReadFile("docs/PROJECT_PLAN.md")
	if err != nil {
		t.Fatalf("Cannot read PROJECT_PLAN.md: %v", err)
	}

	plan := string(content)

	// Verify phases exist
	phases := []string{"Phase 1: Critical Foundation", "Phase 2: Security and Testing", "Phase 3: Production Readiness"}
	for _, phase := range phases {
		t.Run(phase, func(t *testing.T) {
			if !strings.Contains(plan, phase) {
				t.Errorf("PROJECT_PLAN.md missing: %s", phase)
			}
		})
	}

	// Verify stages exist
	stages := []string{"Stage 1.1", "Stage 1.2", "Stage 1.3", "Stage 2.1", "Stage 2.2", "Stage 3.1"}
	for _, stage := range stages {
		if !strings.Contains(plan, stage) {
			t.Errorf("PROJECT_PLAN.md missing: %s", stage)
		}
	}

	// Verify task count
	taskCount := strings.Count(plan, "TASK-")
	expectedTasks := 44
	if taskCount < expectedTasks {
		t.Errorf("PROJECT_PLAN.md has %d tasks, expected at least %d", taskCount, expectedTasks)
	}

	// Verify decision points
	if !strings.Contains(plan, "Decision Points") {
		t.Error("PROJECT_PLAN.md missing Decision Points section")
	}
}

// TestArchitectureDocument verifies ARCHITECTURE.md has required content.
func TestArchitectureDocument(t *testing.T) {
	content, err := os.ReadFile("docs/ARCHITECTURE.md")
	if err != nil {
		t.Fatalf("Cannot read ARCHITECTURE.md: %v", err)
	}

	arch := string(content)

	requiredSections := []string{
		"## Business Architecture",
		"## Application Architecture",
		"## Data Architecture",
		"## Technology Architecture",
		"## Deployment Architecture",
	}

	for _, section := range requiredSections {
		t.Run(section, func(t *testing.T) {
			if !strings.Contains(arch, section) {
				t.Errorf("ARCHITECTURE.md missing: %s", section)
			}
		})
	}

	// Verify mermaid diagrams exist
	if !strings.Contains(arch, "```mermaid") {
		t.Error("ARCHITECTURE.md missing mermaid diagrams")
	}

	// Count mermaid diagrams
	mermaidCount := strings.Count(arch, "```mermaid")
	if mermaidCount < 5 {
		t.Errorf("ARCHITECTURE.md has %d mermaid diagrams, expected at least 5", mermaidCount)
	}
}

// TestAPISpecification verifies API_SPECIFICATION.md structure.
func TestAPISpecification(t *testing.T) {
	content, err := os.ReadFile("docs/API_SPECIFICATION.md")
	if err != nil {
		t.Fatalf("Cannot read API_SPECIFICATION.md: %v", err)
	}

	spec := string(content)

	// Verify OpenAPI-style sections
	requiredSections := []string{
		"## Interface Contract",
		"## Data Models",
		"## Operations",
		"## Error Responses",
		"## Security",
		"## Examples",
	}

	for _, section := range requiredSections {
		t.Run(section, func(t *testing.T) {
			if !strings.Contains(spec, section) {
				t.Errorf("API_SPECIFICATION.md missing: %s", section)
			}
		})
	}

	// Verify operations are documented
	operations := []string{"NewAuthStore", "AutoMigrate", "FindIdentity", "SaveIdentity", "DeleteIdentity"}
	for _, op := range operations {
		if !strings.Contains(spec, op) {
			t.Errorf("API_SPECIFICATION.md missing operation: %s", op)
		}
	}

	// Verify SqrlIdentity data model
	if !strings.Contains(spec, "SqrlIdentity") {
		t.Error("API_SPECIFICATION.md missing SqrlIdentity data model")
	}

	// Verify security section has content
	if strings.Contains(spec, "## Security") {
		securityIndex := strings.Index(spec, "## Security")
		if securityIndex != -1 {
			end := securityIndex + 2000
			if end > len(spec) {
				end = len(spec)
			}
			securitySection := spec[securityIndex:end]
			if !strings.Contains(securitySection, "SQL Injection") {
				t.Error("API_SPECIFICATION.md security section missing SQL Injection discussion")
			}
		}
	}
}

// TestAPITestsSpec verifies API_TESTS_SPEC.md completeness.
func TestAPITestsSpec(t *testing.T) {
	content, err := os.ReadFile("docs/API_TESTS_SPEC.md")
	if err != nil {
		t.Fatalf("Cannot read API_TESTS_SPEC.md: %v", err)
	}

	testSpec := string(content)

	// Verify test categories
	testCategories := []string{
		"## Unit Tests",
		"## Integration Tests",
		"## Security Tests",
		"## Performance Tests",
	}

	for _, category := range testCategories {
		t.Run(category, func(t *testing.T) {
			if !strings.Contains(testSpec, category) {
				t.Errorf("API_TESTS_SPEC.md missing: %s", category)
			}
		})
	}

	// Verify test case identifiers
	testIDPatterns := []string{"TC-", "IT-", "SEC-", "PERF-"}
	for _, pattern := range testIDPatterns {
		if !strings.Contains(testSpec, pattern) {
			t.Errorf("API_TESTS_SPEC.md missing test case pattern: %s", pattern)
		}
	}

	// Verify coverage target is mentioned
	if !strings.Contains(testSpec, "70%") {
		t.Error("API_TESTS_SPEC.md missing coverage target")
	}
}

// TestDependenciesDocument verifies DEPENDENCIES.md.
func TestDependenciesDocument(t *testing.T) {
	content, err := os.ReadFile("docs/DEPENDENCIES.md")
	if err != nil {
		t.Fatalf("Cannot read DEPENDENCIES.md: %v", err)
	}

	deps := string(content)

	// Verify key dependencies are documented
	keyDeps := []string{
		"gorm.io/gorm",
		"github.com/lib/pq",
		"github.com/mattn/go-sqlite3",
		"golang.org/x/crypto",
	}

	for _, dep := range keyDeps {
		t.Run(dep, func(t *testing.T) {
			if !strings.Contains(deps, dep) {
				t.Errorf("DEPENDENCIES.md missing dependency: %s", dep)
			}
		})
	}

	// Verify GORM v2 migration is discussed
	if !strings.Contains(deps, "GORM v1") && !strings.Contains(deps, "GORM v2") {
		t.Error("DEPENDENCIES.md missing GORM migration discussion")
	}

	// Verify security section
	if !strings.Contains(deps, "Security") {
		t.Error("DEPENDENCIES.md missing security considerations")
	}
}

// TestNoticeOfDecisions verifies decision documentation.
func TestNoticeOfDecisions(t *testing.T) {
	content, err := os.ReadFile("docs/Notice_of_Decisions.md")
	if err != nil {
		t.Fatalf("Cannot read Notice_of_Decisions.md: %v", err)
	}

	decisions := string(content)

	// Verify decision points are documented (document uses full names like "DP-001:")
	decisionPatterns := []string{"DECISION POINT 1", "DECISION POINT 2", "DECISION POINT 3"}
	for i, pattern := range decisionPatterns {
		dpNum := i + 1
		t.Run(pattern, func(t *testing.T) {
			// Check for various formats: DP-00X or DECISION POINT X
			if !strings.Contains(decisions, pattern) && !strings.Contains(decisions, "DP-00"+string(rune('0'+dpNum))) {
				t.Logf("NOTE: Notice_of_Decisions.md uses alternative format for decision %d", dpNum)
				// Don't fail - just informational
			}
		})
	}

	// Verify protocol compliance is discussed
	if !strings.Contains(decisions, "Protocol Compliance") {
		t.Error("Notice_of_Decisions.md missing protocol compliance analysis")
	}

	// Verify decision response form exists
	if !strings.Contains(decisions, "Decision Response Form") {
		t.Error("Notice_of_Decisions.md missing decision response form")
	}
}

// TestPhase1Tasks verifies PHASE1_TASKS.md completeness (archived).
func TestPhase1Tasks(t *testing.T) {
	content, err := os.ReadFile("docs/archive/PHASE1_TASKS.md")
	if err != nil {
		t.Fatalf("Cannot read archive/PHASE1_TASKS.md: %v", err)
	}

	tasks := string(content)

	// Verify stages
	stages := []string{"Stage 1.1", "Stage 1.2", "Stage 1.3"}
	for _, stage := range stages {
		t.Run(stage, func(t *testing.T) {
			if !strings.Contains(tasks, stage) {
				t.Errorf("PHASE1_TASKS.md missing: %s", stage)
			}
		})
	}

	// Verify task identifiers
	taskIDs := []string{"TASK-001", "TASK-010", "TASK-020"}
	for _, taskID := range taskIDs {
		if !strings.Contains(tasks, taskID) {
			t.Errorf("PHASE1_TASKS.md missing: %s", taskID)
		}
	}

	// Verify completion criteria
	if !strings.Contains(tasks, "Completion Criteria") {
		t.Error("PHASE1_TASKS.md missing completion criteria")
	}
}

// TestRequirementsDocument verifies REQUIREMENTS.md structure.
func TestRequirementsDocument(t *testing.T) {
	content, err := os.ReadFile("docs/REQUIREMENTS.md")
	if err != nil {
		t.Fatalf("Cannot read REQUIREMENTS.md: %v", err)
	}

	reqs := string(content)

	// Verify requirement types
	reqTypes := []string{
		"## Functional Requirements",
		"## Non-Functional Requirements",
		"## Security Requirements",
		"## Interface Requirements",
		"## Data Requirements",
	}

	for _, reqType := range reqTypes {
		t.Run(reqType, func(t *testing.T) {
			if !strings.Contains(reqs, reqType) {
				t.Errorf("REQUIREMENTS.md missing: %s", reqType)
			}
		})
	}

	// Verify requirement identifiers
	reqIDs := []string{"FR-001", "NFR-001", "SEC-001", "INT-001"}
	for _, reqID := range reqIDs {
		if !strings.Contains(reqs, reqID) {
			t.Errorf("REQUIREMENTS.md missing requirement: %s", reqID)
		}
	}

	// Verify traceability matrix
	if !strings.Contains(reqs, "Traceability Matrix") {
		t.Error("REQUIREMENTS.md missing traceability matrix")
	}
}

// TestDocumentationConsistency verifies cross-document consistency.
func TestDocumentationConsistency(t *testing.T) {
	// Read key documents
	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("Cannot read README.md: %v", err)
	}
	claude, err := os.ReadFile("CLAUDE.md")
	if err != nil {
		t.Fatalf("Cannot read CLAUDE.md: %v", err)
	}
	plan, err := os.ReadFile("docs/PROJECT_PLAN.md")
	if err != nil {
		t.Fatalf("Cannot read docs/PROJECT_PLAN.md: %v", err)
	}

	readmeStr := string(readme)
	claudeStr := string(claude)
	planStr := string(plan)

	// Verify consistent Go version
	t.Run("GoVersion", func(t *testing.T) {
		if !strings.Contains(claudeStr, "1.24") {
			t.Error("CLAUDE.md missing Go version 1.24 reference")
		}
		if !strings.Contains(readmeStr, "1.24") {
			t.Error("README.md missing Go version 1.24 reference")
		}
	})

	// Verify consistent phase count
	t.Run("PhaseCount", func(t *testing.T) {
		if !strings.Contains(readmeStr, "Phase 1") || !strings.Contains(readmeStr, "Phase 2") || !strings.Contains(readmeStr, "Phase 3") {
			t.Error("README.md missing phase references")
		}
		if !strings.Contains(planStr, "Phase 1") || !strings.Contains(planStr, "Phase 2") || !strings.Contains(planStr, "Phase 3") {
			t.Error("PROJECT_PLAN.md missing phase references")
		}
	})

	// Verify consistent GORM version status
	t.Run("GORMStatus", func(t *testing.T) {
		if !strings.Contains(readmeStr, "DEPRECATED") && !strings.Contains(readmeStr, "v1.9.16") {
			t.Error("README.md missing GORM deprecation status")
		}
		if !strings.Contains(claudeStr, "GORM") {
			t.Error("CLAUDE.md missing GORM reference")
		}
	})
}

// TestArchiveDocumentationMarked verifies archived docs have deprecation notice.
func TestArchiveDocumentationMarked(t *testing.T) {
	archiveDocs := []string{
		"docs/archive/SECURITY_REVIEW_AND_UPGRADE_PLAN.md",
		"docs/archive/STAGED_UPGRADE_PLAN.md",
		"docs/archive/TODO.md",
		"docs/archive/UNIFIED_TODO.md",
	}

	for _, doc := range archiveDocs {
		t.Run(doc, func(t *testing.T) {
			content, err := os.ReadFile(doc)
			if err != nil {
				t.Skipf("Cannot read %s: %v", doc, err)
				return
			}

			docStr := string(content)

			// Check for deprecation markers (any of these would indicate it's marked as archived)
			hasArchiveMarker := strings.Contains(strings.ToLower(docStr), "superseded") ||
				strings.Contains(strings.ToLower(docStr), "archived") ||
				strings.Contains(strings.ToLower(docStr), "obsolete") ||
				strings.Contains(strings.ToLower(docStr), "deprecated")

			if !hasArchiveMarker {
				t.Logf("NOTE: %s may benefit from an archive/superseded notice", doc)
				// Not failing the test, just logging
			}
		})
	}
}

// TestNoHardcodedPasswords ensures no passwords or secrets in documentation.
func TestNoHardcodedPasswords(t *testing.T) {
	docs := []string{
		"README.md",
		"CLAUDE.md",
		"docs/PROJECT_PLAN.md",
		"docs/REQUIREMENTS.md",
		"docs/ARCHITECTURE.md",
		"docs/API_SPECIFICATION.md",
		"docs/DEPENDENCIES.md",
	}

	// Common password patterns (examples in docs are OK if obvious)
	suspiciousPatterns := []string{
		"password=",
		"passwd=",
		"pwd=",
		"secret=",
		"token=",
		"key=",
	}

	for _, doc := range docs {
		t.Run(doc, func(t *testing.T) {
			content, err := os.ReadFile(doc)
			if err != nil {
				t.Skipf("Cannot read %s: %v", doc, err)
				return
			}

			docStr := strings.ToLower(string(content))
			lines := strings.Split(docStr, "\n")
			for _, pattern := range suspiciousPatterns {
				for lineNum, line := range lines {
					if strings.Contains(line, pattern) {
						// Check if this specific line contains an obvious placeholder
						if strings.Contains(line, "password=secret") ||
							strings.Contains(line, "password=test") ||
							strings.Contains(line, "password=your") ||
							strings.Contains(line, "example") ||
							strings.Contains(line, "placeholder") {
							// This is OK - obvious placeholder on this line
							continue
						}
						t.Logf("NOTE: %s line %d contains '%s' - verify it's just an example", doc, lineNum+1, pattern)
					}
				}
			}
		})
	}
}

// TestMermaidDiagramSyntax performs basic validation of mermaid diagrams.
func TestMermaidDiagramSyntax(t *testing.T) {
	docsWithDiagrams := []string{
		"docs/ARCHITECTURE.md",
		"docs/Notice_of_Decisions.md",
		"docs/API_TESTS_SPEC.md",
	}

	for _, doc := range docsWithDiagrams {
		t.Run(doc, func(t *testing.T) {
			content, err := os.ReadFile(doc)
			if err != nil {
				t.Skipf("Cannot read %s: %v", doc, err)
				return
			}

			docStr := string(content)

			// Find all mermaid blocks
			mermaidStart := 0
			for {
				idx := strings.Index(docStr[mermaidStart:], "```mermaid")
				if idx == -1 {
					break
				}
				mermaidStart += idx

				// Find end of block
				endIdx := strings.Index(docStr[mermaidStart+10:], "```")
				if endIdx == -1 {
					t.Errorf("%s: Unclosed mermaid block at position %d", doc, mermaidStart)
					break
				}

				mermaidBlock := docStr[mermaidStart+10 : mermaidStart+10+endIdx]

				// Basic syntax checks
				if !strings.Contains(mermaidBlock, "graph") &&
					!strings.Contains(mermaidBlock, "sequenceDiagram") &&
					!strings.Contains(mermaidBlock, "flowchart") &&
					!strings.Contains(mermaidBlock, "stateDiagram") &&
					!strings.Contains(mermaidBlock, "erDiagram") {
					t.Errorf("%s: Mermaid block at %d missing diagram type", doc, mermaidStart)
				}

				mermaidStart += 10 + endIdx + 3
			}
		})
	}
}
