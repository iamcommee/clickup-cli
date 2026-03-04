package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"clickup-cli/internal/api"
	"clickup-cli/internal/output"
)

var (
	statusFilter  string
	limitTasks    int
	outputFormat  string
	includeClosed bool
	openInBrowser bool
	showComments  bool
)

var tasksCmd = &cobra.Command{
	Use:     "tasks [task_id]",
	Aliases: []string{"task", "t"},
	Short:   "List tasks or get task details",
	Long: `List your assigned tasks or get details of a specific task.

Examples:
  clickup tasks                          List all assigned tasks
  clickup tasks --status "in progress"   Filter by status
  clickup tasks --closed                 Include closed tasks
  clickup tasks HGAI-1217                Get task details
  clickup tasks HGAI-1217 --open         Open task in browser
  clickup tasks HGAI-1217 --comments     Show task comments`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTasks,
}

func init() {
	rootCmd.AddCommand(tasksCmd)

	tasksCmd.Flags().StringVarP(&statusFilter, "status", "s", "", "filter by status")
	tasksCmd.Flags().IntVarP(&limitTasks, "limit", "l", 0, "limit number of results")
	tasksCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "output format (table, json)")
	tasksCmd.Flags().BoolVar(&includeClosed, "closed", false, "include closed tasks")
	tasksCmd.Flags().BoolVar(&openInBrowser, "open", false, "open task in browser (only with task_id)")
	tasksCmd.Flags().BoolVar(&showComments, "comments", false, "show task comments (only with task_id)")
}

func runTasks(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return getTask(args[0])
	}
	if showComments {
		return fmt.Errorf("--comments requires a task ID: clickup tasks <task_id> --comments")
	}
	return listTasks()
}

func listTasks() error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	cfg, err := getConfig()
	if err != nil {
		return err
	}

	opts := &api.ListTasksOptions{
		Assignees:     []string{cfg.UserID},
		IncludeClosed: includeClosed,
		Subtasks:      true,
	}

	if statusFilter != "" {
		opts.Statuses = []string{statusFilter}
	}

	tasks, err := client.ListTasks(cfg.WorkspaceID, opts)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	if limitTasks > 0 && len(tasks) > limitTasks {
		tasks = tasks[:limitTasks]
	}

	format := output.Format(outputFormat)
	if format == "" {
		format = output.FormatTable
	}

	return output.GetFormatter(format).FormatTasks(os.Stdout, tasks)
}

func getTask(taskID string) error {
	client, err := getAPIClient()
	if err != nil {
		return err
	}

	cfg, err := getConfig()
	if err != nil {
		return err
	}

	task, err := client.GetTask(taskID, cfg.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	// Fetch parent task info if this is a subtask
	if task.Parent != nil {
		parentTask, err := client.GetTask(*task.Parent, cfg.WorkspaceID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not fetch parent task: %v\n", err)
		} else {
			task.ParentTask = parentTask
		}
	}

	if openInBrowser {
		if err := openURL(task.URL); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not open browser: %v\n", err)
		} else {
			fmt.Printf("Opened in browser: %s\n", task.URL)
			return nil
		}
	}

	format := output.Format(outputFormat)
	if format == "" {
		format = output.FormatDetail
	}

	formatter := output.GetFormatter(format)

	if err := formatter.FormatTask(os.Stdout, task); err != nil {
		return err
	}

	if showComments {
		comments, err := client.GetTaskComments(taskID, cfg.WorkspaceID)
		if err != nil {
			return fmt.Errorf("failed to get comments: %w", err)
		}

		fmt.Fprintln(os.Stdout)
		if err := formatter.FormatComments(os.Stdout, comments); err != nil {
			return err
		}
	}

	return nil
}

func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
