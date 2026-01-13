package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".readme-gen.yaml"

// Config represents the configuration for readme-gen
type Config struct {
	Structure StructureConfig `yaml:"structure"`
}

// StructureConfig configures directory structure scanning
type StructureConfig struct {
	// UseGitignore enables .gitignore pattern matching (default: true)
	UseGitignore bool `yaml:"use_gitignore"`
	// Patterns are additional include/exclude patterns (gitignore syntax)
	// Patterns starting with ! are include patterns (override excludes)
	Patterns []string `yaml:"patterns"`
	// MaxDepth limits directory traversal depth (0 = unlimited)
	MaxDepth int `yaml:"max_depth"`
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Structure: StructureConfig{
			UseGitignore: true,
			Patterns:     []string{},
			MaxDepth:     0,
		},
	}
}

// Load reads configuration from .readme-gen.yaml in the given directory
// If the file doesn't exist, returns default configuration
func Load(root string) (*Config, error) {
	cfg := Default()

	configPath := filepath.Join(root, ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ParsePatterns separates patterns into exclude and include lists
// Patterns starting with ! are include patterns
func (c *StructureConfig) ParsePatterns() (excludes, includes []string) {
	for _, p := range c.Patterns {
		if len(p) > 0 && p[0] == '!' {
			includes = append(includes, p[1:])
		} else {
			excludes = append(excludes, p)
		}
	}
	return
}
