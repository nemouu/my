package config

// This file will contain configuration management for the μ shell.
// It handles loading configuration from YAML files, providing defaults,
// and validating configuration values.
//
// TODO: Implement config system as described in TODO.md Milestone 2 "Configuration System"
// - Define Config struct with nested structs (Ollama, Models, Context, Session, Shell)
// - Implement LoadConfig() to read from ~/.config/my/config.yaml
// - Provide sensible defaults if config file is missing
// - Parse YAML using gopkg.in/yaml.v3
// - Validate required fields

// Config represents the complete configuration for the μ shell
type Config struct {
	// TODO: Add nested configuration structs
	// Ollama  OllamaConfig
	// Models  ModelsConfig
	// Context ContextConfig
	// Session SessionConfig
	// Shell   ShellConfig
}

// LoadConfig loads configuration from ~/.config/my/config.yaml or uses defaults
func LoadConfig() (*Config, error) {
	// TODO: Implement config loading with fallback to defaults
	return nil, nil
}

// SaveDefaultConfig creates a default config file at ~/.config/my/config.yaml
func SaveDefaultConfig() error {
	// TODO: Implement for the --init flag
	return nil
}
