package generator

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ahafin/appteam/internal/config"
	"github.com/ahafin/appteam/internal/templates"
	"github.com/ahafin/appteam/internal/wizard"
)

// generateUUID creates a v4 UUID string using crypto/rand.
func generateUUID() (string, error) {
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		return "", fmt.Errorf("generating UUID: %w", err)
	}
	// Set version 4 (bits 4-7 of byte 6)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant 10 (bits 6-7 of byte 8)
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

type fileSpec struct {
	path     string
	template string
	data     any
}

// Generate creates all output files based on the project configuration.
func Generate(cfg *config.ProjectConfig, color bool) error {
	s := wizard.NewStyler(color)
	base := cfg.TargetDir
	docsDir := filepath.Join(base, "docs")
	specsDir := filepath.Join(docsDir, "specs")

	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		return fmt.Errorf("creating specs directory: %w", err)
	}

	// Common files — generated for both frameworks
	files := []fileSpec{
		{filepath.Join(base, "CLAUDE.md"), templates.ClaudeMDTemplate, cfg},
		{filepath.Join(docsDir, "BACKLOG.md"), templates.BacklogTemplate, cfg},
		{filepath.Join(docsDir, "PROGRESS.md"), templates.ProgressTemplate, cfg},
		{filepath.Join(docsDir, "RELEASENOTES.md"), templates.ReleaseNotesTemplate, cfg},
		{filepath.Join(docsDir, "PIPELINE.md"), templates.PipelineTemplate, cfg},
		{filepath.Join(specsDir, "TEMPLATE.md"), templates.SpecTemplateFile, cfg},
	}

	// Framework-specific files
	if cfg.Framework == "scion" {
		scionFiles, err := generateScionFiles(cfg, s, base)
		if err != nil {
			return err
		}
		files = append(files, scionFiles...)
	} else {
		claudeFiles, err := generateClaudeCodeFiles(cfg, s, base)
		if err != nil {
			return err
		}
		files = append(files, claudeFiles...)
	}

	// Git repository initialization
	if err := setupGitRepo(cfg, s, base); err != nil {
		fmt.Printf("  %s git setup: %v\n", s.Dim("⚠"), err)
	}

	fmt.Printf("\n%s\n", s.Bold("Generating files..."))
	for _, f := range files {
		content, err := templates.Render(filepath.Base(f.path), f.template, f.data)
		if err != nil {
			return fmt.Errorf("rendering %s: %w", f.path, err)
		}
		if err := os.WriteFile(f.path, []byte(content), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", f.path, err)
		}
		fmt.Printf("  %s %s\n", s.Green("✓"), f.path)
	}

	// Save settings for future non-interactive regeneration
	if err := config.SaveSettings(base, cfg); err != nil {
		return fmt.Errorf("saving settings: %w", err)
	}
	fmt.Printf("  %s %s\n", s.Green("✓"), config.SettingsPath(base))

	fmt.Printf("\n%s Generated %d files.\n", s.BoldGreen("Done!"), len(files))
	return nil
}

