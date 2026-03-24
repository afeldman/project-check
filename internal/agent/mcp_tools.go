package agent

import (
	"github.com/afeldman/project-check/internal/llm"
)

// ToolHandler is a function that handles a tool call
type ToolHandler func(args string) string

// BuildTools builds the list of tools and their handlers
func BuildTools(fixMode bool) ([]llm.Tool, map[string]ToolHandler) {
	tools := []llm.Tool{
		{
			Type: "function",
			Function: llm.ToolFunction{
				Name:        "tree",
				Description: "Walk directory recursively (max depth 5) and return indented tree string",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Directory path to walk",
						},
						"max_depth": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum depth to traverse (default: 5)",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.ToolFunction{
				Name:        "read_file",
				Description: "Read file content (max 500 lines) and return as string with line numbers",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "File path to read",
						},
						"max_lines": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum lines to read (default: 500)",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.ToolFunction{
				Name:        "list_files",
				Description: "List files matching a glob pattern, return newline-separated paths",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"pattern": map[string]interface{}{
							"type":        "string",
							"description": "Glob pattern to match files",
						},
					},
					"required": []string{"pattern"},
				},
			},
		},
		{
			Type: "function",
			Function: llm.ToolFunction{
				Name:        "run_command",
				Description: "Run shell command (allowed commands: wc, head, tail, cat, grep, find, ls, go, golangci-lint) with 10s timeout",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]interface{}{
							"type":        "string",
							"description": "Shell command to execute",
						},
					},
					"required": []string{"command"},
				},
			},
		},
	}

	handlers := map[string]ToolHandler{
		"tree":       handleTree,
		"read_file":  handleReadFile,
		"list_files": handleListFiles,
		"run_command": handleRunCommand,
	}

	// Add write_file tool only in fix mode
	if fixMode {
		tools = append(tools, llm.Tool{
			Type: "function",
			Function: llm.ToolFunction{
				Name:        "write_file",
				Description: "Write content to file path",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "File path to write",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "Content to write to file",
						},
					},
					"required": []string{"path", "content"},
				},
			},
		})
		handlers["write_file"] = handleWriteFile
	}

	return tools, handlers
}
