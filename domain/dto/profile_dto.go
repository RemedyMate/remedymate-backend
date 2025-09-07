package dto

// ProfileResponseDTO represents the user profile response
type ProfileResponseDTO struct {
	ID            string          `json:"id"`
	Username      string          `json:"username"`
	Email         string          `json:"email"`
	PersonalInfo  PersonalInfoDTO `json:"personalInfo"`
	IsVerified    bool            `json:"isVerified"`
	IsProfileFull bool            `json:"isProfileFull"`
	IsActive      bool            `json:"isActive"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	LastLogin     string          `json:"lastLogin"`
}

// UpdateProfileDTO represents the update profile request
type UpdateProfileDTO struct {
	Username     string          `json:"username,omitempty"`
	PersonalInfo PersonalInfoDTO `json:"personalInfo,omitempty"`
}

// EditProfileDTO removed; use UpdateProfileDTO for updates.

// DeleteProfileDTO represents the delete profile request
type DeleteProfileDTO struct {
	Password string `json:"password" binding:"required"` // Require password for security
	Reason   string `json:"reason,omitempty"`            // Optional reason for deletion
}

// UserProfilesQueryParams defines parameters for filtering and paginating user profiles
type UserProfilesQueryParams struct {
	PaginationQueryParams
	Search string `form:"search" json:"search"`   // Search by username, email, first name, or last name
	Status string `form:"status" json:"status"`   // Filter by status: "active", "inactive", "verified", "unverified", "all"
	Role   string `form:"role" json:"role"`       // Filter by role: "admin", "superadmin", "all"
	SortBy string `form:"sort_by" json:"sort_by"` // Sort by: "username", "email", "created_at", "last_login"
	Order  string `form:"order" json:"order"`     // Sort order: "asc", "desc"
}
