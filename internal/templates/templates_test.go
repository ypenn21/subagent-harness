package templates_test

import (
	"strings"
	"testing"

	"github.com/ahafin/appteam/internal/config"
	"github.com/ahafin/appteam/internal/templates"
)

// fullConfig returns a fully populated ProjectConfig for template testing.
func fullConfig() *config.ProjectConfig {
	return &config.ProjectConfig{
		ProjectName: "testproject",
		Description: "A test project for unit tests",
		TechStack:   "Go",
		OwnerName:   "Test Owner",
		OwnerEmail:  "test@example.com",
		OwnerGitHub: "testowner",
		ModelName:   "Opus 4.6",
		ModelID:     "claude-opus-4-6",
		SWEs: []config.SWEConfig{
			{Number: 1, Title: "CLI & Wizard", Bullets: []string{"CLI flags", "Interactive prompts"}},
			{Number: 2, Title: "Templates & Generation", Bullets: []string{"Template rendering", "File output"}},
		},
		GCP: config.GCPConfig{
			Enabled:       true,
			ProjectID:     "test-gcp-project",
			ProjectNumber: "123456789",
			Organization:  "test-org",
			Region:        "us-central1",
		},
		IncludeSWETest:  true,
		IncludeSWEQA:    true,
		IncludePlatform: true,
		IncludeReviewer: true,
		Conventions:     []string{"Use stdlib only", "Table-driven tests"},
		TargetDir:       "/tmp/test",
	}
}

func TestAddFunc(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{2, 3, 5},
		{0, 0, 0},
		{-1, 1, 0},
		{10, -3, 7},
	}
	for _, tc := range tests {
		out, err := templates.Render("add", `{{add .A .B}}`, struct{ A, B int }{tc.a, tc.b})
		if err != nil {
			t.Fatalf("add(%d, %d): unexpected error: %v", tc.a, tc.b, err)
		}
		if out != itoa(tc.want) {
			t.Errorf("add(%d, %d) = %q, want %q", tc.a, tc.b, out, itoa(tc.want))
		}
	}
}

func TestSubFunc(t *testing.T) {
	tests := []struct {
		a, b int
		want int
	}{
		{5, 3, 2},
		{0, 0, 0},
		{1, 5, -4},
	}
	for _, tc := range tests {
		out, err := templates.Render("sub", `{{sub .A .B}}`, struct{ A, B int }{tc.a, tc.b})
		if err != nil {
			t.Fatalf("sub(%d, %d): unexpected error: %v", tc.a, tc.b, err)
		}
		if out != itoa(tc.want) {
			t.Errorf("sub(%d, %d) = %q, want %q", tc.a, tc.b, out, itoa(tc.want))
		}
	}
}

func TestSeqFunc(t *testing.T) {
	tests := []struct {
		name       string
		start, end int
		want       string
	}{
		{"1 to 5", 1, 5, "1 2 3 4 5"},
		{"3 to 3", 3, 3, "3"},
		{"1 to 1", 1, 1, "1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpl := `{{range $i, $v := seq .Start .End}}{{if $i}} {{end}}{{$v}}{{end}}`
			out, err := templates.Render("seq", tmpl, struct{ Start, End int }{tc.start, tc.end})
			if err != nil {
				t.Fatalf("seq(%d, %d): unexpected error: %v", tc.start, tc.end, err)
			}
			if out != tc.want {
				t.Errorf("seq(%d, %d) = %q, want %q", tc.start, tc.end, out, tc.want)
			}
		})
	}
}

func TestLinkRangeFunc(t *testing.T) {
	tests := []struct {
		name       string
		start, cnt int
		want       string
	}{
		{"0,3", 0, 3, "0,1,2"},
		{"5,0", 5, 0, ""},
		{"2,1", 2, 1, "2"},
		{"0,1", 0, 1, "0"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := templates.Render("lr", `{{linkRange .Start .Count}}`, struct{ Start, Count int }{tc.start, tc.cnt})
			if err != nil {
				t.Fatalf("linkRange(%d, %d): unexpected error: %v", tc.start, tc.cnt, err)
			}
			if out != tc.want {
				t.Errorf("linkRange(%d, %d) = %q, want %q", tc.start, tc.cnt, out, tc.want)
			}
		})
	}
}