// generateClaudeCodeFiles creates .claude/agents/ and .claude/skills/ files.
func generateClaudeCodeFiles(cfg *config.ProjectConfig, s *wizard.Styler, base string) ([]fileSpec, error) {
	agentsDir := filepath.Join(base, ".claude", "agents")
	skillsDir := filepath.Join(base, ".claude", "skills")

	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating agents directory: %w", err)
	}
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating skills directory: %w", err)
	}

	var files []fileSpec

	// PM always generated
	files = append(files, fileSpec{filepath.Join(agentsDir, "pm.md"), templates.PMTemplate, cfg})

	// TPM — not generated for lean teams
	if cfg.TeamSize != "lean" {
		files = append(files, fileSpec{filepath.Join(agentsDir, "tpm.md"), templates.TPMTemplate, cfg})
	}

	// SWE agents — lean teams only get SWE-1 regardless of config
	if cfg.TeamSize == "lean" {
		swe := cfg.SWEs[0]
		data := templates.SWETemplateData{
			ProjectName: cfg.ProjectName,
			OwnerName:   cfg.OwnerName,
			OwnerEmail:  cfg.OwnerEmail,
			ModelName:   cfg.ModelName,
			Number:      swe.Number,
			Title:       swe.Title,
			Bullets:     swe.Bullets,
			TeamSize:    cfg.TeamSize,
		}
		files = append(files, fileSpec{
			path:     filepath.Join(agentsDir, "swe-1.md"),
			template: templates.SWETemplate,
			data:     data,
		})
	} else {
		for _, swe := range cfg.SWEs {
			data := templates.SWETemplateData{
				ProjectName: cfg.ProjectName,
				OwnerName:   cfg.OwnerName,
				OwnerEmail:  cfg.OwnerEmail,
				ModelName:   cfg.ModelName,
				Number:      swe.Number,
				Title:       swe.Title,
				Bullets:     swe.Bullets,
				TeamSize:    cfg.TeamSize,
			}
			files = append(files, fileSpec{
				path:     filepath.Join(agentsDir, fmt.Sprintf("swe-%d.md", swe.Number)),
				template: templates.SWETemplate,
				data:     data,
			})
		}
	}

	// Optional agents — lean skips all, full enables all, standard uses config
	if cfg.TeamSize == "lean" {
		files = append(files, fileSpec{filepath.Join(agentsDir, "swe-test.md"), templates.SWETestTemplate, cfg})
	} else if cfg.TeamSize == "full" {
		files = append(files, fileSpec{filepath.Join(agentsDir, "swe-test.md"), templates.SWETestTemplate, cfg})
		files = append(files, fileSpec{filepath.Join(agentsDir, "swe-qa.md"), templates.SWEQATemplate, cfg})
		files = append(files, fileSpec{filepath.Join(agentsDir, "platform.md"), templates.PlatformTemplate, cfg})
		files = append(files, fileSpec{filepath.Join(agentsDir, "reviewer.md"), templates.ReviewerTemplate, cfg})
	} else {
		if cfg.IncludeSWETest {
			files = append(files, fileSpec{filepath.Join(agentsDir, "swe-test.md"), templates.SWETestTemplate, cfg})
		}
		if cfg.IncludeSWEQA {
			files = append(files, fileSpec{filepath.Join(agentsDir, "swe-qa.md"), templates.SWEQATemplate, cfg})
		}
		if cfg.IncludePlatform {
			files = append(files, fileSpec{filepath.Join(agentsDir, "platform.md"), templates.PlatformTemplate, cfg})
		}
		if cfg.IncludeReviewer {
			files = append(files, fileSpec{filepath.Join(agentsDir, "reviewer.md"), templates.ReviewerTemplate, cfg})
		}
	}

	// Custom agents
	for _, ca := range cfg.CustomAgents {
		data := templates.CustomAgentTemplateData{
			Name:         ca.Name,
			Title:        ca.Title,
			Description:  ca.Description,
			Instructions: ca.Instructions,
			ProjectName:  cfg.ProjectName,
			OwnerName:    cfg.OwnerName,
			OwnerEmail:   cfg.OwnerEmail,
			ModelName:    cfg.ModelName,
			TeamSize:     cfg.TeamSize,
		}
		files = append(files, fileSpec{
			path:     filepath.Join(agentsDir, ca.Name+".md"),
			template: templates.CustomAgentTemplate,
			data:     data,
		})
	}

	// Skills — conditional generation based on SelectedSkills.
	skillFiles := []struct {
		name     string
		path     string
		template string
	}{
		{"spec", filepath.Join(skillsDir, "spec.md"), templates.SkillSpecTemplate},
		{"release", filepath.Join(skillsDir, "release.md"), templates.SkillReleaseTemplate},
		{"pipeline", filepath.Join(skillsDir, "pipeline.md"), templates.SkillPipelineTemplate},
		{"status", filepath.Join(skillsDir, "status.md"), templates.SkillStatusTemplate},
		{"regenerate", filepath.Join(skillsDir, "regenerate.md"), templates.SkillRegenerateTemplate},
		{"roadmap", filepath.Join(skillsDir, "roadmap.md"), templates.SkillRoadmapTemplate},
		{"brainstorm", filepath.Join(skillsDir, "brainstorm.md"), templates.SkillBrainstormTemplate},
		{"debug", filepath.Join(skillsDir, "debug.md"), templates.SkillDebugTemplate},
		{"test", filepath.Join(skillsDir, "test.md"), templates.SkillTestWriteTemplate},
		{"review", filepath.Join(skillsDir, "review.md"), templates.SkillReviewTemplate},
		{"docs", filepath.Join(skillsDir, "docs.md"), templates.SkillDocsTemplate},
		{"refactor", filepath.Join(skillsDir, "refactor.md"), templates.SkillRefactorTemplate},
		{"hotfix", filepath.Join(skillsDir, "hotfix.md"), templates.SkillHotfixTemplate},
		{"api-design", filepath.Join(skillsDir, "api-design.md"), templates.SkillAPIDesignTemplate},
		{"schema", filepath.Join(skillsDir, "schema.md"), templates.SkillSchemaTemplate},
		{"deploy", filepath.Join(skillsDir, "deploy.md"), templates.SkillDeployTemplate},
		{"security", filepath.Join(skillsDir, "security.md"), templates.SkillSecurityTemplate},
		{"adr", filepath.Join(skillsDir, "adr.md"), templates.SkillADRTemplate},
		{"standup", filepath.Join(skillsDir, "standup.md"), templates.SkillStandupTemplate},
		{"cuj-list", filepath.Join(skillsDir, "cuj-list.md"), templates.SkillCUJListTemplate},
		{"cuj-test", filepath.Join(skillsDir, "cuj-test.md"), templates.SkillCUJTestTemplate},
	}
	allSkills := len(cfg.SelectedSkills) == 0
	for _, sf := range skillFiles {
		if allSkills || cfg.SelectedSkills[sf.name] {
			files = append(files, fileSpec{sf.path, sf.template, cfg})
		}
	}

	// Custom skills
	for _, cs := range cfg.CustomSkills {
		tmpl := fmt.Sprintf("# /%s — %s\n\n## Trigger\n\nUser invokes `/%s`.\n\n## Instructions\n\n%s\n\n## Project Context\n\n- **Project:** %s\n- **Owner:** %s (%s)\n", cs.Name, cs.Description, cs.Name, cs.Description, cfg.ProjectName, cfg.OwnerName, cfg.OwnerEmail)
		path := filepath.Join(skillsDir, cs.Name+".md")
		if err := os.WriteFile(path, []byte(tmpl), 0o644); err != nil {
			return nil, fmt.Errorf("writing custom skill %s: %w", cs.Name, err)
		}
		fmt.Printf("  %s %s\n", s.Green("✓"), path)
	}

	return files, nil
}

