package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestManager_SaveEncryptsToken(t *testing.T) {
	// Create temp directory with config structure
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(configDir, 0700)
	configPath := filepath.Join(configDir, "config.json")

	mgr := NewManager(configPath)
	cfg := &Config{
		APIToken:    "pk_test_token_12345",
		WorkspaceID: "123456",
		UserID:      "789",
	}

	err := mgr.Save(cfg)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Read raw file content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Parse the JSON to check the token is encrypted
	var rawCfg map[string]string
	if err := json.Unmarshal(data, &rawCfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Token in file should be encrypted (have enc: prefix)
	if !IsEncrypted(rawCfg["api_token"]) {
		t.Errorf("token in file should be encrypted, got: %s", rawCfg["api_token"])
	}

	// Token should NOT be plain text
	if rawCfg["api_token"] == "pk_test_token_12345" {
		t.Error("token should not be stored as plain text")
	}
}

func TestManager_LoadDecryptsToken(t *testing.T) {
	// Create temp directory with config structure
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(configDir, 0700)
	configPath := filepath.Join(configDir, "config.json")

	// Save config
	mgr := NewManager(configPath)
	originalCfg := &Config{
		APIToken:    "pk_secret_token_xyz",
		WorkspaceID: "workspace123",
		UserID:      "user456",
	}

	err := mgr.Save(originalCfg)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load config
	loadedCfg, err := mgr.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Token should be decrypted
	if loadedCfg.APIToken != "pk_secret_token_xyz" {
		t.Errorf("expected 'pk_secret_token_xyz', got '%s'", loadedCfg.APIToken)
	}
}

func TestManager_LoadMigratesPlainTextToken(t *testing.T) {
	// Create temp directory with config structure
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "clickup")
	os.MkdirAll(configDir, 0700)
	configPath := filepath.Join(configDir, "config.json")

	// Write config with plain text token (simulating unencrypted config)
	plainCfg := map[string]string{
		"api_token":    "pk_plain_text_token",
		"workspace_id": "ws123",
		"user_id":      "u456",
	}
	data, _ := json.MarshalIndent(plainCfg, "", "  ")
	os.WriteFile(configPath, data, 0600)

	// Load config - should work and migrate to encrypted
	mgr := NewManager(configPath)
	cfg, err := mgr.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Token should be readable
	if cfg.APIToken != "pk_plain_text_token" {
		t.Errorf("expected 'pk_plain_text_token', got '%s'", cfg.APIToken)
	}

	// File should now be encrypted (after migration)
	data, _ = os.ReadFile(configPath)
	var rawCfg map[string]string
	json.Unmarshal(data, &rawCfg)

	if !IsEncrypted(rawCfg["api_token"]) {
		t.Error("token should be encrypted after migration")
	}
}
