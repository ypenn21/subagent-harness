package templates

// ScionAgentYAMLData holds data for rendering the scion-agent.yaml template.
type ScionAgentYAMLData struct {
	Name        string
	Description string
	Harness     string
}

const ScionAgentYAMLTemplate = `schema_version: "1"
description: "{{.Description}}"
agent_instructions: agents.md
system_prompt: system-prompt.md
default_harness_config: {{.Harness}}
`
