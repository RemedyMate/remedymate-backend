package interfaces

import (
	"context"

	"remedymate-backend/domain/dto"
)

type IUserUsecase interface {

	// TODO: Profile management methods
	GetProfile(ctx context.Context, userID string) (*dto.ProfileResponseDTO, error)
	UpdateProfile(ctx context.Context, userID string, updateData dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error)
	DeleteProfile(ctx context.Context, userID string, deleteData dto.DeleteProfileDTO) error

	// Superadmin methods
	GetUserProfilesWithPagination(ctx context.Context, params dto.UserProfilesQueryParams) (*dto.PaginatedResponse, error)
}
