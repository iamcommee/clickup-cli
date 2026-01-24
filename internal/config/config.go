package config

// Config holds the CLI configuration
type Config struct {
	APIToken    string `json:"api_token"`
	WorkspaceID string `json:"workspace_id"`
	UserID      string `json:"user_id"`
}

// Validate checks if the config has required fields
func (c *Config) Validate() error {
	if c.APIToken == "" {
		return ErrMissingAPIToken
	}
	if c.WorkspaceID == "" {
		return ErrMissingWorkspaceID
	}
	if c.UserID == "" {
		return ErrMissingUserID
	}
	return nil
}

// HasToken returns true if the API token is set
func (c *Config) HasToken() bool {
	return c.APIToken != ""
}
