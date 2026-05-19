package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const settingsDir = ".appteam"
const settingsFile = "settings.json"

// SettingsPath returns the full path to settings.json within the given directory.
func SettingsPath(dir string) string {
	return filepath.Join(dir, settingsDir, settingsFile)
}

// SettingsExist returns true if .appteam/settings.json exists in dir.
func SettingsExist(dir string) bool {
	_, err := os.Stat(SettingsPath(dir))
	return err == nil
}

// LoadSettings reads and unmarshals .appteam/settings.json from dir.
// For backward compatibility, if SelectedSkills exists but has fewer than 21 keys,
// missing PM and Tier 1 dev skills are backfilled as true (on by default).
func LoadSettings(dir string) (*ProjectConfig, error) {
	data, err := os.ReadFile(SettingsPath(dir))
	if err != nil {
		return nil, fmt.Errorf("reading settings: %w", err)
	}
	var cfg ProjectConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing settings: %w", err)
	}
	backfillSkillDefaults(&cfg)
	backfillTeamSize(&cfg)
	backfillFramework(&cfg)
	backfillDefaultHarness(&cfg)
	return &cfg, nil
}

// backfillSkillDefaults ensures that SelectedSkills has entries for all 21
// predefined skills when a non-nil map exists but is missing newer keys.
// PM and Tier 1 dev skills default to true; Tier 2 dev skills default to false.
func backfillSkillDefaults(cfg *ProjectConfig) {
	if cfg.SelectedSkills == nil {
		return // nil means "all skills" — backward compat for pre-v0.7.0
	}
	defaults := map[string]bool{
		// PM skills (default on)
		"spec": true, "release": true, "roadmap": true,
		"status": true, "pipeline": true, "regenerate": true, "brainstorm": true,
		// Dev Tier 1 (default on)
		"debug": true, "test": true, "review": true,
		"docs": true, "refactor": true, "hotfix": true,
		// Dev Tier 2 (default off)
		"api-design": false, "schema": false, "deploy": false,
		"security": false, "adr": false, "standup": false,
		"cuj-list": false, "cuj-test": false,
	}
	for name, def := range defaults {
		if _, exists := cfg.SelectedSkills[name]; !exists {
			cfg.SelectedSkills[name] = def
		}
	}
}

// backfillTeamSize defaults TeamSize to "standard" when empty (pre-v0.12.0 compat).
func backfillTeamSize(cfg *ProjectConfig) {
	if cfg.TeamSize == "" {
		cfg.TeamSize = "standard"
	}
}

// backfillFramework defaults Framework to "claude-code" when empty (pre-v0.15.0 compat).
func backfillFramework(cfg *ProjectConfig) {
	if cfg.Framework == "" {
		cfg.Framework = "claude-code"
	}
}

// backfillDefaultHarness defaults DefaultHarness to "claude" when empty (pre-v0.15.0 compat).
func backfillDefaultHarness(cfg *ProjectConfig) {
	if cfg.DefaultHarness == "" {
		cfg.DefaultHarness = "claude"
	}
}

// SaveSettings marshals cfg to JSON and writes it to <dir>/.appteam/settings.json.
// Creates the .appteam/ directory if it doesn't exist.
func SaveSettings(dir string, cfg *ProjectConfig) error {
	dirPath := filepath.Join(dir, settingsDir)
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("creating %s directory: %w", settingsDir, err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(SettingsPath(dir), data, 0o644); err != nil {
		return fmt.Errorf("writing settings: %w", err)
	}
	return nil
}
