package config

// SWEConfig defines a single SWE agent's specialization.
type SWEConfig struct {
	Number  int
	Title   string
	Bullets []string
}

// GCPConfig holds optional GCP project details.
type GCPConfig struct {
	Enabled       bool
	ProjectID     string
	ProjectNumber string
	Organization  string
	Region        string
}

// CustomSkillConfig defines a user-created skill.
type CustomSkillConfig struct {
	Name        string
	Description string
}

// CustomAgentConfig defines a user-created agent role.
type CustomAgentConfig struct {
	Name         string   // kebab-case identifier, used as filename
	Title        string   // display title (e.g., "Frontend Engineer")
	Description  string   // what this agent does
	Instructions []string // specific bullet-point instructions
}

// ProjectConfig holds all configuration gathered by the wizard.
type ProjectConfig struct {
	ProjectName     string
	Description     string
	TechStack       string
	InitGit         bool
	CreateRepo      bool
	RepoURL         string
	GitHubOrg       string
	OwnerName       string
	OwnerEmail      string
	OwnerGitHub     string
	GCP             GCPConfig
	SWEs            []SWEConfig
	IncludePlatform bool
	IncludeReviewer bool
	IncludeSWEQA    bool
	IncludeSWETest  bool
	Conventions    []string
	SelectedSkills map[string]bool
	CustomSkills   []CustomSkillConfig
	CustomAgents   []CustomAgentConfig
	ModelName      string
	ModelID        string
	TargetDir      string
	TeamSize       string
	Framework      string
	DefaultHarness string
}
