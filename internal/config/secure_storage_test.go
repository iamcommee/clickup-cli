package config

import (
	"testing"
)

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
