package dto

import "remedymate-backend/domain/entities"

// ChatRequest represents a complete chat request (combines triage, mapping, and composition)
type ChatRequest struct {
	Text     string `json:"text" binding:"required" validate:"min=3,max=500"`
	Language string `json:"language" binding:"required" validate:"oneof=en am"`
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
