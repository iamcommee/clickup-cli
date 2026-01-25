package config

import (
	"os"
	"path/filepath"
)

const (
	// AppName is the application name used for config directory
	AppName = "clickup"
	// ConfigFileName is the config file name
	ConfigFileName = "config.json"
)

// DefaultConfigDir returns the default config directory path
// Following XDG Base Directory Specification:
// - Uses $XDG_CONFIG_HOME/clickup if XDG_CONFIG_HOME is set
// - Otherwise uses ~/.config/clickup
func DefaultConfigDir() string {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, AppName)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".config", AppName)
}

// DefaultConfigPath returns the full path to the config file
func DefaultConfigPath() string {
	return filepath.Join(DefaultConfigDir(), ConfigFileName)
}

// EnsureConfigDir creates the config directory if it doesn't exist
// with secure permissions (0700 - owner only)
func EnsureConfigDir(dir string) error {
	return os.MkdirAll(dir, 0700)
}

// FindConfigFile searches for config file
// Returns empty string if no config file found
func FindConfigFile(homeDir string) string {
	configPath := filepath.Join(homeDir, ".config", AppName, ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}
	return ""
}
