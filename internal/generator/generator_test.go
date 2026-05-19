package generator_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/ahafin/appteam/internal/config"
	"github.com/ahafin/appteam/internal/generator"
)

// testConfig returns a ProjectConfig suitable for generator testing.
// InitGit and CreateRepo are false to avoid git/gh operations.
func testConfig(targetDir string) *config.ProjectConfig {
	return &config.ProjectConfig{
		ProjectName: "gentest",
		Description: "Generator test project",
		TechStack:   "Go",
		OwnerName:   "Test Owner",
		OwnerEmail:  "test@example.com",
		OwnerGitHub: "testowner",
		ModelName:   "Opus 4.6",
		ModelID:     "claude-opus-4-6",
		SWEs: []config.SWEConfig{
			{Number: 1, Title: "CLI & Wizard", Bullets: []string{"CLI flags"}},
			{Number: 2, Title: "Templates", Bullets: []string{"Template rendering"}},
		},
		GCP: config.GCPConfig{
			Enabled:       true,
			ProjectID:     "test-project",
			ProjectNumber: "123456",
			Organization:  "test-org",
			Region:        "us-central1",
		},
		IncludeSWETest:  true,
		IncludeSWEQA:    true,
		IncludePlatform: true,
		IncludeReviewer: true,
		Conventions:     []string{"Use stdlib only"},
		InitGit:         false,
		CreateRepo:      false,
		TargetDir:       targetDir,
	}
}

func TestGenerateCreatesDirectories(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	dirs := []string{
		filepath.Join(dir, ".claude", "agents"),
		filepath.Join(dir, ".claude", "skills"),
		filepath.Join(dir, "docs", "specs"),
	}
	for _, d := range dirs {
		info, err := os.Stat(d)
		if err != nil {
			t.Errorf("expected directory %s to exist: %v", d, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", d)
		}
	}
}

func TestGenerateCreatesExpectedFiles(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	expectedFiles := []string{
		"CLAUDE.md",
		"docs/BACKLOG.md",
		"docs/PROGRESS.md",
		"docs/RELEASENOTES.md",
		"docs/PIPELINE.md",
		"docs/specs/TEMPLATE.md",
		".claude/agents/pm.md",
		".claude/agents/tpm.md",
		".claude/agents/swe-1.md",
		".claude/agents/swe-2.md",
		".claude/agents/swe-test.md",
		".claude/agents/swe-qa.md",
		".claude/agents/platform.md",
		".claude/agents/reviewer.md",
		".claude/skills/spec.md",
		".claude/skills/release.md",
		".claude/skills/pipeline.md",
		".claude/skills/status.md",
		".claude/skills/regenerate.md",
		".claude/skills/roadmap.md",
		".claude/skills/debug.md",
		".claude/skills/test.md",
		".claude/skills/review.md",
		".claude/skills/docs.md",
		".claude/skills/refactor.md",
		".claude/skills/hotfix.md",
		".claude/skills/api-design.md",
		".claude/skills/schema.md",
		".claude/skills/deploy.md",
		".claude/skills/security.md",
		".claude/skills/adr.md",
		".claude/skills/standup.md",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(dir, f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", f, err)
		}
	}
}

func TestGenerateMinimalConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.ProjectConfig{
		ProjectName:     "minimal",
		Description:     "Minimal project",
		TechStack:       "Go",
		OwnerName:       "Owner",
		OwnerEmail:      "owner@example.com",
		OwnerGitHub:     "owner",
		ModelName:       "Opus 4.6",
		ModelID:         "claude-opus-4-6",
		SWEs:            []config.SWEConfig{{Number: 1, Title: "General"}},
		InitGit:         false,
		CreateRepo:      false,
		TargetDir:       dir,
		IncludeSWETest:  false,
		IncludeSWEQA:    false,
		IncludePlatform: false,
		IncludeReviewer: false,
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Core files must exist
	coreFiles := []string{
		"CLAUDE.md",
		"docs/BACKLOG.md",
		"docs/PROGRESS.md",
		"docs/RELEASENOTES.md",
		"docs/PIPELINE.md",
		"docs/specs/TEMPLATE.md",
		".claude/agents/pm.md",
		".claude/agents/tpm.md",
		".claude/agents/swe-1.md",
		".claude/skills/spec.md",
		".claude/skills/release.md",
		".claude/skills/pipeline.md",
		".claude/skills/status.md",
		".claude/skills/regenerate.md",
		".claude/skills/roadmap.md",
		".claude/skills/debug.md",
		".claude/skills/test.md",
		".claude/skills/review.md",
		".claude/skills/docs.md",
		".claude/skills/refactor.md",
		".claude/skills/hotfix.md",
		".claude/skills/api-design.md",
		".claude/skills/schema.md",
		".claude/skills/deploy.md",
		".claude/skills/security.md",
		".claude/skills/adr.md",
		".claude/skills/standup.md",
	}
	for _, f := range coreFiles {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("expected core file %s to exist: %v", f, err)
		}
	}

	// Optional agent files must NOT exist
	optionalFiles := []string{
		".claude/agents/swe-test.md",
		".claude/agents/swe-qa.md",
		".claude/agents/platform.md",
		".claude/agents/reviewer.md",
	}
	for _, f := range optionalFiles {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Errorf("optional file %s should NOT exist in minimal config", f)
		}
	}
}

