package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/auth"
)

// AuthUsecase implements IAuthUsecase interface
type AuthUsecase struct {
	userRepo        interfaces.IUserRepository
	passwordService *auth.PasswordService
	jwtService      *auth.JWTService
}

// NewAuthUsecase creates a new Auth usecase instance
func NewAuthUsecase(userRepo interfaces.IUserRepository, passwordService *auth.PasswordService, jwtService *auth.JWTService) *AuthUsecase {
	return &AuthUsecase{
		userRepo:        userRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

// Login authenticates a user with email and password
func (uc *AuthUsecase) Login(ctx context.Context, loginData dto.LoginDTO) (*dto.LoginResponseDTO, error) {
	log.Printf("🔐 Login attempt for email: %s", loginData.Email)

	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, loginData.Email)
	if err != nil {
		log.Printf("❌ Error finding user by email: %v", err)
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		log.Printf("❌ User not found for email: %s", loginData.Email)
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		log.Printf("❌ User account is inactive: %s", user.Email)
		return nil, errors.New("account is inactive")
	}

	// Verify password
	log.Printf("🔍 Verifying password for user: %s", user.Email)
	isValid, err := uc.passwordService.VerifyPassword(loginData.Password, user.Password)
	if err != nil {
		log.Printf("❌ Error verifying password: %v", err)
		return nil, errors.New("authentication failed")
	}

	if !isValid {
		log.Printf("❌ Invalid password for user: %s", user.Email)
		return nil, errors.New("invalid credentials")
	}

	// Update last login
	user.LastLogin = time.Now()
	if err := uc.userRepo.UpdateUser(ctx, *user); err != nil {
		log.Printf("⚠️ Failed to update last login: %v", err)
		// Don't fail login for this, just log warning
	}

	// Generate JWT token
	log.Printf("🔑 Generating JWT token for user: %s", user.ID)
	jwtToken, err := uc.jwtService.GenerateToken(*user)
	if err != nil {
		log.Printf("❌ Failed to generate JWT token: %v", err)
		return nil, errors.New("authentication failed")
	}

	log.Printf("✅ Login successful for user: %s (%s)", user.Username, user.Email)
	return &dto.LoginResponseDTO{
		AccessToken:  jwtToken,
		RefreshToken: "", // You can implement refresh tokens later
		User:         user,
		Message:      "Login successful",
	}, nil
}

// Logout invalidates a user's session
func (uc *AuthUsecase) Logout(ctx context.Context, userID string) error {
	log.Printf("🚪 Logout requested for user: %s", userID)

	// In a more advanced implementation, you might:
	// - Add token to blacklist
	// - Update user's last logout time
	// - Clear refresh tokens

	log.Printf("✅ Logout successful for user: %s", userID)
	return nil
}

// ChangePassword changes a user's password
func (uc *AuthUsecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	log.Printf("🔐 Password change requested for user: %s", userID)

	// Find user by ID
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("❌ Error finding user by ID: %v", err)
		return errors.New("user not found")
	}

	if user == nil {
		log.Printf("❌ User not found for ID: %s", userID)
		return errors.New("user not found")
	}

	// Verify old password
	isValid, err := uc.passwordService.VerifyPassword(oldPassword, user.Password)
	if err != nil {
		log.Printf("❌ Error verifying old password: %v", err)
		return errors.New("authentication failed")
	}

	if !isValid {
		log.Printf("❌ Invalid password for user: %s", userID)
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := uc.passwordService.HashPassword(newPassword)
	if err != nil {
		log.Printf("❌ Error hashing new password: %v", err)
		return errors.New("failed to process new password")
	}

	// Update user password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.UpdateUser(ctx, *user); err != nil {
		log.Printf("❌ Failed to update password: %v", err)
		return errors.New("failed to update password")
	}

	log.Printf("✅ Password changed successfully for user: %s", userID)
	return nil
}
