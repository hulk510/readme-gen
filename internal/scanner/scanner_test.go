package scanner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultExcludes(t *testing.T) {
	excludes := DefaultExcludes()

	expected := []string{".git", "node_modules", "vendor", "dist", "build"}
	for _, e := range expected {
		found := false
		for _, ex := range excludes {
			if ex == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %s to be in default excludes", e)
		}
	}
}

func TestScan(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	dirs := []string{
		"src/api",
		"src/models",
		"internal/utils",
		"node_modules/pkg", // should be excluded
	}

	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, d), 0755)
		if err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	result, err := Scan(tmpDir, DefaultExcludes())
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should contain src and internal
	if !strings.Contains(result, "src/") {
		t.Error("expected result to contain 'src/'")
	}
	if !strings.Contains(result, "internal/") {
		t.Error("expected result to contain 'internal/'")
	}

	// Should NOT contain node_modules
	if strings.Contains(result, "node_modules") {
		t.Error("expected result to NOT contain 'node_modules'")
	}
}

func TestDetectProjectInfo_GoMod(t *testing.T) {
	tmpDir := t.TempDir()

	goModContent := `module github.com/example/myproject

go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	if err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	info := DetectProjectInfo(tmpDir)

	if info.Language != "go" {
		t.Errorf("expected language 'go', got '%s'", info.Language)
	}
	if info.ModulePath != "github.com/example/myproject" {
		t.Errorf("expected module path 'github.com/example/myproject', got '%s'", info.ModulePath)
	}
	if info.Name != "myproject" {
		t.Errorf("expected name 'myproject', got '%s'", info.Name)
	}
}

func TestDetectProjectInfo_PackageJson(t *testing.T) {
	tmpDir := t.TempDir()

	pkgContent := `{
  "name": "my-ts-project",
  "description": "A TypeScript project"
}`
	err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgContent), 0644)
	if err != nil {
		t.Fatalf("failed to write package.json: %v", err)
	}

	info := DetectProjectInfo(tmpDir)

	if info.Language != "typescript" {
		t.Errorf("expected language 'typescript', got '%s'", info.Language)
	}
	if info.Name != "my-ts-project" {
		t.Errorf("expected name 'my-ts-project', got '%s'", info.Name)
	}
	if info.Description != "A TypeScript project" {
		t.Errorf("expected description 'A TypeScript project', got '%s'", info.Description)
	}
}

func TestScan_TreeFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested structure
	dirs := []string{
		"cmd/app",
		"internal/handler",
		"internal/service",
		"pkg/utils",
	}

	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, d), 0755)
		if err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	result, err := Scan(tmpDir, DefaultExcludes())
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Check tree format characters
	if !strings.Contains(result, "├── ") && !strings.Contains(result, "└── ") {
		t.Error("expected tree format with ├── or └──")
	}
}
