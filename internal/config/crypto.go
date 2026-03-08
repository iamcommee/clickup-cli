package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// EncryptedPrefix is added to encrypted tokens to identify them
	EncryptedPrefix = "enc:"
	// keyFileName is the name of the file that stores the encryption key
	keyFileName = ".key"
	// keySize is the size of the encryption key in bytes (AES-256)
	keySize = 32
)

// getOrCreateKeyFile returns a 32-byte encryption key from a .key file in dir.
// If the file does not exist, a new random key is generated and written.
func getOrCreateKeyFile(dir string) ([]byte, error) {
	keyPath := filepath.Join(dir, keyFileName)

	data, err := os.ReadFile(keyPath)
	if err == nil && len(data) == keySize {
		return data, nil
	}

	// Generate a new random key
	key := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	if err := os.WriteFile(keyPath, key, 0600); err != nil {
		return nil, err
	}

	return key, nil
}

// getEncryptionKey returns the encryption key from the key file
func getEncryptionKey() ([]byte, error) {
	dir := DefaultConfigDir()
	if dir == "" {
		return nil, errors.New("cannot determine config directory")
	}
	return getOrCreateKeyFile(dir)
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

// Encrypt encrypts plaintext using AES-GCM with the key file
func Encrypt(plaintext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}
	return encryptWithKey(plaintext, key)
}

// Decrypt decrypts ciphertext using AES-GCM with the key file
func Decrypt(ciphertext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}
	return decryptWithKey(ciphertext, key)
}

// IsEncrypted checks if a token is encrypted (has the enc: prefix)
func IsEncrypted(token string) bool {
	return strings.HasPrefix(token, EncryptedPrefix)
}
