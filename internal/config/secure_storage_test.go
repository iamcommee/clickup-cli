package config

import (
	"testing"
)

func TestNormalizeHostname(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"commees-MacBook-Pro.local", "commees-MacBook-Pro"},
		{"commees-MacBook-Pro", "commees-MacBook-Pro"},
		{"my-host.local", "my-host"},
		{"my-host.example.com", "my-host.example.com"},
		{"localhost", "localhost"},
		{"", ""},
	}

	for _, tt := range tests {
		result := normalizeHostname(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeHostname(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDecryptWithFallback(t *testing.T) {
	// Encrypt with current key (normalized hostname)
	plaintext := "pk_test_token_12345"
	encrypted, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Should decrypt normally
	decrypted, err := DecryptWithFallback(encrypted)
	if err != nil {
		t.Fatalf("DecryptWithFallback failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWithFallback_OldKey(t *testing.T) {
	// Encrypt with a legacy key (simulating old hostname with .local)
	plaintext := "pk_legacy_token_67890"
	encrypted, err := encryptWithKey(plaintext, getEncryptionKeyWithHostname("test-host.local"))
	if err != nil {
		t.Fatalf("encryptWithKey failed: %v", err)
	}

	// Register the old hostname as a fallback
	oldFallbackHostnames := fallbackHostnames
	fallbackHostnames = func() []string {
		return []string{"test-host.local"}
	}
	defer func() { fallbackHostnames = oldFallbackHostnames }()

	decrypted, err := DecryptWithFallback(encrypted)
	if err != nil {
		t.Fatalf("DecryptWithFallback should succeed with fallback key: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := "pk_12345678_abcdefghijklmnop"

	encrypted, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Encrypted should be different from plaintext
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
	// Try to decrypt invalid base64
	_, err := Decrypt("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}

	// Try to decrypt valid base64 but too short for AES
	_, err = Decrypt("c2hvcnQ=") // "short" in base64
	if err == nil {
		t.Error("expected error for data too short")
	}
}

func TestIsEncrypted(t *testing.T) {
	// Plain text token
	if IsEncrypted("pk_12345678_abcdefghijklmnop") {
		t.Error("plain text should not be detected as encrypted")
	}

	// Encrypted token
	encrypted, _ := Encrypt("pk_12345678_abcdefghijklmnop")
	if !IsEncrypted(encrypted) {
		t.Error("encrypted text should be detected as encrypted")
	}
}
