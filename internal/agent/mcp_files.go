package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// handleTree handles the tree tool
func handleTree(args string) string {
	var params struct {
		Path     string `json:"path"`
		MaxDepth int    `json:"max_depth"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	if params.MaxDepth == 0 {
		params.MaxDepth = 5
	}
	var result strings.Builder
	err := walkDir(params.Path, "", 0, params.MaxDepth, &result)
	if err != nil {
		return fmt.Sprintf("Error walking directory: %v", err)
	}
	return result.String()
}

// walkDir recursively walks a directory and builds a tree string
func walkDir(path, prefix string, depth, maxDepth int, result *strings.Builder) error {
	if depth > maxDepth {
		return nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for i, entry := range entries {
		isLast := i == len(entries)-1
		var connector string
		if depth == 0 {
			connector = ""
		} else if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}
		result.WriteString(prefix + connector + entry.Name() + "\n")
		if entry.IsDir() {
			var newPrefix string
			if depth == 0 {
				newPrefix = ""
			} else if isLast {
				newPrefix = prefix + "    "
			} else {
				newPrefix = prefix + "│   "
			}
			fullPath := filepath.Join(path, entry.Name())
			if err := walkDir(fullPath, newPrefix, depth+1, maxDepth, result); err != nil {
				return err
			}
		}
	}
	return nil
}

// handleReadFile handles the read_file tool
func handleReadFile(args string) string {
	var params struct {
		Path     string `json:"path"`
		MaxLines int    `json:"max_lines"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	if params.MaxLines == 0 {
		params.MaxLines = 500
	}
	content, err := os.ReadFile(params.Path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) > params.MaxLines {
		lines = lines[:params.MaxLines]
		lines = append(lines, fmt.Sprintf("... (truncated, file has more than %d lines)", params.MaxLines))
	}
	var result strings.Builder
	for i, line := range lines {
		result.WriteString(fmt.Sprintf("%4d: %s\n", i+1, line))
	}
	return result.String()
}

// handleListFiles handles the list_files tool
func handleListFiles(args string) string {
	var params struct {
		Pattern string `json:"pattern"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	matches, err := filepath.Glob(params.Pattern)
	if err != nil {
		return fmt.Sprintf("Error matching pattern: %v", err)
	}
	return strings.Join(matches, "\n")
}

// handleWriteFile handles the write_file tool
func handleWriteFile(args string) string {
	var params struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	dir := filepath.Dir(params.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}
	if err := os.WriteFile(params.Path, []byte(params.Content), 0644); err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}
	return fmt.Sprintf("Successfully wrote to %s", params.Path)
}

// handleWriteFileDryRun handles the write_file tool in dry-run mode:
// prints the diff to stdout but does not modify the file.
func handleWriteFileDryRun(args string) string {
	var params struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return fmt.Sprintf("Error parsing arguments: %v", err)
	}
	existing, err := os.ReadFile(params.Path)
	var oldContent string
	if err == nil {
		oldContent = string(existing)
	}
	diff := unifiedDiff(params.Path, oldContent, params.Content)
	fmt.Print(diff)
	return fmt.Sprintf("Dry-run: showed diff for %s (file not modified)", params.Path)
}