func TestGenerateSettingsPersisted(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	settingsPath := filepath.Join(dir, ".appteam", "settings.json")
	if _, err := os.Stat(settingsPath); err != nil {
		t.Errorf("expected settings.json at %s: %v", settingsPath, err)
	}
}

func TestGeneratedFileContent(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Check CLAUDE.md content
	claudeData, err := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("reading CLAUDE.md: %v", err)
	}
	claudeContent := string(claudeData)
	if !strings.HasPrefix(claudeContent, "# gentest") {
		t.Errorf("CLAUDE.md should start with '# gentest', got %q", claudeContent[:min(50, len(claudeContent))])
	}
	if !strings.Contains(claudeContent, "gentest") {
		t.Error("CLAUDE.md should contain project name 'gentest'")
	}

	// Check SWE file content
	swe1Data, err := os.ReadFile(filepath.Join(dir, ".claude", "agents", "swe-1.md"))
	if err != nil {
		t.Fatalf("reading swe-1.md: %v", err)
	}
	swe1Content := string(swe1Data)
	if !strings.HasPrefix(swe1Content, "---\n") {
		t.Errorf("swe-1.md should start with YAML frontmatter '---', got %q", swe1Content[:min(50, len(swe1Content))])
	}
	if !strings.Contains(swe1Content, "# SWE-1 Agent") {
		t.Error("swe-1.md should contain '# SWE-1 Agent'")
	}

	swe2Data, err := os.ReadFile(filepath.Join(dir, ".claude", "agents", "swe-2.md"))
	if err != nil {
		t.Fatalf("reading swe-2.md: %v", err)
	}
	if !strings.HasPrefix(string(swe2Data), "---\n") {
		t.Errorf("swe-2.md should start with YAML frontmatter '---'")
	}
	if !strings.Contains(string(swe2Data), "# SWE-2 Agent") {
		t.Error("swe-2.md should contain '# SWE-2 Agent'")
	}
}

func TestGenerateConditionalSkillsAllSelected(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.SelectedSkills = map[string]bool{
		"spec": true, "release": true, "pipeline": true,
		"status": true, "regenerate": true, "roadmap": true,
		"debug": true, "test": true, "review": true,
		"docs": true, "refactor": true, "hotfix": true,
		"api-design": true, "schema": true, "deploy": true,
		"security": true, "adr": true, "standup": true,
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	allSkills := []string{
		".claude/skills/spec.md",
		".claude/skills/release.md",
		".claude/skills/pipeline.md",
		".claude/skills/status.md",
		".claude/skills/regenerate.md",
		".claude/skills/roadmap.md",
		".claude/skills/debug.md",
		".claude/skills/test.md",
		".claude/skills/review.md",
		".claude/skills/docs.md",
		".claude/skills/refactor.md",
		".claude/skills/hotfix.md",
		".claude/skills/api-design.md",
		".claude/skills/schema.md",
		".claude/skills/deploy.md",
		".claude/skills/security.md",
		".claude/skills/adr.md",
		".claude/skills/standup.md",
	}
	for _, f := range allSkills {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("expected skill file %s to exist when all selected: %v", f, err)
		}
	}
}

