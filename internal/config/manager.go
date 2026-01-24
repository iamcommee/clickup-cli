package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	// DefaultConfigFileName is the default config file name
	DefaultConfigFileName = ".clickup.json"
)

var (
	// ErrConfigNotFound is returned when no config file is found
	ErrConfigNotFound = errors.New("config file not found. Run 'clickup config init' to create one")
	// ErrMissingAPIToken is returned when API token is missing
	ErrMissingAPIToken = errors.New("API token is required")
	// ErrMissingWorkspaceID is returned when workspace ID is missing
	ErrMissingWorkspaceID = errors.New("workspace ID is required")
	// ErrMissingUserID is returned when user ID is missing
	ErrMissingUserID = errors.New("user ID is required")
)

// Manager handles config loading and saving
type Manager struct {
	configPath string
}

// NewManager creates a new config manager
func NewManager(configPath string) *Manager {
	return &Manager{configPath: configPath}
}

// Load loads config from file with fallback to home directory
func (m *Manager) Load() (*Config, error) {
	cfg := &Config{}

	// Try to load from environment first
	cfg.APIToken = os.Getenv("CLICKUP_API_TOKEN")
	cfg.WorkspaceID = os.Getenv("CLICKUP_WORKSPACE_ID")
	cfg.UserID = os.Getenv("CLICKUP_USER_ID")

	// Determine config file path
	configPath := m.resolveConfigPath()
	if configPath == "" {
		// No config file found, return with env vars only
		if cfg.APIToken == "" {
			return nil, ErrConfigNotFound
		}
		return cfg, nil
	}

	// Load from file
	fileCfg, err := m.loadFromFile(configPath)
	if err != nil {
		return nil, err
	}

	// Merge: env vars take precedence
	if cfg.APIToken == "" {
		cfg.APIToken = fileCfg.APIToken
	}
	if cfg.WorkspaceID == "" {
		cfg.WorkspaceID = fileCfg.WorkspaceID
	}
	if cfg.UserID == "" {
		cfg.UserID = fileCfg.UserID
	}

	return cfg, nil
}

// Save saves config to file
func (m *Manager) Save(cfg *Config) error {
	path := m.configPath
	if path == "" {
		path = DefaultConfigFileName
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// resolveConfigPath finds the config file path
func (m *Manager) resolveConfigPath() string {
	// Custom path specified
	if m.configPath != "" {
		if _, err := os.Stat(m.configPath); err == nil {
			return m.configPath
		}
		return ""
	}

	// Current directory
	if _, err := os.Stat(DefaultConfigFileName); err == nil {
		return DefaultConfigFileName
	}

	// Home directory fallback
	home, err := os.UserHomeDir()
	if err == nil {
		homePath := filepath.Join(home, DefaultConfigFileName)
		if _, err := os.Stat(homePath); err == nil {
			return homePath
		}
	}

	return ""
}

// loadFromFile loads config from a specific file
func (m *Manager) loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetConfigPath returns the path where config will be saved
func (m *Manager) GetConfigPath() string {
	if m.configPath != "" {
		return m.configPath
	}
	return DefaultConfigFileName
}
