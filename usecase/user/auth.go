package user

import (
	"context"
	"log"
	"os"
	"time"

	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	mailInfra "remedymate-backend/infrastructure/mail"
	apputil "remedymate-backend/util"
	"remedymate-backend/util/hash"
	jwtutil "remedymate-backend/util/jwt"
)

// AuthUsecase implements IAuthUsecase interface
type AuthUsecase struct {
	userRepo       interfaces.IUserRepository
	tokenRepo      interfaces.IRefreshTokenRepository
	mailer         interfaces.IMailService
	activationRepo interfaces.IActivationTokenRepository
}

// NewAuthUsecase creates a new Auth usecase instance
func NewAuthUsecase(userRepo interfaces.IUserRepository, tokenRepo interfaces.IRefreshTokenRepository, mailer interfaces.IMailService, activationRepo interfaces.IActivationTokenRepository) interfaces.IAuthUsecase {
	return &AuthUsecase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		mailer:         mailer,
		activationRepo: activationRepo,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, user *entities.User) (*dto.RegisterResponseDTO, error) {
	// Check if email exists
	existing, _ := uc.userRepo.FindByEmail(ctx, user.Email)
	if existing != nil {
		return nil, AppError.ErrEmailAlreadyExist
	}

	// Generate username if empty and ensure uniqueness
	if user.Username == "" {
		// try a handful of times to avoid rare collision
		for i := 0; i < 5; i++ {
			token, _ := apputil.GenerateToken(3) // 6 hex chars
			candidate := "rm_" + token
			if _, err := uc.userRepo.FindByUsername(ctx, candidate); err != nil {
				// not found => use it
				user.Username = candidate
				break
			}
		}
		if user.Username == "" {
			user.Username = "rm_user"
		}
	} else {
		// If provided, ensure it's unique
		if _, err := uc.userRepo.FindByUsername(ctx, user.Username); err == nil {
			return nil, AppError.ErrDuplicateUsername
		}
	}

	// Generate password if empty
	generatedPassword := user.Password
	if generatedPassword == "" {
		pwToken, _ := apputil.GenerateToken(8) // 16 hex chars
		generatedPassword = pwToken
		user.Password = generatedPassword
	}

	// Hash password using the password service
	hashed, err := hash.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashed

	// Initialize user status
	userStatus := &entities.UserStatus{
		IsActive:      false,
		IsProfileFull: false,
		IsVerified:    false,
	}

	err = uc.userRepo.CreateUserWithStatus(ctx, user, userStatus)
	if err != nil {
		return nil, err
	}

	// If mailer and activationRepo are configured, send verification email
	if uc.mailer != nil && uc.activationRepo != nil {
		// generate token and save
		tokenStr, genErr := apputil.GenerateToken(32)
		if genErr == nil {
			at := &entities.ActivationToken{
				Token:     tokenStr,
				UserID:    user.ID,
				Email:     user.Email,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				CreatedAt: time.Now(),
			}
			if err := uc.activationRepo.Create(ctx, at); err != nil {
				return nil, AppError.ErrInternalServer
			}
			// send email using template
			baseURL := os.Getenv("APP_BASE_URL")
			if baseURL == "" {
				log.Printf("APP_BASE_URL environment variable is not set")
				return nil, AppError.ErrInternalServer
			}
			link := baseURL + "/api/v1/auth/verify?token=" + tokenStr
			subject := "Verify your RemedyMate account"

			// Prepare template data
			firstName := ""
			if user.PersonalInfo != nil && user.PersonalInfo.FirstName != nil {
				firstName = *user.PersonalInfo.FirstName
			}
			tplData := struct {
				AppName     string
				FirstName   string
				VerifyLink  string
				ExpiryHours int
				Year        int
			}{
				AppName:     "RemedyMate",
				FirstName:   firstName,
				VerifyLink:  link,
				ExpiryHours: 24,
				Year:        time.Now().Year(),
			}

			body, rendErr := mailInfra.RenderTemplate("./infrastructure/mail/templates/activation_email.html", tplData)
			if rendErr != nil {
				log.Printf("failed to render activation email template: %v", rendErr)
				body = "<p>Hello,</p><p>Please verify your account by clicking the link below:</p><p><a href='" + link + "'>Activate Account</a></p>"
			}

			if err := uc.mailer.Send(user.Email, subject, body); err != nil {
				return nil, AppError.ErrEmailSendFailed
			}
		}
	}

	return &dto.RegisterResponseDTO{
		Message:  "User registered successfully",
		Username: user.Username,
		Password: generatedPassword,
	}, nil
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
		UserID:       user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
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

// Activate activates a user account by email
func (uc *AuthUsecase) Activate(ctx context.Context, email string) error {
	// Ensure user exists
	_, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	return uc.userRepo.ActivateByEmail(ctx, email)
}

// VerifyAccount verifies using a token and activates the user account
func (uc *AuthUsecase) VerifyAccount(ctx context.Context, token string) error {
	if uc.activationRepo == nil {
		return AppError.ErrInternalServer
	}
	at, err := uc.activationRepo.FindValidByToken(ctx, token)
	if err != nil {
		return err
	}
	if err := uc.userRepo.ActivateByEmail(ctx, at.Email); err != nil {
		return err
	}
	_ = uc.activationRepo.MarkUsed(ctx, at.ID)
	return nil
}

// helper to read env with default
func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// ChangePassword changes a user's password
func (uc *AuthUsecase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	log.Printf("🔐 Password change requested for user: %s", userID)

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return AppError.ErrUserNotFound
	}

	// Verify old password
	isValid, verr := hash.VerifyPassword(oldPassword, user.PasswordHash)
	if verr != nil {
		return verr
	}
	if !isValid {
		return AppError.ErrIncorrectPassword
	}

	// Hash new password and save
	hashed, herr := hash.HashPassword(newPassword)
	if herr != nil {
		return herr
	}
	user.PasswordHash = hashed
	user.UpdatedAt = time.Now()
	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		log.Printf("❌ Failed to update password: %v", err)
		return err
	}
	log.Printf("✅ Password changed successfully for user: %s", userID)
	return nil
}
