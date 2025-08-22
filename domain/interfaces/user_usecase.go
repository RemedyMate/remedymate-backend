package interfaces

import (
	"context"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
)

type IUserUsecase interface {
    RegisterUser(ctx context.Context, user entities.User) error
    GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}
