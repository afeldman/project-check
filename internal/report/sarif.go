package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/afeldman/project-check/internal/agent"
	"github.com/afeldman/project-check/internal/rules"
)

// WriteSARIF generates a SARIF v2.1.0 report file
func WriteSARIF(path string, findings []agent.Finding, rules rules.RuleSet) error {
	// Create SARIF rules from the rule set
	sarifRules := make([]sarifRule, 0, len(rules.Rules))
	for _, rule := range rules.Rules {
		// Map severity to SARIF default level
		var defaultLevel string
		switch rule.Severity {
		case "error":
			defaultLevel = "error"
		case "warning":
			defaultLevel = "warning"
		case "info":
			defaultLevel = "note"
		default:
			defaultLevel = "warning"
		}

		sarifRules = append(sarifRules, sarifRule{
			ID:   rule.ID,
			Name: rule.ID,
			ShortDescription: sarifMessage{
				Text: rule.Description,
			},
			FullDescription: sarifMessage{
				Text: rule.Description,
			},
			DefaultLevel: defaultLevel,
			Properties: map[string]string{
				"languages": fmt.Sprintf("%v", rule.Languages),
			},
		})
	}

	// Create SARIF results from findings
	results := make([]sarifResult, 0, len(findings))
	
	// Create a map from rule ID to severity for quick lookup
	ruleSeverity := make(map[string]string)
	for _, rule := range rules.Rules {
		ruleSeverity[rule.ID] = rule.Severity
	}
	
	for _, finding := range findings {
		severity, exists := ruleSeverity[finding.RuleID]
		if !exists {
			severity = "warning"
		}
		
		// Map severity to SARIF level
		var level string
		switch severity {
		case "error":
			level = "error"
		case "warning":
			level = "warning"
		case "info":
			level = "note"
		default:
			level = "warning"
		}
		
		// Create location
		location := sarifLocation{
			PhysicalLocation: sarifPhysicalLocation{
				ArtifactLocation: sarifArtifactLocation{
					URI: finding.File,
				},
			},
		}
		
		// Add region if line number is provided
		if finding.Line > 0 {
			location.PhysicalLocation.Region = sarifRegion{
				StartLine: finding.Line,
			}
		}
		
		result := sarifResult{
			RuleID: finding.RuleID,
			Level:  level,
			Message: sarifMessage{
				Text: finding.Message,
			},
			Locations: []sarifLocation{location},
		}
		
		results = append(results, result)
	}

	// Build the complete SARIF report
	report := sarifReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           "project-check",
						Version:        "1.0.0", // TODO: Get actual version
						InformationURI: "https://github.com/afeldman/project-check",
						Rules:          sarifRules,
					},
				},
				Results: results,
			},
		},
	}

	// Marshal with indentation
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SARIF report: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write SARIF file: %w", err)
	}

	return nil
}
