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
	// LegacyConfigFileName is the old config file name for backwards compatibility
	LegacyConfigFileName = ".clickup.json"
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

// FindConfigFile searches for config file in order of preference:
// 1. ~/.config/clickup/config.json (new XDG location)
// 2. ~/.clickup.json (legacy location)
// Returns empty string if no config file found
func FindConfigFile(homeDir string) string {
	// Check new XDG location first
	newPath := filepath.Join(homeDir, ".config", AppName, ConfigFileName)
	if _, err := os.Stat(newPath); err == nil {
		return newPath
	}

	// Check legacy location
	legacyPath := filepath.Join(homeDir, LegacyConfigFileName)
	if _, err := os.Stat(legacyPath); err == nil {
		return legacyPath
	}

	return ""
}

// IsLegacyPath checks if a path is the legacy config location
func IsLegacyPath(path string) bool {
	return filepath.Base(path) == LegacyConfigFileName
}
