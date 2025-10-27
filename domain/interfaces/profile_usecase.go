package interfaces

import (
	"context"

	"remedymate-backend/domain/entities"
)

type IProfileUsecase interface {
	GetOnboardingStatus(ctx context.Context, userID string) (*entities.OnboardingStatus, error)
	PostConsent(ctx context.Context, userID string, envelope *entities.ConsentEnvelope) error
	PatchConsent(ctx context.Context, userID string, consentType, reason string) error
	PutDemographics(ctx context.Context, userID string, demographics *entities.Demographics) error
	PutMedicalHistory(ctx context.Context, userID string, history *entities.MedicalHistory) error
	PutLifestyle(ctx context.Context, userID string, lifestyle *entities.Lifestyle) error
}
