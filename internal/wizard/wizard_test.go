package wizard

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestAsk(t *testing.T) {
	s := NewStyler(false)
	w := io.Discard

	tests := []struct {
		name       string
		input      string
		defaultVal string
		want       string
	}{
		{"user input", "myproject\n", "", "myproject"},
		{"empty with default", "\n", "fallback", "fallback"},
		{"whitespace trimmed", "  hello  \n", "", "hello"},
		{"user overrides default", "custom\n", "fallback", "custom"},
	}
	for _, tt := range tests {
		scanner := bufio.NewScanner(strings.NewReader(tt.input))
		got := ask(scanner, w, s, "prompt", tt.defaultVal)
		if got != tt.want {
			t.Errorf("ask(%q, default=%q) = %q, want %q", tt.name, tt.defaultVal, got, tt.want)
		}
	}
}

func TestAskBool(t *testing.T) {
	s := NewStyler(false)
	w := io.Discard

	tests := []struct {
		name       string
		input      string
		defaultVal bool
		want       bool
	}{
		{"y", "y\n", false, true},
		{"n", "n\n", true, false},
		{"yes", "yes\n", false, true},
		{"no", "no\n", true, false},
		{"empty default true", "\n", true, true},
		{"empty default false", "\n", false, false},
		{"Y uppercase", "Y\n", false, true},
		{"N uppercase", "N\n", true, false},
	}
	for _, tt := range tests {
		scanner := bufio.NewScanner(strings.NewReader(tt.input))
		got := askBool(scanner, w, s, "prompt", tt.defaultVal)
		if got != tt.want {
			t.Errorf("askBool(%q, default=%v) = %v, want %v", tt.name, tt.defaultVal, got, tt.want)
		}
	}
}

func TestAskInt(t *testing.T) {
	s := NewStyler(false)
	w := io.Discard

	tests := []struct {
		name       string
		input      string
		defaultVal int
		want       int
	}{
		{"valid number", "3\n", 5, 3},
		{"empty default", "\n", 5, 5},
		{"non-numeric", "abc\n", 7, 7},
		{"zero", "0\n", 5, 0},
	}
	for _, tt := range tests {
		scanner := bufio.NewScanner(strings.NewReader(tt.input))
		got := askInt(scanner, w, s, "prompt", tt.defaultVal)
		if got != tt.want {
			t.Errorf("askInt(%q, default=%d) = %d, want %d", tt.name, tt.defaultVal, got, tt.want)
		}
	}
}

