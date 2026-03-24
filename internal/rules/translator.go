package rules

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/afeldman/project-check/internal/llm"
	"gopkg.in/yaml.v3"
)

func Translate(cfg llm.Config, standardsPath, outPath string) error {
	// Read standards file
	content, err := os.ReadFile(standardsPath)
	if err != nil {
		return fmt.Errorf("failed to read standards file %s: %w", standardsPath, err)
	}

	// Create LLM client
	client := llm.New(cfg)

	// Build prompt
	systemPrompt := `You are a code standards analyzer. Convert the following STANDARDS document into a YAML rule set.
Return ONLY valid YAML, no markdown fences, no explanation.

Schema:
version: 1
rules:
  - id: "unique-kebab-case-id"
    description: "Clear description of the rule"
    severity: "error"  # error | warning | info
    languages: ["go"]  # subset of: go, python, bash, typescript
    fix_hint: "Optional: how to fix this"  # omit if not applicable`

	userContent := fmt.Sprintf("STANDARDS document:\n%s", string(content))

	// Call LLM
	ctx := context.Background()
	yamlContent, err := client.Analyze(ctx, systemPrompt, userContent)
	if err != nil {
		return fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Strip markdown fences
	yamlContent = strings.TrimSpace(yamlContent)
	if strings.HasPrefix(yamlContent, "```yaml") {
		yamlContent = strings.TrimPrefix(yamlContent, "```yaml")
		yamlContent = strings.TrimSuffix(yamlContent, "```")
	} else if strings.HasPrefix(yamlContent, "```") {
		yamlContent = strings.TrimPrefix(yamlContent, "```")
		yamlContent = strings.TrimSuffix(yamlContent, "```")
	}
	yamlContent = strings.TrimSpace(yamlContent)

	// Parse YAML
	var rs RuleSet
	if err := yaml.Unmarshal([]byte(yamlContent), &rs); err != nil {
		return fmt.Errorf("failed to parse LLM response as YAML: %w\nLLM response was:\n%s", err, yamlContent)
	}

	// Validate
	if err := Validate(rs); err != nil {
		return fmt.Errorf("LLM generated invalid rules: %w\nCheck LLM response", err)
	}

	// Save
	if err := Save(outPath, rs); err != nil {
		return fmt.Errorf("failed to save rules: %w", err)
	}

	return nil
}
