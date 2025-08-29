package dto

import "remedymate-backend/domain/entities"

// TriageRequest represents the request for symptom triage
type RemedyRequest struct {
	Text     string `json:"text" binding:"required" validate:"min=3,max=500"`
	Language string `json:"language" binding:"required" validate:"oneof=en am"`
}

// TriageResponse represents the response from triage
type TriageResponse struct {
	Level     entities.TriageLevel `json:"level"`
	RedFlags  []string             `json:"red_flags"`
	Message   string               `json:"message"`
	SessionID string               `json:"session_id,omitempty"` // ?
}

type RemedyResponse struct {
	SessionID string                 `json:"session_id"`
	Triage    TriageResponse         `json:"triage"`
	Content   *entities.GuidanceCard `json:"guidance_card,omitempty"`
}