func TestRunIntegration(t *testing.T) {
	// Simulate full wizard input:
	// Step 1: Project Basics (4 lines)
	// Step 2: Team Framework (1 line)
	// Step 3: Git Repository (no .git dir, so asks init git) (1 line)
	// Step 4: Product Owner (3 lines)
	// Step 5: GCP Configuration (1 line)
	// Step 6: Agent Team (1 SWE, with bullets, optional agents, model)
	// Step 7: Conventions (2 lines)
	// Step 8: Skills (18 skill Y/n + custom skills + blank)
	// Step 9: Confirm (1 line)
	lines := []string{
		"testproject",          // Project name
		"A test project",      // Description
		"Go",                  // Tech stack
		t.TempDir(),           // Target directory (temp dir, no .git)
		"",                    // Framework: default (Claude Code)
		"n",                   // Init git? no
		"Test User",           // Owner name
		"test@example.com",    // Owner email
		"testuser",            // Owner GitHub
		"n",                   // Include GCP? no
		"2",                   // Team size: standard
		"1",                   // Number of SWEs
		"Backend Engineer",    // SWE-1 title
		"Go APIs",             // SWE-1 bullet 1
		"",                    // End bullets
		"n",                   // Include Platform?
		"n",                   // Include Reviewer?
		"n",                   // Include SWE-Test?
		"n",                   // Include SWE-QA?
		"y",                   // Add a custom agent?
		"frontend-eng",        // Agent name
		"Frontend Engineer",   // Title
		"React specialist",    // Description
		"Build UI components", // Instruction 1
		"",                    // End instructions
		"n",                   // Add another custom agent?
		"1",                   // Model choice (Opus)
		"Use tabs for indent", // Convention 1
		"",                    // End conventions
		// PM Skills (7)
		"y",                   // /spec
		"y",                   // /release
		"y",                   // /roadmap
		"y",                   // /status
		"y",                   // /pipeline
		"y",                   // /regenerate
		"y",                   // /brainstorm
		// Dev Skills Tier 1 (6)
		"y",                   // /debug
		"y",                   // /test
		"y",                   // /review
		"y",                   // /docs
		"y",                   // /refactor
		"n",                   // /hotfix (deselect)
		// Dev Skills Tier 2 (8)
		"n",                   // /api-design
		"y",                   // /schema (select)
		"n",                   // /deploy
		"n",                   // /security
		"n",                   // /adr
		"n",                   // /standup
		"n",                   // /cuj-list
		"n",                   // /cuj-test
		// Custom skills
		"migrate: Database migration workflow",
		"",                    // End custom skills
		"y",                   // Proceed?
	}

	input := strings.Join(lines, "\n") + "\n"
	var output bytes.Buffer

	cfg, err := Run(strings.NewReader(input), &output, false, "")
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if cfg.ProjectName != "testproject" {
		t.Errorf("ProjectName = %q, want %q", cfg.ProjectName, "testproject")
	}
	if cfg.Description != "A test project" {
		t.Errorf("Description = %q, want %q", cfg.Description, "A test project")
	}
	if cfg.TechStack != "Go" {
		t.Errorf("TechStack = %q, want %q", cfg.TechStack, "Go")
	}
	if cfg.OwnerName != "Test User" {
		t.Errorf("OwnerName = %q, want %q", cfg.OwnerName, "Test User")
	}
	if cfg.OwnerEmail != "test@example.com" {
		t.Errorf("OwnerEmail = %q, want %q", cfg.OwnerEmail, "test@example.com")
	}
	if cfg.OwnerGitHub != "testuser" {
		t.Errorf("OwnerGitHub = %q, want %q", cfg.OwnerGitHub, "testuser")
	}
	if cfg.TeamSize != "standard" {
		t.Errorf("TeamSize = %q, want %q", cfg.TeamSize, "standard")
	}
	if cfg.GCP.Enabled {
		t.Error("GCP should be disabled")
	}
	if len(cfg.SWEs) != 1 {
		t.Fatalf("len(SWEs) = %d, want 1", len(cfg.SWEs))
	}
	if cfg.SWEs[0].Title != "Backend Engineer" {
		t.Errorf("SWEs[0].Title = %q, want %q", cfg.SWEs[0].Title, "Backend Engineer")
	}
	if len(cfg.SWEs[0].Bullets) != 1 || cfg.SWEs[0].Bullets[0] != "Go APIs" {
		t.Errorf("SWEs[0].Bullets = %v, want [Go APIs]", cfg.SWEs[0].Bullets)
	}
	if cfg.IncludePlatform {
		t.Error("IncludePlatform should be false")
	}
	if cfg.IncludeReviewer {
		t.Error("IncludeReviewer should be false")
	}
	if cfg.ModelName != "Opus 4.6" {
		t.Errorf("ModelName = %q, want %q", cfg.ModelName, "Opus 4.6")
	}
	if cfg.ModelID != "claude-opus-4-6" {
		t.Errorf("ModelID = %q, want %q", cfg.ModelID, "claude-opus-4-6")
	}
	if len(cfg.Conventions) != 1 || cfg.Conventions[0] != "Use tabs for indent" {
		t.Errorf("Conventions = %v, want [Use tabs for indent]", cfg.Conventions)
	}
	// Check SelectedSkills (21 keys)
	if cfg.SelectedSkills == nil {
		t.Fatal("SelectedSkills should not be nil")
	}
	if len(cfg.SelectedSkills) != 21 {
		t.Errorf("len(SelectedSkills) = %d, want 21", len(cfg.SelectedSkills))
	}
	expectedSkills := map[string]bool{
		"spec": true, "release": true, "roadmap": true,
		"status": true, "pipeline": true, "regenerate": true, "brainstorm": true,
		"debug": true, "test": true, "review": true,
		"docs": true, "refactor": true, "hotfix": false,
		"api-design": false, "schema": true, "deploy": false,
		"security": false, "adr": false, "standup": false,
		"cuj-list": false, "cuj-test": false,
	}
	for name, want := range expectedSkills {
		if cfg.SelectedSkills[name] != want {
			t.Errorf("SelectedSkills[%q] = %v, want %v", name, cfg.SelectedSkills[name], want)
		}
	}
	// Check CustomSkills
	if len(cfg.CustomSkills) != 1 {
		t.Fatalf("len(CustomSkills) = %d, want 1", len(cfg.CustomSkills))
	}
	if cfg.CustomSkills[0].Name != "migrate" {
		t.Errorf("CustomSkills[0].Name = %q, want %q", cfg.CustomSkills[0].Name, "migrate")
	}
	if cfg.CustomSkills[0].Description != "Database migration workflow" {
		t.Errorf("CustomSkills[0].Description = %q, want %q", cfg.CustomSkills[0].Description, "Database migration workflow")
	}
	// Check CustomAgents
	if len(cfg.CustomAgents) != 1 {
		t.Fatalf("len(CustomAgents) = %d, want 1", len(cfg.CustomAgents))
	}
	if cfg.CustomAgents[0].Name != "frontend-eng" {
		t.Errorf("CustomAgents[0].Name = %q, want %q", cfg.CustomAgents[0].Name, "frontend-eng")
	}
	if cfg.CustomAgents[0].Title != "Frontend Engineer" {
		t.Errorf("CustomAgents[0].Title = %q, want %q", cfg.CustomAgents[0].Title, "Frontend Engineer")
	}
	if cfg.CustomAgents[0].Description != "React specialist" {
		t.Errorf("CustomAgents[0].Description = %q, want %q", cfg.CustomAgents[0].Description, "React specialist")
	}
	if len(cfg.CustomAgents[0].Instructions) != 1 || cfg.CustomAgents[0].Instructions[0] != "Build UI components" {
		t.Errorf("CustomAgents[0].Instructions = %v, want [Build UI components]", cfg.CustomAgents[0].Instructions)
	}
}

