package interfaces

import (
	"context"

	"remedymate-backend/domain/entities"
)

type IUserRepository interface {
	CreateUserWithStatus(ctx context.Context, user *entities.User, userStatus *entities.UserStatus) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	CheckByRole(ctx context.Context, role string) (*bool, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, userID string) (*entities.User, error)
	SoftDeleteUser(ctx context.Context, userID string) error
	// user status
	GetUserStatus(ctx context.Context, userID string) (*entities.UserStatus, error)
	// CreateUserStatus(ctx context.Context, userStatus *entities.UserStatus) error

}

type IRefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, refreshToken *entities.RefreshToken) error
}
