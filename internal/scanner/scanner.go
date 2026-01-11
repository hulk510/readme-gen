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
func Scan(root string, excludes []string) (string, error) {
	excludeMap := make(map[string]bool)
	for _, e := range excludes {
		excludeMap[e] = true
	}

	var builder strings.Builder
	err := walkDir(root, "", &builder, excludeMap, true)
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(builder.String(), "\n"), nil
}

func walkDir(path string, prefix string, builder *strings.Builder, excludes map[string]bool, isRoot bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Filter and sort entries
	var dirs []os.DirEntry
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files and excluded directories
		if strings.HasPrefix(name, ".") && !isRoot {
			continue
		}
		if excludes[name] {
			continue
		}
		if entry.IsDir() {
			dirs = append(dirs, entry)
		}
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
		if err := walkDir(subPath, newPrefix, builder, excludes, false); err != nil {
			return err
		}
	}

	return nil
}

func absPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
