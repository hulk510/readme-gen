package template

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"github.com/hulk510/readme-gen/internal/i18n"
	"github.com/hulk510/readme-gen/internal/marker"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// Data contains the data for rendering templates
type Data struct {
	ProjectName string
	Description string
	Structure   string
	Language    string
	ModulePath  string
	Lang        i18n.Language
}

// Render renders a template with the given data
func Render(templateName string, data Data) (string, error) {
	// Wrap structure with markers
	data.Structure = marker.Wrap(data.Structure)

	// Determine template file based on language
	tmplFile := fmt.Sprintf("templates/%s.md.tmpl", templateName)
	if data.Lang == i18n.Japanese {
		tmplFile = fmt.Sprintf("templates/%s_ja.md.tmpl", templateName)
	}

	// Load template
	tmplContent, err := templatesFS.ReadFile(tmplFile)
	if err != nil {
		// Fallback to English template if Japanese not found
		tmplContent, err = templatesFS.ReadFile(fmt.Sprintf("templates/%s.md.tmpl", templateName))
		if err != nil {
			return "", fmt.Errorf("template '%s' not found", templateName)
		}
	}

	// Parse and execute template
	tmpl, err := template.New(templateName).Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// GetClaudeSkills returns the Claude Code skills content
func GetClaudeSkills(lang i18n.Language) string {
	if lang == i18n.Japanese {
		return `# README更新スキル

プロジェクトの構造や機能が大きく変わったときにREADME.mdを更新する。

## トリガー

- ディレクトリ構造が変わったとき
- 新機能を追加したとき
- ユーザーが「README更新して」と言ったとき

## 手順

1. ` + "`readme-gen check`" + ` で構造の変更を検出
2. 同期されていなければ ` + "`readme-gen structure --update`" + ` を実行
3. 他のセクション（Usage、Descriptionなど）はコードの変更に基づいて更新
4. 更新は簡潔かつ正確に

## ルール

- マーカー外のコンテンツは明示的に必要な場合以外は変更しない
- READMEは簡潔に - 過度なドキュメント化は避ける
- 大きな変更について不明な場合はユーザーに確認する
`
	}

	return `# README Update Skill

Update README.md when project structure or functionality changes significantly.

## Trigger

- Directory structure changed
- New features added
- User requests "update README"

## Steps

1. Run ` + "`readme-gen check`" + ` to detect structure changes
2. If out of sync, run ` + "`readme-gen structure --update`" + `
3. Review and update other sections (Usage, Description) based on code changes
4. Keep updates concise and accurate

## Rules

- Don't modify content outside markers unless explicitly needed
- Keep README concise - avoid over-documentation
- Ask user if unsure about significant changes
`
}
