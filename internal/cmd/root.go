package cmd

import (
	"github.com/hulk510/readme-gen/internal/i18n"
	"github.com/spf13/cobra"
)

var langFlag string

var rootCmd = &cobra.Command{
	Use:   "readme-gen",
	Short: "A CLI tool for managing README.md with structure auto-sync",
	Long: `readme-gen helps you maintain your README.md with automatic
directory structure updates and consistent templates.

Features:
  - Multiple templates (oss, personal, team)
  - Marker-based structure auto-sync
  - Claude Code integration support
  - CI-friendly check command`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set language
		switch langFlag {
		case "ja", "jp", "japanese":
			i18n.SetLanguage(i18n.Japanese)
		case "en", "english":
			i18n.SetLanguage(i18n.English)
		default:
			// Auto-detect from environment
			i18n.SetLanguage(i18n.DetectLanguage())
		}
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "Language (en, ja)")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(structureCmd)
	rootCmd.AddCommand(checkCmd)
}
