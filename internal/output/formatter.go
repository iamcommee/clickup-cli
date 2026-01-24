package output

import (
	"io"
)

// Format represents the output format type
type Format string

const (
	FormatTable  Format = "table"
	FormatJSON   Format = "json"
	FormatDetail Format = "detail"
)

// Formatter defines the interface for output formatters
type Formatter interface {
	// FormatTasks formats a list of tasks
	FormatTasks(w io.Writer, tasks interface{}) error
	// FormatTask formats a single task
	FormatTask(w io.Writer, task interface{}) error
}

// GetFormatter returns the appropriate formatter for the given format
func GetFormatter(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatDetail:
		return &DetailFormatter{}
	default:
		return &TableFormatter{}
	}
}
