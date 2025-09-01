package user

import (
	"context"
	"log"
	"time"

	"remedymate-backend/domain/AppError"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/auth"
	"remedymate-backend/util/hash"
)

type UserUsecase struct {
	UserRepo interfaces.IUserRepository
	AESkey   []byte
}

func NewUserUsecase(repo interfaces.IUserRepository) interfaces.IUserUsecase {
	// Use the helper function from encryption service
	key, err := auth.GetEncryptionKey()
	if err != nil {
		panic("Failed to get encryption key: " + err.Error())
	}

	return &UserUsecase{
		UserRepo: repo,
		AESkey:   key,
	}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, user entities.User) error {
	// Check if email exists
	existing, _ := u.UserRepo.FindByEmail(ctx, user.Email)
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

	err = u.UserRepo.CreateUserWithStatus(ctx, &user, userStatus)
	if err != nil {
		return err
	}

	return nil
}

// GetProfile retrieves user profile by ID
func (u *UserUsecase) GetProfile(ctx context.Context, userID string) (*dto.ProfileResponseDTO, error) {
	log.Printf("👤 Getting profile for user: %s", userID)

	user, err := u.UserRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("error finding user: %v", err)
		return nil, AppError.ErrUserNotFound
	}
	// TODO: make FindByID and GetUserStatus in one transaction
	userStatus, err := u.UserRepo.GetUserStatus(ctx, userID)
	if err != nil {
		log.Printf("error finding user status: %v", err)
		return nil, AppError.ErrUserStatusNotFound
	}

	profile := &dto.ProfileResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		PersonalInfo: dto.PersonalInfoDTO{
			FirstName: *user.PersonalInfo.FirstName,
			LastName:  *user.PersonalInfo.LastName,
			Age:       *user.PersonalInfo.Age,
			Gender:    *user.PersonalInfo.Gender,
		},
		IsVerified:    userStatus.IsVerified,
		IsProfileFull: userStatus.IsProfileFull,
		IsActive:      userStatus.IsActive,
		CreatedAt:     user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     user.UpdatedAt.Format(time.RFC3339),
		LastLogin:     user.LastLogin.Format(time.RFC3339),
	}

	log.Printf("✅ Profile retrieved for user: %s", user.Username)
	return profile, nil
}

// // UpdateProfile updates user profile (basic update)
// func (u *UserUsecase) UpdateProfile(ctx context.Context, userID string, updateData dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error) {
// 	log.Printf("🔄 Updating profile for user: %s", userID)

// 	user, err := u.UserRepo.FindByID(ctx, userID)
// 	if err != nil {
// 		return nil, errors.New("user not found")
// 	}

// 	// Update fields if provided
// 	if updateData.Username != "" {
// 		user.Username = updateData.Username
// 	}

// 	if updateData.PersonalInfo.FirstName != "" {
// 		user.PersonalInfo.FirstName = updateData.PersonalInfo.FirstName
// 	}
// 	if updateData.PersonalInfo.LastName != "" {
// 		user.PersonalInfo.LastName = updateData.PersonalInfo.LastName
// 	}
// 	if updateData.PersonalInfo.Age > 0 {
// 		user.PersonalInfo.Age = updateData.PersonalInfo.Age
// 	}
// 	if updateData.PersonalInfo.Gender != "" {
// 		user.PersonalInfo.Gender = updateData.PersonalInfo.Gender
// 	}

// 	user.UpdatedAt = time.Now()

// 	if err := u.UserRepo.UpdateUser(ctx, *user); err != nil {
// 		log.Printf("❌ Failed to update user: %v", err)
// 		return nil, errors.New("failed to update profile")
// 	}

// 	log.Printf("✅ Profile updated for user: %s", user.Username)
// 	return u.GetProfile(ctx, userID)
// }

// // EditProfile edits user profile (comprehensive edit)
// func (u *UserUsecase) EditProfile(ctx context.Context, userID string, editData dto.EditProfileDTO) (*dto.ProfileResponseDTO, error) {
// 	log.Printf("✏️ Editing profile for user: %s", userID)

// 	user, err := u.UserRepo.FindByID(ctx, userID)
// 	if err != nil {
// 		return nil, errors.New("user not found")
// 	}

// 	// Update fields if provided (similar to update but more comprehensive)
// 	if editData.Username != "" {
// 		user.Username = editData.Username
// 	}

// 	if editData.PersonalInfo.FirstName != "" {
// 		user.PersonalInfo.FirstName = editData.PersonalInfo.FirstName
// 	}
// 	if editData.PersonalInfo.LastName != "" {
// 		user.PersonalInfo.LastName = editData.PersonalInfo.LastName
// 	}
// 	if editData.PersonalInfo.Age > 0 {
// 		user.PersonalInfo.Age = editData.PersonalInfo.Age
// 	}
// 	if editData.PersonalInfo.Gender != "" {
// 		user.PersonalInfo.Gender = editData.PersonalInfo.Gender
// 	}

// 	// Handle health conditions
// 	if editData.HealthConditions != "" {
// 		encrypted, err := auth.Encrypt(editData.HealthConditions, u.AESkey)
// 		if err != nil {
// 			return nil, errors.New("failed to encrypt health conditions")
// 		}
// 		user.HealthConditions = encrypted
// 	}

// 	// Handle profile completeness flag
// 	if editData.IsProfileFull != nil {
// 		user.IsProfileFull = *editData.IsProfileFull
// 	}

// 	user.UpdatedAt = time.Now()

// 	if err := u.UserRepo.UpdateUser(ctx, *user); err != nil {
// 		log.Printf("❌ Failed to edit user: %v", err)
// 		return nil, errors.New("failed to edit profile")
// 	}

// 	log.Printf("✅ Profile edited for user: %s", user.Username)
// 	return u.GetProfile(ctx, userID)
// }

// // DeleteProfile soft deletes user profile
// func (u *UserUsecase) DeleteProfile(ctx context.Context, userID string, deleteData dto.DeleteProfileDTO) error {
// 	log.Printf("🗑️ Deleting profile for user: %s", userID)

// 	user, err := u.UserRepo.FindByID(ctx, userID)
// 	if err != nil {
// 		return errors.New("user not found")
// 	}

// 	// Verify password for security
// 	isValid, err := u.passwordService.VerifyPassword(deleteData.Password, user.Password)
// 	if err != nil {
// 		return errors.New("authentication failed")
// 	}

// 	if !isValid {
// 		log.Printf("❌ Invalid password for profile deletion: %s", userID)
// 		return errors.New("invalid password")
// 	}

// 	// Log deletion reason if provided
// 	if deleteData.Reason != "" {
// 		log.Printf("📝 Deletion reason for user %s: %s", userID, deleteData.Reason)
// 	}

// 	if err := u.UserRepo.SoftDeleteUser(ctx, userID); err != nil {
// 		log.Printf("❌ Failed to delete user: %v", err)
// 		return errors.New("failed to delete profile")
// 	}

// 	log.Printf("✅ Profile deleted for user: %s", userID)
// 	return nil
// }