func TestRenderBasic(t *testing.T) {
	out, err := templates.Render("basic", "Hello, {{.Name}}!", struct{ Name string }{"World"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "Hello, World!" {
		t.Errorf("got %q, want %q", out, "Hello, World!")
	}
}

func TestRenderInvalidTemplate(t *testing.T) {
	_, err := templates.Render("bad", "{{.Foo", nil)
	if err == nil {
		t.Fatal("expected error for invalid template, got nil")
	}
}

func TestRenderMissingField(t *testing.T) {
	_, err := templates.Render("missing", "{{.NoSuchField}}", struct{}{})
	if err == nil {
		t.Fatal("expected error for missing field, got nil")
	}
}

func TestRenderClaudeMD(t *testing.T) {
	cfg := fullConfig()
	cfg.GCP.Enabled = false

	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"# testproject — Project Instructions",
		"Test Owner",
		"test@example.com",
		"testowner",
		"SWE-1",
		"SWE-2",
		"CLI & Wizard",
		"Templates & Generation",
		"Opus 4.6",
		"claude-opus-4-6",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("ClaudeMD output missing %q", s)
		}
	}

	if strings.Contains(out, "## GCP Project") {
		t.Error("ClaudeMD output should NOT contain '## GCP Project' when GCP is disabled")
	}
}

func TestRenderClaudeMDWithGCP(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"## GCP Project",
		"test-gcp-project",
		"123456789",
		"test-org",
		"us-central1",
		"## GCP Free Tier",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("ClaudeMD+GCP output missing %q", s)
		}
	}
}

func TestRenderClaudeMDWithoutGCP(t *testing.T) {
	cfg := fullConfig()
	cfg.GCP.Enabled = false

	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	noContain := []string{
		"## GCP Project",
		"## GCP Free Tier",
		"test-gcp-project",
	}
	for _, s := range noContain {
		if strings.Contains(out, s) {
			t.Errorf("ClaudeMD without GCP should NOT contain %q", s)
		}
	}
}

func TestRenderSWETemplate(t *testing.T) {
	data := templates.SWETemplateData{
		ProjectName: "testproject",
		OwnerName:   "Test Owner",
		OwnerEmail:  "test@example.com",
		ModelName:   "Opus 4.6",
		Number:      1,
		Title:       "CLI & Wizard",
		Bullets:     []string{"CLI flags", "Interactive prompts"},
	}

	out, err := templates.Render("swe.md", templates.SWETemplate, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"# SWE-1 Agent — CLI & Wizard",
		"CLI flags",
		"Interactive prompts",
		"Your specialty is CLI & Wizard",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SWE template output missing %q", s)
		}
	}
}

func TestRenderSWETemplateWithoutBullets(t *testing.T) {
	data := templates.SWETemplateData{
		ProjectName: "testproject",
		OwnerName:   "Test Owner",
		OwnerEmail:  "test@example.com",
		ModelName:   "Opus 4.6",
		Number:      3,
		Title:       "General",
		Bullets:     nil,
	}

	out, err := templates.Render("swe.md", templates.SWETemplate, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"You are additional engineering capacity",
		"General full-stack development",
		"Assigned by TPM based on current workload",
		"Can take on any tasks as assigned",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SWE no-bullets output missing %q", s)
		}
	}
}

