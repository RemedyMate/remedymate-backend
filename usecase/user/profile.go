package user

import (
	"context"
	"log"
	"math"
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

	// Safe unwraps
	fn, ln, age, gender, pfp := "", "", 0, "", ""
	if user.PersonalInfo != nil {
		if user.PersonalInfo.FirstName != nil {
			fn = *user.PersonalInfo.FirstName
		}
		if user.PersonalInfo.LastName != nil {
			ln = *user.PersonalInfo.LastName
		}
		if user.PersonalInfo.Age != nil {
			age = *user.PersonalInfo.Age
		}
		if user.PersonalInfo.Gender != nil {
			gender = *user.PersonalInfo.Gender
		}
		if user.PersonalInfo.ProfilePictureURL != nil {
			pfp = *user.PersonalInfo.ProfilePictureURL
		}
	}

	profile := &dto.ProfileResponseDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		PersonalInfo: dto.PersonalInfoDTO{
			FirstName:         fn,
			LastName:          ln,
			Age:               age,
			Gender:            gender,
			ProfilePictureURL: pfp,
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

// UpdateProfile updates user profile (basic update)
func (u *UserUsecase) UpdateProfile(ctx context.Context, userID string, updateData dto.UpdateProfileDTO) (*dto.ProfileResponseDTO, error) {
	log.Printf("🔄 Updating profile for user: %s", userID)

	user, err := u.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, AppError.ErrUserNotFound
	}

	// Update fields if provided
	if updateData.Username != "" {
		user.Username = updateData.Username
	}

	if user.PersonalInfo == nil {
		user.PersonalInfo = &entities.PersonalInfo{}
	}
	// Apply personal info deltas
	if updateData.PersonalInfo.FirstName != "" {
		user.PersonalInfo.FirstName = &updateData.PersonalInfo.FirstName
	}
	if updateData.PersonalInfo.LastName != "" {
		user.PersonalInfo.LastName = &updateData.PersonalInfo.LastName
	}
	if updateData.PersonalInfo.Age > 0 {
		age := updateData.PersonalInfo.Age
		user.PersonalInfo.Age = &age
	}
	if updateData.PersonalInfo.Gender != "" {
		user.PersonalInfo.Gender = &updateData.PersonalInfo.Gender
	}

	// Apply profile picture URL
	if updateData.PersonalInfo.ProfilePictureURL != "" {
		url := updateData.PersonalInfo.ProfilePictureURL
		user.PersonalInfo.ProfilePictureURL = &url
	}

	user.UpdatedAt = time.Now()

	if err := u.UserRepo.UpdateUser(ctx, user); err != nil {
		log.Printf("❌ Failed to update user: %v", err)
		return nil, AppError.ErrInternalServer
	}

	// Update IsProfileFull status based on completeness
	if err := u.updateIsProfileFull(ctx, userID); err != nil {
		log.Printf("warning: failed to update IsProfileFull: %v", err)
	}

	log.Printf("✅ Profile updated for user: %s", user.Username)
	return u.GetProfile(ctx, userID)
}

// DeleteProfile soft deletes user profile
func (u *UserUsecase) DeleteProfile(ctx context.Context, userID string, deleteData dto.DeleteProfileDTO) error {
	log.Printf("🗑️ Deleting profile for user: %s", userID)

	user, err := u.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return AppError.ErrUserNotFound
	}

	// Verify password for security — compare with PasswordHash
	ok, verr := hash.VerifyPassword(deleteData.Password, user.PasswordHash)
	if verr != nil {
		return verr
	}
	if !ok {
		return AppError.ErrIncorrectPassword
	}

	if err := u.UserRepo.SoftDeleteUser(ctx, userID); err != nil {
		log.Printf("❌ Failed to delete user: %v", err)
		return AppError.ErrInternalServer
	}

	log.Printf("✅ Profile deleted for user: %s", userID)
	return nil
}

