package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hulk510/readme-gen/internal/config"
)

func TestMatcher_DefaultExcludesGit(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := config.Default()
	matcher := NewMatcher(tmpDir, cfg)

	// .git should be excluded by default
	if !matcher.IsExcluded(".git", true) {
		t.Error("expected .git to be excluded by default")
	}
}

func TestMatcher_GitignorePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .gitignore
	gitignoreContent := `node_modules/
coverage/
*.log
`
	err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("failed to write .gitignore: %v", err)
	}

	cfg := config.Default()
	matcher := NewMatcher(tmpDir, cfg)

	// Should exclude patterns from .gitignore
	if !matcher.IsExcluded("node_modules", true) {
		t.Error("expected node_modules to be excluded")
	}
	if !matcher.IsExcluded("coverage", true) {
		t.Error("expected coverage to be excluded")
	}

	// Should not exclude other directories
	if matcher.IsExcluded("src", true) {
		t.Error("expected src to NOT be excluded")
	}
}

func TestMatcher_IncludeOverride(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Structure: config.StructureConfig{
			UseGitignore: true,
			Patterns: []string{
				"!.git", // Override default .git exclusion
			},
		},
	}
	matcher := NewMatcher(tmpDir, cfg)

	// .git should NOT be excluded because of include pattern
	if matcher.IsExcluded(".git", true) {
		t.Error("expected .git to NOT be excluded due to include pattern")
	}
}

func TestMatcher_ExtraExcludes(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Structure: config.StructureConfig{
			UseGitignore: false,
			Patterns: []string{
				"tmp/",
				"cache/",
			},
		},
	}
	matcher := NewMatcher(tmpDir, cfg)

	if !matcher.IsExcluded("tmp", true) {
		t.Error("expected tmp to be excluded")
	}
	if !matcher.IsExcluded("cache", true) {
		t.Error("expected cache to be excluded")
	}
	if matcher.IsExcluded("src", true) {
		t.Error("expected src to NOT be excluded")
	}
}

func TestMatcher_MaxDepth(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &config.Config{
		Structure: config.StructureConfig{
			UseGitignore: true,
			MaxDepth:     3,
		},
	}
	matcher := NewMatcher(tmpDir, cfg)

	if matcher.MaxDepth() != 3 {
		t.Errorf("expected MaxDepth to be 3, got %d", matcher.MaxDepth())
	}
}

func TestScanAuto_WithGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"src/api",
		"src/models",
		"node_modules/pkg",
		"coverage/reports",
	}
	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, d), 0755)
		if err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	// Create .gitignore
	gitignoreContent := `node_modules/
coverage/
`
	err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("failed to write .gitignore: %v", err)
	}

	// Scan
	result, err := ScanAuto(tmpDir)
	if err != nil {
		t.Fatalf("ScanAuto failed: %v", err)
	}

	// Should contain src
	if !contains(result, "src/") {
		t.Error("expected result to contain 'src/'")
	}

	// Should NOT contain node_modules or coverage
	if contains(result, "node_modules") {
		t.Error("expected result to NOT contain 'node_modules'")
	}
	if contains(result, "coverage") {
		t.Error("expected result to NOT contain 'coverage'")
	}
}

func TestScanAuto_WithConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"src/api",
		".git/objects",
		"tmp/cache",
	}
	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, d), 0755)
		if err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	// Create config that includes .git and excludes tmp
	configContent := `structure:
  use_gitignore: true
  patterns:
    - "tmp/"
    - "!.git"
`
	err := os.WriteFile(filepath.Join(tmpDir, ".readme-gen.yaml"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// Scan
	result, err := ScanAuto(tmpDir)
	if err != nil {
		t.Fatalf("ScanAuto failed: %v", err)
	}

	// Should contain .git (override)
	if !contains(result, ".git/") {
		t.Error("expected result to contain '.git/' due to include override")
	}

	// Should NOT contain tmp
	if contains(result, "tmp") {
		t.Error("expected result to NOT contain 'tmp'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