func TestRenderAllTemplates(t *testing.T) {
	cfg := fullConfig()

	allTemplates := map[string]struct {
		tmpl string
		data any
	}{
		"ClaudeMD":        {templates.ClaudeMDTemplate, cfg},
		"PM":              {templates.PMTemplate, cfg},
		"TPM":             {templates.TPMTemplate, cfg},
		"SWETest":         {templates.SWETestTemplate, cfg},
		"SWEQA":           {templates.SWEQATemplate, cfg},
		"Platform":        {templates.PlatformTemplate, cfg},
		"Reviewer":        {templates.ReviewerTemplate, cfg},
		"Backlog":         {templates.BacklogTemplate, cfg},
		"Progress":        {templates.ProgressTemplate, cfg},
		"ReleaseNotes":    {templates.ReleaseNotesTemplate, cfg},
		"Pipeline":        {templates.PipelineTemplate, cfg},
		"SpecTemplate":    {templates.SpecTemplateFile, cfg},
		"SkillSpec":       {templates.SkillSpecTemplate, cfg},
		"SkillRelease":    {templates.SkillReleaseTemplate, cfg},
		"SkillPipeline":   {templates.SkillPipelineTemplate, cfg},
		"SkillStatus":     {templates.SkillStatusTemplate, cfg},
		"SkillRegenerate": {templates.SkillRegenerateTemplate, cfg},
		"SkillRoadmap":    {templates.SkillRoadmapTemplate, cfg},
		"SkillDebug":      {templates.SkillDebugTemplate, cfg},
		"SkillTestWrite":  {templates.SkillTestWriteTemplate, cfg},
		"SkillReview":     {templates.SkillReviewTemplate, cfg},
		"SkillDocs":       {templates.SkillDocsTemplate, cfg},
		"SkillRefactor":   {templates.SkillRefactorTemplate, cfg},
		"SkillHotfix":     {templates.SkillHotfixTemplate, cfg},
		"SkillAPIDesign":  {templates.SkillAPIDesignTemplate, cfg},
		"SkillSchema":     {templates.SkillSchemaTemplate, cfg},
		"SkillDeploy":     {templates.SkillDeployTemplate, cfg},
		"SkillSecurity":   {templates.SkillSecurityTemplate, cfg},
		"SkillADR":        {templates.SkillADRTemplate, cfg},
		"SkillStandup":    {templates.SkillStandupTemplate, cfg},
		"SkillBrainstorm": {templates.SkillBrainstormTemplate, cfg},
		"SkillCUJList":    {templates.SkillCUJListTemplate, cfg},
		"SkillCUJTest":    {templates.SkillCUJTestTemplate, cfg},
		"SWE": {templates.SWETemplate, templates.SWETemplateData{
			ProjectName: cfg.ProjectName,
			OwnerName:   cfg.OwnerName,
			OwnerEmail:  cfg.OwnerEmail,
			ModelName:   cfg.ModelName,
			Number:      1,
			Title:       "CLI & Wizard",
			Bullets:     []string{"CLI flags"},
			TeamSize:    "standard",
		}},
		"CustomAgent": {templates.CustomAgentTemplate, templates.CustomAgentTemplateData{
			Name:         "dba",
			Title:        "Database Administrator",
			Description:  "Manages database schemas and migrations",
			Instructions: []string{"Review schema changes", "Optimize queries"},
			ProjectName:  cfg.ProjectName,
			OwnerName:    cfg.OwnerName,
			OwnerEmail:   cfg.OwnerEmail,
			ModelName:    cfg.ModelName,
			TeamSize:     "standard",
		}},
		"ScionPM":       {templates.ScionPMAgentsMD, cfg},
		"ScionTPM":      {templates.ScionTPMAgentsMD, cfg},
		"ScionSWETest":  {templates.ScionSWETestAgentsMD, cfg},
		"ScionSWEQA":    {templates.ScionSWEQAAgentsMD, cfg},
		"ScionReviewer": {templates.ScionReviewerAgentsMD, cfg},
		"ScionPlatform": {templates.ScionPlatformAgentsMD, cfg},
		"ScionSWE": {templates.ScionSWEAgentsMD, templates.SWETemplateData{
			ProjectName: cfg.ProjectName,
			OwnerName:   cfg.OwnerName,
			OwnerEmail:  cfg.OwnerEmail,
			ModelName:   cfg.ModelName,
			Number:      1,
			Title:       "CLI & Wizard",
			Bullets:     []string{"CLI flags"},
			TeamSize:    "standard",
		}},
		"ScionCustomAgent": {templates.ScionCustomAgentAgentsMD, templates.CustomAgentTemplateData{
			Name:         "dba",
			Title:        "Database Administrator",
			Description:  "Manages database schemas and migrations",
			Instructions: []string{"Review schema changes", "Optimize queries"},
			ProjectName:  cfg.ProjectName,
			OwnerName:    cfg.OwnerName,
			OwnerEmail:   cfg.OwnerEmail,
			ModelName:    cfg.ModelName,
			TeamSize:     "standard",
		}},
	}

	for name, tc := range allTemplates {
		t.Run(name, func(t *testing.T) {
			out, err := templates.Render(name, tc.tmpl, tc.data)
			if err != nil {
				t.Fatalf("template %s render error: %v", name, err)
			}
			if len(out) == 0 {
				t.Errorf("template %s produced empty output", name)
			}
		})
	}
}

