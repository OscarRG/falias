package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Theme string `yaml:"theme"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Theme: "default",
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "falias")
	return filepath.Join(configDir, "config.yaml"), nil
}

// EnsureConfigDir ensures the config directory exists
func EnsureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "falias")
	return os.MkdirAll(configDir, 0755)
}

// Load loads the configuration from disk
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config doesn't exist, create default
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if !IsValidTheme(cfg.Theme) {
		cfg.Theme = "default"
	}

	return &cfg, nil
}

// Save saves the configuration to disk
func (c *Config) Save() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SetTheme sets the theme and saves the config
func (c *Config) SetTheme(theme string) error {
	if !IsValidTheme(theme) {
		return fmt.Errorf("invalid theme: %s (available: %v)", theme, GetAvailableThemes())
	}

	c.Theme = theme
	return c.Save()
}

// IsValidTheme checks if a theme name is valid
func IsValidTheme(theme string) bool {
	themes := GetAvailableThemes()
	for _, t := range themes {
		if t == theme {
			return true
		}
	}
	return false
}

// GetAvailableThemes returns all available theme names
func GetAvailableThemes() []string {
	return []string{
		"default",
		"light",
		"dark",
		"high-contrast",
		"nord",
		"gruvbox",
	}
}
