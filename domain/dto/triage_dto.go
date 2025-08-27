package dto

import "github.com/RemedyMate/remedymate-backend/domain/entities"

// TriageRequest represents the request for symptom triage
type TriageRequest struct {
	Text     string `json:"text" binding:"required" validate:"min=3,max=500"`
	Language string `json:"language" binding:"required" validate:"oneof=en am"`
}

// TriageResponse represents the response from triage
type TriageResponse struct {
	Level     entities.TriageLevel `json:"level"`
	RedFlags  []string             `json:"red_flags"`
	Message   string               `json:"message"`
	SessionID string               `json:"session_id,omitempty"`
}

// ChatRequest represents a complete chat request (combines triage, mapping, and composition)
type ChatRequest struct {
	Text     string `json:"text" binding:"required" validate:"min=3,max=500"`
	Language string `json:"language" binding:"required" validate:"oneof=en am"`
}

// ChatResponse represents a complete chat response
type ChatResponse struct {
	Triage       entities.TriageResult  `json:"triage"`
	GuidanceCard *entities.GuidanceCard `json:"guidance_card,omitempty"`
	SessionID    string                 `json:"session_id"`
	IsOffline    bool                   `json:"is_offline"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// ComposeRequest represents the request for guidance composition
type ComposeRequest struct {
	TopicKey string `json:"topic_key" binding:"required"`
	Language string `json:"language" binding:"required" validate:"oneof=en am"`
}

// ComposeResponse represents the response from guidance composition
type ComposeResponse struct {
	GuidanceCard entities.GuidanceCard `json:"guidance_card"`
	SessionID    string                `json:"session_id,omitempty"`
}
