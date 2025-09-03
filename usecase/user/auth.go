package user

import (
	"context"

	// "errors"
	"log"
	"time"

	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/util/hash"
	jwtutil "remedymate-backend/util/jwt"
)

// AuthUsecase implements IAuthUsecase interface
type AuthUsecase struct {
	userRepo  interfaces.IUserRepository
	tokenRepo interfaces.IRefreshTokenRepository
}

// NewAuthUsecase creates a new Auth usecase instance
func NewAuthUsecase(userRepo interfaces.IUserRepository, tokenRepo interfaces.IRefreshTokenRepository) interfaces.IAuthUsecase {
	return &AuthUsecase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, user *entities.User) error {
	// Check if email exists
	existing, _ := uc.userRepo.FindByEmail(ctx, user.Email)
	if existing != nil {
		return AppError.ErrEmailAlreadyExist
	}

	// Hash password using the password service
	hashed, err := hash.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.PasswordHash = hashed

	// Initialize user status
	userStatus := &entities.UserStatus{
		IsActive:      false,
		IsProfileFull: false,
		IsVerified:    true,
	}

	err = uc.userRepo.CreateUserWithStatus(ctx, user, userStatus)
	if err != nil {
		return err
	}

	return nil
}

// Login authenticates a user with email and password
func (uc *AuthUsecase) Login(ctx context.Context, loginData dto.LoginDTO) (*dto.LoginResponseDTO, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, loginData.Email)
	if err != nil {
		return nil, err
	}

	userStatus, err := uc.userRepo.GetUserStatus(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Check if user is active
	if !userStatus.IsActive {
		log.Printf("❌ User account is inactive: %s", user.Email)
		return nil, AppError.ErrInactiveAccount
	}

	// Verify password
	isValid, err := hash.VerifyPassword(loginData.Password, user.PasswordHash)
	if err != nil {
		log.Printf("error verifying password for user %s: %v", user.Email, err)
		return nil, AppError.ErrInternalServer
	}

	if !isValid {
		log.Printf("invalid password for user: %s", user.Email)
		return nil, AppError.ErrIncorrectPassword
	}

	// Update last login
	user.LastLogin = time.Now()
	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		log.Printf("failed to update last login: %v", err)
		return nil, err
	}

	// Generate JWT access token
	accessTokenString, err := jwtutil.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwtutil.GenerateRefreshToken(user)
	if err != nil {
		log.Printf("failed to produce refresh token: %v", err)
		return nil, err
	}

	// Persist refresh token
	// tokenHash := hashutil.HashToken(refreshToken, us.cfg.HMAC.Secret) Hashing is a good idea

	err = uc.tokenRepo.StoreRefreshToken(ctx, refreshToken)

	if err != nil {
		log.Println("failed to store refresh token: ", err.Error())
		return nil, err
	}

	log.Printf("✅ Login successful for user: %s (%s)", user.Username, user.Email)
	return &dto.LoginResponseDTO{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (uc *AuthUsecase) RefreshToken(ctx context.Context, tokenString string) (*dto.RefreshResponseDTO, error) {
	claims, err := jwtutil.ValidateToken(tokenString, false)
	if err != nil {
		return nil, err
	}

	err = uc.tokenRepo.DeleteRefreshToken(ctx, claims.TokenID)
	if err != nil {
		return nil, err
	}

	newAccessTokenString, err := jwtutil.GenerateAccessToken(&entities.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role})
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwtutil.GenerateRefreshToken(&entities.User{
		ID:       claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		Role:     claims.Role})
	if err != nil {
		return nil, err
	}

	err = uc.tokenRepo.StoreRefreshToken(ctx, newRefreshToken)
	if err != nil {
		return nil, err
	}
	return &dto.RefreshResponseDTO{
		AccessToken:  newAccessTokenString,
		RefreshToken: newRefreshToken.Token,
	}, nil
}

// Logout invalidates a user's session
func (uc *AuthUsecase) Logout(ctx context.Context, userID string) error {
	// TODO: Implement the logout
	return nil
}

// ChangePassword changes a user's password
// func (uc *AuthUsecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
// 	log.Printf("🔐 Password change requested for user: %s", userID)

// 	// Find user by ID
// 	user, err := uc.userRepo.FindByID(ctx, userID)
// 	if err != nil {
// 		log.Printf("❌ Error finding user by ID: %v", err)
// 		return errors.New("user not found")
// 	}

// 	if user == nil {
// 		log.Printf("❌ User not found for ID: %s", userID)
// 		return errors.New("user not found")
// 	}

// 	// Verify old password
// 	isValid, err := uc.passwordService.VerifyPassword(oldPassword, user.Password)
// 	if err != nil {
// 		log.Printf("❌ Error verifying old password: %v", err)
// 		return errors.New("authentication failed")
// 	}

// 	if !isValid {
// 		log.Printf("❌ Invalid password for user: %s", userID)
// 		return errors.New("invalid old password")
// 	}

// 	// Hash new password
// 	hashedPassword, err := uc.passwordService.HashPassword(newPassword)
// 	if err != nil {
// 		log.Printf("❌ Error hashing new password: %v", err)
// 		return errors.New("failed to process new password")
// 	}

// 	// Update user password
// 	user.Password = hashedPassword
// 	user.UpdatedAt = time.Now()

// 	if err := uc.userRepo.UpdateUser(ctx, *user); err != nil {
// 		log.Printf("❌ Failed to update password: %v", err)
// 		return errors.New("failed to update password")
// 	}

// 	log.Printf("✅ Password changed successfully for user: %s", userID)
// 	return nil
// }