// Skill-to-role mapping for Scion agents.md embedding.
var (
	pmSkillNames  = []string{"spec", "release", "roadmap", "status", "pipeline", "regenerate", "brainstorm"}
	devSkillNames = []string{"debug", "test", "review", "docs", "refactor", "hotfix", "api-design", "schema", "deploy", "security", "adr", "standup"}
	qaSkillNames  = []string{"cuj-list", "cuj-test"}
)

// skillTemplateMap maps skill names to their template constants.
var skillTemplateMap = map[string]string{
	"spec":       templates.SkillSpecTemplate,
	"release":    templates.SkillReleaseTemplate,
	"roadmap":    templates.SkillRoadmapTemplate,
	"status":     templates.SkillStatusTemplate,
	"pipeline":   templates.SkillPipelineTemplate,
	"regenerate": templates.SkillRegenerateTemplate,
	"brainstorm": templates.SkillBrainstormTemplate,
	"debug":      templates.SkillDebugTemplate,
	"test":       templates.SkillTestWriteTemplate,
	"review":     templates.SkillReviewTemplate,
	"docs":       templates.SkillDocsTemplate,
	"refactor":   templates.SkillRefactorTemplate,
	"hotfix":     templates.SkillHotfixTemplate,
	"api-design": templates.SkillAPIDesignTemplate,
	"schema":     templates.SkillSchemaTemplate,
	"deploy":     templates.SkillDeployTemplate,
	"security":   templates.SkillSecurityTemplate,
	"adr":        templates.SkillADRTemplate,
	"standup":    templates.SkillStandupTemplate,
	"cuj-list":   templates.SkillCUJListTemplate,
	"cuj-test":   templates.SkillCUJTestTemplate,
}

