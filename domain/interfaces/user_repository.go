package interfaces

import (
	"context"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
)

type IUserRepository interface {
	InsertUser(ctx context.Context, user entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, user entities.User) error
	FindByID(ctx context.Context, userID string) (*entities.User, error)
	SoftDeleteUser(ctx context.Context, userID string) error
}
