package scanner

import (
	"os"
	"path/filepath"

	"github.com/hulk510/readme-gen/internal/config"
	ignore "github.com/sabhiram/go-gitignore"
)

// Matcher handles pattern matching for directory exclusion
type Matcher struct {
	root         string
	gitignore    *ignore.GitIgnore
	extraExclude *ignore.GitIgnore
	includes     map[string]bool
	maxDepth     int
}

// NewMatcher creates a new Matcher from configuration
func NewMatcher(root string, cfg *config.Config) *Matcher {
	m := &Matcher{
		root:     root,
		includes: make(map[string]bool),
		maxDepth: cfg.Structure.MaxDepth,
	}

	// Load .gitignore if enabled
	if cfg.Structure.UseGitignore {
		gitignorePath := filepath.Join(root, ".gitignore")
		if gi, err := ignore.CompileIgnoreFile(gitignorePath); err == nil {
			m.gitignore = gi
		}
	}

	// Parse config patterns
	excludes, includes := cfg.Structure.ParsePatterns()

	// Build include set
	for _, inc := range includes {
		m.includes[inc] = true
	}

	// Build extra exclude matcher
	if len(excludes) > 0 {
		m.extraExclude = ignore.CompileIgnoreLines(excludes...)
	}

	return m
}

// IsExcluded checks if a path should be excluded from the tree
func (m *Matcher) IsExcluded(relPath string, isDir bool) bool {
	name := filepath.Base(relPath)

	// Check if explicitly included (overrides everything)
	if m.includes[name] || m.includes[relPath] {
		return false
	}

	// Always exclude .git by default (unless explicitly included)
	if name == ".git" {
		return true
	}

	// Check .gitignore patterns
	if m.gitignore != nil {
		// go-gitignore expects paths relative to .gitignore location
		checkPath := relPath
		if isDir {
			checkPath = relPath + "/"
		}
		if m.gitignore.MatchesPath(checkPath) {
			return true
		}
	}

	// Check extra exclude patterns from config
	if m.extraExclude != nil {
		checkPath := relPath
		if isDir {
			checkPath = relPath + "/"
		}
		if m.extraExclude.MatchesPath(checkPath) {
			return true
		}
	}

	return false
}

// MaxDepth returns the configured max depth (0 = unlimited)
func (m *Matcher) MaxDepth() int {
	return m.maxDepth
}

// DefaultMatcher returns a matcher with default settings (no config file)
func DefaultMatcher(root string) *Matcher {
	cfg := config.Default()
	return NewMatcher(root, cfg)
}

// LegacyMatcher creates a matcher that behaves like the old DefaultExcludes
// Used for backward compatibility
func LegacyMatcher(root string, excludes []string) *Matcher {
	m := &Matcher{
		root:     root,
		includes: make(map[string]bool),
	}

	if len(excludes) > 0 {
		m.extraExclude = ignore.CompileIgnoreLines(excludes...)
	}

	return m
}

// LoadMatcher loads configuration and creates a matcher
func LoadMatcher(root string) (*Matcher, error) {
	cfg, err := config.Load(root)
	if err != nil {
		return nil, err
	}
	return NewMatcher(root, cfg), nil
}

// gitignoreExists checks if .gitignore file exists
func gitignoreExists(root string) bool {
	_, err := os.Stat(filepath.Join(root, ".gitignore"))
	return err == nil
}
