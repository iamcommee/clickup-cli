package config

import (
	"os"
	"path/filepath"
)

// Cleanup removes all ClickUp CLI config files and returns the list of removed paths
func Cleanup(homeDir string) ([]string, error) {
	var removed []string

	// Remove new XDG config file
	configDir := filepath.Join(homeDir, ".config", AppName)
	configPath := filepath.Join(configDir, ConfigFileName)

	if _, err := os.Stat(configPath); err == nil {
		if err := os.Remove(configPath); err != nil {
			return removed, err
		}
		removed = append(removed, configPath)

		// Try to remove the directory if empty
		_ = os.Remove(configDir)
	}

	// Remove legacy config file
	legacyPath := filepath.Join(homeDir, LegacyConfigFileName)
	if _, err := os.Stat(legacyPath); err == nil {
		if err := os.Remove(legacyPath); err != nil {
			return removed, err
		}
		removed = append(removed, legacyPath)
	}

	return removed, nil
}

// GetConfigFiles returns all config file paths that exist
func GetConfigFiles(homeDir string) []string {
	var files []string

	// Check new XDG config file
	configPath := filepath.Join(homeDir, ".config", AppName, ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		files = append(files, configPath)
	}

	// Check legacy config file
	legacyPath := filepath.Join(homeDir, LegacyConfigFileName)
	if _, err := os.Stat(legacyPath); err == nil {
		files = append(files, legacyPath)
	}

	return files
}
