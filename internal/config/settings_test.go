package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSettingsPath(t *testing.T) {
	tests := []struct {
		dir  string
		want string
	}{
		{"/tmp/project", "/tmp/project/.appteam/settings.json"},
		{".", ".appteam/settings.json"},
		{"/home/user/app", "/home/user/app/.appteam/settings.json"},
	}
	for _, tt := range tests {
		got := SettingsPath(tt.dir)
		want := filepath.FromSlash(tt.want)
		if got != want {
			t.Errorf("SettingsPath(%q) = %q, want %q", tt.dir, got, want)
		}
	}
}

func TestSettingsExist(t *testing.T) {
	dir := t.TempDir()

	// No settings yet
	if SettingsExist(dir) {
		t.Fatal("SettingsExist should return false when no settings file exists")
	}

	// Create settings file
	cfg := &ProjectConfig{ProjectName: "test"}
	if err := SaveSettings(dir, cfg); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	if !SettingsExist(dir) {
		t.Fatal("SettingsExist should return true after saving settings")
	}
}

func TestSaveAndLoadSettings(t *testing.T) {
	dir := t.TempDir()

	original := &ProjectConfig{
		ProjectName:     "myproject",
		Description:     "A test project",
		TechStack:       "Go",
		InitGit:         true,
		CreateRepo:      true,
		RepoURL:         "https://github.com/test/myproject.git",
		GitHubOrg:       "test",
		OwnerName:       "Test User",
		OwnerEmail:      "test@example.com",
		OwnerGitHub:     "testuser",
		IncludePlatform: true,
		IncludeReviewer: true,
		IncludeSWETest:  true,
		IncludeSWEQA:    false,
		ModelName:       "Opus 4.6",
		ModelID:         "claude-opus-4-6",
		TargetDir:       "/tmp/target",
		Conventions:     []string{"Use tabs", "No globals"},
		SelectedSkills: map[string]bool{
			"spec": true, "release": true, "pipeline": true,
			"status": true, "regenerate": false, "roadmap": true,
			"debug": true, "test": true, "review": true,
			"docs": true, "refactor": true, "hotfix": true,
			"api-design": false, "schema": false, "deploy": false,
			"security": true, "adr": false, "standup": false,
			"brainstorm": true,
			"cuj-list": false, "cuj-test": false,
		},
		CustomSkills: []CustomSkillConfig{
			{Name: "migrate", Description: "Database migration workflow"},
			{Name: "perf", Description: "Performance profiling"},
		},
		CustomAgents: []CustomAgentConfig{
			{
				Name:         "frontend-eng",
				Title:        "Frontend Engineer",
				Description:  "React/TypeScript specialist for UI components",
				Instructions: []string{"Build responsive UI components", "Follow design system patterns"},
			},
			{
				Name:         "dba",
				Title:        "Database Administrator",
				Description:  "Database design and optimization",
				Instructions: []string{"Design schemas", "Optimize queries"},
			},
		},
		GCP: GCPConfig{
			Enabled:       true,
			ProjectID:     "my-gcp-project",
			ProjectNumber: "123456789",
			Organization:  "my-org",
			Region:        "us-west1",
		},
		SWEs: []SWEConfig{
			{Number: 1, Title: "Frontend", Bullets: []string{"React", "CSS"}},
			{Number: 2, Title: "Backend", Bullets: []string{"Go", "APIs"}},
		},
	}

	if err := SaveSettings(dir, original); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	loaded, err := LoadSettings(dir)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}

	// Check top-level fields
	if loaded.ProjectName != original.ProjectName {
		t.Errorf("ProjectName = %q, want %q", loaded.ProjectName, original.ProjectName)
	}
	if loaded.Description != original.Description {
		t.Errorf("Description = %q, want %q", loaded.Description, original.Description)
	}
	if loaded.TechStack != original.TechStack {
		t.Errorf("TechStack = %q, want %q", loaded.TechStack, original.TechStack)
	}
	if loaded.InitGit != original.InitGit {
		t.Errorf("InitGit = %v, want %v", loaded.InitGit, original.InitGit)
	}
	if loaded.CreateRepo != original.CreateRepo {
		t.Errorf("CreateRepo = %v, want %v", loaded.CreateRepo, original.CreateRepo)
	}
	if loaded.RepoURL != original.RepoURL {
		t.Errorf("RepoURL = %q, want %q", loaded.RepoURL, original.RepoURL)
	}
	if loaded.OwnerName != original.OwnerName {
		t.Errorf("OwnerName = %q, want %q", loaded.OwnerName, original.OwnerName)
	}
	if loaded.OwnerEmail != original.OwnerEmail {
		t.Errorf("OwnerEmail = %q, want %q", loaded.OwnerEmail, original.OwnerEmail)
	}
	if loaded.OwnerGitHub != original.OwnerGitHub {
		t.Errorf("OwnerGitHub = %q, want %q", loaded.OwnerGitHub, original.OwnerGitHub)
	}
	if loaded.IncludePlatform != original.IncludePlatform {
		t.Errorf("IncludePlatform = %v, want %v", loaded.IncludePlatform, original.IncludePlatform)
	}
	if loaded.IncludeReviewer != original.IncludeReviewer {
		t.Errorf("IncludeReviewer = %v, want %v", loaded.IncludeReviewer, original.IncludeReviewer)
	}
	if loaded.IncludeSWETest != original.IncludeSWETest {
		t.Errorf("IncludeSWETest = %v, want %v", loaded.IncludeSWETest, original.IncludeSWETest)
	}
	if loaded.IncludeSWEQA != original.IncludeSWEQA {
		t.Errorf("IncludeSWEQA = %v, want %v", loaded.IncludeSWEQA, original.IncludeSWEQA)
	}
	if loaded.ModelName != original.ModelName {
		t.Errorf("ModelName = %q, want %q", loaded.ModelName, original.ModelName)
	}
	if loaded.ModelID != original.ModelID {
		t.Errorf("ModelID = %q, want %q", loaded.ModelID, original.ModelID)
	}

	// Check GCP
	if loaded.GCP.Enabled != original.GCP.Enabled {
		t.Errorf("GCP.Enabled = %v, want %v", loaded.GCP.Enabled, original.GCP.Enabled)
	}
	if loaded.GCP.ProjectID != original.GCP.ProjectID {
		t.Errorf("GCP.ProjectID = %q, want %q", loaded.GCP.ProjectID, original.GCP.ProjectID)
	}
	if loaded.GCP.ProjectNumber != original.GCP.ProjectNumber {
		t.Errorf("GCP.ProjectNumber = %q, want %q", loaded.GCP.ProjectNumber, original.GCP.ProjectNumber)
	}
	if loaded.GCP.Organization != original.GCP.Organization {
		t.Errorf("GCP.Organization = %q, want %q", loaded.GCP.Organization, original.GCP.Organization)
	}
	if loaded.GCP.Region != original.GCP.Region {
		t.Errorf("GCP.Region = %q, want %q", loaded.GCP.Region, original.GCP.Region)
	}

	// Check SWEs
	if len(loaded.SWEs) != len(original.SWEs) {
		t.Fatalf("len(SWEs) = %d, want %d", len(loaded.SWEs), len(original.SWEs))
	}
	for i, swe := range loaded.SWEs {
		orig := original.SWEs[i]
		if swe.Number != orig.Number {
			t.Errorf("SWEs[%d].Number = %d, want %d", i, swe.Number, orig.Number)
		}
		if swe.Title != orig.Title {
			t.Errorf("SWEs[%d].Title = %q, want %q", i, swe.Title, orig.Title)
		}
		if len(swe.Bullets) != len(orig.Bullets) {
			t.Errorf("SWEs[%d].Bullets len = %d, want %d", i, len(swe.Bullets), len(orig.Bullets))
		}
	}

	// Check SelectedSkills
	if len(loaded.SelectedSkills) != len(original.SelectedSkills) {
		t.Fatalf("len(SelectedSkills) = %d, want %d", len(loaded.SelectedSkills), len(original.SelectedSkills))
	}
	for name, want := range original.SelectedSkills {
		if loaded.SelectedSkills[name] != want {
			t.Errorf("SelectedSkills[%q] = %v, want %v", name, loaded.SelectedSkills[name], want)
		}
	}

	// Check CustomSkills
	if len(loaded.CustomSkills) != len(original.CustomSkills) {
		t.Fatalf("len(CustomSkills) = %d, want %d", len(loaded.CustomSkills), len(original.CustomSkills))
	}
	for i, cs := range loaded.CustomSkills {
		orig := original.CustomSkills[i]
		if cs.Name != orig.Name {
			t.Errorf("CustomSkills[%d].Name = %q, want %q", i, cs.Name, orig.Name)
		}
		if cs.Description != orig.Description {
			t.Errorf("CustomSkills[%d].Description = %q, want %q", i, cs.Description, orig.Description)
		}
	}

	// Check CustomAgents
	if len(loaded.CustomAgents) != len(original.CustomAgents) {
		t.Fatalf("len(CustomAgents) = %d, want %d", len(loaded.CustomAgents), len(original.CustomAgents))
	}
	for i, ca := range loaded.CustomAgents {
		orig := original.CustomAgents[i]
		if ca.Name != orig.Name {
			t.Errorf("CustomAgents[%d].Name = %q, want %q", i, ca.Name, orig.Name)
		}
		if ca.Title != orig.Title {
			t.Errorf("CustomAgents[%d].Title = %q, want %q", i, ca.Title, orig.Title)
		}
		if ca.Description != orig.Description {
			t.Errorf("CustomAgents[%d].Description = %q, want %q", i, ca.Description, orig.Description)
		}
		if len(ca.Instructions) != len(orig.Instructions) {
			t.Errorf("CustomAgents[%d].Instructions len = %d, want %d", i, len(ca.Instructions), len(orig.Instructions))
		} else {
			for j, inst := range ca.Instructions {
				if inst != orig.Instructions[j] {
					t.Errorf("CustomAgents[%d].Instructions[%d] = %q, want %q", i, j, inst, orig.Instructions[j])
				}
			}
		}
	}

	// Check Conventions
	if len(loaded.Conventions) != len(original.Conventions) {
		t.Fatalf("len(Conventions) = %d, want %d", len(loaded.Conventions), len(original.Conventions))
	}
	for i, c := range loaded.Conventions {
		if c != original.Conventions[i] {
			t.Errorf("Conventions[%d] = %q, want %q", i, c, original.Conventions[i])
		}
	}
}

func TestSaveSettingsCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	appteamDir := filepath.Join(dir, ".appteam")

	// Verify .appteam does not exist yet
	if _, err := os.Stat(appteamDir); err == nil {
		t.Fatal(".appteam directory should not exist yet")
	}

	cfg := &ProjectConfig{ProjectName: "test"}
	if err := SaveSettings(dir, cfg); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	info, err := os.Stat(appteamDir)
	if err != nil {
		t.Fatalf(".appteam directory was not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal(".appteam should be a directory")
	}
}

func TestLoadSettingsNotExist(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadSettings(dir)
	if err == nil {
		t.Fatal("LoadSettings should return an error when file does not exist")
	}
	if !strings.Contains(err.Error(), "reading settings") {
		t.Errorf("error should mention 'reading settings', got: %v", err)
	}
}

func TestLoadSettingsInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	appteamDir := filepath.Join(dir, ".appteam")
	if err := os.MkdirAll(appteamDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(filepath.Join(appteamDir, "settings.json"), []byte("{invalid"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := LoadSettings(dir)
	if err == nil {
		t.Fatal("LoadSettings should return an error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "parsing settings") {
		t.Errorf("error should mention 'parsing settings', got: %v", err)
	}
}

func TestLoadSettingsBackwardCompatNilSkills(t *testing.T) {
	dir := t.TempDir()
	appteamDir := filepath.Join(dir, ".appteam")
	if err := os.MkdirAll(appteamDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	// Pre-v0.7.0 settings.json without SelectedSkills
	oldJSON := `{"ProjectName":"legacy","OwnerName":"Test"}`
	if err := os.WriteFile(filepath.Join(appteamDir, "settings.json"), []byte(oldJSON), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := LoadSettings(dir)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if cfg.ProjectName != "legacy" {
		t.Errorf("ProjectName = %q, want %q", cfg.ProjectName, "legacy")
	}
	// nil SelectedSkills is valid — generator treats nil as all selected
	if cfg.SelectedSkills != nil {
		t.Errorf("SelectedSkills should be nil for pre-v0.7.0 settings, got %v", cfg.SelectedSkills)
	}
}

func TestLoadSettingsBackfillSkillDefaults(t *testing.T) {
	dir := t.TempDir()
	appteamDir := filepath.Join(dir, ".appteam")
	if err := os.MkdirAll(appteamDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	// Pre-v0.8.0 settings.json with only 6 PM skills
	oldJSON := `{"ProjectName":"legacy-v07","SelectedSkills":{"spec":true,"release":true,"pipeline":true,"status":true,"regenerate":false,"roadmap":true}}`
	if err := os.WriteFile(filepath.Join(appteamDir, "settings.json"), []byte(oldJSON), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := LoadSettings(dir)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}

	// Original PM skills should be preserved
	if !cfg.SelectedSkills["spec"] {
		t.Error("spec should be true")
	}
	if cfg.SelectedSkills["regenerate"] {
		t.Error("regenerate should be false (was explicitly set)")
	}

	// Tier 1 dev skills should be backfilled as true
	tier1 := []string{"debug", "test", "review", "docs", "refactor", "hotfix"}
	for _, name := range tier1 {
		if !cfg.SelectedSkills[name] {
			t.Errorf("SelectedSkills[%q] should be true (Tier 1 backfill), got false", name)
		}
	}

	// Tier 2 dev skills should be backfilled as false
	tier2 := []string{"api-design", "schema", "deploy", "security", "adr", "standup"}
	for _, name := range tier2 {
		if cfg.SelectedSkills[name] {
			t.Errorf("SelectedSkills[%q] should be false (Tier 2 backfill), got true", name)
		}
	}

	// brainstorm should be backfilled as true (PM skill)
	if !cfg.SelectedSkills["brainstorm"] {
		t.Error("SelectedSkills[\"brainstorm\"] should be true (PM backfill)")
	}

	// Total should be 21 keys (6 original PM + brainstorm + 14 dev)
	if len(cfg.SelectedSkills) != 21 {
		t.Errorf("len(SelectedSkills) = %d, want 21", len(cfg.SelectedSkills))
	}
}

func TestSavedSettingsJSONFormat(t *testing.T) {
	dir := t.TempDir()
	cfg := &ProjectConfig{
		ProjectName: "formattest",
		OwnerName:   "Test",
		SWEs:        []SWEConfig{{Number: 1, Title: "Dev"}},
	}
	if err := SaveSettings(dir, cfg); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	data, err := os.ReadFile(SettingsPath(dir))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	// Must end with newline
	if !strings.HasSuffix(string(data), "\n") {
		t.Error("settings file should end with a trailing newline")
	}

	// Must be valid JSON
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}

	// Must be indented (contains newline + spaces pattern)
	if !strings.Contains(string(data), "\n  ") {
		t.Error("settings file should be indented with 2 spaces")
	}

	// Check expected keys
	expectedKeys := []string{"ProjectName", "OwnerName", "SWEs"}
	for _, key := range expectedKeys {
		if _, ok := raw[key]; !ok {
			t.Errorf("expected key %q not found in JSON", key)
		}
	}
}

func TestTeamSizeBackfillDefault(t *testing.T) {
	dir := t.TempDir()
	appteamDir := filepath.Join(dir, ".appteam")
	if err := os.MkdirAll(appteamDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	// Pre-v0.12.0 settings.json without TeamSize
	oldJSON := `{"ProjectName":"legacy","OwnerName":"Test"}`
	if err := os.WriteFile(filepath.Join(appteamDir, "settings.json"), []byte(oldJSON), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg, err := LoadSettings(dir)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if cfg.TeamSize != "standard" {
		t.Errorf("TeamSize = %q, want %q (backfill default)", cfg.TeamSize, "standard")
	}
}

func TestTeamSizeRoundTrip(t *testing.T) {
	for _, size := range []string{"lean", "standard", "full"} {
		t.Run(size, func(t *testing.T) {
			dir := t.TempDir()
			original := &ProjectConfig{
				ProjectName: "test-" + size,
				TeamSize:    size,
			}
			if err := SaveSettings(dir, original); err != nil {
				t.Fatalf("SaveSettings: %v", err)
			}
			loaded, err := LoadSettings(dir)
			if err != nil {
				t.Fatalf("LoadSettings: %v", err)
			}
			if loaded.TeamSize != size {
				t.Errorf("TeamSize = %q, want %q", loaded.TeamSize, size)
			}
		})
	}
}
