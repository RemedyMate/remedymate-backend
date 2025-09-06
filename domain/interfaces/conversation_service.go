package interfaces

import (
	"context"

	"remedymate-backend/domain/entities"
)

// ConversationService defines the interface for conversation management
type ConversationService interface {
	// ValidateSymptom validates if the provided symptom is medical and appropriate
	ValidateSymptom(ctx context.Context, symptom, language string) (bool, string, error)

	// GenerateQuestions generates follow-up questions based on the initial symptom
	GenerateQuestions(ctx context.Context, symptom, language string) ([]entities.Question, error)

	// ValidateAnswer validates a user's answer to a question
	ValidateAnswer(ctx context.Context, question entities.Question, answer string) (bool, string, error)

	// GenerateHealthReport creates a structured health report from conversation data
	GenerateHealthReport(ctx context.Context, conversation *entities.Conversation) (*entities.HealthReport, error)
}
