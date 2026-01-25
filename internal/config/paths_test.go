package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigDir(t *testing.T) {
	// Clear XDG_CONFIG_HOME to test default behavior
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Unsetenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	dir := DefaultConfigDir()

	// Should be under home directory
	home, _ := os.UserHomeDir()
	if !strings.HasPrefix(dir, home) {
		t.Errorf("config dir should be under home, got: %s", dir)
	}

	// Should end with clickup
	if !strings.HasSuffix(dir, "clickup") {
		t.Errorf("config dir should end with 'clickup', got: %s", dir)
	}

	// Should contain .config
	if !strings.Contains(dir, ".config") {
		t.Errorf("config dir should contain '.config', got: %s", dir)
	}
}

func TestDefaultConfigDir_WithXDG(t *testing.T) {
	// Set custom XDG_CONFIG_HOME
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", "/custom/config")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	dir := DefaultConfigDir()

	expected := "/custom/config/clickup"
	if dir != expected {
		t.Errorf("expected '%s', got '%s'", expected, dir)
	}
}

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()

	// Should end with config.json
	if !strings.HasSuffix(path, "config.json") {
		t.Errorf("config path should end with 'config.json', got: %s", path)
	}

	// Should contain clickup directory
	if !strings.Contains(path, "clickup") {
		t.Errorf("config path should contain 'clickup', got: %s", path)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "clickup")

	err := EnsureConfigDir(configDir)
	if err != nil {
		t.Fatalf("EnsureConfigDir failed: %v", err)
	}

	// Directory should exist
	info, err := os.Stat(configDir)
	if err != nil {
		t.Fatalf("config dir should exist: %v", err)
	}

	if !info.IsDir() {
		t.Error("config dir should be a directory")
	}

	// Permissions should be 0700 (owner only)
	if info.Mode().Perm() != 0700 {
		t.Errorf("expected permissions 0700, got %o", info.Mode().Perm())
	}
}

func TestFindConfigFile_NewLocation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create new location config
	newConfigDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(newConfigDir, 0700)
	newConfigPath := filepath.Join(newConfigDir, "config.json")
	os.WriteFile(newConfigPath, []byte(`{}`), 0600)

	found := FindConfigFile(tmpDir)
	if found != newConfigPath {
		t.Errorf("expected '%s', got '%s'", newConfigPath, found)
	}
}

func TestFindConfigFile_LegacyLocation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create only legacy location config
	legacyPath := filepath.Join(tmpDir, ".clickup.json")
	os.WriteFile(legacyPath, []byte(`{}`), 0600)

	found := FindConfigFile(tmpDir)
	if found != legacyPath {
		t.Errorf("expected legacy path '%s', got '%s'", legacyPath, found)
	}
}

func TestFindConfigFile_PrefersNewOverLegacy(t *testing.T) {
	tmpDir := t.TempDir()

	// Create both locations
	newConfigDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(newConfigDir, 0700)
	newConfigPath := filepath.Join(newConfigDir, "config.json")
	os.WriteFile(newConfigPath, []byte(`{"new": true}`), 0600)

	legacyPath := filepath.Join(tmpDir, ".clickup.json")
	os.WriteFile(legacyPath, []byte(`{"legacy": true}`), 0600)

	found := FindConfigFile(tmpDir)
	if found != newConfigPath {
		t.Errorf("should prefer new location, expected '%s', got '%s'", newConfigPath, found)
	}
}

func TestFindConfigFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	found := FindConfigFile(tmpDir)
	if found != "" {
		t.Errorf("expected empty string when not found, got '%s'", found)
	}
}
