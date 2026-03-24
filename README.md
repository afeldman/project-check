# project-check

LLM-powered project standards checker. Translates a colloquial `STANDARDS.md` into a machine-readable `rules.yaml`, then checks a project directory against those rules using an LLM agent with MCP tools.

## Install

```bash
gofish install project-check
```

Or build from source:

```bash
git clone https://github.com/afeldman/project-check
cd project-check
go build -o project-check .
```

## Usage

### Step 1: Translate standards to rules

```bash
project-check translate --standards STANDARDS.md --out rules.yaml
```

Sends `STANDARDS.md` to the configured LLM and generates a structured `rules.yaml`.

### Step 2: Check a project

```bash
project-check check --rules rules.yaml --dir /path/to/project
```

```
[ERROR] file-size-limit: internal/server/handler.go:312 — 312 lines (max 150)
[WARN]  no-global-state: pkg/cache/cache.go:14 — global mutable variable
✗ 1 error, 1 warning
```

Exit code 1 on errors, 0 on clean.

### Options

```bash
project-check check --rules rules.yaml [--dir .] [--fix] [--dry-run] [--sarif out.sarif]
```

| Flag | Description |
|------|-------------|
| `--fix` | Auto-rewrite files to comply (shows confirmation prompt) |
| `--dry-run` | With `--fix`: show diff without writing |
| `--sarif` | Write SARIF v2.1.0 output for GitHub Code Scanning |

## Configuration

`~/.config/project-check/config.toml`:

```toml
[llm]
enabled   = true
endpoint  = "http://localhost:11434/v1"  # ollama default
model     = "llama3.2"
api_key   = ""                           # set for cloud APIs (OpenAI, etc.)
timeout_s = 60
```

Works with Ollama, LM Studio, and any OpenAI-compatible API.

## rules.yaml Format

```yaml
version: 1
rules:
  - id: "file-size-limit"
    description: "Max 150 lines per file"
    severity: warning          # error | warning | info
    languages: [go, python, typescript, bash]
    fix_hint: "Split into smaller modules."

  - id: "no-global-state"
    description: "No global mutable state"
    severity: error
    languages: [go]
```

Supported languages: `go`, `python`, `bash`, `typescript`

## License

MIT
