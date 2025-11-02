package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomID generates a secure random ID of n bytes (e.g., 16 = 32-char hex)
func GenerateRandomID(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateVerificationToken creates a random 32-character hex token
func GenerateVerificationToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
