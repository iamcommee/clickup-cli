package models

import (
	"encoding/json"
	"testing"
)

func TestTaskParentFieldDeserialization(t *testing.T) {
	t.Run("subtask has parent ID", func(t *testing.T) {
		data := `{"id":"abc123","name":"Subtask","parent":"86ewrup4j","status":{"status":"open","color":"#000"},"date_created":"1700000000000","date_updated":"1700000000000","creator":{},"list":{},"folder":{},"space":{}}`
		var task Task
		if err := json.Unmarshal([]byte(data), &task); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if task.Parent == nil {
			t.Fatal("expected Parent to be non-nil for subtask")
		}
		if *task.Parent != "86ewrup4j" {
			t.Errorf("expected parent '86ewrup4j', got '%s'", *task.Parent)
		}
	})

	t.Run("top-level task has null parent", func(t *testing.T) {
		data := `{"id":"abc123","name":"Top Level","parent":null,"status":{"status":"open","color":"#000"},"date_created":"1700000000000","date_updated":"1700000000000","creator":{},"list":{},"folder":{},"space":{}}`
		var task Task
		if err := json.Unmarshal([]byte(data), &task); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if task.Parent != nil {
			t.Errorf("expected Parent to be nil for top-level task, got '%s'", *task.Parent)
		}
	})

	t.Run("task without parent field", func(t *testing.T) {
		data := `{"id":"abc123","name":"No Parent Field","status":{"status":"open","color":"#000"},"date_created":"1700000000000","date_updated":"1700000000000","creator":{},"list":{},"folder":{},"space":{}}`
		var task Task
		if err := json.Unmarshal([]byte(data), &task); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if task.Parent != nil {
			t.Error("expected Parent to be nil when field is absent")
		}
	})
}

func TestTaskParentTaskField(t *testing.T) {
	t.Run("ParentTask is not serialized to JSON", func(t *testing.T) {
		parentID := "86ewrup4j"
		task := Task{
			ID:     "abc123",
			Name:   "Subtask",
			Parent: &parentID,
			ParentTask: &Task{
				ID:       "86ewrup4j",
				CustomID: "HGAI-1382",
				Name:     "Parent Task Name",
			},
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}
		// ParentTask should NOT appear in JSON output
		var m map[string]interface{}
		json.Unmarshal(data, &m)
		if _, exists := m["parent_task"]; exists {
			t.Error("ParentTask should not be serialized to JSON")
		}
	})
}
