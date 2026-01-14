package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/hulk510/readme-gen/internal/config"
	"github.com/hulk510/readme-gen/internal/i18n"
	"github.com/hulk510/readme-gen/internal/scanner"
	"github.com/hulk510/readme-gen/internal/template"
	"github.com/hulk510/readme-gen/internal/ui"
	"github.com/spf13/cobra"
)

var (
	templateFlag   string
	nonInteractive bool
	withSkills     bool
	withAI         bool
	noSkills       bool
	noAI           bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize README.md from template",
	Long:  `Create a new README.md file from a template with directory structure.`,
	RunE:  runInit,
}

func init() {
	initCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "Template to use (oss, general)")
	initCmd.Flags().BoolVarP(&nonInteractive, "yes", "y", false, "Non-interactive mode with defaults")
	initCmd.Flags().BoolVar(&withSkills, "with-skills", false, "Add Claude Code skills")
	initCmd.Flags().BoolVar(&withAI, "with-ai", false, "Generate descriptions with AI")
	initCmd.Flags().BoolVar(&noSkills, "no-skills", false, "Skip adding Claude Code skills")
	initCmd.Flags().BoolVar(&noAI, "no-ai", false, "Skip AI generation")
}

func runInit(cmd *cobra.Command, args []string) error {
	msg := i18n.Get()
	fmt.Println(ui.Title())

	// Check if README already exists
	if _, err := os.Stat("README.md"); err == nil {
		fmt.Println(ui.Warn("README.md already exists"))
		var overwrite bool
		if !nonInteractive {
			err := huh.NewConfirm().
				Title(msg.OverwriteConfirm).
				Value(&overwrite).
				Run()
			if err != nil {
				return err
			}
			if !overwrite {
				fmt.Println(ui.Info(msg.Cancelled))
				return nil
			}
		}
	}

	var (
		selectedTemplate string
		projectName      string
		selectedOptions  []string
	)

	// Detect project info
	info := scanner.DetectProjectInfo(".")
	projectName = info.Name

	if nonInteractive {
		// Use flags or defaults
		selectedTemplate = templateFlag
		if selectedTemplate == "" {
			selectedTemplate = "oss"
		}
		if withSkills || !noSkills {
			selectedOptions = append(selectedOptions, "skills")
		}
		if withAI && !noAI {
			selectedOptions = append(selectedOptions, "ai")
		}
	} else {
		// Interactive mode with steps
		var selectedLang string

		// Step 1: Language selection
		fmt.Println()
		fmt.Println(ui.Step(1, 4, msg.StepLanguage))
		langForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(msg.SelectLanguage).
					Options(
						huh.NewOption(msg.LangEnglish, "en"),
						huh.NewOption(msg.LangJapanese, "ja"),
					).
					Value(&selectedLang),
			),
		)
		if err := langForm.Run(); err != nil {
			return err
		}

		// Apply language selection
		if selectedLang == "ja" {
			i18n.SetLanguage(i18n.Japanese)
		} else {
			i18n.SetLanguage(i18n.English)
		}
		msg = i18n.Get() // Refresh messages

		// Step 2: Template selection
		fmt.Println()
		fmt.Println(ui.Step(2, 4, msg.StepTemplate))
		templateForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title(msg.SelectTemplate).
					Options(
						huh.NewOption(msg.TemplateOSS, "oss"),
						huh.NewOption(msg.TemplateGeneral, "general"),
					).
					Value(&selectedTemplate),
			),
		)
		if err := templateForm.Run(); err != nil {
			return err
		}

		// Step 3: Project info
		fmt.Println()
		fmt.Println(ui.Step(3, 4, msg.StepProjectInfo))
		infoForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title(msg.ProjectName).
					Value(&projectName).
					Placeholder(info.Name),
			),
		)
		if err := infoForm.Run(); err != nil {
			return err
		}

		// Step 4: Integrations
		fmt.Println()
		fmt.Println(ui.Step(4, 4, msg.StepIntegration))
		integrationForm := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title(msg.ClaudeCodeIntegration).
					Options(
						huh.NewOption(msg.OptionAddSkills, "skills").Selected(true),
						huh.NewOption(msg.OptionGenerateWithAI, "ai"),
					).
					Value(&selectedOptions),
			),
		)
		if err := integrationForm.Run(); err != nil {
			return err
		}
	}

	// Generate structure
	structure, err := scanner.ScanAuto(".")
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	// Generate README
	data := template.Data{
		ProjectName: projectName,
		Description: info.Description,
		Structure:   structure,
		Language:    info.Language,
		ModulePath:  info.ModulePath,
		Lang:        i18n.Current(),
	}

	content, err := template.Render(selectedTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Write README.md
	if err := os.WriteFile("README.md", []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}
	fmt.Println(ui.Success(msg.CreatedReadme))

	// Add Claude Code skills if requested
	if contains(selectedOptions, "skills") {
		skillsDir := filepath.Join(".claude", "skills")
		if err := os.MkdirAll(skillsDir, 0755); err != nil {
			return fmt.Errorf("failed to create skills directory: %w", err)
		}

		skillsContent := template.GetClaudeSkills()
		skillsPath := filepath.Join(skillsDir, "readme.md")
		if err := os.WriteFile(skillsPath, []byte(skillsContent), 0644); err != nil {
			return fmt.Errorf("failed to write skills file: %w", err)
		}
		fmt.Println(ui.Success(msg.CreatedSkills))
	}

	// Generate descriptions with AI if requested
	if contains(selectedOptions, "ai") {
		if err := generateWithAI(msg); err != nil {
			fmt.Println(ui.Warn(err.Error()))
		}
	}

	// Print helpful info
	fmt.Println()
	fmt.Println(ui.Box(msg.MarkersInfo + ":\n\n<!-- readme-gen:structure:start -->\n<!-- readme-gen:structure:end -->"))
	fmt.Println()
	fmt.Println(ui.Info(msg.RunLaterHint))

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func generateWithAI(msg i18n.Messages) error {
	// Check if claude command exists
	_, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("%s", msg.ClaudeCodeNotFound)
	}

	// Load config for timeout setting
	cfg, err := config.Load(".")
	if err != nil {
		cfg = config.Default()
	}
	aiTimeout := time.Duration(cfg.AI.GetTimeout()) * time.Second

	// Collect additional context
	additionalContext := collectProjectContext()

	// Build prompt based on language
	var prompt string
	if i18n.Current() == i18n.Japanese {
		prompt = fmt.Sprintf(`README.mdを以下の情報を基に更新してください。

%s

更新対象:
1. 冒頭のDescription - プロジェクトの目的を1-2文で簡潔に記載
2. Structureセクション - 各ディレクトリの横に簡潔な説明コメントを追加
   形式例:
   ├── cmd/           # CLIエントリーポイント
   ├── internal/      # 内部パッケージ
3. Usageセクション - 基本的な使い方を記載
4. Developmentセクション - 開発コマンドを記載（mise/make/npm scripts等から）

ルール:
- コードを読んで適切な内容を生成してください
- 簡潔で実用的な内容にしてください
- 既存のマーカー（<!-- readme-gen:structure:start/end -->）は維持してください`, additionalContext)
	} else {
		prompt = fmt.Sprintf(`Update README.md based on the following information.

%s

Sections to update:
1. Description at the top - Write 1-2 sentences about the project's purpose
2. Structure section - Add brief description comments next to each directory
   Example format:
   ├── cmd/           # CLI entry point
   ├── internal/      # Internal packages
3. Usage section - Document basic usage
4. Development section - Document development commands (from mise/make/npm scripts etc.)

Rules:
- Read the code and generate appropriate content
- Keep it concise and practical
- Preserve existing markers (<!-- readme-gen:structure:start/end -->)`, additionalContext)
	}

	// Run claude command with timeout and spinner
	ctx, cancel := context.WithTimeout(context.Background(), aiTimeout)
	defer cancel()

	var runErr error
	action := func() {
		claudeCmd := exec.CommandContext(ctx, "claude", "-p", prompt, "--allowedTools", "Read,Edit,Glob,Grep")
		runErr = claudeCmd.Run()
	}

	spinnerTitle := msg.GeneratingWithAI
	if err := spinner.New().Title(spinnerTitle).Action(action).Run(); err != nil {
		return fmt.Errorf("%s: %w", msg.AIGenerationFailed, err)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("%s: timeout after %v", msg.AIGenerationFailed, aiTimeout)
	}

	if runErr != nil {
		return fmt.Errorf("%s: %w", msg.AIGenerationFailed, runErr)
	}

	fmt.Println(ui.Success(msg.AddedDescriptions))
	return nil
}

// collectProjectContext gathers additional context from project files
func collectProjectContext() string {
	var contexts []string

	// Check for mise.toml
	if content, err := os.ReadFile("mise.toml"); err == nil {
		contexts = append(contexts, fmt.Sprintf("mise.toml found:\n```toml\n%s\n```", truncateContent(string(content), 500)))
	}

	// Check for Makefile
	if content, err := os.ReadFile("Makefile"); err == nil {
		contexts = append(contexts, fmt.Sprintf("Makefile found:\n```make\n%s\n```", truncateContent(string(content), 500)))
	}

	// Check for package.json scripts
	if content, err := os.ReadFile("package.json"); err == nil {
		contexts = append(contexts, fmt.Sprintf("package.json found:\n```json\n%s\n```", truncateContent(string(content), 500)))
	}

	if len(contexts) == 0 {
		return "No additional build configuration files detected."
	}

	return "Detected project files:\n\n" + strings.Join(contexts, "\n\n")
}

// truncateContent truncates content to maxLen characters
func truncateContent(content string, maxLen int) string {
	if len(content) <= maxLen {
		return content
	}
	return content[:maxLen] + "\n... (truncated)"
}
