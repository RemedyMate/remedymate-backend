package dto

// ProfileResponseDTO represents the user profile response
type ProfileResponseDTO struct {
	ID               string          `json:"id"`
	Username         string          `json:"username"`
	Email            string          `json:"email"`
	PersonalInfo     PersonalInfoDTO `json:"personalInfo"`
	HealthConditions string          `json:"healthConditions,omitempty"`
	IsVerified       bool            `json:"isVerified"`
	IsProfileFull    bool            `json:"isProfileFull"`
	IsActive         bool            `json:"isActive"`
	CreatedAt        string          `json:"createdAt"`
	UpdatedAt        string          `json:"updatedAt"`
	LastLogin        string          `json:"lastLogin"`
}

// UpdateProfileDTO represents the update profile request
type UpdateProfileDTO struct {
	Username         string          `json:"username,omitempty"`
	PersonalInfo     PersonalInfoDTO `json:"personalInfo,omitempty"`
	HealthConditions string          `json:"healthConditions,omitempty"`
}

// EditProfileDTO represents the edit profile request (more comprehensive)
type EditProfileDTO struct {
	Username         string          `json:"username,omitempty"`
	PersonalInfo     PersonalInfoDTO `json:"personalInfo,omitempty"`
	HealthConditions string          `json:"healthConditions,omitempty"`
	IsProfileFull    *bool           `json:"isProfileFull,omitempty"` // pointer to allow null/false distinction
}

// DeleteProfileDTO represents the delete profile request
type DeleteProfileDTO struct {
	Password string `json:"password" binding:"required"` // Require password for security
	Reason   string `json:"reason,omitempty"`            // Optional reason for deletion
}
