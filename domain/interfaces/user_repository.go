package interfaces

import (
	"context"

	"remedymate-backend/domain/entities"
)

type IUserRepository interface {
	CreateUserWithStatus(ctx context.Context, user *entities.User, userStatus *entities.UserStatus) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	CheckByRole(ctx context.Context, role string) (*bool, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, userID string) (*entities.User, error)

	SoftDeleteUser(ctx context.Context, userID string) error
	// user status
	GetUserStatus(ctx context.Context, userID string) (*entities.UserStatus, error)
	UpdateUserStatusFields(ctx context.Context, userID string, fields map[string]interface{}) error
	// CreateUserStatus(ctx context.Context, userStatus *entities.UserStatus) error

	// Activate user by email (sets status.isActive=true)
	ActivateByEmail(ctx context.Context, email string) error
}

type IRefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, refreshToken *entities.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, tokenId string) error
}

type IActivationTokenRepository interface {
	Create(ctx context.Context, token *entities.ActivationToken) error
	FindValidByToken(ctx context.Context, token string) (*entities.ActivationToken, error)
	MarkUsed(ctx context.Context, id string) error
}