func TestRunCancelled(t *testing.T) {
	// Same wizard flow but answer "n" to proceed
	lines := []string{
		"cancelproject",       // Project name
		"A cancel test",       // Description
		"Python",              // Tech stack
		t.TempDir(),           // Target directory
		"",                    // Framework: default (Claude Code)
		"n",                   // Init git? no
		"Cancel User",         // Owner name
		"cancel@example.com",  // Owner email
		"canceluser",          // Owner GitHub
		"n",                   // Include GCP? no
		"2",                   // Team size: standard
		"1",                   // Number of SWEs
		"General",             // SWE-1 title
		"",                    // End bullets (no bullets)
		"n",                   // Include Platform?
		"n",                   // Include Reviewer?
		"n",                   // Include SWE-Test?
		"n",                   // Include SWE-QA?
		"n",                   // Add a custom agent? no
		"1",                   // Model
		"",                    // End conventions (no conventions)
		// PM Skills (6)
		"",                    // /spec (default Y)
		"",                    // /release (default Y)
		"",                    // /roadmap (default Y)
		"",                    // /status (default Y)
		"",                    // /pipeline (default Y)
		"",                    // /regenerate (default Y)
		"",                    // /brainstorm (default Y)
		// Dev Skills Tier 1 (6)
		"",                    // /debug (default Y)
		"",                    // /test (default Y)
		"",                    // /review (default Y)
		"",                    // /docs (default Y)
		"",                    // /refactor (default Y)
		"",                    // /hotfix (default Y)
		// Dev Skills Tier 2 (8)
		"",                    // /api-design (default N)
		"",                    // /schema (default N)
		"",                    // /deploy (default N)
		"",                    // /security (default N)
		"",                    // /adr (default N)
		"",                    // /standup (default N)
		"",                    // /cuj-list (default N)
		"",                    // /cuj-test (default N)
		// Custom skills
		"",                    // End custom skills (none)
		"n",                   // Proceed? NO
	}

	input := strings.Join(lines, "\n") + "\n"
	var output bytes.Buffer

	_, err := Run(strings.NewReader(input), &output, false, "")
	if err == nil {
		t.Fatal("Run() should return an error when user cancels")
	}
	if !strings.Contains(err.Error(), "cancelled") {
		t.Errorf("error should mention 'cancelled', got: %v", err)
	}
}