// renderSkillsForRole renders and concatenates skill templates for a given role.
// When includeCustom is true, custom skills are also appended (for SWE agents).
func renderSkillsForRole(cfg *config.ProjectConfig, skillNames []string, includeCustom bool) string {
	allSkills := len(cfg.SelectedSkills) == 0
	var parts []string
	for _, name := range skillNames {
		if allSkills || cfg.SelectedSkills[name] {
			tmpl, ok := skillTemplateMap[name]
			if !ok {
				continue
			}
			content, err := templates.Render(name+".md", tmpl, cfg)
			if err != nil {
				continue
			}
			parts = append(parts, content)
		}
	}
	if includeCustom {
		for _, cs := range cfg.CustomSkills {
			tmpl := fmt.Sprintf("# /%s — %s\n\n## Trigger\n\nUser invokes `/%s`.\n\n## Instructions\n\n%s\n\n## Project Context\n\n- **Project:** %s\n- **Owner:** %s (%s)\n", cs.Name, cs.Description, cs.Name, cs.Description, cfg.ProjectName, cfg.OwnerName, cfg.OwnerEmail)
			parts = append(parts, tmpl)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "\n---\n\n## Embedded Skills\n\n" + strings.Join(parts, "\n---\n\n")
}

// scionRole describes a single Scion agent role for template generation.
type scionRole struct {
	dirName       string // directory name under .scion/templates/
	agentName     string // name field in scion-agent.yaml
	description   string // description field in scion-agent.yaml
	sysTmpl       string // system prompt template constant
	agentsTmpl    string // agents.md template constant
	data          any    // template data
	skillCategory string // "pm", "dev", "qa", or ""
}

// generateScionFiles creates .scion/templates/<role>/ files.
func generateScionFiles(cfg *config.ProjectConfig, s *wizard.Styler, base string) ([]fileSpec, error) {
	scionBase := filepath.Join(base, ".scion", "templates")

	// Build the list of roles based on team size
	var roles []scionRole

	// PM — always present
	roles = append(roles, scionRole{
		dirName:       "pm",
		agentName:     "pm",
		description:   "Product Manager — translates PO feedback into specs and backlog items",
		sysTmpl:       templates.ScionPMSystemPrompt,
		agentsTmpl:    templates.ScionPMAgentsMD,
		data:          cfg,
		skillCategory: "pm",
	})

	// TPM — not generated for lean teams
	if cfg.TeamSize != "lean" {
		roles = append(roles, scionRole{
			dirName:     "tpm",
			agentName:   "tpm",
			description: "Technical Program Manager — coordinates agents, manages backlog and progress",
			sysTmpl:     templates.ScionTPMSystemPrompt,
			agentsTmpl:  templates.ScionTPMAgentsMD,
			data:        cfg,
		})
	}

	// SWE agents
	if cfg.TeamSize == "lean" {
		swe := cfg.SWEs[0]
		data := templates.SWETemplateData{
			ProjectName: cfg.ProjectName,
			OwnerName:   cfg.OwnerName,
			OwnerEmail:  cfg.OwnerEmail,
			ModelName:   cfg.ModelName,
			Number:      swe.Number,
			Title:       swe.Title,
			Bullets:     swe.Bullets,
			TeamSize:    cfg.TeamSize,
		}
		roles = append(roles, scionRole{
			dirName:       "swe-1",
			agentName:     "swe-1",
			description:   fmt.Sprintf("Software Engineer 1 — %s", swe.Title),
			sysTmpl:       templates.ScionSWESystemPrompt,
			agentsTmpl:    templates.ScionSWEAgentsMD,
			data:          data,
			skillCategory: "dev",
		})
	} else {
		for _, swe := range cfg.SWEs {
			data := templates.SWETemplateData{
				ProjectName: cfg.ProjectName,
				OwnerName:   cfg.OwnerName,
				OwnerEmail:  cfg.OwnerEmail,
				ModelName:   cfg.ModelName,
				Number:      swe.Number,
				Title:       swe.Title,
				Bullets:     swe.Bullets,
				TeamSize:    cfg.TeamSize,
			}
			name := fmt.Sprintf("swe-%d", swe.Number)
			roles = append(roles, scionRole{
				dirName:       name,
				agentName:     name,
				description:   fmt.Sprintf("Software Engineer %d — %s", swe.Number, swe.Title),
				sysTmpl:       templates.ScionSWESystemPrompt,
				agentsTmpl:    templates.ScionSWEAgentsMD,
				data:          data,
				skillCategory: "dev",
			})
		}
	}

	// Optional agents — same inclusion rules as Claude Code
	if cfg.TeamSize == "lean" {
		roles = append(roles, scionRole{
			dirName:       "swe-test",
			agentName:     "swe-test",
			description:   "Test Engineer — automated testing and quality assurance",
			sysTmpl:       templates.ScionSWETestSystemPrompt,
			agentsTmpl:    templates.ScionSWETestAgentsMD,
			data:          cfg,
			skillCategory: "qa",
		})
	} else if cfg.TeamSize == "full" {
		roles = append(roles, scionRole{
			dirName:       "swe-test",
			agentName:     "swe-test",
			description:   "Test Engineer — automated testing and quality assurance",
			sysTmpl:       templates.ScionSWETestSystemPrompt,
			agentsTmpl:    templates.ScionSWETestAgentsMD,
			data:          cfg,
			skillCategory: "qa",
		})
		roles = append(roles, scionRole{
			dirName:     "swe-qa",
			agentName:   "swe-qa",
			description: "QA Engineer — CUJ-oriented testing and browser validation",
			sysTmpl:     templates.ScionSWEQASystemPrompt,
			agentsTmpl:  templates.ScionSWEQAAgentsMD,
			data:        cfg,
		})
		roles = append(roles, scionRole{
			dirName:     "platform",
			agentName:   "platform",
			description: "Platform Engineer — infrastructure, deployment, and reliability",
			sysTmpl:     templates.ScionPlatformSystemPrompt,
			agentsTmpl:  templates.ScionPlatformAgentsMD,
			data:        cfg,
		})
		roles = append(roles, scionRole{
			dirName:     "reviewer",
			agentName:   "reviewer",
			description: "Code Reviewer — quality, security, and performance reviews",
			sysTmpl:     templates.ScionReviewerSystemPrompt,
			agentsTmpl:  templates.ScionReviewerAgentsMD,
			data:        cfg,
		})
	} else {
		if cfg.IncludeSWETest {
			roles = append(roles, scionRole{
				dirName:       "swe-test",
				agentName:     "swe-test",
				description:   "Test Engineer — automated testing and quality assurance",
				sysTmpl:       templates.ScionSWETestSystemPrompt,
				agentsTmpl:    templates.ScionSWETestAgentsMD,
				data:          cfg,
				skillCategory: "qa",
			})
		}
		if cfg.IncludeSWEQA {
			roles = append(roles, scionRole{
				dirName:     "swe-qa",
				agentName:   "swe-qa",
				description: "QA Engineer — CUJ-oriented testing and browser validation",
				sysTmpl:     templates.ScionSWEQASystemPrompt,
				agentsTmpl:  templates.ScionSWEQAAgentsMD,
				data:        cfg,
			})
		}
		if cfg.IncludePlatform {
			roles = append(roles, scionRole{
				dirName:     "platform",
				agentName:   "platform",
				description: "Platform Engineer — infrastructure, deployment, and reliability",
				sysTmpl:     templates.ScionPlatformSystemPrompt,
				agentsTmpl:  templates.ScionPlatformAgentsMD,
				data:        cfg,
			})
		}
		if cfg.IncludeReviewer {
			roles = append(roles, scionRole{
				dirName:     "reviewer",
				agentName:   "reviewer",
				description: "Code Reviewer — quality, security, and performance reviews",
				sysTmpl:     templates.ScionReviewerSystemPrompt,
				agentsTmpl:  templates.ScionReviewerAgentsMD,
				data:        cfg,
			})
		}
	}

	// Custom agents
	for _, ca := range cfg.CustomAgents {
		data := templates.CustomAgentTemplateData{
			Name:         ca.Name,
			Title:        ca.Title,
			Description:  ca.Description,
			Instructions: ca.Instructions,
			ProjectName:  cfg.ProjectName,
			OwnerName:    cfg.OwnerName,
			OwnerEmail:   cfg.OwnerEmail,
			ModelName:    cfg.ModelName,
			TeamSize:     cfg.TeamSize,
		}
		roles = append(roles, scionRole{
			dirName:     ca.Name,
			agentName:   ca.Name,
			description: fmt.Sprintf("%s — %s", ca.Title, ca.Description),
			sysTmpl:     templates.ScionCustomAgentSystemPrompt,
			agentsTmpl:  templates.ScionCustomAgentAgentsMD,
			data:        data,
		})
	}

	// Ensure .scion/ is in .gitignore (Scion requires this)
	gitignorePath := filepath.Join(base, ".gitignore")
	if err := ensureGitignoreEntry(gitignorePath, ".scion/agents/"); err != nil {
		fmt.Printf("  %s .gitignore update: %v\n", s.Dim("⚠"), err)
	}

	// Create .scion/agents/ directory (empty, for runtime agent state)
	scionAgentsDir := filepath.Join(base, ".scion", "agents")
	if err := os.MkdirAll(scionAgentsDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating scion agents directory: %w", err)
	}

	// Create .scion/grove-id with a UUID v4
	groveID, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("generating grove-id: %w", err)
	}
	groveIDPath := filepath.Join(base, ".scion", "grove-id")
	if err := os.WriteFile(groveIDPath, []byte(groveID+"\n"), 0o644); err != nil {
		return nil, fmt.Errorf("writing grove-id: %w", err)
	}
	fmt.Printf("  %s %s\n", s.Green("✓"), groveIDPath)

	// Create directories and generate files
	var files []fileSpec
	for _, role := range roles {
		roleDir := filepath.Join(scionBase, role.dirName)
		if err := os.MkdirAll(roleDir, 0o755); err != nil {
			return nil, fmt.Errorf("creating scion template directory %s: %w", role.dirName, err)
		}

		// scion-agent.yaml
		yamlData := templates.ScionAgentYAMLData{
			Name:        role.agentName,
			Description: role.description,
			Harness:     cfg.DefaultHarness,
		}
		files = append(files, fileSpec{
			path:     filepath.Join(roleDir, "scion-agent.yaml"),
			template: templates.ScionAgentYAMLTemplate,
			data:     yamlData,
		})

		// system-prompt.md
		files = append(files, fileSpec{
			path:     filepath.Join(roleDir, "system-prompt.md"),
			template: role.sysTmpl,
			data:     role.data,
		})

		// agents.md — render base template then append skills
		agentsMDContent, err := templates.Render(role.dirName+"-agents.md", role.agentsTmpl, role.data)
		if err != nil {
			return nil, fmt.Errorf("rendering agents.md for %s: %w", role.dirName, err)
		}

		// Embed skills based on role category
		var skillContent string
		switch role.skillCategory {
		case "pm":
			skillContent = renderSkillsForRole(cfg, pmSkillNames, false)
		case "dev":
			skillContent = renderSkillsForRole(cfg, devSkillNames, true)
		case "qa":
			skillContent = renderSkillsForRole(cfg, qaSkillNames, false)
		}

		fullContent := agentsMDContent + skillContent
		agentsMDPath := filepath.Join(roleDir, "agents.md")
		if err := os.WriteFile(agentsMDPath, []byte(fullContent), 0o644); err != nil {
			return nil, fmt.Errorf("writing agents.md for %s: %w", role.dirName, err)
		}
		fmt.Printf("  %s %s\n", s.Green("✓"), agentsMDPath)
	}

	return files, nil
}

// ensureGitignoreEntry appends entry to .gitignore if not already present.
func ensureGitignoreEntry(gitignorePath, entry string) error {
	data, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	content := string(data)
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == entry {
			return nil // already present
		}
	}
	// Append with a preceding newline if file doesn't end with one
	prefix := ""
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		prefix = "\n"
	}
	return os.WriteFile(gitignorePath, []byte(content+prefix+entry+"\n"), 0o644)
}

