package cmd

import (
	"fmt"
	"os"

	"github.com/hulk510/readme-gen/internal/i18n"
	"github.com/hulk510/readme-gen/internal/marker"
	"github.com/hulk510/readme-gen/internal/scanner"
	"github.com/hulk510/readme-gen/internal/ui"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if README structure is up to date",
	Long:  `Verify that the directory structure in README.md matches the current state. Exits with code 1 if out of sync.`,
	RunE:  runCheck,
}

func runCheck(cmd *cobra.Command, args []string) error {
	msg := i18n.Get()

	// Read current README
	content, err := os.ReadFile("README.md")
	if err != nil {
		fmt.Println(ui.Err(msg.ReadmeNotFound))
		return fmt.Errorf("%s", msg.ReadmeNotFound)
	}

	// Extract current structure from README
	readmeStructure, found := marker.Extract(string(content))
	if !found {
		fmt.Println(ui.Warn(msg.NoMarkersFound))
		fmt.Println(ui.Info(msg.AddMarkersHint))
		return nil
	}

	// Scan current directory
	currentStructure, err := scanner.Scan(".", scanner.DefaultExcludes())
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	// Compare
	if readmeStructure == currentStructure {
		fmt.Println(ui.Check(msg.StructureUpToDate))
		return nil
	}

	// Out of sync
	fmt.Println(ui.Warn(msg.StructureOutOfSync))
	fmt.Println()
	fmt.Println(ui.Info(msg.RunUpdateHint))

	// Return error for CI
	os.Exit(1)
	return nil
}
