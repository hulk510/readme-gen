package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestDir creates a temporary directory with test files and returns cleanup function
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "readme-gen-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Save current dir
	origDir, err := os.Getwd()
	if err != nil {
		os.RemoveAll(dir)
		t.Fatalf("failed to get current dir: %v", err)
	}

	// Change to temp dir
	if err := os.Chdir(dir); err != nil {
		os.RemoveAll(dir)
		t.Fatalf("failed to change to temp dir: %v", err)
	}

	cleanup := func() {
		os.Chdir(origDir)
		os.RemoveAll(dir)
	}

	return dir, cleanup
}

// createTestFile creates a file with given content
func createTestFile(t *testing.T, path, content string) {
	t.Helper()
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
}

// readTestFile reads file content
func readTestFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(content)
}

func TestRunStructure_ShowOnly(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// Create test directory structure
	createTestFile(t, "src/main.go", "package main")
	createTestFile(t, "README.md", "# Test")

	// Reset flag
	updateFlag = false

	err := runStructure(nil, nil)
	if err != nil {
		t.Errorf("runStructure() error = %v", err)
	}
}

func TestRunStructure_Update(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// Create test files
	createTestFile(t, "src/main.go", "package main")

	readmeContent := `# Test Project

<!-- readme-gen:structure:start -->
` + "```" + `
old structure
` + "```" + `
<!-- readme-gen:structure:end -->
`
	createTestFile(t, "README.md", readmeContent)

	// Set update flag
	updateFlag = true

	err := runStructure(nil, nil)
	if err != nil {
		t.Errorf("runStructure() error = %v", err)
	}

	// Check README was updated
	content := readTestFile(t, "README.md")
	if !strings.Contains(content, "src/") {
		t.Errorf("README should contain updated structure with src/, got: %s", content)
	}
}

func TestRunStructure_NoReadme(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	updateFlag = true

	err := runStructure(nil, nil)
	if err == nil {
		t.Error("runStructure() should return error when README.md not found")
	}
}

func TestRunCheck_UpToDate(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// Create test directory structure (scanner only shows directories)
	createTestFile(t, "src/main.go", "package main")

	// Create README with current structure (directories only)
	readmeContent := `# Test Project

<!-- readme-gen:structure:start -->
` + "```" + `
└── src/
` + "```" + `
<!-- readme-gen:structure:end -->
`
	createTestFile(t, "README.md", readmeContent)

	// Mock exitFunc to not actually exit
	origExitFunc := exitFunc
	exitFunc = func(code int) {}
	defer func() { exitFunc = origExitFunc }()

	err := runCheck(nil, nil)
	if err != nil {
		t.Errorf("runCheck() error = %v", err)
	}
}

func TestRunCheck_NoMarkers(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// README without markers
	createTestFile(t, "README.md", "# Test Project\n\nNo markers here.")

	err := runCheck(nil, nil)
	if err != nil {
		t.Errorf("runCheck() should not error when no markers found, got: %v", err)
	}
}

func TestRunCheck_OutOfSync(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// Mock exitFunc to prevent os.Exit
	exitCalled := false
	exitCode := 0
	origExitFunc := exitFunc
	exitFunc = func(code int) {
		exitCalled = true
		exitCode = code
	}
	defer func() { exitFunc = origExitFunc }()

	// Create test directories (scanner only shows directories)
	createTestFile(t, "src/main.go", "package main")
	createTestFile(t, "newdir/file.go", "package newdir") // New directory not in README

	// Create README with old structure (missing newdir/)
	readmeContent := `# Test Project

<!-- readme-gen:structure:start -->
` + "```" + `
└── src/
` + "```" + `
<!-- readme-gen:structure:end -->
`
	createTestFile(t, "README.md", readmeContent)

	err := runCheck(nil, nil)
	if err != ErrOutOfSync {
		t.Errorf("runCheck() should return ErrOutOfSync, got: %v", err)
	}
	if !exitCalled {
		t.Error("exitFunc should be called")
	}
	if exitCode != 1 {
		t.Errorf("exit code should be 1, got: %d", exitCode)
	}
}

func TestRunCheck_NoReadme(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	err := runCheck(nil, nil)
	if err == nil {
		t.Error("runCheck() should return error when README.md not found")
	}
}

func TestRunInit_NonInteractive(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// Create go.mod for project detection
	createTestFile(t, "go.mod", "module github.com/test/project\n\ngo 1.21")
	createTestFile(t, "main.go", "package main")

	// Set non-interactive mode
	nonInteractive = true
	templateFlag = "oss"
	noSkills = true
	noAI = true

	err := runInit(nil, nil)
	if err != nil {
		t.Errorf("runInit() error = %v", err)
	}

	// Check README was created
	if _, err := os.Stat("README.md"); os.IsNotExist(err) {
		t.Error("README.md should be created")
	}

	content := readTestFile(t, "README.md")
	if !strings.Contains(content, "project") {
		t.Errorf("README should contain project name, got: %s", content)
	}
	if !strings.Contains(content, "readme-gen:structure:start") {
		t.Errorf("README should contain structure markers, got: %s", content)
	}
}

func TestRunInit_WithSkills(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	createTestFile(t, "go.mod", "module github.com/test/project\n\ngo 1.21")

	nonInteractive = true
	templateFlag = "oss"
	withSkills = true
	noSkills = false
	noAI = true

	err := runInit(nil, nil)
	if err != nil {
		t.Errorf("runInit() error = %v", err)
	}

	// Check skills file was created
	skillsPath := filepath.Join(".claude", "skills", "readme.md")
	if _, err := os.Stat(skillsPath); os.IsNotExist(err) {
		t.Error("skills file should be created")
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		slice []string
		item  string
		want  bool
	}{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
		{[]string{"skills", "ai"}, "skills", true},
	}

	for _, tt := range tests {
		got := contains(tt.slice, tt.item)
		if got != tt.want {
			t.Errorf("contains(%v, %q) = %v, want %v", tt.slice, tt.item, got, tt.want)
		}
	}
}
