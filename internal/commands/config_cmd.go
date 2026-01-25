package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"clickup-cli/internal/config"
)

// configCmd shows the current configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Long: `Display the current configuration settings.

Shows your workspace ID, user ID, and masked API token.
The full API token is never displayed for security.`,
	RunE: runConfigShow,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	fmt.Println("Current Configuration")
	fmt.Println("=====================")
	fmt.Println()

	// Mask the token
	maskedToken := cfg.APIToken
	if len(maskedToken) > 8 {
		maskedToken = maskedToken[:4] + strings.Repeat("*", len(maskedToken)-8) + maskedToken[len(maskedToken)-4:]
	}

	fmt.Printf("API Token:    %s\n", maskedToken)
	fmt.Printf("Workspace ID: %s\n", cfg.WorkspaceID)
	fmt.Printf("User ID:      %s\n", cfg.UserID)

	mgr := config.NewManager(cfgFile)
	fmt.Printf("\nConfig file:  %s\n", mgr.GetConfigPath())

	return nil
}
