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
	installTokenFlag string
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and configure CLI",
	Long: `Install and configure ClickUp CLI.

This will prompt for your API token and automatically fetch your
user ID and workspace ID.

Your API token is encrypted before being stored in:
  ~/.config/clickup/config.json

Get your API token from: https://app.clickup.com/settings/apps

Examples:
  clickup install                     Interactive setup
  clickup install --token pk_xxxxx    Non-interactive with token`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&installTokenFlag, "token", "t", "", "API token")
}

func runInstall(cmd *cobra.Command, args []string) error {
	token := installTokenFlag

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
	fmt.Println("  clickup tasks          List your assigned tasks")
	fmt.Println("  clickup tasks <ID>     View task details")

	return nil
}
