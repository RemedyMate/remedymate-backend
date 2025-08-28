package interfaces

import (
	"context"

	"remedymate-backend/domain/dto"
)

// IAuthUsecase defines the contract for authentication business logic
type IAuthUsecase interface {
	// Login authenticates a user with email and password
	Login(ctx context.Context, loginData dto.LoginDTO) (*dto.LoginResponseDTO, error)

	// Logout invalidates a user's session
	Logout(ctx context.Context, userID string) error

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}
