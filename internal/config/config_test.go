package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if !cfg.Structure.UseGitignore {
		t.Error("expected UseGitignore to be true by default")
	}
	if cfg.Structure.MaxDepth != 0 {
		t.Errorf("expected MaxDepth to be 0, got %d", cfg.Structure.MaxDepth)
	}
	if len(cfg.Structure.Patterns) != 0 {
		t.Errorf("expected Patterns to be empty, got %v", cfg.Structure.Patterns)
	}
}

func TestLoad_NoFile(t *testing.T) {
	tmpDir := t.TempDir()

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should return defaults when no config file
	if !cfg.Structure.UseGitignore {
		t.Error("expected UseGitignore to be true")
	}
}

func TestLoad_WithFile(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `structure:
  use_gitignore: false
  max_depth: 3
  patterns:
    - "coverage/"
    - "!.git"
`
	err := os.WriteFile(filepath.Join(tmpDir, ConfigFileName), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Structure.UseGitignore {
		t.Error("expected UseGitignore to be false")
	}
	if cfg.Structure.MaxDepth != 3 {
		t.Errorf("expected MaxDepth to be 3, got %d", cfg.Structure.MaxDepth)
	}
	if len(cfg.Structure.Patterns) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(cfg.Structure.Patterns))
	}
}

func TestParsePatterns(t *testing.T) {
	cfg := StructureConfig{
		Patterns: []string{
			"coverage/",
			"tmp/",
			"!.git",
			"!node_modules",
		},
	}

	excludes, includes := cfg.ParsePatterns()

	if len(excludes) != 2 {
		t.Errorf("expected 2 excludes, got %d", len(excludes))
	}
	if len(includes) != 2 {
		t.Errorf("expected 2 includes, got %d", len(includes))
	}

	// Check excludes
	if excludes[0] != "coverage/" || excludes[1] != "tmp/" {
		t.Errorf("unexpected excludes: %v", excludes)
	}

	// Check includes (! should be stripped)
	if includes[0] != ".git" || includes[1] != "node_modules" {
		t.Errorf("unexpected includes: %v", includes)
	}
}
