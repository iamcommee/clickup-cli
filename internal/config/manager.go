package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var (
	// ErrConfigNotFound is returned when no config file is found
	ErrConfigNotFound = errors.New("config file not found. Run 'clickup install' to set up")
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

// Save saves config to file with encrypted API token
func (m *Manager) Save(cfg *Config) error {
	path := m.configPath
	if path == "" {
		path = DefaultConfigPath()
	}

	// Ensure config directory exists
	dir := filepath.Dir(path)
	if err := EnsureConfigDir(dir); err != nil {
		return err
	}

	// Create a copy with encrypted token
	saveCfg := &Config{
		WorkspaceID: cfg.WorkspaceID,
		UserID:      cfg.UserID,
	}

	// Encrypt the API token
	if cfg.APIToken != "" {
		encryptedToken, err := Encrypt(cfg.APIToken)
		if err != nil {
			return err
		}
		saveCfg.APIToken = encryptedToken
	}

	data, err := json.MarshalIndent(saveCfg, "", "  ")
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

	// Find config file in standard locations
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return FindConfigFile(home)
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

	// Decrypt the API token if it's encrypted
	if cfg.APIToken != "" {
		if IsEncrypted(cfg.APIToken) {
			// Try primary key first
			decrypted, err := Decrypt(cfg.APIToken)
			if err != nil {
				// Primary key failed, try fallback keys (legacy hostname variants)
				decrypted, err = DecryptWithFallback(cfg.APIToken)
				if err != nil {
					return nil, err
				}
				// Fallback succeeded - re-encrypt with the current stable key
				cfg.APIToken = decrypted
				_ = m.migrateToEncrypted(path, &cfg)
			} else {
				cfg.APIToken = decrypted
			}
		} else {
			// Plain text token found - migrate it to encrypted format
			if err := m.migrateToEncrypted(path, &cfg); err != nil {
				// Log warning but continue - migration failure shouldn't block usage
				// The token still works, just not encrypted yet
			}
		}
	}

	return &cfg, nil
}

// migrateToEncrypted encrypts a plain text token and saves the config
func (m *Manager) migrateToEncrypted(path string, cfg *Config) error {
	encryptedToken, err := Encrypt(cfg.APIToken)
	if err != nil {
		return err
	}

	saveCfg := &Config{
		APIToken:    encryptedToken,
		WorkspaceID: cfg.WorkspaceID,
		UserID:      cfg.UserID,
	}

	data, err := json.MarshalIndent(saveCfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// GetConfigPath returns the path where config will be saved
func (m *Manager) GetConfigPath() string {
	if m.configPath != "" {
		return m.configPath
	}
	return DefaultConfigPath()
}
