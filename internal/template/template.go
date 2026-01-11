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
		return `---
description: README.mdの構造を更新する。/readme または「READMEを更新して」で実行。
user_invocable: true
---

# README更新スキル

プロジェクトの構造変更時にREADME.mdを自動更新します。

## 使い方

ユーザーが以下のいずれかを実行した場合:
- ` + "`/readme`" + ` コマンド
- 「READMEを更新して」などの指示

## 手順

1. ` + "`readme-gen check`" + ` を実行して構造の変更を確認
2. 同期されていなければ ` + "`readme-gen structure --update`" + ` を実行
3. 構造セクション（マーカー間）のみ更新
4. 他のセクションは明示的に依頼された場合のみ変更

## ルール

- マーカー (` + "`<!-- readme-gen:structure:start -->`" + ` / ` + "`<!-- readme-gen:structure:end -->`" + `) 間のみ自動更新
- コメント (` + "`# 説明`" + `) がある場合は保持
- 大きな変更は確認を取る
`
	}

	return `---
description: Update README.md structure. Run with /readme or "update the README".
user_invocable: true
---

# README Update Skill

Automatically update README.md when project structure changes.

## Usage

When user runs:
- ` + "`/readme`" + ` command
- "Update the README" or similar request

## Steps

1. Run ` + "`readme-gen check`" + ` to detect structure changes
2. If out of sync, run ` + "`readme-gen structure --update`" + `
3. Only update structure section (between markers)
4. Modify other sections only when explicitly requested

## Rules

- Only auto-update content between markers (` + "`<!-- readme-gen:structure:start -->`" + ` / ` + "`<!-- readme-gen:structure:end -->`" + `)
- Preserve comments (` + "`# description`" + `) if present
- Confirm with user for significant changes
`
}
