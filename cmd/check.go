package cmd

import (
	"fmt"
	"os"

	"github.com/afeldman/project-check/internal/agent"
	"github.com/afeldman/project-check/internal/llm"
	"github.com/afeldman/project-check/internal/report"
	"github.com/afeldman/project-check/internal/rules"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check project against rules",
	Long:  `Check a project directory against rules.yaml using LLM agent with MCP tools.`,
	Run:   runCheck,
}

var (
	rulesPath string
	dirPath   string
	fixMode   bool
	dryRun    bool
	sarifPath string
)

func init() {
	checkCmd.Flags().StringVar(&rulesPath, "rules", "rules.yaml", "path to rules.yaml")
	checkCmd.Flags().StringVar(&dirPath, "dir", ".", "project directory to check")
	checkCmd.Flags().BoolVar(&fixMode, "fix", false, "enable auto-fix mode")
	checkCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show diff without writing (only meaningful with --fix)")
	checkCmd.Flags().StringVar(&sarifPath, "sarif", "", "path for SARIF output file")

	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) {
	// Load rules
	ruleSet, err := rules.Load(rulesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading rules: %v\n", err)
		fmt.Fprintf(os.Stderr, "Hint: Run 'project-check translate' to generate rules.yaml from STANDARDS.md\n")
		os.Exit(1)
	}

	// Load LLM config from root command
	// For now, use defaults since getLLMConfig() is not implemented
	llmCfg := llm.Config{
		Enabled:  true,
		Endpoint: "http://localhost:11434/v1",
		Model:    "llama3.2",
		TimeoutS: 60,
	}
	
	// Create LLM client
	client := llm.New(llmCfg)
	
	// Create agent config
	cfg := agent.Config{
		LLM:     client,
		Rules:   ruleSet,
		Dir:     dirPath,
		FixMode: fixMode,
		DryRun:  dryRun,
	}

	// Run agent
	findings, err := agent.Run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// Print terminal report
	report.Print(findings, ruleSet)

	// Generate SARIF output if requested
	if sarifPath != "" {
		if err := report.WriteSARIF(sarifPath, findings, ruleSet); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing SARIF file: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("SARIF report written to %s\n", sarifPath)
	}

	// Check for errors in findings
	hasErrors := false
	for _, finding := range findings {
		// Check if this finding is for an error-level rule
		for _, rule := range ruleSet.Rules {
			if rule.ID == finding.RuleID && rule.Severity == "error" {
				hasErrors = true
				break
			}
		}
		if hasErrors {
			break
		}
	}

	if hasErrors {
		os.Exit(1)
	}
	// Clean project, exit with 0
}
