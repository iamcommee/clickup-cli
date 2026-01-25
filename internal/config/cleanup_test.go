package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanup_RemovesConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(configDir, 0700)
	configPath := filepath.Join(configDir, "config.json")
	os.WriteFile(configPath, []byte(`{}`), 0600)

	removed, err := Cleanup(tmpDir)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Should have removed 1 file
	if len(removed) != 1 {
		t.Errorf("expected 1 removed file, got %d", len(removed))
	}

	// Config file should not exist
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Error("config file should be removed")
	}

	// Config directory should also be removed (empty)
	if _, err := os.Stat(configDir); !os.IsNotExist(err) {
		t.Error("empty config directory should be removed")
	}
}

func TestCleanup_RemovesLegacyConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Create legacy config file
	legacyPath := filepath.Join(tmpDir, ".clickup.json")
	os.WriteFile(legacyPath, []byte(`{}`), 0600)

	removed, err := Cleanup(tmpDir)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Should have removed 1 file
	if len(removed) != 1 {
		t.Errorf("expected 1 removed file, got %d", len(removed))
	}

	// Legacy file should not exist
	if _, err := os.Stat(legacyPath); !os.IsNotExist(err) {
		t.Error("legacy config file should be removed")
	}
}

func TestCleanup_RemovesBothConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create both config files
	configDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(configDir, 0700)
	configPath := filepath.Join(configDir, "config.json")
	os.WriteFile(configPath, []byte(`{}`), 0600)

	legacyPath := filepath.Join(tmpDir, ".clickup.json")
	os.WriteFile(legacyPath, []byte(`{}`), 0600)

	removed, err := Cleanup(tmpDir)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Should have removed 2 files
	if len(removed) != 2 {
		t.Errorf("expected 2 removed files, got %d", len(removed))
	}
}

func TestCleanup_NoConfigExists(t *testing.T) {
	tmpDir := t.TempDir()

	removed, err := Cleanup(tmpDir)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Should have removed 0 files
	if len(removed) != 0 {
		t.Errorf("expected 0 removed files, got %d", len(removed))
	}
}

func TestCleanup_PreservesNonEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config directory with extra file
	configDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(configDir, 0700)
	configPath := filepath.Join(configDir, "config.json")
	os.WriteFile(configPath, []byte(`{}`), 0600)
	otherFile := filepath.Join(configDir, "other.txt")
	os.WriteFile(otherFile, []byte(`keep me`), 0600)

	_, err := Cleanup(tmpDir)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Config file should be removed
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Error("config file should be removed")
	}

	// Directory should still exist (not empty)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("config directory should be preserved (has other files)")
	}

	// Other file should still exist
	if _, err := os.Stat(otherFile); os.IsNotExist(err) {
		t.Error("other file should be preserved")
	}
}
