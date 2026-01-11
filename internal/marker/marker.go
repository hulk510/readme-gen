package marker

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// MarkerStart is the opening marker for structure section
	MarkerStart = "<!-- readme-gen:structure:start -->"
	// MarkerEnd is the closing marker for structure section
	MarkerEnd = "<!-- readme-gen:structure:end -->"
)

var markerRegex = regexp.MustCompile(`(?s)` + regexp.QuoteMeta(MarkerStart) + `.*?` + regexp.QuoteMeta(MarkerEnd))

// Extract extracts the structure content between markers
func Extract(content string) (string, bool) {
	startIdx := strings.Index(content, MarkerStart)
	endIdx := strings.Index(content, MarkerEnd)

	if startIdx == -1 || endIdx == -1 || startIdx >= endIdx {
		return "", false
	}

	// Extract content between markers
	between := content[startIdx+len(MarkerStart) : endIdx]

	// Find the code block content
	lines := strings.Split(between, "\n")
	var structureLines []string
	inCodeBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			structureLines = append(structureLines, line)
		}
	}

	result := strings.Join(structureLines, "\n")
	return strings.TrimSpace(result), true
}

// Update updates the structure section between markers
func Update(content string, structure string) (string, error) {
	if !strings.Contains(content, MarkerStart) || !strings.Contains(content, MarkerEnd) {
		return "", fmt.Errorf("markers not found in content")
	}

	newSection := fmt.Sprintf("%s\n```\n%s\n```\n%s", MarkerStart, structure, MarkerEnd)

	result := markerRegex.ReplaceAllString(content, newSection)
	return result, nil
}

// Wrap wraps structure content with markers
func Wrap(structure string) string {
	return fmt.Sprintf("%s\n```\n%s\n```\n%s", MarkerStart, structure, MarkerEnd)
}

// StripComments removes inline comments (# ...) from structure lines for comparison
func StripComments(structure string) string {
	lines := strings.Split(structure, "\n")
	var result []string

	for _, line := range lines {
		// Find comment marker that's not inside the tree structure
		// Comments typically appear after directory/file names with spaces
		if idx := strings.Index(line, "  #"); idx != -1 {
			line = strings.TrimRight(line[:idx], " ")
		} else if idx := strings.Index(line, "\t#"); idx != -1 {
			line = strings.TrimRight(line[:idx], " \t")
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
