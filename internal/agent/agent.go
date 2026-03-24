package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/afeldman/project-check/internal/llm"
	"github.com/afeldman/project-check/internal/rules"
)

// Config holds agent configuration
type Config struct {
	LLM     *llm.Client
	Rules   rules.RuleSet
	Dir     string
	FixMode bool
	DryRun  bool
}

// Finding represents a single finding
type Finding struct {
	RuleID  string `json:"rule_id"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// Run executes the ReAct loop
func Run(cfg Config) ([]Finding, error) {
	// Build tools and handlers
	tools, handlers := BuildTools(cfg.FixMode)
	
	// Build system prompt
	systemPrompt := BuildSystemPrompt(cfg.Rules, cfg.Dir, cfg.FixMode)
	
	// Initialize messages
	messages := []llm.ChatMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
		{
			Role:    "user",
			Content: "Analyze the project.",
		},
	}
	
	var findings []Finding
	maxIterations := 20
	
	for i := 0; i < maxIterations; i++ {
		// Call LLM
		response, err := cfg.LLM.Chat(context.Background(), messages, tools)
		if err != nil {
			return nil, fmt.Errorf("LLM chat failed: %w", err)
		}
		
		// Check for tool calls
		if len(response.ToolCalls) > 0 {
			// Append assistant message with tool calls
			messages = append(messages, llm.ChatMessage{
				Role:      "assistant",
				Content:   response.Content,
				ToolCalls: response.ToolCalls,
			})
			
			// Execute each tool call
			for _, toolCall := range response.ToolCalls {
				handler, ok := handlers[toolCall.Function.Name]
				if !ok {
					// Tool not found - return error result
					messages = append(messages, llm.ChatMessage{
						Role:       "tool",
						Content:    fmt.Sprintf("Error: tool '%s' not found", toolCall.Function.Name),
						ToolCallID: toolCall.ID,
					})
					continue
				}
				
				// Execute handler
				result := handler(toolCall.Function.Arguments)
				
				// Append tool result
				messages = append(messages, llm.ChatMessage{
					Role:       "tool",
					Content:    result,
					ToolCallID: toolCall.ID,
				})
			}
			
			continue
		}
		
		// Check for findings JSON
		content := strings.TrimSpace(response.Content)
		if strings.Contains(content, `{"findings":`) {
			// Try to parse JSON
			var result struct {
				Findings []Finding `json:"findings"`
			}
			
			if err := json.Unmarshal([]byte(content), &result); err != nil {
				// Not valid JSON - append as assistant message and continue
				messages = append(messages, llm.ChatMessage{
					Role:    "assistant",
					Content: response.Content,
				})
				continue
			}
			
			findings = result.Findings
			
			// If in dry-run fix mode, show diffs
			if cfg.DryRun && cfg.FixMode && len(findings) > 0 {
				// Note: In a real implementation, we would need to track
				// which files were modified and show diffs
				fmt.Fprintf(os.Stderr, "Dry run mode: would apply fixes for %d findings\n", len(findings))
			}
			
			return findings, nil
		}
		
		// No tool calls and no findings JSON - append as assistant message
		messages = append(messages, llm.ChatMessage{
			Role:    "assistant",
			Content: response.Content,
		})
	}
	
	// Max iterations reached
	fmt.Fprintf(os.Stderr, "Warning: Reached maximum iterations (%d). Returning partial findings.\n", maxIterations)
	return findings, nil
}
