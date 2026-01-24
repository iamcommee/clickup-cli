package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"clickup-cli/pkg/models"
)

// ListTasksOptions contains options for listing tasks
type ListTasksOptions struct {
	Assignees     []string
	Statuses      []string
	IncludeClosed bool
	Subtasks      bool
	Page          int
}

// ListTasks lists tasks for a team/workspace
func (c *Client) ListTasks(teamID string, opts *ListTasksOptions) ([]models.Task, error) {
	path := fmt.Sprintf("/team/%s/task", teamID)

	query := url.Values{}
	if opts != nil {
		for _, assignee := range opts.Assignees {
			query.Add("assignees[]", assignee)
		}
		for _, status := range opts.Statuses {
			query.Add("statuses[]", status)
		}
		if opts.IncludeClosed {
			query.Set("include_closed", "true")
		}
		if opts.Subtasks {
			query.Set("subtasks", "true")
		}
		if opts.Page > 0 {
			query.Set("page", strconv.Itoa(opts.Page))
		}
	}

	body, err := c.Get(path, query)
	if err != nil {
		return nil, err
	}

	var resp models.TasksResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Tasks, nil
}

// GetTask gets a single task by ID
// teamID is required when using custom task IDs (e.g., HGAI-1217)
func (c *Client) GetTask(taskID string, teamID string) (*models.Task, error) {
	path := fmt.Sprintf("/task/%s", taskID)

	query := url.Values{}
	query.Set("include_subtasks", "true")
	query.Set("include_markdown_description", "true")

	// Check if this looks like a custom task ID (contains non-numeric characters)
	if isCustomTaskID(taskID) && teamID != "" {
		query.Set("custom_task_ids", "true")
		query.Set("team_id", teamID)
	}

	body, err := c.Get(path, query)
	if err != nil {
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal(body, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// isCustomTaskID checks if the ID looks like a custom task ID
// Custom IDs typically contain letters or hyphens (e.g., HGAI-1217)
// Native IDs are alphanumeric strings without hyphens (e.g., 86a3xyzw)
func isCustomTaskID(id string) bool {
	for _, c := range id {
		if c == '-' {
			return true
		}
	}
	return false
}
