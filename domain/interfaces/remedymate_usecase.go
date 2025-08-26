package interfaces

import (
	"context"

	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/domain/entities"
)

// TriageService defines the interface for symptom triage
type TriageService interface {
	ClassifySymptoms(ctx context.Context, input entities.SymptomInput) (*entities.TriageResult, error)
	ValidateInput(input entities.SymptomInput) error
}

// RemedyMateUsecase defines the main use case interface
type RemedyMateUsecase interface {
	GetTriage(ctx context.Context, req dto.TriageRequest) (*dto.TriageResponse, error)
}

// ContentService defines the interface for content management
type ContentService interface {
	GetApprovedBlocks() ([]entities.ApprovedBlock, error)
	GetContentByTopic(topicKey, language string) (*entities.ContentTranslation, error)
}