func TestGenerateConditionalSkillsSomeSelected(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.SelectedSkills = map[string]bool{
		"spec": true, "release": false, "pipeline": true,
		"status": false, "regenerate": true, "roadmap": false,
		"debug": true, "test": false, "review": true,
		"docs": false, "refactor": true, "hotfix": false,
		"api-design": true, "schema": false, "deploy": true,
		"security": false, "adr": true, "standup": false,
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	shouldExist := []string{
		".claude/skills/spec.md",
		".claude/skills/pipeline.md",
		".claude/skills/regenerate.md",
		".claude/skills/debug.md",
		".claude/skills/review.md",
		".claude/skills/refactor.md",
		".claude/skills/api-design.md",
		".claude/skills/deploy.md",
		".claude/skills/adr.md",
	}
	for _, f := range shouldExist {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("expected selected skill file %s to exist: %v", f, err)
		}
	}

	shouldNotExist := []string{
		".claude/skills/release.md",
		".claude/skills/status.md",
		".claude/skills/roadmap.md",
		".claude/skills/test.md",
		".claude/skills/docs.md",
		".claude/skills/hotfix.md",
		".claude/skills/schema.md",
		".claude/skills/security.md",
		".claude/skills/standup.md",
	}
	for _, f := range shouldNotExist {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Errorf("unselected skill file %s should NOT exist", f)
		}
	}

	// Skills directory must still exist even with deselected skills
	if _, err := os.Stat(filepath.Join(dir, ".claude", "skills")); err != nil {
		t.Errorf("skills directory should always exist: %v", err)
	}
}

func TestGenerateConditionalSkillsNoneSelected(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.SelectedSkills = map[string]bool{
		"spec": false, "release": false, "pipeline": false,
		"status": false, "regenerate": false, "roadmap": false,
		"debug": false, "test": false, "review": false,
		"docs": false, "refactor": false, "hotfix": false,
		"api-design": false, "schema": false, "deploy": false,
		"security": false, "adr": false, "standup": false,
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	noSkills := []string{
		".claude/skills/spec.md",
		".claude/skills/release.md",
		".claude/skills/pipeline.md",
		".claude/skills/status.md",
		".claude/skills/regenerate.md",
		".claude/skills/roadmap.md",
		".claude/skills/debug.md",
		".claude/skills/test.md",
		".claude/skills/review.md",
		".claude/skills/docs.md",
		".claude/skills/refactor.md",
		".claude/skills/hotfix.md",
		".claude/skills/api-design.md",
		".claude/skills/schema.md",
		".claude/skills/deploy.md",
		".claude/skills/security.md",
		".claude/skills/adr.md",
		".claude/skills/standup.md",
	}
	for _, f := range noSkills {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Errorf("skill file %s should NOT exist when none selected", f)
		}
	}

	// Skills directory must still exist
	if _, err := os.Stat(filepath.Join(dir, ".claude", "skills")); err != nil {
		t.Errorf("skills directory should always exist: %v", err)
	}
}

func TestGenerateConditionalSkillsNilMap(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.SelectedSkills = nil // pre-v0.7.0 backward compat

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// All 18 skills should be generated when SelectedSkills is nil
	allSkills := []string{
		".claude/skills/spec.md",
		".claude/skills/release.md",
		".claude/skills/pipeline.md",
		".claude/skills/status.md",
		".claude/skills/regenerate.md",
		".claude/skills/roadmap.md",
		".claude/skills/debug.md",
		".claude/skills/test.md",
		".claude/skills/review.md",
		".claude/skills/docs.md",
		".claude/skills/refactor.md",
		".claude/skills/hotfix.md",
		".claude/skills/api-design.md",
		".claude/skills/schema.md",
		".claude/skills/deploy.md",
		".claude/skills/security.md",
		".claude/skills/adr.md",
		".claude/skills/standup.md",
	}
	for _, f := range allSkills {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("expected skill file %s to exist when SelectedSkills is nil (backward compat): %v", f, err)
		}
	}
}

func TestGenerateCustomAgents(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.CustomAgents = []config.CustomAgentConfig{
		{
			Name:         "frontend-eng",
			Title:        "Frontend Engineer",
			Description:  "React/TypeScript specialist for UI components",
			Instructions: []string{"Build responsive UI components", "Follow accessibility standards"},
		},
		{
			Name:         "dba",
			Title:        "Database Administrator",
			Description:  "Manages database schemas and migrations",
			Instructions: []string{"Review schema changes", "Optimize queries"},
		},
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Custom agent files should exist with correct content
	for _, ca := range cfg.CustomAgents {
		path := filepath.Join(dir, ".claude", "agents", ca.Name+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected custom agent file %s.md to exist: %v", ca.Name, err)
			continue
		}
		content := string(data)
		if !strings.Contains(content, ca.Title) {
			t.Errorf("custom agent %s should contain title %q", ca.Name, ca.Title)
		}
		if !strings.Contains(content, ca.Description) {
			t.Errorf("custom agent %s should contain description", ca.Name)
		}
		if !strings.Contains(content, "gentest") {
			t.Errorf("custom agent %s should contain project name", ca.Name)
		}
		for _, inst := range ca.Instructions {
			if !strings.Contains(content, inst) {
				t.Errorf("custom agent %s should contain instruction %q", ca.Name, inst)
			}
		}
	}

	// Built-in agents should still exist
	builtInAgents := []string{"pm.md", "tpm.md", "swe-1.md", "swe-2.md", "swe-test.md", "swe-qa.md", "platform.md", "reviewer.md"}
	for _, f := range builtInAgents {
		path := filepath.Join(dir, ".claude", "agents", f)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("built-in agent %s should still exist alongside custom agents: %v", f, err)
		}
	}
}

