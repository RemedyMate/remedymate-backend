package interfaces

import (
	"context"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

type IUserUsecase interface {
	RegisterUser(ctx context.Context, user entities.User) error
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)

	// Profile management methods
	GetProfile(ctx context.Context, userID string) (*dto.ProfileResponseDTO, error)
	UpdateProfile(ctx context.Context, userID string, updateData dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error)
	EditProfile(ctx context.Context, userID string, editData dto.EditProfileDTO) (*dto.ProfileResponseDTO, error)
	DeleteProfile(ctx context.Context, userID string, deleteData dto.DeleteProfileDTO) error
}
