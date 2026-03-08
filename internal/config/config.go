package config

// Config holds the CLI configuration
type Config struct {
	APIToken    string `json:"api_token"`
	WorkspaceID string `json:"workspace_id"`
	UserID      string `json:"user_id"`
}
