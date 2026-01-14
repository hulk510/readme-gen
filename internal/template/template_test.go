package template

import (
	"strings"
	"testing"

	"github.com/hulk510/readme-gen/internal/i18n"
)

func TestRender_OSS(t *testing.T) {
	data := Data{
		ProjectName: "test-project",
		Description: "A test project",
		Structure:   "src/\n└── main.go",
		Language:    "go",
		ModulePath:  "github.com/example/test-project",
		Lang:        i18n.English,
	}

	result, err := Render("oss", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check project name
	if !strings.Contains(result, "# test-project") {
		t.Error("expected result to contain project name header")
	}

	// Check structure markers
	if !strings.Contains(result, "<!-- readme-gen:structure:start -->") {
		t.Error("expected result to contain structure start marker")
	}
	if !strings.Contains(result, "<!-- readme-gen:structure:end -->") {
		t.Error("expected result to contain structure end marker")
	}

	// Check Go-specific content
	if !strings.Contains(result, "go install") {
		t.Error("expected result to contain 'go install' for Go projects")
	}

	// Check OSS-specific sections
	if !strings.Contains(result, "Contributing") {
		t.Error("expected OSS template to contain Contributing section")
	}
	if !strings.Contains(result, "License") {
		t.Error("expected OSS template to contain License section")
	}
}

func TestRender_General(t *testing.T) {
	data := Data{
		ProjectName: "my-app",
		Description: "",
		Structure:   "src/",
		Language:    "typescript",
		Lang:        i18n.English,
	}

	result, err := Render("general", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// General template should have project name
	if !strings.Contains(result, "# my-app") {
		t.Error("expected result to contain project name")
	}

	// Should have TS-specific content
	if !strings.Contains(result, "bun") {
		t.Error("expected TypeScript template to mention bun")
	}

	// General template should have Getting Started
	if !strings.Contains(result, "Getting Started") {
		t.Error("expected general template to contain Getting Started section")
	}

	// General template should have Development
	if !strings.Contains(result, "Development") {
		t.Error("expected general template to contain Development section")
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	data := Data{
		ProjectName: "test",
		Lang:        i18n.English,
	}

	_, err := Render("nonexistent", data)
	if err == nil {
		t.Error("expected error for invalid template")
	}
}

func TestRender_OSS_Japanese(t *testing.T) {
	data := Data{
		ProjectName: "test-project",
		Description: "テストプロジェクト",
		Structure:   "src/\n└── main.go",
		Language:    "go",
		ModulePath:  "github.com/example/test-project",
		Lang:        i18n.Japanese,
	}

	result, err := Render("oss", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check Japanese content
	if !strings.Contains(result, "コントリビューション") {
		t.Error("expected Japanese OSS template to contain 'コントリビューション'")
	}
	if !strings.Contains(result, "ライセンス") {
		t.Error("expected Japanese OSS template to contain 'ライセンス'")
	}
	if !strings.Contains(result, "インストール") {
		t.Error("expected Japanese OSS template to contain 'インストール'")
	}
}

func TestRender_General_Japanese(t *testing.T) {
	data := Data{
		ProjectName: "my-app",
		Description: "内部ツール",
		Structure:   "cmd/\ninternal/",
		Language:    "go",
		ModulePath:  "github.com/company/my-app",
		Lang:        i18n.Japanese,
	}

	result, err := Render("general", data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check Japanese content
	if !strings.Contains(result, "はじめに") {
		t.Error("expected Japanese general template to contain 'はじめに'")
	}
	if !strings.Contains(result, "開発") {
		t.Error("expected Japanese general template to contain '開発'")
	}
}

func TestGetClaudeSkills(t *testing.T) {
	result := GetClaudeSkills()

	if !strings.Contains(result, "README Update Skill") {
		t.Error("expected skills to contain 'README Update Skill'")
	}
	if !strings.Contains(result, "readme-gen check") {
		t.Error("expected skills to contain 'readme-gen check' command")
	}
	if !strings.Contains(result, "user_invocable: true") {
		t.Error("expected skills to contain 'user_invocable: true'")
	}
}
