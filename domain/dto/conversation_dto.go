package dto

import (
	"remedymate-backend/domain/entities"
)

// StartConversationRequest represents the request to start a new conversation
type StartConversationRequest struct {
	Symptom  string `json:"symptom" binding:"required" validate:"min=3,max=500"`
	Language string `json:"language" binding:"required" validate:"oneof=en am"`
	UserID   string `json:"user_id,omitempty"` // Optional, for unauthenticated users
}

// StartConversationResponse represents the response when starting a conversation
type StartConversationResponse struct {
	ConversationID string            `json:"conversation_id"`
	Question       entities.Question `json:"question"`
	TotalSteps     int               `json:"total_steps"`
	CurrentStep    int               `json:"current_step"`
}

// SubmitAnswerRequest represents the request to submit an answer
type SubmitAnswerRequest struct {
	ConversationID string `json:"conversation_id" binding:"required"`
	Answer         string `json:"answer" binding:"required" validate:"min=1,max=1000"`
}

// SubmitAnswerResponse represents the response when submitting an answer
type SubmitAnswerResponse struct {
	ConversationID string             `json:"conversation_id"`
	Question       *entities.Question `json:"question,omitempty"` // Next question if available
	Message        string             `json:"message,omitempty"`  // Feedback message for invalid answers
	IsComplete     bool               `json:"is_complete"`        // Whether all questions are answered
	CurrentStep    int                `json:"current_step"`
	TotalSteps     int                `json:"total_steps"`
}

// GetReportResponse represents the response for getting the final health report
type GetReportResponse struct {
	ConversationID string                 `json:"conversation_id"`
	Report         *entities.HealthReport `json:"report"`
	Symptom        string                 `json:"symptom"`
	Status         string                 `json:"status"`
}

// ConversationError represents conversation-specific errors
type ConversationError struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}
