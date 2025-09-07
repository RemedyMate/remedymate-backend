package dto

import (
	"remedymate-backend/domain/entities"
)

// REGISTER: request
// TODO: Add validation tags as needed and try to have comprehensive request DTOs
type PersonalInfoDTO struct {
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	Age               int    `json:"age"`
	Gender            string `json:"gender"`
	ProfilePictureURL string `json:"profilePictureUrl,omitempty"` // Optional
}

type RegisterDTO struct {
	Username       string          `json:"username,omitempty"`
	Email          string          `json:"email" binding:"required,email"`
	Password       string          `json:"password,omitempty"`
	Role           entities.Role   `json:"role" binding:"required,oneof=admin superadmin"`
	PersonalInfo   PersonalInfoDTO `json:"personalInfo"`
	FrontendDomain string          `json:"frontendDomain" binding:"required,url"`
}

// REGISTER: response
type RegisterResponseDTO struct {
	Message string `json:"message" example:"User registered successfully"`
}

// LOGIN: request
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LOGIN: response
type LoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	// User related data
	UserID   string        `json:"user_id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Role     entities.Role `json:"role"`
	Message  string        `json:"message,omitempty"`
}

// REFRESH: request
type RefreshDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// REFRESH: response
type RefreshResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Activate: request
type ActivateDTO struct {
	Email string `json:"email" binding:"required,email"`
}

// Activate: response
type ActivateResponseDTO struct {
	Message string `json:"message"`
}
