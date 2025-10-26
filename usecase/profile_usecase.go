package usecase

import (
	"context"
	"encoding/json"
	"time"

	"remedymate-backend/domain/entities"
	"remedymate-backend/infrastructure/cache"
	"remedymate-backend/repository"

	"github.com/redis/go-redis/v9"
)

// ProfileUsecase handles business logic for user profiles and onboarding
type ProfileUsecase struct {
	userProfileRepo *repository.UserProfileRepository
	consentRepo     *repository.ConsentRepository
	redisClient     *redis.Client
}

// NewProfileUsecase creates a new Profile usecase instance
func NewProfileUsecase(userProfileRepo *repository.UserProfileRepository, consentRepo *repository.ConsentRepository, redisClient *redis.Client) *ProfileUsecase {
	return &ProfileUsecase{
		userProfileRepo: userProfileRepo,
		consentRepo:     consentRepo,
		redisClient:     redisClient,
	}
}

// GetOnboardingStatus computes and returns the onboarding status for a user
// Checks Redis cache first (5-minute TTL), computes from database if cache miss
// Calculates completion percentage based on completed profile sections
// Determines next required step in onboarding sequence
func (uc *ProfileUsecase) GetOnboardingStatus(ctx context.Context, userID string) (*entities.OnboardingStatus, error) {
	// Check cache first
	cacheKey := cache.OnboardingStatusCacheKey(userID)
	cachedData, err := cache.Get(ctx, uc.redisClient, cacheKey)
	if err == nil && cachedData != "" {
		// Cache hit - unmarshal and return
		var status entities.OnboardingStatus
		if err := json.Unmarshal([]byte(cachedData), &status); err == nil {
			return &status, nil
		}
		// If unmarshal fails, continue with computation
	}

	// Cache miss or error - compute status
	appUser, err := uc.userProfileRepo.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	status := &entities.OnboardingStatus{
		Sections: make(map[string]entities.SectionStatus),
	}

	// Define the onboarding steps in order
	steps := []string{"consent", "demographics", "medical_history", "lifestyle"}
	stepCompleted := map[string]bool{
		"consent":         appUser.Status.ConsentCompleted,
		"demographics":    appUser.Status.DemographicsCompleted,
		"medical_history": appUser.Status.MedicalHistoryCompleted,
		"lifestyle":       appUser.Status.LifestyleCompleted,
	}

	completedCount := 0
	totalSteps := len(steps)

	for _, step := range steps {
		completed := stepCompleted[step]
		if completed {
			completedCount++
		}
		status.Sections[step] = entities.SectionStatus{
			Complete:    completed,
			LastUpdated: appUser.Status.LastUpdated,
		}
	}

	status.PercentComplete = (completedCount * 100) / totalSteps

	// Determine next step
	for _, step := range steps {
		if !stepCompleted[step] {
			status.NextStep = step
			break
		}
	}

	// If all completed, next_step can be empty or "completed"
	if status.PercentComplete == 100 {
		status.NextStep = ""
	}

	// Cache the result with 5 minute TTL
	if jsonData, err := json.Marshal(status); err == nil {
		cache.SetWithTTL(ctx, uc.redisClient, cacheKey, string(jsonData), 5*time.Minute)
	}

	return status, nil
}

// PostConsent records new consents for a user
func (uc *ProfileUsecase) PostConsent(ctx context.Context, userID string, envelope *entities.ConsentEnvelope) error {
	// Create consent record for audit trail
	err := uc.consentRepo.CreateConsent(ctx, userID, envelope)
	if err != nil {
		return err
	}

	// Update current consent in profile
	err = uc.userProfileRepo.UpsertConsent(ctx, userID, envelope)
	if err != nil {
		return err
	}

	// Invalidate onboarding status cache since consent completion status may have changed
	cache.InvalidateOnboardingStatusCache(ctx, uc.redisClient, userID)

	return nil
}

// PatchConsent appends revocation records for consents
func (uc *ProfileUsecase) PatchConsent(ctx context.Context, userID string, consentType, reason string) error {
	return uc.consentRepo.AppendRevocation(ctx, userID, consentType, reason)
}

// PutDemographics updates user demographics
func (uc *ProfileUsecase) PutDemographics(ctx context.Context, userID string, demographics *entities.Demographics) error {
	err := uc.userProfileRepo.UpsertDemographics(ctx, userID, demographics)
	if err != nil {
		return err
	}

	// Invalidate onboarding status cache since demographics completion status may have changed
	cache.InvalidateOnboardingStatusCache(ctx, uc.redisClient, userID)

	return nil
}

// PutMedicalHistory updates user medical history
func (uc *ProfileUsecase) PutMedicalHistory(ctx context.Context, userID string, history *entities.MedicalHistory) error {
	err := uc.userProfileRepo.UpsertMedicalHistory(ctx, userID, history)
	if err != nil {
		return err
	}

	// Invalidate onboarding status cache since medical history completion status may have changed
	cache.InvalidateOnboardingStatusCache(ctx, uc.redisClient, userID)

	return nil
}

// PutLifestyle updates user lifestyle data
func (uc *ProfileUsecase) PutLifestyle(ctx context.Context, userID string, lifestyle *entities.Lifestyle) error {
	err := uc.userProfileRepo.UpsertLifestyle(ctx, userID, lifestyle)
	if err != nil {
		return err
	}

	// Invalidate onboarding status cache since lifestyle completion status may have changed
	cache.InvalidateOnboardingStatusCache(ctx, uc.redisClient, userID)

	return nil
}
