package bootstrap

import (
	"context"
	"time"

	"os"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/util/hash"
)

func SeedSuperAdmin(userRepo interfaces.IUserRepository) error {
	ctx := context.Background()

	// check if superadmin exists
	exists, err := userRepo.CheckByRole(ctx, "superadmin")
	if err != nil {
		return err
	}

	if *exists {
		return nil
	}

	// credentials from env/config, not hardcoded
	username := os.Getenv("SUPERADMIN_USERNAME")
	email := os.Getenv("SUPERADMIN_EMAIL")
	password := os.Getenv("SUPERADMIN_PASSWORD")
	passwordHash, err := hash.HashPassword(password)
	if err != nil {
		return err
	}

	// call the normal CreateUser use case
	err = userRepo.CreateUserWithStatus(ctx, &entities.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "superadmin",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, &entities.UserStatus{
		IsActive:      true,
		IsProfileFull: false,
		IsVerified:    false,
	})
	return err
}
