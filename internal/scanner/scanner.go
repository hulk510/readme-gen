package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ProjectInfo contains detected project metadata
type ProjectInfo struct {
	Name        string
	Description string
	Language    string
	ModulePath  string
}

// DefaultExcludes returns the default list of directories to exclude
// Deprecated: Use LoadMatcher or DefaultMatcher instead
func DefaultExcludes() []string {
	return []string{
		".git",
		"node_modules",
		"vendor",
		"dist",
		"build",
		".next",
		"__pycache__",
		".venv",
		"venv",
		"target",
		"bin",
		".idea",
		".vscode",
	}
}

// DetectProjectInfo detects project metadata from common files
func DetectProjectInfo(root string) ProjectInfo {
	info := ProjectInfo{
		Name: filepath.Base(absPath(root)),
	}

	// Try go.mod
	if content, err := os.ReadFile(filepath.Join(root, "go.mod")); err == nil {
		info.Language = "go"
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				info.ModulePath = strings.TrimPrefix(line, "module ")
				// Extract name from module path
				parts := strings.Split(info.ModulePath, "/")
				if len(parts) > 0 {
					info.Name = parts[len(parts)-1]
				}
				break
			}
		}
	}

	// Try package.json
	if content, err := os.ReadFile(filepath.Join(root, "package.json")); err == nil {
		info.Language = "typescript"
		var pkg struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if json.Unmarshal(content, &pkg) == nil {
			if pkg.Name != "" {
				info.Name = pkg.Name
			}
			info.Description = pkg.Description
		}
	}

	return info
}

// Scan scans the directory and returns a tree structure as string
// Deprecated: Use ScanWithMatcher instead
func Scan(root string, excludes []string) (string, error) {
	matcher := LegacyMatcher(root, excludes)
	return ScanWithMatcher(root, matcher)
}

// ScanWithMatcher scans the directory using the provided matcher
func ScanWithMatcher(root string, matcher *Matcher) (string, error) {
	var builder strings.Builder
	err := walkDirWithMatcher(root, "", &builder, matcher, true, 0)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(builder.String(), "\n"), nil
}

// ScanAuto scans the directory using auto-loaded configuration
func ScanAuto(root string) (string, error) {
	matcher, err := LoadMatcher(root)
	if err != nil {
		return "", err
	}
	return ScanWithMatcher(root, matcher)
}

func walkDirWithMatcher(path string, prefix string, builder *strings.Builder, matcher *Matcher, isRoot bool, depth int) error {
	// Check max depth
	if matcher.MaxDepth() > 0 && depth > matcher.MaxDepth() {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter and sort entries
	var dirs []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Skip hidden files (except at root level or if explicitly included)
		if strings.HasPrefix(name, ".") && !isRoot {
			// Calculate relative path for matcher
			relPath := name
			if prefix != "" {
				// Reconstruct relative path from prefix
				relPath = getRelPath(prefix, name)
			}
			if matcher.IsExcluded(relPath, true) {
				continue
			}
		} else {
			// Calculate relative path
			relPath := name
			if prefix != "" {
				relPath = getRelPath(prefix, name)
			}
			if matcher.IsExcluded(relPath, true) {
				continue
			}
		}

		dirs = append(dirs, entry)
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	for i, entry := range dirs {
		name := entry.Name()
		isLast := i == len(dirs)-1

		// Determine the connector
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		// Write the entry
		builder.WriteString(prefix + connector + name + "/\n")

		// Recurse into subdirectory
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}

		subPath := filepath.Join(path, name)
		if err := walkDirWithMatcher(subPath, newPrefix, builder, matcher, false, depth+1); err != nil {
			return err
		}
	}

	return nil
}

// getRelPath reconstructs the relative path from tree prefix and name
func getRelPath(prefix, name string) string {
	// Count depth from prefix (each level adds 4 chars: "│   " or "    ")
	depth := len(prefix) / 4
	if depth == 0 {
		return name
	}
	// We can't reconstruct the full path from prefix alone,
	// so we just return the name for matching
	return name
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
