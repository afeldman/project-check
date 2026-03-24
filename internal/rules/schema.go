package rules

import "fmt"

type RuleSet struct {
	Version int    `yaml:"version"`
	Rules   []Rule `yaml:"rules"`
}

type Rule struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Severity    string   `yaml:"severity"` // error | warning | info
	Languages   []string `yaml:"languages"`
	FixHint     string   `yaml:"fix_hint,omitempty"`
}

var LanguageExtensions = map[string][]string{
	"go":         {".go"},
	"python":     {".py"},
	"bash":       {".sh", ".zsh", ".bash"},
	"typescript": {".ts", ".tsx"},
}

func Validate(rs RuleSet) error {
	if rs.Version != 1 {
		return fmt.Errorf("invalid version: expected 1, got %d", rs.Version)
	}

	for i, rule := range rs.Rules {
		if rule.ID == "" {
			return fmt.Errorf("rule %d: missing id", i)
		}
		if rule.Description == "" {
			return fmt.Errorf("rule %d: missing description", i)
		}
		if rule.Severity == "" {
			return fmt.Errorf("rule %d: missing severity", i)
		}
		if rule.Severity != "error" && rule.Severity != "warning" && rule.Severity != "info" {
			return fmt.Errorf("rule %d: invalid severity %q, must be one of: error, warning, info", i, rule.Severity)
		}
		if len(rule.Languages) == 0 {
			return fmt.Errorf("rule %d: languages cannot be empty", i)
		}
		for _, lang := range rule.Languages {
			if _, ok := LanguageExtensions[lang]; !ok {
				return fmt.Errorf("rule %d: unknown language %q, must be one of: go, python, bash, typescript", i, lang)
			}
		}
	}
	return nil
}
