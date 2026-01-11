package i18n

import (
	"os"
	"strings"
)

// Language represents supported languages
type Language string

const (
	English  Language = "en"
	Japanese Language = "ja"
)

var currentLang Language = English

// Messages contains all translatable strings
type Messages struct {
	// CLI
	AppDescription     string
	SelectTemplate     string
	ProjectName        string
	AddClaudeInteg     string
	ClaudeIntegDesc    string
	OverwriteConfirm   string
	Cancelled          string
	UpdatingStructure  string
	ChangesDetected    string
	NoMarkersFound     string
	AddMarkersHint     string
	StructureUpToDate  string
	StructureOutOfSync string
	RunUpdateHint      string
	ReadmeNotFound     string
	RunInitHint        string

	// Template options
	TemplateOSS      string
	TemplatePersonal string
	TemplateTeam     string

	// Success messages
	CreatedReadme string
	CreatedSkills string
	UpdatedReadme string

	// Structure markers info
	MarkersInfo  string
	RunLaterHint string

	// Claude Code integration
	ClaudeCodeIntegration string
	OptionAddSkills       string
	OptionGenerateWithAI  string
	GeneratingWithAI      string
	AddedDescriptions     string
	ClaudeCodeNotFound    string
	AIGenerationFailed    string

	// Steps
	StepLanguage    string
	StepTemplate    string
	StepProjectInfo string
	StepIntegration string

	// Language selection
	SelectLanguage string
	LangEnglish    string
	LangJapanese   string
}

var messages = map[Language]Messages{
	English: {
		AppDescription:     "A CLI tool for managing README.md with structure auto-sync",
		SelectTemplate:     "Select template",
		ProjectName:        "Project name",
		AddClaudeInteg:     "Add Claude Code integration?",
		ClaudeIntegDesc:    "Adds skills to .claude/skills/",
		OverwriteConfirm:   "Overwrite existing README.md?",
		Cancelled:          "Cancelled",
		UpdatingStructure:  "Updating structure...",
		ChangesDetected:    "Changes detected",
		NoMarkersFound:     "No structure markers found in README.md",
		AddMarkersHint:     "Add markers with `readme-gen init` or manually",
		StructureUpToDate:  "README.md is up to date",
		StructureOutOfSync: "Structure out of sync!",
		RunUpdateHint:      "Run `readme-gen structure --update` to fix",
		ReadmeNotFound:     "README.md not found",
		RunInitHint:        "Run `readme-gen init` first",

		TemplateOSS:      "oss - MIT license, contributing guide",
		TemplatePersonal: "personal - Simple, minimal",
		TemplateTeam:     "team - Internal docs style",

		CreatedReadme: "Created README.md",
		CreatedSkills: "Created .claude/skills/readme.md",
		UpdatedReadme: "README.md updated!",

		MarkersInfo:  "Structure will be placed between markers",
		RunLaterHint: "Run `readme-gen structure` to update later",

		ClaudeCodeIntegration: "Claude Code integration",
		OptionAddSkills:       "Add skills to .claude/skills/",
		OptionGenerateWithAI:  "Generate descriptions with AI",
		GeneratingWithAI:      "Generating descriptions with Claude Code...",
		AddedDescriptions:     "Added directory descriptions",
		ClaudeCodeNotFound:    "Claude Code not found. Skipping AI generation.",
		AIGenerationFailed:    "AI generation failed",

		StepLanguage:    "Language",
		StepTemplate:    "Template",
		StepProjectInfo: "Project Info",
		StepIntegration: "Integrations",

		SelectLanguage: "Select language",
		LangEnglish:    "English",
		LangJapanese:   "日本語 (Japanese)",
	},
	Japanese: {
		AppDescription:     "README.mdを構造自動同期で管理するCLIツール",
		SelectTemplate:     "テンプレートを選択",
		ProjectName:        "プロジェクト名",
		AddClaudeInteg:     "Claude Code連携を追加しますか？",
		ClaudeIntegDesc:    ".claude/skills/にスキルを追加します",
		OverwriteConfirm:   "既存のREADME.mdを上書きしますか？",
		Cancelled:          "キャンセルしました",
		UpdatingStructure:  "構造を更新中...",
		ChangesDetected:    "変更を検出",
		NoMarkersFound:     "README.mdに構造マーカーが見つかりません",
		AddMarkersHint:     "`readme-gen init`またはマーカーを手動で追加してください",
		StructureUpToDate:  "README.mdは最新です",
		StructureOutOfSync: "構造が同期されていません！",
		RunUpdateHint:      "`readme-gen structure --update`で修正してください",
		ReadmeNotFound:     "README.mdが見つかりません",
		RunInitHint:        "先に`readme-gen init`を実行してください",

		TemplateOSS:      "oss - MITライセンス、コントリビューションガイド付き",
		TemplatePersonal: "personal - シンプル、最小構成",
		TemplateTeam:     "team - 内部ドキュメント向け",

		CreatedReadme: "README.mdを作成しました",
		CreatedSkills: ".claude/skills/readme.mdを作成しました",
		UpdatedReadme: "README.mdを更新しました！",

		MarkersInfo:  "構造はマーカー間に配置されます",
		RunLaterHint: "`readme-gen structure`で後から更新できます",

		ClaudeCodeIntegration: "Claude Code連携",
		OptionAddSkills:       ".claude/skills/にスキルを追加",
		OptionGenerateWithAI:  "AIで説明を自動生成",
		GeneratingWithAI:      "Claude Codeで説明を生成中...",
		AddedDescriptions:     "ディレクトリの説明を追加しました",
		ClaudeCodeNotFound:    "Claude Codeが見つかりません。AI生成をスキップします。",
		AIGenerationFailed:    "AI生成に失敗しました",

		StepLanguage:    "言語",
		StepTemplate:    "テンプレート",
		StepProjectInfo: "プロジェクト情報",
		StepIntegration: "連携設定",

		SelectLanguage: "言語を選択",
		LangEnglish:    "English (英語)",
		LangJapanese:   "日本語",
	},
}

// SetLanguage sets the current language
func SetLanguage(lang Language) {
	currentLang = lang
}

// DetectLanguage detects language from environment
func DetectLanguage() Language {
	// Check LANG environment variable
	lang := os.Getenv("LANG")
	if strings.HasPrefix(lang, "ja") {
		return Japanese
	}

	// Check LANGUAGE
	lang = os.Getenv("LANGUAGE")
	if strings.HasPrefix(lang, "ja") {
		return Japanese
	}

	// Check LC_ALL
	lang = os.Getenv("LC_ALL")
	if strings.HasPrefix(lang, "ja") {
		return Japanese
	}

	return English
}

// Get returns the current messages
func Get() Messages {
	return messages[currentLang]
}

// Current returns the current language
func Current() Language {
	return currentLang
}
