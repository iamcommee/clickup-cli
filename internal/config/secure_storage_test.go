package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGetOrCreateKeyFile(t *testing.T) {
	dir := t.TempDir()

	key, err := getOrCreateKeyFile(dir)
	if err != nil {
		t.Fatalf("getOrCreateKeyFile failed on first call: %v", err)
	}

	if len(key) != 32 {
		t.Errorf("expected 32-byte key, got %d bytes", len(key))
	}

	keyPath := filepath.Join(dir, ".key")
	if _, statErr := os.Stat(keyPath); os.IsNotExist(statErr) {
		t.Errorf(".key file was not created at %s", keyPath)
	}

	// Second call should return identical key
	key2, err := getOrCreateKeyFile(dir)
	if err != nil {
		t.Fatalf("getOrCreateKeyFile failed on second call: %v", err)
	}

	if !bytes.Equal(key, key2) {
		t.Error("expected the same key on subsequent calls, got a different key")
	}
}

func TestGetOrCreateKeyFile_ExistingKey(t *testing.T) {
	dir := t.TempDir()

	knownKey := make([]byte, 32)
	for i := range knownKey {
		knownKey[i] = byte(i + 1)
	}
	keyPath := filepath.Join(dir, ".key")
	if err := os.WriteFile(keyPath, knownKey, 0600); err != nil {
		t.Fatalf("failed to write pre-existing .key file: %v", err)
	}

	returned, err := getOrCreateKeyFile(dir)
	if err != nil {
		t.Fatalf("getOrCreateKeyFile failed: %v", err)
	}

	if !bytes.Equal(returned, knownKey) {
		t.Errorf("expected pre-existing key to be returned unchanged")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "pk_12345678_abcdefghijklmnop"

	encrypted, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == plaintext {
		t.Error("encrypted text should not equal plaintext")
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected '%s', got '%s'", plaintext, decrypted)
	}
}

func TestEncryptDecrypt_EmptyString(t *testing.T) {
	plaintext := ""

	encrypted, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected empty string, got '%s'", decrypted)
	}
}

func TestDecrypt_InvalidData(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}

	_, err = Decrypt("c2hvcnQ=") // "short" in base64
	if err == nil {
		t.Error("expected error for data too short")
	}
}

func TestIsEncrypted(t *testing.T) {
	if IsEncrypted("pk_12345678_abcdefghijklmnop") {
		t.Error("plain text should not be detected as encrypted")
	}

	encrypted, _ := Encrypt("pk_12345678_abcdefghijklmnop")
	if !IsEncrypted(encrypted) {
		t.Error("encrypted text should be detected as encrypted")
	}
}

func TestEncryptDecryptWithKeyFile(t *testing.T) {
	dir := t.TempDir()

	plaintext := "pk_keyfile_token_abc123"

	key, err := getOrCreateKeyFile(dir)
	if err != nil {
		t.Fatalf("getOrCreateKeyFile failed: %v", err)
	}

	encrypted, err := encryptWithKey(plaintext, key)
	if err != nil {
		t.Fatalf("encryptWithKey failed: %v", err)
	}

	if !IsEncrypted(encrypted) {
		t.Errorf("encrypted value should have the %q prefix", EncryptedPrefix)
	}

	// Re-fetch key to simulate a fresh process
	key2, err := getOrCreateKeyFile(dir)
	if err != nil {
		t.Fatalf("getOrCreateKeyFile failed on second fetch: %v", err)
	}

	decrypted, err := decryptWithKey(encrypted, key2)
	if err != nil {
		t.Fatalf("decryptWithKey failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected %q after decrypt, got %q", plaintext, decrypted)
	}
}