// updateIsProfileFull calculates completeness and stores in user_status
func (u *UserUsecase) updateIsProfileFull(ctx context.Context, userID string) error {
	user, err := u.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	// Define completeness: non-empty FirstName, LastName, Age>0, Gender
	complete := false
	if user.PersonalInfo != nil &&
		user.PersonalInfo.FirstName != nil && *user.PersonalInfo.FirstName != "" &&
		user.PersonalInfo.LastName != nil && *user.PersonalInfo.LastName != "" &&
		user.PersonalInfo.Age != nil && *user.PersonalInfo.Age > 0 &&
		user.PersonalInfo.Gender != nil && *user.PersonalInfo.Gender != "" &&
		user.PersonalInfo.ProfilePictureURL != nil && *user.PersonalInfo.ProfilePictureURL != "" {
		complete = true
	}
	return u.UserRepo.UpdateUserStatusFields(ctx, userID, map[string]interface{}{"isProfileFull": complete})
}

// GetUserProfilesWithPagination retrieves user profiles with pagination, filtering, and sorting (superadmin only)
func (u *UserUsecase) GetUserProfilesWithPagination(ctx context.Context, params dto.UserProfilesQueryParams) (*dto.PaginatedResponse, error) {
	log.Printf("👥 Getting user profiles with pagination - page: %d, limit: %d", params.Page, params.Limit)

	// Get paginated users
	users, total, err := u.UserRepo.FindUsersWithPagination(ctx, params)
	if err != nil {
		log.Printf("error finding users with pagination: %v", err)
		return nil, AppError.ErrInternalServer
	}

	// Extract user IDs for status lookup
	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	// Get user statuses for the current page users
	var userStatuses []*entities.UserStatus
	if len(userIDs) > 0 {
		userStatuses, err = u.UserRepo.GetUserStatusesByUserIDs(ctx, userIDs)
		if err != nil {
			log.Printf("error finding user statuses: %v", err)
			return nil, AppError.ErrInternalServer
		}
	}

	// Create a map of userID to userStatus for efficient lookup
	statusMap := make(map[string]*entities.UserStatus)
	for _, status := range userStatuses {
		statusMap[status.UserID] = status
	}

	// Build profile responses
	var profiles []*dto.ProfileResponseDTO
	for _, user := range users {
		// Get the corresponding user status
		userStatus, exists := statusMap[user.ID]
		if !exists {
			log.Printf("warning: no status found for user %s, using default values", user.ID)
			userStatus = &entities.UserStatus{
				UserID:        user.ID,
				IsActive:      false,
				IsProfileFull: false,
				IsVerified:    false,
			}
		}

		// Safe unwraps for personal info
		fn, ln, age, gender, pfp := "", "", 0, "", ""
		if user.PersonalInfo != nil {
			if user.PersonalInfo.FirstName != nil {
				fn = *user.PersonalInfo.FirstName
			}
			if user.PersonalInfo.LastName != nil {
				ln = *user.PersonalInfo.LastName
			}
			if user.PersonalInfo.Age != nil {
				age = *user.PersonalInfo.Age
			}
			if user.PersonalInfo.Gender != nil {
				gender = *user.PersonalInfo.Gender
			}
			if user.PersonalInfo.ProfilePictureURL != nil {
				pfp = *user.PersonalInfo.ProfilePictureURL
			}
		}

		profile := &dto.ProfileResponseDTO{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			PersonalInfo: dto.PersonalInfoDTO{
				FirstName:         fn,
				LastName:          ln,
				Age:               age,
				Gender:            gender,
				ProfilePictureURL: pfp,
			},
			IsVerified:    userStatus.IsVerified,
			IsProfileFull: userStatus.IsProfileFull,
			IsActive:      userStatus.IsActive,
			CreatedAt:     user.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     user.UpdatedAt.Format(time.RFC3339),
			LastLogin:     user.LastLogin.Format(time.RFC3339),
		}

		profiles = append(profiles, profile)
	}

	// Calculate pagination metadata
	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 {
		limit = 10
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasNext := page < totalPages
	hasPrev := page > 1

	pagination := dto.PaginationMetadata{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	response := &dto.PaginatedResponse{
		Data:       profiles,
		Pagination: pagination,
		Message:    "User profiles retrieved successfully",
	}

	log.Printf("✅ Retrieved %d user profiles (page %d of %d)", len(profiles), page, totalPages)
	return response, nil
}
