package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// FormatTasks formats tasks as JSON
func (f *JSONFormatter) FormatTasks(w io.Writer, tasks interface{}) error {
	return f.writeJSON(w, tasks)
}

// FormatTask formats a single task as JSON
func (f *JSONFormatter) FormatTask(w io.Writer, task interface{}) error {
	return f.writeJSON(w, task)
}

// writeJSON writes the data as formatted JSON
func (f *JSONFormatter) writeJSON(w io.Writer, data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	_, err = fmt.Fprintln(w, string(output))
	return err
}
