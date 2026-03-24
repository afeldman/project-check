package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// allowedCommands is the list of allowed commands for run_command
var allowedCommands = []string{"wc", "head", "tail", "cat", "grep", "find", "ls", "go", "golangci-lint"}

// handleRunCommand handles the run_command tool
func handleRunCommand(args string) string {
	var params struct {
		Command string `json:"command"`
	}
	
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	
	// Parse first token of command
	parts := strings.Fields(params.Command)
	if len(parts) == 0 {
		return "Error: empty command"
	}
	
	firstToken := parts[0]
	
	// Check if command is allowed
	allowed := false
	for _, cmd := range allowedCommands {
		if cmd == firstToken {
			allowed = true
			break
		}
	}
	
	if !allowed {
		return fmt.Sprintf("Error: command '%s' is not in the allowed list. Allowed commands: %v", firstToken, allowedCommands)
	}
	
	// Execute command with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "sh", "-c", params.Command)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Error executing command: %v\nOutput: %s", err, output)
	}
	
	return string(output)
}
