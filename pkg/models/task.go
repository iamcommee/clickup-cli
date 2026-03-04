package models

import (
	"strconv"
	"time"
)

// Task represents a ClickUp task
type Task struct {
	ID                  string     `json:"id"`
	CustomID            string     `json:"custom_id,omitempty"`
	Name                string     `json:"name"`
	Description         string     `json:"description,omitempty"`
	MarkdownDescription string     `json:"markdown_description,omitempty"`
	TextContent         string     `json:"text_content,omitempty"`
	Status              Status     `json:"status"`
	Priority            *Priority  `json:"priority,omitempty"`
	DueDate             *Timestamp `json:"due_date,omitempty"`
	StartDate           *Timestamp `json:"start_date,omitempty"`
	DateCreated         Timestamp  `json:"date_created"`
	DateUpdated         Timestamp  `json:"date_updated"`
	DateClosed          *Timestamp `json:"date_closed,omitempty"`
	Creator             User       `json:"creator"`
	Assignees           []User     `json:"assignees"`
	Tags                []Tag      `json:"tags"`
	List                ListInfo   `json:"list"`
	Folder              FolderInfo `json:"folder"`
	Space               SpaceInfo  `json:"space"`
	URL                 string     `json:"url"`
	Parent              *string    `json:"parent,omitempty"`
	Subtasks            []Task     `json:"subtasks,omitempty"`
	ParentTask          *Task      `json:"-"`
}

// GetDisplayID returns the custom ID if available, otherwise the regular ID
func (t *Task) GetDisplayID() string {
	if t.CustomID != "" {
		return t.CustomID
	}
	return t.ID
}

// GetDescription returns the best available description
func (t *Task) GetDescription() string {
	if t.MarkdownDescription != "" {
		return t.MarkdownDescription
	}
	if t.Description != "" {
		return t.Description
	}
	return t.TextContent
}

// Status represents a task status
type Status struct {
	ID         string `json:"id,omitempty"`
	Status     string `json:"status"`
	Color      string `json:"color"`
	Type       string `json:"type,omitempty"`
	Orderindex int    `json:"orderindex"`
}

// Priority represents a task priority
type Priority struct {
	ID         string `json:"id"`
	Priority   string `json:"priority"`
	Color      string `json:"color"`
	Orderindex string `json:"orderindex"`
}

// Tag represents a task tag
type Tag struct {
	Name    string `json:"name"`
	TagFg   string `json:"tag_fg"`
	TagBg   string `json:"tag_bg"`
	Creator int    `json:"creator"`
}

// ListInfo contains basic list information
type ListInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

// FolderInfo contains basic folder information
type FolderInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
	Access bool   `json:"access"`
}

// SpaceInfo contains basic space information
type SpaceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// Timestamp is a custom type for ClickUp timestamps (milliseconds as string)
type Timestamp string

// Time converts the timestamp to time.Time
func (t Timestamp) Time() time.Time {
	if t == "" {
		return time.Time{}
	}
	ms, err := strconv.ParseInt(string(t), 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.UnixMilli(ms)
}

// IsZero returns true if the timestamp is empty or zero
func (t Timestamp) IsZero() bool {
	return t == "" || t == "0"
}

// Format returns a formatted string representation
func (t Timestamp) Format(layout string) string {
	if t.IsZero() {
		return ""
	}
	return t.Time().Format(layout)
}

// TasksResponse is the API response for listing tasks
type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
