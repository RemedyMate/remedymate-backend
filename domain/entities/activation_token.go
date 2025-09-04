package entities

import "time"

type ActivationToken struct {
	ID        string     `bson:"_id,omitempty"`
	Token     string     `bson:"token"`
	UserID    string     `bson:"userId"`
	Email     string     `bson:"email"`
	ExpiresAt time.Time  `bson:"expiresAt"`
	CreatedAt time.Time  `bson:"createdAt"`
	UsedAt    *time.Time `bson:"usedAt,omitempty"`
}