func TestOptionalAgentTemplates(t *testing.T) {
	cfg := fullConfig()

	tests := []struct {
		name     string
		tmpl     string
		contains []string
	}{
		{
			"SWE-Test",
			templates.SWETestTemplate,
			[]string{"SWE-Test", "Test Engineer", "testproject"},
		},
		{
			"SWE-QA",
			templates.SWEQATemplate,
			[]string{"SWE-QA", "QA Engineer", "testproject", "CUJ", "Puppeteer", "Chromium", "Lighthouse"},
		},
		{
			"Platform",
			templates.PlatformTemplate,
			[]string{"Platform Engineer", "testproject"},
		},
		{
			"Reviewer",
			templates.ReviewerTemplate,
			[]string{"Reviewer", "Code Review", "testproject"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := templates.Render(tc.name, tc.tmpl, cfg)
			if err != nil {
				t.Fatalf("template %s render error: %v", tc.name, err)
			}
			for _, s := range tc.contains {
				if !strings.Contains(out, s) {
					t.Errorf("template %s output missing %q", tc.name, s)
				}
			}
		})
	}
}

func TestConventionsInClaudeMD(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "## Important Conventions") {
		t.Error("output should contain '## Important Conventions' when conventions are set")
	}
	if !strings.Contains(out, "Use stdlib only") {
		t.Error("output missing convention 'Use stdlib only'")
	}
	if !strings.Contains(out, "Table-driven tests") {
		t.Error("output missing convention 'Table-driven tests'")
	}

	// Test without conventions
	cfg.Conventions = nil
	out2, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out2, "## Important Conventions") {
		t.Error("output should NOT contain '## Important Conventions' when conventions are empty")
	}
}

func TestSWECountVariations(t *testing.T) {
	tests := []struct {
		name     string
		sweCount int
	}{
		{"1 SWE", 1},
		{"3 SWEs", 3},
		{"5 SWEs", 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := fullConfig()
			cfg.SWEs = make([]config.SWEConfig, tc.sweCount)
			for i := range cfg.SWEs {
				cfg.SWEs[i] = config.SWEConfig{
					Number:  i + 1,
					Title:   "General",
					Bullets: []string{"General development"},
				}
			}

			out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
			if err != nil {
				t.Fatalf("unexpected error with %d SWEs: %v", tc.sweCount, err)
			}

			for i := 1; i <= tc.sweCount; i++ {
				marker := "SWE-" + itoa(i)
				if !strings.Contains(out, marker) {
					t.Errorf("output with %d SWEs missing %q", tc.sweCount, marker)
				}
			}
		})
	}
}

func TestOptionalAgentsInClaudeMD(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*config.ProjectConfig)
		present string
		absent  string
	}{
		{
			"SWE-Test included",
			func(c *config.ProjectConfig) { c.IncludeSWETest = true },
			"SWE-Test",
			"",
		},
		{
			"SWE-Test excluded",
			func(c *config.ProjectConfig) { c.IncludeSWETest = false },
			"",
			"SWE-Test",
		},
		{
			"SWE-QA included",
			func(c *config.ProjectConfig) { c.IncludeSWEQA = true },
			"SWE-QA",
			"",
		},
		{
			"SWE-QA excluded",
			func(c *config.ProjectConfig) { c.IncludeSWEQA = false },
			"",
			"SWE-QA",
		},
		{
			"Platform included",
			func(c *config.ProjectConfig) { c.IncludePlatform = true },
			"Platform Engineer",
			"",
		},
		{
			"Platform excluded",
			func(c *config.ProjectConfig) { c.IncludePlatform = false },
			"",
			"Platform Engineer (PE)",
		},
		{
			"Reviewer included",
			func(c *config.ProjectConfig) { c.IncludeReviewer = true },
			"Code review, quality/security/performance",
			"",
		},
		{
			"Reviewer excluded",
			func(c *config.ProjectConfig) { c.IncludeReviewer = false },
			"",
			"Code review, quality/security/performance",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := fullConfig()
			// Disable all optional agents first
			cfg.IncludeSWETest = false
			cfg.IncludeSWEQA = false
			cfg.IncludePlatform = false
			cfg.IncludeReviewer = false
			tc.setup(cfg)

			out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.present != "" && !strings.Contains(out, tc.present) {
				t.Errorf("output should contain %q", tc.present)
			}
			if tc.absent != "" && strings.Contains(out, tc.absent) {
				t.Errorf("output should NOT contain %q", tc.absent)
			}
		})
	}
}

func TestTemplateOutputContainsExpectedMarkers(t *testing.T) {
	cfg := fullConfig()

	tests := []struct {
		name     string
		tmpl     string
		markers  []string
	}{
		{"PM", templates.PMTemplate, []string{"PM Agent", "Product Manager", "testproject"}},
		{"TPM", templates.TPMTemplate, []string{"TPM Agent", "Technical Program Manager", "testproject"}},
		{"SWE-Test", templates.SWETestTemplate, []string{"SWE-Test Agent", "testproject"}},
		{"Reviewer", templates.ReviewerTemplate, []string{"Reviewer Agent", "testproject"}},
		{"Platform", templates.PlatformTemplate, []string{"Platform Engineer", "testproject"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := templates.Render(tc.name, tc.tmpl, cfg)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			for _, m := range tc.markers {
				if !strings.Contains(out, m) {
					t.Errorf("template %s output missing marker %q", tc.name, m)
				}
			}
		})
	}
}

