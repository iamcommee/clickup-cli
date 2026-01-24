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

// IsNotFound returns true if the error is a 404
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

// IsForbidden returns true if the error is a 403
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == 403
}

// IsRateLimited returns true if the error is a 429
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == 429
}
