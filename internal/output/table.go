package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"

	"clickup-cli/pkg/models"
)

// wrapText wraps text to the specified width, preserving existing line breaks
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	var result strings.Builder
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(wrapLine(line, width))
	}

	return result.String()
}

// wrapLine wraps a single line of text to the specified width
func wrapLine(line string, width int) string {
	if len(line) <= width {
		return line
	}

	var result strings.Builder
	words := strings.Fields(line)
	currentLen := 0

	for _, word := range words {
		wordLen := len(word)

		if currentLen == 0 {
			result.WriteString(word)
			currentLen = wordLen
		} else if currentLen+1+wordLen <= width {
			result.WriteString(" ")
			result.WriteString(word)
			currentLen += 1 + wordLen
		} else {
			result.WriteString("\n")
			result.WriteString(word)
			currentLen = wordLen
		}
	}

	return result.String()
}

// TableFormatter formats output as a table
type TableFormatter struct{}

// FormatTasks formats a list of tasks as a table, sorted by ID
func (f *TableFormatter) FormatTasks(w io.Writer, tasks interface{}) error {
	taskList, ok := tasks.([]models.Task)
	if !ok {
		return fmt.Errorf("invalid task list type")
	}

	if len(taskList) == 0 {
		fmt.Fprintln(w, "No tasks found.")
		return nil
	}

	sortTasksByID(taskList)

	table := createTable(w, []string{"ID", "Name", "Status", "Priority", "Due Date"})

	for _, task := range taskList {
		appendTaskRow(table, task, false)
		for _, subtask := range task.Subtasks {
			appendTaskRow(table, subtask, true)
		}
	}

	return table.Render()
}

// FormatTask formats a single task with full details
func (f *TableFormatter) FormatTask(w io.Writer, task interface{}) error {
	return (&DetailFormatter{}).FormatTask(w, task)
}

// FormatComments formats a list of comments as a table
func (f *TableFormatter) FormatComments(w io.Writer, comments interface{}) error {
	return (&DetailFormatter{}).FormatComments(w, comments)
}

// DetailFormatter formats a single task with full details
type DetailFormatter struct{}

// FormatTasks delegates to TableFormatter for task lists
func (f *DetailFormatter) FormatTasks(w io.Writer, tasks interface{}) error {
	return (&TableFormatter{}).FormatTasks(w, tasks)
}

// FormatTask formats a single task with full details
func (f *DetailFormatter) FormatTask(w io.Writer, task interface{}) error {
	t, ok := task.(*models.Task)
	if !ok {
		return fmt.Errorf("invalid task type")
	}

	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)

	// Header
	bold.Fprintf(w, "%s: %s\n", t.GetDisplayID(), t.Name)
	fmt.Fprintln(w, strings.Repeat("-", 60))

	// Basic info
	printField(w, cyan, "Status", t.Status.Status)
	printField(w, cyan, "Priority", getPriorityString(t.Priority))
	printField(w, cyan, "Created", t.DateCreated.Format("2006-01-02 15:04"))
	printField(w, cyan, "Updated", t.DateUpdated.Format("2006-01-02 15:04"))

	if t.DueDate != nil && !t.DueDate.IsZero() {
		printField(w, cyan, "Due", t.DueDate.Format("2006-01-02 15:04"))
	}
	if t.StartDate != nil && !t.StartDate.IsZero() {
		printField(w, cyan, "Start", t.StartDate.Format("2006-01-02 15:04"))
	}

	// Location
	printField(w, cyan, "List", t.List.Name)
	if t.Folder.Name != "" && !t.Folder.Hidden {
		printField(w, cyan, "Folder", t.Folder.Name)
	}

	// People
	if len(t.Assignees) > 0 {
		names := make([]string, len(t.Assignees))
		for i, a := range t.Assignees {
			names[i] = a.Username
		}
		printField(w, cyan, "Assignees", strings.Join(names, ", "))
	}

	// Tags
	if len(t.Tags) > 0 {
		tagNames := make([]string, len(t.Tags))
		for i, tag := range t.Tags {
			tagNames[i] = tag.Name
		}
		printField(w, cyan, "Tags", strings.Join(tagNames, ", "))
	}

	// URL
	printField(w, cyan, "URL", t.URL)

	// Description
	if desc := t.GetDescription(); desc != "" {
		fmt.Fprintln(w)
		cyan.Fprintln(w, "Description:")
		fmt.Fprintln(w, desc)
	}

	return nil
}

// FormatComments formats a list of comments
func (f *DetailFormatter) FormatComments(w io.Writer, comments interface{}) error {
	commentList, ok := comments.([]models.Comment)
	if !ok {
		return fmt.Errorf("invalid comment list type")
	}

	if len(commentList) == 0 {
		fmt.Fprintln(w, "No comments.")
		return nil
	}

	// Sort by date ascending (oldest first)
	sort.Slice(commentList, func(i, j int) bool {
		return commentList[i].DateCreated.Time().Before(commentList[j].DateCreated.Time())
	})

	bold := color.New(color.Bold)
	cyan := color.New(color.FgCyan)
	dim := color.New(color.Faint)

	bold.Fprintln(w, "Comments:")
	fmt.Fprintln(w, strings.Repeat("-", 60))

	// Wrap comments to slightly less than terminal width
	wrapWidth := getTerminalWidth() - 4

	for i, c := range commentList {
		// Author and date
		cyan.Fprintf(w, "%s", c.User.Username)
		dim.Fprintf(w, " - %s", c.DateCreated.Format("2006-01-02 15:04"))
		if c.Resolved {
			dim.Fprint(w, " [resolved]")
		}
		fmt.Fprintln(w)

		// Content with word wrapping
		fmt.Fprintln(w, wrapText(c.CommentText, wrapWidth))

		if i < len(commentList)-1 {
			dim.Fprintln(w, strings.Repeat("·", 40))
		}
	}

	return nil
}

// Helper functions

func sortTasksByID(tasks []models.Task) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].GetDisplayID() < tasks[j].GetDisplayID()
	})
}

func createTable(w io.Writer, headers []string) *tablewriter.Table {
	return tablewriter.NewTable(w,
		tablewriter.WithHeader(headers),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithHeaderAlignment(tw.AlignLeft),
		tablewriter.WithBorders(tw.Border{Left: tw.Off, Right: tw.Off, Top: tw.Off, Bottom: tw.Off}),
		tablewriter.WithRendition(tw.Rendition{Symbols: tw.NewSymbols(tw.StyleNone)}),
	)
}

func appendTaskRow(table *tablewriter.Table, task models.Task, isSubtask bool) {
	name := task.Name
	if isSubtask {
		name = "  └─ " + name
	}

	table.Append([]string{
		task.GetDisplayID(),
		truncate(name, 50),
		task.Status.Status,
		getPriorityString(task.Priority),
		formatDueDate(task.DueDate),
	})
}

func printField(w io.Writer, c *color.Color, label, value string) {
	c.Fprintf(w, "%-10s ", label+":")
	fmt.Fprintln(w, value)
}

func getPriorityString(p *models.Priority) string {
	if p != nil {
		return p.Priority
	}
	return "-"
}

func formatDueDate(d *models.Timestamp) string {
	if d != nil && !d.IsZero() {
		return d.Format("2006-01-02")
	}
	return "-"
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
