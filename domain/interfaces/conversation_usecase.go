package interfaces

import (
	"context"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

// ConversationUsecase defines the interface for conversation business logic
type ConversationUsecase interface {
	// ValidateSymptom validates if the provided symptom is medical and appropriate
	ValidateSymptom(ctx context.Context, symptom, language string) (bool, string, error)

	// StartConversation starts a new conversation with the initial symptom
	StartConversation(ctx context.Context, req dto.StartConversationRequest) (*dto.StartConversationResponse, error)

	// SubmitAnswer submits an answer to the current question
	SubmitAnswer(ctx context.Context, req dto.SubmitAnswerRequest) (*dto.SubmitAnswerResponse, error)

	// GetReport retrieves the final health report for a completed conversation
	GetReport(ctx context.Context, conversationID string) (*dto.GetReportResponse, error)

	// GetOfflineHealthTopics retrieves all active health topics
	GetOfflineHealthTopics(ctx context.Context) ([]entities.HealthTopic, error)
}
