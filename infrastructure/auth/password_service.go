package auth

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPassword compares hashed password with plain password
func CheckPassword(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}
