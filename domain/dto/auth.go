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
	Username     string          `json:"username" binding:"required"`
	Email        string          `json:"email" binding:"required,email"`
	Password     string          `json:"password" binding:"required,min=6"`
	Role         entities.Role   `json:"role" binding:"required,oneof=admin superadmin"`
	PersonalInfo PersonalInfoDTO `json:"personalInfo"`
}

// REGISTER: response
type RegisterResponseDTO struct {
	Message string `json:"message" example:"User registered successfully"`
	// ID           string           `json:"id" example:"64f3b8e2d5c123456789abcd"`
	// Username     string           `json:"username" example:"john_admin"`
	// Email        string           `json:"email" example:"john.admin@example.com"`
	// Role         string           `json:"role" example:"admin"`
	// PersonalInfo *PersonalInfoDTO `json:"personalInfo,omitempty"`
	// CreatedBy    string           `json:"createdBy" example:"64f3a9d1e4b123456789abcd"`
	// CreatedAt    time.Time        `json:"createdAt" example:"2025-09-01T11:30:00Z"`
	// UpdatedAt    time.Time        `json:"updatedAt" example:"2025-09-01T11:30:00Z"`
	// LastLogin    *time.Time       `json:"lastLogin,omitempty" example:"null"`
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
	// TODO: Decide on which user data to send.
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
