package interfaces

import (
	"context"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

// IAuthUsecase defines the contract for authentication business logic
type IAuthUsecase interface {
	// Register registers a new admin
	Register(ctx context.Context, user *entities.User, frontendDomain string) (*dto.RegisterResponseDTO, error)

	// Login authenticates a user with email and password
	Login(ctx context.Context, loginData dto.LoginDTO) (*dto.LoginResponseDTO, error)

	// Refresh gives a new refresh token
	RefreshToken(ctx context.Context, tokenString string) (*dto.RefreshResponseDTO, error)

	// Logout invalidates a user's session
	Logout(ctx context.Context, userID string) error

	// Activate activates a user account by email
	Activate(ctx context.Context, email string) error

	// VerifyAccount verifies token and activates
	VerifyAccount(ctx context.Context, token string) error

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error

	// ResendVerificationToken resends the verification email if token expired or not used
	ResendVerificationToken(ctx context.Context, email, frontendDomain string) error
}
