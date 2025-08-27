package entities

import "time"

type OAuthProvider struct {
	Provider string `bson:"provider"`
	ID       string `bson:"id"`
}

type PersonalInfo struct {
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
	Age       int    `bson:"age"`
	Gender    string `bson:"gender"`
}

type User struct {
	ID               string          `bson:"_id,omitempty"`
	Username         string          `bson:"username"`
	Email            string          `bson:"email"`
	Password         string          `bson:"password"`
	PersonalInfo     PersonalInfo    `bson:"personalInfo,omitempty"`
	HealthConditions string          `bson:"healthConditions,omitempty"`
	IsVerified       bool            `bson:"isVerified"`
	IsProfileFull    bool            `bson:"isProfileFull"`
	OAuthProviders   []OAuthProvider `bson:"oauthProviders,omitempty"`
	RefreshToken     string          `bson:"refreshToken,omitempty"`
	IsActive         bool            `bson:"isActive"`
	CreatedAt        time.Time       `bson:"createdAt"`
	UpdatedAt        time.Time       `bson:"updatedAt"`
	LastLogin        time.Time       `bson:"lastLogin,omitempty"`
}
