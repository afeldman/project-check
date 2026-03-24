package agent

import (
	"fmt"
	"strings"

	"github.com/afeldman/project-check/internal/rules"
	"gopkg.in/yaml.v3"
)

// BuildSystemPrompt builds the system prompt from rules and mode
func BuildSystemPrompt(ruleSet rules.RuleSet, dir string, fixMode bool) string {
	// Format rules as YAML
	rulesYAML, err := yaml.Marshal(ruleSet.Rules)
	if err != nil {
		// Fallback to simple formatting if YAML marshaling fails
		var rulesText strings.Builder
		for _, rule := range ruleSet.Rules {
			rulesText.WriteString(fmt.Sprintf("- id: %s\n", rule.ID))
			rulesText.WriteString(fmt.Sprintf("  description: %s\n", rule.Description))
			rulesText.WriteString(fmt.Sprintf("  severity: %s\n", rule.Severity))
			if len(rule.Languages) > 0 {
				rulesText.WriteString(fmt.Sprintf("  languages: %v\n", rule.Languages))
			}
			if rule.FixHint != "" {
				rulesText.WriteString(fmt.Sprintf("  fix_hint: %s\n", rule.FixHint))
			}
			rulesText.WriteString("\n")
		}
		rulesYAML = []byte(rulesText.String())
	}

	mode := "check"
	if fixMode {
		mode = "fix — rewrite files to make them comply"
	}

	prompt := fmt.Sprintf(`You are a project standards checker. Your job is to analyze the project at %s and check it
against the following rules. Use the provided tools to explore the codebase.

When done, respond with a JSON object ONLY (no other text):
{"findings": [{"rule_id": "...", "file": "relative/path.go", "line": 42, "message": "..."}]}

If the project is clean, respond with: {"findings": []}

Rules:
%s

Mode: %s`, dir, string(rulesYAML), mode)

	return prompt
}