func TestRenderSkillRoadmapTemplate(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("roadmap.md", templates.SkillRoadmapTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"/roadmap",
		"Add items to the backlog as future work",
		"docs/BACKLOG.md",
		"F-NNNN",
		"Feature name",
		"Priority",
		"P0",
		"P1",
		"P2",
		"Future Items",
		"testproject",
		"Test Owner",
		"test@example.com",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SkillRoadmap output missing %q", s)
		}
	}
}

func TestRenderSkillDebugTemplate(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("debug.md", templates.SkillDebugTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"/debug",
		"Systematic debugging",
		"Reproduce the issue",
		"root cause",
		"regression test",
		"git bisect",
		"testproject",
		"Test Owner",
		"test@example.com",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SkillDebug output missing %q", s)
		}
	}
}

func TestRenderSkillReviewTemplate(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("review.md", templates.SkillReviewTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"/review",
		"Code review",
		"OWASP",
		"Critical",
		"Warning",
		"Info",
		"testproject",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SkillReview output missing %q", s)
		}
	}
}

func TestRenderCustomAgentTemplate(t *testing.T) {
	data := templates.CustomAgentTemplateData{
		Name:         "frontend-eng",
		Title:        "Frontend Engineer",
		Description:  "React/TypeScript specialist for UI components",
		Instructions: []string{"Build responsive UI components", "Follow accessibility standards"},
		ProjectName:  "testproject",
		OwnerName:    "Test Owner",
		OwnerEmail:   "test@example.com",
		ModelName:    "Opus 4.6",
	}

	out, err := templates.Render("custom_agent.md", templates.CustomAgentTemplate, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"# Frontend Engineer Agent",
		"frontend-eng",
		"React/TypeScript specialist for UI components",
		"Build responsive UI components",
		"Follow accessibility standards",
		"testproject",
		"Test Owner",
		"test@example.com",
		"Opus 4.6",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("CustomAgent template output missing %q", s)
		}
	}
}

func TestRenderSkillBrainstormTemplate(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("brainstorm.md", templates.SkillBrainstormTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"/brainstorm",
		"Brainstorm product ideas with PO",
		"Vision",
		"Ideation",
		"Competitive Analysis",
		"Prioritization",
		"Action Items",
		"/spec",
		"/roadmap",
		"testproject",
		"Test Owner",
		"test@example.com",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SkillBrainstorm output missing %q", s)
		}
	}
}

func TestRenderClaudeMDTeamSizeLean(t *testing.T) {
	cfg := fullConfig()
	cfg.TeamSize = "lean"

	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Lean pipeline: no TPM, PM handles everything
	mustContain := []string{
		"Manages backlog, progress log, and release notes",
		"PM → SWE-1",
		"SWE-1 → Testing",
		"SWE-1 → PM",
		"PM → Progress + Release Notes",
		"Backlog Tracking (Non-Negotiable)",
		"Every piece of work gets a backlog entry",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("lean CLAUDE.md missing %q", s)
		}
	}

	mustNotContain := []string{
		"**TPM**",
		"TPM → SWE",
		"SWE → TPM",
		"TPM → Progress",
		"TPM → PM",
	}
	for _, s := range mustNotContain {
		if strings.Contains(out, s) {
			t.Errorf("lean CLAUDE.md should NOT contain %q", s)
		}
	}
}

func TestRenderClaudeMDTeamSizeStandard(t *testing.T) {
	cfg := fullConfig()
	cfg.TeamSize = "standard"

	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Standard pipeline: has TPM, conditional agents
	mustContain := []string{
		"**TPM**",
		"TPM → SWE",
		"SWE → TPM",
		"TPM → Progress",
		"Backlog Tracking (Non-Negotiable)",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("standard CLAUDE.md missing %q", s)
		}
	}
}

