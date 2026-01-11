package marker

import (
	"strings"
	"testing"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		want     string
		wantOk   bool
	}{
		{
			name: "valid markers with content",
			content: `# README

## Structure

<!-- readme-gen:structure:start -->
` + "```" + `
src/
├── api/
└── models/
` + "```" + `
<!-- readme-gen:structure:end -->

## Usage
`,
			want: `src/
├── api/
└── models/`,
			wantOk: true,
		},
		{
			name:    "no markers",
			content: "# README\n\nSome content",
			want:    "",
			wantOk:  false,
		},
		{
			name: "only start marker",
			content: `# README
<!-- readme-gen:structure:start -->
content
`,
			want:   "",
			wantOk: false,
		},
		{
			name: "only end marker",
			content: `# README
content
<!-- readme-gen:structure:end -->
`,
			want:   "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Extract(tt.content)
			if ok != tt.wantOk {
				t.Errorf("Extract() ok = %v, want %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("Extract() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	content := `# README

## Structure

<!-- readme-gen:structure:start -->
` + "```" + `
old/
└── structure/
` + "```" + `
<!-- readme-gen:structure:end -->

## Usage
`

	newStructure := `new/
├── api/
└── models/`

	result, err := Update(content, newStructure)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Should contain new structure
	if !strings.Contains(result, "new/") {
		t.Error("expected result to contain new structure")
	}

	// Should NOT contain old structure
	if strings.Contains(result, "old/") {
		t.Error("expected result to NOT contain old structure")
	}

	// Should preserve content outside markers
	if !strings.Contains(result, "# README") {
		t.Error("expected result to preserve header")
	}
	if !strings.Contains(result, "## Usage") {
		t.Error("expected result to preserve usage section")
	}
}

func TestUpdate_NoMarkers(t *testing.T) {
	content := "# README\n\nNo markers here"

	_, err := Update(content, "some structure")
	if err == nil {
		t.Error("expected error when no markers present")
	}
}

func TestWrap(t *testing.T) {
	structure := `src/
├── api/
└── models/`

	result := Wrap(structure)

	if !strings.Contains(result, MarkerStart) {
		t.Error("expected result to contain start marker")
	}
	if !strings.Contains(result, MarkerEnd) {
		t.Error("expected result to contain end marker")
	}
	if !strings.Contains(result, structure) {
		t.Error("expected result to contain structure")
	}
}
