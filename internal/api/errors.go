package api

import (
	"fmt"
)

// APIError represents an error from the ClickUp API
type APIError struct {
	StatusCode int
	Message    string
	Err        string `json:"err"`
	ECODE      string `json:"ECODE"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
	}
	if e.Err != "" {
		return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Err)
	}
	return fmt.Sprintf("API error (status %d)", e.StatusCode)
}
