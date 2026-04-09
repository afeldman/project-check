package cmd

import (
	"fmt"
	"os"

	"github.com/afeldman/project-check/internal/rules"
	"github.com/spf13/cobra"
)

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Translate STANDARDS.md to rules.yaml",
	Long:  `Translate a human-readable STANDARDS.md file to a machine-readable rules.yaml file using LLM.`,
	Run:   runTranslate,
}

var (
	standardsPath string
	outPath       string
)

func init() {
	translateCmd.Flags().StringVar(&standardsPath, "standards", "", "path to STANDARDS.md (required)")
	translateCmd.MarkFlagRequired("standards")
	translateCmd.Flags().StringVar(&outPath, "out", "rules.yaml", "output path for rules.yaml")

	rootCmd.AddCommand(translateCmd)
}

func runTranslate(cmd *cobra.Command, args []string) {
	if err := rules.Translate(getLLMConfig(), standardsPath, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}
}
