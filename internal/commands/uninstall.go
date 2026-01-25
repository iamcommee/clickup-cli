package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"clickup-cli/internal/config"
)

var (
	forceUninstall bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove ClickUp CLI configuration",
	Long: `Remove all ClickUp CLI configuration files.

This will delete:
  - ~/.config/clickup/config.json (current config)
  - ~/.clickup.json (legacy config, if exists)

Your API token will be removed. You'll need to run 'clickup install'
again to use the CLI.`,
	RunE: runUninstall,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "skip confirmation prompt")
}

func runUninstall(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check what files exist
	files := config.GetConfigFiles(home)
	if len(files) == 0 {
		fmt.Println("No configuration files found. Nothing to remove.")
		return nil
	}

	// Show what will be removed
	fmt.Println("The following files will be removed:")
	for _, f := range files {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println()

	// Confirm unless --force
	if !forceUninstall {
		fmt.Print("Are you sure you want to remove these files? [y/N] ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	// Remove files
	removed, err := config.Cleanup(home)
	if err != nil {
		return fmt.Errorf("failed to remove config: %w", err)
	}

	fmt.Println()
	fmt.Println("Removed:")
	for _, f := range removed {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println()
	fmt.Println("Configuration removed successfully.")
	fmt.Println("Run 'clickup install' to set up again.")

	return nil
}
