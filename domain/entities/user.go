package entities

import "time"

// Role represents the role of a user in the system.
type Role string

const (
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "superadmin"
)

// PersonalInfo represents a user's personal details.
type PersonalInfo struct {
	FirstName         *string `bson:"firstName"`                   // Optional
	LastName          *string `bson:"lastName"`                    // Optional
	Age               *int    `bson:"age"`                         // Optional
	Gender            *string `bson:"gender"`                      // Optional
	ProfilePictureURL *string `bson:"profilePictureUrl,omitempty"` // Optional
}

// User represents a system user with authentication and profile details.
type User struct {
	ID           string        `bson:"_id,omitempty"`          // MongoDB ObjectID
	Username     string        `bson:"username"`               // Required
	Email        string        `bson:"email"`                  // Required
	Password     string        `bson:"-"`                      // Transient, not stored in DB
	PasswordHash string        `bson:"passwordHash"`           // Required
	PersonalInfo *PersonalInfo `bson:"personalInfo,omitempty"` // Optional
	Role         Role          `bson:"role"`                   // Required, e.g., "admin" or "superadmin"
	CreatedBy    *string       `bson:"createdBy"`              // Optional, ID of the creator
	UpdatedBy    *string       `bson:"updatedBy"`              // Optional, ID of the last updater
	CreatedAt    time.Time     `bson:"createdAt"`              // Required, auto-set during creation
	UpdatedAt    time.Time     `bson:"updatedAt"`              // Required, auto-set during updates
	LastLogin    time.Time     `bson:"lastLogin"`              // Required, updated on login
}

// UserStatus represents the status of a user account.
type UserStatus struct {
	ID            string `bson:"_id,omitempty"` // MongoDB ObjectID
	UserID        string `bson:"userId"`        // Required, references User.ID
	IsActive      bool   `bson:"isActive"`      // Required, default: true
	IsProfileFull bool   `bson:"isProfileFull"` // Required, default: false
	IsVerified    bool   `bson:"isVerified"`    // Required, default: false
}
