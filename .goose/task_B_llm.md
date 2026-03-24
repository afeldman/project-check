# Task B: internal/llm — LLM Client (minimal, no daemon)

## Goal
Implement a minimal LLM client package in `internal/llm/client.go`.
Single file only — no subprocess, no Python daemon, no daemon_unix/windows.
Ollama and LM Studio already provide an OpenAI-compatible HTTP endpoint.

## File to create: internal/llm/client.go

```go
package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Config struct {
    Enabled  bool   `toml:"enabled"`
    Endpoint string `toml:"endpoint"`
    Model    string `toml:"model"`
    APIKey   string `toml:"api_key"`
    TimeoutS int    `toml:"timeout_s"`
}

type Client struct {
    cfg  Config
    http *http.Client
}

func New(cfg Config) *Client {
    if cfg.TimeoutS == 0 {
        cfg.TimeoutS = 60
    }
    return &Client{cfg: cfg, http: &http.Client{Timeout: time.Duration(cfg.TimeoutS) * time.Second}}
}
```

### ChatMessage, ToolCall, Tool, ChatResponse types

```go
type ChatMessage struct {
    Role       string     `json:"role"`
    Content    string     `json:"content,omitempty"`
    ToolCallID string     `json:"tool_call_id,omitempty"`
    ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
    Name       string     `json:"name,omitempty"`
}

type ToolCall struct {
    ID       string `json:"id"`
    Type     string `json:"type"`
    Function struct {
        Name      string `json:"name"`
        Arguments string `json:"arguments"`
    } `json:"function"`
}

type Tool struct {
    Type     string       `json:"type"` // always "function"
    Function ToolFunction `json:"function"`
}

type ToolFunction struct {
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Parameters  interface{} `json:"parameters"`
}

type ChatResponse struct {
    Content   string     // non-empty when LLM returns text
    ToolCalls []ToolCall // non-empty when LLM wants to call tools
}
```

### Methods

```go
// Analyze: single-turn, system + user prompt → response string
func (c *Client) Analyze(ctx context.Context, systemPrompt, userContent string) (string, error)

// Chat: multi-turn with tool_calls support
func (c *Client) Chat(ctx context.Context, messages []ChatMessage, tools []Tool) (*ChatResponse, error)
```

Both methods:
- POST to `c.cfg.Endpoint + "/chat/completions"`
- Set `Authorization: Bearer <api_key>` if api_key != ""
- On HTTP != 200: return error with status + body
- On empty choices: return error

For `Chat`: use OpenAI request format:
```json
{
  "model": "...",
  "messages": [...],
  "tools": [...],
  "tool_choice": "auto",
  "stream": false
}
```

Parse response: if `choices[0].message.tool_calls` present → fill ChatResponse.ToolCalls;
else → fill ChatResponse.Content.

## Constraints
- Single file, max 150 lines
- No CGO
- `go build ./...` must pass