func TestRunLeanTeamSize(t *testing.T) {
	lines := []string{
		"leanproject",         // Project name
		"A lean test",         // Description
		"Go",                  // Tech stack
		t.TempDir(),           // Target directory
		"",                    // Framework: default (Claude Code)
		"n",                   // Init git? no
		"Test User",           // Owner name
		"test@example.com",    // Owner email
		"testuser",            // Owner GitHub
		"n",                   // Include GCP? no
		"1",                   // Team size: lean
		// No SWE count prompt (hardcoded to 1)
		"Solo Engineer",       // SWE-1 title
		"",                    // End bullets
		// No optional agent prompts (all skipped for lean)
		"n",                   // Add a custom agent? no
		"1",                   // Model choice (Opus)
		"",                    // End conventions
		// PM Skills (7)
		"", "", "", "", "", "", "",
		// Dev Skills Tier 1 (6)
		"", "", "", "", "", "",
		// Dev Skills Tier 2 (8)
		"", "", "", "", "", "", "", "",
		// Custom skills
		"",                    // End custom skills
		"y",                   // Proceed?
	}

	input := strings.Join(lines, "\n") + "\n"
	var output bytes.Buffer

	cfg, err := Run(strings.NewReader(input), &output, false, "")
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if cfg.TeamSize != "lean" {
		t.Errorf("TeamSize = %q, want %q", cfg.TeamSize, "lean")
	}
	if len(cfg.SWEs) != 1 {
		t.Fatalf("len(SWEs) = %d, want 1", len(cfg.SWEs))
	}
	if cfg.SWEs[0].Title != "Solo Engineer" {
		t.Errorf("SWEs[0].Title = %q, want %q", cfg.SWEs[0].Title, "Solo Engineer")
	}
	if !cfg.IncludeSWETest {
		t.Error("IncludeSWETest should be true for lean")
	}
	if cfg.IncludePlatform {
		t.Error("IncludePlatform should be false for lean")
	}
	if cfg.IncludeReviewer {
		t.Error("IncludeReviewer should be false for lean")
	}
	if cfg.IncludeSWEQA {
		t.Error("IncludeSWEQA should be false for lean")
	}
}

func TestRunFullTeamSize(t *testing.T) {
	lines := []string{
		"fullproject",         // Project name
		"A full test",         // Description
		"Go",                  // Tech stack
		t.TempDir(),           // Target directory
		"",                    // Framework: default (Claude Code)
		"n",                   // Init git? no
		"Test User",           // Owner name
		"test@example.com",    // Owner email
		"testuser",            // Owner GitHub
		"n",                   // Include GCP? no
		"3",                   // Team size: full
		"",                    // SWE count: default 5
		// 5 SWEs with titles + end bullets
		"Frontend",            // SWE-1 title
		"",                    // End bullets
		"Backend",             // SWE-2 title
		"",                    // End bullets
		"Infra",               // SWE-3 title
		"",                    // End bullets
		"Data",                // SWE-4 title
		"",                    // End bullets
		"Mobile",              // SWE-5 title
		"",                    // End bullets
		// No optional agent prompts (all auto-enabled for full)
		"n",                   // Add a custom agent? no
		"1",                   // Model choice (Opus)
		"",                    // End conventions
		// PM Skills (7)
		"", "", "", "", "", "", "",
		// Dev Skills Tier 1 (6)
		"", "", "", "", "", "",
		// Dev Skills Tier 2 (8)
		"", "", "", "", "", "", "", "",
		// Custom skills
		"",                    // End custom skills
		"y",                   // Proceed?
	}

	input := strings.Join(lines, "\n") + "\n"
	var output bytes.Buffer

	cfg, err := Run(strings.NewReader(input), &output, false, "")
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	if cfg.TeamSize != "full" {
		t.Errorf("TeamSize = %q, want %q", cfg.TeamSize, "full")
	}
	if len(cfg.SWEs) != 5 {
		t.Fatalf("len(SWEs) = %d, want 5", len(cfg.SWEs))
	}
	if cfg.SWEs[0].Title != "Frontend" {
		t.Errorf("SWEs[0].Title = %q, want %q", cfg.SWEs[0].Title, "Frontend")
	}
	if cfg.SWEs[4].Title != "Mobile" {
		t.Errorf("SWEs[4].Title = %q, want %q", cfg.SWEs[4].Title, "Mobile")
	}
	if !cfg.IncludePlatform {
		t.Error("IncludePlatform should be true for full")
	}
	if !cfg.IncludeReviewer {
		t.Error("IncludeReviewer should be true for full")
	}
	if !cfg.IncludeSWETest {
		t.Error("IncludeSWETest should be true for full")
	}
	if !cfg.IncludeSWEQA {
		t.Error("IncludeSWEQA should be true for full")
	}
}
