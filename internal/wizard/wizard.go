package wizard

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ahafin/appteam/internal/config"
)

// Run executes the interactive wizard and returns the resulting configuration.
// If targetDir is non-empty, the target directory prompt is skipped and the provided value is used.
func Run(r io.Reader, w io.Writer, color bool, targetDir string) (*config.ProjectConfig, error) {
	scanner := bufio.NewScanner(r)
	cfg := &config.ProjectConfig{}
	s := NewStyler(color)
	totalSteps := 9

	// Welcome banner
	fmt.Fprintln(w, s.Banner())
	fmt.Fprintln(w)

	// Step 1: Project Basics
	fmt.Fprintln(w, s.StepHeader(1, totalSteps, "Project Basics"))
	cfg.ProjectName = ask(scanner, w, s, "Project name", "")
	cfg.Description = ask(scanner, w, s, "Short description", "")
	cfg.TechStack = ask(scanner, w, s, "Tech stack", "")
	if targetDir != "" {
		cfg.TargetDir = targetDir
		fmt.Fprintf(w, "  %s Target directory: %s %s\n", s.Green("▸"), targetDir, s.Dim("(from -d flag)"))
	} else {
		cfg.TargetDir = ask(scanner, w, s, "Target directory", ".")
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 2: Team Framework
	fmt.Fprintln(w, s.StepHeader(2, totalSteps, "Team Framework"))
	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s\n", s.BoldWhite("Which team framework?"))
	fmt.Fprintf(w, "    %s Claude Code Agent Teams  %s\n", s.Green("1)"), s.Dim("— .claude/agents/ + .claude/skills/"))
	fmt.Fprintf(w, "    %s Scion                    %s\n", s.Green("2)"), s.Dim("— .scion/templates/ with container-based agents"))
	fmt.Fprintln(w)
	frameworkChoice := askInt(scanner, w, s, "Framework", 1)
	if frameworkChoice == 2 {
		cfg.Framework = "scion"
		cfg.DefaultHarness = ask(scanner, w, s, "Default harness (claude/gemini/opencode/codex)", "claude")
	} else {
		cfg.Framework = "claude-code"
		cfg.DefaultHarness = "claude"
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 3: Git Repository
	fmt.Fprintln(w, s.StepHeader(3, totalSteps, "Git Repository"))
	targetForGit := cfg.TargetDir
	gitDir := filepath.Join(targetForGit, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		fmt.Fprintf(w, "  %s %s\n", s.Green("✓"), "Already a git repository")
		repoURL := ask(scanner, w, s, "Remote URL (leave blank to skip)", "")
		if repoURL != "" {
			cfg.RepoURL = repoURL
		}
	} else {
		cfg.InitGit = askBool(scanner, w, s, "Initialize git repository?", true)
		if cfg.InitGit {
			alreadyCreated := askBool(scanner, w, s, "GitHub repo already created?", false)
			if alreadyCreated {
				cfg.RepoURL = ask(scanner, w, s, "Repository URL", "")
			} else {
				cfg.GitHubOrg = ask(scanner, w, s, "GitHub owner (org or username)", "")
				repoName := ask(scanner, w, s, "Repository name", cfg.ProjectName)
				cfg.CreateRepo = true
				cfg.RepoURL = fmt.Sprintf("https://github.com/%s/%s.git", cfg.GitHubOrg, repoName)
			}
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 4: Product Owner
	fmt.Fprintln(w, s.StepHeader(4, totalSteps, "Product Owner"))
	cfg.OwnerName = ask(scanner, w, s, "Name", "Ameer Abbas")
	cfg.OwnerEmail = ask(scanner, w, s, "Email", "ameer00@gmail.com")
	cfg.OwnerGitHub = ask(scanner, w, s, "GitHub username", "ahafin")
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 5: GCP Configuration
	fmt.Fprintln(w, s.StepHeader(5, totalSteps, "GCP Configuration"))
	cfg.GCP.Enabled = askBool(scanner, w, s, "Include GCP project details?", false)
	if cfg.GCP.Enabled {
		cfg.GCP.ProjectID = ask(scanner, w, s, "GCP Project ID", "")
		cfg.GCP.ProjectNumber = ask(scanner, w, s, "GCP Project Number", "")
		cfg.GCP.Organization = ask(scanner, w, s, "GCP Organization", "")
		cfg.GCP.Region = ask(scanner, w, s, "GCP Region", "us-west1")
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 6: Agent Team
	fmt.Fprintln(w, s.StepHeader(6, totalSteps, "Agent Team"))

	fmt.Fprintf(w, "  %s\n", s.BoldWhite("Team size:"))
	fmt.Fprintf(w, "    %s lean     %s\n", s.Green("1."), s.Dim("— PM + 1 SWE + SWE-Test"))
	fmt.Fprintf(w, "    %s standard %s\n", s.Green("2."), s.Dim("— PM + TPM + 2 SWEs + SWE-Test + Reviewer"))
	fmt.Fprintf(w, "    %s full     %s\n", s.Green("3."), s.Dim("— PM + TPM + 5 SWEs + all agents"))
	teamSizeChoice := askInt(scanner, w, s, "Select team size", 2)
	switch teamSizeChoice {
	case 1:
		cfg.TeamSize = "lean"
	case 3:
		cfg.TeamSize = "full"
	default:
		cfg.TeamSize = "standard"
	}

	var sweCount int
	switch cfg.TeamSize {
	case "lean":
		sweCount = 1
	case "full":
		sweCount = askInt(scanner, w, s, "Number of SWE agents (1-5)", 5)
		if sweCount < 1 {
			sweCount = 1
		}
		if sweCount > 5 {
			sweCount = 5
		}
	default: // standard
		sweCount = askInt(scanner, w, s, "Number of SWE agents (1-5)", 3)
		if sweCount < 1 {
			sweCount = 1
		}
		if sweCount > 5 {
			sweCount = 5
		}
	}

	cfg.SWEs = make([]config.SWEConfig, sweCount)
	for i := 0; i < sweCount; i++ {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  %s\n", s.BoldWhite(fmt.Sprintf("SWE-%d", i+1)))
		cfg.SWEs[i].Number = i + 1
		cfg.SWEs[i].Title = ask(scanner, w, s, "Title/specialty", fmt.Sprintf("General Engineer %d", i+1))

		fmt.Fprintf(w, "  %s Specialty bullets %s\n", s.Green("▸"), s.Dim("(blank line to finish)"))
		for {
			fmt.Fprintf(w, "    %s ", s.Dim("│"))
			if !scanner.Scan() {
				break
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				break
			}
			cfg.SWEs[i].Bullets = append(cfg.SWEs[i].Bullets, line)
		}
	}

	fmt.Fprintln(w)
	switch cfg.TeamSize {
	case "lean":
		cfg.IncludeSWETest = true
		cfg.IncludePlatform = false
		cfg.IncludeReviewer = false
		cfg.IncludeSWEQA = false
		fmt.Fprintf(w, "  %s SWE-Test auto-enabled for lean team\n", s.Green("✓"))
	case "full":
		cfg.IncludePlatform = true
		cfg.IncludeReviewer = true
		cfg.IncludeSWETest = true
		cfg.IncludeSWEQA = true
		fmt.Fprintf(w, "  %s All agents auto-enabled: Platform, Reviewer, SWE-Test, SWE-QA\n", s.Green("✓"))
	default: // standard
		cfg.IncludePlatform = askBool(scanner, w, s, "Include Platform Engineer?", false)
		cfg.IncludeReviewer = askBool(scanner, w, s, "Include Code Reviewer?", true)
		cfg.IncludeSWETest = askBool(scanner, w, s, "Include SWE-Test?", true)
		cfg.IncludeSWEQA = askBool(scanner, w, s, "Include SWE-QA?", true)
	}

	fmt.Fprintln(w)
	if askBool(scanner, w, s, "Add a custom agent?", false) {
		for {
			fmt.Fprintln(w)
			name := ask(scanner, w, s, "Agent name (kebab-case)", "")
			if name == "" {
				break
			}
			title := ask(scanner, w, s, "Title", "")
			desc := ask(scanner, w, s, "Description", "")
			fmt.Fprintf(w, "  %s Instructions %s\n", s.Green("▸"), s.Dim("(blank line to finish)"))
			var instructions []string
			for {
				fmt.Fprintf(w, "    %s ", s.Dim("│"))
				if !scanner.Scan() {
					break
				}
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					break
				}
				instructions = append(instructions, line)
			}
			cfg.CustomAgents = append(cfg.CustomAgents, config.CustomAgentConfig{
				Name:         name,
				Title:        title,
				Description:  desc,
				Instructions: instructions,
			})
			if !askBool(scanner, w, s, "Add another custom agent?", false) {
				break
			}
		}
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s\n", s.BoldWhite("Default model for agents:"))
	fmt.Fprintf(w, "    %s Opus 4.6 %s\n", s.Green("1."), s.Dim("(claude-opus-4-6)"))
	fmt.Fprintf(w, "    %s Sonnet 4.6 %s\n", s.Green("2."), s.Dim("(claude-sonnet-4-6)"))
	fmt.Fprintf(w, "    %s Haiku 4.5 %s\n", s.Green("3."), s.Dim("(claude-haiku-4-5-20251001)"))
	modelChoice := askInt(scanner, w, s, "Select model", 1)
	switch modelChoice {
	case 2:
		cfg.ModelName = "Sonnet 4.6"
		cfg.ModelID = "claude-sonnet-4-6"
	case 3:
		cfg.ModelName = "Haiku 4.5"
		cfg.ModelID = "claude-haiku-4-5-20251001"
	default:
		cfg.ModelName = "Opus 4.6"
		cfg.ModelID = "claude-opus-4-6"
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 7: Conventions
	fmt.Fprintln(w, s.StepHeader(7, totalSteps, "Conventions"))
	fmt.Fprintf(w, "  %s App-specific conventions %s\n", s.Green("▸"), s.Dim("(blank line to finish)"))
	for {
		fmt.Fprintf(w, "    %s ", s.Dim("│"))
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		cfg.Conventions = append(cfg.Conventions, line)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 8: Skills
	fmt.Fprintln(w, s.StepHeader(8, totalSteps, "Skills"))

	fmt.Fprintf(w, "\n  %s\n", s.BoldWhite("Project Management Skills:"))
	cfg.SelectedSkills = map[string]bool{
		"spec":       askBool(scanner, w, s, "/spec       — Create product specs", true),
		"release":    askBool(scanner, w, s, "/release    — Release notes, commit, tag, push", true),
		"roadmap":    askBool(scanner, w, s, "/roadmap    — Add items to backlog", true),
		"status":     askBool(scanner, w, s, "/status     — Milestone status summary", true),
		"pipeline":   askBool(scanner, w, s, "/pipeline   — Spin up agent team", true),
		"regenerate": askBool(scanner, w, s, "/regenerate — Regenerate from settings", true),
		"brainstorm": askBool(scanner, w, s, "/brainstorm — Brainstorm product ideas with PO", true),
	}

	fmt.Fprintf(w, "\n  %s\n", s.BoldWhite("Development Skills:"))
	cfg.SelectedSkills["debug"] = askBool(scanner, w, s, "/debug      — Systematic debugging workflow", true)
	cfg.SelectedSkills["test"] = askBool(scanner, w, s, "/test       — Write tests for a module", true)
	cfg.SelectedSkills["review"] = askBool(scanner, w, s, "/review     — Code review checklist", true)
	cfg.SelectedSkills["docs"] = askBool(scanner, w, s, "/docs       — Generate documentation", true)
	cfg.SelectedSkills["refactor"] = askBool(scanner, w, s, "/refactor   — Safe refactoring workflow", true)
	cfg.SelectedSkills["hotfix"] = askBool(scanner, w, s, "/hotfix     — Emergency fix workflow", true)
	cfg.SelectedSkills["api-design"] = askBool(scanner, w, s, "/api-design — Design API endpoints", false)
	cfg.SelectedSkills["schema"] = askBool(scanner, w, s, "/schema     — Database schema design", false)
	cfg.SelectedSkills["deploy"] = askBool(scanner, w, s, "/deploy     — Deployment checklist", false)
	cfg.SelectedSkills["security"] = askBool(scanner, w, s, "/security   — Security audit", false)
	cfg.SelectedSkills["adr"] = askBool(scanner, w, s, "/adr        — Architecture Decision Record", false)
	cfg.SelectedSkills["standup"] = askBool(scanner, w, s, "/standup    — Standup summary", false)
	cfg.SelectedSkills["cuj-list"] = askBool(scanner, w, s, "/cuj-list   — Create CUJ inventory", false)
	cfg.SelectedSkills["cuj-test"] = askBool(scanner, w, s, "/cuj-test   — Run CUJ tests (headless Chromium)", false)

	fmt.Fprintln(w)
	fmt.Fprintf(w, "  %s Custom Skills %s\n", s.Green("▸"), s.Dim("(enter name: description, blank line to finish)"))
	for {
		fmt.Fprintf(w, "    %s ", s.Dim("│"))
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			cfg.CustomSkills = append(cfg.CustomSkills, config.CustomSkillConfig{
				Name:        strings.TrimSpace(parts[0]),
				Description: strings.TrimSpace(parts[1]),
			})
		}
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, s.Divider())

	// Step 9: Confirm
	fmt.Fprintln(w, s.StepHeader(9, totalSteps, "Confirm"))

	// Summary
	printSummary(w, s, cfg)

	if !askBool(scanner, w, s, "Proceed with generation?", true) {
		return nil, fmt.Errorf("generation cancelled by user")
	}

	return cfg, nil
}

func ask(scanner *bufio.Scanner, w io.Writer, s *Styler, prompt, defaultVal string) string {
	if defaultVal != "" {
		fmt.Fprintf(w, "  %s %s %s: ", s.Green("▸"), prompt, s.Dim("["+defaultVal+"]"))
	} else {
		fmt.Fprintf(w, "  %s %s: ", s.Green("▸"), prompt)
	}
	if !scanner.Scan() {
		return defaultVal
	}
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal
	}
	return val
}

func askBool(scanner *bufio.Scanner, w io.Writer, s *Styler, prompt string, defaultVal bool) bool {
	hint := "(y/N)"
	if defaultVal {
		hint = "(Y/n)"
	}
	fmt.Fprintf(w, "  %s %s %s: ", s.Green("▸"), prompt, s.Dim(hint))
	if !scanner.Scan() {
		return defaultVal
	}
	val := strings.TrimSpace(strings.ToLower(scanner.Text()))
	if val == "" {
		return defaultVal
	}
	return val == "y" || val == "yes"
}

func askInt(scanner *bufio.Scanner, w io.Writer, s *Styler, prompt string, defaultVal int) int {
	fmt.Fprintf(w, "  %s %s %s: ", s.Green("▸"), prompt, s.Dim(fmt.Sprintf("[%d]", defaultVal)))
	if !scanner.Scan() {
		return defaultVal
	}
	val := strings.TrimSpace(scanner.Text())
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func printSummary(w io.Writer, s *Styler, cfg *config.ProjectConfig) {
	top := "┌──────────────────────────────────────────────────┐"
	mid := "├──────────────────────────────────────────────────┤"
	bot := "└──────────────────────────────────────────────────┘"

	fmt.Fprintln(w, top)
	fmt.Fprintf(w, "│  %s%s│\n", s.BoldCyan("Configuration Summary"), strings.Repeat(" ", 27))
	fmt.Fprintln(w, mid)
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Project:", 14), cfg.ProjectName)
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Description:", 14), cfg.Description)
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Tech Stack:", 14), cfg.TechStack)
	fmt.Fprintf(w, "│  %s %s <%s> (@%s)\n", s.PadBold("Owner:", 14), cfg.OwnerName, cfg.OwnerEmail, cfg.OwnerGitHub)
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Framework:", 14), cfg.Framework)
	if cfg.Framework == "scion" {
		fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Harness:", 14), cfg.DefaultHarness)
	}
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Team Size:", 14), cfg.TeamSize)
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("GCP:", 14), boolLabel(s, cfg.GCP.Enabled))
	if cfg.GCP.Enabled {
		fmt.Fprintf(w, "│    %s %s\n", s.PadBold("Project ID:", 12), cfg.GCP.ProjectID)
		fmt.Fprintf(w, "│    %s %s\n", s.PadBold("Region:", 12), cfg.GCP.Region)
	}
	fmt.Fprintf(w, "│  %s %d\n", s.PadBold("SWE Agents:", 14), len(cfg.SWEs))
	for _, swe := range cfg.SWEs {
		fmt.Fprintf(w, "│    SWE-%d: %s\n", swe.Number, swe.Title)
	}
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Platform:", 14), boolLabel(s, cfg.IncludePlatform))
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Reviewer:", 14), boolLabel(s, cfg.IncludeReviewer))
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("SWE-Test:", 14), boolLabel(s, cfg.IncludeSWETest))
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("SWE-QA:", 14), boolLabel(s, cfg.IncludeSWEQA))
	if len(cfg.CustomAgents) > 0 {
		fmt.Fprintf(w, "│  %s %d\n", s.PadBold("Custom Agents:", 14), len(cfg.CustomAgents))
		for _, ca := range cfg.CustomAgents {
			fmt.Fprintf(w, "│    %s: %s\n", ca.Name, ca.Title)
		}
	}
	if len(cfg.Conventions) > 0 {
		fmt.Fprintf(w, "│  %s\n", s.Bold("Conventions:"))
		for _, c := range cfg.Conventions {
			fmt.Fprintf(w, "│    • %s\n", c)
		}
	}
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Skills:", 14), skillsSummary(cfg.SelectedSkills, len(cfg.CustomSkills)))
	if len(cfg.CustomSkills) > 0 {
		for _, cs := range cfg.CustomSkills {
			fmt.Fprintf(w, "│    /%s — %s\n", cs.Name, cs.Description)
		}
	}
	fmt.Fprintf(w, "│  %s %s %s\n", s.PadBold("Model:", 14), cfg.ModelName, s.Dim("("+cfg.ModelID+")"))
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Target:", 14), cfg.TargetDir)
	fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Git Init:", 14), boolLabel(s, cfg.InitGit))
	if cfg.RepoURL != "" {
		fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Repo URL:", 14), cfg.RepoURL)
	}
	if cfg.CreateRepo {
		fmt.Fprintf(w, "│  %s %s\n", s.PadBold("Create Repo:", 14), boolLabel(s, cfg.CreateRepo))
	}
	fmt.Fprintln(w, bot)
	fmt.Fprintln(w)
}

func skillsSummary(selected map[string]bool, customCount int) string {
	allSkills := []string{
		"spec", "release", "roadmap", "status", "pipeline", "regenerate", "brainstorm",
		"debug", "test", "review", "docs", "refactor", "hotfix",
		"api-design", "schema", "deploy", "security", "adr", "standup",
		"cuj-list", "cuj-test",
	}
	var active []string
	for _, name := range allSkills {
		if selected[name] {
			active = append(active, "/"+name)
		}
	}
	summary := fmt.Sprintf("%d of %d", len(active), len(allSkills))
	if customCount > 0 {
		summary += fmt.Sprintf(" + %d custom", customCount)
	}
	return summary
}

func boolLabel(s *Styler, v bool) string {
	if v {
		return s.BoldGreen("✓ Yes")
	}
	return s.Dim("✗ No")
}
