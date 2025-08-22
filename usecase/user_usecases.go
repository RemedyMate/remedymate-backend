package usecase

import (
	"context"
	"errors"
	"os"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
	"github.com/RemedyMate/remedymate-backend/infrastructure/auth"
)

type UserUsecase struct {
	UserRepo interfaces.IUserRepository
	AESkey []byte
}

func NewUserUsecase(repo interfaces.IUserRepository) interfaces.IUserUsecase{
	key := []byte(os.Getenv("AES_KEY")) 
	return &UserUsecase{UserRepo: repo, AESkey: key}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, user entities.User) error {
    // Check if email exists
    existing, _ := u.UserRepo.FindByEmail(ctx, user.Email)
    if existing != nil {
        return errors.New("email already exists")
    }

	 // Hash password
    hashed, err := auth.HashPassword(user.Password)
    if err != nil {
        return err
    }
    user.Password = hashed

    // Encrypt healthConditions
    if user.HealthConditions != "" {
        encrypted, err := auth.Encrypt(user.HealthConditions, u.AESkey)
        if err != nil {
			 return err
			 }
        user.HealthConditions = encrypted
    }

    return u.UserRepo.InsertUser(ctx, user)
}

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
    return u.UserRepo.FindByEmail(ctx, email)
}