func TestRenderClaudeMDTeamSizeFull(t *testing.T) {
	cfg := fullConfig()
	cfg.TeamSize = "full"

	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Full pipeline: all agents always present
	mustContain := []string{
		"**TPM**",
		"**SWE-Test**",
		"**SWE-QA**",
		"**Platform Engineer (PE)**",
		"**Reviewer**",
		"Backlog Tracking (Non-Negotiable)",
		"Every piece of work gets a backlog entry",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("full CLAUDE.md missing %q", s)
		}
	}
}

func TestBacklogEnforcementAllTeamSizes(t *testing.T) {
	for _, size := range []string{"lean", "standard", "full"} {
		t.Run(size, func(t *testing.T) {
			cfg := fullConfig()
			cfg.TeamSize = size

			out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !strings.Contains(out, "### Backlog Tracking (Non-Negotiable)") {
				t.Errorf("team size %q: missing backlog enforcement section", size)
			}
			if !strings.Contains(out, "Every piece of work gets a backlog entry") {
				t.Errorf("team size %q: missing backlog enforcement text", size)
			}
		})
	}
}

func TestRenderSWEQATemplate(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("swe-qa.md", templates.SWEQATemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"SWE-QA",
		"CUJ",
		"Critical User Journey",
		"Puppeteer",
		"Chromium",
		"Lighthouse",
		"page.goto",
		"page.waitForSelector",
		"page.screenshot",
		"page.on('console'",
		"page.setViewport",
		"docs/CUJ.md",
		"testproject",
		"Test Owner",
		"test@example.com",
		"Opus 4.6",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SWE-QA template output missing %q", s)
		}
	}
}

func TestRenderSkillCUJListTemplate(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("cuj-list.md", templates.SkillCUJListTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"/cuj-list",
		"CUJ",
		"Critical User Journey",
		"docs/CUJ.md",
		"CUJ-001",
		"P0",
		"P1",
		"P2",
		"testproject",
		"Test Owner",
		"test@example.com",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SkillCUJList output missing %q", s)
		}
	}
}

func TestRenderSkillCUJTestTemplate(t *testing.T) {
	cfg := fullConfig()

	out, err := templates.Render("cuj-test.md", templates.SkillCUJTestTemplate, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mustContain := []string{
		"/cuj-test",
		"CUJ",
		"headless",
		"Chromium",
		"Puppeteer",
		"page.goto",
		"page.screenshot",
		"screenshots/",
		"docs/CUJ.md",
		"PASS",
		"FAIL",
		"Phase 6: Cleanup",
		"browser.close()",
		"Cleanup Report",
		"testproject",
		"Test Owner",
		"test@example.com",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("SkillCUJTest output missing %q", s)
		}
	}
}

func scionConfig() *config.ProjectConfig {
	cfg := fullConfig()
	cfg.Framework = "scion"
	cfg.DefaultHarness = "claude-code"
	return cfg
}

func TestScionPMMessagingProtocol(t *testing.T) {
	tests := []struct {
		name     string
		teamSize string
		contains []string
		absent   []string
	}{
		{
			"lean — PM signals completion for SWE-1",
			"lean",
			[]string{"Signaling Protocol", "sciontool status task_completed", "SWE-1 should begin",
				"Worktree Isolation", "No merge is needed before your work"},
			[]string{},
		},
		{
			"standard — PM signals completion for TPM",
			"standard",
			[]string{"Signaling Protocol", "sciontool status task_completed", "TPM should proceed",
				"Worktree Isolation", "No merge is needed before your work"},
			[]string{},
		},
		{
			"full — PM signals completion for TPM",
			"full",
			[]string{"Signaling Protocol", "sciontool status task_completed", "TPM should proceed",
				"Worktree Isolation", "No merge is needed before your work"},
			[]string{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := scionConfig()
			cfg.TeamSize = tc.teamSize
			out, err := templates.Render("scion-pm", templates.ScionPMAgentsMD, cfg)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			for _, s := range tc.contains {
				if !strings.Contains(out, s) {
					t.Errorf("output missing %q", s)
				}
			}
			for _, s := range tc.absent {
				if strings.Contains(out, s) {
					t.Errorf("output should NOT contain %q", s)
				}
			}
		})
	}
}

