package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"strings"
)

const (
	// EncryptedPrefix is added to encrypted tokens to identify them
	EncryptedPrefix = "enc:"
)

// normalizeHostname strips the .local suffix that macOS appends inconsistently
func normalizeHostname(hostname string) string {
	return strings.TrimSuffix(hostname, ".local")
}

// getEncryptionKeyWithHostname derives a key using a specific hostname
func getEncryptionKeyWithHostname(hostname string) []byte {
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // Windows
	}

	keyMaterial := hostname + ":" + username + ":clickup-cli-secret"
	hash := sha256.Sum256([]byte(keyMaterial))
	return hash[:]
}

// getEncryptionKey derives a key from machine-specific data
func getEncryptionKey() []byte {
	hostname, _ := os.Hostname()
	hostname = normalizeHostname(hostname)
	return getEncryptionKeyWithHostname(hostname)
}

// fallbackHostnames returns hostnames to try when decryption with the primary key fails.
// This handles tokens encrypted before hostname normalization was added.
var fallbackHostnames = func() []string {
	hostname, _ := os.Hostname()
	normalized := normalizeHostname(hostname)

	var fallbacks []string
	// If hostname was normalized, try the original (with .local)
	if normalized != hostname {
		fallbacks = append(fallbacks, hostname)
	} else {
		// If hostname has no .local, try with .local appended
		fallbacks = append(fallbacks, hostname+".local")
	}
	return fallbacks
}

// encryptWithKey encrypts plaintext using a specific key
func encryptWithKey(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return EncryptedPrefix + encoded, nil
}

// decryptWithKey decrypts ciphertext using a specific key
func decryptWithKey(ciphertext string, key []byte) (string, error) {
	ciphertext = strings.TrimPrefix(ciphertext, EncryptedPrefix)

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Encrypt encrypts plaintext using AES-GCM with the current machine key
func Encrypt(plaintext string) (string, error) {
	return encryptWithKey(plaintext, getEncryptionKey())
}

// Decrypt decrypts ciphertext using AES-GCM with the current machine key
func Decrypt(ciphertext string) (string, error) {
	return decryptWithKey(ciphertext, getEncryptionKey())
}

// DecryptWithFallback tries the current key first, then falls back to legacy
// hostname variants (e.g., with/without .local suffix on macOS)
func DecryptWithFallback(ciphertext string) (string, error) {
	// Try current (normalized) key first
	plaintext, err := Decrypt(ciphertext)
	if err == nil {
		return plaintext, nil
	}

	// Try fallback hostnames (handles tokens encrypted before normalization)
	for _, hostname := range fallbackHostnames() {
		key := getEncryptionKeyWithHostname(hostname)
		plaintext, fallbackErr := decryptWithKey(ciphertext, key)
		if fallbackErr == nil {
			return plaintext, nil
		}
	}

	return "", err
}

// IsEncrypted checks if a token is encrypted (has the enc: prefix)
func IsEncrypted(token string) bool {
	return strings.HasPrefix(token, EncryptedPrefix)
}