func TestGenerateCustomAgentsEmpty(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.CustomAgents = nil // no custom agents

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Only built-in agents should exist, no extra files in agents dir
	builtInAgents := []string{"pm.md", "tpm.md", "swe-1.md", "swe-2.md", "swe-test.md", "swe-qa.md", "platform.md", "reviewer.md"}
	for _, f := range builtInAgents {
		if _, err := os.Stat(filepath.Join(dir, ".claude", "agents", f)); err != nil {
			t.Errorf("built-in agent %s should exist: %v", f, err)
		}
	}
}

func TestGenerateTeamSizeLean(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.TeamSize = "lean"
	// Even with multiple SWEs configured, lean only generates SWE-1
	cfg.SWEs = []config.SWEConfig{
		{Number: 1, Title: "CLI & Wizard", Bullets: []string{"CLI flags"}},
		{Number: 2, Title: "Templates", Bullets: []string{"Template rendering"}},
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Lean: PM, SWE-1, SWE-Test only
	mustExist := []string{
		".claude/agents/pm.md",
		".claude/agents/swe-1.md",
		".claude/agents/swe-test.md",
	}
	for _, f := range mustExist {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("lean team should generate %s: %v", f, err)
		}
	}

	// Lean: no TPM, no SWE-2+, no optional agents
	mustNotExist := []string{
		".claude/agents/tpm.md",
		".claude/agents/swe-2.md",
		".claude/agents/swe-qa.md",
		".claude/agents/platform.md",
		".claude/agents/reviewer.md",
	}
	for _, f := range mustNotExist {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Errorf("lean team should NOT generate %s", f)
		}
	}
}

func TestGenerateTeamSizeFull(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.TeamSize = "full"
	cfg.SWEs = []config.SWEConfig{
		{Number: 1, Title: "CLI", Bullets: []string{"CLI"}},
		{Number: 2, Title: "Templates", Bullets: []string{"Templates"}},
		{Number: 3, Title: "API", Bullets: []string{"API"}},
	}
	// Even if optional flags are false, full generates all
	cfg.IncludeSWETest = false
	cfg.IncludeSWEQA = false
	cfg.IncludePlatform = false
	cfg.IncludeReviewer = false

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	mustExist := []string{
		".claude/agents/pm.md",
		".claude/agents/tpm.md",
		".claude/agents/swe-1.md",
		".claude/agents/swe-2.md",
		".claude/agents/swe-3.md",
		".claude/agents/swe-test.md",
		".claude/agents/swe-qa.md",
		".claude/agents/platform.md",
		".claude/agents/reviewer.md",
	}
	for _, f := range mustExist {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("full team should generate %s: %v", f, err)
		}
	}
}

func TestGenerateTeamSizeStandard(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.TeamSize = "standard"
	cfg.IncludeSWETest = true
	cfg.IncludeReviewer = true
	cfg.IncludeSWEQA = false
	cfg.IncludePlatform = false

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	mustExist := []string{
		".claude/agents/pm.md",
		".claude/agents/tpm.md",
		".claude/agents/swe-1.md",
		".claude/agents/swe-2.md",
		".claude/agents/swe-test.md",
		".claude/agents/reviewer.md",
	}
	for _, f := range mustExist {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("standard team should generate %s: %v", f, err)
		}
	}

	mustNotExist := []string{
		".claude/agents/swe-qa.md",
		".claude/agents/platform.md",
	}
	for _, f := range mustNotExist {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Errorf("standard team (no QA/Platform) should NOT generate %s", f)
		}
	}
}

