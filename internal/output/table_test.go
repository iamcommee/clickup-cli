package output

import (
	"bytes"
	"strings"
	"testing"

	"clickup-cli/pkg/models"
)

func TestFormatTaskWithParent(t *testing.T) {
	t.Run("subtask shows parent info", func(t *testing.T) {
		parentID := "86ewrup4j"
		task := &models.Task{
			ID:          "child123",
			CustomID:    "HGAI-1823",
			Name:        "Child Task",
			Status:      models.Status{Status: "open", Color: "#000"},
			DateCreated: "1700000000000",
			DateUpdated: "1700000000000",
			Parent:      &parentID,
			ParentTask: &models.Task{
				ID:          "86ewrup4j",
				CustomID:    "HGAI-1382",
				Name:        "Parent Task Name",
				Description: "Parent task description content",
			},
			List:   models.ListInfo{Name: "Test List"},
			Folder: models.FolderInfo{Name: "Test Folder"},
			Space:  models.SpaceInfo{ID: "1"},
		}

		var buf bytes.Buffer
		f := &DetailFormatter{}
		err := f.FormatTask(&buf, task)
		if err != nil {
			t.Fatalf("FormatTask failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "Parent") {
			t.Error("expected output to contain 'Parent' label")
		}
		if !strings.Contains(output, "HGAI-1382") {
			t.Error("expected output to contain parent custom ID 'HGAI-1382'")
		}
		if !strings.Contains(output, "Parent Task Name") {
			t.Error("expected output to contain parent task name")
		}
		if !strings.Contains(output, "Parent Description (HGAI-1382)") {
			t.Error("expected output to contain 'Parent Description' section")
		}
		if !strings.Contains(output, "Parent task description content") {
			t.Error("expected output to contain parent description text")
		}
	})

	t.Run("top-level task does not show parent", func(t *testing.T) {
		task := &models.Task{
			ID:          "top123",
			CustomID:    "HGAI-1382",
			Name:        "Top Level Task",
			Status:      models.Status{Status: "open", Color: "#000"},
			DateCreated: "1700000000000",
			DateUpdated: "1700000000000",
			List:        models.ListInfo{Name: "Test List"},
			Folder:      models.FolderInfo{Name: "Test Folder"},
			Space:       models.SpaceInfo{ID: "1"},
		}

		var buf bytes.Buffer
		f := &DetailFormatter{}
		err := f.FormatTask(&buf, task)
		if err != nil {
			t.Fatalf("FormatTask failed: %v", err)
		}

		output := buf.String()
		if strings.Contains(output, "Parent") {
			t.Error("top-level task should NOT contain 'Parent' label")
		}
	})
}
