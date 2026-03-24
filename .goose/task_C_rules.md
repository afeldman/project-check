# Task C: internal/rules — Schema, Loader, Translator

## Goal
Implement the rules package: YAML schema, loader/validator, and LLM-based translator.

## Spec Reference
Full spec: `/Users/anton.feldmann/Projects/priv/pkg/docs/superpowers/specs/2026-03-24-project-check-design.md`

## Files to create

### internal/rules/schema.go
Go structs for rules.yaml:

```go
package rules

type RuleSet struct {
    Version int    `yaml:"version"`
    Rules   []Rule `yaml:"rules"`
}

type Rule struct {
    ID          string   `yaml:"id"`
    Description string   `yaml:"description"`
    Severity    string   `yaml:"severity"`   // error | warning | info
    Languages   []string `yaml:"languages"`
    FixHint     string   `yaml:"fix_hint,omitempty"`
}
```

Extension mapping (exported map for agent use):
```go
var LanguageExtensions = map[string][]string{
    "go":         {".go"},
    "python":     {".py"},
    "bash":       {".sh", ".zsh", ".bash"},
    "typescript": {".ts", ".tsx"},
}
```

Validation: `Validate(rs RuleSet) error` — checks version==1, all rules have id/description/severity,
severity is one of error|warning|info, languages are all in LanguageExtensions keys.

### internal/rules/loader.go
- `Load(path string) (RuleSet, error)` — read YAML file + Validate
- `Save(path string, rs RuleSet) error` — write YAML file with version header comment

### internal/rules/translator.go
Translates STANDARDS.md → rules.yaml via LLM.

```go
func Translate(llmClient *llm.Client, standardsPath, outPath string) error
```

Steps:
1. Read standardsPath content
2. Build prompt:
   ```
   You are a code standards analyzer. Convert the following STANDARDS document into a YAML rule set.
   Return ONLY valid YAML, no markdown fences, no explanation.

   Schema:
   version: 1
   rules:
     - id: "unique-kebab-case-id"
       description: "Clear description of the rule"
       severity: "error"  # error | warning | info
       languages: ["go"]  # subset of: go, python, bash, typescript
       fix_hint: "Optional: how to fix this"  # omit if not applicable

   STANDARDS document:
   <content>
   ```
3. Call `llmClient.Analyze(ctx, systemPrompt, userContent)`
4. Strip markdown fences (```yaml ... ```)
5. yaml.Unmarshal into RuleSet
6. Validate
7. Save to outPath
8. On any error: return fmt.Errorf with context + "check LLM response" hint

## Constraints
- `go build ./...` must pass
- Max 150 lines per file