// setupGitRepo handles git init, remote add, and gh repo create based on config.
func setupGitRepo(cfg *config.ProjectConfig, s *wizard.Styler, base string) error {
	gitDir := filepath.Join(base, ".git")
	hasGit := false
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		hasGit = true
	}

	if cfg.InitGit && !hasGit {
		fmt.Printf("\n%s\n", s.Bold("Initializing git repository..."))
		cmd := exec.Command("git", "init")
		cmd.Dir = base
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git init: %w\n    %s", err, strings.TrimSpace(string(out)))
		}
		fmt.Printf("  %s git init\n", s.Green("✓"))
	}

	if cfg.CreateRepo {
		// Extract org/name from RepoURL: https://github.com/org/name.git
		repoSlug := cfg.RepoURL
		repoSlug = strings.TrimPrefix(repoSlug, "https://github.com/")
		repoSlug = strings.TrimSuffix(repoSlug, ".git")

		if _, err := exec.LookPath("gh"); err != nil {
			fmt.Printf("  %s gh CLI not found — skipping repo creation\n", s.Dim("⚠"))
		} else {
			fmt.Printf("%s\n", s.Bold("Creating GitHub repository..."))
			cmd := exec.Command("gh", "repo", "create", repoSlug, "--public", "--source=.", "--remote=origin")
			cmd.Dir = base
			if out, err := cmd.CombinedOutput(); err != nil {
				fmt.Printf("  %s gh repo create: %v\n    %s\n", s.Dim("⚠"), err, strings.TrimSpace(string(out)))
			} else {
				fmt.Printf("  %s gh repo create %s\n", s.Green("✓"), repoSlug)
			}
		}
	} else if cfg.RepoURL != "" {
		// Add remote if URL was provided but repo wasn't just created (gh repo create adds remote)
		cmd := exec.Command("git", "remote", "add", "origin", cfg.RepoURL)
		cmd.Dir = base
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("  %s git remote add: %v\n    %s\n", s.Dim("⚠"), err, strings.TrimSpace(string(out)))
		} else {
			fmt.Printf("  %s git remote add origin %s\n", s.Green("✓"), cfg.RepoURL)
		}
	}

	return nil
}
