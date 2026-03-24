package report

import (
	"fmt"
	"strings"

	"github.com/afeldman/project-check/internal/agent"
	"github.com/afeldman/project-check/internal/rules"
	"github.com/charmbracelet/lipgloss"
)

// Print displays findings in colored terminal output
func Print(findings []agent.Finding, rules rules.RuleSet) {
	if len(findings) == 0 {
		successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(successStyle.Render("✓ Project matches all standards"))
		return
	}

	// Define styles
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	messageStyle := lipgloss.NewStyle()

	// Create a map from rule ID to severity for quick lookup
	ruleSeverity := make(map[string]string)
	for _, rule := range rules.Rules {
		ruleSeverity[rule.ID] = rule.Severity
	}

	// Track statistics
	errorCount := 0
	warningCount := 0
	infoCount := 0
	files := make(map[string]bool)

	// Print each finding
	for _, finding := range findings {
		severity, exists := ruleSeverity[finding.RuleID]
		if !exists {
			// Default to warning if rule not found
			severity = "warning"
		}

		// Update counts
		switch severity {
		case "error":
			errorCount++
		case "warning":
			warningCount++
		case "info":
			infoCount++
		}

		// Track unique files
		if finding.File != "" {
			files[finding.File] = true
		}

		// Format the output line
		var severityLabel string
		var severityStyle lipgloss.Style

		switch severity {
		case "error":
			severityLabel = "ERROR"
			severityStyle = errorStyle
		case "warning":
			severityLabel = "WARN"
			severityStyle = warningStyle
		case "info":
			severityLabel = "INFO"
			severityStyle = infoStyle
		default:
			severityLabel = severity
			severityStyle = warningStyle
		}

		// Format the file location
		location := ""
		if finding.File != "" {
			location = fileStyle.Render(finding.File)
			if finding.Line > 0 {
				location += fmt.Sprintf(":%d", finding.Line)
			}
		}

		// Build the output line
		line := fmt.Sprintf("[%s] %s: %s — %s",
			severityStyle.Render(severityLabel),
			finding.RuleID,
			location,
			messageStyle.Render(finding.Message),
		)

		fmt.Println(line)
	}

	// Print summary
	if len(findings) > 0 {
		fmt.Println()
		
		summaryParts := []string{}
		if errorCount > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d error%s", errorCount, pluralize(errorCount)))
		}
		if warningCount > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d warning%s", warningCount, pluralize(warningCount)))
		}
		if infoCount > 0 {
			summaryParts = append(summaryParts, fmt.Sprintf("%d info%s", infoCount, pluralize(infoCount)))
		}
		
		summary := strings.Join(summaryParts, ", ")
		if len(files) > 0 {
			summary += fmt.Sprintf(" in %d file%s", len(files), pluralize(len(files)))
		}
		
		fmt.Println(summary)
	}
}

// pluralize returns "s" if count != 1, otherwise empty string
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}
