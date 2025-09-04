package interfaces

import (
	"context"

	"remedymate-backend/domain/entities"
)

// ConversationRepository defines the interface for conversation data persistence
type ConversationRepository interface {
	// CreateConversation creates a new conversation
	CreateConversation(ctx context.Context, conversation *entities.Conversation) error

	// GetConversation retrieves a conversation by ID
	GetConversation(ctx context.Context, conversationID string) (*entities.Conversation, error)

	// UpdateConversation updates an existing conversation
	UpdateConversation(ctx context.Context, conversation *entities.Conversation) error

	// AddAnswer adds an answer to a conversation
	AddAnswer(ctx context.Context, conversationID string, answer entities.Answer) error

	// UpdateConversationStatus updates the status of a conversation
	UpdateConversationStatus(ctx context.Context, conversationID string, status entities.ConversationStatus) error

	// SetFinalReport sets the final health report for a conversation
	SetFinalReport(ctx context.Context, conversationID string, report *entities.HealthReport) error

	// DeleteExpiredConversations deletes conversations that have expired
	DeleteExpiredConversations(ctx context.Context, maxAgeHours int) error

	// GetOfflineHealthTopics retrieves all active health topics
	GetOfflineHealthTopics(ctx context.Context) ([]entities.HealthTopic, error)
}
