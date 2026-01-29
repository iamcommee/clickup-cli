package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"clickup-cli/pkg/models"
)

// GetTaskComments gets comments for a task
// teamID is required when using custom task IDs (e.g., HGAI-1217)
func (c *Client) GetTaskComments(taskID string, teamID string) ([]models.Comment, error) {
	path := fmt.Sprintf("/task/%s/comment", taskID)

	query := url.Values{}
	if isCustomTaskID(taskID) && teamID != "" {
		query.Set("custom_task_ids", "true")
		query.Set("team_id", teamID)
	}

	body, err := c.Get(path, query)
	if err != nil {
		return nil, err
	}

	var resp models.CommentsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Comments, nil
}
