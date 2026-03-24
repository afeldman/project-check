# Task A: Project Scaffold + go.mod + main.go

## Goal
Create the basic Go project structure for `project-check` in `/Users/anton.feldmann/Projects/priv/pkg/project-check/`.

## Spec Reference
Full spec: `/Users/anton.feldmann/Projects/priv/pkg/docs/superpowers/specs/2026-03-24-project-check-design.md`

## What to create

### go.mod
Module: `github.com/afeldman/project-check`
Go version: 1.24

Dependencies:
- `github.com/spf13/cobra`
- `github.com/charmbracelet/lipgloss`
- `github.com/charmbracelet/bubbles`
- `gopkg.in/yaml.v3`
- `github.com/sashabaranov/go-openai`
- `github.com/mark3labs/mcp-go`
- `github.com/BurntSushi/toml`

Run `go mod tidy` after creating go.mod.

### main.go
```go
package main

import "github.com/afeldman/project-check/cmd"

func main() {
    cmd.Execute()
}
```

### cmd/root.go
Cobra root command. Loads config from `~/.config/project-check/config.toml` using BurntSushi/toml.
Config struct:
```go
type AppConfig struct {
    LLM llm.Config `toml:"llm"`
}
```
Persistent flags: none (per-command flags only).
Version: injected via ldflags `var Version = "dev"`.

### cmd/translate.go
Cobra sub-command `translate`.
Flags:
- `--standards` (string, required): path to STANDARDS.md
- `--out` (string, default: `rules.yaml`): output path

Calls `internal/rules/Translate(cfg llm.Config, standardsPath, outPath string) error`.
On error: print error + exit code 2.

### cmd/check.go
Cobra sub-command `check`.
Flags:
- `--rules` (string, default: `rules.yaml`): path to rules.yaml
- `--dir` (string, default: `.`): project directory to check
- `--fix` (bool): enable auto-fix mode
- `--dry-run` (bool): show diff without writing (only meaningful with --fix)
- `--sarif` (string): path for SARIF output file

Calls `internal/agent/Run(cfg AgentConfig) ([]Finding, error)`.
On errors: exit code 2.
On findings with severity=error: exit code 1.
On clean: exit code 0.

## Constraints
- Max ~150 lines per file
- No global state
- All files must compile: `go build ./...`
