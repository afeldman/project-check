# Task E: internal/report — Terminal + SARIF Output

## Goal
Implement the report package: colored terminal output and SARIF v2.1.0 generator.

## Spec Reference
Full spec: `/Users/anton.feldmann/Projects/priv/pkg/docs/superpowers/specs/2026-03-24-project-check-design.md`

## Files to create

### internal/report/terminal.go
Uses `github.com/charmbracelet/lipgloss` for colors.

```go
func Print(findings []agent.Finding, rules rules.RuleSet)
```

Output format per finding:
```
[ERROR] file-size-limit: internal/server/handler.go:312 — File has 312 lines (max 150)
[WARN]  no-global-state: pkg/cache/cache.go:14 — Global mutable variable 'globalCache'
```

Colors:
- ERROR: red bold
- WARN: yellow
- INFO: cyan
- File path: dim/gray
- Clean message (0 findings): green "✓ Project matches all standards"

Summary line at end:
```
3 errors, 1 warning in 4 files
```

### internal/report/sarif.go
SARIF v2.1.0 output.

SARIF field mapping (from spec):
| Finding field | SARIF field |
|---|---|
| rule_id | runs[].results[].ruleId |
| severity=error | level: "error" |
| severity=warning | level: "warning" |
| severity=info | level: "note" |
| file | locations[].physicalLocation.artifactLocation.uri |
| line | locations[].physicalLocation.region.startLine |
| message | message.text |

```go
func WriteSARIF(path string, findings []agent.Finding, rules rules.RuleSet) error
```

SARIF structure:
```go
type sarifReport struct {
    Schema  string    `json:"$schema"`
    Version string    `json:"version"`
    Runs    []sarifRun `json:"runs"`
}
// version: "2.1.0"
// $schema: "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json"
```

Populate `runs[0].tool.driver.rules` from the full rule set (not just violated rules).
Use `json.MarshalIndent` with 2-space indent.

## Constraints
- `go build ./...` must pass
- Max 150 lines per file
