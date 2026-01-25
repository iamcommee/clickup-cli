package config

import (
	"os"
	"path/filepath"
)

// Cleanup removes ClickUp CLI config files and returns the list of removed paths
func Cleanup(homeDir string) ([]string, error) {
	var removed []string

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

	return removed, nil
}

// GetConfigFiles returns config file path if it exists
func GetConfigFiles(homeDir string) []string {
	var files []string

	configPath := filepath.Join(homeDir, ".config", AppName, ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		files = append(files, configPath)
	}

	return files
}
