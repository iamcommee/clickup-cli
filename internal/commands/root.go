package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"clickup-cli/internal/api"
	"clickup-cli/internal/config"
)

var (
	// Global flags
	cfgFile     string
	workspaceID string
	debug       bool

	// Shared instances
	cfg       *config.Config
	apiClient *api.Client
)

// rootCmd is the base command
var rootCmd = &cobra.Command{
	Use:   "clickup",
	Short: "ClickUp CLI - Manage tasks from the command line",
	Long: `ClickUp CLI is a command-line tool for interacting with ClickUp.

List your assigned tasks, view task details, and manage your workflow
directly from the terminal.

Get started:
  clickup install       Set up with your API token
  clickup tasks         List your assigned tasks
  clickup tasks <ID>    View task details

Management:
  clickup config        Show current configuration
  clickup uninstall     Remove all configuration

Config: ~/.config/clickup/config.json (encrypted)`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default ~/.config/clickup/config.json)")
	rootCmd.PersistentFlags().StringVarP(&workspaceID, "workspace", "w", "", "override workspace ID")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
}

// loadConfig loads the configuration
func loadConfig() error {
	mgr := config.NewManager(cfgFile)
	var err error
	cfg, err = mgr.Load()
	if err != nil {
		return err
	}

	// Override workspace ID if specified
	if workspaceID != "" {
		cfg.WorkspaceID = workspaceID
	}

	return nil
}

// getAPIClient returns the API client, initializing if needed
func getAPIClient() (*api.Client, error) {
	if apiClient != nil {
		return apiClient, nil
	}

	if err := loadConfig(); err != nil {
		return nil, err
	}

	apiClient = api.NewClient(cfg.APIToken)
	apiClient.SetDebug(debug)

	return apiClient, nil
}

// getConfig returns the loaded config
func getConfig() (*config.Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	if err := loadConfig(); err != nil {
		return nil, err
	}

	return cfg, nil
}