func TestScionTPMMessagingProtocol(t *testing.T) {
	cfg := scionConfig()
	cfg.TeamSize = "standard"
	out, err := templates.Render("scion-tpm", templates.ScionTPMAgentsMD, cfg)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	mustContain := []string{
		"Signaling Protocol",
		"Wait for:",
		"message from PM",
		"sciontool status task_completed",
		"SWEs should begin",
		"PM should proceed with release notes",
		"Worktree Isolation",
		"git merge pm",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestScionSWEMessagingProtocol(t *testing.T) {
	tests := []struct {
		name      string
		teamSize  string
		waitFor   string
		mergeFrom string
	}{
		{"lean — waits for PM", "lean", "message from PM", "git merge pm"},
		{"standard — waits for TPM", "standard", "message from TPM", "git merge tpm"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := templates.SWETemplateData{
				ProjectName: "testproject",
				OwnerName:   "Test Owner",
				OwnerEmail:  "test@example.com",
				ModelName:   "Opus 4.6",
				Number:      1,
				Title:       "General",
				Bullets:     []string{"General dev"},
				TeamSize:    tc.teamSize,
			}
			out, err := templates.Render("scion-swe", templates.ScionSWEAgentsMD, data)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			if !strings.Contains(out, "Signaling Protocol") {
				t.Error("output missing 'Signaling Protocol'")
			}
			if !strings.Contains(out, tc.waitFor) {
				t.Errorf("output missing %q", tc.waitFor)
			}
			if !strings.Contains(out, "sciontool status task_completed") {
				t.Error("output missing 'sciontool status task_completed'")
			}
			if !strings.Contains(out, "Worktree Isolation") {
				t.Error("output missing 'Worktree Isolation'")
			}
			if !strings.Contains(out, tc.mergeFrom) {
				t.Errorf("output missing %q", tc.mergeFrom)
			}
		})
	}
}

func TestScionSWETestMessagingProtocol(t *testing.T) {
	tests := []struct {
		name     string
		teamSize string
		reviewer bool
		target   string
	}{
		{"lean — signals PM", "lean", false, "PM should proceed"},
		{"standard with reviewer", "standard", true, "Reviewer should proceed"},
		{"standard no reviewer", "standard", false, "PM should proceed"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := scionConfig()
			cfg.TeamSize = tc.teamSize
			cfg.IncludeReviewer = tc.reviewer
			out, err := templates.Render("scion-swe-test", templates.ScionSWETestAgentsMD, cfg)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			if !strings.Contains(out, "Signaling Protocol") {
				t.Error("output missing 'Signaling Protocol'")
			}
			if !strings.Contains(out, tc.target) {
				t.Errorf("output missing %q", tc.target)
			}
			if !strings.Contains(out, "Worktree Isolation") {
				t.Error("output missing 'Worktree Isolation'")
			}
			if !strings.Contains(out, "git merge swe-1") {
				t.Error("output missing 'git merge swe-1'")
			}
		})
	}
}

