package hash

import (
	AppError "remedymate-backend/domain/AppError"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Generate hash with cost factor 12 (good balance of security vs performance)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", AppError.ErrInternalServer
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a plain text password against a hashed password
func VerifyPassword(plainPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, AppError.ErrIncorrectPassword
		}
		return false, AppError.ErrInternalServer
	}
	return true, nil
}
