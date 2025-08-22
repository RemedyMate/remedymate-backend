package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password hashing and verification
type PasswordService struct{}

// NewPasswordService creates a new password service instance
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// HashPassword hashes a plain text password using bcrypt
func (ps *PasswordService) HashPassword(password string) (string, error) {
	// Generate hash with cost factor 12 (good balance of security vs performance)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a plain text password against a hashed password
func (ps *PasswordService) VerifyPassword(plainPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil // Password doesn't match
		}
		return false, err // Other error
	}
	return true, nil // Password matches
}