func TestScionSWEQAMessagingProtocol(t *testing.T) {
	cfg := scionConfig()
	cfg.TeamSize = "full"
	out, err := templates.Render("scion-swe-qa", templates.ScionSWEQAAgentsMD, cfg)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	mustContain := []string{
		"Signaling Protocol",
		"Wait for:",
		"message from SWE-Test",
		"sciontool status task_completed",
		"Reviewer should proceed",
		"Worktree Isolation",
		"git merge swe-test",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestScionReviewerMessagingProtocol(t *testing.T) {
	cfg := scionConfig()
	cfg.TeamSize = "full"
	out, err := templates.Render("scion-reviewer", templates.ScionReviewerAgentsMD, cfg)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	mustContain := []string{
		"Signaling Protocol",
		"Wait for:",
		"message from SWE-Test",
		"SWE-QA",
		"sciontool status task_completed",
		"PM should proceed with release notes",
		"Worktree Isolation",
		"git merge swe-qa",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestScionPlatformMessagingProtocol(t *testing.T) {
	cfg := scionConfig()
	cfg.TeamSize = "full"
	out, err := templates.Render("scion-platform", templates.ScionPlatformAgentsMD, cfg)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	mustContain := []string{
		"Signaling Protocol",
		"sciontool status task_completed",
		"TPM should",
		"Worktree Isolation",
		"git merge tpm",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("output missing %q", s)
		}
	}
}

func TestScionCustomAgentMessagingProtocol(t *testing.T) {
	tests := []struct {
		name     string
		teamSize string
		target   string
	}{
		{"lean — signals PM", "lean", "PM should proceed"},
		{"standard — signals TPM", "standard", "TPM should update progress"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := templates.CustomAgentTemplateData{
				Name:         "dba",
				Title:        "Database Admin",
				Description:  "Manages DB",
				Instructions: []string{"Manage schemas"},
				ProjectName:  "testproject",
				OwnerName:    "Test Owner",
				OwnerEmail:   "test@example.com",
				ModelName:    "Opus 4.6",
				TeamSize:     tc.teamSize,
			}
			out, err := templates.Render("scion-custom", templates.ScionCustomAgentAgentsMD, data)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			if !strings.Contains(out, "Signaling Protocol") {
				t.Error("output missing 'Signaling Protocol'")
			}
			if !strings.Contains(out, tc.target) {
				t.Errorf("output missing %q", tc.target)
			}
			if !strings.Contains(out, "Worktree Isolation") {
				t.Error("output missing 'Worktree Isolation'")
			}
			if !strings.Contains(out, "Merge the branch") {
				t.Error("output missing generic merge guidance")
			}
		})
	}
}

func TestScionReviewerMergeConditional(t *testing.T) {
	tests := []struct {
		name      string
		includeQA bool
		mergeFrom string
		absent    string
	}{
		{"with SWE-QA — merges swe-qa", true, "git merge swe-qa", "git merge swe-test"},
		{"without SWE-QA — merges swe-test", false, "git merge swe-test", "git merge swe-qa"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := scionConfig()
			cfg.TeamSize = "standard"
			cfg.IncludeSWEQA = tc.includeQA
			out, err := templates.Render("scion-reviewer", templates.ScionReviewerAgentsMD, cfg)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			if !strings.Contains(out, "Worktree Isolation") {
				t.Error("output missing 'Worktree Isolation'")
			}
			if !strings.Contains(out, tc.mergeFrom) {
				t.Errorf("output missing %q", tc.mergeFrom)
			}
			if strings.Contains(out, tc.absent) {
				t.Errorf("output should NOT contain %q", tc.absent)
			}
		})
	}
}

func TestScionPlatformMergeConditional(t *testing.T) {
	tests := []struct {
		name      string
		teamSize  string
		mergeFrom string
	}{
		{"lean — merges pm", "lean", "git merge pm"},
		{"standard — merges tpm", "standard", "git merge tpm"},
		{"full — merges tpm", "full", "git merge tpm"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := scionConfig()
			cfg.TeamSize = tc.teamSize
			out, err := templates.Render("scion-platform", templates.ScionPlatformAgentsMD, cfg)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			if !strings.Contains(out, "Worktree Isolation") {
				t.Error("output missing 'Worktree Isolation'")
			}
			if !strings.Contains(out, tc.mergeFrom) {
				t.Errorf("output missing %q", tc.mergeFrom)
			}
		})
	}
}

func TestScionPipelineParallelLaunch(t *testing.T) {
	cfg := scionConfig()
	cfg.TeamSize = "standard"
	out, err := templates.Render("pipeline", templates.SkillPipelineTemplate, cfg)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	mustContain := []string{
		"Launch ALL agents simultaneously",
		"scion start pm",
		"scion start tpm",
		"scion start swe-1",
		"scion start swe-test",
		"Orchestrator relays messages",
		"sciontool status task_completed",
		"scion message",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("pipeline output missing %q", s)
		}
	}
}

func TestScionClaudeMDParallelOrchestration(t *testing.T) {
	for _, size := range []string{"lean", "standard", "full"} {
		t.Run(size, func(t *testing.T) {
			cfg := scionConfig()
			cfg.TeamSize = size
			out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			mustContain := []string{
				"Agent Orchestration via Scion",
				"Launch all agents simultaneously",
				"Agents signal via",
				"sciontool status task_completed",
				"Orchestrator relays messages",
				"Monitor progress",
				"Stop all agents",
			}
			for _, s := range mustContain {
				if !strings.Contains(out, s) {
					t.Errorf("team size %q: output missing %q", size, s)
				}
			}
		})
	}
}

func TestClaudeCodeTemplatesUnaffectedByScion(t *testing.T) {
	cfg := fullConfig()
	cfg.Framework = "claude-code"
	out, err := templates.Render("CLAUDE.md", templates.ClaudeMDTemplate, cfg)
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	mustContain := []string{
		"Interactive Agent Teams via Tmux",
		"TeamCreate",
		"SendMessage",
		"shutdown_request",
		"TeamDelete",
	}
	for _, s := range mustContain {
		if !strings.Contains(out, s) {
			t.Errorf("claude-code CLAUDE.md missing %q", s)
		}
	}
	if strings.Contains(out, "Launch all agents simultaneously") {
		t.Error("claude-code CLAUDE.md should NOT contain Scion parallel launch instructions")
	}
}

// itoa converts int to string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	if neg {
		digits = append(digits, '-')
	}
	// reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
