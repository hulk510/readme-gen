package i18n

import (
	"os"
	"testing"
)

func TestSetLanguage(t *testing.T) {
	// Set to Japanese
	SetLanguage(Japanese)
	if Current() != Japanese {
		t.Errorf("expected Japanese, got %s", Current())
	}

	// Set to English
	SetLanguage(English)
	if Current() != English {
		t.Errorf("expected English, got %s", Current())
	}
}

func TestGet(t *testing.T) {
	// Test English messages
	SetLanguage(English)
	msg := Get()
	if msg.CreatedReadme != "Created README.md" {
		t.Errorf("expected 'Created README.md', got '%s'", msg.CreatedReadme)
	}

	// Test Japanese messages
	SetLanguage(Japanese)
	msg = Get()
	if msg.CreatedReadme != "README.mdを作成しました" {
		t.Errorf("expected 'README.mdを作成しました', got '%s'", msg.CreatedReadme)
	}
}

func TestDetectLanguage(t *testing.T) {
	// Save original env
	origLang := os.Getenv("LANG")
	origLanguage := os.Getenv("LANGUAGE")
	origLcAll := os.Getenv("LC_ALL")

	defer func() {
		os.Setenv("LANG", origLang)
		os.Setenv("LANGUAGE", origLanguage)
		os.Setenv("LC_ALL", origLcAll)
	}()

	tests := []struct {
		name     string
		lang     string
		language string
		lcAll    string
		want     Language
	}{
		{
			name: "Japanese LANG",
			lang: "ja_JP.UTF-8",
			want: Japanese,
		},
		{
			name:     "Japanese LANGUAGE",
			language: "ja",
			want:     Japanese,
		},
		{
			name:  "Japanese LC_ALL",
			lcAll: "ja_JP.UTF-8",
			want:  Japanese,
		},
		{
			name: "English LANG",
			lang: "en_US.UTF-8",
			want: English,
		},
		{
			name: "No env set",
			want: English,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv("LANG")
			os.Unsetenv("LANGUAGE")
			os.Unsetenv("LC_ALL")

			if tt.lang != "" {
				os.Setenv("LANG", tt.lang)
			}
			if tt.language != "" {
				os.Setenv("LANGUAGE", tt.language)
			}
			if tt.lcAll != "" {
				os.Setenv("LC_ALL", tt.lcAll)
			}

			got := DetectLanguage()
			if got != tt.want {
				t.Errorf("DetectLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}
