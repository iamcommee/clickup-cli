package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"clickup-cli/internal/api"
	"clickup-cli/internal/config"
)

var (
	// Config init flags
	tokenFlag string
)

// configCmd is the parent command for config operations
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Initialize and view ClickUp CLI configuration.",
}

// configInitCmd initializes the configuration
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long: `Initialize ClickUp CLI configuration.

This will prompt for your API token and automatically fetch your
user ID and workspace ID.

Get your API token from: https://app.clickup.com/settings/apps

Examples:
  clickup config init                     Interactive setup
  clickup config init --token pk_xxxxx    Non-interactive with token`,
	RunE: runConfigInit,
}

// configShowCmd shows the current configuration
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration settings.",
	RunE:  runConfigShow,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)

	configInitCmd.Flags().StringVarP(&tokenFlag, "token", "t", "", "API token")
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	token := tokenFlag

	// Prompt for token if not provided
	if token == "" {
		fmt.Println("ClickUp CLI Configuration")
		fmt.Println("==========================")
		fmt.Println()
		fmt.Println("Get your API token from: https://app.clickup.com/settings/apps")
		fmt.Println()
		fmt.Print("Enter your API token: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		token = strings.TrimSpace(input)
	}

	if token == "" {
		return fmt.Errorf("API token is required")
	}

	// Validate token by fetching user info
	fmt.Println()
	fmt.Println("Validating token...")

	client := api.NewClient(token)
	client.SetDebug(debug)

	user, err := client.GetUser()
	if err != nil {
		return fmt.Errorf("failed to validate token: %w", err)
	}

	fmt.Printf("Authenticated as: %s (%s)\n", user.Username, user.Email)

	// Fetch workspaces
	teams, err := client.GetTeams()
	if err != nil {
		return fmt.Errorf("failed to fetch workspaces: %w", err)
	}

	if len(teams) == 0 {
		return fmt.Errorf("no workspaces found for this user")
	}

	// Select workspace
	var selectedTeam string
	if len(teams) == 1 {
		selectedTeam = teams[0].ID
		fmt.Printf("Using workspace: %s (%s)\n", teams[0].Name, teams[0].ID)
	} else {
		fmt.Println()
		fmt.Println("Available workspaces:")
		for i, team := range teams {
			fmt.Printf("  [%d] %s (%s)\n", i+1, team.Name, team.ID)
		}
		fmt.Print("Select workspace (number): ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 1 || choice > len(teams) {
			return fmt.Errorf("invalid selection")
		}

		selectedTeam = teams[choice-1].ID
	}

	// Create config
	newCfg := &config.Config{
		APIToken:    token,
		WorkspaceID: selectedTeam,
		UserID:      strconv.Itoa(user.ID),
	}

	// Save config
	mgr := config.NewManager(cfgFile)
	if err := mgr.Save(newCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	fmt.Printf("Configuration saved to: %s\n", mgr.GetConfigPath())
	fmt.Println()
	fmt.Println("You can now use the CLI:")
	fmt.Println("  clickup tasks list     List your assigned tasks")
	fmt.Println("  clickup tasks get ID   View task details")

	return nil
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
