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

var updateFlag bool

var structureCmd = &cobra.Command{
	Use:   "structure",
	Short: "Show or update directory structure",
	Long:  `Display current directory structure or update the structure section in README.md.`,
	RunE:  runStructure,
}

func init() {
	structureCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "Update README.md structure section")
}

func runStructure(cmd *cobra.Command, args []string) error {
	msg := i18n.Get()

	// Scan directory
	structure, err := scanner.Scan(".", scanner.DefaultExcludes())
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if !updateFlag {
		// Just print structure
		fmt.Println(structure)
		return nil
	}

	// Update README.md
	fmt.Println(ui.Title())
	fmt.Printf("%s %s\n", ui.IconSync, msg.UpdatingStructure)

	// Read current README
	content, err := os.ReadFile("README.md")
	if err != nil {
		return fmt.Errorf("%s. %s", msg.ReadmeNotFound, msg.RunInitHint)
	}

	// Get current structure from README
	oldStructure, found := marker.Extract(string(content))

	// Update markers
	newContent, err := marker.Update(string(content), structure)
	if err != nil {
		return fmt.Errorf("failed to update structure: %w", err)
	}

	// Show diff if there were changes (strip comments for comparison)
	if found && marker.StripComments(oldStructure) != structure {
		fmt.Println()
		fmt.Println(ui.Box(fmt.Sprintf("%s:\n\nOld:\n%s\n\nNew:\n%s", msg.ChangesDetected, oldStructure, structure)))
		fmt.Println()
	}

	// Write updated README
	if err := os.WriteFile("README.md", []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	fmt.Println(ui.Check(msg.UpdatedReadme))
	return nil
}
