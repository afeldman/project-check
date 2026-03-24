package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/afeldman/project-check/internal/llm"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	Version = "dev"
)

type AppConfig struct {
	LLM llm.Config `toml:"llm"`
}

var appConfig AppConfig

var rootCmd = &cobra.Command{
	Use:     "project-check",
	Short:   "Check projects against company standards",
	Long:    `project-check is a CLI tool that checks projects against company standards.`,
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/project-check/config.toml)")
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			return
		}
		cfgFile = filepath.Join(home, ".config", "project-check", "config.toml")
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		// Config file doesn't exist, use defaults
		appConfig = AppConfig{
			LLM: llm.Config{
				Enabled:  true,
				Endpoint: "http://localhost:11434/v1",
				Model:    "llama3.2",
				TimeoutS: 60,
			},
		}
		return
	}

	var config AppConfig
	if _, err := toml.DecodeFile(cfgFile, &config); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file %s: %v\n", cfgFile, err)
		// Use defaults on error
		appConfig = AppConfig{
			LLM: llm.Config{
				Enabled:  true,
				Endpoint: "http://localhost:11434/v1",
				Model:    "llama3.2",
				TimeoutS: 60,
			},
		}
		return
	}

	appConfig = config
}

// getLLMConfig returns the LLM configuration
func getLLMConfig() llm.Config {
	return appConfig.LLM
}