func TestGenerateTeamSizeLeanWithCustomAgents(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.TeamSize = "lean"
	cfg.CustomAgents = []config.CustomAgentConfig{
		{
			Name:         "dba",
			Title:        "Database Administrator",
			Description:  "Manages databases",
			Instructions: []string{"Optimize queries"},
		},
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Custom agents still generated in lean mode
	if _, err := os.Stat(filepath.Join(dir, ".claude", "agents", "dba.md")); err != nil {
		t.Error("lean team should still generate custom agents")
	}
	// But no TPM
	if _, err := os.Stat(filepath.Join(dir, ".claude", "agents", "tpm.md")); err == nil {
		t.Error("lean team should NOT generate tpm.md")
	}
}

func TestGenerateCustomSkills(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir)
	cfg.SelectedSkills = map[string]bool{
		"spec": true, "release": false, "pipeline": false,
		"status": false, "regenerate": false, "roadmap": false,
		"debug": false, "test": false, "review": false,
		"docs": false, "refactor": false, "hotfix": false,
		"api-design": false, "schema": false, "deploy": false,
		"security": false, "adr": false, "standup": false,
	}
	cfg.CustomSkills = []config.CustomSkillConfig{
		{Name: "lint", Description: "Run linter and fix issues"},
		{Name: "migrate", Description: "Create database migration"},
	}

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// Custom skill files should exist
	for _, name := range []string{"lint", "migrate"} {
		path := filepath.Join(dir, ".claude", "skills", name+".md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("expected custom skill file %s.md to exist: %v", name, err)
			continue
		}
		content := string(data)
		if !strings.Contains(content, "/"+name) {
			t.Errorf("custom skill %s should contain trigger /%s", name, name)
		}
		if !strings.Contains(content, "gentest") {
			t.Errorf("custom skill %s should contain project name", name)
		}
	}
}

// scionTestConfig returns a Scion framework ProjectConfig for testing.
func scionTestConfig(targetDir string) *config.ProjectConfig {
	cfg := testConfig(targetDir)
	cfg.Framework = "scion"
	cfg.DefaultHarness = "claude"
	return cfg
}

func TestGenerateScionGroveID(t *testing.T) {
	dir := t.TempDir()
	cfg := scionTestConfig(dir)

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// grove-id must exist
	groveIDPath := filepath.Join(dir, ".scion", "grove-id")
	data, err := os.ReadFile(groveIDPath)
	if err != nil {
		t.Fatalf("expected .scion/grove-id to exist: %v", err)
	}

	// Must contain a valid UUID v4 (8-4-4-4-12 hex)
	content := strings.TrimSpace(string(data))
	uuidPattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidPattern.MatchString(content) {
		t.Errorf("grove-id should contain a valid UUID v4, got %q", content)
	}
}

func TestGenerateScionAgentsDir(t *testing.T) {
	dir := t.TempDir()
	cfg := scionTestConfig(dir)

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// .scion/agents/ directory must exist
	agentsDir := filepath.Join(dir, ".scion", "agents")
	info, err := os.Stat(agentsDir)
	if err != nil {
		t.Fatalf("expected .scion/agents/ to exist: %v", err)
	}
	if !info.IsDir() {
		t.Error(".scion/agents/ should be a directory")
	}
}

func TestGenerateClaudeCodeNoGroveID(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig(dir) // framework defaults to claude-code

	if err := generator.Generate(cfg, false); err != nil {
		t.Fatalf("Generate() error: %v", err)
	}

	// grove-id must NOT exist for claude-code framework
	groveIDPath := filepath.Join(dir, ".scion", "grove-id")
	if _, err := os.Stat(groveIDPath); err == nil {
		t.Error(".scion/grove-id should NOT exist for claude-code framework projects")
	}
}

func TestGenerateScionGroveIDUnique(t *testing.T) {
	// Generate two Scion projects and verify UUIDs are different
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	cfg1 := scionTestConfig(dir1)
	cfg2 := scionTestConfig(dir2)

	if err := generator.Generate(cfg1, false); err != nil {
		t.Fatalf("Generate() dir1 error: %v", err)
	}
	if err := generator.Generate(cfg2, false); err != nil {
		t.Fatalf("Generate() dir2 error: %v", err)
	}

	data1, _ := os.ReadFile(filepath.Join(dir1, ".scion", "grove-id"))
	data2, _ := os.ReadFile(filepath.Join(dir2, ".scion", "grove-id"))
	uuid1 := strings.TrimSpace(string(data1))
	uuid2 := strings.TrimSpace(string(data2))

	if uuid1 == uuid2 {
		t.Errorf("grove-id should be unique across projects, both got %q", uuid1)
	}
}
