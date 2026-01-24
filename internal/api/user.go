package api

import (
	"encoding/json"

	"clickup-cli/pkg/models"
)

// GetUser gets the current authenticated user
func (c *Client) GetUser() (*models.User, error) {
	body, err := c.Get("/user", nil)
	if err != nil {
		return nil, err
	}

	var resp models.UserResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp.User, nil
}

// GetTeams gets all workspaces/teams for the current user
func (c *Client) GetTeams() ([]models.Team, error) {
	body, err := c.Get("/team", nil)
	if err != nil {
		return nil, err
	}

	var resp models.TeamsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Teams, nil
}
