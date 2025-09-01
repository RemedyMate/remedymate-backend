package bootstrap

import (
	"context"
	"os"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
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

	// call the normal CreateUser use case
	err = userRepo.CreateUserWithStatus(ctx, &entities.User{
		Username: username,
		Email:    email,
		Password: password,
		Role:     "superadmin",
	}, &entities.UserStatus{
		IsActive:      true,
		IsProfileFull: false,
		IsVerified:    false,
	})
	return err
}
