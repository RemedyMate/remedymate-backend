package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"os"
)

// Encrypt encrypts sensitive data using AES
func Encrypt(plaintext string, key []byte) (string, error) {
	// Validate key size
	if len(key) == 0 {
		return "", errors.New("encryption key cannot be empty")
	}

	// AES requires key size of 16, 24, or 32 bytes
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("encryption key must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCFBEncrypter(block, key[:block.BlockSize()])
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts AES encrypted data
func Decrypt(ciphertext string, key []byte) (string, error) {
	// Validate key size
	if len(key) == 0 {
		return "", errors.New("encryption key cannot be empty")
	}

	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("encryption key must be 16, 24, or 32 bytes")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(data))
	stream := cipher.NewCFBDecrypter(block, key[:block.BlockSize()])
	stream.XORKeyStream(plaintext, data)

	return string(plaintext), nil
}

// GetEncryptionKey returns the encryption key from environment
func GetEncryptionKey() ([]byte, error) {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		return nil, errors.New("ENCRYPTION_KEY environment variable not set")
	}

	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return nil, errors.New("ENCRYPTION_KEY must be 16, 24, or 32 characters")
	}

	return keyBytes, nil
}
