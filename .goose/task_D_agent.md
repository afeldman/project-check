# Task D: internal/agent — ReAct Loop + MCP Tools

## Goal
Implement the LLM agent with tool-calling (ReAct pattern) and MCP tool handlers.

## Spec Reference
Full spec: `/Users/anton.feldmann/Projects/priv/pkg/docs/superpowers/specs/2026-03-24-project-check-design.md`

## Files to create

### internal/agent/mcp.go
Tool definitions and handlers. All tools take a JSON arguments string, return a string result.

**Tool list:**
- `tree`: Walk directory recursively (max depth 5), return indented tree string
- `read_file`: Read file content (max 500 lines), return as string with line numbers
- `list_files`: List files matching a glob pattern, return newline-separated paths
- `run_command`: Run shell command — FIRST check against allowlist:
  ```go
  var allowedCommands = []string{"wc", "head", "tail", "cat", "grep", "find", "ls", "go", "golangci-lint"}
  ```
  Parse first token of command, if not in list: return error string (do NOT exec).
  If in list: exec with 10s timeout, return stdout+stderr.
- `write_file`: Write content to file path. Only registered when fixMode=true.

Export: `func BuildTools(fixMode bool) ([]llm.Tool, map[string]ToolHandler)`
where `type ToolHandler func(args string) string`

### internal/agent/prompt.go
Builds the system prompt from rules + mode:

```go
func BuildSystemPrompt(rules rules.RuleSet, dir string, fixMode bool) string
```

Prompt structure:
```
You are a project standards checker. Your job is to analyze the project at <dir> and check it
against the following rules. Use the provided tools to explore the codebase.

When done, respond with a JSON object ONLY (no other text):
{"findings": [{"rule_id": "...", "file": "relative/path.go", "line": 42, "message": "..."}]}

If the project is clean, respond with: {"findings": []}

Rules:
<yaml-formatted rule list>

Mode: check  (or: fix — rewrite files to make them comply)
```

### internal/agent/agent.go
ReAct loop:

```go
type Config struct {
    LLM     *llm.Client
    Rules   rules.RuleSet
    Dir     string
    FixMode bool
    DryRun  bool
}

type Finding struct {
    RuleID  string `json:"rule_id"`
    File    string `json:"file"`
    Line    int    `json:"line"`
    Message string `json:"message"`
}

func Run(cfg Config) ([]Finding, error)
```

Loop algorithm:
1. Build tools + system prompt
2. Initialize messages: [{role:system, content:systemPrompt}, {role:user, content:"Analyze the project."}]
3. Loop (max 20 iterations to avoid infinite loops):
   a. Call `cfg.LLM.Chat(ctx, messages, tools)`
   b. If response has ToolCalls:
      - For each tool call: execute handler, append tool result message
      - Append assistant message with tool_calls
      - Continue loop
   c. If response.Content contains `{"findings":`:
      - Parse JSON findings
      - If DryRun + FixMode: show diff (read original, compare)
      - Return findings, nil
   d. If no tool calls and no findings JSON: append as assistant message, continue
4. If max iterations reached: return partial findings with warning logged to stderr

## Constraints
- `go build ./...` must pass
- Max 150 lines per file (split agent.go if needed)
- No global mutable state
