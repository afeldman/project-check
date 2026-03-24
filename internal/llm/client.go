package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type Client struct{ cfg Config; http *http.Client }

func New(cfg Config) *Client {
	if cfg.TimeoutS == 0 { cfg.TimeoutS = 60 }
	return &Client{cfg: cfg, http: &http.Client{Timeout: time.Duration(cfg.TimeoutS) * time.Second}}
}

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
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

type ChatResponse struct{ Content string; ToolCalls []ToolCall }

type openAIRequest struct {
	Model      string        `json:"model"`
	Messages   []ChatMessage `json:"messages"`
	Tools      []Tool        `json:"tools,omitempty"`
	ToolChoice string        `json:"tool_choice,omitempty"`
	Stream     bool          `json:"stream"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
	} `json:"choices"`
	Error struct{ Message string } `json:"error,omitempty"`
}

func (c *Client) Analyze(ctx context.Context, systemPrompt, userContent string) (string, error) {
	messages := []ChatMessage{{Role: "system", Content: systemPrompt}, {Role: "user", Content: userContent}}
	resp, err := c.Chat(ctx, messages, nil)
	if err != nil { return "", err }
	if resp.Content == "" { return "", fmt.Errorf("LLM returned empty content") }
	return resp.Content, nil
}

func (c *Client) Chat(ctx context.Context, messages []ChatMessage, tools []Tool) (*ChatResponse, error) {
	if !c.cfg.Enabled { return nil, fmt.Errorf("LLM client is disabled") }
	if c.cfg.Endpoint == "" { return nil, fmt.Errorf("LLM endpoint is not configured") }

	reqBody := openAIRequest{Model: c.cfg.Model, Messages: messages, Tools: tools, Stream: false}
	if tools != nil { reqBody.ToolChoice = "auto" }

	jsonBody, err := json.Marshal(reqBody)
	if err != nil { return nil, fmt.Errorf("failed to marshal request: %w", err) }

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.Endpoint+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil { return nil, fmt.Errorf("failed to create request: %w", err) }
	req.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" { req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey) }

	resp, err := c.http.Do(req)
	if err != nil { return nil, fmt.Errorf("HTTP request failed: %w", err) }
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp openAIResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.Error.Message != "" {
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp.Error.Message)
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if len(openAIResp.Choices) == 0 { return nil, fmt.Errorf("no choices in response") }

	return &ChatResponse{
		Content:   openAIResp.Choices[0].Message.Content,
		ToolCalls: openAIResp.Choices[0].Message.ToolCalls,
	}, nil
}
