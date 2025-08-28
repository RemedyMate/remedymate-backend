package interfaces

import (
	"context"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

// TriageService defines the interface for symptom triage
type TriageService interface {
	ClassifySymptoms(ctx context.Context, input entities.SymptomInput) (*entities.TriageResult, error)
	ValidateInput(input entities.SymptomInput) error
}

// RemedyMateUsecase defines the main use case interface
type RemedyMateUsecase interface {
	GetTriage(ctx context.Context, req dto.TriageRequest) (*dto.TriageResponse, error)
	GetContent(ctx context.Context, topicKey, language string) (*entities.ContentTranslation, error)
	ComposeGuidance(ctx context.Context, req dto.ComposeRequest) (*dto.ComposeResponse, error)
}

// ContentService defines the interface for content management
type ContentService interface {
	GetApprovedBlocks() ([]entities.ApprovedBlock, error)
	GetContentByTopic(topicKey, language string) (*entities.ContentTranslation, error)
}

// GuidanceComposerService defines the interface for composing guidance cards
type GuidanceComposerService interface {
	ComposeGuidance(ctx context.Context, topicKey, language string) (*entities.GuidanceCard, error)
	ComposeFromBlocks(ctx context.Context, topicKey, language string, blocks entities.ContentTranslation) (*entities.GuidanceCard, error)
}

type RemedyAIRepository interface {
	MapSymptomToTopic(ctx context.Context, userInput string, availableTopics []string) (string, error)
	BuildMapTopicPrompt(userInput string, availableTopics []string) string
	CreatePayload(prompt string) map[string]any
	ExecuteAPIRequest(ctx context.Context, body []byte) ([]byte, error)
	ExtractTopicKeyResponse(respBody []byte) (string, error)
}
